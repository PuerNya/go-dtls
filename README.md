# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-16T13:24:06Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `bc9c0790253382e73198689dbef57b5c5c4792a5 (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 586.776 us/op | - | 101322 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 867.68 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 432.105 us/op | - | 123193 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.09 ms/op | - | 123643 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.497 ms/op | - | 142418 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 681.272 us/op | - | 121347 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 919.177 us/op | - | 125649 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 919.884 us/op | - | 125745 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 368.562 us/op | - | 103774 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.039 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.036 ms/op | - | 126045 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.723 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.725 ms/op | - | 162447 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 985.938 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 986.57 us/op | - | 154179 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 875.528 us/op | - | 149311 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 828.247 us/op | - | 152656 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.129 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 134.456 us/op | 0.95 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 145.969 us/op | 16.44 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 135.645 us/op | 0.94 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 150.414 us/op | 15.96 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 22.081 us/op | 5.8 MB/s | 422 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 20.399 us/op | 117.65 MB/s | 1541 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 18.702 us/op | 6.84 MB/s | 420 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 21.359 us/op | 112.36 MB/s | 1548 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 20.16 us/op | 6.35 MB/s | 452 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 21.158 us/op | 113.43 MB/s | 1666 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 19.34 us/op | 6.62 MB/s | 459 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 21.15 us/op | 113.48 MB/s | 1676 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 20.848 us/op | 115.12 MB/s | 1670 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.416 ms/op | - | 88615 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 614.436 us/op | - | 90550 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.128 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 426.153 us/op | - | 95503 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 413.51 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 418.675 us/op | - | 109177 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.479 us/op | 36.79 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.613 us/op | 362.95 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.985 ms/conn | 2619384 B/op | 21138 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.996 ms/conn | 2631496 B/op | 21477 allocs/op |
| 完整 mTLS 握手 | 5 | 3.528 ms/conn | 3274160 B/op | 29157 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.573 ms/conn | 3730672 B/op | 30901 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5454 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.972 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.112 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.061 ms/conn | 4054000 B/op | 39397 allocs/op |
| 会话恢复握手 | 5 | 0.6342 ms/conn | 4650808 B/op | 38432 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.7982 ms/conn | 6671728 B/op | 49952 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6799 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.346 ms/conn | 3432256 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.342 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.583 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.611 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.695 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.176 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.023 ms/conn | 1495872 B/op | 11931 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.508 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8067 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.519 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.677 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.96 ms/conn | 1547152 B/op | 12664 allocs/op |
| 会话恢复握手 | 5 | 1.13 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.158 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.84 ms/conn | 2006024 B/op | 11252 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.063 ms/conn | 2017880 B/op | 11366 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.416 ms/conn | 1303936 B/op | 9259 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.087 ms/conn | 2458632 B/op | 10766 allocs/op |
| 完整 mTLS 握手 | 5 | 6.117 ms/conn | 1703344 B/op | 14581 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.07 ms/conn | 1906240 B/op | 15101 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.339 ms/conn | 1369344 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.964 ms/conn | 1026736 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.089 ms/conn | 2472840 B/op | 11327 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.141 ms/conn | 2634104 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.69 ms/conn | 3308512 B/op | 22331 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3599360 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242240 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.615 ms/conn | 1572432 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.687 ms/conn | 1592912 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.973 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.72 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63712 B/op | 57 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.719 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.804 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.4 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.9967 ms/conn | 452648 B/op | 3437 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.065 ms/conn | 518920 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.013 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9459 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.124 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.026 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.284 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.1 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.463 ms/conn | 537096 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.045 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.042 ms/conn | 488008 B/op | 3665 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.233 ms/conn | 534040 B/op | 4204 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6686 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.71 ms/pair | 868936 B/op | 7181 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.158 ms/pair | 1047048 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.695 ms/pair | 868480 B/op | 7130 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.396 ms/conn | 511624 B/op | 3989 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.183 ms/pair | 947768 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.372 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.397 ms/conn | 215664 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.481 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.096 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.408 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.144 ms/conn | 215536 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.471 ms/conn | 275056 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.33 ms/conn | 222800 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.345 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.304 ms/conn | 213648 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.382 ms/conn | 194608 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8039 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.484 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.398 ms/pair | 459248 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.527 ms/conn | 203632 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.785 ms/pair | 420848 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.264 ms/conn | 481288 B/op | 1943 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.482 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.731 ms/conn | 453192 B/op | 1923 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.148 ms/conn | 484072 B/op | 1965 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.558 ms/conn | 486408 B/op | 1997 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.327 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.788 ms/conn | 538456 B/op | 2039 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.654 ms/conn | 514504 B/op | 2123 allocs/op |
| 完整 mTLS 握手 | 5 | 5.03 ms/conn | 523720 B/op | 2509 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.369 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.473 ms/conn | 516232 B/op | 2269 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.508 ms/conn | 506440 B/op | 2149 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.528 ms/conn | 460072 B/op | 1837 allocs/op |
| 会话恢复握手 | 5 | 5.72 ms/pair | 703560 B/op | 3824 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.148 ms/pair | 802760 B/op | 4632 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.348 ms/pair | 706000 B/op | 3787 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.441 ms/conn | 517224 B/op | 2299 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.38 ms/pair | 740488 B/op | 4196 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.524 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.826 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.771 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.892 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.529 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 6.179 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 7.713 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.704 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.475 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.55 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.44 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.503 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.626 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.628 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 5.251 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.451 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.034 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.093 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.061 ms/conn | 389976 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9457 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.115 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.035 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.276 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.036 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.464 ms/conn | 537224 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.997 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.067 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.231 ms/conn | 534168 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.841 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6236 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.707 ms/pair | 868168 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.123 ms/pair | 1047560 B/op | 8531 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.403 ms/conn | 650264 B/op | 5440 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.389 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.186 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8244 ms/conn | 180192 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9267 ms/conn | 213344 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8007 ms/conn | 145248 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.8225 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.383 ms/conn | 185280 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.9656 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.9663 ms/conn | 270944 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.7902 ms/conn | 180640 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 1.11 ms/conn | 220320 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8243 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.856 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.003 ms/conn | 229472 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.699 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.609 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.355 ms/pair | 388512 B/op | 3279 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.753 ms/pair | 474208 B/op | 3655 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.058 ms/conn | 241280 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.8869 ms/conn | 201248 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.563 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.523 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.568 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.574 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.51 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.105 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.614 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.748 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.521 ms/conn | 509864 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.841 ms/conn | 516552 B/op | 2341 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.71 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.612 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.817 ms/conn | 502152 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.471 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.259 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.237 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.486 ms/pair | 789384 B/op | 4370 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.916 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.16 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.248 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.196 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.184 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.389 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.221 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.239 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.204 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.39 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.178 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.232 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.318 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.15 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.031 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.692 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.901 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.309 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.001 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.046 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.026 ms/conn | 390104 B/op | 3400 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.019 ms/conn | 517128 B/op | 4122 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9384 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.132 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.009 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.267 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.048 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.408 ms/conn | 537352 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.01 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.041 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.209 ms/conn | 532888 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8338 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6809 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.724 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.171 ms/pair | 1047176 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.708 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.924 ms/conn | 641992 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.114 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.098 ms/conn | 245088 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.097 ms/conn | 176224 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.142 ms/conn | 264672 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.171 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.412 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.201 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.463 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.138 ms/conn | 211872 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 1.454 ms/conn | 267008 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.146 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.209 ms/conn | 247424 B/op | 1988 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.484 ms/conn | 269984 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.021 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9282 ms/conn | 199328 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 2.127 ms/pair | 470944 B/op | 3877 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.731 ms/pair | 592752 B/op | 4441 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.266 ms/conn | 270528 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.303 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.508 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.41 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.234 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.774 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.135 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.596 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.653 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.246 ms/conn | 514248 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.6 ms/conn | 521320 B/op | 2381 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.583 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.732 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.487 ms/conn | 506088 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.981 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.115 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 6.589 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.672 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 6.272 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.287 ms/conn | 574600 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.504 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.743 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.33 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.4 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.768 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.251 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.478 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.971 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.367 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.516 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.589 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.592 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.552 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.179 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.195 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.708 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.008 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.648 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 16.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 16.55 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 17.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 16.52 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 12.19 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 27.77 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 29.62 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 23.03 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 13.41 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 35.21 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 33.13 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.734 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.335 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 39.62 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 18.14 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 2.029 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 4.041 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 61.25 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 28.49 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 232.7 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 188.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 198 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.379 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 198.2 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 363.8 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 36.28 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 44.35 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 521.2 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.007 us/op | 2041.17 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.668 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 211 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 173.4 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 939.6 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.57 us/op | 1147.29 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 319.7 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 89.2 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 343.8 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 19.58 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.39 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.681 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.154 us/op | 557.05 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 863.4 ns/op | 1389.83 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 16.49 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.225 us/op | 539.21 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.27 us/op | 944.75 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.026 us/op | 170.79 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.114 us/op | 567.69 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.222 us/op | 540.13 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 2.96 us/op | 405.41 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 920 ns/op | 1304.36 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.044 us/op | 394.27 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 844.6 ns/op | 1420.82 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 899.8 ns/op | 1333.6 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.26 us/op | 952.49 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
