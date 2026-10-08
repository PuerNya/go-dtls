package dtls13

import (
	"bytes"
	"testing"
	"time"
)

func TestHandshakeFragmentRoundTrip(t *testing.T) {
	f := handshakeFragment{typ: 1, messageSequence: 7, length: 10, offset: 3, body: []byte("abcd")}
	for _, into := range []bool{false, true} {
		var wire []byte
		var err error
		if into {
			wire = bytes.Repeat([]byte{0xa5}, handshakeHeaderLen+len(f.body))
			err = marshalHandshakeFragmentInto(wire, f)
		} else {
			wire, err = marshalHandshakeFragment(f)
		}
		if err != nil {
			t.Fatal(err)
		}
		got, err := parseHandshakeFragments(wire)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].typ != f.typ || got[0].messageSequence != f.messageSequence || got[0].length != f.length || got[0].offset != f.offset || !bytes.Equal(got[0].body, f.body) {
			t.Fatalf("into=%t round trip mismatch: %#v", into, got)
		}
	}
}

func TestMarshalHandshakeFragmentIntoRejectsInvalidInputAtomically(t *testing.T) {
	tests := []struct {
		name     string
		fragment handshakeFragment
		dstLen   int
	}{
		{name: "length", fragment: handshakeFragment{length: 1 << 24}, dstLen: handshakeHeaderLen},
		{name: "offset", fragment: handshakeFragment{length: 1, offset: 1 << 24}, dstLen: handshakeHeaderLen},
		{name: "range", fragment: handshakeFragment{length: 1, offset: 1, body: []byte{1}}, dstLen: handshakeHeaderLen + 1},
		{name: "destination", fragment: handshakeFragment{length: 1, body: []byte{1}}, dstLen: handshakeHeaderLen},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dst := bytes.Repeat([]byte{0xa5}, test.dstLen)
			before := append([]byte(nil), dst...)
			if err := marshalHandshakeFragmentInto(dst, test.fragment); err == nil {
				t.Fatal("accepted invalid fragment")
			}
			if !bytes.Equal(dst, before) {
				t.Fatalf("destination changed: %x", dst)
			}
		})
	}
}

func TestHandshakeFragmentParserOwnershipModes(t *testing.T) {
	wire, err := marshalHandshakeFragment(handshakeFragment{typ: 1, messageSequence: 1, length: 4, body: []byte("body")})
	if err != nil {
		t.Fatal(err)
	}
	fragments, err := parseHandshakeFragments(wire)
	if err != nil {
		t.Fatal(err)
	}
	fragments[0].body[0] ^= 0xff
	if wire[handshakeHeaderLen] != 'b' {
		t.Fatal("parseHandshakeFragments body aliases input")
	}
	fragments, err = parseHandshakeFragmentsView(wire)
	if err != nil {
		t.Fatal(err)
	}
	fragments[0].body[0] ^= 0xff
	if wire[handshakeHeaderLen] != fragments[0].body[0] {
		t.Fatal("parseHandshakeFragmentsView did not reuse input")
	}
	var storage [1]handshakeFragment
	fragments, err = parseHandshakeFragmentsViewInto(wire, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	if len(fragments) != 1 || &fragments[0] != &storage[0] || &fragments[0].body[0] != &wire[handshakeHeaderLen] {
		t.Fatal("parseHandshakeFragmentsViewInto did not reuse destination and input")
	}
}
func TestReassemblyOutOfOrderAndOverlap(t *testing.T) {
	r := newReassembler()
	whole := []byte("abcdefghij")
	parts := []handshakeFragment{{typ: 1, messageSequence: 2, length: 10, offset: 5, body: whole[5:]}, {typ: 1, messageSequence: 2, length: 10, offset: 3, body: whole[3:7]}, {typ: 1, messageSequence: 2, length: 10, offset: 0, body: whole[:5]}}
	for i, p := range parts {
		got, done, err := r.add(p)
		if err != nil {
			t.Fatal(err)
		}
		if i < 2 && done {
			t.Fatal("completed early")
		}
		if i == 2 && (!done || !bytes.Equal(got, whole)) {
			t.Fatalf("bad result %q", got)
		}
	}
}

func TestReassemblySingleFragmentFastPathOwnsBodyAndRecordRange(t *testing.T) {
	r := newReassemblerWithLimit(4)
	input := []byte("body")
	number := recordNumber{epoch: 3, sequence: 9}
	body, complete, first, last, err := r.addProtectedRecord(handshakeFragment{
		typ: handshakeTypeCertificate, messageSequence: 1, length: 4, body: input,
	}, number)
	if err != nil || !complete || first != number || last != number {
		t.Fatalf("complete=%v first=%v last=%v err=%v", complete, first, last, err)
	}
	if len(r.messages) != 0 || r.allocated != 0 {
		t.Fatalf("fast path retained partial state: messages=%d allocated=%d", len(r.messages), r.allocated)
	}
	input[0] = 'x'
	if string(body) != "body" {
		t.Fatalf("completed body aliases input: %q", body)
	}
	body[1] = 'x'
	if string(input) != "xody" {
		t.Fatalf("input aliases completed body: %q", input)
	}
}

func TestReassemblySingleFragmentFastPathPreservesResourceLimits(t *testing.T) {
	fragment := handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 2, length: 2, body: []byte{1, 2}}
	if _, _, err := newReassemblerWithLimits(2, 0, 2).add(fragment); err == nil {
		t.Fatal("complete fragment bypassed incomplete-message count limit")
	}
	if _, _, err := newReassemblerWithLimits(2, 1, 1).add(fragment); err == nil {
		t.Fatal("complete fragment bypassed aggregate byte limit")
	}
	r := newReassemblerWithLimits(2, 1, 2)
	if _, complete, err := r.add(handshakeFragment{typ: fragment.typ, messageSequence: fragment.messageSequence, length: 2, body: []byte{1}}); err != nil || complete {
		t.Fatalf("partial fragment complete=%v err=%v", complete, err)
	}
	if _, _, err := r.addProtected(fragment, 3); err == nil {
		t.Fatal("complete fragment bypassed existing partial epoch check")
	}
}
func TestReassemblyRejectsConflictingOverlap(t *testing.T) {
	r := newReassembler()
	_, _, _ = r.add(handshakeFragment{typ: 1, messageSequence: 1, length: 3, offset: 0, body: []byte("ab")})
	if _, _, err := r.add(handshakeFragment{typ: 1, messageSequence: 1, length: 3, offset: 1, body: []byte("x")}); err == nil {
		t.Fatal("accepted conflicting overlap")
	}
}

func TestReassemblyBitmapWordBoundaries(t *testing.T) {
	const size = 130
	want := bytes.Repeat([]byte{0x5a}, size)
	r := newReassemblerWithLimit(size)
	fragments := []handshakeFragment{
		{typ: 1, messageSequence: 1, length: size, offset: 1, body: want[1:64]},
		{typ: 1, messageSequence: 1, length: size, offset: 64, body: want[64:129]},
		{typ: 1, messageSequence: 1, length: size, offset: 0, body: want[:1]},
		{typ: 1, messageSequence: 1, length: size, offset: 129, body: want[129:]},
	}
	for i, fragment := range fragments {
		got, complete, err := r.add(fragment)
		if err != nil {
			t.Fatal(err)
		}
		if i != len(fragments)-1 && complete {
			t.Fatalf("fragment %d completed early", i)
		}
		if i == len(fragments)-1 && (!complete || !bytes.Equal(got, want)) {
			t.Fatalf("final reassembly complete=%v body=%x", complete, got)
		}
	}
}

func TestReassemblyRejectsFragmentsSpanningKeyChange(t *testing.T) {
	r := newReassembler()
	first := handshakeFragment{typ: handshakeTypeNewSessionTicket, messageSequence: 7, length: 2, body: []byte{1}}
	if _, complete, err := r.addProtected(first, 3); err != nil || complete {
		t.Fatalf("first fragment complete=%v err=%v", complete, err)
	}
	second := handshakeFragment{typ: first.typ, messageSequence: first.messageSequence, length: first.length, offset: 1, body: []byte{2}}
	if _, _, err := r.addProtected(second, 4); err == nil {
		t.Fatal("reassembled one handshake message across a key change")
	} else if description, ok := protocolAlert(err); !ok || description != alertUnexpectedMessage {
		t.Fatalf("cross-epoch fragment alert=%d ok=%v err=%v", description, ok, err)
	}
	if body, complete, err := r.addProtected(second, 3); err != nil || !complete || !equalBytes(body, []byte{1, 2}) {
		t.Fatalf("same-epoch reassembly body=%x complete=%v err=%v", body, complete, err)
	}
}

func TestReassemblyLimit(t *testing.T) {
	r := newReassemblerWithLimit(2)
	_, _, err := r.add(handshakeFragment{typ: 1, messageSequence: 1, length: 3, body: []byte("a")})
	if err == nil {
		t.Fatal("accepted a handshake message over the configured limit")
	}
}

func TestReassemblyAggregateLimits(t *testing.T) {
	r := newReassemblerWithLimits(10, 2, 6)
	for seq := range uint16(2) {
		if _, _, err := r.add(handshakeFragment{typ: 1, messageSequence: seq, length: 3, body: []byte{1}}); err != nil {
			t.Fatal(err)
		}
	}
	if _, _, err := r.add(handshakeFragment{typ: 1, messageSequence: 2, length: 1, body: []byte{1}}); err == nil {
		t.Fatal("accepted too many incomplete messages")
	}
	r2 := newReassemblerWithLimits(10, 3, 5)
	if _, _, err := r2.add(handshakeFragment{typ: 1, messageSequence: 0, length: 3, body: []byte{1}}); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r2.add(handshakeFragment{typ: 1, messageSequence: 1, length: 3, body: []byte{1}}); err == nil {
		t.Fatal("exceeded aggregate byte limit")
	}
}

func TestReassemblyUsesCompactReceivedBitmap(t *testing.T) {
	r := newReassemblerWithLimit(1024)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 1, length: 1024, body: []byte{1}}
	if _, complete, err := r.add(fragment); err != nil || complete {
		t.Fatalf("initial fragment complete=%v err=%v", complete, err)
	}
	partial := r.messages[fragment.messageSequence]
	if partial == nil || len(partial.received) != 16 {
		t.Fatalf("received bitmap words=%d, want 16", len(partial.received))
	}
}

func TestServerHandshakeMessageOrder(t *testing.T) {
	valid := []struct {
		resumed bool
		types   []uint8
	}{
		{false, []uint8{handshakeTypeEncryptedExtensions, handshakeTypeCertificate, handshakeTypeCertificateVerify, handshakeTypeFinished}},
		{false, []uint8{handshakeTypeEncryptedExtensions, handshakeTypeCompressedCertificate, handshakeTypeCertificateVerify, handshakeTypeFinished}},
		{false, []uint8{handshakeTypeEncryptedExtensions, handshakeTypeCertificateRequest, handshakeTypeCertificate, handshakeTypeCertificateVerify, handshakeTypeFinished}},
		{false, []uint8{handshakeTypeEncryptedExtensions, handshakeTypeCertificateRequest, handshakeTypeCompressedCertificate, handshakeTypeCertificateVerify, handshakeTypeFinished}},
		{true, []uint8{handshakeTypeEncryptedExtensions, handshakeTypeFinished}},
	}
	for _, test := range valid {
		stage := serverExpectEncryptedExtensions
		for _, typ := range test.types {
			if err := stage.accept(typ, test.resumed); err != nil {
				t.Fatalf("valid sequence %v: %v", test.types, err)
			}
		}
		if stage != serverHandshakeComplete {
			t.Fatalf("sequence %v ended at stage %d", test.types, stage)
		}
	}
	for _, types := range [][]uint8{
		{handshakeTypeCertificate, handshakeTypeEncryptedExtensions},
		{handshakeTypeEncryptedExtensions, handshakeTypeEncryptedExtensions},
		{handshakeTypeEncryptedExtensions, handshakeTypeFinished},
		{handshakeTypeEncryptedExtensions, handshakeTypeCertificate, handshakeTypeCertificate},
	} {
		stage := serverExpectEncryptedExtensions
		var err error
		for _, typ := range types {
			if err = stage.accept(typ, false); err != nil {
				break
			}
		}
		if err == nil {
			t.Fatalf("accepted invalid server sequence %v", types)
		}
	}
}

func TestHandshakeRejectsWrongHelloMessage(t *testing.T) {
	for _, role := range []string{"client", "server"} {
		t.Run(role, func(t *testing.T) {
			config, err := (&Config{}).normalized()
			if err != nil {
				t.Fatal(err)
			}
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			deadline := time.Now().Add(time.Second)
			_ = right.SetDeadline(deadline)
			flight, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeFinished, body: []byte{1}}}, config.MTU, 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			if _, err = left.Write(flight.records[0].wire); err != nil {
				t.Fatal(err)
			}
			conn := &Conn{conn: right, config: config, handshakeDeadline: deadline}
			if role == "client" {
				err = conn.clientReceiveServerHello(&clientHandshakeState{})
			} else {
				err = conn.serverReceiveClientHello(&serverHandshakeState{preValidationConn: &amplificationConn{Conn: right, guard: &amplificationGuard{}}})
			}
			if description, ok := protocolAlert(err); !ok || description != alertUnexpectedMessage {
				t.Fatalf("wrong Hello message: alert=%d ok=%v err=%v, want unexpected_message", description, ok, err)
			}
		})
	}
}

func TestHandshakeFlightRejectsMessagesAfterFinished(t *testing.T) {
	for _, role := range []string{"client", "server"} {
		for _, trailing := range []struct {
			name string
			typ  uint8
		}{
			{"none", 0},
			{"Finished", handshakeTypeFinished},
			{"CertificateVerify", handshakeTypeCertificateVerify},
			{"NewSessionTicket", handshakeTypeNewSessionTicket},
		} {
			t.Run(role+"/"+trailing.name, func(t *testing.T) {
				config, err := (&Config{}).normalized()
				if err != nil {
					t.Fatal(err)
				}
				left, right := memoryDatagramPair()
				defer left.Close()
				defer right.Close()
				deadline := time.Now().Add(time.Second)
				_ = right.SetDeadline(deadline)
				conn := &Conn{conn: right, config: config, isClient: role == "client", handshakeDeadline: deadline}
				suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
				sender, receiver := recordCipherPair(t, suite.id, 2)
				ackSender, _ := recordCipherPair(t, suite.id, 2)
				secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
				schedule := &keySchedule{suite: suite, clientHandshakeTraffic: secret, serverHandshakeTraffic: secret}
				transcript := newTranscriptHash(suite.hash.New())
				peerTranscript := transcript.clone()
				var messages []handshakeMessage
				if conn.isClient {
					body := []byte{0, 0} // Empty EncryptedExtensions for a PSK handshake.
					messages = append(messages, handshakeMessage{typ: handshakeTypeEncryptedExtensions, sequence: 0, body: body})
					_ = peerTranscript.add(handshakeTypeEncryptedExtensions, 0, body)
				}
				var digest [maxSupportedHashSize]byte
				finished := schedule.finishedVerifyData(secret, peerTranscript.sumInto(digest[:0]))
				messages = append(messages, handshakeMessage{typ: handshakeTypeFinished, sequence: uint16(len(messages)), body: finished})
				if trailing.typ != 0 {
					messages = append(messages, handshakeMessage{typ: trailing.typ, sequence: uint16(len(messages)), body: finished})
				}
				// Put Finished and its successor in one authenticated record so
				// the receive loop delivers both in the same batch.
				var payload []byte
				for _, message := range messages {
					fragment, marshalErr := marshalHandshakeFragment(handshakeFragment{typ: message.typ, messageSequence: message.sequence, length: uint32(len(message.body)), body: message.body})
					if marshalErr != nil {
						t.Fatal(marshalErr)
					}
					payload = append(payload, fragment...)
				}
				wire, err := sender.seal(recordTypeHandshake, payload)
				if err != nil {
					t.Fatal(err)
				}
				if _, err = left.Write(wire); err != nil {
					t.Fatal(err)
				}
				if conn.isClient {
					err = conn.clientProcessServerFlight(&clientHandshakeState{
						hello: &clientHello{}, sh: &serverHello{}, suite: suite, usingPSK: true,
						schedule: schedule, transcript: transcript, receiveCipher: receiver, sendCipher: ackSender,
						inbox: newHandshakeInbox(0, config.MaxHandshakeMessage, config.MaxBufferedHandshakeMessages, config.MaxBufferedHandshakeBytes),
					})
				} else {
					err = conn.serverProcessClientFlight(&serverHandshakeState{
						suite: suite, usingPSK: true, hrrUsed: true, schedule: schedule,
						transcript: transcript, clientCipher: receiver, serverCipher: ackSender,
					})
				}
				if trailing.typ == 0 {
					if err != nil {
						t.Fatalf("valid Finished flight: %v", err)
					}
				} else if description, ok := protocolAlert(err); !ok || description != alertUnexpectedMessage {
					t.Fatalf("message after Finished: alert=%d ok=%v err=%v, want unexpected_message", description, ok, err)
				}
			})
		}
	}
}

func TestClientCertificateTypeSelectionInPSKHandshake(t *testing.T) {
	for _, test := range []struct {
		name      string
		extension uint16
		pha       bool
		wantError bool
	}{
		{name: "client without CertificateRequest", extension: extClientCertificateType, wantError: true},
		{name: "client for PHA", extension: extClientCertificateType, pha: true},
		{name: "server in PSK handshake", extension: extServerCertificateType},
	} {
		t.Run(test.name, func(t *testing.T) {
			config, err := (&Config{
				ServerCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
				ClientCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
			}).normalized()
			if err != nil {
				t.Fatal(err)
			}
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			deadline := time.Now().Add(time.Second)
			_ = right.SetDeadline(deadline)
			conn := &Conn{conn: right, config: config, isClient: true, handshakeDeadline: deadline}
			suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
			sender, receiver := recordCipherPair(t, suite.id, 2)
			ackSender, _ := recordCipherPair(t, suite.id, 2)
			secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
			schedule := &keySchedule{suite: suite, clientHandshakeTraffic: secret, serverHandshakeTraffic: secret}
			transcript := newTranscriptHash(suite.hash.New())
			peerTranscript := transcript.clone()
			hello := &clientHello{postHandshakeAuth: test.pha, unknownExtensions: map[uint16][]byte{test.extension: {1, byte(CertificateTypeRawPublicKey)}}}
			eeBody, err := (&encryptedExtensions{extensions: map[uint16][]byte{test.extension: {byte(CertificateTypeRawPublicKey)}}}).marshal()
			if err != nil {
				t.Fatal(err)
			}
			_ = peerTranscript.add(handshakeTypeEncryptedExtensions, 0, eeBody)
			var digest [maxSupportedHashSize]byte
			finished := schedule.finishedVerifyData(secret, peerTranscript.sumInto(digest[:0]))
			flight, err := buildProtectedFlight([]handshakeMessage{
				{typ: handshakeTypeEncryptedExtensions, sequence: 0, body: eeBody},
				{typ: handshakeTypeFinished, sequence: 1, body: finished},
			}, config.MTU, sender)
			if err != nil {
				t.Fatal(err)
			}
			if err = conn.writeFlight(left, flight); err != nil {
				t.Fatal(err)
			}
			err = conn.clientProcessServerFlight(&clientHandshakeState{
				hello: hello, sh: &serverHello{}, suite: suite, usingPSK: true,
				schedule: schedule, transcript: transcript, receiveCipher: receiver, sendCipher: ackSender,
				inbox: newHandshakeInbox(0, config.MaxHandshakeMessage, config.MaxBufferedHandshakeMessages, config.MaxBufferedHandshakeBytes),
			})
			if test.wantError {
				if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
					t.Fatalf("certificate type without CertificateRequest or PHA: %v", err)
				}
			} else if err != nil {
				t.Fatalf("valid PSK certificate type selection: %v", err)
			}
		})
	}
}
