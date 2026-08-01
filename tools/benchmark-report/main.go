package main

import (
	"bufio"
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
	name    string
	samples int
	byUnit  map[string]*metric
}

type report struct {
	goVersion       string
	goos            string
	goarch          string
	cpu             string
	wolfSSL         string
	commit          string
	workloadChanged bool
	benchmarks      []*benchmark
	byName          map[string]*benchmark
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
	sections        [7]string
	realUDP         [4]string
	realUDPTiming   [2]string
	columns         [6]string
}

var sectionAnchors = [...]string{
	"section-connection-lifecycle",
	"section-real-udp-interoperability",
	"section-record-layer-and-reliability",
	"section-key-schedule-and-cryptography",
	"section-wire-encoding-and-parsing",
	"section-certificate-compression",
	"section-other-benchmarks",
}

var realUDPAnchors = [...]string{
	"real-udp-go-dtls-client-go-dtls-server",
	"real-udp-go-dtls-client-wolfssl-server",
	"real-udp-wolfssl-client-go-dtls-server",
	"real-udp-wolfssl-client-wolfssl-server",
}

var reportLanguages = map[string]reportLanguage{
	"en": {
		labelIndex:      0,
		title:           "Automated benchmark results",
		summary:         "%d results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.\n",
		details:         "Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.",
		quickNavigation: "Quick navigation",
		metadata:        [5]string{"Commit", "Generated", "Go", "Platform", "wolfSSL"},
		sections: [7]string{
			"Connection lifecycle",
			"Real UDP interoperability",
			"Record layer and reliability",
			"Key schedule and cryptography",
			"Wire encoding and parsing",
			"Certificate compression",
			"Other benchmarks",
		},
		realUDP: [4]string{
			"go-dtls client -> go-dtls server",
			"go-dtls client -> wolfSSL server",
			"wolfSSL client -> go-dtls server",
			"wolfSSL client -> wolfSSL server",
		},
		realUDPTiming: [2]string{
			"Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.",
			"Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.",
		},
		columns: [6]string{"Benchmark", "Samples", "Median time", "Throughput", "Harness memory", "Harness allocations"},
	},
	"zh-CN": {
		labelIndex:      1,
		title:           "自动化基准测试结果",
		summary:         "共 %d 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。\n",
		details:         "工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。",
		quickNavigation: "快速跳转",
		metadata:        [5]string{"提交", "生成时间", "Go", "平台", "wolfSSL"},
		sections: [7]string{
			"连接生命周期",
			"真实 UDP 互通",
			"记录层与可靠性",
			"密钥调度与密码学",
			"报文编码与解析",
			"证书压缩",
			"其他基准测试",
		},
		realUDP: [4]string{
			"go-dtls 客户端 -> go-dtls 服务端",
			"go-dtls 客户端 -> wolfSSL 服务端",
			"wolfSSL 客户端 -> go-dtls 服务端",
			"wolfSSL 客户端 -> wolfSSL 服务端",
		},
		realUDPTiming: [2]string{
			"中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。",
			"中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。",
		},
		columns: [6]string{"基准测试", "样本数", "中位耗时", "吞吐量", "测试框架内存", "测试框架分配次数"},
	},
	"ru": {
		labelIndex:      2,
		title:           "Результаты автоматических бенчмарков",
		summary:         "%d результатов сгруппированы по нагрузке и упорядочены по функции, затем по имени бенчмарка. Значения являются медианами выборок последнего запуска.\n",
		details:         "Специализированные метрики соединения имеют приоритет над временем стенда Go. Память и число аллокаций указаны на одну операцию Go-бенчмарка. Точный исходный вывод сохранен в артефакте Workflow.",
		quickNavigation: "Быстрая навигация",
		metadata:        [5]string{"Коммит", "Сформировано", "Go", "Платформа", "wolfSSL"},
		sections: [7]string{
			"Жизненный цикл соединения",
			"Совместимость через реальный UDP",
			"Уровень записей и надежность",
			"Расписание ключей и криптография",
			"Кодирование и разбор данных",
			"Сжатие сертификатов",
			"Другие бенчмарки",
		},
		realUDP: [4]string{
			"Клиент go-dtls -> сервер go-dtls",
			"Клиент go-dtls -> сервер wolfSSL",
			"Клиент wolfSSL -> сервер go-dtls",
			"Клиент wolfSSL -> сервер wolfSSL",
		},
		realUDPTiming: [2]string{
			"Медианное время измеряет клиент go-dtls; `ms/conn` означает одну полную операцию соединения.",
			"Медианное время измеряет клиент wolfSSL; `ms/conn` означает одно соединение, а `ms/pair` означает пару из полного и возобновленного соединений с учетом встроенного ожидания клиента.",
		},
		columns: [6]string{"Бенчмарк", "Замеры", "Медианное время", "Пропускная способность", "Память стенда", "Аллокации стенда"},
	},
}

var workloadMetrics = [...]struct {
	unit  string
	label string
}{
	{"go_ms/conn", "ms/conn"},
	{"wolfssl_ms/conn", "ms/conn"},
	{"wolf_process_ms/conn", "ms/conn"},
	{"wolf_process_ms/pair", "ms/pair"},
}

var benchmarkPartLabels = map[string]string{
	"1301":                                   "AES-128-GCM",
	"1302":                                   "AES-256-GCM",
	"AES128CCM":                              "AES-128-CCM",
	"AES128GCM":                              "AES-128-GCM",
	"AES256GCM":                              "AES-256-GCM",
	"ApplicationDataRoundTrip":               "Application data",
	"BuildPlainACKRecords":                   "Plain ACK build",
	"BuildProtectedACKRecords":               "Protected ACK build",
	"CertificateAES128GCM":                   "Certificate / AES-128-GCM",
	"CertificateCompression":                 "Certificate compression",
	"ChaCha20Poly1305":                       "ChaCha20-Poly1305",
	"Client1":                                "1 key share",
	"Client4":                                "4 key shares",
	"Client9":                                "9 key shares",
	"ConnectionID":                           "CID",
	"EarlyData":                              "0-RTT",
	"HandshakeInboxSequentialSingleFragment": "Inbox / Single fragment",
	"HandshakeInboxSequentialSingleFragmentBatch": "Inbox / Fragment batch",
	"HandshakeInboxSequentialSingleFragmentReuse": "Inbox / Fragment reuse",
	"HybridKeyExchange":                           "Hybrid KEX",
	"KeyScheduleSideDerivations":                  "Key derivation",
	"MarshalHandshakeMessages":                    "Handshake marshal",
	"MutualTLS":                                   "mTLS",
	"MutualTLSSessionResumption":                  "mTLS resumption",
	"ParseKeyShares":                              "Key share parse",
	"PostHandshakeAuthentication":                 "Post-handshake authentication",
	"ProtectedRecordReceiveErrorUnauthenticated":  "Reject unauthenticated record",
	"ProtectedRecordRoundTrip":                    "Record round trip",
	"ProtectedRecordSeal":                         "Record seal",
	"ReceivingTrafficKeyUpdate":                   "Receive KeyUpdate",
	"SecP256r1MLKEM768":                           "SecP256r1MLKEM768",
	"SecP384r1MLKEM1024":                          "SecP384r1MLKEM1024",
	"SendingTrafficKeyUpdate":                     "Send KeyUpdate",
	"ServerCertificate":                           "Server cert",
	"SessionResumption":                           "Session resumption",
}

var benchmarkFeatureLabels = map[string][3]string{
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
	"Plain ACK build":                      {"Plain ACK build", "明文 ACK 构建", "Построение открытого ACK"},
	"Protected ACK build":                  {"Protected ACK build", "受保护 ACK 构建", "Построение защищенного ACK"},
	"Build Plain Flight":                   {"Build Plain Flight", "构建明文握手报文组", "Построение открытой группы сообщений"},
	"Build Protected Flight":               {"Build Protected Flight", "构建受保护握手报文组", "Построение защищенной группы сообщений"},
	"Combine Flights":                      {"Combine Flights", "合并握手报文组", "Объединение групп сообщений"},
	"Flight First Refresh":                 {"Flight First Refresh", "握手报文组首次刷新", "Первое обновление группы сообщений"},
	"Flight Initial History Batch":         {"Flight Initial History Batch", "握手报文组初始历史批次", "Начальная история группы сообщений"},
	"Flight Pending Indices":               {"Flight Pending Indices", "握手报文组待处理索引", "Индексы ожидания группы сообщений"},
	"Flight Wire Window":                   {"Flight Wire Window", "握手报文组传输窗口", "Окно передачи группы сообщений"},
	"Inbox / Single fragment":              {"Inbox / Single fragment", "接收缓存 / 单分片", "Буфер приема / один фрагмент"},
	"Inbox / Fragment batch":               {"Inbox / Fragment batch", "接收缓存 / 分片批次", "Буфер приема / пакет фрагментов"},
	"Inbox / Fragment reuse":               {"Inbox / Fragment reuse", "接收缓存 / 分片复用", "Буфер приема / повторное использование фрагмента"},
	"Handshake Reassembly":                 {"Handshake Reassembly", "握手重组", "Сборка рукопожатия"},
	"Handshake Reassembly Single Fragment": {"Handshake Reassembly Single Fragment", "握手重组单分片", "Сборка рукопожатия одним фрагментом"},
	"Parse ACK":                            {"Parse ACK", "解析 ACK", "Разбор ACK"},
	"Protected Record CID":                 {"Protected Record CID", "受保护记录 CID", "Защищенная запись CID"},
	"Reject unauthenticated record":        {"Reject unauthenticated record", "拒绝未认证记录", "Отклонение неаутентифицированной записи"},
	"Record round trip":                    {"Record round trip", "记录往返", "Обмен записью туда и обратно"},
	"Protected Record Round Trip In Place": {"Protected Record Round Trip In Place", "受保护记录原地往返", "Обмен защищенной записью на месте"},
	"Record seal":                          {"Record seal", "记录封装", "Защита записи"},
	"Calculate PSK Binder":                 {"Calculate PSK Binder", "计算 PSK 绑定值", "Расчет привязки PSK"},
	"Derive Traffic Keys":                  {"Derive Traffic Keys", "派生流量密钥", "Вывод ключей трафика"},
	"Derive Traffic Keys Into":             {"Derive Traffic Keys Into", "派生流量密钥并写入", "Вывод ключей трафика в буфер"},
	"Empty Transcript Hash":                {"Empty Transcript Hash", "空握手转录哈希", "Хеш пустого транскрипта"},
	"Finished Verify Data":                 {"Finished Verify Data", "Finished 验证数据", "Данные проверки Finished"},
	"Install Application Keys":             {"Install Application Keys", "安装应用密钥", "Установка ключей приложения"},
	"Key Schedule Derivation":              {"Key Schedule Derivation", "密钥调度派生", "Вывод расписания ключей"},
	"Key derivation":                       {"Key derivation", "密钥派生", "Вывод ключей"},
	"New Record Cipher":                    {"New Record Cipher", "新建记录密码器", "Новый шифр записи"},
	"Receive KeyUpdate":                    {"Receive KeyUpdate", "接收 KeyUpdate", "Получение KeyUpdate"},
	"Send KeyUpdate":                       {"Send KeyUpdate", "发送 KeyUpdate", "Отправка KeyUpdate"},
	"Transcript Clone":                     {"Transcript Clone", "握手转录克隆", "Клонирование транскрипта"},
	"Transcript Sum":                       {"Transcript Sum", "握手转录求和", "Сумма транскрипта"},
	"Marshal Extensions":                   {"Marshal Extensions", "编码扩展", "Кодирование расширений"},
	"Handshake marshal":                    {"Handshake marshal", "编码握手", "Кодирование рукопожатия"},
	"Parse Extensions":                     {"Parse Extensions", "解析扩展", "Разбор расширений"},
	"Parse Handshake Fragment":             {"Parse Handshake Fragment", "解析握手分片", "Разбор фрагмента рукопожатия"},
	"Key share parse":                      {"Key share parse", "解析密钥份额", "Разбор доли ключа"},
	"Parse Plain Record":                   {"Parse Plain Record", "解析明文记录", "Разбор открытой записи"},
	"Compress":                             {"Compress", "压缩", "Сжатие"},
	"Decompress":                           {"Decompress", "解压", "Распаковка"},
	"Empty":                                {"Empty", "空", "Пустой"},
	"Single":                               {"Single", "单条", "Один"},
	"Sorted64":                             {"Sorted64", "已排序 64", "Отсортированный 64"},
	"Reversed64":                           {"Reversed64", "逆序 64", "Обратный 64"},
	"Single Reuse":                         {"Single Reuse", "单条复用", "Повторное использование одного"},
	"Reuse Single":                         {"Reuse Single", "单条复用", "Повторное использование одного"},
	"Allocated":                            {"Allocated", "已分配", "Выделенный"},
	"Reuse Window":                         {"Reuse Window", "复用窗口", "Окно повторного использования"},
	"Pending":                              {"Pending", "待处理", "Ожидание"},
	"Retransmit":                           {"Retransmit", "重传", "Повторная передача"},
	"Round Trip":                           {"Round Trip", "往返", "Обмен туда и обратно"},
	"Seal":                                 {"Seal", "封装", "Защита"},
	"Owned":                                {"Owned", "独占", "Владеемое"},
	"Reuse":                                {"Reuse", "复用", "Повторное использование"},
	"View":                                 {"View", "视图", "Представление"},
	"View Into":                            {"View Into", "写入视图", "Представление в буфере"},
	"Ordered View":                         {"Ordered View", "有序视图", "Упорядоченное представление"},
	"Certificate":                          {"Certificate", "证书", "Сертификат"},
	"Certificate Verify":                   {"Certificate Verify", "证书验证", "Проверка сертификата"},
	"Client Hello":                         {"Client Hello", "客户端 Hello", "Client Hello"},
	"Hello Retry Request":                  {"Hello Retry Request", "Hello 重试请求", "Hello Retry Request"},
	"New Connection ID":                    {"New Connection ID", "新连接 ID", "Новый ID соединения"},
	"New Session Ticket":                   {"New Session Ticket", "新会话票据", "Новый билет сеанса"},
	"Resumption Client Hello":              {"Resumption Client Hello", "恢复 Client Hello", "Client Hello возобновления"},
	"Server Hello":                         {"Server Hello", "服务端 Hello", "Server Hello"},
	"Session Ticket State":                 {"Session Ticket State", "会话票据状态", "Состояние билета сеанса"},
	"Early Traffic":                        {"Early Traffic", "早期流量", "Ранний трафик"},
	"Exporter":                             {"Exporter", "导出器", "Экспортер"},
	"Exporter Zero":                        {"Exporter Zero", "零值导出器", "Нулевой экспортер"},
	"Resumption PSK":                       {"Resumption PSK", "恢复 PSK", "PSK возобновления"},
	"Traffic Update":                       {"Traffic Update", "流量更新", "Обновление трафика"},
	"1 key share":                          {"1 key share", "1 个密钥份额", "1 доля ключа"},
	"4 key shares":                         {"4 key shares", "4 个密钥份额", "4 доли ключа"},
	"9 key shares":                         {"9 key shares", "9 个密钥份额", "9 долей ключа"},
}

var benchmarkDisplayOrders = map[string]int{
	"ConnectionHandshakeLifecycle":                                     10,
	"MutualTLSHandshakeLifecycle/Full":                                 20,
	"MutualTLSHandshakeLifecycle/Resumed":                              30,
	"SessionTicketRequestHandshakeLifecycle":                           35,
	"ExternalPSKHandshakeLifecycle":                                    40,
	"CertificateCompressionHandshakeLifecycle/ServerCertificate/Plain": 50,
	"CertificateCompressionHandshakeLifecycle/ServerCertificate/Zlib":  60,
	"CertificateCompressionHandshakeLifecycle/MutualTLS/Plain":         70,
	"CertificateCompressionHandshakeLifecycle/MutualTLS/Zlib":          80,
	"ECHHandshakeLifecycle/Direct":                                     90,
	"ECHHandshakeLifecycle/HRR":                                        100,
	"HybridKeyExchangeHandshakeLifecycle/X25519MLKEM768":               110,
	"HybridKeyExchangeHandshakeLifecycle/SecP256r1MLKEM768":            120,
	"HybridKeyExchangeHandshakeLifecycle/SecP384r1MLKEM1024":           130,

	"WolfSSLFeatureRealUDP/CertificateAES128GCM":        10,
	"WolfSSLFeatureRealUDP/ApplicationDataRoundTrip":    20,
	"WolfSSLFeatureRealUDP/MutualTLS":                   30,
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
		if *baseline == "" || *candidate == "" || *input != "" {
			flag.Usage()
			os.Exit(2)
		}
		baselineReport, err := readReportFile(*baseline)
		if err != nil {
			fatal(err)
		}
		candidateReport, err := readReportFile(*candidate)
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

	result, err := readReportFile(*input)
	if err != nil {
		fatal(err)
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

func readReportFile(path string) (*report, error) {
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
		case strings.HasPrefix(line, "commit: ") && result.commit == "":
			result.commit = strings.TrimSpace(strings.TrimPrefix(line, "commit: "))
		case strings.HasPrefix(line, "workload-changed: "):
			changed, err := strconv.ParseBool(strings.TrimSpace(strings.TrimPrefix(line, "workload-changed: ")))
			if err != nil {
				return nil, fmt.Errorf("line %d: parse workload change marker: %w", lineNumber, err)
			}
			result.workloadChanged = changed
		case strings.HasPrefix(line, "go version ") && result.goVersion == "":
			result.goVersion = line
		case strings.HasPrefix(line, "goos: ") && result.goos == "":
			result.goos = strings.TrimSpace(strings.TrimPrefix(line, "goos: "))
		case strings.HasPrefix(line, "goarch: ") && result.goarch == "":
			result.goarch = strings.TrimSpace(strings.TrimPrefix(line, "goarch: "))
		case strings.HasPrefix(line, "cpu: ") && result.cpu == "":
			result.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu: "))
		case strings.HasPrefix(line, "wolfssl: ") && result.wolfSSL == "":
			result.wolfSSL = strings.TrimSpace(strings.TrimPrefix(line, "wolfssl: "))
		case strings.HasPrefix(line, "Benchmark"):
			if err := result.addBenchmarkLine(lineNumber, line); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read benchmark output: %w", err)
	}
	if len(result.benchmarks) == 0 {
		return nil, errors.New("benchmark output contains no results")
	}
	return result, nil
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
		if err := validateBenchmarkCoverage(fields[0]); err != nil {
			return fmt.Errorf("line %d: %w", lineNumber, err)
		}
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
	case "CertificateCompression":
		parts = parts[1:]
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
	name := strings.TrimSuffix(filepath.Base(output), filepath.Ext(output))
	name = strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, ".zh-CN"), ".en"), ".ru")
	return fmt.Sprintf("[简体中文](%s.md) | [English](%s.en.md) | [Русский](%s.ru.md)", name, name, name)
}

func writeReport(writer io.Writer, result *report, commit, wolfSSL, generated string, text reportLanguage, languageLinks string) error {
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
				if _, err := fmt.Fprintf(writer, "\n  - [%s (%d)](#%s)", text.realUDP[direction], len(group), realUDPAnchors[direction]); err != nil {
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
	case strings.Contains(root, "HandshakeLifecycle"):
		return 0
	case strings.Contains(root, "RealUDP"):
		return 1
	case strings.Contains(root, "ProtectedRecord"), strings.Contains(root, "Flight"), strings.Contains(root, "ACK"),
		strings.Contains(root, "HandshakeReassembly"), strings.Contains(root, "HandshakeInbox"):
		return 2
	case strings.Contains(root, "Transcript"), strings.Contains(root, "TrafficKey"), strings.Contains(root, "KeySchedule"),
		strings.Contains(root, "Finished"), strings.Contains(root, "PSKBinder"), strings.Contains(root, "RecordCipher"),
		strings.Contains(root, "ApplicationKeys"):
		return 3
	case strings.HasPrefix(root, "Parse"), strings.HasPrefix(root, "Marshal"):
		return 4
	case strings.Contains(root, "CertificateCompression"):
		return 5
	default:
		return 6
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
		if _, err := fmt.Fprintf(writer, "%s<a id=\"%s\"></a>\n### %s\n\n%s\n\n", separator, realUDPAnchors[index], text.realUDP[index], text.realUDPTiming[index/2]); err != nil {
			return err
		}
		if err := writeBenchmarkTable(writer, group, text); err != nil {
			return err
		}
	}
	return nil
}

func groupRealUDP(items []*benchmark) ([4][]*benchmark, error) {
	var groups [4][]*benchmark
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
	switch {
	case strings.Contains(name, "/GoClient/GoServer"):
		return 0
	case strings.Contains(name, "/GoClient/WolfSSLServer"):
		return 1
	case strings.Contains(name, "/WolfSSLClient/GoServer"):
		return 2
	case strings.Contains(name, "/WolfSSLClient/WolfSSLServer"):
		return 3
	default:
		return -1
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
	case "HybridKeyExchangeRealUDP":
		parts[0] = "HybridKeyExchange"
	case "CertificateCompression":
		parts = parts[1:]
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
	if order, ok := benchmarkDisplayOrders[strings.Join(benchmarkNameParts(name), "/")]; ok {
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
	case "WolfSSLFeatureRealUDP":
		if len(parts) == 2 {
			key = parts[1]
		}
	case "HybridKeyExchangeRealUDP", "HybridKeyExchangeHandshakeLifecycle":
		if len(parts) == 2 {
			key, group = "HybridKeyExchange", parts[1]
		}
	case "CertificateCompressionHandshakeLifecycle":
		if len(parts) == 3 {
			key = "CertificateCompression" + parts[1] + parts[2]
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
	return implementation == "Go" || implementation == "WolfSSL"
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
	return strings.ReplaceAll(value, "`", "&#96;")
}
