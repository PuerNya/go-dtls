package dtls13

import (
	"context"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ConnectionState records negotiated parameters and authenticated peer
// identity for a DTLS association. A zero value describes a connection whose
// handshake has not completed.
//
// The certificate slices refer to immutable certificate objects owned by the
// connection and must not be modified. Connection ID slices returned by
// [Conn.ConnectionState] are copies.
type ConnectionState struct {
	// Version is the negotiated DTLS version. It is VersionDTLS13 after a
	// successful handshake.
	Version uint16
	// HandshakeComplete is true after the initial handshake has completed and
	// installed application traffic keys.
	HandshakeComplete bool
	// DidResume is true when the connection used a PSK from a session ticket
	// instead of performing a full certificate handshake.
	DidResume bool
	// CipherSuite is the negotiated TLS 1.3 cipher-suite identifier.
	CipherSuite uint16
	// NegotiatedProtocol is the ALPN protocol selected by the server, or the
	// empty string when ALPN was not negotiated.
	NegotiatedProtocol string
	// ServerName is the name sent by a client and used for certificate hostname
	// verification when built-in verification is enabled. It is empty in
	// server-side state.
	ServerName string
	// PeerCertificates contains the certificate chain presented by the peer,
	// with the leaf first. Resumed connections restore this state from the
	// session. It can be empty when a server did not request a client certificate.
	PeerCertificates []*x509.Certificate
	// VerifiedChains contains the chains built during certificate verification.
	// Resumed connections restore chains after current built-in policy checks.
	// It is nil when built-in verification was skipped or the peer did not send
	// a certificate.
	VerifiedChains [][]*x509.Certificate
	// LocalConnectionID is the CID the peer currently uses when sending
	// protected records to this endpoint. It is empty when no non-empty CID is
	// active in that direction.
	LocalConnectionID []byte
	// PeerConnectionID is the CID this endpoint currently places in protected
	// records sent to the peer. It is empty when no non-empty CID is active in
	// that direction.
	PeerConnectionID []byte
	// ReturnRoutabilityCheck is true when RFC 9853 was negotiated together with
	// Connection ID for this association.
	ReturnRoutabilityCheck bool
	// RecordSizeLimitNegotiated is true when RFC 8449 record_size_limit was
	// negotiated for this connection.
	RecordSizeLimitNegotiated bool
	// LocalRecordSizeLimit is the maximum complete protected plaintext this
	// endpoint accepts. It is the DTLS 1.3 protocol maximum when the extension
	// was not negotiated.
	LocalRecordSizeLimit uint16
	// PeerRecordSizeLimit is the maximum complete protected plaintext this
	// endpoint sends. It is the DTLS 1.3 protocol maximum when the extension
	// was not negotiated or the peer advertised a larger future value.
	PeerRecordSizeLimit uint16
	exporter            *exporterState
}

// ExportKeyingMaterial returns exporter output for the completed connection,
// following RFC 8446 section 7.5 with the DTLS 1.3 label prefix required by
// RFC 9147. A nil and an empty context are equivalent.
//
// The label is application-defined and must fit the TLS HKDF label encoding.
// The maximum output is 255 times the negotiated hash length. The method
// returns an error before handshake completion, after the connection's secrets
// have been cleared, or for an invalid label or length.
func (s ConnectionState) ExportKeyingMaterial(label string, context []byte, length int) ([]byte, error) {
	if s.exporter == nil {
		return nil, errors.New("dtls13: exporter is unavailable before handshake completion")
	}
	return s.exporter.export(label, context, length)
}

// Conn represents one DTLS association over a connected datagram transport.
// It exposes unreliable, message-oriented I/O and intentionally implements
// neither [net.Conn] nor [net.PacketConn].
//
// Unless Dial completed it already, the first call to Handshake,
// ReadDatagram, WriteDatagram, or another operation that requires traffic keys
// performs the handshake. Conn methods synchronize access to protocol state;
// reads and writes may run concurrently. Concurrent reads are serialized, as
// are concurrent writes.
type Conn struct {
	conn     net.Conn
	config   *Config
	isClient bool

	handshakeOnce                    sync.Once
	handshakeErr                     error
	mu                               sync.RWMutex
	state                            ConnectionState
	readMu                           sync.Mutex
	dispatchMu                       sync.Mutex
	inputOnce                        sync.Once
	inputMu                          sync.Mutex
	readNotify                       chan struct{}
	readErr                          error
	peerReadClosed                   bool
	readerMu                         sync.Mutex
	readerRunning                    bool
	readerClosed                     bool
	writeMu                          sync.Mutex
	applicationDatagrams             []applicationDatagram
	bufferedApplicationBytes         int
	sendCipher                       *recordCipher
	receiveEpochs                    *epochSet
	closure                          closureState
	handshakeDeadline                time.Time
	sendingTraffic                   *sendingTraffic
	receivingTraffic                 *receivingTraffic
	finishedACKCipher                *recordCipher
	finishedFlightStart              uint16
	finishedMessageSequence          uint16
	completedPeerFlightStart         uint16
	completedPeerFlightEnd           uint16
	hasCompletedPeerFlight           bool
	resumptionSuite                  *cipherSuite
	resumptionMasterSecret           []byte
	ticketFlight                     *flight
	postHandshakeReassembly          *reassembler
	protectedHandshakeRanges         []protectedHandshakeRecordRange
	recentApplicationRecords         []recordNumber
	pendingHandshakeApplications     []pendingPostAuthApplication
	pendingHandshakeApplicationBytes int
	postHandshakeTranscript          *transcriptHash
	postHandshakeAuthOffered         bool
	postHandshakeAuthCounter         atomic.Uint64
	postHandshakeAuthState           *postHandshakeAuthState
	clientAuthRequestFlight          *flight
	clientAuthResponseFlight         *flight
	lastClientAuthRequestSeq         uint16
	hasClientAuthRequestSeq          bool
	completedClientAuthStart         uint16
	completedClientAuthEnd           uint16
	hasCompletedClientAuth           bool
	sendConnectionID                 []byte
	receiveConnectionID              []byte
	connectionIDNegotiated           bool
	returnRoutabilityCheckNegotiated bool
	recordSizeLimitNegotiated        bool
	localRecordSizeLimit             uint16
	peerRecordSizeLimit              uint16
	localCIDUpdatesAllowed           bool
	peerCIDUpdatesAllowed            bool
	newConnectionIDFlight            *flight
	requestCIDFlight                 *flight
	connectionIDRequestOpen          bool
	keyUpdateResponsePending         bool
	peerSpareConnectionIDs           [][]byte
	lastNewCIDSequence               uint16
	hasNewCIDSequence                bool
	lastRequestCIDSequence           uint16
	hasRequestCIDSequence            bool
	cidGenMu                         sync.Mutex
	earlyMu                          sync.Mutex
	earlyPending                     []byte
	earlySignaled                    bool
	earlyReadDatagrams               [][]byte
	earlyReadBytes                   int
	earlyAccepted                    bool
	earlyDataLimit                   uint32
	earlySent                        bool
	earlyRejected                    bool
	pathMTU                          atomic.Int64
	plainSendSequence                atomic.Uint64
	retransmitNanos                  atomic.Int64
	lastRTTSampleUnixNano            atomic.Int64
	returnRoutability                *returnRoutabilityState
}

type applicationDatagram struct {
	payload []byte
	from    net.Addr
}

const maxDatagramSize = 65535

var datagramBufferPool = sync.Pool{New: func() any {
	return new([maxDatagramSize]byte)
}}

func acquireDatagramBuffer() *[maxDatagramSize]byte {
	return datagramBufferPool.Get().(*[maxDatagramSize]byte)
}

func releaseDatagramBuffer(buffer *[maxDatagramSize]byte) {
	datagramBufferPool.Put(buffer)
}

// DatagramInfo describes one authenticated Application Data record consumed by
// [Conn.ReadDatagram].
type DatagramInfo struct {
	// Source is the network address from which the authenticated record was
	// received. It is informational: Connection IDs do not by themselves
	// validate a changed network path or alter the address used for replies.
	Source net.Addr
	// FullLength is the length of the complete plaintext datagram before it was
	// copied into the caller's buffer.
	FullLength int
	// Truncated reports whether the caller's buffer was shorter than FullLength.
	// The unread remainder has already been discarded.
	Truncated bool
}

type protectedHandshakeRecordRange struct {
	first recordNumber
	last  recordNumber
}

func (c *Conn) flightInterval() time.Duration {
	current := time.Duration(c.retransmitNanos.Load())
	if current <= 0 {
		return c.config.FlightInterval
	}
	last := time.Unix(0, c.lastRTTSampleUnixNano.Load())
	if c.config.Time().Sub(last) >= 10*current {
		c.retransmitNanos.Store(0)
		return c.config.FlightInterval
	}
	return current
}

func (c *Conn) observeFlightRTT(flight *flight) {
	now := c.config.Time()
	sample, ok := flight.rttSample(now)
	if !ok {
		return
	}
	interval := sample + sample/2
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	if interval > c.config.MaxFlightInterval {
		interval = c.config.MaxFlightInterval
	}
	c.retransmitNanos.Store(int64(interval))
	c.lastRTTSampleUnixNano.Store(now.UnixNano())
}

func (c *Conn) currentMTU() int {
	if mtu := c.pathMTU.Load(); mtu != 0 {
		return int(mtu)
	}
	mtu := c.config.MTU
	if c.pathMTU.CompareAndSwap(0, int64(mtu)) {
		return mtu
	}
	return int(c.pathMTU.Load())
}

func (c *Conn) reducePathMTU() (int, bool) {
	for {
		current := c.currentMTU()
		floor := c.pathMTUFloor()
		if current <= floor {
			return current, false
		}
		next := current * 3 / 4
		if next < floor {
			next = floor
		}
		if c.pathMTU.CompareAndSwap(int64(current), int64(next)) {
			return next, true
		}
	}
}

func (c *Conn) pathMTUFloor() int {
	if c.conn != nil {
		if udp, ok := c.conn.LocalAddr().(*net.UDPAddr); ok {
			if udp.IP.To4() != nil {
				return 576 - 20 - 8
			}
			return 1280 - 40 - 8
		}
	}
	return 256
}

// PathMTU returns the current maximum transport datagram size used by the
// connection, including DTLS record framing. It starts at Config.MTU and can
// decrease in response to path-MTU write errors or repeated handshake
// timeouts. It returns zero before an effective configuration is available.
//
// For an established connection, PathMTU()-RecordOverhead() is the current
// upper bound for an application datagram, subject also to the DTLS 2^14-byte
// record-content limit. The path can change after this check, so callers must
// still handle [ErrDatagramTooLarge] from [Conn.WriteDatagram].
// [Config.IgnorePathMTU] makes application writes ignore this estimate;
// handshake traffic continues to use it.
func (c *Conn) PathMTU() int {
	if c.config == nil {
		return 0
	}
	return c.currentMTU()
}

// RecordOverhead returns the number of bytes added around one application
// datagram in the current sending epoch, including record framing, the inner
// content type, and the AEAD tag. Before handshake keys exist it returns the
// plaintext record-header size and must not be used to size application data.
func (c *Conn) RecordOverhead() int {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.sendCipher == nil {
		return plainRecordHeaderLen
	}
	return c.sendCipher.headerLen16() + c.sendCipher.aead.Overhead() + 1
}

func isMessageTooLong(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, ErrDatagramTooLarge) || errors.Is(err, syscall.EMSGSIZE) || errors.Is(err, syscall.Errno(10040)) || strings.Contains(strings.ToLower(err.Error()), "message too long")
}

func normalizeDatagramWriteError(err error, addr net.Addr) error {
	if isMessageTooLong(err) {
		return datagramTooLargeError(addr)
	}
	return err
}

func (c *Conn) writeRecord(wire []byte) error {
	_, err := c.conn.Write(wire)
	if err == nil {
		return nil
	}
	return normalizeDatagramWriteError(err, c.RemoteAddr())
}

// Client returns a client-side DTLS association over conn. The underlying
// connection must be a connected datagram transport that preserves one write
// as one datagram, normally a *net.UDPConn; a stream connection is invalid.
//
// Client does not perform network I/O. The first operation that needs traffic
// keys performs the handshake. Close closes conn. A nil config selects default
// settings, but a client that verifies certificates usually needs RootCAs and
// ServerName.
func Client(conn net.Conn, config *Config) *Conn {
	return &Conn{conn: conn, config: config, isClient: true}
}

// Server returns a server-side DTLS association over conn. The underlying
// connection must be a connected datagram transport that preserves one write
// as one datagram. Server defers the handshake until the first operation that
// needs traffic keys. Close closes conn.
func Server(conn net.Conn, config *Config) *Conn { return &Conn{conn: conn, config: config} }

// Dial connects to address on a UDP network and performs a client handshake.
// Network must be "udp", "udp4", or "udp6". If Config.ServerName is empty,
// Dial derives it from the host portion of address. The returned Conn owns the
// underlying socket and must be closed by the caller.
func Dial(network, address string, config *Config) (*Conn, error) {
	return DialWithDialer(&net.Dialer{}, network, address, config)
}

// DialWithDialer is like Dial but uses dialer to create the connected UDP
// transport. A nil dialer is equivalent to a zero net.Dialer. A positive
// dialer Timeout also bounds the DTLS handshake; Config.HandshakeTimeout still
// applies when it is shorter.
func DialWithDialer(dialer *net.Dialer, network, address string, config *Config) (*Conn, error) {
	if network != "udp" && network != "udp4" && network != "udp6" {
		return nil, &ConfigError{"Dial network must be udp, udp4, or udp6"}
	}
	if dialer == nil {
		dialer = &net.Dialer{}
	}
	clientConfig := config
	if clientConfig == nil {
		clientConfig = &Config{}
	} else {
		clientConfig = clientConfig.Clone()
	}
	if clientConfig.ServerName == "" {
		host, _, splitErr := net.SplitHostPort(address)
		if splitErr != nil {
			return nil, splitErr
		}
		clientConfig.ServerName = host
	}
	raw, err := dialer.Dial(network, address)
	if err != nil {
		return nil, err
	}
	c := Client(raw, clientConfig)
	ctx := context.Background()
	if dialer.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dialer.Timeout)
		defer cancel()
	}
	if err = c.HandshakeContext(ctx); err != nil {
		_ = raw.Close()
		return nil, err
	}
	return c, nil
}

// Handshake performs the DTLS handshake if it has not already run. Subsequent
// calls return the result of the first handshake attempt.
func (c *Conn) Handshake() error { return c.HandshakeContext(context.Background()) }

// HandshakeContext performs the DTLS handshake if it has not already run. The
// handshake is bounded by the earlier of ctx's deadline and
// Config.HandshakeTimeout. Canceling ctx interrupts transport I/O.
//
// The first call controls the handshake and all later calls return its result;
// a canceled or failed handshake is not retried. The context must be non-nil.
func (c *Conn) HandshakeContext(ctx context.Context) error {
	c.handshakeOnce.Do(func() {
		if c.conn == nil {
			c.handshakeErr = &ConfigError{"nil underlying connection"}
			return
		}
		cfg, err := c.config.normalized()
		if err != nil {
			c.handshakeErr = err
			return
		}
		c.config = cfg
		c.handshakeErr = c.runHandshake(ctx)
	})
	if c.handshakeErr == nil {
		c.startRecordReader()
	}
	return c.handshakeErr
}

func (c *Conn) initInput() {
	c.inputOnce.Do(func() { c.readNotify = make(chan struct{}, 1) })
}

func (c *Conn) notifyRead() {
	select {
	case c.readNotify <- struct{}{}:
	default:
	}
}

func (c *Conn) startRecordReader() {
	c.initInput()
	c.readerMu.Lock()
	if c.readerRunning || c.readerClosed {
		c.readerMu.Unlock()
		return
	}
	c.readerRunning = true
	c.readerMu.Unlock()
	go c.readRecords()
}

func (c *Conn) finishRecordReader(err error) {
	c.readerMu.Lock()
	c.readerRunning = false
	closed := c.readerClosed
	c.readerMu.Unlock()
	if err == nil || closed {
		return
	}
	c.inputMu.Lock()
	if c.readErr == nil {
		c.readErr = err
	}
	c.inputMu.Unlock()
	c.notifyRead()
}

func (c *Conn) queueApplicationData(content []byte, from net.Addr) error {
	c.inputMu.Lock()
	if len(c.applicationDatagrams) >= c.config.MaxBufferedApplicationDatagrams {
		c.inputMu.Unlock()
		return &ProtocolError{"buffered application datagram limit exceeded"}
	}
	if len(content) > c.config.MaxBufferedApplicationData-c.bufferedApplicationBytes {
		c.inputMu.Unlock()
		return &ProtocolError{"buffered application data limit exceeded"}
	}
	c.applicationDatagrams = append(c.applicationDatagrams, applicationDatagram{
		payload: append([]byte(nil), content...), from: from,
	})
	c.bufferedApplicationBytes += len(content)
	c.inputMu.Unlock()
	c.notifyRead()
	return nil
}

func (c *Conn) rememberApplicationRecordLocked(number recordNumber) error {
	for _, recordRange := range c.protectedHandshakeRanges {
		if recordNumberLess(recordRange.first, number) && recordNumberLess(number, recordRange.last) {
			return alertError(alertUnexpectedMessage, &ProtocolError{"application data interleaved with a protected handshake message"})
		}
	}
	c.recentApplicationRecords = append(c.recentApplicationRecords, number)
	c.trimRecordOrderingHistoryLocked()
	return nil
}

func (c *Conn) rememberProtectedHandshakeRangeLocked(first, last recordNumber) error {
	for _, number := range c.recentApplicationRecords {
		if recordNumberLess(first, number) && recordNumberLess(number, last) {
			return alertError(alertUnexpectedMessage, &ProtocolError{"application data interleaved with a protected handshake message"})
		}
	}
	for _, existing := range c.protectedHandshakeRanges {
		if existing.first == first && existing.last == last {
			return nil
		}
	}
	c.protectedHandshakeRanges = append(c.protectedHandshakeRanges, protectedHandshakeRecordRange{first: first, last: last})
	c.trimRecordOrderingHistoryLocked()
	return nil
}

func (c *Conn) rememberProtectedHandshakeRange(first, last recordNumber) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.rememberProtectedHandshakeRangeLocked(first, last); err != nil {
		return err
	}
	if c.postHandshakeReassembly != nil && c.postHandshakeReassembly.hasIncompleteProtected() {
		return nil
	}
	for _, application := range c.pendingHandshakeApplications {
		if err := c.queueApplicationData(application.content, application.from); err != nil {
			return err
		}
	}
	c.pendingHandshakeApplications = nil
	c.pendingHandshakeApplicationBytes = 0
	return nil
}

func (c *Conn) bufferIncompleteHandshakeApplicationLocked(content []byte, number recordNumber, from net.Addr) (bool, error) {
	if c.postHandshakeReassembly == nil || !c.postHandshakeReassembly.hasIncompleteProtected() {
		return false, nil
	}
	if len(c.pendingHandshakeApplications) >= c.maxPendingOrderingRecords() {
		return false, &ProtocolError{"too many application records buffered during protected-handshake reassembly"}
	}
	if len(content) > c.config.MaxBufferedApplicationData-c.pendingHandshakeApplicationBytes {
		return false, &ProtocolError{"buffered protected-handshake application data limit exceeded"}
	}
	c.pendingHandshakeApplications = append(c.pendingHandshakeApplications, pendingPostAuthApplication{
		number: number, content: append([]byte(nil), content...), from: from,
	})
	c.pendingHandshakeApplicationBytes += len(content)
	return true, nil
}

func (c *Conn) maxPendingOrderingRecords() int {
	limit := c.config.ReplayWindow * c.config.MaxBufferedHandshakeMessages
	if limit < 8 {
		limit = 8
	}
	return limit
}

func (c *Conn) trimRecordOrderingHistoryLocked() {
	limit := 2 * c.config.ReplayWindow
	if limit < 2 {
		limit = 2
	}
	if len(c.recentApplicationRecords) > limit {
		c.recentApplicationRecords = c.recentApplicationRecords[len(c.recentApplicationRecords)-limit:]
	}
	if len(c.protectedHandshakeRanges) > limit {
		c.protectedHandshakeRanges = c.protectedHandshakeRanges[len(c.protectedHandshakeRanges)-limit:]
	}
}

// queueEarlyApplicationData buffers authenticated epoch-1 application data
// until the handshake reaches epoch 3. Authenticated bytes beyond the
// advertised allowance terminate the connection with unexpected_message;
// they are never exposed to the application.
func (c *Conn) queueEarlyApplicationData(content []byte) error {
	if len(content) == 0 {
		return nil
	}
	c.earlyMu.Lock()
	defer c.earlyMu.Unlock()
	if !c.earlyAccepted || c.earlyDataLimit == 0 {
		return nil
	}
	if len(c.earlyReadDatagrams) >= c.config.MaxBufferedApplicationDatagrams {
		return alertError(alertUnexpectedMessage, &ProtocolError{"too many early application datagrams"})
	}
	if uint64(c.earlyReadBytes) >= uint64(c.earlyDataLimit) {
		return alertError(alertUnexpectedMessage, &ProtocolError{"early data exceeds the advertised limit"})
	}
	remaining := int(uint64(c.earlyDataLimit) - uint64(c.earlyReadBytes))
	if len(content) > remaining {
		return alertError(alertUnexpectedMessage, &ProtocolError{"early data exceeds the advertised limit"})
	}
	c.earlyReadDatagrams = append(c.earlyReadDatagrams, append([]byte(nil), content...))
	c.earlyReadBytes += len(content)
	return nil
}

func (c *Conn) promoteEarlyApplicationData() error {
	c.earlyMu.Lock()
	datagrams := c.earlyReadDatagrams
	c.earlyReadDatagrams = nil
	c.earlyReadBytes = 0
	c.earlyMu.Unlock()
	for _, datagram := range datagrams {
		if err := c.queueApplicationData(datagram, c.conn.RemoteAddr()); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) readRecords() {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)
	datagram := buffer[:]
	for {
		n, err := c.conn.Read(datagram)
		if err != nil {
			c.finishRecordReader(err)
			return
		}
		from := c.conn.RemoteAddr()
		if source, ok := c.conn.(interface{ lastReadSource() net.Addr }); ok {
			if actual := source.lastReadSource(); actual != nil {
				from = actual
			}
		}
		if err = c.dispatchDatagramFrom(datagram[:n], from); err != nil {
			c.failConnection(err)
			return
		}
	}
}

func (c *Conn) failConnection(err error) {
	if err == nil {
		return
	}
	if description, ok := outboundAlert(err); ok {
		c.sendFatalAlert(description)
	}
	c.clearTrafficSecrets(err)
	c.inputMu.Lock()
	if c.readErr == nil {
		c.readErr = err
	}
	c.inputMu.Unlock()
	c.notifyRead()
	_ = c.conn.Close()
}

func (c *Conn) clearTrafficSecrets(failure error) {
	c.dispatchMu.Lock()
	defer c.dispatchMu.Unlock()
	c.writeMu.Lock()
	c.clearReturnRoutabilityLocked()
	c.sendCipher = nil
	if c.sendingTraffic != nil {
		c.sendingTraffic.clearSecrets()
	}
	c.sendingTraffic = nil
	if c.receivingTraffic != nil {
		c.receivingTraffic.clearSecrets()
	}
	c.receivingTraffic = nil
	c.finishedACKCipher = nil
	clear(c.resumptionMasterSecret)
	c.resumptionMasterSecret = nil
	c.resumptionSuite = nil
	c.ticketFlight = nil
	c.postHandshakeReassembly = nil
	c.protectedHandshakeRanges = nil
	c.recentApplicationRecords = nil
	c.pendingHandshakeApplications = nil
	c.pendingHandshakeApplicationBytes = 0
	c.clientAuthRequestFlight = nil
	c.clientAuthResponseFlight = nil
	c.newConnectionIDFlight = nil
	c.requestCIDFlight = nil
	c.keyUpdateResponsePending = false
	if c.postHandshakeAuthState != nil && c.postHandshakeAuthState.done != nil {
		select {
		case c.postHandshakeAuthState.done <- failure:
		default:
		}
	}
	c.postHandshakeAuthState = nil
	c.writeMu.Unlock()
	c.inputMu.Lock()
	c.applicationDatagrams = nil
	c.bufferedApplicationBytes = 0
	c.inputMu.Unlock()
	c.earlyMu.Lock()
	c.earlyPending = nil
	c.earlyReadDatagrams = nil
	c.earlyReadBytes = 0
	c.earlyMu.Unlock()
	if c.receiveEpochs != nil {
		c.receiveEpochs.clear()
	}
	c.mu.Lock()
	if c.state.exporter != nil {
		c.state.exporter.clear()
	}
	c.state.exporter = nil
	c.mu.Unlock()
}

func (c *Conn) sendFatalAlert(description uint8) {
	body, err := (alertMessage{level: alertLevelFatal, description: description}).marshal()
	if err != nil {
		return
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.sendCipher == nil {
		return
	}
	wire, err := c.sendCipher.seal(recordTypeAlert, body)
	if err == nil {
		_, _ = c.conn.Write(wire)
	}
}

func (c *Conn) sendPlainFatalAlert(description uint8) {
	body, err := (alertMessage{level: alertLevelFatal, description: description}).marshal()
	if err != nil {
		return
	}
	sequence := c.plainSendSequence.Add(1) - 1
	wire, err := marshalPlainRecord(record{typ: recordTypeAlert, sequence: sequence, payload: body})
	if err == nil {
		_, _ = c.conn.Write(wire)
	}
}

func (c *Conn) dispatchDatagram(datagram []byte) error {
	return c.dispatchDatagramFrom(datagram, c.conn.RemoteAddr())
}

func (c *Conn) dispatchDatagramFrom(datagram []byte, from net.Addr) error {
	c.dispatchMu.Lock()
	defer c.dispatchMu.Unlock()
	for len(datagram) > 0 {
		content, typ, epoch, consumed, openErr := c.receiveEpochs.openInPlace(datagram)
		if openErr != nil {
			if fatalErr := protectedRecordReceiveError(openErr); fatalErr != nil {
				return fatalErr
			}
			if c.receiveEpochs.shouldRequestKeyUpdateForAuthFailures(datagram[0]) {
				if err := c.requestKeyUpdateAfterAuthFailures(); err != nil {
					return err
				}
			}
			return nil
		}
		cipher, selectErr := c.receiveEpochs.selectCipher(datagram[0])
		if selectErr != nil {
			return nil
		}
		number := recordNumber{epoch: epoch, sequence: cipher.lastOpened}
		if cipher.hasConnectionID {
			c.mu.Lock()
			if !equalBytes(c.state.LocalConnectionID, cipher.lastConnectionID) {
				c.state.LocalConnectionID = append([]byte{}, cipher.lastConnectionID...)
			}
			c.mu.Unlock()
		}
		if c.closure.ignore(number) {
			datagram = datagram[consumed:]
			continue
		}
		if cipher.hasConnectionID {
			if err := c.observeReturnRoutabilityRecord(from, consumed); err != nil {
				return err
			}
		}
		switch typ {
		case recordTypeApplicationData:
			c.writeMu.Lock()
			bufferErr := c.rememberApplicationRecordLocked(number)
			buffered := false
			if bufferErr == nil {
				buffered, bufferErr = c.bufferIncompleteHandshakeApplicationLocked(content, number, from)
			}
			if bufferErr == nil && !buffered {
				buffered, bufferErr = c.bufferPostHandshakeAuthApplicationLocked(content, number, from)
			}
			c.writeMu.Unlock()
			if bufferErr != nil {
				return bufferErr
			}
			if buffered {
				break
			}
			if err := c.queueApplicationData(content, from); err != nil {
				return err
			}
		case recordTypeAlert:
			alert, parseErr := parseAlert(content)
			if parseErr != nil {
				description, _ := protocolAlert(parseErr)
				return alertError(description, parseErr)
			}
			if alert.isUserCanceled() {
				break
			}
			if alert.isCloseNotify() {
				c.closure.receive(number)
				c.inputMu.Lock()
				c.peerReadClosed = true
				c.inputMu.Unlock()
				c.notifyRead()
				break
			}
			return AlertError(alert.description)
		case recordTypeACK:
			var scratch [1]recordNumber
			numbers, parseErr := parseACKInto(content, scratch[:0])
			if parseErr != nil {
				description, _ := protocolAlert(parseErr)
				return alertError(description, parseErr)
			}
			if parseErr = validateACKEpoch(numbers, epoch); parseErr != nil {
				return parseErr
			}
			startKeyUpdateResponse := false
			c.writeMu.Lock()
			if c.sendingTraffic != nil {
				if c.sendingTraffic.processACK(numbers) {
					c.sendCipher = c.sendingTraffic.cipher
					if c.keyUpdateResponsePending && c.sendingTraffic.canBeginKeyUpdate() {
						wire, _, beginErr := c.sendingTraffic.beginKeyUpdate(false)
						if beginErr != nil {
							c.writeMu.Unlock()
							return beginErr
						}
						if beginErr = c.writeRecord(wire); beginErr != nil {
							c.writeMu.Unlock()
							return beginErr
						}
						c.keyUpdateResponsePending = false
						startKeyUpdateResponse = true
					} else if c.keyUpdateResponsePending && c.sendingTraffic.cipher.epoch >= maxSendingEpoch {
						c.keyUpdateResponsePending = false
					}
				}
			}
			if c.ticketFlight != nil {
				c.ticketFlight.ack(numbers)
				if c.ticketFlight.complete() {
					c.observeFlightRTT(c.ticketFlight)
					c.ticketFlight = nil
				} else if err := c.retransmitPartialFlight(c.conn, c.ticketFlight); err != nil {
					c.writeMu.Unlock()
					return err
				}
			}
			if c.clientAuthRequestFlight != nil {
				c.clientAuthRequestFlight.ack(numbers)
				if c.clientAuthRequestFlight.complete() {
					c.observeFlightRTT(c.clientAuthRequestFlight)
					c.clientAuthRequestFlight = nil
				} else if err := c.retransmitPartialFlight(c.conn, c.clientAuthRequestFlight); err != nil {
					c.writeMu.Unlock()
					return err
				}
			}
			if c.clientAuthResponseFlight != nil {
				c.clientAuthResponseFlight.ack(numbers)
				if c.clientAuthResponseFlight.complete() {
					c.observeFlightRTT(c.clientAuthResponseFlight)
					c.clientAuthResponseFlight = nil
				} else if err := c.retransmitPartialFlight(c.conn, c.clientAuthResponseFlight); err != nil {
					c.writeMu.Unlock()
					return err
				}
			}
			if err := c.processCIDACKsLocked(numbers); err != nil {
				c.writeMu.Unlock()
				return err
			}
			c.writeMu.Unlock()
			if startKeyUpdateResponse {
				go c.startKeyUpdateRetransmission()
			}
		case recordTypeReturnRoutability:
			if err := c.handleReturnRoutability(content, from); err != nil {
				return err
			}
		case recordTypeHandshake:
			var fragmentScratch [1]handshakeFragment
			fragments, parseErr := parseHandshakeFragmentsViewInto(content, fragmentScratch[:0])
			if parseErr != nil {
				return parseErr
			}
			for _, fragment := range fragments {
				if fragment.typ == handshakeTypeKeyUpdate && (len(fragments) != 1 || fragment.offset != 0 || int(fragment.length) != len(fragment.body)) {
					return alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate is not aligned to a record boundary"})
				}
			}
			if epoch == 2 && c.hasCompletedPeerFlight {
				completedFlightRecord := len(fragments) > 0
				for _, fragment := range fragments {
					if fragment.messageSequence < c.completedPeerFlightStart || fragment.messageSequence > c.completedPeerFlightEnd {
						completedFlightRecord = false
						break
					}
				}
				if completedFlightRecord {
					c.writeMu.Lock()
					var ackScratch [1][]byte
					acks, _, ackErr := buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
					if ackErr == nil {
						for _, wire := range acks {
							_, _ = c.conn.Write(wire)
						}
					}
					c.writeMu.Unlock()
				}
				break
			}
			if c.receivingTraffic != nil {
				for _, fragment := range fragments {
					if fragment.typ == handshakeTypeCertificateRequest && c.isClient {
						if err := c.ackProtectedRecord(number); err != nil {
							return err
						}
						if c.postHandshakeReassembly == nil {
							c.postHandshakeReassembly = newReassemblerWithLimits(c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
						}
						body, complete, firstRecord, lastRecord, reassemblyErr := c.postHandshakeReassembly.addProtectedRecord(fragment, number)
						if reassemblyErr != nil {
							return reassemblyErr
						}
						if complete {
							if reassemblyErr = c.rememberProtectedHandshakeRange(firstRecord, lastRecord); reassemblyErr != nil {
								return reassemblyErr
							}
							if err := c.processPostHandshakeCertificateRequest(fragment.messageSequence, body); err != nil {
								return err
							}
						}
						continue
					}
					if !c.isClient && (fragment.typ == handshakeTypeCertificate || fragment.typ == handshakeTypeCompressedCertificate || fragment.typ == handshakeTypeCertificateVerify || fragment.typ == handshakeTypeFinished) {
						if err := c.processPostHandshakeAuthFragment(fragment, number); err != nil {
							return err
						}
						continue
					}
					if fragment.typ == handshakeTypeNewSessionTicket && c.isClient {
						if err := c.ackProtectedRecord(number); err != nil {
							return err
						}
						if c.postHandshakeReassembly == nil {
							c.postHandshakeReassembly = newReassemblerWithLimits(c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
						}
						body, complete, firstRecord, lastRecord, reassemblyErr := c.postHandshakeReassembly.addProtectedRecord(fragment, number)
						if reassemblyErr != nil {
							return reassemblyErr
						}
						if complete {
							if reassemblyErr = c.rememberProtectedHandshakeRange(firstRecord, lastRecord); reassemblyErr != nil {
								return reassemblyErr
							}
							if err := c.processNewSessionTicket(body); err != nil {
								return err
							}
						}
						continue
					}
					if fragment.typ == handshakeTypeNewConnectionID || fragment.typ == handshakeTypeRequestConnectionID {
						if c.postHandshakeReassembly == nil {
							c.postHandshakeReassembly = newReassemblerWithLimits(c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
						}
						body, complete, firstRecord, lastRecord, reassemblyErr := c.postHandshakeReassembly.addProtectedRecord(fragment, number)
						if reassemblyErr != nil {
							return reassemblyErr
						}
						var requestCount uint8
						respond := false
						if complete {
							if reassemblyErr = c.rememberProtectedHandshakeRange(firstRecord, lastRecord); reassemblyErr != nil {
								return reassemblyErr
							}
							if fragment.typ == handshakeTypeNewConnectionID {
								reassemblyErr = c.processNewConnectionID(fragment.messageSequence, body)
							} else {
								requestCount, respond, reassemblyErr = c.processRequestConnectionID(fragment.messageSequence, body)
							}
							if reassemblyErr != nil {
								return reassemblyErr
							}
						}
						if err := c.ackProtectedRecord(number); err != nil {
							return err
						}
						if respond {
							go c.respondToConnectionIDRequest(requestCount)
						}
						continue
					}
					if fragment.typ == handshakeTypeKeyUpdate && fragment.offset == 0 && int(fragment.length) == len(fragment.body) {
						startRetransmission := false
						var responseErr error
						c.writeMu.Lock()
						if c.postHandshakeAuthState != nil && c.postHandshakeAuthState.hasResponseEpoch {
							c.writeMu.Unlock()
							return alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate interleaved with post-handshake authentication response"})
						}
						if c.hasIncompleteProtectedHandshakeLocked() {
							c.writeMu.Unlock()
							return alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate followed an incomplete handshake message"})
						}
						message, updated, updateErr := c.receivingTraffic.processKeyUpdate(fragment.messageSequence, fragment.body)
						if updateErr == nil {
							var acks [][]byte
							var ackScratch [1][]byte
							acks, _, responseErr = buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
							if responseErr == nil {
								for _, wire := range acks {
									if responseErr = c.writeRecord(wire); responseErr != nil {
										break
									}
								}
							}
							if responseErr == nil && updated && message.requestUpdate {
								if c.sendingTraffic.canBeginKeyUpdate() {
									wire, _, beginErr := c.sendingTraffic.beginKeyUpdate(false)
									responseErr = beginErr
									if responseErr == nil {
										c.sendCipher = c.sendingTraffic.cipher
										if responseErr = c.writeRecord(wire); responseErr == nil {
											startRetransmission = true
										}
									}
								} else if c.sendingTraffic.cipher.epoch < maxSendingEpoch {
									c.keyUpdateResponsePending = true
								}
							}
						}
						c.writeMu.Unlock()
						if updateErr != nil {
							return updateErr
						}
						if responseErr != nil {
							return responseErr
						}
						if startRetransmission {
							c.startKeyUpdateRetransmission()
						}
						continue
					}
					return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake message"})
				}
			}
		}
		datagram = datagram[consumed:]
	}
	return nil
}

func (c *Conn) hasIncompleteProtectedHandshakeLocked() bool {
	if c.postHandshakeReassembly != nil && c.postHandshakeReassembly.hasIncompleteProtected() {
		return true
	}
	return c.postHandshakeAuthState != nil && c.postHandshakeAuthState.inbox.hasIncompleteProtected()
}

func (c *Conn) requestKeyUpdateAfterAuthFailures() error {
	c.writeMu.Lock()
	if !c.sendingTraffic.canBeginKeyUpdate() {
		c.writeMu.Unlock()
		return nil
	}
	wire, _, err := c.sendingTraffic.beginKeyUpdate(true)
	if err == nil {
		err = c.writeRecord(wire)
	}
	c.writeMu.Unlock()
	if err != nil {
		return err
	}
	go c.startKeyUpdateRetransmission()
	return nil
}

// ReadDatagram reads and consumes one authenticated DTLS Application Data
// record. It performs the handshake first if necessary. The returned n is the
// number of plaintext bytes copied into p.
//
// If p is too small, the unread remainder is discarded,
// DatagramInfo.FullLength reports the original length, and
// DatagramInfo.Truncated is true. A subsequent call reads the next record, not
// the remainder. Empty records are returned with n and FullLength both zero.
// A peer close_notify is returned as io.EOF after all earlier deliverable
// records have been consumed.
func (c *Conn) ReadDatagram(p []byte) (int, DatagramInfo, error) {
	if err := c.Handshake(); err != nil {
		return 0, DatagramInfo{}, err
	}
	c.readMu.Lock()
	defer c.readMu.Unlock()
	for {
		c.inputMu.Lock()
		if len(c.applicationDatagrams) > 0 {
			datagram := c.applicationDatagrams[0]
			c.applicationDatagrams = c.applicationDatagrams[1:]
			c.bufferedApplicationBytes -= len(datagram.payload)
			n := copy(p, datagram.payload)
			c.inputMu.Unlock()
			return n, DatagramInfo{Source: datagram.from, FullLength: len(datagram.payload), Truncated: n < len(datagram.payload)}, nil
		}
		if c.readErr != nil {
			err := c.readErr
			c.readErr = nil
			c.inputMu.Unlock()
			return 0, DatagramInfo{}, err
		}
		if c.peerReadClosed {
			c.inputMu.Unlock()
			return 0, DatagramInfo{}, io.EOF
		}
		c.inputMu.Unlock()
		c.startRecordReader()
		<-c.readNotify
	}
}

// WriteDatagram sends p as exactly one DTLS Application Data record to the
// association's authenticated peer. It performs the handshake first if
// necessary. Application data is not internally fragmented, retransmitted, or
// reordered. A nil or empty p sends a valid empty application datagram.
//
// By default, if p exceeds the current path MTU or the DTLS record-content
// limit, WriteDatagram returns an error matching [ErrDatagramTooLarge] and n
// == 0, without a partial record on the wire. [Config.IgnorePathMTU] skips the
// library's PMTU check but not the record-content limit. The transport may
// still reject the complete record with ErrDatagramTooLarge. Otherwise n is
// len(p) if the complete record was handed to the underlying transport.
func (c *Conn) WriteDatagram(p []byte) (int, error) {
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	addr := c.RemoteAddr()
	if addr == nil {
		return 0, &net.OpError{Op: "write", Net: "dtls", Err: errors.New("missing destination address")}
	}
	c.writeMu.Lock()
	if c.sendCipher == nil {
		c.writeMu.Unlock()
		return 0, &ProtocolError{"application write keys are not installed"}
	}
	for {
		maximum := c.maxApplicationDatagramLocked()
		if maximum < 0 || len(p) > maximum {
			c.writeMu.Unlock()
			return 0, datagramTooLargeError(addr)
		}
		wire, err := c.sendCipher.seal(recordTypeApplicationData, p)
		if err != nil {
			c.writeMu.Unlock()
			return 0, err
		}
		if err = c.writeRecord(wire); err != nil {
			if !c.config.IgnorePathMTU && isMessageTooLong(err) {
				if _, reduced := c.reducePathMTU(); reduced {
					continue
				}
			}
			c.writeMu.Unlock()
			return 0, err
		}
		startUpdate, err := c.maybeStartAutomaticKeyUpdateLocked()
		c.writeMu.Unlock()
		if err != nil {
			return len(p), err
		}
		if startUpdate {
			go c.startKeyUpdateRetransmission()
		}
		return len(p), nil
	}
}

func (c *Conn) maxApplicationDatagramLocked() int {
	return c.maxApplicationDatagramForCipher(c.sendCipher)
}

func (c *Conn) maxApplicationDatagramForCipher(cipher *recordCipher) int {
	maximum := cipher.maxContent()
	if c.config.IgnorePathMTU {
		return maximum
	}
	pathMaximum := c.currentMTU() - cipher.headerLen16() - cipher.aead.Overhead() - 1
	if pathMaximum < maximum {
		maximum = pathMaximum
	}
	return maximum
}

func (c *Conn) maybeStartAutomaticKeyUpdateLocked() (bool, error) {
	if c.sendingTraffic == nil || !c.sendingTraffic.update.canUseNewKeys() || c.sendCipher.epoch >= maxSendingEpoch {
		return false, nil
	}
	margin := uint64(1024)
	if c.sendCipher.recordLimit/4 < margin {
		margin = c.sendCipher.recordLimit / 4
	}
	if margin < 1 {
		margin = 1
	}
	if c.sendCipher.nextSequence < c.sendCipher.recordLimit-margin {
		return false, nil
	}
	wire, _, err := c.sendingTraffic.beginKeyUpdate(false)
	if err != nil {
		return false, err
	}
	if err = c.writeRecord(wire); err != nil {
		return false, err
	}
	return true, nil
}

// WriteEarlyData attempts to send p as one client 0-RTT Application Data
// record and completes the handshake. It can be called at most once and only
// on a client created with a usable cached session whose ticket permits early
// data. A nil or empty p is a no-op.
//
// The method returns ErrEarlyDataUnavailable when no eligible early-data
// session exists, and ErrEarlyDataRejected when the record was sent but the
// server completed the handshake without accepting it. In both cases the
// caller may use the established connection for 1-RTT data. Retrying p is an
// application decision because 0-RTT data is replayable. Oversized data
// returns [ErrDatagramTooLarge] without a partial record.
// [Config.IgnorePathMTU] skips the library's PMTU check but not the DTLS
// record-content or ticket limits; the transport may still reject the complete
// record as too large.
func (c *Conn) WriteEarlyData(p []byte) (int, error) {
	if !c.isClient {
		return 0, ErrEarlyDataUnavailable
	}
	if len(p) == 0 {
		return 0, nil
	}
	c.earlyMu.Lock()
	if c.earlySignaled {
		c.earlyMu.Unlock()
		return 0, ErrEarlyDataUnavailable
	}
	c.earlyPending = append(c.earlyPending, p...)
	c.earlyMu.Unlock()
	if err := c.Handshake(); err != nil {
		return 0, err
	}
	c.earlyMu.Lock()
	sent, rejected := c.earlySent, c.earlyRejected
	c.earlySignaled = true
	c.earlyMu.Unlock()
	if !sent {
		return 0, ErrEarlyDataUnavailable
	}
	if rejected {
		return 0, ErrEarlyDataRejected
	}
	return len(p), nil
}

// Close sends close_notify when application sending keys are available,
// clears retained traffic, resumption, and exporter secrets, stops background
// protocol work, and closes the underlying transport. It does not wait for the
// peer to acknowledge close_notify.
func (c *Conn) Close() error {
	c.writeMu.Lock()
	if c.sendCipher != nil {
		if body, err := (alertMessage{level: alertLevelWarning, description: alertCloseNotify}).marshal(); err == nil {
			if wire, sealErr := c.sendCipher.seal(recordTypeAlert, body); sealErr == nil {
				_, _ = c.conn.Write(wire)
			}
		}
	}
	c.writeMu.Unlock()
	c.clearTrafficSecrets(net.ErrClosed)
	c.readerMu.Lock()
	c.readerClosed = true
	c.readerMu.Unlock()
	c.initInput()
	c.inputMu.Lock()
	if c.readErr == nil {
		c.readErr = net.ErrClosed
	}
	c.inputMu.Unlock()
	c.notifyRead()
	return c.conn.Close()
}

// LocalAddr returns the local network address of the underlying transport.
func (c *Conn) LocalAddr() net.Addr { return c.conn.LocalAddr() }

// RemoteAddr returns the current destination address. When Connection ID and
// RFC 9853 Return Routability Check are negotiated on a Listener association,
// it changes only after the new path has completed validation.
func (c *Conn) RemoteAddr() net.Addr { return c.conn.RemoteAddr() }

// SetDeadline sets both read and write deadlines on the underlying transport.
// A zero value disables the deadlines. The deadline applies to currently
// blocked and future I/O, following net.Conn semantics. An initial handshake
// temporarily installs its own deadline and clears transport deadlines when it
// finishes.
func (c *Conn) SetDeadline(t time.Time) error { return c.conn.SetDeadline(t) }

// SetReadDeadline sets the deadline for future and currently blocked reads
// from the underlying transport. A zero value disables the deadline.
func (c *Conn) SetReadDeadline(t time.Time) error { return c.conn.SetReadDeadline(t) }

// SetWriteDeadline sets the deadline for future and currently blocked writes
// to the underlying transport. A zero value disables the deadline.
func (c *Conn) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }

// ConnectionState returns a snapshot of the current negotiated state. Before
// handshake completion its fields have zero values. Post-handshake client
// authentication and Connection ID changes are reflected in later snapshots.
func (c *Conn) ConnectionState() ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	state := c.state
	state.LocalConnectionID = append([]byte(nil), state.LocalConnectionID...)
	state.PeerConnectionID = append([]byte(nil), state.PeerConnectionID...)
	return state
}

// SendKeyUpdate reliably sends a post-handshake KeyUpdate. Application writes
// continue with the old sending epoch until the KeyUpdate record is
// acknowledged, as required by RFC 9147 section 8. If requestPeer is true, the
// peer is asked to update its sending keys as well.
//
// Only one locally initiated KeyUpdate may await acknowledgement at a time.
// The package also initiates updates automatically before the negotiated AEAD
// usage limit is reached.
func (c *Conn) SendKeyUpdate(requestPeer bool) error {
	if err := c.Handshake(); err != nil {
		return err
	}
	c.writeMu.Lock()
	if c.sendingTraffic == nil {
		c.writeMu.Unlock()
		return &ProtocolError{"application traffic state is not installed"}
	}
	wire, _, err := c.sendingTraffic.beginKeyUpdate(requestPeer)
	if err != nil {
		c.writeMu.Unlock()
		return err
	}
	err = c.writeRecord(wire)
	c.writeMu.Unlock()
	if err == nil {
		c.startKeyUpdateRetransmission()
	}
	return err
}

func (c *Conn) startKeyUpdateRetransmission() {
	go func() {
		interval := c.config.FlightInterval
		if interval <= 0 {
			interval = time.Second
		}
		max := c.config.MaxFlightInterval
		if max < interval {
			max = interval
		}
		timer := time.NewTimer(interval)
		defer timer.Stop()
		for {
			<-timer.C
			c.readerMu.Lock()
			closed := c.readerClosed
			c.readerMu.Unlock()
			if closed {
				return
			}
			c.writeMu.Lock()
			if c.sendingTraffic == nil || c.sendingTraffic.update.canUseNewKeys() {
				c.writeMu.Unlock()
				return
			}
			wire, _, err := c.sendingTraffic.retransmitKeyUpdate()
			if err == nil {
				err = c.writeRecord(wire)
			}
			c.writeMu.Unlock()
			if err != nil {
				c.failConnection(err)
				return
			}
			if interval < max {
				interval *= 2
				if interval > max {
					interval = max
				}
			}
			timer.Reset(interval)
		}
	}()
}

func (c *Conn) installApplicationKeys(suite *cipherSuite, clientSecret, serverSecret []byte) error {
	return c.installApplicationKeysAt(suite, clientSecret, serverSecret, 0)
}
func (c *Conn) installApplicationKeysAt(suite *cipherSuite, clientSecret, serverSecret []byte, messageSequence uint16) error {
	var sendSecret, receiveSecret []byte
	if c.isClient {
		sendSecret, receiveSecret = clientSecret, serverSecret
	} else {
		sendSecret, receiveSecret = serverSecret, clientSecret
	}
	hashSize := suite.hash.Size()
	secretStorage := make([]byte, 2*hashSize)
	ownedSendSecret := secretStorage[:hashSize:hashSize]
	ownedReceiveSecret := secretStorage[hashSize : 2*hashSize : 2*hashSize]
	copy(ownedSendSecret, sendSecret)
	copy(ownedReceiveSecret, receiveSecret)
	sending, err := newSendingTrafficWithOwnedSecret(suite, ownedSendSecret, 3, messageSequence, c.config.ReplayWindow)
	if err != nil {
		clear(secretStorage)
		return err
	}
	receiving, err := newReceivingTrafficWithOwnedSecret(suite, ownedReceiveSecret, 3, c.config.ReplayWindow)
	if err != nil {
		clear(secretStorage)
		return err
	}
	sending.cipher.setPlaintextLimit(c.peerRecordSizeLimit)
	receiving.setPlaintextLimit(c.localRecordSizeLimit)
	if c.connectionIDNegotiated {
		if err = sending.setConnectionID(c.sendConnectionID); err != nil {
			sending.clearSecrets()
			receiving.clearSecrets()
			return err
		}
		if err = receiving.setConnectionID(c.receiveConnectionID); err != nil {
			sending.clearSecrets()
			receiving.clearSecrets()
			return err
		}
	}
	c.sendingTraffic = sending
	c.receivingTraffic = receiving
	c.sendCipher = sending.cipher
	c.receiveEpochs = receiving.epochs
	return nil
}
