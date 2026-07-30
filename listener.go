package dtls13

import (
	"bytes"
	"errors"
	"net"
	"sync"
	"time"
)

// Listen creates and owns a UDP socket on address, then returns a DTLS
// association listener. Network must be "udp", "udp4", or "udp6". A nil
// config selects defaults. Close releases the socket and all associations
// demultiplexed from it.
func Listen(network, address string, config *Config) (*Listener, error) {
	if network != "udp" && network != "udp4" && network != "udp6" {
		return nil, &ConfigError{"Listen network must be udp, udp4, or udp6"}
	}
	normalized, err := config.normalized()
	if err != nil {
		return nil, err
	}
	if err = ensureSessionTicketKey(normalized); err != nil {
		return nil, err
	}
	if err = ensureCookieProtector(normalized); err != nil {
		return nil, err
	}
	packetConn, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}
	return newListener(packetConn, normalized, nil), nil
}

// NewListener creates a DTLS association listener over inner. Ownership of
// inner is transferred to the Listener; Close closes it. inner must preserve
// datagram boundaries and provide source addresses, normally by being a UDP
// net.PacketConn.
//
// NewListener cannot return an error for compatibility with construction over
// an existing transport. A nil transport and configuration initialization
// errors are retained and returned by Accept.
func NewListener(inner net.PacketConn, config *Config) *Listener {
	normalized, err := config.normalized()
	if err == nil {
		err = ensureSessionTicketKey(normalized)
	}
	if err == nil {
		err = ensureCookieProtector(normalized)
	}
	return newListener(inner, normalized, err)
}

// Listener demultiplexes DTLS associations from one packet transport. It
// intentionally does not implement [net.Listener] because [Conn] exposes a
// connected datagram API rather than a byte stream.
//
// Accept returns a Conn after receiving an initial datagram for a new peer;
// the DTLS handshake and peer authentication occur when Handshake or the first
// datagram I/O method is called. Listener bounds pending state and per-peer
// queues according to Config.
type Listener struct {
	conn        net.PacketConn
	config      *Config
	configErr   error
	accept      chan *packetSession
	done        chan struct{}
	closeOnce   sync.Once
	mu          sync.Mutex
	closed      bool
	readErr     error
	sessions    map[string]*packetSession
	pending     map[string]*packetSession
	cidSessions map[string]*packetSession
	cidGenMu    sync.Mutex
}

type packetListener = Listener

func newListener(conn net.PacketConn, config *Config, configErr error) *packetListener {
	capacity := 1
	if config != nil && config.MaxPendingConnections > 0 {
		capacity = config.MaxPendingConnections
	}
	l := &packetListener{
		conn: conn, config: config, configErr: configErr,
		accept: make(chan *packetSession, capacity), done: make(chan struct{}),
		sessions: make(map[string]*packetSession), pending: make(map[string]*packetSession), cidSessions: make(map[string]*packetSession),
	}
	if conn == nil && l.configErr == nil {
		l.configErr = &ConfigError{"nil packet connection"}
	}
	if l.configErr == nil {
		go l.readLoop()
	}
	return l
}

func sessionKey(address net.Addr) string {
	return address.Network() + "\x00" + address.String()
}

func (l *packetListener) readLoop() {
	buffer := make([]byte, 65535)
	for {
		n, address, err := l.conn.ReadFrom(buffer)
		if err != nil {
			l.shutdown(err, false)
			return
		}
		if n == 0 || address == nil {
			continue
		}
		key := sessionKey(address)
		l.mu.Lock()
		var targets []*packetSession
		if isUnifiedCIDRecord(buffer[:n]) {
			session := l.sessionForCIDLocked(buffer[:n])
			if session != nil {
				targets = append(targets, session)
			} else {
				// An empty negotiated CID still sets C=1 but has no bytes to
				// index globally. Route it by the original tuple; record
				// authentication remains the final association check.
				for _, candidate := range []*packetSession{l.sessions[key], l.pending[key]} {
					if candidate != nil && candidate.acceptsEmptyConnectionID() {
						targets = append(targets, candidate)
					}
				}
			}
		} else {
			active := l.sessions[key]
			candidate := l.pending[key]
			initial := isInitialClientHelloDatagram(buffer[:n], l.config.MaxHandshakeMessage)
			if active == nil && candidate == nil && !l.closed && len(l.sessions)+len(l.pending) < l.config.MaxPendingConnections && initial {
				active = newPacketSession(l.conn, address, l.config.MaxSessionQueueDatagrams, l.removeSession, l.rebindSession, l.registerSessionCIDs, l.unregisterSessionCIDs, l.connectionValidated)
				l.sessions[key] = active
				select {
				case l.accept <- active:
				default:
					delete(l.sessions, key)
					active.closeWithoutCallback()
					active = nil
				}
			}
			if active != nil && active.validated && candidate == nil && !l.closed && len(l.sessions)+len(l.pending) < l.config.MaxPendingConnections && initial {
				candidate = newPacketSession(l.conn, address, l.config.MaxSessionQueueDatagrams, l.removeSession, l.rebindSession, l.registerSessionCIDs, l.unregisterSessionCIDs, l.connectionValidated)
				l.pending[key] = candidate
				select {
				case l.accept <- candidate:
				default:
					delete(l.pending, key)
					candidate.closeWithoutCallback()
					candidate = nil
				}
			}
			if active != nil {
				targets = append(targets, active)
			}
			if candidate != nil {
				targets = append(targets, candidate)
			}
		}
		l.mu.Unlock()
		for _, session := range targets {
			session.enqueue(buffer[:n], address)
		}
	}
}

func isUnifiedRecord(datagram []byte) bool {
	return len(datagram) > 0 && datagram[0]&0xe0 == unifiedHeaderFixed
}

func isUnifiedCIDRecord(datagram []byte) bool {
	return isUnifiedRecord(datagram) && datagram[0]&unifiedHeaderCID != 0
}

func isInitialClientHelloDatagram(datagram []byte, maxMessage int) bool {
	if len(datagram) < plainRecordHeaderLen || datagram[0] != recordTypeHandshake || datagram[3] != 0 || datagram[4] != 0 {
		return false
	}
	recordLength := int(datagram[11])<<8 | int(datagram[12])
	if recordLength < 12 || recordLength > 1<<14 || len(datagram) < plainRecordHeaderLen+recordLength {
		return false
	}
	fragment := datagram[plainRecordHeaderLen : plainRecordHeaderLen+recordLength]
	messageLength := int(fragment[1])<<16 | int(fragment[2])<<8 | int(fragment[3])
	messageSequence := int(fragment[4])<<8 | int(fragment[5])
	fragmentOffset := int(fragment[6])<<16 | int(fragment[7])<<8 | int(fragment[8])
	fragmentLength := int(fragment[9])<<16 | int(fragment[10])<<8 | int(fragment[11])
	return fragment[0] == handshakeTypeClientHello && messageSequence == 0 && messageLength > 0 && messageLength <= maxMessage &&
		fragmentLength > 0 && fragmentOffset <= messageLength-fragmentLength && len(fragment) >= 12+fragmentLength
}

func (l *packetListener) sessionForCIDLocked(datagram []byte) *packetSession {
	if len(datagram) < 2 || datagram[0]&unifiedHeaderCID == 0 {
		return nil
	}
	for cid, session := range l.cidSessions {
		if len(datagram) >= 1+len(cid) && bytes.Equal(datagram[1:1+len(cid)], []byte(cid)) {
			return session
		}
	}
	return nil
}

func (l *packetListener) registerSessionCID(session *packetSession, cid []byte) error {
	session.mu.Lock()
	session.connectionIDNegotiated = true
	session.localCID = append([]byte(nil), cid...)
	session.mu.Unlock()
	if len(cid) == 0 {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for existing, owner := range l.cidSessions {
		if bytes.Equal(cid, []byte(existing)) {
			if owner == session {
				return nil
			}
			return &ConfigError{"Listener connection IDs must be unique and prefix-free"}
		}
		if bytes.HasPrefix(cid, []byte(existing)) || bytes.HasPrefix([]byte(existing), cid) {
			return &ConfigError{"Listener connection IDs must be unique and prefix-free"}
		}
	}
	l.cidSessions[string(cid)] = session
	return nil
}

func (l *packetListener) registerSessionCIDs(session *packetSession, cids [][]byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	merged := make([][]byte, 0, len(l.cidSessions)+len(cids))
	for existing := range l.cidSessions {
		merged = append(merged, []byte(existing))
	}
	if _, err := mergeConnectionIDs(merged, cids); err != nil {
		return &ConfigError{"Listener connection IDs must be unique and prefix-free"}
	}
	for _, cid := range cids {
		if len(cid) == 0 {
			continue
		}
		if owner := l.cidSessions[string(cid)]; owner != nil && owner != session {
			return &ConfigError{"Listener connection IDs must be unique and prefix-free"}
		}
	}
	for _, cid := range cids {
		if len(cid) == 0 {
			continue
		}
		l.cidSessions[string(cid)] = session
	}
	return nil
}

func (l *packetListener) unregisterSessionCIDs(session *packetSession, cids [][]byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, cid := range cids {
		key := string(cid)
		if l.cidSessions[key] == session {
			delete(l.cidSessions, key)
		}
	}
}

func (l *packetListener) removeSession(session *packetSession) {
	l.mu.Lock()
	for key, candidate := range l.sessions {
		if candidate == session {
			delete(l.sessions, key)
		}
	}
	for key, candidate := range l.pending {
		if candidate == session {
			delete(l.pending, key)
		}
	}
	for cid, candidate := range l.cidSessions {
		if candidate == session {
			delete(l.cidSessions, cid)
		}
	}
	l.mu.Unlock()
}

// connectionValidated promotes a pending association only after its Finished
// has been authenticated. The previous association remains usable until this
// point, as required by RFC 9147 section 5.11.
func (l *packetListener) connectionValidated(session *packetSession) {
	if session == nil || session.RemoteAddr() == nil {
		return
	}
	key := sessionKey(session.RemoteAddr())
	var old *packetSession
	l.mu.Lock()
	session.validated = true
	if l.pending[key] == session {
		old = l.sessions[key]
		l.sessions[key] = session
		delete(l.pending, key)
	}
	l.mu.Unlock()
	if old != nil && old != session {
		_ = old.Close()
	}
}

func (l *packetListener) rebindSession(session *packetSession, newAddress net.Addr) {
	if session == nil || newAddress == nil {
		return
	}
	l.mu.Lock()
	session.mu.Lock()
	oldAddress := session.remote
	if oldAddress == nil || sessionKey(oldAddress) == sessionKey(newAddress) {
		session.mu.Unlock()
		l.mu.Unlock()
		return
	}
	oldKey, newKey := sessionKey(oldAddress), sessionKey(newAddress)
	if l.sessions[oldKey] == session {
		delete(l.sessions, oldKey)
	}
	if displaced := l.sessions[newKey]; displaced != nil && displaced != session {
		delete(l.sessions, newKey)
		for cid, candidate := range l.cidSessions {
			if candidate == displaced {
				delete(l.cidSessions, cid)
			}
		}
		displaced.closeWithoutCallback()
	}
	l.sessions[newKey] = session
	session.remote = newAddress
	session.mu.Unlock()
	l.mu.Unlock()
}

// Accept waits for the next candidate DTLS association and returns it as a
// server-side Conn. The returned connection has not necessarily completed its
// handshake; call Conn.Handshake to authenticate it before handing it to code
// that assumes an authenticated peer. Closing the returned Conn removes its
// association state from the Listener.
//
// Accept returns a retained NewListener configuration error, the underlying
// packet read error, or net.ErrClosed after the Listener is closed.
func (l *Listener) Accept() (*Conn, error) {
	if l.configErr != nil {
		return nil, l.configErr
	}
	l.mu.Lock()
	closed, readErr := l.closed, l.readErr
	l.mu.Unlock()
	if closed {
		if readErr != nil {
			return nil, readErr
		}
		return nil, net.ErrClosed
	}
	select {
	case session := <-l.accept:
		l.mu.Lock()
		closed, readErr = l.closed, l.readErr
		l.mu.Unlock()
		if closed {
			_ = session.Close()
			if readErr != nil {
				return nil, readErr
			}
			return nil, net.ErrClosed
		}
		config := l.config.Clone()
		if config.GetConnectionID != nil {
			l.cidGenMu.Lock()
			cid, err := config.GetConnectionID()
			l.cidGenMu.Unlock()
			if err != nil {
				_ = session.Close()
				return nil, err
			}
			if len(cid) > 255 {
				_ = session.Close()
				return nil, &ConfigError{"GetConnectionID returned more than 255 bytes"}
			}
			if cid == nil {
				config.ConnectionID = nil
			} else {
				config.ConnectionID = append([]byte{}, cid...)
			}
		}
		if config.ConnectionID != nil {
			if err := l.registerSessionCID(session, config.ConnectionID); err != nil {
				_ = session.Close()
				return nil, err
			}
		}
		return Server(session, config), nil
	case <-l.done:
		l.mu.Lock()
		err := l.readErr
		l.mu.Unlock()
		if err != nil {
			return nil, err
		}
		return nil, net.ErrClosed
	}
}

func (l *packetListener) shutdown(err error, closePacketConn bool) {
	l.closeOnce.Do(func() {
		l.mu.Lock()
		l.closed = true
		l.readErr = err
		sessions := make([]*packetSession, 0, len(l.sessions)+len(l.pending))
		for _, session := range l.sessions {
			sessions = append(sessions, session)
		}
		for _, session := range l.pending {
			sessions = append(sessions, session)
		}
		l.sessions = make(map[string]*packetSession)
		l.pending = make(map[string]*packetSession)
		l.cidSessions = make(map[string]*packetSession)
		l.mu.Unlock()
		close(l.done)
		for _, session := range sessions {
			session.closeWithoutCallback()
		}
		if closePacketConn && l.conn != nil {
			_ = l.conn.Close()
		}
	})
}

// Close closes the packet transport, unblocks Accept, and closes all
// associations owned by the Listener. Calls blocked in Accept return
// net.ErrClosed unless an earlier packet-read error is available.
func (l *Listener) Close() error {
	if l.conn == nil {
		l.shutdown(net.ErrClosed, false)
		return nil
	}
	err := l.conn.Close()
	l.shutdown(net.ErrClosed, false)
	return err
}

// Addr returns the packet transport's local network address. It returns nil
// when the Listener has no underlying transport.
func (l *Listener) Addr() net.Addr {
	if l.conn == nil {
		return nil
	}
	return l.conn.LocalAddr()
}

type packetSession struct {
	conn                   net.PacketConn
	remote                 net.Addr
	in                     chan packetDatagram
	done                   chan struct{}
	deadline               chan struct{}
	closeOnce              sync.Once
	onClose                func(*packetSession)
	onRebind               func(*packetSession, net.Addr)
	onNewCIDs              func(*packetSession, [][]byte) error
	onRemoveCIDs           func(*packetSession, [][]byte)
	onValidated            func(*packetSession)
	mu                     sync.Mutex
	readBy                 time.Time
	writeBy                time.Time
	lastReadFrom           net.Addr
	localCID               []byte
	connectionIDNegotiated bool
	validated              bool
}

type packetDatagram struct {
	payload []byte
	from    net.Addr
}

func newPacketSession(conn net.PacketConn, remote net.Addr, queueSize int, onClose func(*packetSession), onRebind func(*packetSession, net.Addr), onNewCIDs func(*packetSession, [][]byte) error, onRemoveCIDs func(*packetSession, [][]byte), onValidated func(*packetSession)) *packetSession {
	return &packetSession{
		conn: conn, remote: remote, in: make(chan packetDatagram, queueSize), done: make(chan struct{}),
		deadline: make(chan struct{}, 1), onClose: onClose, onRebind: onRebind, onNewCIDs: onNewCIDs, onRemoveCIDs: onRemoveCIDs, onValidated: onValidated,
	}
}

func (s *packetSession) enqueue(datagram []byte, from net.Addr) {
	copyOfDatagram := append([]byte(nil), datagram...)
	select {
	case <-s.done:
	case s.in <- packetDatagram{payload: copyOfDatagram, from: from}:
	default:
	}
}

func (s *packetSession) Read(p []byte) (int, error) {
	for {
		s.mu.Lock()
		deadline := s.readBy
		s.mu.Unlock()
		var timer *time.Timer
		var timeout <-chan time.Time
		if !deadline.IsZero() {
			delay := time.Until(deadline)
			if delay < 0 {
				delay = 0
			}
			timer = time.NewTimer(delay)
			timeout = timer.C
		}
		select {
		case datagram := <-s.in:
			if timer != nil {
				timer.Stop()
			}
			s.mu.Lock()
			s.lastReadFrom = datagram.from
			s.mu.Unlock()
			return copy(p, datagram.payload), nil
		case <-s.done:
			if timer != nil {
				timer.Stop()
			}
			return 0, net.ErrClosed
		case <-s.deadline:
			if timer != nil {
				timer.Stop()
			}
			continue
		case <-timeout:
			return 0, &sessionTimeoutError{}
		}
	}
}

func (s *packetSession) Write(p []byte) (int, error) {
	select {
	case <-s.done:
		return 0, net.ErrClosed
	default:
	}
	s.mu.Lock()
	deadline := s.writeBy
	remote := s.remote
	s.mu.Unlock()
	if !deadline.IsZero() && !time.Now().Before(deadline) {
		return 0, &sessionTimeoutError{}
	}
	return s.conn.WriteTo(p, remote)
}

func (s *packetSession) WriteTo(p []byte, addr net.Addr) (int, error) {
	select {
	case <-s.done:
		return 0, net.ErrClosed
	default:
	}
	if addr == nil {
		return 0, &net.OpError{Op: "write", Net: "dtls", Err: errors.New("missing destination address")}
	}
	s.mu.Lock()
	deadline := s.writeBy
	s.mu.Unlock()
	if !deadline.IsZero() && !time.Now().Before(deadline) {
		return 0, &sessionTimeoutError{}
	}
	return s.conn.WriteTo(p, addr)
}

func (s *packetSession) closeWithoutCallback() {
	s.closeOnce.Do(func() { close(s.done) })
}

func (s *packetSession) Close() error {
	closed := false
	s.closeOnce.Do(func() {
		closed = true
		close(s.done)
		if s.onClose != nil {
			s.onClose(s)
		}
	})
	if !closed {
		return net.ErrClosed
	}
	return nil
}

func (s *packetSession) LocalAddr() net.Addr { return s.conn.LocalAddr() }
func (s *packetSession) RemoteAddr() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remote
}

func (s *packetSession) lastReadSource() net.Addr {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastReadFrom
}

func (s *packetSession) rebindRemote(newAddress net.Addr) {
	if newAddress == nil {
		return
	}
	if s.onRebind != nil {
		s.onRebind(s, newAddress)
		return
	}
	s.mu.Lock()
	if !sameNetworkAddress(s.remote, newAddress) {
		s.remote = newAddress
	}
	s.mu.Unlock()
}

func (s *packetSession) handshakeValidated() {
	if s.onValidated != nil {
		s.onValidated(s)
	}
}

func (s *packetSession) registerConnectionIDs(cids [][]byte) error {
	if s.onNewCIDs == nil {
		return nil
	}
	return s.onNewCIDs(s, cids)
}

func (s *packetSession) unregisterConnectionIDs(cids [][]byte) {
	if s.onRemoveCIDs != nil {
		s.onRemoveCIDs(s, cids)
	}
}

func (s *packetSession) acceptsEmptyConnectionID() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.connectionIDNegotiated && len(s.localCID) == 0
}

func (s *packetSession) signalDeadline() {
	select {
	case s.deadline <- struct{}{}:
	default:
	}
}

func (s *packetSession) SetDeadline(t time.Time) error {
	s.mu.Lock()
	s.readBy, s.writeBy = t, t
	s.mu.Unlock()
	s.signalDeadline()
	return nil
}

func (s *packetSession) SetReadDeadline(t time.Time) error {
	s.mu.Lock()
	s.readBy = t
	s.mu.Unlock()
	s.signalDeadline()
	return nil
}

func (s *packetSession) SetWriteDeadline(t time.Time) error {
	s.mu.Lock()
	s.writeBy = t
	s.mu.Unlock()
	return nil
}

type sessionTimeoutError struct{}

func (*sessionTimeoutError) Error() string   { return "i/o timeout" }
func (*sessionTimeoutError) Timeout() bool   { return true }
func (*sessionTimeoutError) Temporary() bool { return true }

var _ net.Conn = (*packetSession)(nil)
