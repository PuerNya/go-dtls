# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- 生成时间: `2026-10-13T15:16:37Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `83bca78b7b545eed982e6775b73750064408847e (Linux Release static)`
- NSS: `b9283c1e4ee232b4450d36115272af670d38ec81 (Linux Release)`
- BoringSSL: `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 (Linux Release)`
- OpenSSL: `41eecc7c7cd4069b1b29bca49b9d97825c6bf648 (Linux Release)`

共 347 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (34)](#section-connection-lifecycle)
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
- [记录层与可靠性 (29)](#section-record-layer-and-reliability)

<a id="section-connection-lifecycle"></a>
## 连接生命周期

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 564.502 us/op | - | 101322 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 857.481 us/op | - | 116335 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 440.226 us/op | - | 123191 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.09 ms/op | - | 123643 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.458 ms/op | - | 142421 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 655.454 us/op | - | 121348 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 897.342 us/op | - | 125649 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 899.638 us/op | - | 125745 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 366.209 us/op | - | 103774 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.007 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.013 ms/op | - | 126045 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.691 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.689 ms/op | - | 162384 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 986.173 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 991.95 us/op | - | 154179 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 845.072 us/op | - | 149311 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 800.527 us/op | - | 152656 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.225 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 13.345 us/op | 9.59 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 16.075 us/op | 149.3 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 13.371 us/op | 9.57 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 16.113 us/op | 148.95 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 14.445 us/op | 8.86 MB/s | 596 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 17.285 us/op | 138.85 MB/s | 2996 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 15.233 us/op | 8.4 MB/s | 593 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 18.047 us/op | 132.99 MB/s | 2998 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 18.966 us/op | 6.75 MB/s | 593 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 22.535 us/op | 106.5 MB/s | 3004 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 26.53 us/op | 4.82 MB/s | 619 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 29.538 us/op | 81.25 MB/s | 3027 B/op | 12 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 409.543 us/op | - | 109279 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 407.796 us/op | - | 109116 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 3.307 us/op | 38.71 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 6.618 us/op | 362.66 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.892 ms/conn | 2671704 B/op | 21758 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.89 ms/conn | 2682536 B/op | 22097 allocs/op |
| 完整 mTLS 握手 | 5 | 3.435 ms/conn | 3375328 B/op | 30103 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.481 ms/conn | 3832400 B/op | 31826 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.912 ms/conn | 3149576 B/op | 25608 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4499 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.885 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.001 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.976 ms/conn | 4180208 B/op | 40352 allocs/op |
| 会话恢复握手 | 5 | 0.5367 ms/conn | 4735608 B/op | 39433 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6501 ms/conn | 6867720 B/op | 51433 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.5498 ms/conn | 298872 B/op | 1959 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.224 ms/conn | 3543168 B/op | 22081 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.217 ms/conn | 3585088 B/op | 22421 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.562 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.065 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.185 ms/conn | 1353584 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.934 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.927 ms/conn | 1495024 B/op | 11922 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.981 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9988 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.021 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.207 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.78 ms/conn | 1546960 B/op | 12671 allocs/op |
| 会话恢复握手 | 5 | 1.049 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.063 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.431 ms/conn | 2005320 B/op | 11245 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.567 ms/conn | 2017720 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.226 ms/conn | 1326784 B/op | 9519 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.709 ms/conn | 2480824 B/op | 11025 allocs/op |
| 完整 mTLS 握手 | 5 | 5.529 ms/conn | 1770936 B/op | 15066 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.519 ms/conn | 1973680 B/op | 15589 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.197 ms/conn | 1391680 B/op | 10439 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.79 ms/conn | 1048848 B/op | 7679 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.704 ms/conn | 2494552 B/op | 11586 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.759 ms/conn | 2655632 B/op | 12665 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.55 ms/conn | 3403664 B/op | 22837 allocs/op |
| 会话恢复握手 | 5 | 1010 ms/pair | 3642320 B/op | 20507 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1010 ms/pair | 244432 B/op | 1014 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.381 ms/conn | 1628128 B/op | 10019 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.546 ms/conn | 1648448 B/op | 10159 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.491 ms/conn | 1807040 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.392 ms/conn | 63664 B/op | 56 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.086 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 7.674 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 7.222 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.38 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.165 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.294 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.031 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 10.19 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.377 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.656 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 11.31 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8956 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9477 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8863 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.8277 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.144 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.9063 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.106 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.8916 ms/conn | 517120 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 1.256 ms/conn | 549480 B/op | 4699 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8775 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.8797 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.108 ms/conn | 540920 B/op | 4304 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.532 ms/conn | 423192 B/op | 3252 allocs/op |
| 会话恢复握手 | 5 | 1.4 ms/pair | 884792 B/op | 7357 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.797 ms/pair | 1065352 B/op | 8755 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.452 ms/pair | 882184 B/op | 7303 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.258 ms/conn | 519608 B/op | 4088 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.883 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.275 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.284 ms/conn | 215280 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.206 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.697 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.321 ms/conn | 186512 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.852 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.09 ms/conn | 274736 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 1.965 ms/conn | 222800 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.175 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.186 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.182 ms/conn | 194288 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7842 ms/conn | 158896 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 2.068 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.942 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 1.305 ms/conn | 203312 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 2.415 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 3.485 ms/conn | 485672 B/op | 1995 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 3.602 ms/conn | 517288 B/op | 1994 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 3.491 ms/conn | 457576 B/op | 1974 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.05 ms/conn | 488488 B/op | 2018 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 6.602 ms/conn | 490856 B/op | 2053 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.104 ms/conn | 518600 B/op | 2191 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.431 ms/conn | 555944 B/op | 2121 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 3.484 ms/conn | 518696 B/op | 2179 allocs/op |
| 完整 mTLS 握手 | 5 | 4.057 ms/conn | 528488 B/op | 2565 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 3.564 ms/conn | 487080 B/op | 2025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.595 ms/conn | 520584 B/op | 2319 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.578 ms/conn | 510376 B/op | 2205 allocs/op |
| 直接外部 PSK 握手 | 5 | 2.711 ms/conn | 464520 B/op | 1893 allocs/op |
| 会话恢复握手 | 5 | 4.372 ms/pair | 712392 B/op | 3932 allocs/op |
| mTLS 会话恢复握手 | 5 | 5.008 ms/pair | 811592 B/op | 4734 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 4.506 ms/pair | 714272 B/op | 3886 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 4.415 ms/conn | 521640 B/op | 2353 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 5.136 ms/pair | 749224 B/op | 4298 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 3.69 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 3.81 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 3.743 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.864 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.582 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.996 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 6.394 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 4.663 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 3.69 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.755 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.722 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 2.842 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 5.529 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.57 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 4.326 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 6.237 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8735 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9598 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8973 ms/conn | 398008 B/op | 3500 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.8157 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.132 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.9041 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.123 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.9108 ms/conn | 517248 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 1.271 ms/conn | 549848 B/op | 4709 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8829 ms/conn | 462504 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.8924 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.069 ms/conn | 542200 B/op | 4304 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7005 ms/conn | 445432 B/op | 3472 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.5323 ms/conn | 425944 B/op | 3260 allocs/op |
| 会话恢复握手 | 5 | 1.405 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.83 ms/pair | 1065480 B/op | 8755 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.24 ms/conn | 658296 B/op | 5540 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 1.236 ms/conn | 519480 B/op | 4088 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.884 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.6671 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.7261 ms/conn | 213024 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.6737 ms/conn | 144928 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.642 ms/conn | 182880 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.399 ms/conn | 184960 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.8389 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.9577 ms/conn | 270496 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.691 ms/conn | 180320 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 1.005 ms/conn | 220000 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.7751 ms/conn | 180064 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.6773 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.9643 ms/conn | 229152 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.5814 ms/conn | 164064 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.4908 ms/conn | 160544 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 1.155 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.475 ms/pair | 473568 B/op | 3655 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 0.9356 ms/conn | 240960 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.8434 ms/conn | 200800 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.481 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.201 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.225 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.084 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 2.21 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.746 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 2.107 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.301 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.189 ms/conn | 514184 B/op | 2097 allocs/op |
| 完整 mTLS 握手 | 5 | 2.478 ms/conn | 521064 B/op | 2393 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.219 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.187 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 2.249 ms/conn | 506536 B/op | 2133 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 2.123 ms/conn | 484072 B/op | 2009 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 1.9 ms/conn | 461416 B/op | 1813 allocs/op |
| 会话恢复握手 | 5 | 2.767 ms/pair | 701816 B/op | 3762 allocs/op |
| mTLS 会话恢复握手 | 5 | 3.153 ms/pair | 798152 B/op | 4474 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.49 ms/conn | 562664 B/op | 2793 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.892 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.927 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.882 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.889 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.823 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.991 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 2.036 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.862 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 2.067 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.868 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.838 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.926 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.883 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 1.773 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 2.341 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.551 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 2.028 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8888 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.9584 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.8709 ms/conn | 398136 B/op | 3500 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 0.8861 ms/conn | 525208 B/op | 4223 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.8067 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.111 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.8833 ms/conn | 535520 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.127 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.8937 ms/conn | 517248 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 1.267 ms/conn | 549096 B/op | 4704 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.8758 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.898 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.057 ms/conn | 542200 B/op | 4304 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.7233 ms/conn | 445432 B/op | 3472 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.51 ms/conn | 423192 B/op | 3252 allocs/op |
| 会话恢复握手 | 5 | 1.403 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.872 ms/pair | 1064976 B/op | 8752 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.464 ms/pair | 881936 B/op | 7300 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.718 ms/conn | 653728 B/op | 6121 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.017 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.006 ms/conn | 244768 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.9631 ms/conn | 175904 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 0.9874 ms/conn | 264352 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.003 ms/conn | 220064 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 3.424 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.056 ms/conn | 223392 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.289 ms/conn | 331552 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.064 ms/conn | 211552 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 1.291 ms/conn | 266560 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.9569 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.9611 ms/conn | 246912 B/op | 1986 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.151 ms/conn | 270496 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.8587 ms/conn | 193952 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7365 ms/conn | 199008 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 1.545 ms/pair | 470496 B/op | 3878 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.282 ms/pair | 592400 B/op | 4443 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.089 ms/conn | 270576 B/op | 2310 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.665 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.675 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.471 ms/conn | 457704 B/op | 1921 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.57 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.561 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 7.049 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.671 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 4.886 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.693 ms/conn | 518632 B/op | 2121 allocs/op |
| 完整 mTLS 握手 | 5 | 4.865 ms/conn | 525704 B/op | 2433 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.567 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.658 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.815 ms/conn | 510984 B/op | 2157 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.122 ms/conn | 485480 B/op | 2001 allocs/op |
| 直接外部 PSK 握手 | 5 | 4.34 ms/conn | 464360 B/op | 1825 allocs/op |
| 会话恢复握手 | 5 | 5.545 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS 会话恢复握手 | 5 | 5.941 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 5.484 ms/pair | 711744 B/op | 3784 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.425 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 4.615 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 4.644 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 4.637 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.653 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 4.75 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 8.159 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.879 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 5.097 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 4.618 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 4.884 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 4.634 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.681 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 4.766 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 5.138 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 4.492 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 5.639 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 6.112 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.953 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 35.95 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.34 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 519.4 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.983 us/op | 2065.19 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.667 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 200.7 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 160.9 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 916.1 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.369 us/op | 1215.88 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 324.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 89.62 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 346.8 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 28.01 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 24.97 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.334 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.358 us/op | 508.9 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 840.4 ns/op | 1427.86 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 12.21 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.441 us/op | 491.61 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.247 us/op | 961.97 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 6.894 us/op | 174.06 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.337 us/op | 513.46 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.43 us/op | 493.79 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.535 us/op | 339.46 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 937.7 ns/op | 1279.73 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 2.928 us/op | 409.9 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 834.2 ns/op | 1438.56 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 902.8 ns/op | 1329.23 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.478 us/op | 811.71 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
