# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `8010fee127101d213fb208cbfbff51779d6c6954`
- 生成时间: `2026-10-06T03:54:10Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `b874f7167fc592b331c79a33674bc51c2f0286b0 (Linux Release static)`

共 181 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (18)](#section-connection-lifecycle)
- [真实 UDP 互通 (60)](#section-real-udp-interoperability)
  - [go-dtls 客户端 -> go-dtls 服务端 (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls 客户端 -> wolfSSL 服务端 (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL 客户端 -> go-dtls 服务端 (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL 客户端 -> wolfSSL 服务端 (15)](#real-udp-wolfssl-client-wolfssl-server)
- [记录层与可靠性 (37)](#section-record-layer-and-reliability)
- [密钥调度与密码学 (38)](#section-key-schedule-and-cryptography)
- [报文编码与解析 (26)](#section-wire-encoding-and-parsing)
- [证书压缩 (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## 连接生命周期

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 394.977 us/op | 97817 B/op | 709 allocs/op |
| 完整 mTLS 握手 | 5 | 607.17 us/op | 111975 B/op | 891 allocs/op |
| mTLS 会话恢复握手 | 5 | 322.684 us/op | 119186 B/op | 822 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 798.085 us/op | 119589 B/op | 1064 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.07 ms/op | 138634 B/op | 1358 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 472.564 us/op | 117039 B/op | 932 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 652.52 us/op | 121225 B/op | 963 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 652.658 us/op | 121225 B/op | 963 allocs/op |
| 直接外部 PSK 握手 | 5 | 271.591 us/op | 101264 B/op | 743 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 711.875 us/op | 129632 B/op | 992 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 710.437 us/op | 121751 B/op | 970 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.196 ms/op | 169029 B/op | 1425 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.191 ms/op | 156479 B/op | 1384 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 703.72 us/op | 147133 B/op | 1209 allocs/op |
| ECH 握手 / 经 HRR | 5 | 721.386 us/op | 149910 B/op | 1230 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 612.875 us/op | 145793 B/op | 742 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 595.158 us/op | 149138 B/op | 773 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1.537 ms/op | 174579 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.226 ms/conn | 2482048 B/op | 21235 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.254 ms/conn | 2490720 B/op | 21555 allocs/op |
| 完整 mTLS 握手 | 5 | 2.173 ms/conn | 3163456 B/op | 29333 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.369 ms/conn | 3658944 B/op | 31346 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.235 ms/conn | 2971168 B/op | 25166 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.2951 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.235 ms/conn | 2502016 B/op | 22353 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.337 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 3.251 ms/conn | 3956952 B/op | 39455 allocs/op |
| 会话恢复握手 | 5 | 0.3481 ms/conn | 4432304 B/op | 37841 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.4032 ms/conn | 6449872 B/op | 49583 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.3905 ms/conn | 286952 B/op | 1893 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.482 ms/conn | 3353512 B/op | 21558 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1.477 ms/conn | 3395432 B/op | 21898 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.425 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 3.099 ms/conn | 1186064 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 3.144 ms/conn | 1204784 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 4.071 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.044 ms/conn | 1341744 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 3.086 ms/conn | 1424464 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.624 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 3.074 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.198 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.008 ms/conn | 1396816 B/op | 12385 allocs/op |
| 会话恢复握手 | 5 | 0.8169 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8138 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 3.234 ms/conn | 1857800 B/op | 10965 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.985 ms/conn | 1870200 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.53 ms/conn | 1305808 B/op | 9220 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 3.476 ms/conn | 2406104 B/op | 10725 allocs/op |
| 完整 mTLS 握手 | 5 | 3.938 ms/conn | 1742840 B/op | 14706 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.941 ms/conn | 1946448 B/op | 15229 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.419 ms/conn | 1370720 B/op | 10139 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.516 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 3.457 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.522 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 8.498 ms/conn | 3314640 B/op | 22483 allocs/op |
| 会话恢复握手 | 5 | 1007 ms/pair | 3492336 B/op | 19887 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1007 ms/pair | 238256 B/op | 990 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.596 ms/conn | 1555648 B/op | 9719 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.471 ms/conn | 1575808 B/op | 9859 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 4.727 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.956 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.197 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 5.22 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.279 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.849 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.684 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.24 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.352 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 7.796 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 1190848 B/op | 1202 allocs/op |
| mTLS 会话恢复握手 | 5 | 1011 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.918 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 4.652 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.703 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 29.95 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 37.73 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 524 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.934 us/op | 2117.67 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.423 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 157.1 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 131.7 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 788.7 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.019 us/op | 1356.74 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 299.4 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 105 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 301.8 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 62.38 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 24.97 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 41.81 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 23.35 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 459.7 ns/op | 2610.42 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 433.8 ns/op | 2766.3 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 436 ns/op | 2752.08 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 22.06 us/op | 2970.83 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 405.4 ns/op | 2960.25 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 21.43 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 3.494 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.038 us/op | 588.73 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 767.8 ns/op | 1562.9 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 9.2 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.232 us/op | 537.66 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.07 us/op | 1121.42 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 5.91 us/op | 203.03 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.026 us/op | 592.29 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.162 us/op | 555.07 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.086 us/op | 388.82 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 868 ns/op | 1382.56 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 2.466 us/op | 486.69 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 782.4 ns/op | 1533.82 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 836.8 ns/op | 1434.09 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.307 us/op | 918.14 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.256 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 5.501 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.261 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 2.959 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.229 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 2.926 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 0.9168 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.139 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 2.783 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 5.268 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 8.756 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 5.788 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 13.935 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 559.6 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.499 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 16.87 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 588.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 589.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.379 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 3.555 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 16.86 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.423 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.421 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 1.92 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.204 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 3.928 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.509 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.142 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 5.689 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.002 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 5.554 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 243.8 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 485.4 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 97.96 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 70.34 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 262.1 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 225.1 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 277.5 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 475.7 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 43.11 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 422.2 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 72.36 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 36.04 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 59.18 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 601.7 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 73.34 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 55.9 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 52.77 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 512 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 327.5 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 10.37 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 45.38 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 73.87 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 45.51 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 22.68 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 240.5 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 128.7 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 62.62 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 905.3 ns/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 647.5 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 647.6 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 9.252 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 45.18 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 4.173 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 5.129 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
