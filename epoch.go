package dtls13

import (
	"errors"
	"sync"
)

// epochSet owns receive ciphers and their independent replay windows. The
// current epoch is tried when its low bits match; otherwise the most recent
// installed past epoch with matching low bits is selected (RFC 9147 4.2.2).
type epochSet struct {
	mu      sync.RWMutex
	current uint64
	ciphers map[uint64]*recordCipher
}

func newEpochSet() *epochSet { return &epochSet{ciphers: make(map[uint64]*recordCipher)} }
func (s *epochSet) install(cipher *recordCipher) error {
	if cipher == nil {
		return errors.New("dtls13: cannot install nil epoch cipher")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.ciphers[cipher.epoch]; exists {
		return errors.New("dtls13: epoch cipher already installed")
	}
	s.ciphers[cipher.epoch] = cipher
	return nil
}
func (s *epochSet) setCurrent(epoch uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if epoch < s.current {
		return errors.New("dtls13: epoch cannot move backwards")
	}
	if _, ok := s.ciphers[epoch]; !ok {
		return errors.New("dtls13: current epoch cipher is not installed")
	}
	s.current = epoch
	return nil
}
func (s *epochSet) selectCipher(first byte) (*recordCipher, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectCipherLocked(first)
}

func (s *epochSet) selectCipherLocked(first byte) (*recordCipher, error) {
	bits := uint64(first & unifiedHeaderEpochMask)
	if cipher := s.ciphers[s.current]; cipher != nil && s.current&unifiedHeaderEpochMask == bits {
		return cipher, nil
	}
	var selected *recordCipher
	var selectedEpoch uint64
	for epoch, cipher := range s.ciphers {
		if epoch >= s.current || epoch&unifiedHeaderEpochMask != bits {
			continue
		}
		if selected == nil || epoch > selectedEpoch {
			selected, selectedEpoch = cipher, epoch
		}
	}
	if selected == nil {
		return nil, &ProtocolError{"no receive keys for record epoch"}
	}
	return selected, nil
}
func (s *epochSet) open(datagram []byte) (content []byte, contentType uint8, epoch uint64, consumed int, err error) {
	return s.openRecord(datagram, false)
}

func (s *epochSet) openInPlace(datagram []byte) (content []byte, contentType uint8, epoch uint64, consumed int, err error) {
	return s.openRecord(datagram, true)
}

func (s *epochSet) openRecord(datagram []byte, inPlace bool) (content []byte, contentType uint8, epoch uint64, consumed int, err error) {
	if len(datagram) == 0 {
		return nil, 0, 0, 0, &ProtocolError{"empty datagram"}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	cipher, err := s.selectCipherLocked(datagram[0])
	if err != nil {
		return nil, 0, 0, 0, err
	}
	if inPlace {
		content, contentType, consumed, err = cipher.openInPlace(datagram)
	} else {
		content, contentType, consumed, err = cipher.open(datagram)
	}
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return content, contentType, cipher.epoch, consumed, nil
}

func (s *epochSet) shouldRequestKeyUpdateForAuthFailures(first byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cipher, err := s.selectCipherLocked(first)
	return err == nil && cipher.shouldRequestKeyUpdateForAuthFailures()
}

func (s *epochSet) discardBefore(epoch uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for candidate := range s.ciphers {
		if candidate < epoch {
			delete(s.ciphers, candidate)
		}
	}
}

func (s *epochSet) clear() {
	s.mu.Lock()
	s.current = 0
	s.ciphers = make(map[uint64]*recordCipher)
	s.mu.Unlock()
}
