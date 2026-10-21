# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-21T13:53:55Z`
- Go: `go version go1.27.3 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `4b4e21ad2335e47b0016be89648b076e1effba18 (Linux Release static)`
- NSS: `b9283c1e4ee232b4450d36115272af670d38ec81 (Linux Release)`
- BoringSSL: `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 (Linux Release)`
- OpenSSL: `41eecc7c7cd4069b1b29bca49b9d97825c6bf648 (Linux Release)`

共 378 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数按每次操作统计；Listener UDP 汇总独立客户端和服务端进程，排除中转器。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (39)](#section-connection-lifecycle)
- [真实 UDP 互通 (284)](#section-real-udp-interoperability)
  - [go-dtls 客户端 -> go-dtls 服务端 (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls 客户端 -> wolfSSL 服务端 (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL 客户端 -> go-dtls 服务端 (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL 客户端 -> wolfSSL 服务端 (15)](#real-udp-wolfssl-client-wolfssl-server)
  - [NSS / go-dtls 客户端 -> go-dtls 服务端 (18)](#real-udp-nss-go-go)
  - [NSS / go-dtls 客户端 -> NSS 服务端 (18)](#real-udp-nss-go-peer)
  - [NSS / NSS 客户端 -> go-dtls 服务端 (18)](#real-udp-nss-peer-go)
  - [NSS / NSS 客户端 -> NSS 服务端 (18)](#real-udp-nss-peer-peer)
  - [BoringSSL / go-dtls 客户端 -> go-dtls 服务端 (19)](#real-udp-boringssl-go-go)
  - [BoringSSL / go-dtls 客户端 -> BoringSSL 服务端 (19)](#real-udp-boringssl-go-peer)
  - [BoringSSL / BoringSSL 客户端 -> go-dtls 服务端 (19)](#real-udp-boringssl-peer-go)
  - [BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端 (19)](#real-udp-boringssl-peer-peer)
  - [OpenSSL / go-dtls 客户端 -> go-dtls 服务端 (19)](#real-udp-openssl-go-go)
  - [OpenSSL / go-dtls 客户端 -> OpenSSL 服务端 (19)](#real-udp-openssl-go-peer)
  - [OpenSSL / OpenSSL 客户端 -> go-dtls 服务端 (19)](#real-udp-openssl-peer-go)
  - [OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端 (19)](#real-udp-openssl-peer-peer)
- [记录层与可靠性 (55)](#section-record-layer-and-reliability)

<a id="section-connection-lifecycle"></a>
## 连接生命周期

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 596.043 us/op | - | 101386 B/op | 733 allocs/op |
| 完整 mTLS 握手 | 5 | 896.335 us/op | - | 116462 B/op | 931 allocs/op |
| mTLS 会话恢复握手 | 5 | 459.03 us/op | - | 123258 B/op | 856 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.125 ms/op | - | 123773 B/op | 1110 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.529 ms/op | - | 142620 B/op | 1407 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 691.95 us/op | - | 121411 B/op | 963 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 947.321 us/op | - | 125777 B/op | 1002 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 954.001 us/op | - | 125873 B/op | 1002 allocs/op |
| 直接外部 PSK 握手 | 5 | 382.144 us/op | - | 103774 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.054 ms/op | - | 136156 B/op | 1018 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.047 ms/op | - | 126110 B/op | 994 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.746 ms/op | - | 177281 B/op | 1471 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.749 ms/op | - | 162512 B/op | 1424 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.011 ms/op | - | 151467 B/op | 1233 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.012 ms/op | - | 154243 B/op | 1254 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 892.215 us/op | - | 149376 B/op | 766 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 850.877 us/op | - | 152720 B/op | 797 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.158 ms/op | - | 178167 B/op | 815 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 141.147 us/op | 0.91 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 154.398 us/op | 15.54 MB/s | 1534 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 143.09 us/op | 0.89 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 155.798 us/op | 15.4 MB/s | 1534 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 20.96 us/op | 6.11 MB/s | 420 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 19.269 us/op | 124.55 MB/s | 1544 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 20.805 us/op | 6.15 MB/s | 417 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 19.669 us/op | 122.02 MB/s | 1542 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 19.734 us/op | 6.49 MB/s | 457 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 19.389 us/op | 123.78 MB/s | 1659 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 20.221 us/op | 6.33 MB/s | 456 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 19.469 us/op | 123.27 MB/s | 1664 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 19.473 us/op | 123.25 MB/s | 1659 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.463 ms/op | - | 88679 B/op | 768 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 667.718 us/op | - | 90601 B/op | 768 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.173 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 451.532 us/op | - | 95776 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 430.883 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 443.565 us/op | - | 109177 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.313 us/op | 38.64 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.616 us/op | 362.76 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.002 ms/conn | 2619384 B/op | 21138 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.993 ms/conn | 2631544 B/op | 21478 allocs/op |
| 完整 mTLS 握手 | 5 | 3.589 ms/conn | 3275440 B/op | 29177 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.596 ms/conn | 3731904 B/op | 30920 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.005 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5608 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.986 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.154 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.137 ms/conn | 4056512 B/op | 39435 allocs/op |
| 会话恢复握手 | 5 | 0.6588 ms/conn | 4650808 B/op | 38436 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.7939 ms/conn | 6674272 B/op | 49993 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6507 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.328 ms/conn | 3432256 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.304 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.554 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.824 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.885 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.37 ms/conn | 1496624 B/op | 11938 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.366 ms/conn | 1495176 B/op | 11924 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.8 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.027 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.768 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.954 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.29 ms/conn | 1547088 B/op | 12663 allocs/op |
| 会话恢复握手 | 5 | 1.185 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.197 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.895 ms/conn | 2007048 B/op | 11263 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.089 ms/conn | 2017720 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.491 ms/conn | 1304032 B/op | 9259 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.227 ms/conn | 2459224 B/op | 10766 allocs/op |
| 完整 mTLS 握手 | 5 | 6.206 ms/conn | 1704784 B/op | 14601 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.215 ms/conn | 1907344 B/op | 15121 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.431 ms/conn | 1368896 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.991 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.169 ms/conn | 2472312 B/op | 11328 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.283 ms/conn | 2633712 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.99 ms/conn | 3310928 B/op | 22371 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3599600 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.611 ms/conn | 1573072 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.671 ms/conn | 1592432 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.83 ms/conn | 1751504 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.032 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.962 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.953 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.64 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64464 B/op | 57 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.015 ms/conn | 452904 B/op | 3441 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.06 ms/conn | 519144 B/op | 3441 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.9946 ms/conn | 390360 B/op | 3404 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9253 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.098 ms/conn | 464568 B/op | 3560 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.016 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.26 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.026 ms/conn | 507488 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.433 ms/conn | 537736 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.9992 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.028 ms/conn | 487800 B/op | 3669 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.206 ms/conn | 533528 B/op | 4212 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6409 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.653 ms/pair | 870296 B/op | 7185 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.165 ms/pair | 1047944 B/op | 8543 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.689 ms/pair | 868736 B/op | 7134 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.352 ms/conn | 511832 B/op | 3992 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.129 ms/pair | 948152 B/op | 7915 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.432 ms/conn | 182128 B/op | 1538 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.364 ms/conn | 215856 B/op | 1538 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.422 ms/conn | 147408 B/op | 1518 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.929 ms/conn | 185136 B/op | 1562 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.337 ms/conn | 187216 B/op | 1598 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.169 ms/conn | 215920 B/op | 1814 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.458 ms/conn | 275440 B/op | 1650 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.182 ms/conn | 223248 B/op | 1914 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.307 ms/conn | 182416 B/op | 1550 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.299 ms/conn | 213904 B/op | 1722 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.237 ms/conn | 195248 B/op | 1646 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.771 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.372 ms/pair | 379920 B/op | 3145 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.555 ms/pair | 459696 B/op | 3521 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.542 ms/conn | 203888 B/op | 1718 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.636 ms/pair | 421232 B/op | 3481 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.166 ms/conn | 481240 B/op | 1939 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.295 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.213 ms/conn | 453160 B/op | 1921 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.748 ms/conn | 484040 B/op | 1963 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.16 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.891 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.279 ms/conn | 538488 B/op | 2040 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.163 ms/conn | 514280 B/op | 2125 allocs/op |
| 完整 mTLS 握手 | 5 | 4.786 ms/conn | 523976 B/op | 2513 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.138 ms/conn | 482696 B/op | 1973 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.286 ms/conn | 516264 B/op | 2271 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.234 ms/conn | 505960 B/op | 2151 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.37 ms/conn | 460120 B/op | 1839 allocs/op |
| 会话恢复握手 | 5 | 5.113 ms/pair | 703688 B/op | 3832 allocs/op |
| mTLS 会话恢复握手 | 5 | 5.773 ms/pair | 803336 B/op | 4644 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.377 ms/pair | 706000 B/op | 3786 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.066 ms/conn | 517224 B/op | 2299 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.062 ms/pair | 740456 B/op | 4194 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.291 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.492 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.367 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.564 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.12 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.772 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 7.36 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.363 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.31 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.393 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.307 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.395 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.444 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.431 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 4.982 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.129 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.009 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.033 ms/conn | 519096 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.9877 ms/conn | 390360 B/op | 3404 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9431 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.114 ms/conn | 464440 B/op | 3560 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.244 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.04 ms/conn | 507488 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.397 ms/conn | 537736 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.9907 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.012 ms/conn | 485384 B/op | 3660 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.201 ms/conn | 533528 B/op | 4212 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8216 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6424 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.661 ms/pair | 868552 B/op | 7183 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.098 ms/pair | 1047944 B/op | 8543 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.359 ms/conn | 650520 B/op | 5444 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.347 ms/conn | 511704 B/op | 3992 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.105 ms/pair | 948152 B/op | 7915 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.842 ms/conn | 180576 B/op | 1528 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.902 ms/conn | 213728 B/op | 1528 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8353 ms/conn | 145632 B/op | 1507 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.7718 ms/conn | 183584 B/op | 1552 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.415 ms/conn | 185536 B/op | 1588 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.8626 ms/conn | 204256 B/op | 1760 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.9968 ms/conn | 271200 B/op | 1596 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.8085 ms/conn | 180896 B/op | 1560 allocs/op |
| 完整 mTLS 握手 | 5 | 1.099 ms/conn | 220576 B/op | 1904 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.7825 ms/conn | 180768 B/op | 1540 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.8131 ms/conn | 212352 B/op | 1712 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.028 ms/conn | 229344 B/op | 2076 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7685 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5189 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.309 ms/pair | 388896 B/op | 3283 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.568 ms/pair | 474336 B/op | 3659 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.036 ms/conn | 241536 B/op | 2256 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.9255 ms/conn | 201376 B/op | 1708 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.51 ms/pair | 432512 B/op | 3635 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.485 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.513 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.452 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.427 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.062 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.45 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.628 ms/conn | 530072 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.405 ms/conn | 509800 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.678 ms/conn | 516808 B/op | 2345 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.436 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.459 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.543 ms/conn | 501640 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.36 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.214 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.186 ms/pair | 693096 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.495 ms/pair | 789800 B/op | 4378 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.837 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.11 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.193 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.115 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.121 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.204 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.203 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.25 ms/conn | 212512 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.103 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.267 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.126 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.101 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.233 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.107 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 1.977 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.645 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.829 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.283 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.027 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.044 ms/conn | 519096 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.9887 ms/conn | 390360 B/op | 3404 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.004 ms/conn | 517384 B/op | 4126 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9381 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.082 ms/conn | 464568 B/op | 3560 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.996 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.229 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.035 ms/conn | 507488 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.417 ms/conn | 537736 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.9972 ms/conn | 454728 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.021 ms/conn | 485256 B/op | 3660 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.183 ms/conn | 533528 B/op | 4212 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8258 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6052 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.646 ms/pair | 868552 B/op | 7183 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.094 ms/pair | 1048200 B/op | 8543 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.649 ms/pair | 868736 B/op | 7134 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.885 ms/conn | 642248 B/op | 5982 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.092 ms/conn | 211424 B/op | 1780 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.141 ms/conn | 245472 B/op | 1780 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.145 ms/conn | 176480 B/op | 1759 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.061 ms/conn | 264928 B/op | 2320 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.084 ms/conn | 220768 B/op | 1848 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.534 ms/conn | 220544 B/op | 1864 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.109 ms/conn | 224096 B/op | 1896 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.386 ms/conn | 332128 B/op | 1884 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.208 ms/conn | 212128 B/op | 1792 allocs/op |
| 完整 mTLS 握手 | 5 | 1.455 ms/conn | 267264 B/op | 2244 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.225 ms/conn | 214816 B/op | 1816 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.137 ms/conn | 247408 B/op | 1991 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.334 ms/conn | 270624 B/op | 2424 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.026 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9665 ms/conn | 199328 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 2.046 ms/pair | 471488 B/op | 3883 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.689 ms/pair | 592160 B/op | 4444 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.363 ms/conn | 270656 B/op | 2308 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.057 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.107 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.078 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.091 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.061 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.4 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.102 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.271 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.032 ms/conn | 514312 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.422 ms/conn | 521576 B/op | 2385 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.061 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.195 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.179 ms/conn | 506088 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.534 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 4.764 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 5.95 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.391 ms/pair | 801832 B/op | 4514 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.912 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.865 ms/conn | 574856 B/op | 3093 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.057 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.133 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.074 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.03 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.161 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.533 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.184 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.39 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.059 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.271 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.081 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.213 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.26 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.616 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 4.903 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.193 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.691 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.47 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 16.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 15.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 17.2 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 15.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 17.29 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 11.81 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 26.78 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 30.72 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 22.22 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 13.31 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 35.51 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 33.1 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.714 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.255 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 41.37 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 22.24 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 1.999 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 4.169 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 62.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 19.71 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 237.6 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 192.2 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 200.6 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.432 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 195.2 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 358.5 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 36.86 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.15 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 546 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.096 us/op | 1954.31 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.667 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 209.6 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 168.6 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 944.6 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.65 us/op | 1122.18 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 331.6 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 90.03 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 347.8 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 22.53 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.06 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.372 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.299 us/op | 522.01 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 897.5 ns/op | 1337.1 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.59 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.365 us/op | 507.45 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.382 us/op | 868.05 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.145 us/op | 167.95 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.282 us/op | 525.93 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.364 us/op | 507.72 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.086 us/op | 388.85 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 959.8 ns/op | 1250.22 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.088 us/op | 388.55 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 897.2 ns/op | 1337.54 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 935.6 ns/op | 1282.62 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.323 us/op | 906.83 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
