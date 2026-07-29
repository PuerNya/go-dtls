package dtls13

import "errors"

const maxSendingEpoch uint64 = 1<<48 - 1

// sendingTraffic manages the ACK-gated sending side of RFC 9147 section 8.
// A KeyUpdate is protected with the current epoch; the next epoch is not used
// until an ACK covers that KeyUpdate record.
type sendingTraffic struct {
	suite             *cipherSuite
	secret            []byte
	cipher            *recordCipher
	nextSecret        []byte
	pendingCipher     *recordCipher
	update            keyUpdateState
	messageSequence   uint16
	sequenceExhausted bool
	replaySize        int
	pendingFragment   []byte
	connectionID      []byte
	hasConnectionID   bool
}

func (s *sendingTraffic) clearSecrets() {
	clear(s.secret)
	clear(s.nextSecret)
	s.secret = nil
	s.nextSecret = nil
}

func (s *sendingTraffic) canAllocateMessageSequences(count uint32) error {
	if count == 0 {
		return nil
	}
	if s.sequenceExhausted || uint32(s.messageSequence)+count > 1<<16 {
		return errors.New("dtls13: handshake message sequence exhausted")
	}
	return nil
}

func (s *sendingTraffic) commitMessageSequences(count uint32) {
	next := uint32(s.messageSequence) + count
	if next == 1<<16 {
		s.sequenceExhausted = true
		return
	}
	s.messageSequence = uint16(next)
}

func (s *sendingTraffic) canBeginKeyUpdate() bool {
	return s != nil && s.update.canUseNewKeys() && s.cipher.epoch < maxSendingEpoch
}

func newSendingTraffic(suite *cipherSuite, secret []byte, epoch uint64, messageSequence uint16, replaySize int) (*sendingTraffic, error) {
	cipher, err := newRecordCipher(suite, secret, epoch, replaySize)
	if err != nil {
		return nil, err
	}
	return newSendingTrafficWithCipher(suite, append([]byte(nil), secret...), cipher, messageSequence, replaySize), nil
}

func newSendingTrafficWithOwnedSecret(suite *cipherSuite, secret []byte, epoch uint64, messageSequence uint16, replaySize int) (*sendingTraffic, error) {
	cipher, err := newRecordCipher(suite, secret, epoch, replaySize)
	if err != nil {
		return nil, err
	}
	return newSendingTrafficWithCipher(suite, secret, cipher, messageSequence, replaySize), nil
}

func newSendingTrafficWithCipher(suite *cipherSuite, secret []byte, cipher *recordCipher, messageSequence uint16, replaySize int) *sendingTraffic {
	return &sendingTraffic{suite: suite, secret: secret, cipher: cipher, messageSequence: messageSequence, replaySize: replaySize}
}
func (s *sendingTraffic) setConnectionID(connectionID []byte) error {
	if err := s.cipher.setConnectionID(connectionID); err != nil {
		return err
	}
	s.connectionID = append([]byte(nil), connectionID...)
	s.hasConnectionID = true
	if s.pendingCipher != nil {
		return s.pendingCipher.setConnectionID(connectionID)
	}
	return nil
}
func (s *sendingTraffic) beginKeyUpdate(requestUpdate bool) (wire []byte, number recordNumber, err error) {
	if !s.update.canUseNewKeys() {
		return nil, recordNumber{}, errors.New("dtls13: KeyUpdate is already pending")
	}
	if s.cipher.epoch >= maxSendingEpoch {
		return nil, recordNumber{}, errors.New("dtls13: sending epoch limit reached")
	}
	if err = s.canAllocateMessageSequences(1); err != nil {
		return nil, recordNumber{}, err
	}
	var body [1]byte
	if requestUpdate {
		body[0] = 1
	}
	fragmentLen := handshakeHeaderLen + len(body)
	if cap(s.pendingFragment) < fragmentLen {
		s.pendingFragment = make([]byte, fragmentLen)
	} else {
		s.pendingFragment = s.pendingFragment[:fragmentLen]
	}
	fragment := handshakeFragment{typ: handshakeTypeKeyUpdate, messageSequence: s.messageSequence, length: uint32(len(body)), body: body[:]}
	if err = marshalHandshakeFragmentInto(s.pendingFragment, fragment); err != nil {
		return nil, recordNumber{}, err
	}
	sequence := s.cipher.nextSequence
	wire, err = s.cipher.seal(recordTypeHandshake, s.pendingFragment)
	if err != nil {
		return nil, recordNumber{}, err
	}
	number = recordNumber{epoch: s.cipher.epoch, sequence: sequence}
	if err = s.update.begin(number); err != nil {
		return nil, recordNumber{}, err
	}
	if len(s.nextSecret) != s.suite.hash.Size() {
		s.nextSecret = make([]byte, s.suite.hash.Size())
	}
	scheduleTrafficUpdateInto(s.suite, s.secret, s.nextSecret)
	s.pendingCipher, err = newRecordCipher(s.suite, s.nextSecret, s.cipher.epoch+1, s.replaySize)
	if err != nil {
		return nil, recordNumber{}, err
	}
	s.pendingCipher.setPlaintextLimit(s.cipher.plaintextLimit)
	if s.hasConnectionID {
		if err = s.pendingCipher.setConnectionID(s.connectionID); err != nil {
			return nil, recordNumber{}, err
		}
	}
	s.commitMessageSequences(1)
	return wire, number, nil
}
func (s *sendingTraffic) retransmitKeyUpdate() (wire []byte, number recordNumber, err error) {
	if s.update.canUseNewKeys() || len(s.pendingFragment) == 0 {
		return nil, recordNumber{}, errors.New("dtls13: no KeyUpdate is pending")
	}
	sequence := s.cipher.nextSequence
	wire, err = s.cipher.seal(recordTypeHandshake, s.pendingFragment)
	if err != nil {
		return nil, recordNumber{}, err
	}
	number = recordNumber{epoch: s.cipher.epoch, sequence: sequence}
	if err = s.update.addRetransmission(number); err != nil {
		return nil, recordNumber{}, err
	}
	return wire, number, nil
}
func (s *sendingTraffic) processACK(numbers []recordNumber) bool {
	if !s.update.ack(numbers) {
		return false
	}
	oldSecret := s.secret
	s.secret = s.nextSecret
	clear(oldSecret)
	s.nextSecret = oldSecret
	s.cipher = s.pendingCipher
	s.pendingCipher = nil
	s.pendingFragment = s.pendingFragment[:0]
	return true
}
func scheduleTrafficUpdate(suite *cipherSuite, secret []byte) []byte {
	out := make([]byte, suite.hash.Size())
	scheduleTrafficUpdateInto(suite, secret, out)
	return out
}

func scheduleTrafficUpdateInto(suite *cipherSuite, secret, out []byte) {
	deriveSecretInto(suite, secret, labelTrafficUpdate, nil, out)
}

type receivingTraffic struct {
	suite              *cipherSuite
	secret             []byte
	epochs             *epochSet
	current            uint64
	replaySize         int
	lastUpdateSequence uint16
	hasUpdateSequence  bool
	connectionID       []byte
	hasConnectionID    bool
	plaintextLimit     uint16
	acceptedCIDs       [][]byte
	secrets            map[uint64][]byte
}

func (r *receivingTraffic) clearSecrets() {
	for epoch, secret := range r.secrets {
		clear(secret)
		delete(r.secrets, epoch)
	}
	clear(r.secret)
	r.secret = nil
	r.secrets = nil
}

func newReceivingTraffic(suite *cipherSuite, secret []byte, epoch uint64, replaySize int) (*receivingTraffic, error) {
	cipher, err := newRecordCipher(suite, secret, epoch, replaySize)
	if err != nil {
		return nil, err
	}
	ownedSecret := make([]byte, len(secret))
	copy(ownedSecret, secret)
	return newReceivingTrafficWithCipher(suite, ownedSecret, cipher, epoch, replaySize)
}

func newReceivingTrafficWithOwnedSecret(suite *cipherSuite, secret []byte, epoch uint64, replaySize int) (*receivingTraffic, error) {
	cipher, err := newRecordCipher(suite, secret, epoch, replaySize)
	if err != nil {
		return nil, err
	}
	return newReceivingTrafficWithCipher(suite, secret, cipher, epoch, replaySize)
}

func newReceivingTrafficWithCipher(suite *cipherSuite, secret []byte, cipher *recordCipher, epoch uint64, replaySize int) (*receivingTraffic, error) {
	epochs := newEpochSet()
	if err := epochs.install(cipher); err != nil {
		return nil, err
	}
	if err := epochs.setCurrent(epoch); err != nil {
		return nil, err
	}
	return &receivingTraffic{
		suite: suite, secret: secret, epochs: epochs, current: epoch, replaySize: replaySize,
		plaintextLimit: cipher.plaintextLimit, secrets: map[uint64][]byte{epoch: secret},
	}, nil
}

func (r *receivingTraffic) setPlaintextLimit(limit uint16) {
	r.plaintextLimit = limit
	r.epochs.mu.Lock()
	defer r.epochs.mu.Unlock()
	for _, cipher := range r.epochs.ciphers {
		cipher.setPlaintextLimit(limit)
	}
}
func (r *receivingTraffic) setConnectionID(connectionID []byte) error {
	r.connectionID = append([]byte(nil), connectionID...)
	r.hasConnectionID = true
	r.acceptedCIDs = [][]byte{append([]byte{}, connectionID...)}
	r.epochs.mu.Lock()
	defer r.epochs.mu.Unlock()
	for _, cipher := range r.epochs.ciphers {
		if err := cipher.setConnectionID(connectionID); err != nil {
			return err
		}
	}
	return nil
}

func (r *receivingTraffic) addConnectionIDs(connectionIDs [][]byte) error {
	if !r.hasConnectionID {
		return &ProtocolError{"connection ID was not negotiated"}
	}
	r.epochs.mu.Lock()
	defer r.epochs.mu.Unlock()
	for _, cipher := range r.epochs.ciphers {
		if err := cipher.addAcceptedConnectionIDs(connectionIDs); err != nil {
			return err
		}
	}
	for _, cid := range connectionIDs {
		found := false
		for _, existing := range r.acceptedCIDs {
			if equalBytes(existing, cid) {
				found = true
				break
			}
		}
		if !found {
			r.acceptedCIDs = append(r.acceptedCIDs, append([]byte{}, cid...))
		}
	}
	return nil
}

func (r *receivingTraffic) restoreConnectionIDs(connectionIDs [][]byte) {
	r.epochs.mu.Lock()
	defer r.epochs.mu.Unlock()
	r.acceptedCIDs = cloneConnectionIDs(connectionIDs)
	for _, cipher := range r.epochs.ciphers {
		cipher.acceptedCIDs = cloneConnectionIDs(connectionIDs)
	}
}

func cloneConnectionIDs(connectionIDs [][]byte) [][]byte {
	cloned := make([][]byte, len(connectionIDs))
	for i, cid := range connectionIDs {
		cloned[i] = append([]byte{}, cid...)
	}
	return cloned
}
func (r *receivingTraffic) processKeyUpdate(sequence uint16, body []byte) (keyUpdateMessage, bool, error) {
	message, err := parseKeyUpdate(body)
	if err != nil {
		return keyUpdateMessage{}, false, err
	}
	if r.hasUpdateSequence {
		if sequence <= r.lastUpdateSequence {
			return message, false, nil
		}
	}
	nextSecret := scheduleTrafficUpdate(r.suite, r.secret)
	cipher, err := newRecordCipher(r.suite, nextSecret, r.current+1, r.replaySize)
	if err != nil {
		return keyUpdateMessage{}, false, err
	}
	cipher.setPlaintextLimit(r.plaintextLimit)
	if r.hasConnectionID {
		if err = cipher.setConnectionID(r.connectionID); err != nil {
			return keyUpdateMessage{}, false, err
		}
		if err = cipher.addAcceptedConnectionIDs(r.acceptedCIDs[1:]); err != nil {
			return keyUpdateMessage{}, false, err
		}
	}
	if err = r.epochs.install(cipher); err != nil {
		return keyUpdateMessage{}, false, err
	}
	if err = r.epochs.setCurrent(r.current + 1); err != nil {
		return keyUpdateMessage{}, false, err
	}
	r.secret = nextSecret
	r.current++
	r.secrets[r.current] = nextSecret
	if r.current > 1 {
		minimum := r.current - 1
		r.epochs.discardBefore(minimum)
		for epoch, secret := range r.secrets {
			if epoch < minimum {
				clear(secret)
				delete(r.secrets, epoch)
			}
		}
	}
	r.lastUpdateSequence = sequence
	r.hasUpdateSequence = true
	return message, true, nil
}

func (r *receivingTraffic) secretForEpoch(epoch uint64) ([]byte, bool) {
	secret, ok := r.secrets[epoch]
	return secret, ok
}
