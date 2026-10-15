# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-15T18:58:30Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `d139ee0c3af846be4bb5925ff8a4ca90391b4460 (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 619.971 us/op | - | 101322 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 924.182 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 462.229 us/op | - | 123192 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.169 ms/op | - | 123642 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.599 ms/op | - | 142422 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 722.426 us/op | - | 121345 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 989.072 us/op | - | 125649 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 986.011 us/op | - | 125745 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 384.214 us/op | - | 103773 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.096 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.097 ms/op | - | 126046 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.819 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.822 ms/op | - | 162448 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.06 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.064 ms/op | - | 154179 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 932.834 us/op | - | 149311 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 878.21 us/op | - | 152655 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.351 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 107.528 us/op | 1.19 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 115.619 us/op | 20.76 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 108.664 us/op | 1.18 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 116.449 us/op | 20.61 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 20.714 us/op | 6.18 MB/s | 416 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 21.877 us/op | 109.7 MB/s | 1546 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 20.753 us/op | 6.17 MB/s | 417 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 21.82 us/op | 109.99 MB/s | 1546 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 20.801 us/op | 6.15 MB/s | 460 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 21.389 us/op | 112.2 MB/s | 1673 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 20.907 us/op | 6.12 MB/s | 461 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 21.235 us/op | 113.02 MB/s | 1673 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 21.933 us/op | 109.42 MB/s | 1666 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.357 ms/op | - | 88614 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 645.128 us/op | - | 90353 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.018 ms/op | - | 93941 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 427.294 us/op | - | 95409 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 439.438 us/op | - | 109309 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 440.655 us/op | - | 109177 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.394 us/op | 37.71 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.902 us/op | 347.75 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.154 ms/conn | 2619432 B/op | 21139 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.18 ms/conn | 2631496 B/op | 21477 allocs/op |
| 完整 mTLS 握手 | 5 | 3.931 ms/conn | 3274208 B/op | 29158 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.921 ms/conn | 3730624 B/op | 30900 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.215 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5428 ms/conn | 1752648 B/op | 14597 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.185 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.335 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.629 ms/conn | 4053952 B/op | 39391 allocs/op |
| 会话恢复握手 | 5 | 0.6854 ms/conn | 4650808 B/op | 38436 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8345 ms/conn | 6671712 B/op | 49952 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6484 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.574 ms/conn | 3432272 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.498 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.934 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.091 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.202 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.712 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.627 ms/conn | 1495024 B/op | 11922 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.044 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.003 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.056 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.252 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.648 ms/conn | 1547280 B/op | 12666 allocs/op |
| 会话恢复握手 | 5 | 1.086 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.112 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.347 ms/conn | 2005320 B/op | 11245 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.39 ms/conn | 2017720 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.626 ms/conn | 1304160 B/op | 9259 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.856 ms/conn | 2459544 B/op | 10765 allocs/op |
| 完整 mTLS 握手 | 5 | 6.76 ms/conn | 1703280 B/op | 14581 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.629 ms/conn | 1906304 B/op | 15101 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.512 ms/conn | 1369120 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.993 ms/conn | 1026288 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.893 ms/conn | 2472840 B/op | 11327 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 6.009 ms/conn | 2634424 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 14.51 ms/conn | 3308592 B/op | 22331 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3599248 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.792 ms/conn | 1572432 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.92 ms/conn | 1592752 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.305 ms/conn | 1751504 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.757 ms/conn | 1169864 B/op | 1203 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 9.017 ms/conn | 1183560 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.743 ms/conn | 1170184 B/op | 1203 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.91 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.068 ms/conn | 452696 B/op | 3437 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.116 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.036 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.015 ms/conn | 459256 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.283 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.055 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.299 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.064 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.499 ms/conn | 537480 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.032 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.075 ms/conn | 487368 B/op | 3664 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.256 ms/conn | 533016 B/op | 4204 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7045 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.697 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.194 ms/pair | 1047176 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.738 ms/pair | 868480 B/op | 7130 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.412 ms/conn | 511448 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.236 ms/pair | 947768 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.373 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.404 ms/conn | 215600 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.38 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.977 ms/conn | 184752 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.546 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.257 ms/conn | 215536 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.606 ms/conn | 275056 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.351 ms/conn | 222992 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.346 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.377 ms/conn | 213520 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.388 ms/conn | 194608 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8156 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.454 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.364 ms/pair | 459440 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.52 ms/conn | 203504 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.709 ms/pair | 420848 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.521 ms/conn | 481192 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.558 ms/conn | 512840 B/op | 1939 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.499 ms/conn | 453160 B/op | 1921 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.127 ms/conn | 484040 B/op | 1963 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.779 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.289 ms/conn | 514168 B/op | 2134 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.891 ms/conn | 538456 B/op | 2038 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.477 ms/conn | 514216 B/op | 2121 allocs/op |
| 完整 mTLS 握手 | 5 | 5.176 ms/conn | 523880 B/op | 2510 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.512 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.589 ms/conn | 516264 B/op | 2271 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.517 ms/conn | 505928 B/op | 2149 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.491 ms/conn | 460072 B/op | 1837 allocs/op |
| 会话恢复握手 | 5 | 5.495 ms/pair | 703624 B/op | 3828 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.078 ms/pair | 802792 B/op | 4634 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.459 ms/pair | 705936 B/op | 3782 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.496 ms/conn | 517192 B/op | 2297 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.454 ms/pair | 740424 B/op | 4192 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.703 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.833 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.725 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 6.079 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.911 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 6.312 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 8.104 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.852 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.712 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.821 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.721 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.656 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.8 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.925 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 5.48 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.644 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.021 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.071 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.042 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9531 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.3 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.083 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.333 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.117 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.471 ms/conn | 537224 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.019 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.074 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.264 ms/conn | 534040 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8531 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6199 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.718 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.215 ms/pair | 1047048 B/op | 8531 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.47 ms/conn | 650136 B/op | 5440 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.401 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.207 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.7855 ms/conn | 180192 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9079 ms/conn | 213472 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.83 ms/conn | 145376 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.7928 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.662 ms/conn | 185408 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.8495 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.005 ms/conn | 270816 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.7917 ms/conn | 180640 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 1.118 ms/conn | 220320 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8214 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.8895 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.9799 ms/conn | 229344 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7558 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5663 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.309 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.704 ms/pair | 473968 B/op | 3656 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.046 ms/conn | 241152 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.8678 ms/conn | 201248 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.554 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.64 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.687 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.635 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.599 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.465 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.621 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.811 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.625 ms/conn | 509864 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.903 ms/conn | 516552 B/op | 2341 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.637 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.623 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.728 ms/conn | 502152 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.522 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.308 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.292 ms/pair | 693064 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.623 ms/pair | 789256 B/op | 4370 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.979 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.304 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.366 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.281 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.272 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.678 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.36 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.436 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.283 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.44 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.298 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.26 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.39 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.292 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.097 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.76 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.984 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.421 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.015 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.123 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.019 ms/conn | 390104 B/op | 3400 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.09 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.949 ms/conn | 459256 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.285 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.066 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.315 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.102 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.494 ms/conn | 537352 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.034 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.077 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.269 ms/conn | 533144 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8171 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6358 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.722 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.206 ms/pair | 1047224 B/op | 8532 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.75 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 2.001 ms/conn | 641480 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.072 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.17 ms/conn | 245088 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.104 ms/conn | 176224 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.127 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.137 ms/conn | 220384 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.732 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.076 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.454 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.035 ms/conn | 211872 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 1.419 ms/conn | 267008 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.018 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.116 ms/conn | 247136 B/op | 1986 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.418 ms/conn | 270240 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.9617 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7956 ms/conn | 199184 B/op | 1679 allocs/op |
| 会话恢复握手 | 5 | 1.963 ms/pair | 471360 B/op | 3878 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.667 ms/pair | 592896 B/op | 4442 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.388 ms/conn | 270400 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.465 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.556 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.489 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.432 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.489 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.062 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.575 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.664 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.526 ms/conn | 514184 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.786 ms/conn | 521320 B/op | 2381 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.504 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.587 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.668 ms/conn | 506088 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.985 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.165 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 6.393 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.855 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 6.338 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.268 ms/conn | 574344 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.473 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.513 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.482 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.476 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.561 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.232 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.584 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.797 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.464 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.773 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.513 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.588 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.712 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.028 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.268 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.477 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.123 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.801 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 14.82 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 16.86 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 15.52 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 16.84 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 15.91 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 12.36 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 25.88 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 31.18 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 19.47 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 12.84 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 31.6 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 33.62 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.443 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 2.737 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 37.38 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 18.13 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 2.094 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.869 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 51.92 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 22.95 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 246.3 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 207.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 211 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.475 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 211.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 376.2 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 36.02 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.38 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 511.9 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.122 us/op | 1930.49 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.602 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 202.6 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 175.7 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 903.2 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.642 us/op | 1124.63 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 321.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 91 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 359 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 20.82 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.24 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.937 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.174 us/op | 552.05 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 868 ns/op | 1382.42 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.5 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.242 us/op | 535.16 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.317 us/op | 910.88 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.424 us/op | 161.63 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.146 us/op | 559.12 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.24 us/op | 535.63 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.183 us/op | 377.06 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 955.9 ns/op | 1255.4 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.255 us/op | 368.63 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 865.9 ns/op | 1385.92 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 912 ns/op | 1315.84 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.363 us/op | 880.41 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
