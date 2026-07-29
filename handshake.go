package dtls13

import "encoding/binary"

const handshakeHeaderLen = 12

type handshakeFragment struct {
	typ             uint8
	messageSequence uint16
	length          uint32
	offset          uint32
	body            []byte
}

func marshalHandshakeFragment(f handshakeFragment) ([]byte, error) {
	if err := validateHandshakeFragment(f); err != nil {
		return nil, err
	}
	b := make([]byte, handshakeHeaderLen+len(f.body))
	putHandshakeFragment(b, f)
	return b, nil
}

func marshalHandshakeFragmentInto(dst []byte, f handshakeFragment) error {
	if err := validateHandshakeFragment(f); err != nil {
		return err
	}
	if len(dst) != handshakeHeaderLen+len(f.body) {
		return &ProtocolError{"invalid handshake fragment destination length"}
	}
	putHandshakeFragment(dst, f)
	return nil
}

func validateHandshakeFragment(f handshakeFragment) error {
	if f.length >= 1<<24 || f.offset >= 1<<24 || uint64(f.offset)+uint64(len(f.body)) > uint64(f.length) {
		return &ProtocolError{"invalid handshake fragment bounds"}
	}
	return nil
}

func putHandshakeFragment(dst []byte, f handshakeFragment) {
	dst[0] = f.typ
	putUint24(dst[1:4], f.length)
	binary.BigEndian.PutUint16(dst[4:6], f.messageSequence)
	putUint24(dst[6:9], f.offset)
	putUint24(dst[9:12], uint32(len(f.body)))
	copy(dst[12:], f.body)
}

func parseHandshakeFragments(b []byte) ([]handshakeFragment, error) {
	return parseHandshakeFragmentsMode(b, true, nil)
}

func parseHandshakeFragmentsView(b []byte) ([]handshakeFragment, error) {
	return parseHandshakeFragmentsMode(b, false, nil)
}

func parseHandshakeFragmentsViewInto(b []byte, dst []handshakeFragment) ([]handshakeFragment, error) {
	return parseHandshakeFragmentsMode(b, false, dst)
}

func parseHandshakeFragmentsMode(b []byte, copyBody bool, dst []handshakeFragment) ([]handshakeFragment, error) {
	if len(b) == 0 {
		return nil, nil
	}
	out := dst[:0]
	for len(b) > 0 {
		if len(b) < handshakeHeaderLen {
			return nil, &ProtocolError{"truncated handshake fragment header"}
		}
		l, off, n := getUint24(b[1:4]), getUint24(b[6:9]), getUint24(b[9:12])
		if uint64(off)+uint64(n) > uint64(l) || int(n) > len(b)-handshakeHeaderLen {
			return nil, &ProtocolError{"invalid handshake fragment"}
		}
		var body []byte
		if n > 0 {
			body = b[12 : 12+int(n)]
		}
		if copyBody && len(body) > 0 {
			body = append([]byte(nil), body...)
		}
		out = append(out, handshakeFragment{typ: b[0], length: l, messageSequence: binary.BigEndian.Uint16(b[4:6]), offset: off, body: body})
		b = b[12+int(n):]
	}
	return out, nil
}
func putUint24(b []byte, v uint32) { b[0] = byte(v >> 16); b[1] = byte(v >> 8); b[2] = byte(v) }
func getUint24(b []byte) uint32    { return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]) }

type reassembler struct {
	messages             map[uint16]*partialMessage
	maxSize              int
	maxMessages          int
	maxBytes             int
	allocated            int
	completedFirstRecord recordNumber
	completedLastRecord  recordNumber
}
type partialMessage struct {
	typ         uint8
	epoch       uint64
	hasEpoch    bool
	firstRecord recordNumber
	lastRecord  recordNumber
	hasRecord   bool
	data        []byte
	received    []uint64
	count       int
}

func bitmapRangeClear(bitmap []uint64, start, end int) bool {
	for start < end {
		word := start >> 6
		wordEnd := (word + 1) << 6
		if wordEnd > end {
			wordEnd = end
		}
		mask := ^uint64(0) << uint(start&63)
		if high := wordEnd & 63; high != 0 {
			mask &= (uint64(1) << uint(high)) - 1
		}
		if bitmap[word]&mask != 0 {
			return false
		}
		start = wordEnd
	}
	return true
}

func setBitmapRange(bitmap []uint64, start, end int) {
	for start < end {
		word := start >> 6
		wordEnd := (word + 1) << 6
		if wordEnd > end {
			wordEnd = end
		}
		mask := ^uint64(0) << uint(start&63)
		if high := wordEnd & 63; high != 0 {
			mask &= (uint64(1) << uint(high)) - 1
		}
		bitmap[word] |= mask
		start = wordEnd
	}
}

func newReassembler() *reassembler { return newReassemblerWithLimits(1<<20, 8, 4<<20) }
func newReassemblerWithLimit(maxSize int) *reassembler {
	return newReassemblerWithLimits(maxSize, 8, 4*maxSize)
}
func newReassemblerWithLimits(maxSize, maxMessages, maxBytes int) *reassembler {
	return &reassembler{messages: make(map[uint16]*partialMessage), maxSize: maxSize, maxMessages: maxMessages, maxBytes: maxBytes}
}
func (r *reassembler) add(f handshakeFragment) ([]byte, bool, error) {
	return r.addAtEpoch(f, 0, false)
}

func (r *reassembler) addProtected(f handshakeFragment, epoch uint64) ([]byte, bool, error) {
	return r.addAtEpoch(f, epoch, true)
}

func (r *reassembler) addProtectedRecord(f handshakeFragment, number recordNumber) ([]byte, bool, recordNumber, recordNumber, error) {
	body, complete, err := r.addAtEpochAndRecord(f, number.epoch, true, number, true)
	if err != nil || !complete {
		return body, complete, recordNumber{}, recordNumber{}, err
	}
	return body, true, r.completedFirstRecord, r.completedLastRecord, nil
}

func (r *reassembler) addAtEpoch(f handshakeFragment, epoch uint64, protected bool) ([]byte, bool, error) {
	return r.addAtEpochAndRecord(f, epoch, protected, recordNumber{}, false)
}

func (r *reassembler) addAtEpochAndRecord(f handshakeFragment, epoch uint64, protected bool, number recordNumber, hasRecord bool) ([]byte, bool, error) {
	if int64(f.length) > int64(r.maxSize) {
		return nil, false, &ProtocolError{"handshake message exceeds configured limit"}
	}
	p := r.messages[f.messageSequence]
	if p == nil {
		if len(r.messages) >= r.maxMessages {
			return nil, false, &ProtocolError{"too many incomplete handshake messages"}
		}
		if int64(r.allocated)+int64(f.length) > int64(r.maxBytes) {
			return nil, false, &ProtocolError{"handshake reassembly memory limit exceeded"}
		}
		if f.offset == 0 && len(f.body) == int(f.length) {
			r.completedFirstRecord = number
			r.completedLastRecord = number
			if len(f.body) == 0 {
				return []byte{}, true, nil
			}
			return append([]byte(nil), f.body...), true, nil
		}
		p = &partialMessage{typ: f.typ, epoch: epoch, hasEpoch: protected, data: make([]byte, int(f.length)), received: make([]uint64, (int(f.length)+63)/64)}
		r.messages[f.messageSequence] = p
		r.allocated += int(f.length)
	}
	if p.typ != f.typ || len(p.data) != int(f.length) {
		return nil, false, &ProtocolError{"inconsistent handshake fragments"}
	}
	if p.hasEpoch != protected || protected && p.epoch != epoch {
		return nil, false, alertError(alertUnexpectedMessage, &ProtocolError{"handshake message fragments span a key change"})
	}
	if hasRecord {
		if !p.hasRecord {
			p.firstRecord = number
			p.lastRecord = number
			p.hasRecord = true
		} else {
			if recordNumberLess(number, p.firstRecord) {
				p.firstRecord = number
			}
			if recordNumberLess(p.lastRecord, number) {
				p.lastRecord = number
			}
		}
	}
	start := int(f.offset)
	end := start + len(f.body)
	if bitmapRangeClear(p.received, start, end) {
		copy(p.data[start:end], f.body)
		setBitmapRange(p.received, start, end)
		p.count += len(f.body)
	} else {
		for i, v := range f.body {
			at := start + i
			word := at >> 6
			bit := uint64(1) << uint(at&63)
			seen := p.received[word]&bit != 0
			if seen && p.data[at] != v {
				return nil, false, &ProtocolError{"overlapping handshake fragments differ"}
			}
			if !seen {
				p.received[word] |= bit
				p.count++
			}
			p.data[at] = v
		}
	}
	if p.count == len(p.data) {
		delete(r.messages, f.messageSequence)
		r.allocated -= len(p.data)
		r.completedFirstRecord = p.firstRecord
		r.completedLastRecord = p.lastRecord
		return p.data, true, nil
	}
	return nil, false, nil
}

func (r *reassembler) hasIncompleteProtected() bool {
	for _, message := range r.messages {
		if message.hasEpoch {
			return true
		}
	}
	return false
}
