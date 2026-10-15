# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- 生成时间: `2026-10-15T13:25:54Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V45 96-Core Processor`
- wolfSSL: `f46cdde15bf815f58d5f93dbdd155d342cba0a22 (Linux Release static)`
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
| 证书认证完整握手 / AES-128-GCM | 5 | 304.837 us/op | - | 101322 B/op | 732 allocs/op |
| 完整 mTLS 握手 | 5 | 461.05 us/op | - | 116339 B/op | 929 allocs/op |
| mTLS 会话恢复握手 | 5 | 244.357 us/op | - | 123190 B/op | 855 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 572.253 us/op | - | 123643 B/op | 1108 allocs/op |
| 握手后认证的多证书选择 | 5 | 813.331 us/op | - | 142410 B/op | 1404 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 374.611 us/op | - | 121350 B/op | 962 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 508.787 us/op | - | 125653 B/op | 1000 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 499.318 us/op | - | 125744 B/op | 1000 allocs/op |
| 直接外部 PSK 握手 | 5 | 200.703 us/op | - | 103777 B/op | 758 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 544.771 us/op | - | 136088 B/op | 1017 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 545.455 us/op | - | 126038 B/op | 993 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 925.865 us/op | - | 177150 B/op | 1469 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 925.067 us/op | - | 162375 B/op | 1422 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 534.361 us/op | - | 151409 B/op | 1232 allocs/op |
| ECH 握手 / 经 HRR | 5 | 527.437 us/op | - | 154186 B/op | 1253 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 476.526 us/op | - | 149307 B/op | 765 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 454.236 us/op | - | 152651 B/op | 796 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1.199 ms/op | - | 178093 B/op | 814 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 64 字节 | 5 | 15.227 us/op | 8.41 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / 无 CID / 1200 字节 | 5 | 17.056 us/op | 140.71 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 64 字节 | 5 | 15.238 us/op | 8.4 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP 往返 / 1 个连接 / CID / 1200 字节 | 5 | 17.094 us/op | 140.4 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 64 字节 | 5 | 15.475 us/op | 8.27 MB/s | 595 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / 无 CID / 1200 字节 | 5 | 17.178 us/op | 139.72 MB/s | 2996 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 64 字节 | 5 | 15.848 us/op | 8.08 MB/s | 597 B/op | 12 allocs/op |
| Listener UDP 往返 / 64 个连接 / CID / 1200 字节 | 5 | 17.661 us/op | 135.89 MB/s | 2998 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 64 字节 | 5 | 20.991 us/op | 6.1 MB/s | 601 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / 无 CID / 1200 字节 | 5 | 23.556 us/op | 101.88 MB/s | 3009 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 64 字节 | 5 | 26.499 us/op | 4.83 MB/s | 619 B/op | 12 allocs/op |
| Listener UDP 往返 / 1024 个连接 / CID / 1200 字节 | 5 | 28.683 us/op | 83.67 MB/s | 3024 B/op | 12 allocs/op |
| 0-RTT / 初始 64 条重放记录 | 5 | 236.834 us/op | - | 109289 B/op | 796 allocs/op |
| 0-RTT / 初始 4096 条重放记录 | 5 | 242.809 us/op | - | 109121 B/op | 796 allocs/op |
| 内存 datagram 往返 / 64 字节 | 5 | 2.03 us/op | 63.04 MB/s | 384 B/op | 4 allocs/op |
| 内存 datagram 往返 / 1200 字节 | 5 | 3.856 us/op | 622.36 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.128 ms/conn | 2671752 B/op | 21759 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.13 ms/conn | 2682584 B/op | 22098 allocs/op |
| 完整 mTLS 握手 | 5 | 2.061 ms/conn | 3374256 B/op | 30089 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 2.054 ms/conn | 3833624 B/op | 31841 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.131 ms/conn | 3149576 B/op | 25608 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.2811 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.197 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.319 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 3.151 ms/conn | 4180504 B/op | 40354 allocs/op |
| 会话恢复握手 | 5 | 0.3818 ms/conn | 4735608 B/op | 39420 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.4569 ms/conn | 6868216 B/op | 51436 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.4057 ms/conn | 299120 B/op | 1961 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.386 ms/conn | 3543168 B/op | 22081 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1.437 ms/conn | 3585088 B/op | 22421 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.223 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.81 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.874 ms/conn | 1353584 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 3.664 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.421 ms/conn | 1495024 B/op | 11922 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.65 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.6755 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.916 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.025 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 3.682 ms/conn | 1546576 B/op | 12665 allocs/op |
| 会话恢复握手 | 5 | 0.7997 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.4846 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.789 ms/conn | 2005240 B/op | 11244 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.47 ms/conn | 2017800 B/op | 11365 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.324 ms/conn | 1326976 B/op | 9519 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 3.031 ms/conn | 2480248 B/op | 11025 allocs/op |
| 完整 mTLS 握手 | 5 | 3.322 ms/conn | 1771584 B/op | 15069 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.33 ms/conn | 1973480 B/op | 15586 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.344 ms/conn | 1391232 B/op | 10439 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.539 ms/conn | 1048624 B/op | 7679 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 3.295 ms/conn | 2494456 B/op | 11586 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.325 ms/conn | 2656984 B/op | 12665 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 7.683 ms/conn | 3403920 B/op | 22840 allocs/op |
| 会话恢复握手 | 5 | 1006 ms/pair | 3641920 B/op | 20507 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1006 ms/pair | 244504 B/op | 1014 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.535 ms/conn | 1627968 B/op | 10019 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.057 ms/conn | 1647968 B/op | 10159 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 4.06 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.494 ms/conn | 1169792 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 1024 ms/conn | 63680 B/op | 57 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.875 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.897 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.955 ms/conn | 1172800 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1008 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS 会话恢复握手 | 5 | 1009 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.58 ms/conn | 460600 B/op | 3538 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.6069 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.598 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.5493 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 1.237 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.5828 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.7637 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.5849 ms/conn | 517248 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 0.8672 ms/conn | 549024 B/op | 4702 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.6218 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.6299 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.7255 ms/conn | 541048 B/op | 4304 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.3723 ms/conn | 423192 B/op | 3252 allocs/op |
| 会话恢复握手 | 5 | 0.9956 ms/pair | 885336 B/op | 7362 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.262 ms/pair | 1065728 B/op | 8758 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.022 ms/pair | 882064 B/op | 7299 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.8102 ms/conn | 519608 B/op | 4088 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.289 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls 客户端 -> NSS 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.8598 ms/conn | 181424 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.7978 ms/conn | 215280 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.7303 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.119 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.587 ms/conn | 186640 B/op | 1594 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.207 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.546 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制 NSS 服务端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 1.327 ms/conn | 222672 B/op | 1910 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.76 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.7646 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.7566 ms/conn | 194288 B/op | 1638 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5357 ms/conn | 158896 B/op | 1306 allocs/op |
| 会话恢复握手 | 5 | 1.384 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.968 ms/pair | 458928 B/op | 3517 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 0.8664 ms/conn | 203312 B/op | 1714 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.598 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 2.969 ms/conn | 485608 B/op | 1990 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 2.956 ms/conn | 517192 B/op | 1989 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 2.941 ms/conn | 457544 B/op | 1973 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 3.29 ms/conn | 488424 B/op | 2015 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.738 ms/conn | 490840 B/op | 2051 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 3.361 ms/conn | 518536 B/op | 2187 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 4.287 ms/conn | 555912 B/op | 2119 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 2.856 ms/conn | 518600 B/op | 2173 allocs/op |
| 完整 mTLS 握手 | 5 | 3.366 ms/conn | 528136 B/op | 2563 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 2.976 ms/conn | 487080 B/op | 2025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.017 ms/conn | 520616 B/op | 2321 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.062 ms/conn | 510376 B/op | 2205 allocs/op |
| 直接外部 PSK 握手 | 5 | 2.417 ms/conn | 464456 B/op | 1889 allocs/op |
| 会话恢复握手 | 5 | 3.609 ms/pair | 712296 B/op | 3926 allocs/op |
| mTLS 会话恢复握手 | 5 | 4.085 ms/pair | 811592 B/op | 4732 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 3.812 ms/pair | 714568 B/op | 3882 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 3.654 ms/conn | 521576 B/op | 2349 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 4.228 ms/pair | 749128 B/op | 4292 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS 客户端 -> NSS 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 3.088 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 3.136 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 3.088 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 3.886 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 6.041 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 4.005 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 4.838 ms/conn | 212448 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | - | 跳过: NSS 夹具无法强制两个 NSS 对端的握手分片；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 完整 mTLS 握手 | 5 | 3.864 ms/conn | 212192 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 3.147 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.191 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.2 ms/conn | 212448 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 2.518 ms/conn | 212448 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 4.399 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS 会话恢复握手 | 5 | 5.237 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 此夹具中的 NSS 服务端无法接收 DTLS 1.3 早期数据；该限制最后验证于 NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| 委派凭证握手 + 数据报交换 | 5 | 3.7 ms/conn | 212448 B/op | 176 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 5.008 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.6069 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.6317 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.6236 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.5807 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 1.295 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.5837 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.7781 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.62 ms/conn | 517248 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 0.8418 ms/conn | 549352 B/op | 4703 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.6036 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.6182 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.6973 ms/conn | 541048 B/op | 4304 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.5014 ms/conn | 445432 B/op | 3472 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.3718 ms/conn | 425944 B/op | 3260 allocs/op |
| 会话恢复握手 | 5 | 0.9864 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.22 ms/pair | 1065224 B/op | 8755 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 0.8056 ms/conn | 658296 B/op | 5540 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.8255 ms/conn | 519608 B/op | 4088 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 1.217 ms/pair | 961848 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls 客户端 -> BoringSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.4375 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.4727 ms/conn | 213024 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.4309 ms/conn | 145056 B/op | 1503 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.4098 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 1.398 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.4628 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.5918 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.4474 ms/conn | 180320 B/op | 1556 allocs/op |
| 完整 mTLS 握手 | 5 | 0.5695 ms/conn | 220000 B/op | 1900 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.4731 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.4778 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.5557 ms/conn | 228512 B/op | 2068 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.3997 ms/conn | 164064 B/op | 1364 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 0.2893 ms/conn | 160544 B/op | 1308 allocs/op |
| 会话恢复握手 | 5 | 0.7368 ms/pair | 388144 B/op | 3280 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8547 ms/pair | 473344 B/op | 3654 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 0.5364 ms/conn | 240960 B/op | 2252 allocs/op |
| 委派凭证握手 + 数据报交换 | 5 | 0.474 ms/conn | 200928 B/op | 1704 allocs/op |
| 委派凭证握手 + 会话恢复 | 5 | 0.7874 ms/pair | 431488 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.814 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.836 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.808 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.805 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.734 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.747 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.921 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.768 ms/conn | 514248 B/op | 2097 allocs/op |
| 完整 mTLS 握手 | 5 | 1.935 ms/conn | 521064 B/op | 2393 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.778 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.774 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.837 ms/conn | 506024 B/op | 2133 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.718 ms/conn | 484072 B/op | 2009 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 1.631 ms/conn | 461416 B/op | 1813 allocs/op |
| 会话恢复握手 | 5 | 2.12 ms/pair | 701816 B/op | 3762 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.301 ms/pair | 798008 B/op | 4474 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.999 ms/conn | 562664 B/op | 2793 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL 客户端 -> BoringSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 1.614 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 1.634 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 1.605 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 1.6 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 2.728 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 1.628 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 1.678 ms/conn | 212768 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 1.556 ms/conn | 212768 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 1.68 ms/conn | 212512 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 1.582 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.563 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 1.626 ms/conn | 212512 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 1.59 ms/conn | 212512 B/op | 176 allocs/op |
| 导入式外部 PSK 握手 + 数据报交换 | 5 | 1.498 ms/conn | 212512 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 1.831 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 2.001 ms/pair | 210848 B/op | 192 allocs/op |
| ECH 接受握手 + 数据报交换 | 5 | 1.724 ms/conn | 212256 B/op | 176 allocs/op |
| 委派凭证握手 + 数据报交换 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| 委派凭证握手 + 会话恢复 | - | 跳过: BoringSSL 客户端不认证 DTLS 1.3 委派凭证；该限制最后验证于 BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.5591 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.592 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.6119 ms/conn | 398008 B/op | 3500 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 0.636 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.5606 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 1.275 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.6199 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.7786 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.6416 ms/conn | 517248 B/op | 3939 allocs/op |
| 完整 mTLS 握手 | 5 | 0.8727 ms/conn | 549616 B/op | 4704 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.6248 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.6337 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.7231 ms/conn | 542072 B/op | 4304 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.5223 ms/conn | 445432 B/op | 3472 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.3787 ms/conn | 423192 B/op | 3252 allocs/op |
| 会话恢复握手 | 5 | 0.9297 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.285 ms/pair | 1065592 B/op | 8760 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1.08 ms/pair | 882064 B/op | 7299 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 1.139 ms/conn | 653728 B/op | 6121 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls 客户端 -> OpenSSL 服务端

中位耗时由 go-dtls 客户端计时并包含认证数据报交换；`ms/pair` 包含获取票据的完整连接。OpenSSL 服务端方向有意丢弃客户端的 epoch-2 ACK。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 0.5557 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 0.5819 ms/conn | 244768 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 0.6165 ms/conn | 175776 B/op | 1755 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 0.5875 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 0.6215 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 1.958 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 0.6546 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 0.8311 ms/conn | 331552 B/op | 1880 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 0.6131 ms/conn | 211552 B/op | 1788 allocs/op |
| 完整 mTLS 握手 | 5 | 0.8386 ms/conn | 266816 B/op | 2240 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 0.6244 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 0.5944 ms/conn | 247008 B/op | 1987 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 0.7873 ms/conn | 270368 B/op | 2416 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 0.533 ms/conn | 193952 B/op | 1600 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4222 ms/conn | 199008 B/op | 1680 allocs/op |
| 会话恢复握手 | 5 | 0.9922 ms/pair | 470496 B/op | 3878 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.497 ms/pair | 589888 B/op | 4436 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 0.7146 ms/conn | 270328 B/op | 2307 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL 客户端 -> go-dtls 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 3.374 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 3.433 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 3.621 ms/conn | 457704 B/op | 1921 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 3.623 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 3.613 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 4.951 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 3.717 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 3.893 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 3.705 ms/conn | 518632 B/op | 2121 allocs/op |
| 完整 mTLS 握手 | 5 | 3.93 ms/conn | 525832 B/op | 2433 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 3.682 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.75 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.788 ms/conn | 510984 B/op | 2157 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 4.019 ms/conn | 485480 B/op | 2001 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.162 ms/conn | 464360 B/op | 1825 allocs/op |
| 会话恢复握手 | 5 | 3.785 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS 会话恢复握手 | 5 | 4.515 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 4.225 ms/pair | 711744 B/op | 3784 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.134 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL 客户端 -> OpenSSL 服务端

中位耗时包含对端客户端进程启动和认证数据报交换；`ms/pair` 包含完整连接及恢复连接。进程耗时不能直接与仅客户端计时比较。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证 AES-128-GCM 连接 + 数据报交换 | 5 | 3.374 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 连接 + 数据报交换 | 5 | 3.675 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 连接 + 数据报交换 | 5 | 3.627 ms/conn | 212896 B/op | 176 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 3.687 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 密钥交换 + 数据报交换 | 5 | 3.673 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 密钥交换 + 数据报交换 | 5 | 5.67 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest 握手 + 数据报交换 | 5 | 3.773 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 混合握手 + 数据报交换 | 5 | 3.953 ms/conn | 212576 B/op | 176 allocs/op |
| MTU 256 分片握手 + 数据报交换 | 5 | 3.738 ms/conn | 212832 B/op | 176 allocs/op |
| 完整 mTLS 握手 | 5 | 3.919 ms/conn | 212576 B/op | 176 allocs/op |
| 协商 ALPN 的连接 + 数据报交换 | 5 | 3.713 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.738 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP 装订握手 + 数据报交换 | 5 | 3.788 ms/conn | 212576 B/op | 176 allocs/op |
| 原始公钥握手 + 数据报交换 | 5 | 4.055 ms/conn | 212576 B/op | 176 allocs/op |
| 直接外部 PSK 握手 | 5 | 3.246 ms/conn | 212576 B/op | 176 allocs/op |
| 会话恢复握手 | 5 | 4.215 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS 会话恢复握手 | 5 | 4.722 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: OpenSSL 监听器无法进入 DTLS 早期数据读取状态；该限制最后验证于 OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 应用数据 1-RTT 往返 | 5 | 3.988 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 20.35 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 24.11 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 262.5 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.003 us/op | 4082.85 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 888.2 ns/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 124.3 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 109.2 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 503.7 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 1.928 us/op | 2124.77 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 165.2 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 40.02 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 178.4 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 16.45 ns/op | - | 0 B/op | 0 allocs/op |
| 解析 ACK / 独占 | 5 | 14.02 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 2.197 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 1.184 us/op | 1013.38 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 478.7 ns/op | 2506.75 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 7.89 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 1.201 us/op | 999.28 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 732.5 ns/op | 1638.34 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 4.246 us/op | 282.63 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 1.154 us/op | 1039.48 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 1.213 us/op | 989.1 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 2.061 us/op | 582.3 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 557.4 ns/op | 2152.86 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 1.887 us/op | 636.07 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 479.4 ns/op | 2503.36 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 505.8 ns/op | 2372.37 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 859.4 ns/op | 1396.26 MB/s | 1280 B/op | 1 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
