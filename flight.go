package dtls13

import (
	"context"
	"sync"
	"time"
)

type handshakeMessage struct {
	typ      uint8
	sequence uint16
	body     []byte
}

type flightRecord struct {
	number         recordNumber
	priorNumber    recordNumber
	earlierNumbers []recordNumber
	wire           []byte
	acked          bool
	sent           bool
	hasPrior       bool
}

func (r *flightRecord) replaceNumber(number recordNumber) {
	if r.hasPrior {
		r.earlierNumbers = append(r.earlierNumbers, r.priorNumber)
	}
	r.priorNumber = r.number
	r.hasPrior = true
	r.number = number
}

func (r *flightRecord) hasNumber(number recordNumber) bool {
	if r.number == number {
		return true
	}
	if r.hasPrior && r.priorNumber == number {
		return true
	}
	for _, prior := range r.earlierNumbers {
		if prior == number {
			return true
		}
	}
	return false
}

func (r *flightRecord) historyLength() int {
	if r.hasPrior {
		return 1 + len(r.earlierNumbers)
	}
	return 0
}

func (r *flightRecord) acknowledgedBy(numbers []recordNumber) bool {
	for _, number := range numbers {
		if r.hasNumber(number) {
			return true
		}
	}
	return false
}

const initialRecordHistoryCapacity = 3

// reserveInitialRecordHistory batches the first history allocation for all
// records rebuilt in one flight generation. Each zero-length slice gets a
// disjoint window, so replaceNumber can append independently.
func reserveInitialRecordHistory(records []flightRecord, indices flightRecordIndices) {
	count := 0
	for i := 0; i < indices.count; i++ {
		record := &records[indices.at(i)]
		if record.hasPrior && cap(record.earlierNumbers) == 0 {
			count++
		}
	}
	if count == 0 {
		return
	}
	storage := make([]recordNumber, count*initialRecordHistoryCapacity)
	next := 0
	for i := 0; i < indices.count; i++ {
		record := &records[indices.at(i)]
		if !record.hasPrior || cap(record.earlierNumbers) != 0 {
			continue
		}
		record.earlierNumbers = storage[next : next : next+initialRecordHistoryCapacity]
		next += initialRecordHistoryCapacity
	}
}

type flightState uint8

const (
	flightPreparing flightState = iota
	flightSending
	flightWaiting
	flightFinished
)

type flightEventKind uint8

const (
	flightEventACK flightEventKind = iota
	flightEventNextFlight
	flightEventPeerRetransmit
)

type flightEvent struct {
	kind    flightEventKind
	numbers []recordNumber
}

type flightRecordIndices struct {
	inline   [10]int
	overflow []int
	count    int
}

func (s *flightRecordIndices) add(index, total int) {
	if s.count < len(s.inline) {
		s.inline[s.count] = index
	} else {
		if s.overflow == nil {
			s.overflow = make([]int, 0, total-len(s.inline))
		}
		s.overflow = append(s.overflow, index)
	}
	s.count++
}

func (s *flightRecordIndices) at(index int) int {
	if index < len(s.inline) {
		return s.inline[index]
	}
	return s.overflow[index-len(s.inline)]
}

// flight retains handshake fragments and their current wire images. Each
// retransmission is rebuilt with a fresh record sequence number while keeping
// the original handshake message_seq and bytes.
type flight struct {
	mu              sync.Mutex
	records         []flightRecord
	initialInterval time.Duration
	maxInterval     time.Duration
	state           flightState
	rebuildPending  func(flightRecordIndices) error
	resizeForMTU    func(int) error
	firstSentAt     time.Time
	retransmitted   bool
}

func (f *flight) noteSend(now time.Time, retransmission bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.firstSentAt.IsZero() {
		f.firstSentAt = now
	}
	if retransmission {
		f.retransmitted = true
	}
}

func (f *flight) rttSample(now time.Time) (time.Duration, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.firstSentAt.IsZero() || f.retransmitted || now.Before(f.firstSentAt) {
		return 0, false
	}
	return now.Sub(f.firstSentAt), true
}

func retainHandshakeMessages(messages []handshakeMessage) []handshakeMessage {
	return append([]handshakeMessage(nil), messages...)
}

func (f *flight) resizePlain(original []handshakeMessage, mtu int, epoch uint16) error {
	first := f.nextRecordSequence()
	rebuilt, _, err := buildPlainFlight(original, mtu, epoch, first)
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.records = rebuilt.records
	f.rebuildPending = rebuilt.rebuildPending
	f.mu.Unlock()
	return nil
}

func (f *flight) installPlainResize(messages []handshakeMessage, epoch uint16) {
	if len(messages) == 1 {
		original := messages[0]
		f.resizeForMTU = func(mtu int) error {
			var storage [1]handshakeMessage
			storage[0] = original
			return f.resizePlain(storage[:], mtu, epoch)
		}
		return
	}
	original := retainHandshakeMessages(messages)
	f.resizeForMTU = func(mtu int) error { return f.resizePlain(original, mtu, epoch) }
}

func (f *flight) resizeProtected(original []handshakeMessage, mtu int, cipher *recordCipher) error {
	rebuilt, err := buildProtectedFlight(original, mtu, cipher)
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.records = rebuilt.records
	f.rebuildPending = rebuilt.rebuildPending
	f.mu.Unlock()
	return nil
}

func (f *flight) installProtectedResize(messages []handshakeMessage, cipher *recordCipher) {
	if len(messages) == 1 {
		original := messages[0]
		f.resizeForMTU = func(mtu int) error {
			var storage [1]handshakeMessage
			storage[0] = original
			return f.resizeProtected(storage[:], mtu, cipher)
		}
		return
	}
	original := retainHandshakeMessages(messages)
	f.resizeForMTU = func(mtu int) error { return f.resizeProtected(original, mtu, cipher) }
}

// Flight builders retain message bodies as immutable input for future MTU
// refragmentation. Callers may reuse the message slice, but not body bytes.
func buildPlainFlight(messages []handshakeMessage, mtu int, epoch uint16, firstRecordSequence uint64) (*flight, uint64, error) {
	if epoch != 0 {
		return nil, firstRecordSequence, &ProtocolError{"plaintext handshake flight must use epoch zero"}
	}
	maxFragment := mtu - plainRecordHeaderLen - handshakeHeaderLen
	if recordMaximum := maxRecordContent - handshakeHeaderLen; maxFragment > recordMaximum {
		maxFragment = recordMaximum
	}
	if maxFragment < 1 {
		return nil, firstRecordSequence, &ConfigError{"MTU is too small for a handshake fragment"}
	}
	fragmentCount, err := countHandshakeFragments(messages, maxFragment)
	if err != nil {
		return nil, firstRecordSequence, err
	}
	f := &flight{state: flightPreparing, records: make([]flightRecord, 0, fragmentCount)}
	recordSequence := firstRecordSequence
	var fragmentStorage [10]handshakeFragment
	for _, message := range messages {
		fragments := fragmentHandshakeMessageInto(fragmentStorage[:0], message, maxFragment)
		for _, fragment := range fragments {
			payloadLen := handshakeHeaderLen + len(fragment.body)
			wire, err := allocatePlainRecordWire(recordTypeHandshake, epoch, recordSequence, payloadLen)
			if err != nil {
				return nil, firstRecordSequence, err
			}
			if err = marshalHandshakeFragmentInto(wire[plainRecordHeaderLen:], fragment); err != nil {
				return nil, firstRecordSequence, err
			}
			number := recordNumber{epoch: uint64(epoch), sequence: recordSequence}
			f.records = append(f.records, flightRecord{number: number, wire: wire})
			recordSequence++
		}
	}
	nextSequence := recordSequence
	f.rebuildPending = func(indices flightRecordIndices) error {
		reserveInitialRecordHistory(f.records, indices)
		for i := 0; i < indices.count; i++ {
			index := indices.at(i)
			var recordScratch [1]record
			records, err := parsePlainRecordsViewInto(f.records[index].wire, recordScratch[:0])
			if err != nil {
				return err
			}
			wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: nextSequence, payload: records[0].payload})
			if err != nil {
				return err
			}
			number := recordNumber{epoch: 0, sequence: nextSequence}
			f.records[index].wire = wire
			f.records[index].replaceNumber(number)
			nextSequence++
		}
		return nil
	}
	f.installPlainResize(messages, epoch)
	return f, recordSequence, nil
}

func buildProtectedFlight(messages []handshakeMessage, mtu int, cipher *recordCipher) (*flight, error) {
	if cipher == nil {
		return nil, &ProtocolError{"missing flight record cipher"}
	}
	maxFragment := mtu - cipher.headerLen16() - cipher.aead.Overhead() - 1 - handshakeHeaderLen
	if recordMaximum := maxRecordContent - handshakeHeaderLen; maxFragment > recordMaximum {
		maxFragment = recordMaximum
	}
	if maxFragment < 1 {
		return nil, &ConfigError{"MTU is too small for a protected handshake fragment"}
	}
	fragmentCount, err := countHandshakeFragments(messages, maxFragment)
	if err != nil {
		return nil, err
	}
	recordOverhead := cipher.headerLen16() + cipher.aead.Overhead() + 1 + handshakeHeaderLen
	maxInt := int(^uint(0) >> 1)
	if recordOverhead > 0 && fragmentCount > maxInt/recordOverhead {
		return nil, &ProtocolError{"protected flight size overflow"}
	}
	wireBytes := fragmentCount * recordOverhead
	for _, message := range messages {
		if len(message.body) > maxInt-wireBytes {
			return nil, &ProtocolError{"protected flight size overflow"}
		}
		wireBytes += len(message.body)
	}
	f := &flight{state: flightPreparing, records: make([]flightRecord, 0, fragmentCount)}
	// Keep fragment descriptors instead of a separately marshaled payload for
	// every record. The record cipher writes each descriptor directly into the
	// final plaintext window before sealing.
	fragments := make([]handshakeFragment, 0, fragmentCount)
	wireStorage := make([]byte, wireBytes)
	wireOffset := 0
	var fragmentStorage [10]handshakeFragment
	for _, message := range messages {
		for _, fragment := range fragmentHandshakeMessageInto(fragmentStorage[:0], message, maxFragment) {
			sequence := cipher.nextSequence
			recordLen := recordOverhead + len(fragment.body)
			window := wireStorage[wireOffset : wireOffset : wireOffset+recordLen]
			wire, err := cipher.sealHandshakeFragmentInto(window, fragment)
			if err != nil {
				return nil, err
			}
			if len(wire) > mtu {
				return nil, &ProtocolError{"protected handshake record exceeds MTU"}
			}
			fragments = append(fragments, fragment)
			number := recordNumber{epoch: cipher.epoch, sequence: sequence}
			f.records = append(f.records, flightRecord{number: number, wire: wire})
			wireOffset += recordLen
		}
	}
	f.rebuildPending = func(indices flightRecordIndices) error {
		reserveInitialRecordHistory(f.records, indices)
		recordOverhead := cipher.headerLen16() + cipher.aead.Overhead() + 1 + handshakeHeaderLen
		maxInt := int(^uint(0) >> 1)
		wireBytes := 0
		for i := 0; i < indices.count; i++ {
			recordLen := recordOverhead + len(fragments[indices.at(i)].body)
			if recordLen > maxInt-wireBytes {
				return &ProtocolError{"protected flight size overflow"}
			}
			wireBytes += recordLen
		}
		wireStorage := make([]byte, wireBytes)
		wireOffset := 0
		for i := 0; i < indices.count; i++ {
			index := indices.at(i)
			sequence := cipher.nextSequence
			fragment := fragments[index]
			recordLen := recordOverhead + len(fragment.body)
			window := wireStorage[wireOffset : wireOffset : wireOffset+recordLen]
			wire, err := cipher.sealHandshakeFragmentInto(window, fragment)
			if err != nil {
				return err
			}
			if len(wire) > mtu {
				return &ProtocolError{"retransmitted protected handshake record exceeds MTU"}
			}
			number := recordNumber{epoch: cipher.epoch, sequence: sequence}
			f.records[index].wire = wire
			f.records[index].replaceNumber(number)
			wireOffset += recordLen
		}
		return nil
	}
	f.installProtectedResize(messages, cipher)
	return f, nil
}

func combineFlightRecords(dst []flightRecord, first, second *flight) []flightRecord {
	total := len(first.records) + len(second.records)
	if cap(dst) < total {
		dst = make([]flightRecord, total)
	} else {
		dst = dst[:total]
	}
	offset := copy(dst, first.records)
	copy(dst[offset:], second.records)
	return dst
}

func combineFlights(first, second *flight) *flight {
	combined := &flight{state: flightPreparing}
	combined.records = combineFlightRecords(nil, first, second)
	combined.rebuildPending = func(_ flightRecordIndices) error {
		if err := first.refreshPending(); err != nil {
			return err
		}
		if err := second.refreshPending(); err != nil {
			return err
		}
		combined.records = combineFlightRecords(combined.records, first, second)
		return nil
	}
	combined.resizeForMTU = func(mtu int) error {
		if err := first.resize(mtu); err != nil {
			return err
		}
		if err := second.resize(mtu); err != nil {
			return err
		}
		combined.mu.Lock()
		combined.records = combineFlightRecords(combined.records, first, second)
		combined.mu.Unlock()
		return nil
	}
	return combined
}

func (f *flight) resize(mtu int) error {
	if f == nil || f.resizeForMTU == nil {
		return &ProtocolError{"flight cannot be refragmented"}
	}
	return f.resizeForMTU(mtu)
}

func countHandshakeFragments(messages []handshakeMessage, maxFragment int) (int, error) {
	total := 0
	maxInt := int(^uint(0) >> 1)
	for _, message := range messages {
		length := len(message.body)
		if length >= 1<<24 {
			return 0, &ProtocolError{"handshake message exceeds 24-bit length"}
		}
		count := 1
		if length > 0 {
			count = (length + maxFragment - 1) / maxFragment
		}
		if count > maxInt-total {
			return 0, &ProtocolError{"flight has too many handshake fragments"}
		}
		total += count
	}
	return total, nil
}

func fragmentHandshakeMessageInto(dst []handshakeFragment, message handshakeMessage, maxFragment int) []handshakeFragment {
	length := len(message.body)
	if length == 0 {
		return append(dst[:0], handshakeFragment{typ: message.typ, messageSequence: message.sequence, length: 0})
	}
	count := (length + maxFragment - 1) / maxFragment
	var out []handshakeFragment
	if cap(dst) < count {
		out = make([]handshakeFragment, 0, count)
	} else {
		out = dst[:0]
	}
	for off := 0; off < length; off += maxFragment {
		end := off + maxFragment
		if end > length {
			end = length
		}
		out = append(out, handshakeFragment{typ: message.typ, messageSequence: message.sequence, length: uint32(length), offset: uint32(off), body: message.body[off:end]})
	}
	return out
}

func (f *flight) setIntervals(initial, max time.Duration) {
	f.initialInterval = initial
	f.maxInterval = max
}
func (f *flight) setState(state flightState) { f.mu.Lock(); f.state = state; f.mu.Unlock() }
func (f *flight) currentState() flightState  { f.mu.Lock(); defer f.mu.Unlock(); return f.state }
func (f *flight) ackAll() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := range f.records {
		f.records[i].acked = true
	}
}
func (f *flight) ack(numbers []recordNumber) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	changed := 0
	for i := range f.records {
		if f.records[i].acked {
			continue
		}
		if f.records[i].acknowledgedBy(numbers) {
			f.records[i].acked = true
			changed++
		}
	}
	return changed
}
func (f *flight) pendingIndices(dst []int) []int {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := dst[:0]
	for index, record := range f.records {
		if !record.acked {
			if len(out) == cap(out) {
				grown := make([]int, len(out), len(f.records))
				copy(grown, out)
				out = grown
			}
			out = append(out, index)
		}
	}
	return out
}
func (f *flight) nextRecordSequence() uint64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	var next uint64
	for _, record := range f.records {
		if record.number.epoch == 0 && record.number.sequence >= next {
			next = record.number.sequence + 1
		}
	}
	return next
}
func (f *flight) refreshPending() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	var indices flightRecordIndices
	for index, record := range f.records {
		if !record.acked {
			indices.add(index, len(f.records))
		}
	}
	if indices.count == 0 || f.rebuildPending == nil {
		return nil
	}
	return f.rebuildPending(indices)
}
func (f *flight) complete() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, record := range f.records {
		if !record.acked {
			return false
		}
	}
	return true
}
func (f *flight) hasAcknowledgedRecord() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, record := range f.records {
		if record.acked {
			return true
		}
	}
	return false
}

// Wire windows alias immutable flight storage and are valid for synchronous
// reads. Rebuilds replace, rather than modify, published backing arrays.
func (f *flight) pendingWire(dst [][]byte) [][]byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := dst[:0]
	for _, record := range f.records {
		if !record.acked {
			if len(out) == cap(out) {
				grown := make([][]byte, len(out), len(f.records))
				copy(grown, out)
				out = grown
			}
			out = append(out, record.wire)
		}
	}
	return out
}

func (f *flight) nextUnsentWire(maxOutstanding int, dst [][]byte) [][]byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	outstanding := 0
	protected := false
	for _, record := range f.records {
		if record.number.epoch > 0 {
			protected = true
		}
		if record.sent && !record.acked {
			outstanding++
		}
	}
	capacity := len(f.records)
	if protected {
		capacity = maxOutstanding - outstanding
		if capacity < 0 {
			capacity = 0
		}
	}
	var out [][]byte
	if cap(dst) < capacity {
		out = make([][]byte, 0, capacity)
	} else {
		out = dst[:0]
	}
	for i := range f.records {
		if capacity == 0 {
			break
		}
		if f.records[i].acked || f.records[i].sent {
			continue
		}
		f.records[i].sent = true
		out = append(out, f.records[i].wire)
		capacity--
	}
	return out
}

func (f *flight) retransmitWire(maxRecords int, dst [][]byte) [][]byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out [][]byte
	if cap(dst) < maxRecords {
		out = make([][]byte, 0, maxRecords)
	} else {
		out = dst[:0]
	}
	for _, record := range f.records {
		if record.acked || !record.sent {
			continue
		}
		out = append(out, record.wire)
		if len(out) == maxRecords {
			break
		}
	}
	return out
}

// transmit sends all unacknowledged records immediately and after every
// timeout. ACKs may be applied concurrently through ack. The caller owns ACK
// parsing because ACK records arrive through the normal record layer.
func (f *flight) transmit(ctx context.Context, send func([]byte) error, acked <-chan struct{}) error {
	f.setState(flightSending)
	interval := f.initialInterval
	if interval <= 0 {
		interval = time.Second
	}
	max := f.maxInterval
	if max < interval {
		max = interval
	}
	sendPending := func(refresh bool) error {
		if refresh {
			if err := f.refreshPending(); err != nil {
				return err
			}
		}
		var storage [10][]byte
		for _, wire := range f.pendingWire(storage[:0]) {
			if err := send(wire); err != nil {
				return err
			}
		}
		return nil
	}
	if err := sendPending(false); err != nil {
		return err
	}
	if f.complete() {
		f.setState(flightFinished)
		return nil
	}
	f.setState(flightWaiting)
	timer := time.NewTimer(interval)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-acked:
			if f.complete() {
				f.setState(flightFinished)
				return nil
			}
			if err := sendPending(true); err != nil {
				return err
			}
			resetTimer(timer, interval)
		case <-timer.C:
			if err := sendPending(true); err != nil {
				return err
			}
			if interval < max {
				interval *= 2
				if interval > max {
					interval = max
				}
			}
			timer.Reset(interval)
		}
	}
}

func resetTimer(timer *time.Timer, interval time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(interval)
}

// runStateMachine drives RFC 9147's SENDING/WAITING portion of the flight
// state machine. For flights that are implicitly acknowledged, receipt of the
// next flight completes this flight. Peer retransmissions trigger an immediate
// retransmission without waiting for the timer.
func (f *flight) runStateMachine(ctx context.Context, send func([]byte) error, events <-chan flightEvent, explicitACKRequired bool) error {
	interval := f.initialInterval
	if interval <= 0 {
		interval = time.Second
	}
	max := f.maxInterval
	if max < interval {
		max = interval
	}
	sendPending := func(refresh bool) error {
		f.setState(flightSending)
		if refresh {
			if err := f.refreshPending(); err != nil {
				return err
			}
		}
		var storage [10][]byte
		for _, wire := range f.pendingWire(storage[:0]) {
			if err := send(wire); err != nil {
				return err
			}
		}
		f.setState(flightWaiting)
		return nil
	}
	if err := sendPending(false); err != nil {
		return err
	}
	if f.complete() {
		f.setState(flightFinished)
		return nil
	}
	timer := time.NewTimer(interval)
	defer timer.Stop()
	reset := func(d time.Duration) {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(d)
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event := <-events:
			switch event.kind {
			case flightEventACK:
				f.ack(event.numbers)
				if f.complete() {
					f.setState(flightFinished)
					return nil
				}
				if err := sendPending(true); err != nil {
					return err
				}
				reset(interval)
			case flightEventNextFlight:
				if !explicitACKRequired {
					f.ackAll()
					f.setState(flightFinished)
					return nil
				}
			case flightEventPeerRetransmit:
				if err := sendPending(true); err != nil {
					return err
				}
				reset(interval)
			}
		case <-timer.C:
			if err := sendPending(true); err != nil {
				return err
			}
			if interval < max {
				interval *= 2
				if interval > max {
					interval = max
				}
			}
			timer.Reset(interval)
		}
	}
}
