package dtls13

import (
	"crypto/tls"
	"testing"
)

func TestCertificateMessageContextAndExtensionValidation(t *testing.T) {
	valid := &certificateMessage{requestContext: []byte{1}, certificates: []certificateEntry{{data: []byte{2}, extensions: map[uint16][]byte{}}}}
	if err := validateCertificateMessage(valid, []byte{1}); err != nil {
		t.Fatal(err)
	}
	if err := validateCertificateMessage(valid, nil); err == nil {
		t.Fatal("accepted mismatched certificate_request_context")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("context alert=%d ok=%v", description, ok)
	}
	withExtension := &certificateMessage{certificates: []certificateEntry{{data: []byte{2}, extensions: map[uint16][]byte{0xffa5: {1}}}}}
	if err := validateCertificateMessage(withExtension, nil); err == nil {
		t.Fatal("accepted unsolicited CertificateEntry extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertUnsupportedExtension {
		t.Fatalf("extension alert=%d ok=%v", description, ok)
	}
}

func TestValidateEncryptedExtensions(t *testing.T) {
	alpn, err := marshalALPN([]string{"h3"})
	if err != nil {
		t.Fatal(err)
	}
	groups, err := marshalSupportedGroups([]tls.CurveID{tls.X25519}, nil)
	if err != nil {
		t.Fatal(err)
	}
	hello := &clientHello{
		serverName:         "example.test",
		alpn:               []string{"coap", "h3"},
		earlyData:          true,
		supportedGroups:    []tls.CurveID{tls.X25519},
		recordSizeLimit:    512,
		hasRecordSizeLimit: true,
	}
	message := &encryptedExtensions{extensions: map[uint16][]byte{
		extServerName:      nil,
		extALPN:            alpn,
		extEarlyData:       nil,
		extSupportedGroups: groups,
		extRecordSizeLimit: {0x01, 0x00},
	}}
	protocol, earlyData, _, err := validateEncryptedExtensions(hello, message)
	if err != nil {
		t.Fatal(err)
	}
	if protocol != "h3" || !earlyData || !message.hasRecordSizeLimit || message.recordSizeLimit != 256 {
		t.Fatalf("unexpected negotiation result: protocol=%q early_data=%v", protocol, earlyData)
	}
}

func TestValidateEncryptedExtensionsRejectsInvalidExtensions(t *testing.T) {
	alpn, err := marshalALPN([]string{"h3"})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		hello      *clientHello
		extensions map[uint16][]byte
		wantAlert  uint8
	}{
		{
			name:       "record size limit with max fragment length",
			hello:      &clientHello{hasRecordSizeLimit: true},
			extensions: map[uint16][]byte{extRecordSizeLimit: {0, 64}, extMaxFragmentLength: {1}},
			wantAlert:  alertIllegalParameter,
		},
		{
			name:       "unoffered record size limit",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{extRecordSizeLimit: {0, 64}},
			wantAlert:  alertUnsupportedExtension,
		},
		{
			name:       "short record size limit",
			hello:      &clientHello{hasRecordSizeLimit: true},
			extensions: map[uint16][]byte{extRecordSizeLimit: {64}},
			wantAlert:  alertDecodeError,
		},
		{
			name:       "record size limit below 64",
			hello:      &clientHello{hasRecordSizeLimit: true},
			extensions: map[uint16][]byte{extRecordSizeLimit: {0, 63}},
			wantAlert:  alertIllegalParameter,
		},
		{
			name:       "record size limit above DTLS 1.3 maximum",
			hello:      &clientHello{hasRecordSizeLimit: true},
			extensions: map[uint16][]byte{extRecordSizeLimit: {0xff, 0xff}},
			wantAlert:  alertIllegalParameter,
		},
		{
			name:       "unoffered ALPN",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{extALPN: alpn},
			wantAlert:  alertUnsupportedExtension,
		},
		{
			name:       "unoffered early data",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{extEarlyData: nil},
			wantAlert:  alertUnsupportedExtension,
		},
		{
			name:       "connection ID in EncryptedExtensions",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{extConnectionID: nil},
			wantAlert:  alertIllegalParameter,
		},
		{
			name:       "certificate compression in EncryptedExtensions",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{extCompressCertificate: {2, 0, 1}},
			wantAlert:  alertIllegalParameter,
		},
		{
			name:       "unknown extension",
			hello:      &clientHello{},
			extensions: map[uint16][]byte{0xffff: nil},
			wantAlert:  alertUnsupportedExtension,
		},
		{
			name:       "malformed supported groups",
			hello:      &clientHello{supportedGroups: []tls.CurveID{tls.X25519}},
			extensions: map[uint16][]byte{extSupportedGroups: {0, 3, 0, 29, 0}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message := &encryptedExtensions{extensions: test.extensions}
			if _, _, _, err := validateEncryptedExtensions(test.hello, message); err == nil {
				t.Fatal("accepted invalid EncryptedExtensions")
			} else if test.wantAlert != 0 {
				if description, ok := protocolAlert(err); !ok || description != test.wantAlert {
					t.Fatalf("alert=%d ok=%v err=%v", description, ok, err)
				}
			}
		})
	}
}

func TestEarlyDataAcceptanceRequiresFirstPSKIdentity(t *testing.T) {
	zero, one := uint16(0), uint16(1)
	for _, selected := range []*uint16{nil, &one} {
		err := validateEarlyDataSelection(true, selected)
		if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
			t.Fatalf("selected_identity=%v returned alert=%d ok=%v err=%v", selected, description, ok, err)
		}
	}
	if err := validateEarlyDataSelection(true, &zero); err != nil {
		t.Fatal(err)
	}
	if err := validateEarlyDataSelection(false, nil); err != nil {
		t.Fatal(err)
	}
}

func TestCertificateVerifyWireAndSignatureErrorClassification(t *testing.T) {
	if _, err := parseCertificateVerify([]byte{0}); err == nil {
		t.Fatal("accepted truncated CertificateVerify")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("truncated CertificateVerify alert=%d ok=%v err=%v", description, ok, err)
	}
	message, err := parseCertificateVerify([]byte{0x08, 0x07, 0, 0})
	if err != nil || len(message.signature) != 0 {
		t.Fatalf("empty CertificateVerify signature message=%#v err=%v", message, err)
	}
}
