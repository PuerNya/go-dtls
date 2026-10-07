# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- 生成时间: `2026-10-07T15:17:54Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `e09a10ecd5072c6138504b77fc710607cdcc6771 (Linux Release static)`

共 178 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (18)](#section-connection-lifecycle)
- [真实 UDP 互通 (60)](#section-real-udp-interoperability)
  - [go-dtls 客户端 -> go-dtls 服务端 (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls 客户端 -> wolfSSL 服务端 (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL 客户端 -> go-dtls 服务端 (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL 客户端 -> wolfSSL 服务端 (15)](#real-udp-wolfssl-client-wolfssl-server)
- [记录层与可靠性 (34)](#section-record-layer-and-reliability)
- [密钥调度与密码学 (38)](#section-key-schedule-and-cryptography)
- [报文编码与解析 (26)](#section-wire-encoding-and-parsing)
- [证书压缩 (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## 连接生命周期

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 545.622 us/op | 97993 B/op | 710 allocs/op |
| 完整 mTLS 握手 | 5 | 835.169 us/op | 112213 B/op | 894 allocs/op |
| mTLS 会话恢复握手 | 5 | 418.255 us/op | 119463 B/op | 823 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.06 ms/op | 120009 B/op | 1068 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.428 ms/op | 138880 B/op | 1361 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 639.841 us/op | 117569 B/op | 934 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 876.374 us/op | 121529 B/op | 966 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 879.722 us/op | 121529 B/op | 966 allocs/op |
| 直接外部 PSK 握手 | 5 | 354.386 us/op | 101438 B/op | 744 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.04 ms/op | 130098 B/op | 993 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.031 ms/op | 122315 B/op | 971 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.72 ms/op | 169844 B/op | 1428 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.71 ms/op | 157394 B/op | 1387 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 951.206 us/op | 147499 B/op | 1210 allocs/op |
| ECH 握手 / 经 HRR | 5 | 964.75 us/op | 150276 B/op | 1231 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 875.034 us/op | 145972 B/op | 743 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 840.046 us/op | 149316 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.239 ms/op | 174757 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.823 ms/conn | 2487488 B/op | 21255 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.837 ms/conn | 2496160 B/op | 21575 allocs/op |
| 完整 mTLS 握手 | 5 | 3.26 ms/conn | 3170176 B/op | 29391 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.508 ms/conn | 3667392 B/op | 31404 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.859 ms/conn | 2976560 B/op | 25185 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4213 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.823 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.939 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.895 ms/conn | 3965008 B/op | 39501 allocs/op |
| 会话恢复握手 | 5 | 0.5073 ms/conn | 4443392 B/op | 37893 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.5459 ms/conn | 6462392 B/op | 49655 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.4801 ms/conn | 287904 B/op | 1898 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.156 ms/conn | 3358952 B/op | 21578 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.144 ms/conn | 3400872 B/op | 21918 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.452 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.852 ms/conn | 1187344 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.871 ms/conn | 1207344 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 6.638 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.579 ms/conn | 1343344 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.746 ms/conn | 1425744 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9284 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.739 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.839 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.525 ms/conn | 1398928 B/op | 12383 allocs/op |
| 会话恢复握手 | 5 | 0.9656 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.9488 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.114 ms/conn | 1859080 B/op | 10965 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.251 ms/conn | 1871480 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.192 ms/conn | 1291136 B/op | 9299 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.611 ms/conn | 2444504 B/op | 10805 allocs/op |
| 完整 mTLS 握手 | 5 | 5.48 ms/conn | 1729328 B/op | 14829 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.518 ms/conn | 1931776 B/op | 15343 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.164 ms/conn | 1355936 B/op | 10219 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.772 ms/conn | 1014064 B/op | 7499 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.629 ms/conn | 2458680 B/op | 11366 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.748 ms/conn | 2620712 B/op | 12445 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.43 ms/conn | 3354936 B/op | 22599 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3572368 B/op | 20067 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1010 ms/pair | 241376 B/op | 994 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.392 ms/conn | 1592448 B/op | 9799 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.442 ms/conn | 1612768 B/op | 9939 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.486 ms/conn | 1771808 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.355 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.128 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 7.662 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 7.637 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.366 ms/conn | 63712 B/op | 57 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.983 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.014 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.187 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 10.2 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190232 B/op | 1203 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.402 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.512 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 10.86 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 42.68 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 56.52 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 708.2 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.658 us/op | 1540.9 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 2.211 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 229.8 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 179.1 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.191 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 4.027 us/op | 1017.12 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 416.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 99.18 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 410.3 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 27.97 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 645.2 ns/op | 1859.81 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 608.2 ns/op | 1972.91 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 615.5 ns/op | 1949.58 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 28.536 us/op | 2296.59 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 591.5 ns/op | 2028.76 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 30.65 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.364 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.709 us/op | 323.57 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.265 us/op | 948.66 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 12.34 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.889 us/op | 308.58 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.736 us/op | 691.4 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.428 us/op | 127.28 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.67 us/op | 327 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.785 us/op | 317.01 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 5.006 us/op | 239.7 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.356 us/op | 884.67 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.599 us/op | 333.38 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.331 us/op | 901.53 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.377 us/op | 871.3 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.956 us/op | 613.48 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 3.139 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 7.303 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.826 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 4.085 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.743 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.848 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.154 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.564 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.61 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 7.349 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 12.024 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 8.494 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 18.192 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 789.1 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 2.124 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 18.57 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 835.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 812.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.817 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.704 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 18.57 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.862 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.869 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.629 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 3 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 5.259 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 2.111 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 4.325 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 7.682 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 4.16 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 7.591 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 396.5 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 750.7 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 127.8 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 77.84 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 331.2 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 262.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 344.5 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 641.4 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 58.05 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 611 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 103.7 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 49.22 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 83.18 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 857.9 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 110 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 78.12 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 66.14 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 704.8 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 471.7 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 10.69 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 69.33 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 105.1 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 67 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 26.31 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 327.9 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 169.7 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 75.24 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.21 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 883.9 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 890.1 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 18.48 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 73.23 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.043 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.105 us/op | 4264 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
