# go-dtls

[![Go Reference](https://pkg.go.dev/badge/github.com/puernya/go-dtls.svg)](https://pkg.go.dev/github.com/puernya/go-dtls)
[![CI](https://github.com/puernya/go-dtls/actions/workflows/ci.yml/badge.svg)](https://github.com/puernya/go-dtls/actions/workflows/ci.yml)
[![License: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](LICENSE)

`go-dtls` 是一个使用 Go 实现的 DTLS 1.3 库，协议行为以 [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147) 为准。模块导入路径为 `github.com/puernya/go-dtls`，包名为 `dtls13`。

当前实现完成了 RFC 9147 的强制语义，并覆盖 RFC 9147 列出的全部 11 项直接规范性引用。TLS 1.3 基线按取代 RFC 8446 的 [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) 审计；其中适用于现有功能的主体行为已经实现，但仍有下文明确列出的推荐行为和后续规范缺口。本文记录 2026-07-29 当前源码的实现状态、测试证据和性能基线。

> 本项目仍处于活跃开发阶段。RFC 审计完成不等同于第三方合规认证。RFC 9325 的证书安全下限默认强制执行：RSA 服务端认证 leaf 不得小于 2048 位，SHA-1/MD5 证书不能通过自签名、信任锚、`InsecureSkipVerify` 或 session resumption 绕过。

## 核心语义

DTLS 是不可靠报文协议，不是 TLS 字节流：

- `Conn` 有意不实现 `net.Conn` 或 `net.PacketConn`。
- 一次 `WriteDatagram` 发送一条 DTLS Application Data record；应用报文不会被内部拆分、排序或重传。
- 一次 `ReadDatagram` 消费一条认证后的 Application Data record；缓冲区过小时会丢弃余部，并通过 `DatagramInfo.Truncated` 报告。
- 默认情况下，超过当前路径 MTU 或 RFC record 上限的应用报文返回 `ErrDatagramTooLarge`，不会部分发送；`IgnorePathMTU` 可只跳过前一项检查。
- 握手消息仍按 RFC 9147 执行分片、ACK、丢包恢复和指数退避；这些可靠性机制不会改变应用数据的报文语义。
- `Listener` 从 UDP socket 接受经过认证的 DTLS association，并返回强类型 `*dtls13.Conn`。

## 环境要求

| 项目 | 要求 |
| --- | --- |
| Go | 最低 Go 1.26；当前审计环境为 Go 1.26.3 `windows/amd64` |
| Transport | `udp`、`udp4` 或 `udp6`；不接受 TCP |
| Windows race | 仓库脚本需要 Zig 0.17 和可用的 CGO 工具链 |
| wolfSSL 互通测试 | 可选；需要设置 `WOLFSSL_ROOT` 指向兼容的 wolfSSL 源码/构建目录 |

安装：

```sh
go get github.com/puernya/go-dtls
```

导入路径末段包含连字符，因此 Go 源码中使用声明的包名 `dtls13`：

```go
import dtls13 "github.com/puernya/go-dtls"
```

## 快速开始

### 客户端

`Dial` 创建 UDP 连接并在返回前完成握手。下面的客户端验证服务端证书、发送一条报文并读取一条响应：

```go
package main

import (
	"crypto/x509"
	"log"
	"os"
	"time"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	caPEM, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		log.Fatal("ca.pem does not contain a certificate")
	}

	conn, err := dtls13.Dial("udp", "127.0.0.1:4444", &dtls13.Config{
		RootCAs:    roots,
		ServerName: "server.example",
		NextProtos: []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.WriteDatagram([]byte("ping")); err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Fatal(err)
	}
	if info.Truncated {
		log.Fatalf("response needs %d bytes", info.FullLength)
	}
	log.Printf("received %q from %s", buffer[:n], info.Source)
}
```

### 服务端

`Listen` 创建 UDP listener。`Accept` 在收到新 association 后返回 `*Conn`，应用可显式调用 `Handshake`，也可让首次 `ReadDatagram` 或 `WriteDatagram` 触发握手。

```go
package main

import (
	"crypto/tls"
	"log"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := dtls13.Listen("udp", ":4444", &dtls13.Config{
		Certificates: []tls.Certificate{certificate},
		NextProtos:   []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(conn)
	}
}

func serve(conn *dtls13.Conn) {
	defer conn.Close()
	if err := conn.Handshake(); err != nil {
		log.Printf("handshake: %v", err)
		return
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Printf("read: %v", err)
		return
	}
	if info.Truncated {
		log.Printf("discarded truncated datagram: need %d bytes", info.FullLength)
		return
	}
	if _, err = conn.WriteDatagram(buffer[:n]); err != nil {
		log.Printf("write: %v", err)
	}
}
```

## 常用 API

| API | 用途 |
| --- | --- |
| `Dial` / `DialWithDialer` | 创建 UDP 客户端连接并完成握手；`DialWithDialer` 可配置本地地址和拨号超时 |
| `Listen` | 创建并拥有一个 UDP listener |
| `NewListener` | 在已有 `net.PacketConn` 上创建 DTLS listener；配置错误由 `Accept` 返回 |
| `Listener.Accept` | 接受一个新 association，返回 `*Conn` |
| `Client` / `Server` | 在已有 connected `net.Conn` 上包装 DTLS 客户端或服务端；首次 I/O 前可显式握手 |
| `Handshake` / `HandshakeContext` | 显式执行一次握手；后续调用返回同一结果 |
| `ReadDatagram` | 读取一条认证报文，并返回来源、完整长度和截断状态 |
| `WriteDatagram` | 向 association 的认证对端发送一条报文 |
| `SetDeadline` / `SetReadDeadline` / `SetWriteDeadline` | 设置底层报文 I/O deadline |
| `ConnectionState` | 获取版本、套件、ALPN、证书、恢复状态和当前 CID |
| `Close` | 发送 `close_notify`、清除 traffic/resumption secrets 并关闭底层连接 |

### 报文大小与截断

应用应在写入前使用 `PathMTU` 和 `RecordOverhead` 估算当前最大 payload，并处理 PMTU 在连接生命周期内下降的情况：

```go
maximum := conn.PathMTU() - conn.RecordOverhead()
if len(payload) > maximum {
	// 拆分属于应用协议职责；每个分片仍是一条独立 datagram。
}
```

即使预检查通过，路径状态仍可能在检查后变化，因此必须处理写入错误：

```go
if _, err := conn.WriteDatagram(payload); errors.Is(err, dtls13.ErrDatagramTooLarge) {
	// 根据应用协议缩小或拆分 payload，然后作为新 datagram 重试。
}
```

需要由应用主动探测 PMTU 或明确依赖 IP 分片时，可以配置
`IgnorePathMTU: true`。此时 `WriteDatagram` 和 `WriteEarlyData` 跳过库内
PMTU payload 检查，把单个完整 DTLS record 直接交给底层 transport；握手、
ACK 和 post-handshake flight 仍按 PMTU 分片、重传和退避。该选项不会放宽
`2^14` 字节 record content 上限，底层仍可能分片、丢包或返回
`ErrDatagramTooLarge`，库不会降低 PMTU 或自动重发该应用 datagram。

`ReadDatagram` 的短缓冲不是流式读取。`Truncated=true` 时，该 record 的未读部分已经丢弃，下一次读取会返回下一条 record。

### 错误模型

| 错误 | 含义 |
| --- | --- |
| `*ConfigError` | 本地配置不合法，例如 transport、MTU、套件或资源上限无效 |
| `*ProtocolError` | 收到或生成了不符合协议状态/格式的内容 |
| `AlertError` | 对端返回 fatal TLS alert；可用 `errors.As` 获取 alert description 数值 |
| `ErrDatagramTooLarge` | 应用 datagram 超过当前 PMTU 或 record 上限；开启 `IgnorePathMTU` 后底层仍可能返回；可用 `errors.Is` 判断 |
| `ErrEarlyDataUnavailable` | 没有可用的 early-data ticket，或该连接不能发送 0-RTT |
| `ErrEarlyDataRejected` | 握手完成，但对端因 HRR、重放或策略拒绝了已发送的 0-RTT |
| `io.EOF` | 对端发送了合法 `close_notify`，读取方向已关闭 |

deadline、socket 关闭和底层 UDP 错误沿 Go `net` 错误模型返回；调用方不应依赖错误字符串。

## 高级能力

| 能力 | API / 配置 | 说明 |
| --- | --- | --- |
| Session resumption | `ClientSessionCache`、`NewLRUClientSessionCache` | 客户端缓存 NewSessionTicket；服务端由 `SessionTicketKey` 和 ticket 配置控制；支持保留客户端认证状态的 mTLS 恢复 |
| 0-RTT | `WriteEarlyData`、`MaxEarlyData`、`EarlyDataReplayCache` | 仅恢复连接可用；调用方必须处理 `ErrEarlyDataUnavailable` 和 `ErrEarlyDataRejected`，且 early data 必须具备可重放语义 |
| KeyUpdate | `SendKeyUpdate(requestPeer)` | 可靠发送并在 ACK 后切换发送 epoch；接近 AEAD 使用上限时也会自动触发 |
| CID | `ConnectionID`、`GetConnectionID`、`SendNewConnectionIDs`、`RequestConnectionIDs`、`UseNextConnectionID` | 支持 RFC 9146 协商、动态更新和 Listener 路由；CID 认证不等于新网络路径已验证 |
| 握手内客户端认证 | `ClientAuth`、`ClientCAs`、`Certificates` | 使用 `crypto/tls` 的客户端证书策略 |
| 握手后客户端认证 | `PostHandshakeAuth`、`RequestClientCertificate` | 客户端先声明支持，服务端再发起 PHA |
| Exporter | `ConnectionState().ExportKeyingMaterial` | 提供 RFC 8446 section 7.5、使用 DTLS `dtls13` label 的导出材料 |

## 主要配置

`Config` 在 TLS 1.3 语义相同的字段上尽量沿用 `crypto/tls.Config`。配置可以复用，但首次使用后不得修改；需要派生配置时使用 `Clone`。

| 配置 | 默认值 / 行为 |
| --- | --- |
| `Certificates` / `GetCertificate` | 服务端证书；RSA leaf 必须至少 2048 位，整条证书链不得使用 SHA-1/MD5 |
| `RootCAs` / `ServerName` | 客户端服务端证书验证；`Dial` 未设置 `ServerName` 时使用目标主机名 |
| `ClientCAs` / `ClientAuth` | 服务端客户端证书验证策略 |
| `VerifyPeerCertificate` | 在完整握手的标准证书处理后执行附加验证；与 `crypto/tls` 一样，恢复连接不会再次调用 |
| `InsecureSkipVerify` | 默认 `false`；生产环境不应依赖它跳过身份验证 |
| `NextProtos` | ALPN 协议列表 |
| `CipherSuites` | AES-128-GCM、AES-256-GCM、ChaCha20-Poly1305、AES-128-CCM |
| `CurvePreferences` | X25519、P-256 |
| `MTU` | 1200 字节 UDP payload；最小 256 |
| `IgnorePathMTU` | 默认 `false`；仅让 Application Data 跳过库内 PMTU 检查，握手不受影响 |
| `FlightInterval` | 1 秒初始握手重传间隔 |
| `MaxFlightInterval` | 60 秒指数退避上限 |
| `HandshakeTimeout` | 30 秒 |
| `ReplayWindow` | 每个 epoch 64 条 record |
| `MaxHandshakeMessage` | 1 MiB，最大可配置到 `2^24-1` |
| `MaxBufferedApplicationData` | 1 MiB |
| `MaxBufferedApplicationDatagrams` | 1024 条 |
| `MaxPendingConnections` | 128 个 Listener session |
| `MaxSessionQueueDatagrams` | 每个 Listener session 64 条 |
| `SessionTicketLifetime` | 24 小时，最大 7 天 |
| `MaxEarlyData` | 0，即默认关闭 0-RTT |
| `MaxConnectionIDs` | 每个方向 8 个 CID |

### mTLS 快速恢复

客户端和服务端沿用普通 mTLS 配置，只需同时启用客户端 session cache 和服务端 ticket：

```go
clientConfig := &dtls13.Config{
	Certificates:       []tls.Certificate{clientCertificate},
	RootCAs:            serverRoots,
	ServerName:         "server.example",
	ClientSessionCache: dtls13.NewLRUClientSessionCache(64),
}
serverConfig := &dtls13.Config{
	Certificates:          []tls.Certificate{serverCertificate},
	ClientAuth:            tls.RequireAndVerifyClientCert,
	ClientCAs:             clientRoots,
	SessionTicketKey:      ticketKey,
	SessionTicketLifetime: time.Hour,
}
```

首次连接执行完整 `CertificateRequest -> Certificate -> CertificateVerify` 客户端认证。后续连接用 RFC 9147/RFC 9846 的 PSK 握手恢复，不重新发送证书；服务端从经过 AES-256-GCM 认证加密的 ticket 恢复客户端证书和已验证链，并按当前 `ClientAuth`、`ClientCAs`、证书有效期重新决定是否接受。策略不满足时忽略 ticket 并回退完整握手。

旧格式 ticket 没有保存客户端身份，无法证明原会话是否匿名；因此只在 `ClientAuth == tls.NoClientCert` 时恢复。配置任何客户端证书策略都会回退完整握手，避免升级后静默丢失已有身份。

续签 ticket 保留最初在线 `CertificateVerify` 的时间，`SessionTicketLifetime` 同时限制 ticket 寿命和这次客户端认证的总寿命。`VerifyPeerCertificate` 不在恢复时重跑；应用身份策略变化时应更换 `SessionTicketKey` 或禁用 session ticket。应用必须定期轮换显式 ticket key；本配置只有一个活动 key，更换后旧 ticket 立即失效。握手后的 PHA 不自动签发补充 ticket，需要 PHA 身份进入后续恢复时应重新建立完整 mTLS 连接。

支持的 cipher suite：

| 常量 | ID | 状态 |
| --- | --- | --- |
| `TLS_AES_128_GCM_SHA256` | `0x1301` | 支持 |
| `TLS_AES_256_GCM_SHA384` | `0x1302` | 支持 |
| `TLS_CHACHA20_POLY1305_SHA256` | `0x1303` | 支持 |
| `TLS_AES_128_CCM_SHA256` | `0x1304` | 支持 |
| `TLS_AES_128_CCM_8_SHA256` | `0x1305` | 显式禁用；通用库无法保证 RFC 9147 要求的部署级额外防伪措施 |

## RFC 9147 完成度

规范关键字按 BCP 14 解释。项目验收要求 `MUST`、`MUST NOT`、`REQUIRED`、`SHALL` 和 `SHALL NOT` 在客户端与服务端的发送、接收方向都严格执行；`SHOULD` 类要求由本客户端主动实现，本服务端只宽容不削弱认证、机密性、反重放、放大限制和状态一致性的对端偏差。当前未达到该目标的推荐项明确列为缺口。`MAY` 和 `OPTIONAL` 能力未启用时不构成缺陷，但一旦协商，其条件性强制要求必须完整执行。

### 总体状态

| 规范 | 状态 | 当前结论 |
| --- | --- | --- |
| [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147) | 完成 | Record、Handshake、epoch、ACK、KeyUpdate、CID update、Application Data 与适用安全要求完成；推荐行为在已启用范围完成 |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | 完成 | CID 协商、方向性 CID、更新、Listener 路由、错误处理和地址保持完成；DTLS 1.2 专属细节不适用 |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | 基本完成 | 当前仍把 `user_canceled(90)` 返回为 `AlertError`，未按 SHOULD 忽略并等待 `close_notify`；`general_error(117)` 尚无命名常量和主动发送策略 |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | 部分实现 | PFS、AEAD、SNI/ALPN、ticket、0-RTT、KeyUpdate 和证书安全下限已覆盖；缺少 OCSP stapling，且本模块有意不实现该 BCP 要求的 DTLS 1.2 |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | 部分实现 | Go X.509 与 `ServerName` 覆盖 DNS-ID/IP-ID；URI-ID、SRV-ID 和应用 service identity 由调用方验证回调承担 |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | 未实现 | 当前不会未经路径验证自动 rebind，因此没有安全降级；但提供 CID 的客户端尚未按 SHOULD 同时提供 Return Routability Check |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | 未实现 | 尚不协商 `record_size_limit`；当前仅有 RFC record 上限和动态 PMTU，两者不能替代该扩展 |

### 章节覆盖

| RFC 9147 章节 | 状态 | 实现摘要 |
| --- | --- | --- |
| section 1-2 引言和术语 | 不适用 | 按 RFC 2119/8174 解释规范关键字，无独立 wire 功能 |
| section 3 设计目标 | 完成 | 丢包、乱序、重复、延迟、分片、反重放、动态 PMTU 和资源回收 |
| section 4 Record Layer | 完成 | DTLSPlaintext、统一头、截断序号、epoch、AEAD、序号保护、anti-replay、CID demux 和使用上限 |
| section 5 Handshake | 完成 | TLS 1.3 握手、HRR cookie、认证、分片重组、flight、ACK、重传、超时和同五元组新关联 |
| section 6 Epochs | 完成 | epoch 0/1/2/3/4+、0-RTT、KeyUpdate 及旧 key 有界保留和清除 |
| section 7 ACK | 完成 | content type 26、空 ACK、即时 ACK、部分 ACK、滑动窗口和 post-handshake 可靠 ACK |
| section 8 KeyUpdate | 完成 | ACK-gated 更新、重传、`update_requested`、旧 epoch 保留和极限处理 |
| section 9 CID Update | 完成 | New/RequestConnectionId、动态路由、更新 ACK、资源上限和 prefix-free 校验 |
| section 10 Application Data | 完成 | connected-datagram API、报文边界、无排序/重传、显式截断和 deadline |
| section 11 Security | 完成 | cookie 轮换、3 倍放大限制、anti-replay、AEAD limit、地址保持和有界状态 |
| section 12 DTLS 1.2 differences | 完成 | 所有适用差异已进入记录层、握手、epoch、ACK 和 CID 实现 |
| section 13 DTLS 1.2 updates | 不适用 | 本模块只实现 DTLS 1.3 |
| section 14 IANA | 不适用 | 使用已分配编号，本库不执行注册表操作 |

### 直接规范性引用

该表对应 RFC Editor 发布的 RFC 9147 XML 中全部 11 项 `Normative References`。下层协议由 Go 和操作系统实现，不表示本库重新实现 UDP/TCP/IP 栈。

| 规范 | RFC 9147 中的用途 | 覆盖状态 |
| --- | --- | --- |
| [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439) | ChaCha20 序号保护的 counter、nonce 和 block function | 完成 |
| [RFC 768](https://www.rfc-editor.org/rfc/rfc768) | UDP transport | 下层完成；库限制为 UDP，保持 datagram 边界 |
| [RFC 793](https://www.rfc-editor.org/rfc/rfc793) | 可靠传输背景和旧 epoch 的 MSL 参考 | 适用语义完成；本库不接受 TCP transport |
| [RFC 1191](https://www.rfc-editor.org/rfc/rfc1191) | IPv4 PMTU discovery | 下层交互完成；内部处理平台 MTU 错误，对外统一为 `ErrDatagramTooLarge`；应用可用 `IgnorePathMTU` 主动探测 |
| [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) | BCP 14 规范关键字 | 审计规则完成 |
| [RFC 4443](https://www.rfc-editor.org/rfc/rfc4443) | IPv6 ICMP Packet Too Big | 下层交互完成；由 OS 处理 ICMPv6，库消费写入错误反馈 |
| [RFC 4821](https://www.rfc-editor.org/rfc/rfc4821) | Packetization Layer PMTU Discovery | 完成；发送错误和连续黑洞超时触发退避与重新分片 |
| [RFC 6298](https://www.rfc-editor.org/rfc/rfc6298) | RTO 初值、指数退避和上限 | 完成；默认 1 秒、最大 60 秒，并使用符合 Karn 原则的 RTT 样本 |
| [RFC 8174](https://www.rfc-editor.org/rfc/rfc8174) | BCP 14 大小写限定 | 审计规则完成 |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | DTLS Connection ID | 完成；协商、更新、路由、隔离和资源上限 |
| [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446) | TLS 1.3 基础协议 | RFC 9147 继承的握手、认证、key schedule、PSK/0-RTT、KeyUpdate、PHA 和 exporter 完成；现行修订按 RFC 9846 的状态单独列出 |

### 后续规范

| 规范 | 与本实现的关系 | 状态 |
| --- | --- | --- |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | 取代 RFC 8446 的 TLS 1.3 规范；修订 KeyShare、PSK/HRR、NST、AEAD limit、KeyUpdate、alert 和 vector 边界 | 主体变更完成；mTLS 恢复保留认证状态、当前策略/CA/有效期及总认证寿命；`user_canceled` 和 `general_error` 差异见总体状态 |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | TLS/DTLS 部署安全 BCP | ticket 使用 AES-256-GCM，寿命限制为 1 秒至 7 天；RSA 2048 位及 SHA-1/MD5 证书下限在完整握手、信任锚和恢复路径统一执行；OCSP 和 DTLS 1.2 范围例外见总体状态 |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | 服务身份校验 | DNS-ID/IP-ID 默认严格验证；其他 reference identifier 需要调用方实现应用语义 |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | 双向 record size 协商 | 未实现；PMTU API 只解决路径报文大小，不表示对端接收限制 |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | 更新 RFC 9146/9147，为 CID 地址变化定义 Return Routability Check | 未实现；库不会仅凭新地址上的合法 CID record 自动 rebind，应用仍需自己的路径验证策略 |
| [RFC 8701](https://www.rfc-editor.org/rfc/rfc8701) | GREASE 抗僵化 | 接收端容忍合法未知值并保持 HRR 不变量；发送端不主动生成 GREASE |
| [RFC 8879](https://www.rfc-editor.org/rfc/rfc8879) | Certificate Compression | 未实现；大证书仍可使用现有 DTLS handshake 分片和重传 |
| [RFC 9149](https://www.rfc-editor.org/rfc/rfc9149) | Ticket Requests | 未实现；服务端当前每连接最多签发一个 outstanding ticket |
| [RFC 9257](https://www.rfc-editor.org/rfc/rfc9257) / [RFC 9258](https://www.rfc-editor.org/rfc/rfc9258) | external PSK / PSK importer | 未实现；当前只支持本库签发的 resumption PSK |
| [RFC 9849](https://www.rfc-editor.org/rfc/rfc9849) | Encrypted ClientHello | 未实现 |
| [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954) | Hybrid Key Exchange | 未实现；当前仅支持 X25519 和 P-256 |

### 范围边界

以下项目不降低 RFC 9147 强制语义完成度，但使用者应明确其边界：

- 本模块只实现 DTLS 1.3，不提供 DTLS 1.2 回退。
- 因上一条范围选择，本模块不声称完整符合 RFC 9325 对通用实现支持 DTLS 1.2 的要求。
- Heartbeat 的 record demux 已实现；完整 Heartbeat 协议由 RFC 6520 定义，不属于 RFC 9147 范围。
- 发送端采用一条 record 一个 UDP datagram 的合法模式，未暴露可选的多 record 聚合 API。
- 未暴露并行多个 NewSessionTicket 或 PHA 请求；RFC 允许但不要求这些并行能力。
- CID 可从新来源地址路由和认证 record，但不会自动更新回包地址；RFC 9853 RRC 尚未实现。
- `record_size_limit`、Certificate Compression、external PSK/importer、ECH、hybrid key exchange 和 Ticket Requests 等后续扩展尚未实现。
- wolfSSL 5.9.2 互通构建未启用 CID、0-RTT、session ticket 或 `SESSION_CERTS`，因此 mTLS 恢复等组合当前只有本实现端到端测试和标准库 TLS 1.3 协议层对照。

### 开发优先级

| 优先级 | 工作项 | 完成条件摘要 |
| --- | --- | --- |
| P0 | RFC 9846 alert 语义 | 正确处理 `user_canceled`，并为 `general_error` 增加命名值和发送策略 |
| P1 | RFC 9853 RRC | CID 客户端默认提供 RRC，完整实现路径挑战、放大限制、计时和 NAT rebind 测试 |
| P1 | RFC 8449 `record_size_limit` | 完成 CH/EE 协商、双向限制、HRR/恢复保持和超限处理 |
| P1 | 第三方互通矩阵 | 在对端实际支持时补齐 CID、KeyUpdate、0-RTT、恢复和 PHA 的双向证据 |

其余扩展按实际部署需求和第三方支持推进，不为尚无调用方或标准算法的能力预建空框架。

## Benchmark

以下为 2026-07-29 的代表值，环境为 AMD Ryzen 7 7435H、Go 1.26.3、Windows/amd64。数值用于观察回归，不是跨机器性能保证。

| 场景 | 代表结果 |
| --- | --- |
| 普通完整证书握手和关闭 `BenchmarkConnectionHandshakeLifecycle` | A/B 交替 5 轮中位数约 860 us/op，117465 B/op，877 allocs/op；修改前冻结二进制约 911 us/op，117251 B/op，877 allocs/op |
| 完整 mTLS `BenchmarkMutualTLSHandshakeLifecycle/Full` | 10 轮中位数约 1.265 ms/op，141801 B/op，1171 allocs/op |
| mTLS 恢复 `BenchmarkMutualTLSHandshakeLifecycle/Resumed` | 10 轮中位数约 815 us/op，118941 B/op，833 allocs/op；相对完整 mTLS 约少 35.6% 时间、16.1% 字节和 28.9% 分配 |
| AES-128-GCM 1200 B seal | 约 1.50-1.51 GB/s，1 alloc/op |
| AES-128-GCM 1200 B in-place round trip | 约 862-954 MB/s，1 alloc/op |
| 未认证 record 错误分类 | 约 12.4-13.7 ns/op，0 allocs/op |
| Extension marshal | 约 316.5-358.0 ns/op，128 B/op，1 alloc/op |
| Extension ordered-view parse | 约 69.8-80.1 ns/op，0 B/op，0 allocs/op |
| 单 key_share caller-storage parse | 约 33.5-35.0 ns/op，0 B/op，0 allocs/op |
| ClientHello marshal | 约 436-519 ns/op，424 B/op，7 allocs/op |
| ServerHello marshal | 约 72-91 ns/op，112 B/op，1 alloc/op |
| 1200 B 单分片握手重组 | 约 0.47-0.60 us/op，1280 B/op，1 alloc/op |
| 4 KiB / MTU 1200 protected flight 构造 | 约 3.00-3.26 us/op，5616 B/op，6 allocs/op |
| 4 KiB / MTU 1200 plain flight 构造 | 约 2.15-2.58 us/op，5040 B/op，9 allocs/op |

完整连接 benchmark 是性能验收的主门禁。本轮普通完整连接相对冻结基线没有增加分配次数，字节中位数差约 0.2% 且落在运行抖动内，时延中位数约低 5.7%；没有出现“微基准更好、完整连接更差”。

运行全部 benchmark：

```sh
go test -run '^$' -bench . -benchmem
```

单独复测完整连接与记录层：

```sh
go test -run '^$' -bench '^BenchmarkConnectionHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkMutualTLSHandshakeLifecycle/(Full|Resumed)$' -benchmem -count=10
go test -run '^$' -bench '^BenchmarkProtectedRecord(Seal|RoundTripInPlace)$' -benchmem -count=5
```

仓库还包含 cipher suite、ACK、record/parser、transcript、key schedule、KeyUpdate、握手消息、重组和 flight 构造等细分 benchmark。比较修改前后结果时，应固定 Go 版本、CPU、`-cpu` 和 `-benchtime`，交替运行多个样本，并同时观察 `ns/op`、`B/op`、`allocs/op` 与完整连接 profile。

## 验证

常规回归：

```sh
go test ./... -count=10 -timeout=10m
go vet ./...
golangci-lint run ./...
```

GitHub Actions 会在每次 push 和 pull request 上分别执行 lint、test、vet 和 race。CI 固定使用 `go.mod` 指定的 Go 版本及 golangci-lint v2.12.2。

Windows race：

```powershell
.\tools\test-race.ps1
```

wolfSSL 双向 UDP 互通：

```powershell
$env:WOLFSSL_ROOT = 'C:\path\to\wolfssl'
go test -run TestInteropWolfSSL -v -count=10
```

最近记录的验证结果：

- RFC 9325 证书策略专项覆盖服务端配置、客户端接收、自签名、未发送信任锚、`InsecureSkipVerify`、普通恢复和 mTLS 恢复，并与 `crypto/x509` 对 1024 位 RSA/SHA-1 trust anchor 的行为做差分。
- 全量测试 `-count=10`、`go vet`、golangci-lint v2.12.2 和 Windows race 于 2026-07-29 通过。
- 综合双向丢包、延迟、乱序和重复，以及 CH/SH/Finished/ACK/HRR/mTLS 恢复组合，连续 100 次通过。
- mTLS 完整握手、PSK 恢复、0-RTT、CA/策略回退、ticket 续签认证寿命以及现有 CID、KeyUpdate、exporter 和 PHA 测试通过。
- parser/record fuzz 累计覆盖数百万输入，并包含四套 AEAD 的复制与原地解密差分。
- wolfSSL 5.9.2 双向互通于 2026-07-29 四个方向/套件用例各 `-count=10` 通过，覆盖 HRR、RSA-PSS 证书握手、Finished ACK、应用数据、AES-GCM 和 AES-128-CCM；该构建未启用 session ticket/`SESSION_CERTS`，未强行声称 mTLS 恢复互通。

协议行为变更必须同步更新本文的完成度、验证证据、性能数据、已知差异和第三方互通边界。

## 许可证

本项目以 [GNU General Public License v3.0](LICENSE)（GPL-3.0-only）发布。
