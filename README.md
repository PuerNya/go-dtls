# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12a5008a3c8aa71239d32fd386029316e649de0e`
- 生成时间: `2026-10-16T02:26:53Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `eb113cac674a9f7da61a70a6307a7db959753ccf (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 592.084 us/op | - | 101321 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 891.064 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 454.692 us/op | - | 123191 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.12 ms/op | - | 123641 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.528 ms/op | - | 142428 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 692.969 us/op | - | 121347 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 945.083 us/op | - | 125649 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 941.737 us/op | - | 125745 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 378.804 us/op | - | 103774 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.037 ms/op | - | 136092 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.034 ms/op | - | 126047 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.726 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.736 ms/op | - | 162385 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.007 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.008 ms/op | - | 154180 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 886.731 us/op | - | 149311 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 837.811 us/op | - | 152656 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.144 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 141.33 us/op | 0.91 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 154.932 us/op | 15.49 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 141.944 us/op | 0.9 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 157.811 us/op | 15.21 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 21.73 us/op | 5.89 MB/s | 421 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 20.834 us/op | 115.2 MB/s | 1543 B/op | 14 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 20.621 us/op | 6.21 MB/s | 419 B/op | 15 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 21.382 us/op | 112.25 MB/s | 1544 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 22.167 us/op | 5.77 MB/s | 465 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 21.392 us/op | 112.19 MB/s | 1669 B/op | 14 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 21.64 us/op | 5.91 MB/s | 494 B/op | 15 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 21.269 us/op | 112.84 MB/s | 1680 B/op | 14 allocs/op |
| Listener UDP 往返 / 服务端无读期限 / 1024 个连接 / CID / 1200 字节 | 5 | 21.53 us/op | 111.47 MB/s | 1676 B/op | 14 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 1 | 5 | 1.449 ms/op | - | 88615 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 完整握手 / 并发 64 | 5 | 649.226 us/op | - | 90259 B/op | 767 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 1 | 5 | 1.149 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP 握手、票据与关闭 / 会话恢复 / 并发 64 | 5 | 439.07 us/op | - | 95363 B/op | 805 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 428.143 us/op | - | 109339 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 427.446 us/op | - | 109181 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.507 us/op | 36.5 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.923 us/op | 346.68 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.016 ms/conn | 2619384 B/op | 21138 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.032 ms/conn | 2631496 B/op | 21477 allocs/op |
| 完整 mTLS 握手 | 5 | 3.556 ms/conn | 3274208 B/op | 29158 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.58 ms/conn | 3730672 B/op | 30901 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.016 ms/conn | 3097256 B/op | 24988 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5547 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.98 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.146 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.114 ms/conn | 4053952 B/op | 39396 allocs/op |
| 会话恢复握手 | 5 | 0.6512 ms/conn | 4650808 B/op | 38436 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8137 ms/conn | 6671712 B/op | 49952 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6663 ms/conn | 294944 B/op | 1912 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.35 ms/conn | 3432256 B/op | 21425 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.348 ms/conn | 3474176 B/op | 21765 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.585 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.801 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.92 ms/conn | 1354224 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.447 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.405 ms/conn | 1495176 B/op | 11924 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.777 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9938 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.817 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.989 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.308 ms/conn | 1547408 B/op | 12668 allocs/op |
| 会话恢复握手 | 5 | 1.171 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.163 ms/conn | 2538360 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.099 ms/conn | 2006104 B/op | 11253 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.11 ms/conn | 2018504 B/op | 11372 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.422 ms/conn | 1304192 B/op | 9259 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.32 ms/conn | 2458312 B/op | 10765 allocs/op |
| 完整 mTLS 握手 | 5 | 6.109 ms/conn | 1703472 B/op | 14581 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.131 ms/conn | 1906000 B/op | 15101 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.446 ms/conn | 1368672 B/op | 10179 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.963 ms/conn | 1026288 B/op | 7419 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.142 ms/conn | 2473032 B/op | 11327 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.195 ms/conn | 2632824 B/op | 12405 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.88 ms/conn | 3308912 B/op | 22331 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3599120 B/op | 19987 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.632 ms/conn | 1572912 B/op | 9681 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.686 ms/conn | 1592912 B/op | 9821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.023 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.975 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63712 B/op | 57 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.877 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.896 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.55 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1191240 B/op | 1204 allocs/op |
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
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.061 ms/conn | 452696 B/op | 3438 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.092 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.021 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9582 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.109 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.039 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.312 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.051 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.448 ms/conn | 537096 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.068 ms/conn | 454472 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.069 ms/conn | 487448 B/op | 3665 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.224 ms/conn | 534040 B/op | 4204 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6489 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.694 ms/pair | 868808 B/op | 7182 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.177 ms/pair | 1047432 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.711 ms/pair | 868480 B/op | 7130 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.379 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.21 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.426 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.504 ms/conn | 215664 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.342 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.954 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.497 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.193 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.501 ms/conn | 275056 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 2.349 ms/conn | 222800 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.447 ms/conn | 182032 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.387 ms/conn | 213520 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.483 ms/conn | 194608 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9326 ms/conn | 159216 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.532 ms/pair | 379536 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.357 ms/pair | 459568 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.509 ms/conn | 203632 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.791 ms/pair | 420976 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.61 ms/conn | 481192 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.711 ms/conn | 512856 B/op | 1939 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.542 ms/conn | 453224 B/op | 1925 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.038 ms/conn | 484008 B/op | 1961 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.699 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.223 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.827 ms/conn | 538456 B/op | 2039 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.489 ms/conn | 514280 B/op | 2125 allocs/op |
| 完整 mTLS 握手 | 5 | 5.079 ms/conn | 523752 B/op | 2511 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.402 ms/conn | 482728 B/op | 1975 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.586 ms/conn | 516200 B/op | 2267 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.558 ms/conn | 506472 B/op | 2151 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.744 ms/conn | 460072 B/op | 1837 allocs/op |
| 会话恢复握手 | 5 | 5.641 ms/pair | 703688 B/op | 3832 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.418 ms/pair | 802824 B/op | 4628 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.71 ms/pair | 706032 B/op | 3790 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 5.511 ms/conn | 517224 B/op | 2299 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.551 ms/pair | 740456 B/op | 4194 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.705 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.782 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.636 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.978 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.522 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 6.184 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 7.722 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 5.734 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.585 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.804 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.608 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.707 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.807 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.921 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 5.414 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 7.558 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.997 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.047 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.016 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9461 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.452 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.021 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.305 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.04 ms/conn | 507232 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.44 ms/conn | 537352 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.003 ms/conn | 454472 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.028 ms/conn | 485000 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.206 ms/conn | 533016 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8388 ms/conn | 437400 B/op | 3372 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.6412 ms/conn | 419896 B/op | 3184 allocs/op |
| 会话恢复握手 | 5 | 1.682 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.201 ms/pair | 1047176 B/op | 8531 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.388 ms/conn | 650264 B/op | 5440 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.354 ms/conn | 511576 B/op | 3988 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.153 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8145 ms/conn | 180320 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.8773 ms/conn | 213472 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8492 ms/conn | 145376 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.821 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.584 ms/conn | 185408 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.9745 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.092 ms/conn | 270816 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.8458 ms/conn | 180640 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 1.067 ms/conn | 220320 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8518 ms/conn | 180384 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.805 ms/conn | 211968 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.073 ms/conn | 228832 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7133 ms/conn | 164384 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5845 ms/conn | 160864 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.345 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.623 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.045 ms/conn | 241280 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.9334 ms/conn | 201248 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.587 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.632 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.679 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.578 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.689 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.426 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.753 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.736 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.667 ms/conn | 509864 B/op | 2045 allocs/op |
| 完整 mTLS 握手 | 5 | 2.877 ms/conn | 516552 B/op | 2341 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.657 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.574 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.687 ms/conn | 501640 B/op | 2081 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.495 ms/conn | 479688 B/op | 1957 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.275 ms/conn | 457032 B/op | 1761 allocs/op |
| 会话恢复握手 | 5 | 3.308 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.813 ms/pair | 789240 B/op | 4370 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 3.008 ms/conn | 558280 B/op | 2741 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.265 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.31 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.194 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.235 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.396 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.335 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.372 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.269 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.344 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.201 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.212 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.313 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.233 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 2.046 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.825 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.981 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.346 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.9889 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.029 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.9899 ms/conn | 390104 B/op | 3400 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.036 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.9387 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.13 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.035 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.299 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.034 ms/conn | 507104 B/op | 3815 allocs/op |
| 完整 mTLS 握手 | 5 | 1.423 ms/conn | 537480 B/op | 4556 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.016 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.022 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.205 ms/conn | 534168 B/op | 4204 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8302 ms/conn | 437400 B/op | 3372 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6271 ms/conn | 417144 B/op | 3176 allocs/op |
| 会话恢复握手 | 5 | 1.714 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.2 ms/pair | 1047176 B/op | 8531 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.704 ms/pair | 868352 B/op | 7130 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.943 ms/conn | 641480 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.155 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.188 ms/conn | 245088 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.201 ms/conn | 176224 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.094 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.2 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.612 ms/conn | 220176 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.172 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.487 ms/conn | 331760 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.18 ms/conn | 211744 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 1.524 ms/conn | 267136 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.175 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.216 ms/conn | 247424 B/op | 1988 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.5 ms/conn | 270816 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.087 ms/conn | 194272 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9758 ms/conn | 199184 B/op | 1679 allocs/op |
| 会话恢复握手 | 5 | 2.172 ms/pair | 471232 B/op | 3879 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.731 ms/pair | 590464 B/op | 4437 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.394 ms/conn | 270400 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.618 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.667 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.436 ms/conn | 453320 B/op | 1869 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.738 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.548 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.194 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.748 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.757 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.578 ms/conn | 514184 B/op | 2069 allocs/op |
| 完整 mTLS 握手 | 5 | 5.888 ms/conn | 521448 B/op | 2381 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.83 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.741 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.8 ms/conn | 506600 B/op | 2105 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.231 ms/conn | 481096 B/op | 1949 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.328 ms/conn | 459976 B/op | 1773 allocs/op |
| 会话恢复握手 | 5 | 6.833 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.157 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 6.585 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.83 ms/conn | 574344 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 5.557 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 5.621 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 5.534 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.557 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 5.5 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 9.179 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 5.754 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.038 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 5.416 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 5.816 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 5.504 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.661 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 5.867 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 6.096 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 5.389 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 6.724 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 7.336 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.956 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID 路由 / 1 个连接 / 固定长度 / 命中 | 5 | 16.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 固定长度 / 未命中 | 5 | 16.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 命中 | 5 | 17.16 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1 个连接 / 混合长度 / 未命中 | 5 | 16.52 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 命中 | 5 | 17.14 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 固定长度 / 未命中 | 5 | 12.19 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 命中 | 5 | 27.84 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 64 个连接 / 混合长度 / 未命中 | 5 | 30.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 命中 | 5 | 22.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 固定长度 / 未命中 | 5 | 13.41 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 命中 | 5 | 35.42 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 混合长度 / 未命中 | 5 | 35.7 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 命中 | 5 | 1.795 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 4–255 字节长度 / 未命中 | 5 | 3.502 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 命中 | 5 | 41.45 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 252 个连接 / 255 字节 CID / 未命中 | 5 | 18.18 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 命中 | 5 | 1.983 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 4–255 字节长度 / 未命中 | 5 | 4.063 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 命中 | 5 | 64.24 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID 路由 / 1024 个连接 / 255 字节 CID / 未命中 | 5 | 24.05 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 64 字节 | 5 | 237.5 ns/op | - | 64 B/op | 1 allocs/op |
| Listener 报文入队 / 1200 字节 | 5 | 192.7 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 2048 字节 | 5 | 202.4 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 报文入队 / 4096 字节 | 5 | 1.47 us/op | - | 4096 B/op | 1 allocs/op |
| Listener 读期限 / 固定期限 | 5 | 200.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener 读期限 / 更新期限 | 5 | 366.3 ns/op | - | 0 B/op | 0 allocs/op |
| 明文 ACK 构建 / 空 | 5 | 36.88 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 44.87 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 549.4 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.051 us/op | 1996.89 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.702 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 216.3 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 174.1 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 961.8 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.648 us/op | 1122.81 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 341.8 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 89.27 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 349.6 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 20.05 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.48 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.677 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.242 us/op | 535.19 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 884.2 ns/op | 1357.2 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 16.51 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.287 us/op | 524.71 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.301 us/op | 922.5 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.184 us/op | 167.03 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.197 us/op | 546.27 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.312 us/op | 518.94 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.062 us/op | 391.86 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 942.9 ns/op | 1272.63 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.097 us/op | 387.51 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 870.9 ns/op | 1377.91 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 928.3 ns/op | 1292.63 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.285 us/op | 934.19 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
