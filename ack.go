package dtls13

import (
	"encoding/binary"
	"slices"
)

// recordNumber is the expanded 128-bit representation from RFC 9147 section
// 4: a uint64 epoch followed by a uint64 sequence number.
type recordNumber struct {
	epoch    uint64
	sequence uint64
}

func marshalACK(numbers []recordNumber) ([]byte, error) {
	numbers = canonicalRecordNumbers(numbers)
	return marshalCanonicalACK(numbers)
}

func marshalCanonicalACK(numbers []recordNumber) ([]byte, error) {
	payloadLen, err := canonicalACKPayloadLen(numbers)
	if err != nil {
		return nil, err
	}
	b := make([]byte, payloadLen)
	marshalCanonicalACKInto(b, numbers)
	return b, nil
}

func canonicalACKPayloadLen(numbers []recordNumber) (int, error) {
	if len(numbers) > 4095 {
		return 0, &ProtocolError{"too many record numbers in ACK"}
	}
	return 2 + 16*len(numbers), nil
}

func marshalCanonicalACKInto(b []byte, numbers []recordNumber) {
	binary.BigEndian.PutUint16(b, uint16(16*len(numbers)))
	for i, n := range numbers {
		off := 2 + i*16
		binary.BigEndian.PutUint64(b[off:off+8], n.epoch)
		binary.BigEndian.PutUint64(b[off+8:off+16], n.sequence)
	}
}

func canonicalRecordNumbers(numbers []recordNumber) []recordNumber {
	canonical := true
	for i := 1; i < len(numbers); i++ {
		if !recordNumberLess(numbers[i-1], numbers[i]) {
			canonical = false
			break
		}
	}
	if canonical {
		return numbers
	}
	numbers = append([]recordNumber(nil), numbers...)
	slices.SortFunc(numbers, func(a, b recordNumber) int {
		if recordNumberLess(a, b) {
			return -1
		}
		if a != b {
			return 1
		}
		return 0
	})
	dedup := numbers[:0]
	for _, number := range numbers {
		if len(dedup) == 0 || dedup[len(dedup)-1] != number {
			dedup = append(dedup, number)
		}
	}
	return dedup
}

func buildACKRecords(numbers []recordNumber, mtu int, plainSequence uint64, cipher *recordCipher) (records [][]byte, nextPlainSequence uint64, err error) {
	return buildACKRecordsInto(nil, numbers, mtu, plainSequence, cipher)
}

func buildACKRecordsInto(dst [][]byte, numbers []recordNumber, mtu int, plainSequence uint64, cipher *recordCipher) (records [][]byte, nextPlainSequence uint64, err error) {
	records = dst[:0]
	overhead := plainRecordHeaderLen + 2
	if cipher != nil {
		overhead = cipher.headerLen16() + cipher.aead.Overhead() + 1 + 2
	}
	maxEntries := (mtu - overhead) / 16
	recordContentMaximum := maxRecordContent
	if cipher != nil {
		recordContentMaximum = cipher.maxContent()
	}
	if recordMaximum := (recordContentMaximum - 2) / 16; maxEntries > recordMaximum {
		maxEntries = recordMaximum
	}
	if maxEntries < 1 {
		return nil, plainSequence, &ConfigError{"MTU is too small for an ACK record"}
	}
	numbers = canonicalRecordNumbers(numbers)
	if len(numbers) == 0 {
		if cipher == nil {
			wire, marshalErr := marshalPlainACKRecord(nil, plainSequence)
			if marshalErr != nil {
				return nil, plainSequence, marshalErr
			}
			return append(records, wire), plainSequence + 1, nil
		}
		wire, marshalErr := cipher.sealWithHeaderBuilder(recordTypeACK, 2, true, true, func(dst []byte) {
			marshalCanonicalACKInto(dst, nil)
		})
		if marshalErr != nil {
			return nil, plainSequence, marshalErr
		}
		return append(records, wire), plainSequence, nil
	}
	for start := 0; start < len(numbers); start += maxEntries {
		end := start + maxEntries
		if end > len(numbers) {
			end = len(numbers)
		}
		var wire []byte
		var marshalErr error
		if cipher == nil {
			wire, marshalErr = marshalPlainACKRecord(numbers[start:end], plainSequence)
			plainSequence++
		} else {
			chunk := numbers[start:end]
			wire, marshalErr = cipher.sealWithHeaderBuilder(recordTypeACK, 2+16*len(chunk), true, true, func(dst []byte) {
				marshalCanonicalACKInto(dst, chunk)
			})
		}
		if marshalErr != nil {
			return nil, plainSequence, marshalErr
		}
		if len(wire) > mtu {
			return nil, plainSequence, &ProtocolError{"ACK record exceeds MTU"}
		}
		records = append(records, wire)
	}
	return records, plainSequence, nil
}

func marshalPlainACKRecord(numbers []recordNumber, sequence uint64) ([]byte, error) {
	payloadLen, err := canonicalACKPayloadLen(numbers)
	if err != nil {
		return nil, err
	}
	wire, err := allocatePlainRecordWire(recordTypeACK, 0, sequence, payloadLen)
	if err != nil {
		return nil, err
	}
	marshalCanonicalACKInto(wire[plainRecordHeaderLen:], numbers)
	return wire, nil
}

func parseACK(b []byte) ([]recordNumber, error) {
	return parseACKInto(b, nil)
}

func parseACKInto(b []byte, dst []recordNumber) ([]recordNumber, error) {
	if len(b) < 2 {
		return nil, &ProtocolError{"truncated ACK"}
	}
	n := int(binary.BigEndian.Uint16(b[:2]))
	if n%16 != 0 || n != len(b)-2 {
		return nil, &ProtocolError{"malformed ACK vector length"}
	}
	count := n / 16
	if count == 0 {
		return nil, nil
	}
	if cap(dst) < count {
		dst = make([]recordNumber, count)
	} else {
		dst = dst[:count]
	}
	for i := range dst {
		off := 2 + i*16
		dst[i] = recordNumber{epoch: binary.BigEndian.Uint64(b[off : off+8]), sequence: binary.BigEndian.Uint64(b[off+8 : off+16])}
		if i > 0 && !recordNumberLess(dst[i-1], dst[i]) {
			return nil, &ProtocolError{"ACK record numbers are not in increasing order"}
		}
	}
	return dst, nil
}

func validateACKEpoch(numbers []recordNumber, ackEpoch uint64) error {
	for _, number := range numbers {
		if number.epoch > ackEpoch {
			return alertError(alertIllegalParameter, &ProtocolError{"ACK epoch is lower than an acknowledged record epoch"})
		}
	}
	return nil
}
