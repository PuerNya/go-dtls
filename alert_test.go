package dtls13

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
)

func TestAlertAndClosureOrdering(t *testing.T) {
	wire, err := (alertMessage{level: alertLevelWarning, description: alertCloseNotify}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseAlert(wire)
	if err != nil || got.description != alertCloseNotify {
		t.Fatalf("got=%v err=%v", got, err)
	}
	var closure closureState
	closure.receive(recordNumber{epoch: 3, sequence: 9})
	if closure.ignore(recordNumber{epoch: 3, sequence: 8}) {
		t.Fatal("ignored record before closure")
	}
	if !closure.ignore(recordNumber{epoch: 3, sequence: 10}) || !closure.ignore(recordNumber{epoch: 4, sequence: 0}) {
		t.Fatal("accepted record after closure")
	}
}

func TestUnknownAlertLevelIsAnErrorAlert(t *testing.T) {
	for _, description := range []uint8{alertCloseNotify, alertUserCanceled} {
		alert, err := parseAlert([]byte{255, description})
		if err != nil {
			t.Fatal(err)
		}
		if alert.isCloseNotify() || alert.isUserCanceled() {
			t.Fatalf("unknown alert level was accepted for description %d", description)
		}
	}
}

// TestProtocolAlertClassification pins the rule that the alert comes from the
// error's type and explicit alert, never from the wording of its Reason.
func TestProtocolAlertClassification(t *testing.T) {
	tests := []struct {
		err  error
		want uint8
		ok   bool
	}{
		{alertError(alertUnexpectedMessage, &ProtocolError{"unexpected handshake message"}), alertUnexpectedMessage, true},
		{alertError(alertDecodeError, &ProtocolError{"truncated handshake fragment header"}), alertDecodeError, true},
		{&ProtocolError{"server selected an unoffered cipher suite"}, alertIllegalParameter, true},
		{&ProtocolError{"unexpected handshake message"}, alertIllegalParameter, true},
		{&ProtocolError{"truncated handshake fragment header"}, alertIllegalParameter, true},
		{alertError(alertTooManyCIDsRequest, &ProtocolError{"too many"}), alertTooManyCIDsRequest, true},
		{alertError(alertDecryptError, errors.New("bad signature")), alertDecryptError, true},
		{errors.New("network failure"), 0, false},
	}
	for _, test := range tests {
		got, ok := protocolAlert(test.err)
		if got != test.want || ok != test.ok {
			t.Fatalf("protocolAlert(%v)=(%d,%v), want (%d,%v)", test.err, got, ok, test.want, test.ok)
		}
	}
}

// TestVectorOverflowErrorKeepsProtocolClassification pins the invariant that
// makes vectorOverflowError safe to use: wrapping does not hide the underlying
// protocol failure from type-based classification, and the width of the
// overflowed vector survives for callers that must react to one width.
func TestVectorOverflowErrorKeepsProtocolClassification(t *testing.T) {
	for _, bits := range []int{8, 16, 24} {
		wrapped := &vectorOverflowError{bits: bits, err: &ProtocolError{"vector overflow"}}
		var protocol *ProtocolError
		if !errors.As(wrapped, &protocol) {
			t.Fatalf("bits=%d: errors.As did not resolve the wrapped *ProtocolError", bits)
		}
		var overflow *vectorOverflowError
		if !errors.As(wrapped, &overflow) || overflow.bits != bits {
			t.Fatalf("bits=%d: errors.As did not resolve the overflow width", bits)
		}
		want := alertIllegalParameter
		if got, ok := protocolAlert(wrapped); !ok || got != want {
			t.Fatalf("bits=%d: protocolAlert=(%d,%v), want (%d,true)", bits, got, ok, want)
		}
	}
}

// TestProtocolAlertIgnoresReasonWording guards the property the classifier was
// rewritten to provide: rewording a diagnostic must not change the wire alert.
func TestProtocolAlertIgnoresReasonWording(t *testing.T) {
	for _, reason := range []string{
		"unexpected handshake message",
		"truncated handshake fragment header",
		"duplicate ClientHello extension",
		"malformed certificate entry",
		"invalid record length",
		"the peer omitted CertificateVerify",
		"",
	} {
		if got, ok := protocolAlert(&ProtocolError{reason}); !ok || got != alertIllegalParameter {
			t.Fatalf("protocolAlert(%q)=(%d,%v), want (%d,true)", reason, got, ok, alertIllegalParameter)
		}
	}
}

func TestOutboundAlertClassification(t *testing.T) {
	tests := []struct {
		err  error
		want uint8
		ok   bool
	}{
		{nil, 0, false},
		{errors.New("local callback failed"), alertGeneralError, true},
		{alertError(alertUnexpectedMessage, &ProtocolError{"unexpected handshake message"}), alertUnexpectedMessage, true},
		{&ProtocolError{"unexpected handshake message"}, alertIllegalParameter, true},
		{AlertError(alertIllegalParameter), 0, false},
		{io.EOF, 0, false},
		{context.Canceled, 0, false},
		{&net.OpError{Op: "read", Net: "udp", Err: errors.New("network failure")}, 0, false},
	}
	for _, test := range tests {
		got, ok := outboundAlert(test.err)
		if got != test.want || ok != test.ok {
			t.Fatalf("outboundAlert(%v)=(%d,%v), want (%d,%v)", test.err, got, ok, test.want, test.ok)
		}
	}
}

func TestTLSWireShapeErrorsUseDecodeError(t *testing.T) {
	tests := []struct {
		name string
		err  func() error
	}{
		{"alert length", func() error { _, err := parseAlert([]byte{alertLevelFatal}); return err }},
		{"Finished length", func() error { _, err := parseFinished(make([]byte, 31), 32); return err }},
		{"KeyUpdate length", func() error { _, err := parseKeyUpdate(nil); return err }},
		{"signature algorithms length", func() error { _, err := parseSignatureSchemes([]byte{0, 1, 4}); return err }},
		{"supported groups length", func() error { _, err := parseSupportedGroups([]byte{0, 1, 0}); return err }},
		{"empty cookie", func() error { _, err := parseCookie([]byte{0, 0}); return err }},
		{"PSK modes length", func() error { _, err := parsePSKKeyExchangeModes([]byte{0}); return err }},
		{"empty ALPN list", func() error { _, err := parseALPN([]byte{0, 0}); return err }},
		{"empty server key share", func() error { _, err := parseKeyShares(nil, false); return err }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.err()
			if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
				t.Fatalf("alert=%d ok=%v err=%v; want decode_error", description, ok, err)
			}
		})
	}
	_, err := parseKeyUpdate([]byte{2})
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("invalid KeyUpdate enum alert=%d ok=%v err=%v; want illegal_parameter", description, ok, err)
	}
}
