package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAndWriteReport(t *testing.T) {
	input := `wolfssl: 0123456789abcdef (Linux Release static)
go version go1.26.0 linux/amd64
goos: linux
goarch: amd64
cpu: Example CPU
BenchmarkProtectedRecordSeal-2 10 1500 ns/op 800 MB/s 1280 B/op 1 allocs/op
BenchmarkBuildPlainFlight-2 10 9 ns/op 30 B/op 2 allocs/op
BenchmarkBuildPlainFlight-2 10 3 ns/op 10 B/op 0 allocs/op
BenchmarkBuildPlainFlight-2 10 6 ns/op 20 B/op 1 allocs/op
BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer-2 1 2000000 ns/op 2.5 go_ms/conn
BenchmarkWolfSSLFeatureRealUDP/SessionResumption/WolfSSLClient/GoServer-2 1 3000000 ns/op 3 wolf_process_ms/pair
BenchmarkMutualTLSHandshakeLifecycle/Full-2 1 2000000 ns/op 100000 B/op 900 allocs/op
BenchmarkConnectionHandshakeLifecycle-2 1 1500000 ns/op 90000 B/op 800 allocs/op
`
	report, err := parseReport(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err = writeReport(&output, report, "abc123", "", "2026-07-31T00:00:00Z", reportLanguages["en"], reportLanguageLinks("benchmark.en.md")); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, want := range []string{
		"[Connection lifecycle (2)](#section-connection-lifecycle)",
		"[Real UDP interoperability (2)](#section-real-udp-interoperability)",
		"  - [go-dtls client -> wolfSSL server (1)](#real-udp-go-dtls-client-wolfssl-server)",
		"  - [wolfSSL client -> go-dtls server (1)](#real-udp-wolfssl-client-go-dtls-server)",
		"[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)",
		"| Benchmark | Samples | Median time | Harness memory | Harness allocations |",
		"| --- | :---: | :---: | :---: | :---: |",
		"| --- | :---: | :---: | :---: | :---: | :---: |",
		"## Connection lifecycle",
		"Certificate-authenticated full handshake / AES-128-GCM | 1 | 1.5 ms/op | 90000 B/op | 800 allocs/op",
		"Full mTLS handshake | 1 | 2 ms/op | 100000 B/op | 900 allocs/op",
		"## Real UDP interoperability",
		"### go-dtls client -> wolfSSL server",
		"Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.",
		"KeyUpdate + 1-RTT application-data round trip | 1 | 2.5 ms/conn | - | -",
		"Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.",
		"Session resumption handshake | 1 | 3 ms/pair | - | -",
		"## Record layer and reliability",
		"Record seal | 1 | 1.5 us/op | 800 MB/s | 1280 B/op | 1 allocs/op",
		"Build Plain Flight | 3 | 6 ns/op | - | 20 B/op | 1 allocs/op",
		"- Platform: `linux/amd64, Example CPU`",
		"- wolfSSL: `0123456789abcdef (Linux Release static)`",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("report does not contain %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "<details>") {
		t.Fatalf("report still contains collapsed sections:\n%s", text)
	}
	for _, marker := range []string{"Implementation markers:", " (G)", " (W)"} {
		if strings.Contains(text, marker) {
			t.Fatalf("report still contains implementation marker %q:\n%s", marker, text)
		}
	}
	if strings.Contains(text, "RFC ") {
		t.Fatalf("report benchmark labels still contain RFC identifiers:\n%s", text)
	}
	if strings.Contains(text, "| `") {
		t.Fatalf("report benchmark labels still use code spans:\n%s", text)
	}
	connectionSection := text[strings.Index(text, "## Connection lifecycle"):strings.Index(text, "## Real UDP interoperability")]
	if strings.Contains(connectionSection, "Throughput") {
		t.Fatalf("section without throughput metrics contains a throughput column:\n%s", connectionSection)
	}
	_, recordSection, _ := strings.Cut(text, "## Record layer and reliability")
	if !strings.Contains(recordSection, "Throughput") {
		t.Fatalf("section with throughput metrics does not contain a throughput column:\n%s", recordSection)
	}
	connection := strings.Index(text, "## Connection lifecycle")
	realUDP := strings.Index(text, "## Real UDP interoperability")
	record := strings.Index(text, "## Record layer and reliability")
	if connection >= realUDP || realUDP >= record {
		t.Fatalf("report sections are not in the expected order:\n%s", text)
	}
	connectionBenchmark := strings.Index(text, "Certificate-authenticated full handshake / AES-128-GCM")
	mutualTLSBenchmark := strings.Index(text, "Full mTLS handshake")
	if connectionBenchmark >= mutualTLSBenchmark {
		t.Fatalf("benchmarks are not sorted by feature:\n%s", text)
	}
}

func TestValidateBenchmarkCoverageRejectsUnknownNames(t *testing.T) {
	for _, name := range []string{
		"BenchmarkTranscriptClone/1301", "BenchmarkParseExtensions/Owned",
		"BenchmarkCertificateCompression/Decompress", "BenchmarkHandshakeInboxSequentialSingleFragment",
		"BenchmarkNewFeature-2", "BenchmarkFlightWireWindow/NewMode-2", "BenchmarkGREASEHandshakeLifecycle-2",
		"BenchmarkListenerCIDRouting/Connections2/Fixed/Hit", "BenchmarkListenerCIDRouting/Connections1/Unknown/Miss",
	} {
		if err := validateBenchmarkCoverage(name); err == nil {
			t.Errorf("uncovered benchmark %q was accepted", name)
		}
	}
}

func TestAllowPartialReportShowsSamplesAndSkips(t *testing.T) {
	const skipped = "BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer"
	var input bytes.Buffer
	input.WriteString("wolfssl: " + reviewedWolfSSLCommit + " (Linux Release static)\n")
	input.WriteString("BenchmarkNewFeatureRealUDP-2 1 100 ns/op 1 B/op 1 allocs/op\n")
	input.WriteString(skipped + " 1 200 ns/op 1 B/op 1 allocs/op\n")
	encoder := json.NewEncoder(&input)
	for _, event := range []goTestEvent{
		{Action: "output", Test: skipped, Output: "    interop_test.go:1: <img src=x onerror=alert(1)>\n"},
		{Action: "skip", Test: skipped},
	} {
		if err := encoder.Encode(event); err != nil {
			t.Fatal(err)
		}
	}
	path := filepath.Join(t.TempDir(), "benchmark.txt")
	if err := os.WriteFile(path, input.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readReportFile(path, "", false); err == nil {
		t.Fatal("strict report accepted incomplete and unmapped output")
	}
	report, err := readReportFile(path, "", true)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := writeReport(&output, report, "", "", "", reportLanguages["en"], ""); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, want := range []string{"100 ns/op", "200 ns/op", "Skipped: &lt;img src=x onerror=alert(1)>"} {
		if !strings.Contains(text, want) {
			t.Fatalf("display report lacks %q:\n%s", want, text)
		}
	}
	if got := strings.Count(text, "KeyUpdate + 1-RTT application-data round trip"); got != 2 {
		t.Fatalf("mixed result and skip produced %d rows, want 2", got)
	}
	if strings.Contains(text, "last verified against wolfSSL commit") {
		t.Fatal("unknown skip claims a verified peer limit")
	}
}

func TestAllowPartialReportAllowsIncompleteMatrix(t *testing.T) {
	path := filepath.Join(t.TempDir(), "benchmark.txt")
	const name = "BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer"
	input := "wolfssl: " + reviewedWolfSSLCommit + " (Linux Release static)\n" +
		`{"Action":"run","Test":"` + name + `"}` + "\n" +
		name + "-2 1 100 ns/op 1 B/op 1 allocs/op\n"
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readReportFile(path, "", false); err == nil {
		t.Fatal("strict report accepted an incomplete real UDP matrix")
	}
	if _, err := readReportFile(path, "", true); err != nil {
		t.Fatalf("display report rejected observed samples: %v", err)
	}
}

func TestParseReportRejectsMalformedInput(t *testing.T) {
	for name, input := range map[string]string{
		"empty":             "",
		"json":              "{\n",
		"unfinished output": `{"Action":"output","Output":"BenchmarkNewFeature-2"}`,
		"iterations":        "BenchmarkNewFeature-2 invalid 100 ns/op\n",
		"metric pair":       "BenchmarkNewFeature-2 1 100 ns/op 1\n",
		"NaN":               "BenchmarkNewFeature-2 1 NaN ns/op\n",
		"Inf":               "BenchmarkNewFeature-2 1 +Inf ns/op\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseReport(strings.NewReader(input)); err == nil {
				t.Fatal("malformed input was accepted")
			}
		})
	}
}

func TestMarkdownTextEscapesArtifactInput(t *testing.T) {
	const input = "<img src=x onerror=alert(1)> | `value` & [link](https://example.invalid)"
	const want = "&lt;img src=x onerror=alert(1)> &#124; &#96;value&#96; &amp; &#91;link&#93;(https://example.invalid)"
	if got := markdownText(input); got != want {
		t.Fatalf("markdownText(%q) = %q, want %q", input, got, want)
	}
}

func TestCompareReports(t *testing.T) {
	baselineInput := `commit: base
BenchmarkConnectionHandshakeLifecycle-1 100 100 ns/op 1000 B/op 10 allocs/op
BenchmarkConnectionHandshakeLifecycle-1 100 101 ns/op 1000 B/op 10 allocs/op
BenchmarkConnectionHandshakeLifecycle-1 100 99 ns/op 1000 B/op 10 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 200 ns/op 2000 B/op 20 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 201 ns/op 2000 B/op 20 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 199 ns/op 2000 B/op 20 allocs/op
`
	candidateInput := `commit: head
BenchmarkConnectionHandshakeLifecycle-1 100 103 ns/op 1000 B/op 10 allocs/op
BenchmarkConnectionHandshakeLifecycle-1 100 104 ns/op 1000 B/op 10 allocs/op
BenchmarkConnectionHandshakeLifecycle-1 100 102 ns/op 1000 B/op 10 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 220 ns/op 2001 B/op 21 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 221 ns/op 2001 B/op 21 allocs/op
BenchmarkMutualTLSHandshakeLifecycle/Full-1 100 219 ns/op 2001 B/op 21 allocs/op
`
	baseline, err := parseReport(strings.NewReader(baselineInput))
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := parseReport(strings.NewReader(candidateInput))
	if err != nil {
		t.Fatal(err)
	}
	comparison, err := compareReports(baseline, candidate, 3, 5, false)
	if err != nil {
		t.Fatal(err)
	}
	if !comparison.failed {
		t.Fatal("comparison passed despite a stable time and allocation regression")
	}
	if len(comparison.items) != 2 || len(comparison.items[0].failures) != 0 || len(comparison.items[1].failures) != 3 {
		t.Fatalf("unexpected comparison: %+v", comparison.items)
	}
	var output bytes.Buffer
	if err = writeComparisonReport(&output, comparison, reportLanguages["zh-CN"]); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"结论：`FAIL`", "Base：`base`", "Head：`head`", "完整 mTLS 握手", "time +10.000%", "B/op 2000 -> 2001", "allocs/op 20 -> 21"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("comparison report does not contain %q:\n%s", want, output.String())
		}
	}
	passing, err := compareReports(baseline, baseline, 3, 5, false)
	if err != nil {
		t.Fatal(err)
	}
	output.Reset()
	if err = writeComparisonReport(&output, passing, reportLanguages["zh-CN"]); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "结论：`PASS`") || !strings.Contains(output.String(), "| 0/3 |") {
		t.Fatalf("unexpected passing comparison report:\n%s", output.String())
	}
}

func TestCompareReportsRejectsUnapprovedWorkloadChanges(t *testing.T) {
	const baselineInput = `commit: base
BenchmarkConnectionHandshakeLifecycle-1 100 100 ns/op 1000 B/op 10 allocs/op
`
	const candidateInput = `commit: head
workload-changed: true
BenchmarkConnectionHandshakeLifecycle-1 100 100 ns/op 1000 B/op 10 allocs/op
`
	baseline, err := parseReport(strings.NewReader(baselineInput))
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := parseReport(strings.NewReader(candidateInput))
	if err != nil {
		t.Fatal(err)
	}
	comparison, err := compareReports(baseline, candidate, 1, 5, false)
	if err != nil {
		t.Fatal(err)
	}
	if !comparison.failed {
		t.Fatal("unapproved workload source change passed")
	}
	comparison, err = compareReports(baseline, candidate, 1, 5, true)
	if err != nil {
		t.Fatal(err)
	}
	if comparison.failed {
		t.Fatal("approved workload source change failed without metric regressions")
	}
	if len(comparison.items) != 0 {
		t.Fatalf("approved workload change produced incomparable metric rows: %+v", comparison.items)
	}
	var output bytes.Buffer
	if err = writeComparisonReport(&output, comparison, reportLanguages["zh-CN"]); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "结论：`PASS`") {
		t.Fatalf("passing comparison report does not contain PASS:\n%s", output.String())
	}
	if strings.Contains(output.String(), "| Benchmark |") {
		t.Fatalf("incomparable workload report contains a metric table:\n%s", output.String())
	}
}

func TestCompareReportsRejectsIncompletePairs(t *testing.T) {
	tests := []struct {
		name      string
		baseline  string
		candidate string
		want      string
	}{
		{
			name:      "missing benchmark",
			baseline:  "BenchmarkConnectionHandshakeLifecycle-1 1 100 ns/op 1000 B/op 10 allocs/op\n",
			candidate: "BenchmarkECHHandshakeLifecycle/Direct-1 1 100 ns/op 1000 B/op 10 allocs/op\n",
			want:      "missing baseline benchmark",
		},
		{
			name:      "unequal sample count",
			baseline:  "BenchmarkConnectionHandshakeLifecycle-1 1 100 ns/op 1000 B/op 10 allocs/op\nBenchmarkConnectionHandshakeLifecycle-1 1 100 ns/op 1000 B/op 10 allocs/op\n",
			candidate: "BenchmarkConnectionHandshakeLifecycle-1 1 100 ns/op 1000 B/op 10 allocs/op\n",
			want:      "sample count differs",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			baseline, err := parseReport(strings.NewReader(test.baseline))
			if err != nil {
				t.Fatal(err)
			}
			candidate, err := parseReport(strings.NewReader(test.candidate))
			if err != nil {
				t.Fatal(err)
			}
			_, err = compareReports(baseline, candidate, 1, 5, false)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("compareReports error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestStableResourceIncrease(t *testing.T) {
	baseline := []float64{1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000, 1000}
	if stableResourceIncrease(baseline, []float64{1001, 1001, 1001, 1001, 1001, 1001, 999, 999, 999}) {
		t.Fatal("six of nine higher resource samples were treated as stable")
	}
	if !stableResourceIncrease(baseline, []float64{1001, 1001, 1001, 1001, 1001, 1001, 1001, 999, 999}) {
		t.Fatal("seven of nine higher resource samples were not treated as stable")
	}
}

func TestBenchmarkLabel(t *testing.T) {
	tests := []struct {
		language string
		input    string
		want     string
	}{
		{"en", "BenchmarkConnectionHandshakeLifecycle-2", "Certificate-authenticated full handshake / AES-128-GCM"},
		{"en", "BenchmarkSessionTicketRequestHandshakeLifecycle-2", "Full handshake + 4 acknowledged session tickets"},
		{"zh-CN", "BenchmarkGREASEHandshakeLifecycle/Disabled-2", "完整 mTLS 握手 + 会话票据 / GREASE 关闭"},
		{"en", "BenchmarkGREASEHandshakeLifecycle/Enabled-2", "Full mTLS handshake + session ticket / GREASE enabled"},
		{"ru", "BenchmarkGREASEHandshakeLifecycle/Enabled-2", "Полное рукопожатие mTLS + билет сеанса / GREASE включен"},
		{"zh-CN", "BenchmarkWolfSSLFeatureRealUDP/GREASE/GoClient/WolfSSLServer-2", "GREASE 兼容性 / 完整 mTLS 握手 + 会话票据"},
		{"en", "BenchmarkMutualTLSHandshakeLifecycle/Resumed-2", "mTLS session resumption handshake"},
		{"zh-CN", "BenchmarkCertificateSelectionHandshakeLifecycle/InitialMutualTLS-2", "按 CA 与 OID filters 选择多证书的 mTLS 握手"},
		{"ru", "BenchmarkCertificateSelectionHandshakeLifecycle/PostHandshakeAuthentication-2", "Выбор сертификата для аутентификации после рукопожатия"},
		{"en", "BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/WolfSSLClient/GoServer-2", "Post-quantum hybrid key exchange / X25519MLKEM768"},
		{"en", "BenchmarkWolfSSLFeatureRealUDP/MutualTLSSessionResumption/GoClient/WolfSSLServer-2", "mTLS session resumption handshake"},
		{"zh-CN", "BenchmarkCertificateCompressionHandshakeLifecycle/MutualTLS/Zlib-2", "zlib mTLS 证书压缩握手"},
		{"zh-CN", "BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer-2", "KeyUpdate + 应用数据 1-RTT 往返"},
		{"ru", "BenchmarkWolfSSLFeatureRealUDP/ApplicationDataRoundTrip/WolfSSLClient/GoServer-2", "Обмен прикладными данными 1-RTT"},
		{"en", "BenchmarkProtectedRecordRoundTripSuites/ChaCha20Poly1305-2", "Record round trip / ChaCha20-Poly1305"},
		{"en", "BenchmarkListenerDatagramRoundTrip/Connections1024/NoCID/Bytes1200-4", "Listener UDP round trip / 1024 connections / No CID / 1200 bytes"},
		{"en", "BenchmarkMemoryDatagramRoundTrip/Bytes64-4", "In-memory datagram round trip / 64 bytes"},
		{"en", "BenchmarkListenerCIDRouting/Connections1024/Fixed/Hit-1", "Listener CID routing / 1024 connections / Fixed length / Hit"},
		{"en", "BenchmarkListenerCIDRouting/Connections252/Wide/Miss-1", "Listener CID routing / 252 connections / Lengths 4–255 / Miss"},
		{"en", "BenchmarkListenerCIDRouting/Connections1024/Long/Hit-1", "Listener CID routing / 1024 connections / 255-byte CIDs / Hit"},
		{"zh-CN", "BenchmarkListenerCIDRouting/Connections64/Mixed/Miss-1", "Listener CID 路由 / 64 个连接 / 混合长度 / 未命中"},
		{"ru", "BenchmarkListenerCIDRouting/Connections1/Fixed/Hit-1", "Маршрутизация CID через Listener / 1 соединение / Фиксированная длина / Совпадение"},
		{"zh-CN", "BenchmarkListenerDatagramRoundTrip/Connections64/CID/Bytes64-4", "Listener UDP 往返 / 64 个连接 / CID / 64 字节"},
		{"ru", "BenchmarkEarlyDataHandshakeLifecycle/Entries4096-4", "0-RTT / 4096 начальных записей повторов"},
		{"zh-CN", "BenchmarkProtectedRecordRoundTripSuites/ChaCha20Poly1305-2", "记录往返 / ChaCha20-Poly1305"},
		{"ru", "BenchmarkBuildProtectedACKRecords/SingleReuse-2", "Построение защищенного ACK / Повторное использование одного"},
	}
	for _, test := range tests {
		if got := benchmarkLabel(test.input, reportLanguages[test.language]); got != test.want {
			t.Errorf("benchmarkLabel(%q, %q) = %q, want %q", test.input, test.language, got, test.want)
		}
	}
}

func TestLocalizedReports(t *testing.T) {
	report, err := parseReport(strings.NewReader("BenchmarkProtectedRecordSeal-2 1 100 ns/op 10 B/op 1 allocs/op\nBenchmarkBuildProtectedACKRecords/Single-2 1 100 ns/op 10 B/op 1 allocs/op\nBenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer-2 1 100 ns/op 1 go_ms/conn\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		language string
		want     []string
	}{
		{"zh-CN", []string{"# 自动化基准测试结果", "## 快速跳转", "  - [go-dtls 客户端 -> wolfSSL 服务端 (1)](#real-udp-go-dtls-client-wolfssl-server)", "## 真实 UDP 互通", "### go-dtls 客户端 -> wolfSSL 服务端", "中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。", "KeyUpdate + 应用数据 1-RTT 往返", "## 记录层与可靠性", "受保护 ACK 构建 / 单条", "| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |"}},
		{"ru", []string{"# Результаты автоматических бенчмарков", "## Быстрая навигация", "  - [Клиент go-dtls -> сервер wolfSSL (1)](#real-udp-go-dtls-client-wolfssl-server)", "## Совместимость через реальный UDP", "### Клиент go-dtls -> сервер wolfSSL", "Медианное время измеряет клиент go-dtls; `ms/conn` означает одну полную операцию соединения.", "KeyUpdate + обмен прикладными данными 1-RTT", "## Уровень записей и надежность", "Построение защищенного ACK / Один", "| Бенчмарк | Замеры | Медианное время | Память стенда | Аллокации стенда |"}},
	} {
		var output bytes.Buffer
		if err := writeReport(&output, report, "", "", "", reportLanguages[test.language], ""); err != nil {
			t.Fatal(err)
		}
		for _, want := range test.want {
			if !strings.Contains(output.String(), want) {
				t.Errorf("%s report does not contain %q:\n%s", test.language, want, output.String())
			}
		}
	}
}

func TestReportLanguageLinks(t *testing.T) {
	tests := map[string]string{
		"README.md":            "[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)",
		"benchmark.en.md":      "[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)",
		"benchmark.ru.md":      "[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)",
		"benchmark-preview.md": "[简体中文](benchmark-preview.md) | [English](benchmark-preview.en.md) | [Русский](benchmark-preview.ru.md)",
	}
	for output, want := range tests {
		if got := reportLanguageLinks(output); got != want {
			t.Errorf("reportLanguageLinks(%q) = %q, want %q", output, got, want)
		}
	}
}

func TestBenchmarkSection(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"BenchmarkConnectionHandshakeLifecycle", 0},
		{"BenchmarkListenerDatagramRoundTrip/Connections1/CID/Bytes64", 0},
		{"BenchmarkListenerCIDRouting/Connections1024/Mixed/Miss", 2},
		{"BenchmarkEarlyDataHandshakeLifecycle/Entries64", 0},
		{"BenchmarkCertificateSelectionHandshakeLifecycle/InitialMutualTLS", 0},
		{"BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/GoServer", 1},
		{"BenchmarkProtectedRecordSeal", 2},
		{"BenchmarkSomethingElse", 3},
	}
	for _, test := range tests {
		if got := benchmarkSection(test.name); got != test.want {
			t.Errorf("benchmarkSection(%q) = %d, want %d", test.name, got, test.want)
		}
	}
}

func TestBenchmarkDisplayOrder(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"BenchmarkConnectionHandshakeLifecycle-2", 10},
		{"BenchmarkSessionTicketRequestHandshakeLifecycle-2", 35},
		{"BenchmarkGREASEHandshakeLifecycle/Disabled-2", 38},
		{"BenchmarkGREASEHandshakeLifecycle/Enabled-2", 39},
		{"BenchmarkMutualTLSHandshakeLifecycle/Resumed-2", 30},
		{"BenchmarkCertificateSelectionHandshakeLifecycle/InitialMutualTLS-2", 32},
		{"BenchmarkCertificateSelectionHandshakeLifecycle/PostHandshakeAuthentication-2", 33},
		{"BenchmarkCertificateCompressionHandshakeLifecycle/MutualTLS/Zlib-2", 80},
		{"BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/GoClient/WolfSSLServer-2", 10},
		{"BenchmarkWolfSSLFeatureRealUDP/GREASE/GoClient/WolfSSLServer-2", 35},
		{"BenchmarkWolfSSLFeatureRealUDP/EarlyData/WolfSSLClient/GoServer-2", 110},
		{"BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/WolfSSLClient/WolfSSLServer-2", 120},
		{"BenchmarkUnknown-2", 1000},
	}
	for _, test := range tests {
		if got := benchmarkDisplayOrder(test.name); got != test.want {
			t.Errorf("benchmarkDisplayOrder(%q) = %d, want %d", test.name, got, test.want)
		}
	}
}

func TestRealUDPDirection(t *testing.T) {
	tests := map[string]int{
		"BenchmarkFeature/GoClient/GoServer":                                         0,
		"BenchmarkFeature/GoClient/WolfSSLServer":                                    1,
		"BenchmarkFeature/WolfSSLClient/GoServer":                                    2,
		"BenchmarkFeature/WolfSSLClient/WolfSSLServer":                               3,
		"BenchmarkFeature":                                                           -1,
		"BenchmarkDTLS13PeerRealUDP/NSS/ALPN/GoClient/GoServer-2":                    4,
		"BenchmarkDTLS13PeerRealUDP/NSS/ALPN/NSSClient/NSSServer-2":                  7,
		"BenchmarkDTLS13PeerRealUDP/BORINGSSL/ECH/BORINGSSLClient/GoServer-2":        10,
		"BenchmarkDTLS13PeerRealUDP/OPENSSL/AES128CCM/OPENSSLClient/OPENSSLServer-2": 15,
	}
	for name, want := range tests {
		if got := realUDPDirection(name); got != want {
			t.Errorf("realUDPDirection(%q) = %d, want %d", name, got, want)
		}
	}
}

func TestValidatePeerRealUDPMatrix(t *testing.T) {
	report := completePeerRealUDPReport(t)
	if err := report.validatePeerRealUDPMatrix(); err != nil {
		t.Fatal(err)
	}
	want := 0
	for _, features := range peerRealUDPWorkloads {
		want += len(features) * 4
	}
	if len(report.benchmarks) != want {
		t.Fatalf("peer matrix entries=%d want=%d", len(report.benchmarks), want)
	}
	for _, item := range report.benchmarks {
		if err := validateBenchmarkCoverage(item.name); err != nil {
			t.Fatal(err)
		}
	}
	for _, test := range []struct{ language, label, metadata, skip string }{
		{"en", "NSS client -> NSS server", "- NSS: `" + peerRealUDPCommits["NSS"] + "`", "last verified against NSS commit "},
		{"zh-CN", "BoringSSL 客户端 -> BoringSSL 服务端", "- BoringSSL: `" + peerRealUDPCommits["BORINGSSL"] + "`", "该限制最后验证于 BoringSSL commit "},
		{"ru", "Клиент OpenSSL -> сервер OpenSSL", "- OpenSSL: `" + peerRealUDPCommits["OPENSSL"] + "`", "последнее подтверждение для OpenSSL commit "},
	} {
		var output bytes.Buffer
		if err := writeReport(&output, report, "", "", "", reportLanguages[test.language], ""); err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{test.label, test.metadata, test.skip, "2 ms/conn", "2 ms/pair", "linux/amd64, Example CPU", "`abc123`"} {
			if !strings.Contains(output.String(), want) {
				t.Errorf("%s report missing %q", test.language, want)
			}
		}
	}
}

func TestValidatePeerRealUDPMatrixRejectsIncompleteEvidence(t *testing.T) {
	name := "BenchmarkDTLS13PeerRealUDP/NSS/CertificateAES128GCM/GoClient/GoServer"
	for _, test := range []struct {
		name   string
		change func(*report)
	}{
		{"missing result", func(r *report) { removeRealUDPResult(r, name) }},
		{"missing peer SHA", func(r *report) { delete(r.peerCommits, "NSS") }},
		{"wrong peer SHA", func(r *report) { r.peerCommits["NSS"] = strings.Repeat("0", 40) }},
		{"unregistered skip", func(r *report) {
			removeRealUDPResult(r, name)
			r.realUDPSkips[name], r.realUDPSamples[name] = "unsupported", 1
		}},
		{"recovered skip", func(r *report) { r.addBenchmarkForTest("BenchmarkDTLS13PeerRealUDP/NSS/EarlyData/GoClient/NSSServer") }},
		{"missing sample", func(r *report) { r.byName[name+"-2"].samples-- }},
	} {
		t.Run(test.name, func(t *testing.T) {
			report := completePeerRealUDPReport(t)
			test.change(report)
			if err := report.validatePeerRealUDPMatrix(); err == nil {
				t.Fatal("accepted incomplete peer matrix")
			}
		})
	}
}

func completePeerRealUDPReport(t *testing.T) *report {
	t.Helper()
	var input bytes.Buffer
	input.WriteString("commit: abc123\ngo version go1.27.0 linux/amd64\n")
	encoder := json.NewEncoder(&input)
	encode := func(event goTestEvent) {
		t.Helper()
		if err := encoder.Encode(event); err != nil {
			t.Fatal(err)
		}
	}
	encode(goTestEvent{Action: "output", Output: "goos: linux\ngoarch: amd64\ncpu: Example CPU\n"})
	for _, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		encode(goTestEvent{Action: "output", Output: strings.ToLower(peer) + ": " + peerRealUDPCommits[peer] + "\n"})
		for _, feature := range peerRealUDPWorkloads[peer] {
			for _, direction := range []string{"GoClient/GoServer", "GoClient/" + peer + "Server", peer + "Client/GoServer", peer + "Client/" + peer + "Server"} {
				name := "BenchmarkDTLS13PeerRealUDP/" + peer + "/" + feature + "/" + direction
				if allowance, ok := realUDPSkipAllowances[name]; ok {
					encode(goTestEvent{Action: "output", Test: name, Output: "    peer_benchmark_test.go:1: " + allowance.output + "\n"})
					encode(goTestEvent{Action: "skip", Test: name})
				} else {
					unit := "go_ms/conn"
					if strings.HasPrefix(direction, peer+"Client") {
						unit = "peer_process_ms/conn"
					}
					if strings.Contains(feature, "Resumption") || feature == "EarlyData" {
						unit = strings.TrimSuffix(unit, "/conn") + "/pair"
					}
					for range realUDPMatrixSamples {
						encode(goTestEvent{Action: "output", Test: name, Output: name + "-2\t"})
						encode(goTestEvent{Action: "output", Output: "1 1000000 ns/op 2 " + unit + " 10 B/op 1 allocs/op\n"})
					}
				}
			}
		}
	}
	r, err := parseReport(&input)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestParseRealUDPJSON(t *testing.T) {
	const skipped = "BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/GoClient/WolfSSLServer"
	const result = "BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/GoClient/GoServer-2"
	var input bytes.Buffer
	input.WriteString("BenchmarkBuildPlainFlight-2 1 1 ns/op 1 B/op 1 allocs/op\n")
	encoder := json.NewEncoder(&input)
	for _, event := range []goTestEvent{
		{Action: "output", Output: "BenchmarkBuildPlainFlight-2\t"},
		{Action: "output", Output: "1 3 ns/op 3 B/op 3 allocs/op\n"},
		{Action: "output", Test: skipped, Output: "    interop_test.go:1079: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096\n"},
		{Action: "skip", Test: skipped},
		{Action: "output", Test: result, Output: result + "\t"},
		{Action: "output", Output: "1 100 ns/op 10 B/op 1 allocs/op\n"},
	} {
		if err := encoder.Encode(event); err != nil {
			t.Fatal(err)
		}
	}
	report, err := parseReport(&input)
	if err != nil {
		t.Fatal(err)
	}
	if !report.realUDPJSON || report.realUDPSkips[skipped] != "wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096" || report.realUDPSamples[skipped] != 1 {
		t.Fatalf("structured skip was not preserved: %#v", report)
	}
	if len(report.benchmarks) != 2 || report.byName[result] == nil || report.byName[result].samples != 1 {
		t.Fatalf("structured benchmark result was not preserved: %#v", report.benchmarks)
	}
	if item := report.byName["BenchmarkBuildPlainFlight-2"]; item == nil || item.samples != 2 || item.timeMetric() != "2 ns/op" {
		t.Fatalf("mixed plain and JSON samples were not preserved: %#v", item)
	}
}

func TestParseReportRejectsConflictingPeerMetadata(t *testing.T) {
	for _, peer := range []string{"nss", "boringssl", "openssl"} {
		t.Run(peer, func(t *testing.T) {
			var input bytes.Buffer
			input.WriteString(peer + ": " + strings.Repeat("a", 40) + "\n")
			encoder := json.NewEncoder(&input)
			if err := encoder.Encode(goTestEvent{Action: "output", Output: peer + ": " + strings.Repeat("b", 40) + "\n"}); err != nil {
				t.Fatal(err)
			}
			if _, err := parseReport(&input); err == nil || !strings.Contains(err.Error(), "conflicts with") {
				t.Fatalf("conflicting peer metadata error = %v", err)
			}
		})
	}
}

func TestValidateRealUDPMatrix(t *testing.T) {
	if workloads := expectedRealUDPWorkloads(); len(workloads) != 15 {
		t.Fatalf("real UDP workloads = %d, want 15", len(workloads))
	}
	if len(realUDPSkipAllowances) != 14 {
		t.Fatalf("real UDP skip allowances = %d, want 14", len(realUDPSkipAllowances))
	}
	report := completeRealUDPReport()
	if err := report.validateRealUDPMatrix(); err != nil {
		t.Fatal(err)
	}
	if len(report.benchmarks) != 15*len(realUDPDirectionPaths) {
		t.Fatalf("matrix entries = %d, want %d", len(report.benchmarks), 15*len(realUDPDirectionPaths))
	}
	var output bytes.Buffer
	if err := writeReport(&output, report, "", "", "", reportLanguages["zh-CN"], ""); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲") {
		t.Fatalf("report does not expose the allowlisted reason:\n%s", output.String())
	}
	if !strings.Contains(output.String(), "该限制最后验证于 wolfSSL commit "+reviewedWolfSSLCommit) {
		t.Fatalf("report does not expose the last verified wolfSSL commit:\n%s", output.String())
	}
	for language, want := range map[string]string{
		"en": "last verified against wolfSSL commit " + reviewedWolfSSLCommit,
		"ru": "последнее подтверждение для wolfSSL commit " + reviewedWolfSSLCommit,
	} {
		output.Reset()
		if err := writeReport(&output, report, "", "", "", reportLanguages[language], ""); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), want) {
			t.Fatalf("%s report does not expose the last verified wolfSSL commit:\n%s", language, output.String())
		}
	}
}

func TestReadReportFileAppliesWolfSSLMetadataBeforeMatrixValidation(t *testing.T) {
	const futureCommit = "0000000000000000000000000000000000000000"
	path := filepath.Join(t.TempDir(), "real-udp.json")
	file, err := os.Create(path) // #nosec G304 -- path is inside the test-owned temporary directory.
	if err != nil {
		t.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	for _, workload := range expectedRealUDPWorkloads() {
		for _, direction := range realUDPDirectionPaths {
			name := workload + "/" + direction
			if allowance, ok := realUDPSkipAllowances[name]; ok {
				if err = encoder.Encode(goTestEvent{Action: "output", Test: name, Output: "    interop_test.go:1: " + allowance.output + "\n"}); err == nil {
					err = encoder.Encode(goTestEvent{Action: "skip", Test: name})
				}
			} else {
				result := name + "-2"
				for range realUDPMatrixSamples {
					if err = encoder.Encode(goTestEvent{Action: "output", Test: result, Output: result + "\t"}); err == nil {
						err = encoder.Encode(goTestEvent{Action: "output", Output: "1 100 ns/op 10 B/op 1 allocs/op\n"})
					}
					if err != nil {
						break
					}
				}
			}
			if err != nil {
				_ = file.Close()
				t.Fatal(err)
			}
		}
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}

	report, err := readReportFile(path, futureCommit+" (future build)", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.benchmarks) != len(expectedRealUDPWorkloads())*len(realUDPDirectionPaths) {
		t.Fatalf("matrix entries = %d, want %d", len(report.benchmarks), len(expectedRealUDPWorkloads())*len(realUDPDirectionPaths))
	}
	unsupported := 0
	for _, item := range report.benchmarks {
		if item.unsupported[0] != "" {
			unsupported++
			if !strings.Contains(item.unsupported[0], reviewedWolfSSLCommit) {
				t.Fatalf("unsupported reason omitted last verified commit: %q", item.unsupported[0])
			}
		}
	}
	if unsupported != 4 {
		t.Fatalf("unsupported wolfSSL entries = %d, want 4", unsupported)
	}
}

func TestReadReportFileRejectsConflictingWolfSSLMetadata(t *testing.T) {
	path := filepath.Join(t.TempDir(), "benchmark.txt")
	if err := os.WriteFile(path, []byte("wolfssl: "+reviewedWolfSSLCommit+" (input build)\nBenchmarkProtectedRecordSeal-2 1 100 ns/op\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readReportFile(path, strings.Repeat("0", 40)+" (explicit build)", false); err == nil || !strings.Contains(err.Error(), "does not match explicit value") {
		t.Fatalf("conflicting wolfSSL metadata error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsMissingDirection(t *testing.T) {
	report := completeRealUDPReport()
	name := "BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/GoClient/GoServer"
	removeRealUDPResult(report, name)
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "produced no benchmark result") {
		t.Fatalf("missing direction error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsUnregisteredSkip(t *testing.T) {
	report := completeRealUDPReport()
	name := "BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/GoClient/GoServer"
	removeRealUDPResult(report, name)
	report.realUDPSkips[name] = "unexpected peer limitation"
	report.realUDPSamples[name] = realUDPMatrixSamples
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "unregistered skip reason") {
		t.Fatalf("unregistered skip error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsChangedSkipReason(t *testing.T) {
	report := completeRealUDPReport()
	name := "BenchmarkWolfSSLFeatureRealUDP/EarlyData/GoClient/WolfSSLServer"
	report.realUDPSkips[name] = "different reason"
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "skip reason") {
		t.Fatalf("changed skip error = %v", err)
	}
}

func TestValidateRealUDPMatrixAcceptsNewWolfSSLCommit(t *testing.T) {
	report := completeRealUDPReport()
	report.wolfSSL = "0000000000000000000000000000000000000000 (Linux Release static)"
	if err := report.validateRealUDPMatrix(); err != nil {
		t.Fatalf("new wolfSSL commit: %v", err)
	}
}

func TestValidateCurrentWolfSSLSkipVerification(t *testing.T) {
	report := completeRealUDPReport()
	if err := report.validateCurrentWolfSSLSkipVerification(); err != nil {
		t.Fatalf("current wolfSSL skip verification: %v", err)
	}

	const name = "BenchmarkWolfSSLFeatureRealUDP/EarlyData/GoClient/WolfSSLServer"
	original := realUDPSkipAllowances[name]
	allowance := original
	allowance.verifiedCommit = "0000000000000000000000000000000000000000"
	realUDPSkipAllowances[name] = allowance
	t.Cleanup(func() {
		realUDPSkipAllowances[name] = original
	})
	if err := report.validateCurrentWolfSSLSkipVerification(); err == nil || !strings.Contains(err.Error(), "was last verified against wolfSSL commit") {
		t.Fatalf("stale wolfSSL skip verification error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsMalformedWolfSSLCommit(t *testing.T) {
	report := completeRealUDPReport()
	report.wolfSSL = "not-a-commit (Linux Release static)"
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "is not a full SHA") {
		t.Fatalf("malformed wolfSSL commit error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsMissingWolfSSLMetadata(t *testing.T) {
	report := completeRealUDPReport()
	report.wolfSSL = ""
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "has no wolfSSL metadata") {
		t.Fatalf("missing wolfSSL metadata error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsIncompleteSkipAllowance(t *testing.T) {
	const name = "BenchmarkWolfSSLFeatureRealUDP/EarlyData/GoClient/WolfSSLServer"
	original := realUDPSkipAllowances[name]
	allowance := original
	allowance.verifiedCommit = "short"
	realUDPSkipAllowances[name] = allowance
	t.Cleanup(func() {
		realUDPSkipAllowances[name] = original
	})
	if err := completeRealUDPReport().validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "is incomplete") {
		t.Fatalf("incomplete allowance error = %v", err)
	}
}

func TestValidateRealUDPMatrixRejectsRecoveredAllowlistedDirection(t *testing.T) {
	report := completeRealUDPReport()
	name := "BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/GoClient/WolfSSLServer"
	delete(report.realUDPSkips, name)
	delete(report.realUDPSamples, name)
	report.addBenchmarkForTest(name)
	if err := report.validateRealUDPMatrix(); err == nil || !strings.Contains(err.Error(), "now produces a result") {
		t.Fatalf("recovered direction error = %v", err)
	}
}

func completeRealUDPReport() *report {
	report := &report{
		wolfSSL:        reviewedWolfSSLCommit + " (Linux Release static)",
		byName:         make(map[string]*benchmark),
		realUDPJSON:    true,
		realUDPSkips:   make(map[string]string),
		realUDPSamples: make(map[string]int),
	}
	for _, workload := range expectedRealUDPWorkloads() {
		for _, direction := range realUDPDirectionPaths {
			name := workload + "/" + direction
			if allowance, ok := realUDPSkipAllowances[name]; ok {
				report.realUDPSkips[name] = allowance.output
				report.realUDPSamples[name] = realUDPSkipSamples
				continue
			}
			report.addBenchmarkForTest(name)
		}
	}
	return report
}

func (r *report) addBenchmarkForTest(name string) {
	timing := &metric{values: []float64{100, 100, 100, 100, 100}}
	item := &benchmark{name: name + "-2", samples: realUDPMatrixSamples, byUnit: map[string]*metric{"ns/op": timing}}
	r.byName[item.name] = item
	r.benchmarks = append(r.benchmarks, item)
}

func removeRealUDPResult(report *report, name string) {
	for index, item := range report.benchmarks {
		if realUDPBenchmarkName(item.name) != name {
			continue
		}
		report.benchmarks = append(report.benchmarks[:index], report.benchmarks[index+1:]...)
		delete(report.byName, item.name)
		return
	}
}
