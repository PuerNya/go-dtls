# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-19T21:23:10Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `5e42831be0bd6d1655fb9033952e5d00015438be (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 612.24 us/op | - | 101321 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 917.917 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 455.507 us/op | - | 123192 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.156 ms/op | - | 123645 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.574 ms/op | - | 142419 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 711.715 us/op | - | 121348 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 977.519 us/op | - | 125649 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 975.678 us/op | - | 125745 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 379.513 us/op | - | 103773 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.085 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.084 ms/op | - | 126045 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.808 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.818 ms/op | - | 162384 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.052 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.051 ms/op | - | 154179 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 926.726 us/op | - | 149311 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 870.28 us/op | - | 152655 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.326 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 104.275 us/op | 1.23 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 115.013 us/op | 20.87 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 108.649 us/op | 1.18 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 114.783 us/op | 20.91 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 20.625 us/op | 6.21 MB/s | 416 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 21.586 us/op | 111.18 MB/s | 1547 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 20.695 us/op | 6.19 MB/s | 415 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 20.987 us/op | 114.35 MB/s | 1540 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 20.595 us/op | 6.22 MB/s | 458 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 21.75 us/op | 110.35 MB/s | 1671 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 21.151 us/op | 6.05 MB/s | 459 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 21.045 us/op | 114.04 MB/s | 1675 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 21.513 us/op | 111.56 MB/s | 1668 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.336 ms/op | - | 88614 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 603.272 us/op | - | 91622 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.012 ms/op | - | 93941 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 403.604 us/op | - | 95306 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 430.938 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 432.6 us/op | - | 109210 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.262 us/op | 39.24 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.686 us/op | 358.94 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.179 ms/conn | 2619384 B/op | 21138 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.165 ms/conn | 2631496 B/op | 21477 allocs/op |
| 完整 mTLS 握手 | 5 | 3.843 ms/conn | 3274160 B/op | 29157 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.883 ms/conn | 3730672 B/op | 30901 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.17 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5532 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.156 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.297 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.548 ms/conn | 4054048 B/op | 39395 allocs/op |
| 会话恢复握手 | 5 | 0.6595 ms/conn | 4651288 B/op | 38434 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8101 ms/conn | 6671712 B/op | 49953 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6671 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.538 ms/conn | 3432256 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.488 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.906 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.004 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.12 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.666 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.655 ms/conn | 1495024 B/op | 11922 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.033 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.002 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.064 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.233 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.601 ms/conn | 1547152 B/op | 12664 allocs/op |
| 会话恢复握手 | 5 | 1.098 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.113 ms/conn | 2538360 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.373 ms/conn | 2006024 B/op | 11252 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.38 ms/conn | 2017720 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.578 ms/conn | 1303936 B/op | 9259 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.515 ms/conn | 2458584 B/op | 10765 allocs/op |
| 完整 mTLS 握手 | 5 | 6.619 ms/conn | 1703696 B/op | 14581 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.507 ms/conn | 1906192 B/op | 15101 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.592 ms/conn | 1368928 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.99 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.537 ms/conn | 2473224 B/op | 11327 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.614 ms/conn | 2634104 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.96 ms/conn | 3308288 B/op | 22331 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3598800 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242560 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.806 ms/conn | 1572752 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.887 ms/conn | 1592592 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.287 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.3 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.385 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.325 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.42 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.038 ms/conn | 452696 B/op | 3438 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.1 ms/conn | 518888 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.074 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9928 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.277 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.078 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.326 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.127 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.529 ms/conn | 537352 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.032 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.122 ms/conn | 487880 B/op | 3665 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.286 ms/conn | 534168 B/op | 4204 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6569 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.719 ms/pair | 869448 B/op | 7181 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.191 ms/pair | 1047432 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.737 ms/pair | 868480 B/op | 7130 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.452 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.278 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.373 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.371 ms/conn | 215664 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.366 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.969 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.683 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.162 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.633 ms/conn | 275184 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.347 ms/conn | 222928 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.385 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.357 ms/conn | 213520 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.484 ms/conn | 194736 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8359 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.414 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.478 ms/pair | 459568 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.561 ms/conn | 203632 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 3.012 ms/pair | 420976 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.568 ms/conn | 481192 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.787 ms/conn | 512872 B/op | 1938 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.59 ms/conn | 453160 B/op | 1921 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.257 ms/conn | 484040 B/op | 1962 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.911 ms/conn | 486472 B/op | 2001 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.478 ms/conn | 514184 B/op | 2134 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 7.072 ms/conn | 538424 B/op | 2037 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.608 ms/conn | 514472 B/op | 2121 allocs/op |
| 完整 mTLS 握手 | 5 | 5.121 ms/conn | 523944 B/op | 2514 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.563 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.687 ms/conn | 516264 B/op | 2271 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.686 ms/conn | 506472 B/op | 2151 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.599 ms/conn | 460072 B/op | 1837 allocs/op |
| 会话恢复握手 | 5 | 5.59 ms/pair | 703592 B/op | 3826 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.378 ms/pair | 802888 B/op | 4632 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.708 ms/pair | 705904 B/op | 3782 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.59 ms/conn | 517224 B/op | 2299 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.736 ms/pair | 740392 B/op | 4190 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.845 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.889 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.806 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 6.17 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 10.01 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 6.424 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 8.197 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.956 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.807 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.811 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.828 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.687 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.914 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 8.089 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 5.601 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.791 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.079 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.128 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.038 ms/conn | 389976 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9958 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.338 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.057 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.333 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.083 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.51 ms/conn | 537224 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.074 ms/conn | 454472 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.108 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.282 ms/conn | 534168 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8612 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6174 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.761 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.213 ms/pair | 1047176 B/op | 8531 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.491 ms/conn | 650264 B/op | 5440 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.432 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.24 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8522 ms/conn | 180192 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.8168 ms/conn | 213472 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8539 ms/conn | 145248 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.7522 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.699 ms/conn | 185280 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.9177 ms/conn | 203872 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.082 ms/conn | 270944 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.8878 ms/conn | 180640 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 1.141 ms/conn | 220320 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8102 ms/conn | 180384 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.7846 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.064 ms/conn | 229472 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7498 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5217 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.466 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.63 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1 ms/conn | 241280 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.9112 ms/conn | 201248 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.725 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.718 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.744 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.646 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.635 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.518 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.635 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.905 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.692 ms/conn | 509864 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.933 ms/conn | 516552 B/op | 2341 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.662 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.677 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.801 ms/conn | 502152 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.565 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.379 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.337 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.706 ms/pair | 789272 B/op | 4370 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 3.097 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.311 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.386 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.307 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.329 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.636 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.408 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.478 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.319 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.471 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.315 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.309 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.425 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.331 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.174 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.804 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.014 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.466 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.079 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.11 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.018 ms/conn | 390104 B/op | 3400 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.055 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.992 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.302 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.047 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.422 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.1 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.496 ms/conn | 537480 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.018 ms/conn | 454472 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.07 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.263 ms/conn | 533016 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8398 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6754 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.749 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.197 ms/pair | 1047432 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.755 ms/pair | 868352 B/op | 7130 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 2.011 ms/conn | 641608 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.069 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.187 ms/conn | 245216 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.117 ms/conn | 176224 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.135 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.133 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.602 ms/conn | 220192 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.08 ms/conn | 223712 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.474 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.023 ms/conn | 211872 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 1.381 ms/conn | 267136 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.156 ms/conn | 214560 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.075 ms/conn | 247120 B/op | 1985 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.384 ms/conn | 270112 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.978 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.827 ms/conn | 199328 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 1.953 ms/pair | 471648 B/op | 3879 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.68 ms/pair | 593024 B/op | 4442 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.31 ms/conn | 270400 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.651 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.791 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.626 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.751 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.667 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.314 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.663 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.954 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.667 ms/conn | 514312 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.99 ms/conn | 521448 B/op | 2381 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.573 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.78 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.816 ms/conn | 506088 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.124 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.292 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 6.846 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.034 ms/pair | 801320 B/op | 4506 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 6.495 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.469 ms/conn | 574344 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.635 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.76 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.71 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.621 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.681 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.384 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.696 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.965 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.598 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.815 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.599 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.683 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.838 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.128 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.408 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.861 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.344 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.02 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 14.79 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 16.82 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 15.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 16.83 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 15.89 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 12.34 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 25.83 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 31.31 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 19.92 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 12.85 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 31.99 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 33.58 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.475 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 2.805 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 36.33 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 18.13 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 2.11 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.843 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 50.77 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 19.73 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 236.1 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 202 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 207.6 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.406 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 207.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 374.4 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 35.86 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.07 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 502.6 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.096 us/op | 1954.42 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.588 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 205.9 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 174.8 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 896.4 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.635 us/op | 1126.67 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 318.5 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 92.87 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 354.2 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 20.77 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.26 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.931 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.151 us/op | 557.85 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 859.8 ns/op | 1395.7 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.5 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.244 us/op | 534.66 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.285 us/op | 933.66 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.372 us/op | 162.77 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.119 us/op | 566.23 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.231 us/op | 537.84 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.137 us/op | 382.48 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 933.4 ns/op | 1285.58 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.233 us/op | 371.12 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 854.1 ns/op | 1405.04 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 901.5 ns/op | 1331.12 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.351 us/op | 888.03 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
