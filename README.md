# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-20T20:18:35Z`
- Go: `go version go1.27.3 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `3846e798c1b52b9bbcc12a1e296ae40eb2a9afbf (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 599.807 us/op | - | 101386 B/op | 733 allocs/op |
| 完整 mTLS 握手 | 5 | 889.543 us/op | - | 116462 B/op | 931 allocs/op |
| mTLS 会话恢复握手 | 5 | 458.3 us/op | - | 123259 B/op | 856 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.136 ms/op | - | 123770 B/op | 1110 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.53 ms/op | - | 142620 B/op | 1407 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 698.593 us/op | - | 121412 B/op | 963 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 947.816 us/op | - | 125777 B/op | 1002 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 951.679 us/op | - | 125873 B/op | 1002 allocs/op |
| 直接外部 PSK 握手 | 5 | 377.718 us/op | - | 103774 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.064 ms/op | - | 136155 B/op | 1018 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.079 ms/op | - | 126110 B/op | 994 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.76 ms/op | - | 177281 B/op | 1471 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.776 ms/op | - | 162512 B/op | 1424 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.012 ms/op | - | 151467 B/op | 1233 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.008 ms/op | - | 154243 B/op | 1254 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 911.417 us/op | - | 149375 B/op | 766 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 860.328 us/op | - | 152720 B/op | 797 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.187 ms/op | - | 178166 B/op | 815 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 153.179 us/op | 0.84 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 164.572 us/op | 14.58 MB/s | 1534 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 154.337 us/op | 0.83 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 164.818 us/op | 14.56 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 22.016 us/op | 5.81 MB/s | 421 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 20.156 us/op | 119.07 MB/s | 1547 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 20.732 us/op | 6.17 MB/s | 421 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 19.889 us/op | 120.67 MB/s | 1547 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 21.391 us/op | 5.98 MB/s | 456 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 20.882 us/op | 114.93 MB/s | 1661 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 20.27 us/op | 6.31 MB/s | 457 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 19.615 us/op | 122.36 MB/s | 1657 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 19.843 us/op | 120.95 MB/s | 1658 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.483 ms/op | - | 88678 B/op | 768 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 652.443 us/op | - | 90332 B/op | 768 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.179 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 449.139 us/op | - | 95407 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 432.87 us/op | - | 109339 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 437.38 us/op | - | 109179 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.432 us/op | 37.3 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.7 us/op | 358.21 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.01 ms/conn | 2619432 B/op | 21139 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.014 ms/conn | 2631544 B/op | 21478 allocs/op |
| 完整 mTLS 握手 | 5 | 3.608 ms/conn | 3275440 B/op | 29177 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.621 ms/conn | 3731904 B/op | 30920 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.009 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5559 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.983 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.177 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.086 ms/conn | 4056560 B/op | 39433 allocs/op |
| 会话恢复握手 | 5 | 0.6607 ms/conn | 4651288 B/op | 38434 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8337 ms/conn | 6674288 B/op | 49992 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.7078 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.374 ms/conn | 3432256 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.337 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.618 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.809 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.913 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.5 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.456 ms/conn | 1495176 B/op | 11924 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.784 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.049 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.844 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.004 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.347 ms/conn | 1547280 B/op | 12666 allocs/op |
| 会话恢复握手 | 5 | 1.207 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.177 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.145 ms/conn | 2006024 B/op | 11252 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.146 ms/conn | 2018504 B/op | 11372 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.459 ms/conn | 1304160 B/op | 9260 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.281 ms/conn | 2458584 B/op | 10765 allocs/op |
| 完整 mTLS 握手 | 5 | 6.201 ms/conn | 1704304 B/op | 14601 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.168 ms/conn | 1907696 B/op | 15121 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.474 ms/conn | 1368672 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.98 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.58 ms/conn | 2473032 B/op | 11327 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.24 ms/conn | 2633784 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.09 ms/conn | 3312208 B/op | 22372 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3600400 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.694 ms/conn | 1572912 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.697 ms/conn | 1592912 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.946 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.012 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.244 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.026 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.71 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1191192 B/op | 1203 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.003 ms/conn | 452776 B/op | 3441 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.095 ms/conn | 519144 B/op | 3441 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.043 ms/conn | 390360 B/op | 3404 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9721 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.102 ms/conn | 464488 B/op | 3561 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.04 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.309 ms/conn | 622520 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.079 ms/conn | 507488 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.469 ms/conn | 537736 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.063 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.089 ms/conn | 487752 B/op | 3668 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.257 ms/conn | 534680 B/op | 4212 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6339 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.695 ms/pair | 869704 B/op | 7185 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.148 ms/pair | 1048200 B/op | 8543 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.762 ms/pair | 868608 B/op | 7134 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.375 ms/conn | 511832 B/op | 3992 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.231 ms/pair | 948024 B/op | 7915 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.472 ms/conn | 182000 B/op | 1538 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.451 ms/conn | 215920 B/op | 1538 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.452 ms/conn | 147408 B/op | 1518 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.959 ms/conn | 185136 B/op | 1562 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.43 ms/conn | 187088 B/op | 1598 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.12 ms/conn | 215920 B/op | 1814 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.686 ms/conn | 275312 B/op | 1650 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.306 ms/conn | 223248 B/op | 1914 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.404 ms/conn | 182416 B/op | 1550 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.449 ms/conn | 213904 B/op | 1722 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.44 ms/conn | 195248 B/op | 1646 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8814 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.466 ms/pair | 379920 B/op | 3145 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.563 ms/pair | 459824 B/op | 3521 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.505 ms/conn | 203888 B/op | 1718 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.841 ms/pair | 421104 B/op | 3481 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.574 ms/conn | 481016 B/op | 1941 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.458 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.433 ms/conn | 453208 B/op | 1923 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.991 ms/conn | 484008 B/op | 1961 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.421 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.172 ms/conn | 514152 B/op | 2135 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.934 ms/conn | 538424 B/op | 2037 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.542 ms/conn | 514216 B/op | 2121 allocs/op |
| 完整 mTLS 握手 | 5 | 4.972 ms/conn | 524072 B/op | 2519 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.58 ms/conn | 482696 B/op | 1973 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.765 ms/conn | 516232 B/op | 2269 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.55 ms/conn | 506504 B/op | 2153 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.552 ms/conn | 460072 B/op | 1837 allocs/op |
| 会话恢复握手 | 5 | 5.536 ms/pair | 703656 B/op | 3830 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.158 ms/pair | 803400 B/op | 4640 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.577 ms/pair | 706032 B/op | 3790 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.377 ms/conn | 517192 B/op | 2297 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.385 ms/pair | 740392 B/op | 4190 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.674 ms/conn | 212448 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.638 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.565 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.777 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.483 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 6.161 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 7.764 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.718 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.678 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.747 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.543 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.575 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.574 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.678 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 5.187 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.446 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.067 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.079 ms/conn | 518968 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.001 ms/conn | 390360 B/op | 3404 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9561 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.121 ms/conn | 464440 B/op | 3560 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.035 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.254 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.07 ms/conn | 507360 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.459 ms/conn | 537992 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.032 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.04 ms/conn | 485256 B/op | 3660 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.241 ms/conn | 533400 B/op | 4212 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8791 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6188 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.692 ms/pair | 868424 B/op | 7183 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.187 ms/pair | 1048072 B/op | 8543 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.445 ms/conn | 650520 B/op | 5444 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.401 ms/conn | 511832 B/op | 3992 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.185 ms/pair | 948152 B/op | 7915 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.7979 ms/conn | 180576 B/op | 1528 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9344 ms/conn | 213600 B/op | 1528 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.7706 ms/conn | 145632 B/op | 1507 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.8269 ms/conn | 183584 B/op | 1552 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.484 ms/conn | 185536 B/op | 1588 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.8693 ms/conn | 204256 B/op | 1760 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.026 ms/conn | 271200 B/op | 1596 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.9511 ms/conn | 180768 B/op | 1560 allocs/op |
| 完整 mTLS 握手 | 5 | 1.071 ms/conn | 220704 B/op | 1904 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.7898 ms/conn | 180768 B/op | 1540 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.9039 ms/conn | 212224 B/op | 1712 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.9891 ms/conn | 229216 B/op | 2076 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7938 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5352 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.45 ms/pair | 388768 B/op | 3283 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.855 ms/pair | 474336 B/op | 3659 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 0.9991 ms/conn | 241536 B/op | 2256 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.9496 ms/conn | 201504 B/op | 1708 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.844 ms/pair | 432512 B/op | 3636 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.52 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.57 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.495 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.57 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.126 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.509 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.814 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.529 ms/conn | 509656 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.795 ms/conn | 516936 B/op | 2345 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.54 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.569 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.652 ms/conn | 501640 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.406 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.296 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.354 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.875 ms/pair | 789912 B/op | 4378 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 3.093 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.206 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.254 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.17 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.172 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.244 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.27 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.317 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.208 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.319 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.2 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.143 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.298 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.161 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.044 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.734 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.146 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.426 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.05 ms/conn | 452728 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.09 ms/conn | 518968 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.039 ms/conn | 390360 B/op | 3404 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.028 ms/conn | 517512 B/op | 4126 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9547 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.154 ms/conn | 464568 B/op | 3560 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.026 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.243 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.074 ms/conn | 507488 B/op | 3819 allocs/op |
| 完整 mTLS 握手 | 5 | 1.456 ms/conn | 537608 B/op | 4564 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.048 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.053 ms/conn | 485256 B/op | 3660 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.229 ms/conn | 534424 B/op | 4212 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8707 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.658 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.712 ms/pair | 868552 B/op | 7183 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.148 ms/pair | 1047944 B/op | 8543 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.713 ms/pair | 868608 B/op | 7134 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.919 ms/conn | 642504 B/op | 5982 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.213 ms/conn | 211296 B/op | 1780 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.163 ms/conn | 245344 B/op | 1780 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.267 ms/conn | 176480 B/op | 1759 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.075 ms/conn | 265056 B/op | 2320 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.121 ms/conn | 220768 B/op | 1848 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.648 ms/conn | 220544 B/op | 1864 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.24 ms/conn | 224096 B/op | 1896 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.527 ms/conn | 332128 B/op | 1884 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.158 ms/conn | 212128 B/op | 1792 allocs/op |
| 完整 mTLS 握手 | 5 | 1.546 ms/conn | 267136 B/op | 2244 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.221 ms/conn | 214944 B/op | 1816 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.222 ms/conn | 247552 B/op | 1992 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.496 ms/conn | 271072 B/op | 2424 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.055 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9268 ms/conn | 199328 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 2.13 ms/pair | 471344 B/op | 3882 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.812 ms/pair | 591024 B/op | 4444 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.411 ms/conn | 270656 B/op | 2308 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.781 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.834 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.701 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.509 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.441 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.981 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.348 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.734 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.569 ms/conn | 514248 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.957 ms/conn | 521576 B/op | 2385 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.382 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.519 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.625 ms/conn | 506600 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.696 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 4.965 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 6.259 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.586 ms/pair | 801704 B/op | 4514 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 6.134 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.11 ms/conn | 574856 B/op | 3093 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.677 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.949 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.451 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.347 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.459 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.156 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.465 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.638 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.429 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.783 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.425 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.415 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.522 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.776 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.007 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.303 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.998 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.688 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 16.56 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 15.92 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 17.19 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 15.93 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 17.79 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 11.81 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 28.92 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 29.66 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 23.06 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 13.34 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 35.94 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 41.08 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.703 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.273 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 41.37 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 18.01 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 2.04 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 4.152 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 66.25 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 29.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 235.7 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 192.6 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 202.1 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.433 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 195.4 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 360.4 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 37.62 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.24 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 540.9 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.171 us/op | 1887.02 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.679 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 208.1 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 169.6 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 936.7 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.632 us/op | 1127.78 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 331.7 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 89.45 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 351 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 22.51 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.37 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.385 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.217 us/op | 541.27 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 868.8 ns/op | 1381.23 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.61 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.371 us/op | 506.04 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.348 us/op | 890.42 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.035 us/op | 170.58 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.205 us/op | 544.14 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.327 us/op | 515.76 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.052 us/op | 393.14 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 953 ns/op | 1259.14 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.037 us/op | 395.11 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 869.7 ns/op | 1379.74 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 929.5 ns/op | 1290.99 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.303 us/op | 921.19 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
