package dtls13

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"sync"
)

var (
	certificateZlibWriters          sync.Pool
	certificateZlibReaders          sync.Pool
	certificateDecompressionBuffers sync.Pool
	configStateMu                   sync.Mutex
)

const maxPooledCertificate = 64 << 10

type certificateCompressionCache struct {
	mu         sync.Mutex
	key        [sha256.Size]byte
	compressed []byte
	set        bool
}

func (c *Config) ensureCertificateCompressionCache() {
	configStateMu.Lock()
	if c.state == nil {
		c.state = new(configState)
	}
	if c.state.certificateCompression == nil {
		c.state.certificateCompression = new(certificateCompressionCache)
	}
	configStateMu.Unlock()
}

func ensureConfigState(c *Config) *configState {
	configStateMu.Lock()
	if c.state == nil {
		c.state = new(configState)
	}
	state := c.state
	configStateMu.Unlock()
	return state
}

func (c *Config) certificateCompressionCache() *certificateCompressionCache {
	if c.state == nil {
		return nil
	}
	return c.state.certificateCompression
}

func (c *certificateCompressionCache) compress(certificate []byte) ([]byte, error) {
	key := sha256.Sum256(certificate)
	c.mu.Lock()
	if c.set && c.key == key {
		compressed := c.compressed
		c.mu.Unlock()
		return compressed, nil
	}
	cacheable := !c.set
	c.mu.Unlock()
	compressed, err := marshalCompressedCertificate(certificate)
	if err != nil || !cacheable {
		return compressed, err
	}
	c.mu.Lock()
	if !c.set {
		c.key = key
		c.compressed = compressed
		c.set = true
	} else if c.key == key {
		compressed = c.compressed
	}
	c.mu.Unlock()
	return compressed, nil
}

const (
	handshakeTypeCompressedCertificate uint8  = 25
	certificateCompressionZlib         uint16 = 1
)

type certificateCompressionAlgorithms string

var certificateCompressionZlibOffer = certificateCompressionAlgorithms("\x02\x00\x01")

func marshalCertificateCompressionAlgorithms(algorithms *certificateCompressionAlgorithms) ([]byte, error) {
	if algorithms == nil {
		return nil, &ProtocolError{"invalid certificate compression algorithm list"}
	}
	raw := string(*algorithms)
	if len(raw) < 3 || len(raw) > 255 || int(raw[0]) != len(raw)-1 || raw[0]%2 != 0 {
		return nil, &ProtocolError{"invalid certificate compression algorithm list"}
	}
	return []byte(raw), nil
}

func parseCertificateCompressionAlgorithms(raw []byte) (*certificateCompressionAlgorithms, error) {
	if len(raw) < 3 || int(raw[0]) != len(raw)-1 || raw[0]%2 != 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid compress_certificate extension"})
	}
	if len(raw) == 3 && raw[0] == 2 && raw[1] == 0 && raw[2] == byte(certificateCompressionZlib) {
		return &certificateCompressionZlibOffer, nil
	}
	algorithms := certificateCompressionAlgorithms(string(raw))
	return &algorithms, nil
}

func supportsCertificateCompression(algorithms *certificateCompressionAlgorithms, algorithm uint16) bool {
	if algorithms == nil {
		return false
	}
	raw := string(*algorithms)
	for offset := 1; offset < len(raw); offset += 2 {
		if uint16(raw[offset])<<8|uint16(raw[offset+1]) == algorithm {
			return true
		}
	}
	return false
}

func (h *clientHello) certificateCompressionAlgorithms() *certificateCompressionAlgorithms {
	if h == nil || !h.certificateCompressionOffered {
		return nil
	}
	if raw, ok := h.unknownExtensions[extCompressCertificate]; ok {
		algorithms, _ := parseCertificateCompressionAlgorithms(raw)
		return algorithms
	}
	return &certificateCompressionZlibOffer
}

func certificateHandshakeMessage(body []byte, algorithms *certificateCompressionAlgorithms, enabled bool, cache *certificateCompressionCache) (uint8, []byte, error) {
	if !enabled || !supportsCertificateCompression(algorithms, certificateCompressionZlib) {
		return handshakeTypeCertificate, body, nil
	}
	var compressed []byte
	var err error
	if cache == nil {
		compressed, err = marshalCompressedCertificate(body)
	} else {
		compressed, err = cache.compress(body)
	}
	if err != nil {
		return 0, nil, err
	}
	if len(compressed) >= len(body) {
		return handshakeTypeCertificate, body, nil
	}
	return handshakeTypeCompressedCertificate, compressed, nil
}

func marshalCompressedCertificate(certificate []byte) ([]byte, error) {
	if len(certificate) >= 1<<24 {
		return nil, &ProtocolError{"Certificate message exceeds protocol limit"}
	}
	var compressed bytes.Buffer
	writer, _ := certificateZlibWriters.Get().(*zlib.Writer)
	if writer == nil {
		var err error
		writer, err = zlib.NewWriterLevel(&compressed, zlib.BestSpeed)
		if err != nil {
			return nil, err
		}
	} else {
		writer.Reset(&compressed)
	}
	_, writeErr := writer.Write(certificate)
	closeErr := writer.Close()
	certificateZlibWriters.Put(writer)
	if writeErr != nil {
		return nil, writeErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	w := newWireBuilder(8 + compressed.Len())
	w.u16(int(certificateCompressionZlib))
	w.u24(len(certificate))
	w.bytes24(compressed.Bytes())
	return w.b, w.err
}

func parseCertificateHandshakeMessage(typ uint8, body []byte, algorithms *certificateCompressionAlgorithms, maxSize int) (*certificateMessage, error) {
	if typ == handshakeTypeCertificate {
		return parseCertificateMessage(body, maxSize)
	}
	if typ != handshakeTypeCompressedCertificate {
		return nil, alertError(alertUnexpectedMessage, &ProtocolError{"unexpected certificate handshake message"})
	}
	certificate, pooled, err := decompressCertificate(body, algorithms, maxSize)
	if err != nil {
		return nil, err
	}
	message, parseErr := parseCertificateMessage(certificate, maxSize)
	if pooled {
		releaseCertificateDecompressionBuffer(certificate)
	}
	return message, parseErr
}

func decompressCertificate(body []byte, algorithms *certificateCompressionAlgorithms, maxSize int) ([]byte, bool, error) {
	if len(body) < 8 {
		return nil, false, alertError(alertDecodeError, &ProtocolError{"invalid CompressedCertificate message"})
	}
	algorithm := binary.BigEndian.Uint16(body)
	uncompressedLength := int(body[2])<<16 | int(body[3])<<8 | int(body[4])
	compressedLength := int(body[5])<<16 | int(body[6])<<8 | int(body[7])
	if compressedLength == 0 || compressedLength != len(body)-8 {
		return nil, false, alertError(alertDecodeError, &ProtocolError{"invalid CompressedCertificate vector"})
	}
	if algorithm != certificateCompressionZlib || !supportsCertificateCompression(algorithms, algorithm) {
		return nil, false, alertError(alertBadCertificate, &ProtocolError{"unoffered certificate compression algorithm"})
	}
	if uncompressedLength > maxSize {
		return nil, false, alertError(alertBadCertificate, &ProtocolError{"decompressed Certificate exceeds configured limit"})
	}
	source := bytes.NewReader(body[8:])
	reader, _ := certificateZlibReaders.Get().(zlibReadCloser)
	if reader == nil {
		fresh, err := zlib.NewReader(source)
		if err != nil {
			return nil, false, alertError(alertBadCertificate, &ProtocolError{"invalid compressed Certificate"})
		}
		var ok bool
		reader, ok = fresh.(zlibReadCloser)
		if !ok {
			_ = fresh.Close()
			return nil, false, alertError(alertBadCertificate, &ProtocolError{"zlib reader cannot be reset"})
		}
	} else if err := reader.Reset(source, nil); err != nil {
		return nil, false, alertError(alertBadCertificate, &ProtocolError{"invalid compressed Certificate"})
	}
	decompressed, pooled := certificateDecompressionBuffer(uncompressedLength)
	if _, err := io.ReadFull(reader, decompressed); err != nil {
		_ = reader.Close()
		certificateZlibReaders.Put(reader)
		if pooled {
			releaseCertificateDecompressionBuffer(decompressed)
		}
		return nil, false, alertError(alertBadCertificate, &ProtocolError{"compressed Certificate length mismatch"})
	}
	var extra [1]byte
	n, readErr := reader.Read(extra[:])
	closeErr := reader.Close()
	certificateZlibReaders.Put(reader)
	if n != 0 || readErr != io.EOF || closeErr != nil || source.Len() != 0 {
		if pooled {
			releaseCertificateDecompressionBuffer(decompressed)
		}
		return nil, false, alertError(alertBadCertificate, &ProtocolError{"compressed Certificate exceeds declared length"})
	}
	return decompressed, pooled, nil
}

func certificateDecompressionBuffer(length int) ([]byte, bool) {
	if length > maxPooledCertificate {
		return make([]byte, length), false
	}
	buffer, _ := certificateDecompressionBuffers.Get().(*[maxPooledCertificate]byte)
	if buffer == nil {
		buffer = new([maxPooledCertificate]byte)
	}
	return buffer[:length], true
}

func releaseCertificateDecompressionBuffer(buffer []byte) {
	certificateDecompressionBuffers.Put((*[maxPooledCertificate]byte)(buffer[:maxPooledCertificate]))
}

type zlibReadCloser interface {
	io.ReadCloser
	zlib.Resetter
}
