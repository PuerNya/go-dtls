package dtls13

import "encoding/binary"

const (
	recordTypeChangeCipherSpec uint8  = 20
	recordTypeAlert            uint8  = 21
	recordTypeHandshake        uint8  = 22
	recordTypeApplicationData  uint8  = 23
	recordTypeHeartbeat        uint8  = 24
	recordTypeACK              uint8  = 26
	dtlsLegacyVersion          uint16 = 0xfefd
	plainRecordHeaderLen              = 13
	maxRecordContent                  = 1 << 14
)

type record struct {
	typ      uint8
	epoch    uint16
	sequence uint64 // low 48 bits are encoded
	payload  []byte
}

func marshalPlainRecord(r record) ([]byte, error) {
	b, err := allocatePlainRecordWire(r.typ, r.epoch, r.sequence, len(r.payload))
	if err != nil {
		return nil, err
	}
	copy(b[plainRecordHeaderLen:], r.payload)
	return b, nil
}

func allocatePlainRecordWire(typ uint8, epoch uint16, sequence uint64, payloadLen int) ([]byte, error) {
	if !validPlainContentType(typ) {
		return nil, &ProtocolError{"invalid DTLS 1.3 plaintext content type"}
	}
	if epoch != 0 {
		return nil, &ProtocolError{"DTLSPlaintext epoch must be zero"}
	}
	if sequence >= 1<<48 {
		return nil, &ProtocolError{"record sequence number exhausted"}
	}
	if payloadLen < 0 || payloadLen > maxRecordContent {
		return nil, &ProtocolError{"record payload exceeds 2^14 bytes"}
	}
	b := make([]byte, plainRecordHeaderLen+payloadLen)
	b[0] = typ
	binary.BigEndian.PutUint16(b[1:3], dtlsLegacyVersion)
	binary.BigEndian.PutUint16(b[3:5], epoch)
	putUint48(b[5:11], sequence)
	binary.BigEndian.PutUint16(b[11:13], uint16(payloadLen))
	return b, nil
}

func parsePlainRecords(datagram []byte) ([]record, error) {
	return parsePlainRecordsMode(datagram, true, nil)
}

func parsePlainRecordsView(datagram []byte) ([]record, error) {
	return parsePlainRecordsMode(datagram, false, nil)
}

func parsePlainRecordsViewInto(datagram []byte, dst []record) ([]record, error) {
	return parsePlainRecordsMode(datagram, false, dst)
}

func parsePlainRecordsMode(datagram []byte, copyPayload bool, dst []record) ([]record, error) {
	if len(datagram) == 0 {
		return nil, nil
	}
	out := dst[:0]
	for len(datagram) != 0 {
		if len(datagram) < plainRecordHeaderLen {
			return nil, &ProtocolError{"truncated record header"}
		}
		if !validPlainContentType(datagram[0]) {
			return nil, &ProtocolError{"invalid DTLS 1.3 plaintext content type"}
		}
		n := int(binary.BigEndian.Uint16(datagram[11:13]))
		if n > maxRecordContent || len(datagram) < plainRecordHeaderLen+n {
			return nil, &ProtocolError{"invalid record length"}
		}
		epoch := binary.BigEndian.Uint16(datagram[3:5])
		if epoch != 0 {
			return nil, &ProtocolError{"DTLSPlaintext epoch must be zero"}
		}
		var payload []byte
		if n > 0 {
			payload = datagram[13 : 13+n]
		}
		if copyPayload && len(payload) > 0 {
			payload = append([]byte(nil), payload...)
		}
		out = append(out, record{typ: datagram[0], epoch: epoch, sequence: getUint48(datagram[5:11]), payload: payload})
		datagram = datagram[13+n:]
	}
	return out, nil
}

func validPlainContentType(contentType uint8) bool {
	return contentType == recordTypeAlert || contentType == recordTypeHandshake || contentType == recordTypeACK
}

func putUint48(dst []byte, v uint64) {
	dst[0] = byte(v >> 40)
	dst[1] = byte(v >> 32)
	dst[2] = byte(v >> 24)
	dst[3] = byte(v >> 16)
	dst[4] = byte(v >> 8)
	dst[5] = byte(v)
}
func getUint48(src []byte) uint64 {
	return uint64(src[0])<<40 | uint64(src[1])<<32 | uint64(src[2])<<24 | uint64(src[3])<<16 | uint64(src[4])<<8 | uint64(src[5])
}
