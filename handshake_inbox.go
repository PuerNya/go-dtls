package dtls13

type completedHandshake struct {
	typ      uint8
	sequence uint16
	body     []byte
}

type completedHandshakeBatch struct {
	inline [1]completedHandshake
	values []completedHandshake
	count  int
}

func completedHandshakeBatchFor(dst []completedHandshake) completedHandshakeBatch {
	return completedHandshakeBatch{values: dst, count: len(dst)}
}

func (b *completedHandshakeBatch) add(message completedHandshake) {
	if b.values != nil {
		b.values = append(b.values, message)
		b.count++
		return
	}
	if b.count == 0 {
		b.inline[0] = message
		b.count = 1
		return
	}
	b.values = make([]completedHandshake, 1, 2)
	b.values[0] = b.inline[0]
	b.values = append(b.values, message)
	b.count++
}

func (b *completedHandshakeBatch) len() int { return b.count }

func (b *completedHandshakeBatch) at(index int) completedHandshake {
	if b.values != nil {
		return b.values[index]
	}
	return b.inline[index]
}

func (b *completedHandshakeBatch) slice() []completedHandshake {
	if b.values != nil {
		return b.values
	}
	return b.inline[:b.count]
}

// handshakeInbox applies RFC 9147's next_receive_seq ordering on top of
// fragment reassembly. Future complete messages remain buffered until every
// preceding message has arrived.
type handshakeInbox struct {
	expected    uint16
	reassembler *reassembler
	ready       map[uint16]completedHandshake
	maxReady    int
	maxBytes    int
	readyBytes  int
}

func newHandshakeInbox(expected uint16, maxMessage, maxMessages, maxBytes int) *handshakeInbox {
	return &handshakeInbox{expected: expected, reassembler: newReassemblerWithLimits(maxMessage, maxMessages, maxBytes), ready: make(map[uint16]completedHandshake), maxReady: maxMessages, maxBytes: maxBytes}
}
func (i *handshakeInbox) add(fragment handshakeFragment) ([]completedHandshake, error) {
	return i.addAtEpochInto(nil, fragment, 0, false)
}

func (i *handshakeInbox) addProtected(fragment handshakeFragment, epoch uint64) ([]completedHandshake, error) {
	return i.addAtEpochInto(nil, fragment, epoch, true)
}

func (i *handshakeInbox) addInto(dst []completedHandshake, fragment handshakeFragment) ([]completedHandshake, error) {
	return i.addAtEpochInto(dst, fragment, 0, false)
}

func (i *handshakeInbox) addAtEpochInto(dst []completedHandshake, fragment handshakeFragment, epoch uint64, protected bool) ([]completedHandshake, error) {
	if dst == nil {
		dst = make([]completedHandshake, 0, 1)
	}
	batch := completedHandshakeBatchFor(dst)
	if err := i.addAtEpochBatch(&batch, fragment, epoch, protected); err != nil {
		return nil, err
	}
	return batch.values, nil
}

func (i *handshakeInbox) addBatch(dst *completedHandshakeBatch, fragment handshakeFragment) error {
	return i.addAtEpochBatch(dst, fragment, 0, false)
}

func (i *handshakeInbox) addProtectedBatch(dst *completedHandshakeBatch, fragment handshakeFragment, epoch uint64) error {
	return i.addAtEpochBatch(dst, fragment, epoch, true)
}

func (i *handshakeInbox) addAtEpochBatch(dst *completedHandshakeBatch, fragment handshakeFragment, epoch uint64, protected bool) error {
	if fragment.messageSequence < i.expected {
		return nil
	}
	if _, exists := i.ready[fragment.messageSequence]; exists {
		return nil
	}
	var body []byte
	var done bool
	var err error
	if protected {
		body, done, err = i.reassembler.addProtected(fragment, epoch)
	} else {
		body, done, err = i.reassembler.add(fragment)
	}
	if err != nil {
		return err
	}
	if !done {
		return nil
	}
	if len(i.ready) >= i.maxReady {
		return &ProtocolError{"too many completed out-of-order handshake messages"}
	}
	if i.readyBytes+len(body) > i.maxBytes {
		return &ProtocolError{"completed handshake message memory limit exceeded"}
	}
	message := completedHandshake{typ: fragment.typ, sequence: fragment.messageSequence, body: body}
	if fragment.messageSequence != i.expected {
		i.ready[fragment.messageSequence] = message
		i.readyBytes += len(body)
		return nil
	}
	dst.add(message)
	i.expected++
	for {
		message, ok := i.ready[i.expected]
		if !ok {
			break
		}
		delete(i.ready, i.expected)
		i.readyBytes -= len(message.body)
		dst.add(message)
		i.expected++
	}
	return nil
}

func (i *handshakeInbox) hasIncompleteProtected() bool {
	return i != nil && i.reassembler.hasIncompleteProtected()
}
