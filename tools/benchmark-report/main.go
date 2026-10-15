package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type metric struct {
	values []float64
}

type benchmark struct {
	name        string
	samples     int
	byUnit      map[string]*metric
	unsupported [3]string
}

type report struct {
	goVersion        string
	goos             string
	goarch           string
	cpu              string
	wolfSSL          string
	peerCommits      map[string]string
	commit           string
	workloadChanged  bool
	benchmarks       []*benchmark
	byName           map[string]*benchmark
	realUDPJSON      bool
	realUDPPending   map[string]string
	realUDPSkips     map[string]string
	realUDPSamples   map[string]int
	goTestOutput     string
	goTestOutputTest string
}

type benchmarkComparison struct {
	name            string
	pairs           int
	baselineTime    float64
	candidateTime   float64
	pairedTimeDelta float64
	worseTimePairs  int
	baselineBytes   float64
	candidateBytes  float64
	baselineAllocs  float64
	candidateAllocs float64
	failures        []string
}

type comparisonReport struct {
	baselineCommit        string
	candidateCommit       string
	timeRegressionPercent float64
	workloadChanged       bool
	workloadChangeAllowed bool
	failed                bool
	items                 []benchmarkComparison
}

type reportLanguage struct {
	labelIndex      int
	title           string
	summary         string
	details         string
	quickNavigation string
	metadata        [5]string
	sections        [4]string
	realUDP         [4]string
	peerUDP         [4]string
	peerMetadata    [3]string
	realUDPTiming   [4]string
	columns         [6]string
	unsupported     string
}

var sectionAnchors = [...]string{
	"section-connection-lifecycle",
	"section-real-udp-interoperability",
	"section-record-layer-and-reliability",
	"section-other-benchmarks",
}

var realUDPAnchors = [...]string{
	"real-udp-go-dtls-client-go-dtls-server",
	"real-udp-go-dtls-client-wolfssl-server",
	"real-udp-wolfssl-client-go-dtls-server",
	"real-udp-wolfssl-client-wolfssl-server",
}

var realUDPDirectionPaths = [...]string{
	"GoClient/GoServer",
	"GoClient/WolfSSLServer",
	"WolfSSLClient/GoServer",
	"WolfSSLClient/WolfSSLServer",
}

var peerRealUDPWorkloads = map[string][]string{
	"NSS":       {"CertificateAES128GCM", "AES256GCM", "ChaCha20Poly1305", "P256", "P384", "HelloRetryRequest", "X25519MLKEM768", "FragmentedHandshake", "MutualTLS", "ALPN", "KeyUpdate", "OCSPStapling", "ExternalPSK", "SessionResumption", "MutualTLSSessionResumption", "EarlyData", "DelegatedCredential", "DelegatedCredentialResumption"},
	"BORINGSSL": {"CertificateAES128GCM", "AES256GCM", "ChaCha20Poly1305", "P256", "P384", "HelloRetryRequest", "X25519MLKEM768", "FragmentedHandshake", "MutualTLS", "ALPN", "KeyUpdate", "OCSPStapling", "RawPublicKey", "ImportedPSK", "SessionResumption", "MutualTLSSessionResumption", "ECH", "DelegatedCredential", "DelegatedCredentialResumption"},
	"OPENSSL":   {"CertificateAES128GCM", "AES256GCM", "ChaCha20Poly1305", "AES128CCM", "P256", "P384", "HelloRetryRequest", "X25519MLKEM768", "FragmentedHandshake", "MutualTLS", "ALPN", "KeyUpdate", "OCSPStapling", "RawPublicKey", "ExternalPSK", "SessionResumption", "MutualTLSSessionResumption", "EarlyData", "PostHandshakeAuthentication"},
}

var peerRealUDPCommits = map[string]string{
	"NSS":       "b9283c1e4ee232b4450d36115272af670d38ec81",
	"BORINGSSL": "4f3b183c5d4a8f8fcaba4c0ab4435df181df6978",
	"OPENSSL":   "41eecc7c7cd4069b1b29bca49b9d97825c6bf648",
}

const (
	realUDPMatrixSamples  = 5
	realUDPSkipSamples    = 1
	reviewedWolfSSLCommit = "7499fc5b6c99eb39055a59f538bb3709971341c9"
)

type realUDPSkipAllowance struct {
	output         string
	reason         [3]string
	evidence       string
	verifiedCommit string
	verifiedPeer   string
}

var realUDPSkipAllowances = map[string]realUDPSkipAllowance{
	"BenchmarkDTLS13PeerRealUDP/NSS/FragmentedHandshake/GoClient/NSSServer": {
		output:   "NSS fixture cannot enforce a fragmented handshake without a configurable DTLS MTU",
		reason:   [3]string{"NSS fixture cannot force handshake fragmentation with an NSS server", "NSS 夹具无法强制 NSS 服务端的握手分片", "Стенд NSS не может принудительно фрагментировать рукопожатие с сервером NSS"},
		evidence: "Go ClientHello at MTU 256 is unfragmented and NSS exposes no public DTLS MTU setter", verifiedCommit: peerRealUDPCommits["NSS"], verifiedPeer: "NSS",
	},
	"BenchmarkDTLS13PeerRealUDP/NSS/FragmentedHandshake/NSSClient/NSSServer": {
		output:   "NSS fixture cannot enforce a fragmented handshake without a configurable DTLS MTU",
		reason:   [3]string{"NSS fixture cannot force handshake fragmentation for two NSS peers", "NSS 夹具无法强制两个 NSS 对端的握手分片", "Стенд NSS не может принудительно фрагментировать рукопожатие двух узлов NSS"},
		evidence: "NSS exposes no DTLS MTU setter in the pinned public API", verifiedCommit: peerRealUDPCommits["NSS"], verifiedPeer: "NSS",
	},
	"BenchmarkDTLS13PeerRealUDP/NSS/EarlyData/GoClient/NSSServer": {
		output:   "NSS server does not support DTLS 1.3 early data in this fixture",
		reason:   [3]string{"NSS server cannot accept DTLS 1.3 early data in this fixture", "此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据", "Сервер NSS в этом стенде не принимает ранние данные DTLS 1.3"},
		evidence: "testdata/interop/nss.c rejects native server early data", verifiedCommit: peerRealUDPCommits["NSS"], verifiedPeer: "NSS",
	},
	"BenchmarkDTLS13PeerRealUDP/NSS/EarlyData/NSSClient/NSSServer": {
		output:   "NSS server does not support DTLS 1.3 early data in this fixture",
		reason:   [3]string{"NSS server cannot accept DTLS 1.3 early data in this fixture", "此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据", "Сервер NSS в этом стенде не принимает ранние данные DTLS 1.3"},
		evidence: "testdata/interop/nss.c rejects native server early data", verifiedCommit: peerRealUDPCommits["NSS"], verifiedPeer: "NSS",
	},
	"BenchmarkDTLS13PeerRealUDP/OPENSSL/EarlyData/GoClient/OPENSSLServer": {
		output:   "OPENSSL server does not support DTLS 1.3 early data in this fixture",
		reason:   [3]string{"OpenSSL listener cannot enter the DTLS early-data read state", "OpenSSL 监听器无法进入 DTLS 早期数据读取状态", "Слушатель OpenSSL не может перейти к чтению ранних данных DTLS"},
		evidence: "testdata/interop/openssl.c rejects native server early data", verifiedCommit: peerRealUDPCommits["OPENSSL"], verifiedPeer: "OPENSSL",
	},
	"BenchmarkDTLS13PeerRealUDP/OPENSSL/EarlyData/OPENSSLClient/OPENSSLServer": {
		output:   "OPENSSL server does not support DTLS 1.3 early data in this fixture",
		reason:   [3]string{"OpenSSL listener cannot enter the DTLS early-data read state", "OpenSSL 监听器无法进入 DTLS 早期数据读取状态", "Слушатель OpenSSL не может перейти к чтению ранних данных DTLS"},
		evidence: "testdata/interop/openssl.c rejects native server early data", verifiedCommit: peerRealUDPCommits["OPENSSL"], verifiedPeer: "OPENSSL",
	},
	"BenchmarkDTLS13PeerRealUDP/BORINGSSL/DelegatedCredential/BORINGSSLClient/GoServer": {
		output:   "BoringSSL client does not authenticate delegated credentials in DTLS 1.3",
		reason:   [3]string{"BoringSSL client does not authenticate delegated credentials in DTLS 1.3", "BoringSSL 客户端不认证 DTLS 1.3 委派凭证", "Клиент BoringSSL не проверяет делегированные учётные данные в DTLS 1.3"},
		evidence: "TestInteropDTLS13Peers excludes BoringSSL client delegated-credential mode", verifiedCommit: peerRealUDPCommits["BORINGSSL"], verifiedPeer: "BORINGSSL",
	},
	"BenchmarkDTLS13PeerRealUDP/BORINGSSL/DelegatedCredential/BORINGSSLClient/BORINGSSLServer": {
		output:   "BoringSSL client does not authenticate delegated credentials in DTLS 1.3",
		reason:   [3]string{"BoringSSL client does not authenticate delegated credentials in DTLS 1.3", "BoringSSL 客户端不认证 DTLS 1.3 委派凭证", "Клиент BoringSSL не проверяет делегированные учётные данные в DTLS 1.3"},
		evidence: "TestInteropDTLS13Peers excludes BoringSSL client delegated-credential mode", verifiedCommit: peerRealUDPCommits["BORINGSSL"], verifiedPeer: "BORINGSSL",
	},
	"BenchmarkDTLS13PeerRealUDP/BORINGSSL/DelegatedCredentialResumption/BORINGSSLClient/GoServer": {
		output:   "BoringSSL client does not authenticate delegated credentials in DTLS 1.3",
		reason:   [3]string{"BoringSSL client does not authenticate delegated credentials in DTLS 1.3", "BoringSSL 客户端不认证 DTLS 1.3 委派凭证", "Клиент BoringSSL не проверяет делегированные учётные данные в DTLS 1.3"},
		evidence: "TestInteropDTLS13Peers excludes BoringSSL client delegated-credential mode", verifiedCommit: peerRealUDPCommits["BORINGSSL"], verifiedPeer: "BORINGSSL",
	},
	"BenchmarkDTLS13PeerRealUDP/BORINGSSL/DelegatedCredentialResumption/BORINGSSLClient/BORINGSSLServer": {
		output:   "BoringSSL client does not authenticate delegated credentials in DTLS 1.3",
		reason:   [3]string{"BoringSSL client does not authenticate delegated credentials in DTLS 1.3", "BoringSSL 客户端不认证 DTLS 1.3 委派凭证", "Клиент BoringSSL не проверяет делегированные учётные данные в DTLS 1.3"},
		evidence: "TestInteropDTLS13Peers excludes BoringSSL client delegated-credential mode", verifiedCommit: peerRealUDPCommits["BORINGSSL"], verifiedPeer: "BORINGSSL",
	},
	"BenchmarkWolfSSLFeatureRealUDP/MutualTLSSessionResumption/WolfSSLClient/GoServer": {
		output: "wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket",
		reason: [3]string{
			"wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket",
			"wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello",
			"Клиент wolfSSL отказывается фрагментировать первый ClientHello с билетом mTLS",
		},
		evidence:       "TestInteropWolfSSLClientMutualTLSSessionResumption: wolfSSL_connect resume error -328, malformed buffer input error",
		verifiedCommit: reviewedWolfSSLCommit,
	},
	"BenchmarkWolfSSLFeatureRealUDP/EarlyData/GoClient/WolfSSLServer": {
		output: "wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest",
		reason: [3]string{
			"this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT",
			"该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT",
			"В этом сценарии включён cookie HRR, требующий отклонить 0-RTT go-dtls",
		},
		evidence:       "RFC 9846 sections 4.2.2 and 4.3.10; the benchmark server uses cookie HRR",
		verifiedCommit: reviewedWolfSSLCommit,
	},
	"BenchmarkWolfSSLFeatureRealUDP/EarlyData/WolfSSLClient/WolfSSLServer": {
		output: "wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest",
		reason: [3]string{
			"this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT",
			"该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT",
			"В этом сценарии включён cookie HRR, требующий отклонить 0-RTT клиента wolfSSL",
		},
		evidence:       "RFC 9846 sections 4.2.2 and 4.3.10; the benchmark server uses cookie HRR",
		verifiedCommit: reviewedWolfSSLCommit,
	},
	"BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/GoClient/WolfSSLServer": {
		output: "wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096",
		reason: [3]string{
			"wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096",
			"MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲",
			"При MTU 4096 ClientHello P384 превышает приёмные буферы сервера wolfSSL",
		},
		evidence:       "wolfSSL default receive buffer is 1900 bytes; the cookie-bearing P384 ClientHello at MTU 4096 is 1913 bytes (Windows MSG_PEEK is limited to 1500 bytes)",
		verifiedCommit: reviewedWolfSSLCommit,
	},
}

var reportLanguages = map[string]reportLanguage{
	"en": {
		labelIndex:      0,
		title:           "Automated benchmark results",
		summary:         "%d results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.\n",
		details:         "Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per operation; Listener UDP sums the separate client/server processes and excludes the relay. Exact raw output remains in the workflow artifact.",
		quickNavigation: "Quick navigation",
		metadata:        [5]string{"Commit", "Generated", "Go", "Platform", "wolfSSL"},
		sections: [4]string{
			"Connection lifecycle",
			"Real UDP interoperability",
			"Record layer and reliability",
			"Other benchmarks",
		},
		realUDP: [4]string{
			"go-dtls client -> go-dtls server",
			"go-dtls client -> wolfSSL server",
			"wolfSSL client -> go-dtls server",
			"wolfSSL client -> wolfSSL server",
		},
		peerUDP:      [4]string{"go-dtls client -> go-dtls server", "go-dtls client -> %s server", "%s client -> go-dtls server", "%s client -> %s server"},
		peerMetadata: [3]string{"NSS", "BoringSSL", "OpenSSL"},
		realUDPTiming: [4]string{
			"Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.",
			"Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.",
			"Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.",
			"Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.",
		},
		columns:     [6]string{"Benchmark", "Samples", "Median time", "Throughput", "Harness memory", "Harness allocations"},
		unsupported: "Skipped",
	},
	"zh-CN": {
		labelIndex:      1,
		title:           "自动化基准测试结果",
		summary:         "共 %d 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。\n",
		details:         "工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数按每次操作统计；Listener UDP 汇总独立客户端和服务端进程，排除中转器。精确原始输出保留在 Workflow Artifact 中。",
		quickNavigation: "快速跳转",
		metadata:        [5]string{"提交", "生成时间", "Go", "平台", "wolfSSL"},
		sections: [4]string{
			"连接生命周期",
			"真实 UDP 互通",
			"记录层与可靠性",
			"其他基准测试",
		},
		realUDP: [4]string{
			"go-dtls 客户端 -> go-dtls 服务端",
			"go-dtls 客户端 -> wolfSSL 服务端",
			"wolfSSL 客户端 -> go-dtls 服务端",
			"wolfSSL 客户端 -> wolfSSL 服务端",
		},
		peerUDP:      [4]string{"go-dtls 客户端 -> go-dtls 服务端", "go-dtls 客户端 -> %s 服务端", "%s 客户端 -> go-dtls 服务端", "%s 客户端 -> %s 服务端"},
		peerMetadata: [3]string{"NSS", "BoringSSL", "OpenSSL"},
		realUDPTiming: [4]string{
			"中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。",
			"中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。",
			"中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。",
			"中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。",
		},
		columns:     [6]string{"基准测试", "样本数", "中位耗时", "吞吐量", "测试框架内存", "测试框架分配次数"},
		unsupported: "跳过",
	},
	"ru": {
		labelIndex:      2,
		title:           "Результаты автоматических бенчмарков",
		summary:         "%d результатов сгруппированы по нагрузке и упорядочены по функции, затем по имени бенчмарка. Значения являются медианами выборок последнего запуска.\n",
		details:         "Специализированные метрики соединения имеют приоритет над временем стенда Go. Память и аллокации указаны на операцию; Listener UDP суммирует отдельные процессы клиента и сервера без ретранслятора. Точный исходный вывод сохранен в артефакте Workflow.",
		quickNavigation: "Быстрая навигация",
		metadata:        [5]string{"Коммит", "Сформировано", "Go", "Платформа", "wolfSSL"},
		sections: [4]string{
			"Жизненный цикл соединения",
			"Совместимость через реальный UDP",
			"Уровень записей и надежность",
			"Другие бенчмарки",
		},
		realUDP: [4]string{
			"Клиент go-dtls -> сервер go-dtls",
			"Клиент go-dtls -> сервер wolfSSL",
			"Клиент wolfSSL -> сервер go-dtls",
			"Клиент wolfSSL -> сервер wolfSSL",
		},
		peerUDP:      [4]string{"Клиент go-dtls -> сервер go-dtls", "Клиент go-dtls -> сервер %s", "Клиент %s -> сервер go-dtls", "Клиент %s -> сервер %s"},
		peerMetadata: [3]string{"NSS", "BoringSSL", "OpenSSL"},
		realUDPTiming: [4]string{
			"Медианное время измеряет клиент go-dtls; `ms/conn` означает одну полную операцию соединения.",
			"Медианное время измеряет клиент wolfSSL; `ms/conn` означает одно соединение, а `ms/pair` означает пару из полного и возобновленного соединений с учетом встроенного ожидания клиента.",
			"Медианное время включает запуск процесса стороннего клиента и обмен аутентифицированными датаграммами; `ms/pair` включает полное и возобновленное соединения. Это время нельзя напрямую сравнивать со временем только клиента.",
			"Медианное время измеряет клиент go-dtls и включает обмен аутентифицированными датаграммами; `ms/pair` включает полное соединение для получения билета. При сервере OpenSSL ACK клиента эпохи 2 намеренно отбрасывается.",
		},
		columns:     [6]string{"Бенчмарк", "Замеры", "Медианное время", "Пропускная способность", "Память стенда", "Аллокации стенда"},
		unsupported: "Пропущено",
	},
}

var workloadMetrics = [...]struct {
	unit  string
	label string
}{
	{"go_ms/conn", "ms/conn"},
	{"go_ms/pair", "ms/pair"},
	{"wolfssl_ms/conn", "ms/conn"},
	{"wolf_process_ms/conn", "ms/conn"},
	{"wolf_process_ms/pair", "ms/pair"},
	{"peer_process_ms/conn", "ms/conn"},
	{"peer_process_ms/pair", "ms/pair"},
}

var benchmarkPartLabels = map[string]string{
	"1301":                        "AES-128-GCM",
	"1302":                        "AES-256-GCM",
	"AES128CCM":                   "AES-128-CCM",
	"AES128GCM":                   "AES-128-GCM",
	"AES256GCM":                   "AES-256-GCM",
	"ApplicationDataRoundTrip":    "Application data",
	"BuildPlainACKRecords":        "Plain ACK build",
	"BuildProtectedACKRecords":    "Protected ACK build",
	"CertificateAES128GCM":        "Certificate / AES-128-GCM",
	"ChaCha20Poly1305":            "ChaCha20-Poly1305",
	"ConnectionID":                "CID",
	"EarlyData":                   "0-RTT",
	"HelloRetryRequest":           "HRR",
	"HybridKeyExchange":           "Hybrid KEX",
	"MutualTLS":                   "mTLS",
	"MutualTLSSessionResumption":  "mTLS resumption",
	"PostHandshakeAuthentication": "Post-handshake authentication",
	"ProtectedRecordReceiveErrorUnauthenticated": "Reject unauthenticated record",
	"ProtectedRecordRoundTrip":                   "Record round trip",
	"ProtectedRecordSeal":                        "Record seal",
	"SecP256r1MLKEM768":                          "SecP256r1MLKEM768",
	"SecP384r1MLKEM1024":                         "SecP384r1MLKEM1024",
	"ServerCertificate":                          "Server cert",
	"SessionResumption":                          "Session resumption",
}

var benchmarkFeatureLabels = map[string][3]string{
	"ALPN":                          {"ALPN-negotiated connection + datagram exchange", "协商 ALPN 的连接 + 数据报交换", "Соединение с ALPN и обмен датаграммами"},
	"HelloRetryRequest":             {"HelloRetryRequest handshake + datagram exchange", "HelloRetryRequest 握手 + 数据报交换", "Рукопожатие с HelloRetryRequest и обмен датаграммами"},
	"ImportedPSK":                   {"Imported external PSK handshake + datagram exchange", "导入式外部 PSK 握手 + 数据报交换", "Рукопожатие с импортированным внешним PSK и обмен датаграммами"},
	"DelegatedCredentialResumption": {"Delegated-credential handshake + session resumption", "委派凭证握手 + 会话恢复", "Рукопожатие с делегированными учётными данными и возобновление сеанса"},
	"AES256GCM":                     {"AES-256-GCM/SHA-384 connection + datagram exchange", "AES-256-GCM/SHA-384 连接 + 数据报交换", "Соединение AES-256-GCM/SHA-384 и обмен датаграммами"},
	"ChaCha20Poly1305":              {"ChaCha20-Poly1305 connection + datagram exchange", "ChaCha20-Poly1305 连接 + 数据报交换", "Соединение ChaCha20-Poly1305 и обмен датаграммами"},
	"P256":                          {"P-256 key exchange + datagram exchange", "P-256 密钥交换 + 数据报交换", "Обмен ключами P-256 и датаграммами"},
	"P384":                          {"P-384 key exchange + datagram exchange", "P-384 密钥交换 + 数据报交换", "Обмен ключами P-384 и датаграммами"},
	"X25519MLKEM768":                {"X25519MLKEM768 hybrid handshake + datagram exchange", "X25519MLKEM768 混合握手 + 数据报交换", "Гибридное рукопожатие X25519MLKEM768 и обмен датаграммами"},
	"FragmentedHandshake":           {"MTU-256 fragmented handshake + datagram exchange", "MTU 256 分片握手 + 数据报交换", "Фрагментированное рукопожатие при MTU 256 и обмен датаграммами"},
	"OCSPStapling":                  {"OCSP-stapled handshake + datagram exchange", "OCSP 装订握手 + 数据报交换", "Рукопожатие с OCSP-ответом и обмен датаграммами"},
	"RawPublicKey":                  {"Raw-public-key handshake + datagram exchange", "原始公钥握手 + 数据报交换", "Рукопожатие с открытым ключом и обмен датаграммами"},
	"ECH":                           {"ECH-accepted handshake + datagram exchange", "ECH 接受握手 + 数据报交换", "Рукопожатие с принятым ECH и обмен датаграммами"},
	"DelegatedCredential":           {"Delegated-credential handshake + datagram exchange", "委派凭证握手 + 数据报交换", "Рукопожатие с делегированными данными и обмен датаграммами"},
	"PeerCertificateAES128GCM":      {"Certificate-authenticated AES-128-GCM connection + datagram exchange", "证书认证 AES-128-GCM 连接 + 数据报交换", "Соединение AES-128-GCM с сертификатом и обмен датаграммами"},
	"AES128CCM": {
		"Certificate-authenticated full handshake / AES-128-CCM",
		"证书认证完整握手 / AES-128-CCM",
		"Полное рукопожатие с сертификатом / AES-128-CCM",
	},
	"ApplicationDataRoundTrip": {
		"1-RTT application-data round trip",
		"应用数据 1-RTT 往返",
		"Обмен прикладными данными 1-RTT",
	},
	"CertificateAES128GCM": {
		"Certificate-authenticated full handshake / AES-128-GCM",
		"证书认证完整握手 / AES-128-GCM",
		"Полное рукопожатие с сертификатом / AES-128-GCM",
	},
	"CertificateSelectionInitialMutualTLS": {
		"Multi-certificate mTLS selection by CA and OID filters",
		"按 CA 与 OID filters 选择多证书的 mTLS 握手",
		"Выбор сертификата mTLS по УЦ и фильтрам OID",
	},
	"CertificateSelectionPostHandshakeAuthentication": {
		"Multi-certificate post-handshake authentication selection",
		"握手后认证的多证书选择",
		"Выбор сертификата для аутентификации после рукопожатия",
	},
	"CertificateCompressionMutualTLSPlain": {
		"Full mTLS handshake / uncompressed certificates",
		"完整 mTLS 握手 / 证书未压缩",
		"Полное рукопожатие mTLS / сертификаты без сжатия",
	},
	"CertificateCompressionMutualTLSZlib": {
		"zlib-compressed mTLS handshake",
		"zlib mTLS 证书压缩握手",
		"Рукопожатие mTLS со сжатием сертификатов zlib",
	},
	"CertificateCompressionServerCertificatePlain": {
		"Full server-certificate handshake / uncompressed certificate",
		"服务器证书完整握手 / 证书未压缩",
		"Полное рукопожатие с сертификатом сервера / без сжатия",
	},
	"CertificateCompressionServerCertificateZlib": {
		"zlib-compressed server-certificate handshake",
		"zlib 服务器证书压缩握手",
		"Рукопожатие со сжатием сертификата сервера zlib",
	},
	"Connection": {
		"Certificate-authenticated full handshake / AES-128-GCM",
		"证书认证完整握手 / AES-128-GCM",
		"Полное рукопожатие с сертификатом / AES-128-GCM",
	},
	"ConnectionID": {
		"CID + 1-RTT application-data round trip",
		"CID + 应用数据 1-RTT 往返",
		"CID + обмен прикладными данными 1-RTT",
	},
	"ECHDirect": {
		"ECH handshake / direct (no HRR)",
		"ECH 握手 / 直接（无 HRR）",
		"Рукопожатие ECH / прямое (без HRR)",
	},
	"ECHHRR": {
		"ECH handshake / via HRR",
		"ECH 握手 / 经 HRR",
		"Рукопожатие ECH / через HRR",
	},
	"EarlyData": {
		"0-RTT + 1-RTT application-data round trip",
		"0-RTT + 应用数据 1-RTT 往返",
		"0-RTT + обмен прикладными данными 1-RTT",
	},
	"ExternalPSK": {
		"Direct external PSK handshake",
		"直接外部 PSK 握手",
		"Рукопожатие с прямым внешним PSK",
	},
	"GREASEDisabled": {
		"Full mTLS handshake + session ticket / GREASE disabled",
		"完整 mTLS 握手 + 会话票据 / GREASE 关闭",
		"Полное рукопожатие mTLS + билет сеанса / GREASE отключен",
	},
	"GREASEEnabled": {
		"Full mTLS handshake + session ticket / GREASE enabled",
		"完整 mTLS 握手 + 会话票据 / GREASE 启用",
		"Полное рукопожатие mTLS + билет сеанса / GREASE включен",
	},
	"GREASERealUDP": {
		"GREASE compatibility / full mTLS handshake + session ticket",
		"GREASE 兼容性 / 完整 mTLS 握手 + 会话票据",
		"Совместимость с GREASE / полное рукопожатие mTLS + билет сеанса",
	},
	"HybridKeyExchange": {
		"Post-quantum hybrid key exchange",
		"后量子混合密钥交换",
		"Гибридный постквантовый обмен ключами",
	},
	"KeyUpdate": {
		"KeyUpdate + 1-RTT application-data round trip",
		"KeyUpdate + 应用数据 1-RTT 往返",
		"KeyUpdate + обмен прикладными данными 1-RTT",
	},
	"MutualTLS": {
		"Full mTLS handshake",
		"完整 mTLS 握手",
		"Полное рукопожатие mTLS",
	},
	"MutualTLSFull": {
		"Full mTLS handshake",
		"完整 mTLS 握手",
		"Полное рукопожатие mTLS",
	},
	"MutualTLSResumed": {
		"mTLS session resumption handshake",
		"mTLS 会话恢复握手",
		"Рукопожатие с возобновлением сеанса mTLS",
	},
	"MutualTLSSessionResumption": {
		"mTLS session resumption handshake",
		"mTLS 会话恢复握手",
		"Рукопожатие с возобновлением сеанса mTLS",
	},
	"PostHandshakeAuthentication": {
		"PHA + 1-RTT application-data round trip",
		"PHA + 应用数据 1-RTT 往返",
		"PHA + обмен прикладными данными 1-RTT",
	},
	"SessionResumption": {
		"Session resumption handshake",
		"会话恢复握手",
		"Рукопожатие с возобновлением сеанса",
	},
	"SessionTicketRequest": {
		"Full handshake + 4 acknowledged session tickets",
		"完整握手 + 4 个已确认会话票据",
		"Полное рукопожатие + 4 подтвержденных билета сеанса",
	},
}

var benchmarkPartTranslations = map[string][3]string{
	"Packet Session Ingress":               {"Listener packet ingress", "Listener 报文入队", "Приём пакетов через Listener"},
	"Listener CID Routing":                 {"Listener CID routing", "Listener CID 路由", "Маршрутизация CID через Listener"},
	"Fixed":                                {"Fixed length", "固定长度", "Фиксированная длина"},
	"Mixed":                                {"Mixed lengths", "混合长度", "Разные длины"},
	"Wide":                                 {"Lengths 4–255", "4–255 字节长度", "Длины 4–255"},
	"Long":                                 {"255-byte CIDs", "255 字节 CID", "CID по 255 байт"},
	"Hit":                                  {"Hit", "命中", "Совпадение"},
	"Miss":                                 {"Miss", "未命中", "Нет совпадения"},
	"Memory Datagram Round Trip":           {"In-memory datagram round trip", "内存 datagram 往返", "Обмен датаграммами в памяти"},
	"Listener Datagram Round Trip":         {"Listener UDP round trip", "Listener UDP 往返", "Обмен UDP через Listener"},
	"No Server Read Deadline":              {"No server read deadline", "服务端无读期限", "Без срока чтения на сервере"},
	"Connections1":                         {"1 connection", "1 个连接", "1 соединение"},
	"Connections64":                        {"64 connections", "64 个连接", "64 соединения"},
	"Connections252":                       {"252 connections", "252 个连接", "252 соединения"},
	"Connections1024":                      {"1024 connections", "1024 个连接", "1024 соединения"},
	"No CID":                               {"No CID", "无 CID", "Без CID"},
	"CID":                                  {"CID", "CID", "CID"},
	"Bytes64":                              {"64 bytes", "64 字节", "64 байта"},
	"Bytes1200":                            {"1200 bytes", "1200 字节", "1200 байт"},
	"Bytes2048":                            {"2048 bytes", "2048 字节", "2048 байт"},
	"Bytes4096":                            {"4096 bytes", "4096 字节", "4096 байт"},
	"Entries64":                            {"64 initial replay entries", "初始 64 条重放记录", "64 начальные записи повторов"},
	"Entries4096":                          {"4096 initial replay entries", "初始 4096 条重放记录", "4096 начальных записей повторов"},
	"Plain ACK build":                      {"Plain ACK build", "明文 ACK 构建", "Построение открытого ACK"},
	"Protected ACK build":                  {"Protected ACK build", "受保护 ACK 构建", "Построение защищенного ACK"},
	"Build Plain Flight":                   {"Build Plain Flight", "构建明文握手报文组", "Построение открытой группы сообщений"},
	"Build Protected Flight":               {"Build Protected Flight", "构建受保护握手报文组", "Построение защищенной группы сообщений"},
	"Combine Flights":                      {"Combine Flights", "合并握手报文组", "Объединение групп сообщений"},
	"Flight First Refresh":                 {"Flight First Refresh", "握手报文组首次刷新", "Первое обновление группы сообщений"},
	"Flight Initial History Batch":         {"Flight Initial History Batch", "握手报文组初始历史批次", "Начальная история группы сообщений"},
	"Flight Wire Window":                   {"Flight Wire Window", "握手报文组传输窗口", "Окно передачи группы сообщений"},
	"Parse ACK":                            {"Parse ACK", "解析 ACK", "Разбор ACK"},
	"Protected Record CID":                 {"Protected Record CID", "受保护记录 CID", "Защищенная запись CID"},
	"Reject unauthenticated record":        {"Reject unauthenticated record", "拒绝未认证记录", "Отклонение неаутентифицированной записи"},
	"Record round trip":                    {"Record round trip", "记录往返", "Обмен записью туда и обратно"},
	"Protected Record Round Trip In Place": {"Protected Record Round Trip In Place", "受保护记录原地往返", "Обмен защищенной записью на месте"},
	"Record seal":                          {"Record seal", "记录封装", "Защита записи"},
	"Empty":                                {"Empty", "空", "Пустой"},
	"Single":                               {"Single", "单条", "Один"},
	"Sorted64":                             {"Sorted64", "已排序 64", "Отсортированный 64"},
	"Reversed64":                           {"Reversed64", "逆序 64", "Обратный 64"},
	"Single Reuse":                         {"Single Reuse", "单条复用", "Повторное использование одного"},
	"Reuse Single":                         {"Reuse Single", "单条复用", "Повторное использование одного"},
	"Retransmit":                           {"Retransmit", "重传", "Повторная передача"},
	"Round Trip":                           {"Round Trip", "往返", "Обмен туда и обратно"},
	"Seal":                                 {"Seal", "封装", "Защита"},
	"Owned":                                {"Owned", "独占", "Владеемое"},
	"Reuse":                                {"Reuse", "复用", "Повторное использование"},
	"Certificate":                          {"Certificate", "证书", "Сертификат"},
}

var benchmarkDisplayOrders = map[string]int{
	"ConnectionHandshakeLifecycle":                                       10,
	"MutualTLSHandshakeLifecycle/Full":                                   20,
	"MutualTLSHandshakeLifecycle/Resumed":                                30,
	"CertificateSelectionHandshakeLifecycle/InitialMutualTLS":            32,
	"CertificateSelectionHandshakeLifecycle/PostHandshakeAuthentication": 33,
	"SessionTicketRequestHandshakeLifecycle":                             35,
	"ExternalPSKHandshakeLifecycle":                                      40,
	"GREASEHandshakeLifecycle/Disabled":                                  38,
	"GREASEHandshakeLifecycle/Enabled":                                   39,
	"CertificateCompressionHandshakeLifecycle/ServerCertificate/Plain":   50,
	"CertificateCompressionHandshakeLifecycle/ServerCertificate/Zlib":    60,
	"CertificateCompressionHandshakeLifecycle/MutualTLS/Plain":           70,
	"CertificateCompressionHandshakeLifecycle/MutualTLS/Zlib":            80,
	"ECHHandshakeLifecycle/Direct":                                       90,
	"ECHHandshakeLifecycle/HRR":                                          100,
	"HybridKeyExchangeHandshakeLifecycle/X25519MLKEM768":                 110,
	"HybridKeyExchangeHandshakeLifecycle/SecP256r1MLKEM768":              120,
	"HybridKeyExchangeHandshakeLifecycle/SecP384r1MLKEM1024":             130,
	"ListenerDatagramRoundTrip/Connections1/NoCID/Bytes64":               140,
	"ListenerDatagramRoundTrip/Connections1/NoCID/Bytes1200":             141,
	"ListenerDatagramRoundTrip/Connections1/CID/Bytes64":                 142,
	"ListenerDatagramRoundTrip/Connections1/CID/Bytes1200":               143,
	"ListenerDatagramRoundTrip/Connections64/NoCID/Bytes64":              144,
	"ListenerDatagramRoundTrip/Connections64/NoCID/Bytes1200":            145,
	"ListenerDatagramRoundTrip/Connections64/CID/Bytes64":                146,
	"ListenerDatagramRoundTrip/Connections64/CID/Bytes1200":              147,
	"ListenerDatagramRoundTrip/Connections1024/NoCID/Bytes64":            148,
	"ListenerDatagramRoundTrip/Connections1024/NoCID/Bytes1200":          149,
	"ListenerDatagramRoundTrip/Connections1024/CID/Bytes64":              150,
	"ListenerDatagramRoundTrip/Connections1024/CID/Bytes1200":            151,
	"EarlyDataHandshakeLifecycle/Entries64":                              160,
	"EarlyDataHandshakeLifecycle/Entries4096":                            161,
	"MemoryDatagramRoundTrip/Bytes64":                                    162,
	"MemoryDatagramRoundTrip/Bytes1200":                                  163,
	"ListenerCIDRouting/Connections1/Fixed/Hit":                          180,
	"ListenerCIDRouting/Connections1/Fixed/Miss":                         181,
	"ListenerCIDRouting/Connections1/Mixed/Hit":                          182,
	"ListenerCIDRouting/Connections1/Mixed/Miss":                         183,
	"ListenerCIDRouting/Connections64/Fixed/Hit":                         184,
	"ListenerCIDRouting/Connections64/Fixed/Miss":                        185,
	"ListenerCIDRouting/Connections64/Mixed/Hit":                         186,
	"ListenerCIDRouting/Connections64/Mixed/Miss":                        187,
	"ListenerCIDRouting/Connections1024/Fixed/Hit":                       188,
	"ListenerCIDRouting/Connections1024/Fixed/Miss":                      189,
	"ListenerCIDRouting/Connections1024/Mixed/Hit":                       190,
	"ListenerCIDRouting/Connections1024/Mixed/Miss":                      191,
	"ListenerCIDRouting/Connections252/Wide/Hit":                         192,
	"ListenerCIDRouting/Connections252/Wide/Miss":                        193,
	"ListenerCIDRouting/Connections252/Long/Hit":                         194,
	"ListenerCIDRouting/Connections252/Long/Miss":                        195,
	"ListenerCIDRouting/Connections1024/Wide/Hit":                        196,
	"ListenerCIDRouting/Connections1024/Wide/Miss":                       197,
	"ListenerCIDRouting/Connections1024/Long/Hit":                        198,
	"ListenerCIDRouting/Connections1024/Long/Miss":                       199,
	"PacketSessionIngress/Bytes64":                                       200,
	"PacketSessionIngress/Bytes1200":                                     201,
	"PacketSessionIngress/Bytes2048":                                     202,
	"PacketSessionIngress/Bytes4096":                                     203,

	"ListenerDatagramRoundTrip/NoServerReadDeadline/Connections1024/CID/Bytes1200": 152,

	"WolfSSLFeatureRealUDP/CertificateAES128GCM":        10,
	"WolfSSLFeatureRealUDP/ApplicationDataRoundTrip":    20,
	"WolfSSLFeatureRealUDP/MutualTLS":                   30,
	"WolfSSLFeatureRealUDP/GREASE":                      35,
	"WolfSSLFeatureRealUDP/AES128CCM":                   40,
	"WolfSSLFeatureRealUDP/ExternalPSK":                 50,
	"WolfSSLFeatureRealUDP/ConnectionID":                60,
	"WolfSSLFeatureRealUDP/KeyUpdate":                   70,
	"WolfSSLFeatureRealUDP/PostHandshakeAuthentication": 80,
	"WolfSSLFeatureRealUDP/SessionResumption":           90,
	"WolfSSLFeatureRealUDP/MutualTLSSessionResumption":  100,
	"WolfSSLFeatureRealUDP/EarlyData":                   110,
	"HybridKeyExchangeRealUDP/X25519MLKEM768":           120,
	"HybridKeyExchangeRealUDP/SecP256r1MLKEM768":        130,
	"HybridKeyExchangeRealUDP/SecP384r1MLKEM1024":       140,
}

func main() {
	input := flag.String("input", "", "go test benchmark output")
	baseline := flag.String("baseline", "", "baseline go test benchmark output")
	candidate := flag.String("candidate", "", "candidate go test benchmark output")
	output := flag.String("output", "", "Markdown report path")
	commit := flag.String("commit", "", "tested commit")
	wolfSSL := flag.String("wolfssl", "", "wolfSSL commit and build")
	generated := flag.String("generated", time.Now().UTC().Format(time.RFC3339), "generation time")
	language := flag.String("language", "zh-CN", "report language: en, zh-CN, or ru")
	languageLinks := flag.Bool("language-links", true, "include links to sibling language reports")
	minimumPairs := flag.Int("minimum-pairs", 9, "minimum paired samples required for each benchmark")
	timeRegressionPercent := flag.Float64("time-regression-percent", 5, "paired median time regression that fails the gate")
	allowWorkloadChange := flag.Bool("allow-workload-change", false, "allow benchmark workload source changes")
	failOnRegression := flag.Bool("fail-on-regression", true, "exit unsuccessfully when the comparison fails")
	requireCurrentWolfSSLSkipVerification := flag.Bool("require-current-wolfssl-skip-verification", false, "require every real UDP skip to be verified against the report's wolfSSL commit")
	allowPartial := flag.Bool("allow-partial", false, "allow unmapped benchmarks and incomplete real UDP workloads in a single report")
	flag.Parse()
	if *output == "" {
		flag.Usage()
		os.Exit(2)
	}
	text, ok := reportLanguages[*language]
	if !ok {
		fatal(fmt.Errorf("unsupported report language %q", *language))
	}
	if *baseline != "" || *candidate != "" {
		if *baseline == "" || *candidate == "" || *input != "" || *allowPartial {
			flag.Usage()
			os.Exit(2)
		}
		baselineReport, err := readReportFile(*baseline, "", false)
		if err != nil {
			fatal(err)
		}
		candidateReport, err := readReportFile(*candidate, "", false)
		if err != nil {
			fatal(err)
		}
		comparison, err := compareReports(baselineReport, candidateReport, *minimumPairs, *timeRegressionPercent, *allowWorkloadChange)
		if err != nil {
			fatal(err)
		}
		if err = writeComparisonFile(*output, comparison, text); err != nil {
			fatal(err)
		}
		if comparison.failed && *failOnRegression {
			os.Exit(1)
		}
		return
	}
	if *input == "" {
		flag.Usage()
		os.Exit(2)
	}

	result, err := readReportFile(*input, *wolfSSL, *allowPartial)
	if err != nil {
		fatal(err)
	}
	if *requireCurrentWolfSSLSkipVerification {
		if err = result.validateCurrentWolfSSLSkipVerification(); err != nil {
			fatal(err)
		}
	}

	out, err := os.Create(*output) // #nosec G304 -- path is an explicit local command argument.
	if err != nil {
		fatal(err)
	}
	links := ""
	if *languageLinks {
		links = reportLanguageLinks(*output)
	}
	err = writeReport(out, result, *commit, *wolfSSL, *generated, text, links)
	closeErr := out.Close()
	if err != nil {
		fatal(err)
	}
	if closeErr != nil {
		fatal(closeErr)
	}
}

func readReportFile(path, wolfSSL string, allowPartial bool) (*report, error) {
	in, err := os.Open(path) // #nosec G304 -- path is an explicit local command argument.
	if err != nil {
		return nil, err
	}
	result, parseErr := parseReport(in)
	closeErr := in.Close()
	if parseErr != nil {
		return nil, parseErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if wolfSSL != "" {
		if result.wolfSSL != "" && result.wolfSSL != wolfSSL {
			return nil, fmt.Errorf("wolfSSL metadata %q does not match explicit value %q", result.wolfSSL, wolfSSL)
		}
		result.wolfSSL = wolfSSL
	}
	if allowPartial {
		for name, reason := range result.realUDPSkips {
			allowance := realUDPSkipAllowances[name]
			if allowance.output != reason {
				allowance = realUDPSkipAllowance{reason: [3]string{reason, reason, reason}}
			}
			result.addUnsupportedRealUDP(name, allowance)
		}
	} else {
		for _, item := range result.benchmarks {
			if err := validateBenchmarkCoverage(item.name); err != nil {
				return nil, err
			}
		}
		if err := result.validateRealUDPMatrix(); err != nil {
			return nil, err
		}
		if err := result.validatePeerRealUDPMatrix(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func writeComparisonFile(path string, comparison *comparisonReport, text reportLanguage) error {
	out, err := os.Create(path) // #nosec G304 -- path is an explicit local command argument.
	if err != nil {
		return err
	}
	writeErr := writeComparisonReport(out, comparison, text)
	closeErr := out.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

func fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func parseReport(reader io.Reader) (*report, error) {
	result := &report{byName: make(map[string]*benchmark)}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "{"):
			if err := result.addGoTestJSONLine(lineNumber, line); err != nil {
				return nil, err
			}
		case strings.HasPrefix(line, "Benchmark"):
			if err := result.addBenchmarkLine(lineNumber, line); err != nil {
				return nil, err
			}
		default:
			if err := result.addMetadataLine(lineNumber, line); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read benchmark output: %w", err)
	}
	if result.goTestOutput != "" {
		return nil, errors.New("go test JSON contains an incomplete output line")
	}
	if len(result.benchmarks) == 0 {
		return nil, errors.New("benchmark output contains no results")
	}
	return result, nil
}

func (r *report) addMetadataLine(lineNumber int, line string) error {
	switch {
	case strings.HasPrefix(line, "commit: ") && r.commit == "":
		r.commit = strings.TrimSpace(strings.TrimPrefix(line, "commit: "))
	case strings.HasPrefix(line, "workload-changed: "):
		changed, err := strconv.ParseBool(strings.TrimSpace(strings.TrimPrefix(line, "workload-changed: ")))
		if err != nil {
			return fmt.Errorf("line %d: parse workload change marker: %w", lineNumber, err)
		}
		r.workloadChanged = changed
	case strings.HasPrefix(line, "go version ") && r.goVersion == "":
		r.goVersion = line
	case strings.HasPrefix(line, "goos: ") && r.goos == "":
		r.goos = strings.TrimSpace(strings.TrimPrefix(line, "goos: "))
	case strings.HasPrefix(line, "goarch: ") && r.goarch == "":
		r.goarch = strings.TrimSpace(strings.TrimPrefix(line, "goarch: "))
	case strings.HasPrefix(line, "cpu: ") && r.cpu == "":
		r.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu: "))
	case strings.HasPrefix(line, "wolfssl: ") && r.wolfSSL == "":
		r.wolfSSL = strings.TrimSpace(strings.TrimPrefix(line, "wolfssl: "))
	case strings.HasPrefix(line, "nss: "):
		if err := r.setPeerCommit("NSS", strings.TrimSpace(strings.TrimPrefix(line, "nss: "))); err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}
	case strings.HasPrefix(line, "boringssl: "):
		if err := r.setPeerCommit("BORINGSSL", strings.TrimSpace(strings.TrimPrefix(line, "boringssl: "))); err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}
	case strings.HasPrefix(line, "openssl: "):
		if err := r.setPeerCommit("OPENSSL", strings.TrimSpace(strings.TrimPrefix(line, "openssl: "))); err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}
	}
	return nil
}

func (r *report) setPeerCommit(peer, value string) error {
	if r.peerCommits == nil {
		r.peerCommits = make(map[string]string)
	}
	if old := r.peerCommits[peer]; old != "" && old != value {
		return fmt.Errorf("%s metadata %q conflicts with %q", peer, old, value)
	}
	r.peerCommits[peer] = value
	return nil
}

type goTestEvent struct {
	Action string
	Test   string
	Output string
}

func (r *report) addGoTestJSONLine(lineNumber int, line string) error {
	var event goTestEvent
	if err := json.Unmarshal([]byte(line), &event); err != nil {
		return fmt.Errorf("line %d: parse go test JSON: %w", lineNumber, err)
	}
	if event.Action == "output" {
		return r.addGoTestOutput(lineNumber, event.Test, event.Output)
	}
	if !isRealUDPBenchmark(event.Test) {
		return nil
	}
	r.realUDPJSON = true
	if event.Action != "skip" {
		return nil
	}
	reason := r.realUDPPending[event.Test]
	if reason == "" {
		return fmt.Errorf("line %d: skipped real UDP benchmark %q has no structured reason", lineNumber, event.Test)
	}
	if r.realUDPSkips == nil {
		r.realUDPSkips = make(map[string]string)
		r.realUDPSamples = make(map[string]int)
	}
	if previous := r.realUDPSkips[event.Test]; previous != "" && previous != reason {
		return fmt.Errorf("line %d: skipped real UDP benchmark %q changed reason from %q to %q", lineNumber, event.Test, previous, reason)
	}
	r.realUDPSkips[event.Test] = reason
	r.realUDPSamples[event.Test]++
	return nil
}

func (r *report) addGoTestOutput(lineNumber int, test, output string) error {
	if r.goTestOutput == "" && test != "" {
		r.goTestOutputTest = test
	}
	r.goTestOutput += output
	for {
		index := strings.IndexByte(r.goTestOutput, '\n')
		if index < 0 {
			return nil
		}
		line := strings.TrimSpace(r.goTestOutput[:index])
		r.goTestOutput = r.goTestOutput[index+1:]
		lineTest := r.goTestOutputTest
		r.goTestOutputTest = ""
		if r.goTestOutput != "" && test != "" {
			r.goTestOutputTest = test
		}
		if err := r.addMetadataLine(lineNumber, line); err != nil {
			return err
		}
		if isBenchmarkResultLine(line) {
			r.realUDPJSON = r.realUDPJSON || isRealUDPBenchmark(strings.Fields(line)[0])
			if err := r.addBenchmarkLine(lineNumber, line); err != nil {
				return err
			}
		}
		if !isRealUDPBenchmark(lineTest) {
			continue
		}
		r.realUDPJSON = true
		if reason, ok := goTestSkipReason(line); ok {
			if r.realUDPPending == nil {
				r.realUDPPending = make(map[string]string)
			}
			r.realUDPPending[lineTest] = reason
		}
	}
}

func isBenchmarkResultLine(line string) bool {
	fields := strings.Fields(line)
	if len(fields) < 4 || !strings.HasPrefix(fields[0], "Benchmark") {
		return false
	}
	_, err := strconv.ParseUint(fields[1], 10, 64)
	return err == nil
}

func goTestSkipReason(output string) (string, bool) {
	source, reason, ok := strings.Cut(strings.TrimSpace(output), ": ")
	if !ok || !strings.Contains(source, ".go:") || strings.TrimSpace(reason) == "" {
		return "", false
	}
	return strings.TrimSpace(reason), true
}

func isRealUDPBenchmark(name string) bool {
	return strings.HasPrefix(name, "BenchmarkWolfSSLFeatureRealUDP") || strings.HasPrefix(name, "BenchmarkHybridKeyExchangeRealUDP") || strings.HasPrefix(name, "BenchmarkDTLS13PeerRealUDP/")
}

func (r *report) validateRealUDPMatrix() error {
	if r.wolfSSL == "" {
		wolfResult := false
		for _, item := range r.benchmarks {
			wolfResult = wolfResult || isWolfSSLRealUDP(item.name)
		}
		for name := range r.realUDPSkips {
			wolfResult = wolfResult || isWolfSSLRealUDP(name)
		}
		if wolfResult {
			return errors.New("structured real UDP matrix has no wolfSSL metadata")
		}
		return nil
	}
	fields := strings.Fields(r.wolfSSL)
	if len(fields) == 0 {
		return errors.New("wolfSSL metadata has no commit")
	}
	if !isFullCommit(fields[0]) {
		return fmt.Errorf("wolfSSL metadata commit %q is not a full SHA", fields[0])
	}
	if !r.realUDPJSON {
		return errors.New("wolfSSL benchmark report has no structured real UDP matrix")
	}

	expected := make(map[string]struct{})
	workloads := expectedRealUDPWorkloads()
	for _, workload := range workloads {
		for _, direction := range realUDPDirectionPaths {
			expected[workload+"/"+direction] = struct{}{}
		}
	}
	for name, allowance := range realUDPSkipAllowances {
		if allowance.verifiedPeer != "" {
			continue
		}
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("real UDP skip allowlist contains unknown direction %q", name)
		}
		if allowance.output == "" || allowance.evidence == "" || !isFullCommit(allowance.verifiedCommit) || allowance.reason[0] == "" || allowance.reason[1] == "" || allowance.reason[2] == "" {
			return fmt.Errorf("real UDP skip allowlist entry %q is incomplete", name)
		}
	}

	results := make(map[string]*benchmark)
	for _, item := range r.benchmarks {
		name := realUDPBenchmarkName(item.name)
		if !isWolfSSLRealUDP(name) {
			continue
		}
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("real UDP matrix produced unknown direction %q", name)
		}
		if results[name] != nil {
			return fmt.Errorf("real UDP matrix produced duplicate direction %q", name)
		}
		results[name] = item
	}
	for name := range r.realUDPSkips {
		if !isWolfSSLRealUDP(name) {
			continue
		}
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("real UDP matrix skipped unknown direction %q", name)
		}
	}

	for _, workload := range workloads {
		for _, direction := range realUDPDirectionPaths {
			name := workload + "/" + direction
			item, hasResult := results[name]
			reason, skipped := r.realUDPSkips[name]
			allowance, allowed := realUDPSkipAllowances[name]
			if hasResult && skipped {
				return fmt.Errorf("real UDP direction %q produced both a result and a skip", name)
			}
			if allowed {
				if hasResult {
					return fmt.Errorf("real UDP direction %q is allowlisted but now produces a result; revalidate wolfSSL capability", name)
				}
				if !skipped {
					return fmt.Errorf("real UDP direction %q is missing its allowlisted skip", name)
				}
				if reason != allowance.output {
					return fmt.Errorf("real UDP direction %q skip reason = %q, want %q", name, reason, allowance.output)
				}
				if samples := r.realUDPSamples[name]; samples != realUDPSkipSamples {
					return fmt.Errorf("real UDP direction %q has %d skip samples, want %d", name, samples, realUDPSkipSamples)
				}
				r.addUnsupportedRealUDP(name, allowance)
				continue
			}
			if skipped {
				return fmt.Errorf("real UDP direction %q has unregistered skip reason %q", name, reason)
			}
			if !hasResult {
				return fmt.Errorf("real UDP direction %q produced no benchmark result", name)
			}
			if item.samples != realUDPMatrixSamples {
				return fmt.Errorf("real UDP direction %q has %d samples, want %d", name, item.samples, realUDPMatrixSamples)
			}
		}
	}
	return nil
}

func isWolfSSLRealUDP(name string) bool {
	return strings.HasPrefix(name, "BenchmarkWolfSSLFeatureRealUDP") || strings.HasPrefix(name, "BenchmarkHybridKeyExchangeRealUDP")
}

func (r *report) validatePeerRealUDPMatrix() error {
	results := make(map[string]*benchmark)
	for _, item := range r.benchmarks {
		name := realUDPBenchmarkName(item.name)
		if strings.HasPrefix(name, "BenchmarkDTLS13PeerRealUDP/") {
			if results[name] != nil {
				return fmt.Errorf("peer real UDP matrix produced duplicate direction %q", name)
			}
			results[name] = item
		}
	}
	if len(results) == 0 && len(r.peerCommits) == 0 {
		for name := range r.realUDPSkips {
			if strings.HasPrefix(name, "BenchmarkDTLS13PeerRealUDP/") {
				return errors.New("peer real UDP matrix has no peer metadata")
			}
		}
		return nil
	}
	for _, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		fields := strings.Fields(r.peerCommits[peer])
		if len(fields) == 0 || !isFullCommit(fields[0]) {
			return fmt.Errorf("%s metadata has no full commit SHA", peer)
		}
		if fields[0] != peerRealUDPCommits[peer] {
			return fmt.Errorf("%s benchmark commit %q differs from pinned %q", peer, fields[0], peerRealUDPCommits[peer])
		}
	}
	expected := make(map[string]struct{})
	for _, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		for _, workload := range peerRealUDPWorkloads[peer] {
			for _, direction := range []string{"GoClient/GoServer", "GoClient/" + peer + "Server", peer + "Client/GoServer", peer + "Client/" + peer + "Server"} {
				expected["BenchmarkDTLS13PeerRealUDP/"+peer+"/"+workload+"/"+direction] = struct{}{}
			}
		}
	}
	for name, allowance := range realUDPSkipAllowances {
		if allowance.verifiedPeer == "" {
			continue
		}
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("peer real UDP skip allowlist contains unknown direction %q", name)
		}
		if allowance.output == "" || allowance.evidence == "" || allowance.verifiedCommit != peerRealUDPCommits[allowance.verifiedPeer] || allowance.reason[0] == "" || allowance.reason[1] == "" || allowance.reason[2] == "" {
			return fmt.Errorf("peer real UDP skip allowlist entry %q is incomplete", name)
		}
	}
	for name := range results {
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("peer real UDP matrix produced unknown direction %q", name)
		}
	}
	for name := range r.realUDPSkips {
		if !strings.HasPrefix(name, "BenchmarkDTLS13PeerRealUDP/") {
			continue
		}
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("peer real UDP matrix skipped unknown direction %q", name)
		}
	}
	for name := range expected {
		item, hasResult := results[name]
		reason, skipped := r.realUDPSkips[name]
		allowance, allowed := realUDPSkipAllowances[name]
		if hasResult && skipped {
			return fmt.Errorf("peer real UDP direction %q produced both a result and a skip", name)
		}
		if allowed {
			if hasResult || !skipped || reason != allowance.output || r.realUDPSamples[name] != realUDPSkipSamples {
				return fmt.Errorf("peer real UDP direction %q does not match its verified skip", name)
			}
			r.addUnsupportedRealUDP(name, allowance)
		} else if skipped || !hasResult || item.samples != realUDPMatrixSamples {
			return fmt.Errorf("peer real UDP direction %q is missing %d result samples or has an unregistered skip", name, realUDPMatrixSamples)
		}
	}
	return nil
}

func (r *report) validateCurrentWolfSSLSkipVerification() error {
	fields := strings.Fields(r.wolfSSL)
	if len(fields) == 0 {
		return errors.New("wolfSSL metadata has no commit")
	}
	for name, allowance := range realUDPSkipAllowances {
		if allowance.verifiedPeer != "" {
			continue
		}
		if allowance.verifiedCommit != fields[0] {
			return fmt.Errorf("real UDP skip allowance %q was last verified against wolfSSL commit %q, want current commit %q", name, allowance.verifiedCommit, fields[0])
		}
	}
	return nil
}

func isFullCommit(commit string) bool {
	if len(commit) != 40 {
		return false
	}
	_, err := hex.DecodeString(commit)
	return err == nil
}

func realUDPBenchmarkName(name string) string {
	if before, after, ok := strings.CutLast(name, "-"); ok {
		if _, err := strconv.ParseUint(after, 10, 64); err == nil {
			return before
		}
	}
	return name
}

func expectedRealUDPWorkloads() []string {
	workloads := make([]string, 0, len(benchmarkDisplayOrders))
	for name := range benchmarkDisplayOrders {
		if strings.Contains(name, "RealUDP") {
			workloads = append(workloads, "Benchmark"+name)
		}
	}
	sort.Slice(workloads, func(left, right int) bool {
		return benchmarkDisplayOrder(workloads[left]) < benchmarkDisplayOrder(workloads[right])
	})
	return workloads
}

func (r *report) addUnsupportedRealUDP(name string, allowance realUDPSkipAllowance) {
	reason := allowance.reason
	if allowance.verifiedCommit != "" {
		peer := allowance.verifiedPeer
		if peer == "" {
			peer = "wolfSSL"
		}
		switch peer {
		case "BORINGSSL":
			peer = "BoringSSL"
		case "OPENSSL":
			peer = "OpenSSL"
		}
		reason[0] += "; last verified against " + peer + " commit " + allowance.verifiedCommit
		reason[1] += "；该限制最后验证于 " + peer + " commit " + allowance.verifiedCommit
		reason[2] += "; последнее подтверждение для " + peer + " commit " + allowance.verifiedCommit
	}
	item := &benchmark{name: name, samples: r.realUDPSamples[name], byUnit: make(map[string]*metric), unsupported: reason}
	r.byName[name] = item
	r.benchmarks = append(r.benchmarks, item)
}

func compareReports(baseline, candidate *report, minimumPairs int, timeRegressionPercent float64, allowWorkloadChange bool) (*comparisonReport, error) {
	if minimumPairs < 1 {
		return nil, errors.New("minimum paired samples must be positive")
	}
	if timeRegressionPercent < 0 {
		return nil, errors.New("time regression percentage must not be negative")
	}
	result := &comparisonReport{
		baselineCommit:        baseline.commit,
		candidateCommit:       candidate.commit,
		timeRegressionPercent: timeRegressionPercent,
		workloadChanged:       candidate.workloadChanged,
		workloadChangeAllowed: allowWorkloadChange,
	}
	if result.workloadChanged {
		result.failed = !result.workloadChangeAllowed
		return result, nil
	}
	candidateByName, err := canonicalBenchmarks(candidate)
	if err != nil {
		return nil, err
	}
	for _, baselineItem := range baseline.benchmarks {
		name := canonicalBenchmarkName(baselineItem.name)
		candidateItem := candidateByName[name]
		if candidateItem == nil {
			return nil, fmt.Errorf("candidate output is missing baseline benchmark %q", name)
		}
		item, compareErr := compareBenchmark(baselineItem, candidateItem, minimumPairs, timeRegressionPercent)
		if compareErr != nil {
			return nil, fmt.Errorf("compare %s: %w", name, compareErr)
		}
		if len(item.failures) != 0 {
			result.failed = true
		}
		result.items = append(result.items, item)
	}
	sort.Slice(result.items, func(left, right int) bool {
		leftOrder := benchmarkDisplayOrder(result.items[left].name)
		rightOrder := benchmarkDisplayOrder(result.items[right].name)
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		return result.items[left].name < result.items[right].name
	})
	return result, nil
}

func canonicalBenchmarks(result *report) (map[string]*benchmark, error) {
	benchmarks := make(map[string]*benchmark, len(result.benchmarks))
	for _, item := range result.benchmarks {
		name := canonicalBenchmarkName(item.name)
		if benchmarks[name] != nil {
			return nil, fmt.Errorf("benchmark output contains duplicate canonical name %q", name)
		}
		benchmarks[name] = item
	}
	return benchmarks, nil
}

func canonicalBenchmarkName(name string) string {
	return "Benchmark" + strings.Join(benchmarkNameParts(name), "/")
}

func compareBenchmark(baseline, candidate *benchmark, minimumPairs int, timeRegressionPercent float64) (benchmarkComparison, error) {
	item := benchmarkComparison{name: canonicalBenchmarkName(baseline.name)}
	baselineTime, candidateTime, err := pairedMetric(baseline, candidate, "ns/op", minimumPairs)
	if err != nil {
		return item, err
	}
	baselineBytes, candidateBytes, err := pairedMetric(baseline, candidate, "B/op", minimumPairs)
	if err != nil {
		return item, err
	}
	baselineAllocs, candidateAllocs, err := pairedMetric(baseline, candidate, "allocs/op", minimumPairs)
	if err != nil {
		return item, err
	}
	item.pairs = len(baselineTime)
	item.baselineTime = median(baselineTime)
	item.candidateTime = median(candidateTime)
	item.baselineBytes = median(baselineBytes)
	item.candidateBytes = median(candidateBytes)
	item.baselineAllocs = median(baselineAllocs)
	item.candidateAllocs = median(candidateAllocs)
	deltas := make([]float64, len(baselineTime))
	for index := range baselineTime {
		if baselineTime[index] <= 0 {
			return item, errors.New("baseline time must be positive")
		}
		deltas[index] = (candidateTime[index] - baselineTime[index]) * 100 / baselineTime[index]
		if candidateTime[index] > baselineTime[index] {
			item.worseTimePairs++
		}
	}
	item.pairedTimeDelta = median(deltas)
	if item.pairedTimeDelta > timeRegressionPercent && item.worseTimePairs > item.pairs/2 {
		item.failures = append(item.failures, fmt.Sprintf("time +%.3f%%", item.pairedTimeDelta))
	}
	if item.candidateBytes > item.baselineBytes && stableResourceIncrease(baselineBytes, candidateBytes) {
		item.failures = append(item.failures, fmt.Sprintf("B/op %s -> %s", formatNumber(item.baselineBytes), formatNumber(item.candidateBytes)))
	}
	if item.candidateAllocs > item.baselineAllocs && stableResourceIncrease(baselineAllocs, candidateAllocs) {
		item.failures = append(item.failures, fmt.Sprintf("allocs/op %s -> %s", formatNumber(item.baselineAllocs), formatNumber(item.candidateAllocs)))
	}
	return item, nil
}

func stableResourceIncrease(baseline, candidate []float64) bool {
	worse := 0
	for index := range baseline {
		if candidate[index] > baseline[index] {
			worse++
		}
	}
	return worse*4 >= len(baseline)*3
}

func pairedMetric(baseline, candidate *benchmark, unit string, minimumPairs int) ([]float64, []float64, error) {
	baselineMetric := baseline.byUnit[unit]
	candidateMetric := candidate.byUnit[unit]
	if baselineMetric == nil || candidateMetric == nil {
		return nil, nil, fmt.Errorf("missing %s metric", unit)
	}
	if len(baselineMetric.values) != len(candidateMetric.values) {
		return nil, nil, fmt.Errorf("%s sample count differs: baseline=%d candidate=%d", unit, len(baselineMetric.values), len(candidateMetric.values))
	}
	if len(baselineMetric.values) < minimumPairs {
		return nil, nil, fmt.Errorf("%s has %d paired samples, need at least %d", unit, len(baselineMetric.values), minimumPairs)
	}
	return baselineMetric.values, candidateMetric.values, nil
}

func writeComparisonReport(writer io.Writer, comparison *comparisonReport, text reportLanguage) error {
	decision := "PASS"
	if comparison.failed {
		decision = "FAIL"
	}
	if _, err := fmt.Fprintln(writer, "## PR base/head 配对性能门禁"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "\n- 结论：`%s`\n- Base：`%s`\n- Head：`%s`\n- 耗时自动失败阈值：`+%s%%`，且多数配对同向\n- 资源自动失败条件：中位数增加，且至少 75%% 配对同向\n", decision, markdownText(comparison.baselineCommit), markdownText(comparison.candidateCommit), formatNumber(comparison.timeRegressionPercent)); err != nil {
		return err
	}
	if comparison.workloadChanged {
		status := "未批准，门禁失败"
		if comparison.workloadChangeAllowed {
			status = "已由维护者批准；新旧 workload 不可直接比较，不生成差值结论"
		}
		if _, err := fmt.Fprintf(writer, "- Benchmark workload 源码发生变化：%s\n", status); err != nil {
			return err
		}
	}
	if len(comparison.items) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(writer, "\n| Benchmark | Base 中位耗时 | Head 中位耗时 | 配对中位差 | 较慢配对 | B/op（Base -> Head） | allocs/op（Base -> Head） | 判定 |\n| --- | ---: | ---: | ---: | :---: | :---: | :---: | :---: |"); err != nil {
		return err
	}
	for _, item := range comparison.items {
		status := "通过"
		if len(item.failures) != 0 {
			status = strings.Join(item.failures, "；")
		}
		if _, err := fmt.Fprintf(writer, "| %s | %s | %s | %+.3f%% | %d/%d | %s -> %s | %s -> %s | %s |\n",
			markdownText(benchmarkLabel(item.name, text)), formatComparisonTime(item.baselineTime), formatComparisonTime(item.candidateTime), item.pairedTimeDelta,
			item.worseTimePairs, item.pairs, formatNumber(item.baselineBytes), formatNumber(item.candidateBytes),
			formatNumber(item.baselineAllocs), formatNumber(item.candidateAllocs), markdownText(status)); err != nil {
			return err
		}
	}
	return nil
}

func formatComparisonTime(value float64) string {
	switch {
	case value >= 1e9:
		return formatDecimal(value/1e9) + " s/op"
	case value >= 1e6:
		return formatDecimal(value/1e6) + " ms/op"
	case value >= 1e3:
		return formatDecimal(value/1e3) + " us/op"
	default:
		return formatNumber(value) + " ns/op"
	}
}

func (r *report) addBenchmarkLine(lineNumber int, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return fmt.Errorf("line %d: malformed benchmark result %q", lineNumber, line)
	}
	if _, err := strconv.ParseUint(fields[1], 10, 64); err != nil {
		return fmt.Errorf("line %d: parse iteration count: %w", lineNumber, err)
	}
	if len(fields[2:])%2 != 0 {
		return fmt.Errorf("line %d: metric value/unit pairs are incomplete", lineNumber)
	}

	item := r.byName[fields[0]]
	if item == nil {
		item = &benchmark{name: fields[0], byUnit: make(map[string]*metric)}
		r.byName[item.name] = item
		r.benchmarks = append(r.benchmarks, item)
	}
	item.samples++
	for index := 2; index < len(fields); index += 2 {
		value, err := strconv.ParseFloat(fields[index], 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("line %d: invalid metric value %q", lineNumber, fields[index])
		}
		unit := fields[index+1]
		entry := item.byUnit[unit]
		if entry == nil {
			entry = new(metric)
			item.byUnit[unit] = entry
		}
		entry.values = append(entry.values, value)
	}
	return nil
}

func validateBenchmarkCoverage(name string) error {
	if benchmarkSection(name) == len(sectionAnchors)-1 {
		return fmt.Errorf("benchmark %q is not covered by the report generator (no mapped section)", name)
	}
	if benchmarkDisplayOrder(name) != 1000 {
		return nil
	}

	parts := benchmarkNameParts(name)
	if len(parts) == 0 {
		return errors.New("benchmark has no report label")
	}
	switch parts[0] {
	case "GREASEHandshakeLifecycle":
		return fmt.Errorf("benchmark %q is not covered by the report generator", name)
	default:
		parts[0] = strings.TrimSuffix(parts[0], "HandshakeLifecycle")
		parts[0] = strings.TrimSuffix(parts[0], "Suites")
	}
	for _, part := range parts {
		if _, ok := benchmarkPartLabels[part]; ok {
			continue
		}
		if _, ok := benchmarkPartTranslations[humanizeBenchmarkPart(part)]; !ok {
			return fmt.Errorf("benchmark %q is not covered by the report generator (unmapped part %q)", name, part)
		}
	}
	return nil
}

func reportLanguageLinks(output string) string {
	switch filepath.Base(output) {
	case "README.md", "benchmark.en.md", "benchmark.ru.md":
		return "[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)"
	}
	name := strings.TrimSuffix(filepath.Base(output), filepath.Ext(output))
	name = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".zh-CN"), ".en"), ".ru")
	return fmt.Sprintf("[简体中文](%s.md) | [English](%s.en.md) | [Русский](%s.ru.md)", name, name, name)
}

func writeReport(writer io.Writer, result *report, commit, wolfSSL, generated string, text reportLanguage, languageLinks string) error {
	if commit == "" {
		commit = result.commit
	}
	if wolfSSL == "" {
		wolfSSL = result.wolfSSL
	}
	sections := make([][]*benchmark, len(sectionAnchors))
	for _, item := range result.benchmarks {
		index := benchmarkSection(item.name)
		sections[index] = append(sections[index], item)
	}
	if _, err := fmt.Fprintf(writer, "# %s\n", text.title); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	if languageLinks != "" {
		if _, err := fmt.Fprintln(writer, languageLinks); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}
	platform := result.goos + "/" + result.goarch
	if result.goos == "" && result.goarch == "" {
		platform = ""
	}
	if platform != "" && result.cpu != "" {
		platform += ", "
	}
	platform += result.cpu
	metadata := [][2]string{
		{text.metadata[0], commit},
		{text.metadata[1], generated},
		{text.metadata[2], result.goVersion},
		{text.metadata[3], platform},
		{text.metadata[4], wolfSSL},
	}
	for index, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		metadata = append(metadata, [2]string{text.peerMetadata[index], result.peerCommits[peer]})
	}
	for _, item := range metadata {
		if item[1] == "" {
			continue
		}
		if _, err := fmt.Fprintf(writer, "- %s: `%s`\n", item[0], markdownText(item[1])); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, text.summary, len(result.benchmarks)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, text.details); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(writer, "\n## %s\n", text.quickNavigation); err != nil {
		return err
	}
	for index, items := range sections {
		if len(items) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(writer, "\n- [%s (%d)](#%s)", text.sections[index], len(items), sectionAnchors[index]); err != nil {
			return err
		}
		if index == 1 {
			groups, err := groupRealUDP(items)
			if err != nil {
				return err
			}
			for direction, group := range groups {
				if len(group) == 0 {
					continue
				}
				if _, err := fmt.Fprintf(writer, "\n  - [%s (%d)](#%s)", realUDPDirectionLabel(text, direction), len(group), realUDPAnchor(direction)); err != nil {
					return err
				}
			}
		}
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	for index, items := range sections {
		if len(items) == 0 {
			continue
		}
		sort.Slice(items, func(left, right int) bool {
			leftOrder := benchmarkDisplayOrder(items[left].name)
			rightOrder := benchmarkDisplayOrder(items[right].name)
			if leftOrder != rightOrder {
				return leftOrder < rightOrder
			}
			return items[left].name < items[right].name
		})
		if _, err := fmt.Fprintf(writer, "\n<a id=\"%s\"></a>\n## %s\n\n", sectionAnchors[index], text.sections[index]); err != nil {
			return err
		}
		if index == 1 {
			if err := writeRealUDPTables(writer, items, text); err != nil {
				return err
			}
			continue
		}
		if err := writeBenchmarkTable(writer, items, text); err != nil {
			return err
		}
	}
	return nil
}

func benchmarkSection(name string) int {
	root, _, _ := strings.Cut(strings.TrimPrefix(name, "Benchmark"), "/")
	switch {
	case strings.Contains(root, "HandshakeLifecycle"), root == "ListenerDatagramRoundTrip", root == "MemoryDatagramRoundTrip":
		return 0
	case strings.Contains(root, "RealUDP") && realUDPDirection(name) >= 0:
		return 1
	case strings.Contains(root, "ProtectedRecord"), strings.Contains(root, "Flight"), strings.Contains(root, "ACK"), root == "ListenerCIDRouting", root == "PacketSessionIngress":
		return 2
	default:
		return 3
	}
}

func writeRealUDPTables(writer io.Writer, items []*benchmark, text reportLanguage) error {
	groups, err := groupRealUDP(items)
	if err != nil {
		return err
	}
	first := true
	for index, group := range groups {
		if len(group) == 0 {
			continue
		}
		separator := ""
		if !first {
			separator = "\n"
		}
		first = false
		timing := 0
		if index == 2 || index == 3 {
			timing = 1
		}
		if index >= 4 && index%4 >= 2 {
			timing = 2
		}
		if index >= 4 && index%4 < 2 {
			timing = 3
		}
		if _, err := fmt.Fprintf(writer, "%s<a id=\"%s\"></a>\n### %s\n\n%s\n\n", separator, realUDPAnchor(index), realUDPDirectionLabel(text, index), text.realUDPTiming[timing]); err != nil {
			return err
		}
		if err := writeBenchmarkTable(writer, group, text); err != nil {
			return err
		}
	}
	return nil
}

func groupRealUDP(items []*benchmark) ([16][]*benchmark, error) {
	var groups [16][]*benchmark
	for _, item := range items {
		index := realUDPDirection(item.name)
		if index < 0 {
			return groups, fmt.Errorf("real UDP benchmark %q has no recognized client/server direction", item.name)
		}
		groups[index] = append(groups[index], item)
	}
	return groups, nil
}

func realUDPDirection(name string) int {
	if strings.HasPrefix(name, "BenchmarkDTLS13PeerRealUDP/") {
		parts := strings.Split(realUDPBenchmarkName(name), "/")
		if len(parts) != 5 {
			return -1
		}
		for peerIndex, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
			if parts[1] != peer {
				continue
			}
			direction := parts[3] + "/" + parts[4]
			for index, want := range []string{"GoClient/GoServer", "GoClient/" + peer + "Server", peer + "Client/GoServer", peer + "Client/" + peer + "Server"} {
				if direction == want {
					return 4 + peerIndex*4 + index
				}
			}
		}
		return -1
	}
	for index, direction := range realUDPDirectionPaths {
		if strings.Contains(name, "/"+direction) {
			return index
		}
	}
	return -1
}

func realUDPAnchor(index int) string {
	if index < len(realUDPAnchors) {
		return realUDPAnchors[index]
	}
	peer := []string{"nss", "boringssl", "openssl"}[(index-4)/4]
	return "real-udp-" + peer + "-" + []string{"go-go", "go-peer", "peer-go", "peer-peer"}[(index-4)%4]
}

func realUDPDirectionLabel(text reportLanguage, index int) string {
	if index < len(text.realUDP) {
		return text.realUDP[index]
	}
	peer := []string{"NSS", "BoringSSL", "OpenSSL"}[(index-4)/4]
	switch (index - 4) % 4 {
	case 1, 2:
		return peer + " / " + fmt.Sprintf(text.peerUDP[(index-4)%4], peer)
	case 3:
		return peer + " / " + fmt.Sprintf(text.peerUDP[3], peer, peer)
	default:
		return peer + " / " + text.peerUDP[0]
	}
}

func writeBenchmarkTable(writer io.Writer, items []*benchmark, text reportLanguage) error {
	hasThroughput := false
	for _, item := range items {
		hasThroughput = hasThroughput || item.byUnit["MB/s"] != nil
	}
	header := fmt.Sprintf("| %s | %s | %s | %s | %s |\n| --- | :---: | :---: | :---: | :---: |", text.columns[0], text.columns[1], text.columns[2], text.columns[4], text.columns[5])
	if hasThroughput {
		header = fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n| --- | :---: | :---: | :---: | :---: | :---: |", text.columns[0], text.columns[1], text.columns[2], text.columns[3], text.columns[4], text.columns[5])
	}
	if _, err := fmt.Fprintln(writer, header); err != nil {
		return err
	}
	for _, item := range items {
		if reason := item.unsupported[text.labelIndex]; reason != "" {
			if hasThroughput {
				if _, err := fmt.Fprintf(writer, "| %s | - | %s: %s | - | - | - |\n", markdownText(benchmarkLabel(item.name, text)), text.unsupported, markdownText(reason)); err != nil {
					return err
				}
			} else if _, err := fmt.Fprintf(writer, "| %s | - | %s: %s | - | - |\n", markdownText(benchmarkLabel(item.name, text)), text.unsupported, markdownText(reason)); err != nil {
				return err
			}
			continue
		}
		format := "| %s | %d | %s | %s | %s |\n"
		metrics := []any{markdownText(benchmarkLabel(item.name, text)), item.samples, item.timeMetric(), item.metric("B/op"), item.metric("allocs/op")}
		if hasThroughput {
			format = "| %s | %d | %s | %s | %s | %s |\n"
			metrics = []any{markdownText(benchmarkLabel(item.name, text)), item.samples, item.timeMetric(), item.metric("MB/s"), item.metric("B/op"), item.metric("allocs/op")}
		}
		if _, err := fmt.Fprintf(writer, format, metrics...); err != nil {
			return err
		}
	}
	return nil
}

func benchmarkLabel(name string, text reportLanguage) string {
	parts := benchmarkNameParts(name)
	if label, ok := benchmarkFeatureLabel(parts, text.labelIndex); ok {
		return label
	}
	switch parts[0] {
	case "WolfSSLFeatureRealUDP":
		parts = parts[1:]
	case "DTLS13PeerRealUDP":
		parts = parts[2:]
	case "HybridKeyExchangeRealUDP":
		parts[0] = "HybridKeyExchange"
	default:
		parts[0] = strings.TrimSuffix(parts[0], "HandshakeLifecycle")
		parts[0] = strings.TrimSuffix(parts[0], "Suites")
	}
	for index, part := range parts {
		parts[index] = localizeBenchmarkPart(humanizeBenchmarkPart(part), text.labelIndex)
	}
	return strings.Join(parts, " / ")
}

func localizeBenchmarkPart(value string, language int) string {
	if labels, ok := benchmarkPartTranslations[value]; ok {
		return labels[language]
	}
	return value
}

func benchmarkNameParts(name string) []string {
	name = strings.TrimPrefix(name, "Benchmark")
	if index := strings.LastIndexByte(name, '-'); index >= 0 {
		if _, err := strconv.ParseUint(name[index+1:], 10, 64); err == nil {
			name = name[:index]
		}
	}
	parts := strings.Split(name, "/")
	if len(parts) >= 2 && benchmarkPeer(parts[len(parts)-2], "Client") && benchmarkPeer(parts[len(parts)-1], "Server") {
		parts = parts[:len(parts)-2]
	}
	return parts
}

func benchmarkDisplayOrder(name string) int {
	parts := benchmarkNameParts(name)
	if len(parts) == 3 && parts[0] == "DTLS13PeerRealUDP" {
		for peerIndex, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
			if parts[1] != peer {
				continue
			}
			for index, feature := range peerRealUDPWorkloads[peer] {
				if parts[2] == feature {
					return 200 + peerIndex*200 + index*10
				}
			}
		}
	}
	if order, ok := benchmarkDisplayOrders[strings.Join(parts, "/")]; ok {
		return order
	}
	return 1000
}

func benchmarkFeatureLabel(parts []string, language int) (string, bool) {
	if len(parts) == 0 {
		return "", false
	}
	key := ""
	group := ""
	switch parts[0] {
	case "DTLS13PeerRealUDP":
		if len(parts) == 3 {
			key = parts[2]
			if key == "CertificateAES128GCM" {
				key = "PeerCertificateAES128GCM"
			}
		}
	case "WolfSSLFeatureRealUDP":
		if len(parts) == 2 {
			key = parts[1]
			if key == "GREASE" {
				key = "GREASERealUDP"
			}
		}
	case "HybridKeyExchangeRealUDP", "HybridKeyExchangeHandshakeLifecycle":
		if len(parts) == 2 {
			key, group = "HybridKeyExchange", parts[1]
		}
	case "CertificateCompressionHandshakeLifecycle":
		if len(parts) == 3 {
			key = "CertificateCompression" + parts[1] + parts[2]
		}
	case "CertificateSelectionHandshakeLifecycle":
		if len(parts) == 2 {
			key = "CertificateSelection" + parts[1]
		}
	case "MutualTLSHandshakeLifecycle":
		if len(parts) == 2 {
			key = "MutualTLS" + parts[1]
		}
	case "ECHHandshakeLifecycle":
		if len(parts) == 2 {
			key = "ECH" + parts[1]
		}
	case "ConnectionHandshakeLifecycle":
		key = "Connection"
	case "ExternalPSKHandshakeLifecycle":
		key = "ExternalPSK"
	case "GREASEHandshakeLifecycle":
		if len(parts) == 2 {
			key = "GREASE" + parts[1]
		}
	case "SessionTicketRequestHandshakeLifecycle":
		key = "SessionTicketRequest"
	}
	labels, ok := benchmarkFeatureLabels[key]
	if !ok {
		return "", false
	}
	label := labels[language]
	if group != "" {
		label += " / " + group
	}
	return label, true
}

func benchmarkPeer(value, role string) bool {
	implementation, found := strings.CutSuffix(value, role)
	if !found {
		return false
	}
	return implementation == "Go" || implementation == "WolfSSL" || implementation == "NSS" || implementation == "BORINGSSL" || implementation == "OPENSSL"
}

func humanizeBenchmarkPart(value string) string {
	if label := benchmarkPartLabels[value]; label != "" {
		return label
	}
	var label strings.Builder
	label.Grow(len(value))
	for index := range len(value) {
		if index > 0 && isUpper(value[index]) && (isLower(value[index-1]) || index+1 < len(value) && isUpper(value[index-1]) && isLower(value[index+1])) {
			label.WriteByte(' ')
		}
		label.WriteByte(value[index])
	}
	return label.String()
}

func isUpper(value byte) bool { return value >= 'A' && value <= 'Z' }
func isLower(value byte) bool { return value >= 'a' && value <= 'z' }

func (b *benchmark) timeMetric() string {
	for _, item := range workloadMetrics {
		if entry := b.byUnit[item.unit]; entry != nil {
			return formatNumber(median(entry.values)) + " " + item.label
		}
	}
	entry := b.byUnit["ns/op"]
	if entry == nil {
		return "-"
	}
	value := median(entry.values)
	switch {
	case value >= 1e9:
		return formatDecimal(value/1e9) + " s/op"
	case value >= 1e6:
		return formatDecimal(value/1e6) + " ms/op"
	case value >= 1e3:
		return formatDecimal(value/1e3) + " us/op"
	default:
		return formatNumber(value) + " ns/op"
	}
}

func (b *benchmark) metric(unit string) string {
	entry := b.byUnit[unit]
	if entry == nil {
		return "-"
	}
	return formatNumber(median(entry.values)) + " " + unit
}

func median(values []float64) float64 {
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 != 0 {
		return ordered[middle]
	}
	return (ordered[middle-1] + ordered[middle]) / 2
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatDecimal(value float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 3, 64), "0"), ".")
}

func markdownText(value string) string {
	value = html.EscapeString(value)
	value = strings.ReplaceAll(value, "&gt;", ">")
	value = strings.ReplaceAll(value, "|", "&#124;")
	value = strings.ReplaceAll(value, "[", "&#91;")
	value = strings.ReplaceAll(value, "]", "&#93;")
	return strings.ReplaceAll(value, "`", "&#96;")
}
