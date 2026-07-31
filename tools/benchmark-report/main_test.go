package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseAndWriteReport(t *testing.T) {
	input := `wolfssl: 0123456789abcdef (Linux Release static)
go version go1.26.0 linux/amd64
goos: linux
goarch: amd64
cpu: Example CPU
BenchmarkHandshake-2 10 9 ns/op 30 B/op 2 allocs/op
BenchmarkHandshake-2 10 3 ns/op 10 B/op 0 allocs/op
BenchmarkHandshake-2 10 6 ns/op 20 B/op 1 allocs/op
BenchmarkPeer/Go|Wolf-2 1 2.5 go_ms/conn
`
	report, err := parseReport(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err = writeReport(&output, report, "abc123", "", "2026-07-31T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, want := range []string{
		"`BenchmarkHandshake-2` | 3 | 6 ns/op<br>20 B/op<br>1 allocs/op",
		"`BenchmarkPeer/Go&#124;Wolf-2` | 1 | 2.5 go_ms/conn",
		"- Platform: `linux/amd64, Example CPU`",
		"- wolfSSL: `0123456789abcdef (Linux Release static)`",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("report does not contain %q:\n%s", want, text)
		}
	}
}
