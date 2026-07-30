package dtls13

import (
	"bytes"
	"testing"
)

func TestHandshakeFragmentRoundTrip(t *testing.T) {
	f := handshakeFragment{typ: 1, messageSequence: 7, length: 10, offset: 3, body: []byte("abcd")}
	b, err := marshalHandshakeFragment(f)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseHandshakeFragments(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].typ != f.typ || got[0].messageSequence != f.messageSequence || got[0].length != f.length || got[0].offset != f.offset || !bytes.Equal(got[0].body, f.body) {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestMarshalHandshakeFragmentIntoMatchesAllocated(t *testing.T) {
	fragment := handshakeFragment{typ: 1, messageSequence: 7, length: 10, offset: 3, body: []byte("abcd")}
	want, err := marshalHandshakeFragment(fragment)
	if err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(want))
	if err = marshalHandshakeFragmentInto(got, fragment); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("into=%x allocated=%x", got, want)
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
	for seq := uint16(0); seq < 2; seq++ {
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
