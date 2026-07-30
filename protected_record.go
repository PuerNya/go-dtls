package dtls13

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"
)

// authenticatedRecordError reports a semantic error discovered only after
// AEAD authentication succeeded. Callers may safely send an alert for these
// errors without turning unauthenticated UDP input into an amplification path.
type authenticatedRecordError struct {
	description uint8
	err         error
}

func (e *authenticatedRecordError) Error() string { return e.err.Error() }
func (e *authenticatedRecordError) Unwrap() error { return e.err }

func authenticatedRecordAlert(description uint8, err error) error {
	return &authenticatedRecordError{description: description, err: err}
}

func protectedRecordReceiveError(err error) error {
	if errors.Is(err, errAEADAuthenticationFailureLimit) {
		return err
	}
	for current := err; current != nil; {
		if authenticatedErr, ok := current.(*authenticatedRecordError); ok {
			description := authenticatedErr.description
			if description == 0 {
				description = alertUnexpectedMessage
			}
			return alertError(description, authenticatedErr)
		}
		unwrapper, ok := current.(interface{ Unwrap() error })
		if !ok {
			break
		}
		current = unwrapper.Unwrap()
	}
	return nil
}

const (
	unifiedHeaderFixed           = 0x20
	unifiedHeaderCID             = 0x10
	unifiedHeaderSequence16      = 0x08
	unifiedHeaderLength          = 0x04
	unifiedHeaderEpochMask       = 0x03
	unifiedHeaderLen16           = 5
	maxCiphertextLen             = (1 << 14) + 256
	sequenceProtectionSampleSize = 16
	maxTrafficKeyLen             = 32
	maxTrafficIVLen              = 16
)

var errAEADAuthenticationFailureLimit = errors.New("dtls13: AEAD authentication failure limit reached")

// recordCipher protects records for one epoch using the RFC 9147 unified
// header, a two-byte truncated sequence number, and an explicit length.
type recordCipher struct {
	epoch            uint64
	aead             cipher.AEAD
	iv               [maxTrafficIVLen]byte
	ivLen            int
	sequenceMask     sequenceMasker
	nextSequence     uint64
	lastOpened       uint64
	replay           replayWindow
	authFailures     uint64
	connectionID     []byte
	hasConnectionID  bool
	plaintextLimit   uint16
	acceptedCIDs     [][]byte
	lastConnectionID []byte
	recordLimit      uint64
	authFailureLimit uint64
	nonceScratch     [16]byte
	headerScratch    [unifiedHeaderLen16 + 255]byte
}

func (c *recordCipher) setConnectionID(connectionID []byte) error {
	if len(connectionID) > 255 {
		return &ConfigError{"connection ID exceeds 255 bytes"}
	}
	c.connectionID = append([]byte{}, connectionID...)
	c.hasConnectionID = true
	c.acceptedCIDs = [][]byte{append([]byte{}, connectionID...)}
	return nil
}

func (c *recordCipher) addAcceptedConnectionIDs(connectionIDs [][]byte) error {
	if !c.hasConnectionID {
		return &ProtocolError{"connection ID was not negotiated"}
	}
	merged, err := mergeConnectionIDs(c.acceptedCIDs, connectionIDs)
	if err != nil {
		return err
	}
	c.acceptedCIDs = merged
	return nil
}

func mergeConnectionIDs(existing, additions [][]byte) ([][]byte, error) {
	merged := make([][]byte, len(existing))
	for i, cid := range existing {
		merged[i] = append([]byte{}, cid...)
	}
	for _, cid := range additions {
		if len(cid) > 255 {
			return nil, &ConfigError{"connection ID exceeds 255 bytes"}
		}
		duplicate := false
		for _, current := range merged {
			if equalBytes(current, cid) {
				duplicate = true
				break
			}
			shorter := len(current)
			if len(cid) < shorter {
				shorter = len(cid)
			}
			if equalBytes(current[:shorter], cid[:shorter]) {
				return nil, &ConfigError{"connection IDs must be prefix-free"}
			}
		}
		if !duplicate {
			merged = append(merged, append([]byte{}, cid...))
		}
	}
	return merged, nil
}

func (c *recordCipher) matchConnectionID(datagram []byte) ([]byte, bool) {
	for _, cid := range c.acceptedCIDs {
		if len(datagram) >= len(cid) && equalBytes(datagram[:len(cid)], cid) {
			return cid, true
		}
	}
	return nil, false
}

func (c *recordCipher) headerLen16() int { return unifiedHeaderLen16 + len(c.connectionID) }

func (c *recordCipher) shouldRequestKeyUpdateForAuthFailures() bool {
	if c.authFailureLimit == 0 {
		return false
	}
	margin := uint64(1024)
	if c.authFailureLimit/4 < margin {
		margin = c.authFailureLimit / 4
	}
	if margin < 1 {
		margin = 1
	}
	return c.authFailures >= c.authFailureLimit-margin
}

func newRecordCipher(suite *cipherSuite, trafficSecret []byte, epoch uint64, replaySize int) (*recordCipher, error) {
	if suite.keyLen > maxTrafficKeyLen || suite.ivLen > maxTrafficIVLen {
		return nil, &ConfigError{"cipher suite key material exceeds internal storage"}
	}
	var material [2*maxTrafficKeyLen + maxTrafficIVLen]byte
	key := material[:suite.keyLen]
	iv := material[maxTrafficKeyLen : maxTrafficKeyLen+suite.ivLen]
	sn := material[maxTrafficKeyLen+maxTrafficIVLen : maxTrafficKeyLen+maxTrafficIVLen+suite.keyLen]
	deriveTrafficKeysInto(suite, trafficSecret, key, iv, sn)
	defer clear(material[:])

	aead, err := suite.newAEAD(key)
	if err != nil {
		return nil, err
	}
	mask, err := suite.newSequenceMask(sn)
	if err != nil {
		return nil, err
	}
	record := &recordCipher{epoch: epoch, aead: aead, ivLen: suite.ivLen, sequenceMask: mask, replay: newReplayWindow(replaySize), recordLimit: suite.recordLimit, authFailureLimit: suite.authFailureLimit}
	copy(record.iv[:record.ivLen], iv)
	return record, nil
}

func (c *recordCipher) setPlaintextLimit(limit uint16) {
	limit = effectiveRecordSizeLimit(limit)
	if limit == defaultRecordSizeLimit {
		limit = 0
	}
	c.plaintextLimit = limit
}

func (c *recordCipher) maxContent() int {
	if c.plaintextLimit == 0 {
		return maxRecordContent
	}
	return int(c.plaintextLimit) - 1
}

func (c *recordCipher) nonce(sequence uint64) []byte {
	return c.nonceInto(sequence, make([]byte, c.ivLen))
}

func (c *recordCipher) nonceInto(sequence uint64, nonce []byte) []byte {
	nonce = nonce[:c.ivLen]
	copy(nonce, c.iv[:c.ivLen])
	for i := range 8 {
		nonce[len(nonce)-1-i] ^= byte(sequence >> (8 * i))
	}
	return nonce
}

func (c *recordCipher) seal(contentType uint8, content []byte) ([]byte, error) {
	return c.sealWithHeader(contentType, content, true, true)
}

func (c *recordCipher) sealWithHeader(contentType uint8, content []byte, sequence16, lengthPresent bool) ([]byte, error) {
	return c.sealWithHeaderBuilder(contentType, len(content), sequence16, lengthPresent, func(dst []byte) {
		copy(dst, content)
	})
}

func (c *recordCipher) sealHandshakeFragmentInto(dst []byte, fragment handshakeFragment) ([]byte, error) {
	if err := validateHandshakeFragment(fragment); err != nil {
		return nil, err
	}
	return c.sealWithHeaderBuilderInto(dst, recordTypeHandshake, handshakeHeaderLen+len(fragment.body), true, true, func(plain []byte) {
		putHandshakeFragment(plain, fragment)
	})
}

func (c *recordCipher) sealWithHeaderBuilder(contentType uint8, contentLen int, sequence16, lengthPresent bool, fill func([]byte)) ([]byte, error) {
	return c.sealWithHeaderBuilderInto(nil, contentType, contentLen, sequence16, lengthPresent, fill)
}

func (c *recordCipher) sealWithHeaderBuilderInto(dst []byte, contentType uint8, contentLen int, sequence16, lengthPresent bool, fill func([]byte)) ([]byte, error) {
	if c.epoch > maxSendingEpoch {
		return nil, &ProtocolError{"sending epoch exceeds 2^48-1"}
	}
	if c.nextSequence >= 1<<48 {
		return nil, &ProtocolError{"record sequence number exhausted"}
	}
	if c.nextSequence >= c.recordLimit {
		return nil, &ProtocolError{"AEAD record usage limit reached; KeyUpdate required"}
	}
	if !validInnerContentType(contentType) {
		return nil, &ProtocolError{"invalid protected content type"}
	}
	if contentLen == 0 && (contentType == recordTypeHandshake || contentType == recordTypeAlert) {
		return nil, &ProtocolError{"zero-length protected Handshake or Alert content"}
	}
	if contentLen < 0 || contentLen > maxRecordContent || (c.plaintextLimit != 0 && contentLen >= int(c.plaintextLimit)) {
		return nil, &ProtocolError{"protected record content exceeds record_size_limit"}
	}
	sequence := c.nextSequence
	plainLen := contentLen + 1
	ciphertextLen := plainLen + c.aead.Overhead()
	if ciphertextLen > maxCiphertextLen {
		return nil, &ProtocolError{"protected record is too large"}
	}
	sequenceBytes := 1
	if sequence16 {
		sequenceBytes = 2
	}
	cidLength := len(c.connectionID)
	sequenceOffset := 1 + cidLength
	headerLen := sequenceOffset + sequenceBytes
	if lengthPresent {
		headerLen += 2
	}
	packetLen := headerLen + ciphertextLen
	var packet []byte
	if cap(dst) < packetLen {
		packet = make([]byte, headerLen+plainLen, packetLen)
	} else {
		packet = dst[: headerLen+plainLen : packetLen]
	}
	header := packet[:headerLen]
	header[0] = unifiedHeaderFixed | byte(c.epoch&unifiedHeaderEpochMask)
	if c.hasConnectionID {
		header[0] |= unifiedHeaderCID
		copy(header[1:sequenceOffset], c.connectionID)
	}
	if sequence16 {
		header[0] |= unifiedHeaderSequence16
		binary.BigEndian.PutUint16(header[sequenceOffset:sequenceOffset+2], uint16(sequence))
	} else {
		header[sequenceOffset] = byte(sequence)
	}
	if lengthPresent {
		header[0] |= unifiedHeaderLength
		binary.BigEndian.PutUint16(header[sequenceOffset+sequenceBytes:], uint16(ciphertextLen))
	}
	plain := packet[headerLen : headerLen+plainLen]
	fill(plain[:contentLen])
	plain[contentLen] = contentType
	packet = c.aead.Seal(header, c.nonceInto(sequence, c.nonceScratch[:]), plain, header)
	ciphertext := packet[headerLen:]
	if len(ciphertext) < sequenceProtectionSampleSize {
		return nil, errors.New("dtls13: ciphertext too short for sequence protection")
	}
	mask, err := c.sequenceMask.mask(ciphertext[:sequenceProtectionSampleSize])
	if err != nil {
		return nil, err
	}
	for i := 0; i < sequenceBytes; i++ {
		header[sequenceOffset+i] ^= mask[i]
	}
	c.nextSequence++
	return packet, nil
}

// open authenticates the first protected record in datagram and returns its
// content, inner type, and number of consumed bytes. Replay state is committed
// only after successful AEAD authentication.
func (c *recordCipher) open(datagram []byte) (content []byte, contentType uint8, consumed int, err error) {
	return c.openRecord(datagram, false)
}

// openInPlace authenticates and decrypts into datagram. The returned content
// is valid only while datagram remains owned by the caller.
func (c *recordCipher) openInPlace(datagram []byte) (content []byte, contentType uint8, consumed int, err error) {
	return c.openRecord(datagram, true)
}

func recordCipherMatchesUnifiedEpoch(cipher *recordCipher, first byte) bool {
	return cipher != nil && first&0xe0 == unifiedHeaderFixed && uint64(first&unifiedHeaderEpochMask) == cipher.epoch&unifiedHeaderEpochMask
}

func (c *recordCipher) openRecord(datagram []byte, inPlace bool) (content []byte, contentType uint8, consumed int, err error) {
	if len(datagram) < 2 {
		return nil, 0, 0, &ProtocolError{"truncated unified record header"}
	}
	first := datagram[0]
	if first&0xe0 != unifiedHeaderFixed {
		return nil, 0, 0, &ProtocolError{"unsupported unified record header"}
	}
	if !recordCipherMatchesUnifiedEpoch(c, first) {
		return nil, 0, 0, &ProtocolError{"record epoch does not match cipher epoch"}
	}
	sequenceBytes := 1
	if first&unifiedHeaderSequence16 != 0 {
		sequenceBytes = 2
	}
	hasCID := first&unifiedHeaderCID != 0
	if hasCID != c.hasConnectionID {
		return nil, 0, 0, &ProtocolError{"record connection ID presence does not match epoch"}
	}
	cidLength := 0
	var matchedCID []byte
	if hasCID {
		var ok bool
		matchedCID, ok = c.matchConnectionID(datagram[1:])
		if !ok {
			return nil, 0, 0, &ProtocolError{"record connection ID does not match"}
		}
		cidLength = len(matchedCID)
	}
	sequenceOffset := 1 + cidLength
	headerLen := sequenceOffset + sequenceBytes
	if first&unifiedHeaderLength != 0 {
		headerLen += 2
	}
	if len(datagram) < headerLen {
		return nil, 0, 0, &ProtocolError{"truncated unified record header"}
	}
	n := len(datagram) - headerLen
	if first&unifiedHeaderLength != 0 {
		n = int(binary.BigEndian.Uint16(datagram[sequenceOffset+sequenceBytes : headerLen]))
	}
	if n < c.aead.Overhead()+1 || n > maxCiphertextLen || len(datagram) < headerLen+n {
		return nil, 0, 0, &ProtocolError{"invalid protected record length"}
	}
	ciphertext := datagram[headerLen : headerLen+n]
	if len(ciphertext) < sequenceProtectionSampleSize {
		return nil, 0, 0, &ProtocolError{"ciphertext too short for sequence protection"}
	}
	header := c.headerScratch[:headerLen]
	copy(header, datagram[:headerLen])
	mask, err := c.sequenceMask.mask(ciphertext[:sequenceProtectionSampleSize])
	if err != nil {
		return nil, 0, 0, err
	}
	for i := 0; i < sequenceBytes; i++ {
		header[sequenceOffset+i] ^= mask[i]
	}
	var truncated uint64
	if sequenceBytes == 2 {
		truncated = uint64(binary.BigEndian.Uint16(header[sequenceOffset : sequenceOffset+2]))
	} else {
		truncated = uint64(header[sequenceOffset])
	}
	sequence := reconstructSequence(c.replay.nextExpected(), truncated, uint(sequenceBytes*8))
	if sequence >= 1<<48 {
		return nil, 0, 0, &ProtocolError{"invalid protected record sequence number"}
	}
	var destination []byte
	if inPlace {
		destination = ciphertext[:0]
	}
	plain, openErr := c.aead.Open(destination, c.nonceInto(sequence, c.nonceScratch[:]), ciphertext, header)
	if openErr != nil {
		if c.authFailures < ^uint64(0) {
			c.authFailures++
		}
		if c.authFailures >= c.authFailureLimit {
			return nil, 0, 0, errAEADAuthenticationFailureLimit
		}
		return nil, 0, 0, &ProtocolError{"protected record authentication failed"}
	}
	if len(plain) > maxRecordContent+1 || (c.plaintextLimit != 0 && len(plain) > int(c.plaintextLimit)) {
		return nil, 0, 0, authenticatedRecordAlert(alertRecordOverflow, &ProtocolError{"DTLSInnerPlaintext exceeds record_size_limit"})
	}
	i := len(plain) - 1
	for i >= 0 && plain[i] == 0 {
		i--
	}
	if i < 0 {
		return nil, 0, 0, &authenticatedRecordError{err: &ProtocolError{"protected record has no inner content type"}}
	}
	contentType = plain[i]
	if !validInnerContentType(contentType) {
		return nil, 0, 0, &authenticatedRecordError{err: &ProtocolError{"invalid inner content type"}}
	}
	if i == 0 && (contentType == recordTypeHandshake || contentType == recordTypeAlert) {
		return nil, 0, 0, authenticatedRecordAlert(alertUnexpectedMessage, &ProtocolError{"zero-length protected Handshake or Alert content"})
	}
	if !c.replay.check(sequence) {
		return nil, 0, 0, &ProtocolError{"duplicate or stale protected record"}
	}
	c.replay.accept(sequence)
	c.lastOpened = sequence
	if hasCID {
		c.lastConnectionID = matchedCID
	}
	if inPlace {
		return plain[:i], contentType, headerLen + n, nil
	}
	return append([]byte(nil), plain[:i]...), contentType, headerLen + n, nil
}

func equalBytes(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

func validInnerContentType(contentType uint8) bool {
	return contentType == recordTypeAlert || contentType == recordTypeHandshake || contentType == recordTypeApplicationData || contentType == recordTypeHeartbeat || contentType == recordTypeACK || contentType == recordTypeReturnRoutability
}

func reconstructSequence(expected, truncated uint64, bits uint) uint64 {
	window := uint64(1) << bits
	half := window / 2
	mask := window - 1
	candidate := (expected &^ mask) | truncated
	if candidate+half <= expected && candidate+window < 1<<48 {
		candidate += window
	} else if candidate > expected+half && candidate >= window {
		candidate -= window
	}
	return candidate
}
