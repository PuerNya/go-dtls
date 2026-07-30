package dtls13

import (
	"crypto/hmac"
	"errors"
	"sync"
)

const keyScheduleSecretSlots = 6

const (
	keyScheduleRootSlot = iota
	keyScheduleClientHandshakeSlot
	keyScheduleServerHandshakeSlot
	keyScheduleClientApplicationSlot
	keyScheduleServerApplicationSlot
	keyScheduleExporterSlot
)

// keySchedule implements the TLS 1.3 key schedule as modified by RFC 9147
// section 5.10. The caller supplies transcript hashes at the precise protocol
// boundaries; keeping transcript mutation outside this type makes those
// boundaries explicit and testable.
type keySchedule struct {
	suite                    *cipherSuite
	earlySecret              []byte
	handshakeSecret          []byte
	masterSecret             []byte
	clientHandshakeTraffic   []byte
	serverHandshakeTraffic   []byte
	clientApplicationTraffic []byte
	serverApplicationTraffic []byte
	exporterMasterSecret     []byte
	resumptionMasterSecret   []byte
	secretStorage            []byte
}

func newKeySchedule(suite *cipherSuite, psk []byte) *keySchedule {
	if psk == nil {
		psk = zeroHashSecret(suite)
	}
	hashSize := suite.hash.Size()
	storage := make([]byte, keyScheduleSecretSlots*hashSize)
	earlySecret := storage[:hashSize:hashSize]
	hkdfExtractInto(suite.hash.New, psk, nil, earlySecret)
	return &keySchedule{suite: suite, earlySecret: earlySecret, secretStorage: storage}
}

func (k *keySchedule) secretWindow(slot int) []byte {
	hashSize := k.suite.hash.Size()
	start := slot * hashSize
	return k.secretStorage[start : start+hashSize : start+hashSize]
}

// earlyTrafficSecret derives the client 0-RTT traffic secret.  DTLS 1.3
// uses the TLS 1.3 key schedule, with the DTLS-specific label encoding from
// RFC 9147 section 5.9 supplied by deriveSecret.
func (k *keySchedule) earlyTrafficSecret(clientHelloTranscriptHash []byte) []byte {
	secret := k.secretWindow(keyScheduleClientHandshakeSlot)
	deriveSecretInto(k.suite, k.earlySecret, labelClientEarlyTraffic, clientHelloTranscriptHash, secret)
	return secret
}

func (k *keySchedule) deriveHandshake(sharedSecret, helloTranscriptHash []byte) error {
	if len(k.earlySecret) == 0 {
		return errors.New("dtls13: handshake secrets already derived")
	}
	handshakeSecret := k.secretWindow(keyScheduleRootSlot)
	deriveSecretInto(k.suite, k.earlySecret, labelDerived, emptyTranscriptHash(k.suite), handshakeSecret)
	hkdfExtractInto(k.suite.hash.New, sharedSecret, handshakeSecret, handshakeSecret)
	k.earlySecret = nil
	k.handshakeSecret = handshakeSecret
	expander := newHKDFExpander(k.suite.hash.New, k.handshakeSecret)
	k.clientHandshakeTraffic = k.secretWindow(keyScheduleClientHandshakeSlot)
	k.serverHandshakeTraffic = k.secretWindow(keyScheduleServerHandshakeSlot)
	expandLabelWithInto(&expander, "c hs traffic", helloTranscriptHash, k.clientHandshakeTraffic)
	expandLabelWithInto(&expander, "s hs traffic", helloTranscriptHash, k.serverHandshakeTraffic)
	return nil
}

func (k *keySchedule) finishedVerifyData(trafficSecret, transcriptHash []byte) []byte {
	return computeFinishedVerifyData(k.suite, trafficSecret, transcriptHash)
}

func computeFinishedVerifyData(suite *cipherSuite, baseSecret, transcriptHash []byte) []byte {
	verifyData := make([]byte, suite.hash.Size())
	computeFinishedVerifyDataInto(suite, baseSecret, transcriptHash, verifyData)
	return verifyData
}

func computeFinishedVerifyDataInto(suite *cipherSuite, baseSecret, transcriptHash, out []byte) {
	if len(out) != suite.hash.Size() {
		panic("dtls13: Finished output must equal Hash.length")
	}
	expandLabelHashInto(suite, baseSecret, labelFinished, nil, out)
	m := hmac.New(suite.hash.New, out)
	_, _ = m.Write(transcriptHash)
	if sum := m.Sum(out[:0]); len(sum) != len(out) {
		panic("dtls13: unexpected Finished output length")
	}
}
func (k *keySchedule) verifyFinished(trafficSecret, transcriptHash, received []byte) bool {
	return hmac.Equal(k.finishedVerifyData(trafficSecret, transcriptHash), received)
}

func (k *keySchedule) deriveApplication(serverFinishedTranscriptHash []byte) error {
	if len(k.handshakeSecret) == 0 {
		return errors.New("dtls13: handshake secrets not derived")
	}
	if len(k.masterSecret) != 0 {
		return errors.New("dtls13: application secrets already derived")
	}
	masterSecret := k.secretWindow(keyScheduleRootSlot)
	deriveSecretInto(k.suite, k.handshakeSecret, labelDerived, emptyTranscriptHash(k.suite), masterSecret)
	hkdfExtractInto(k.suite.hash.New, zeroHashSecret(k.suite), masterSecret, masterSecret)
	k.handshakeSecret = nil
	k.masterSecret = masterSecret
	expander := newHKDFExpander(k.suite.hash.New, k.masterSecret)
	k.clientApplicationTraffic = k.secretWindow(keyScheduleClientApplicationSlot)
	k.serverApplicationTraffic = k.secretWindow(keyScheduleServerApplicationSlot)
	k.exporterMasterSecret = k.secretWindow(keyScheduleExporterSlot)
	expandLabelWithInto(&expander, "c ap traffic", serverFinishedTranscriptHash, k.clientApplicationTraffic)
	expandLabelWithInto(&expander, "s ap traffic", serverFinishedTranscriptHash, k.serverApplicationTraffic)
	expandLabelWithInto(&expander, "exp master", serverFinishedTranscriptHash, k.exporterMasterSecret)
	return nil
}

type exporterState struct {
	mu          sync.Mutex
	suite       *cipherSuite
	secret      []byte
	externalPSK *externalPSKSelection
}

func newExporter(suite *cipherSuite, exporterMasterSecret []byte) *exporterState {
	return &exporterState{suite: suite, secret: append([]byte(nil), exporterMasterSecret...)}
}

func (e *exporterState) export(label string, context []byte, length int) ([]byte, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.suite == nil || len(e.secret) != e.suite.hash.Size() {
		return nil, errors.New("dtls13: exporter is unavailable")
	}
	if len(label) > 255-len("dtls13") || length < 0 || length > 255*e.suite.hash.Size() {
		return nil, errors.New("dtls13: invalid exporter label or length")
	}
	if length == 0 {
		return []byte{}, nil
	}
	hashSize := e.suite.hash.Size()
	scratch := make([]byte, 2*hashSize)
	labelSecret := scratch[:hashSize:hashSize]
	contextHash := scratch[hashSize : 2*hashSize : 2*hashSize]
	expandLabelHashStringInto(e.suite, e.secret, label, emptyTranscriptHash(e.suite), labelSecret)
	h := e.suite.hash.New()
	_, _ = h.Write(context)
	if sum := h.Sum(contextHash[:0]); len(sum) != len(contextHash) {
		panic("dtls13: unexpected exporter context hash length")
	}
	out := make([]byte, length)
	expandLabelWithScratchInto(e.suite, labelSecret, "exporter", contextHash, out, scratch)
	return out, nil
}

func (e *exporterState) clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i := range e.secret {
		e.secret[i] = 0
	}
	e.secret = nil
	e.suite = nil
}

func (k *keySchedule) deriveResumption(clientFinishedTranscriptHash []byte) error {
	if len(k.masterSecret) == 0 {
		return errors.New("dtls13: master secret not derived")
	}
	if len(k.resumptionMasterSecret) != 0 {
		return errors.New("dtls13: resumption secret already derived")
	}
	resumptionMasterSecret := k.secretWindow(keyScheduleRootSlot)
	deriveSecretInto(k.suite, k.masterSecret, labelResumptionMaster, clientFinishedTranscriptHash, resumptionMasterSecret)
	k.masterSecret = nil
	k.resumptionMasterSecret = resumptionMasterSecret
	return nil
}

func (k *keySchedule) resumptionPSK(nonce []byte) ([]byte, error) {
	if len(k.resumptionMasterSecret) == 0 {
		return nil, errors.New("dtls13: resumption secret not derived")
	}
	return deriveResumptionPSK(k.suite, k.resumptionMasterSecret, nonce), nil
}

func deriveResumptionPSK(suite *cipherSuite, resumptionMasterSecret, nonce []byte) []byte {
	return deriveSecret(suite, resumptionMasterSecret, labelResumption, nonce)
}
