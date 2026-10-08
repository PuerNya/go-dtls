# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- 生成时间: `2026-10-08T03:38:02Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `799bd17483efd76f30bd0f17cb080a0fc0ea92da (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 560.167 us/op | 97993 B/op | 710 allocs/op |
| 完整 mTLS 握手 | 5 | 843.364 us/op | 112215 B/op | 894 allocs/op |
| mTLS 会话恢复握手 | 5 | 416.825 us/op | 119419 B/op | 823 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.085 ms/op | 120038 B/op | 1068 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.455 ms/op | 138942 B/op | 1362 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 647.931 us/op | 117567 B/op | 934 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 891.881 us/op | 121529 B/op | 966 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 895.227 us/op | 121529 B/op | 966 allocs/op |
| 直接外部 PSK 握手 | 5 | 355.699 us/op | 101438 B/op | 744 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.014 ms/op | 130098 B/op | 993 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.007 ms/op | 122315 B/op | 971 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.683 ms/op | 169843 B/op | 1428 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.687 ms/op | 157393 B/op | 1387 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 957.264 us/op | 147499 B/op | 1210 allocs/op |
| ECH 握手 / 经 HRR | 5 | 958.16 us/op | 150275 B/op | 1231 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 861.323 us/op | 145972 B/op | 743 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 814.283 us/op | 149316 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.08 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.931 ms/conn | 2487488 B/op | 21255 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.965 ms/conn | 2496160 B/op | 21575 allocs/op |
| 完整 mTLS 握手 | 5 | 3.351 ms/conn | 3170176 B/op | 29384 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.593 ms/conn | 3666696 B/op | 31403 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.94 ms/conn | 2976560 B/op | 25185 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5253 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.93 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.073 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.017 ms/conn | 3964168 B/op | 39514 allocs/op |
| 会话恢复握手 | 5 | 0.5858 ms/conn | 4443080 B/op | 37899 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6384 ms/conn | 6461928 B/op | 49652 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6372 ms/conn | 287904 B/op | 1898 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.295 ms/conn | 3358952 B/op | 21578 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.271 ms/conn | 3400872 B/op | 21918 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.493 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.503 ms/conn | 1187344 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.578 ms/conn | 1207344 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 5.996 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.963 ms/conn | 1343344 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.477 ms/conn | 1425744 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7904 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.496 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.667 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.904 ms/conn | 1398992 B/op | 12384 allocs/op |
| 会话恢复握手 | 5 | 0.9694 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.9356 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.747 ms/conn | 1859000 B/op | 10964 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.76 ms/conn | 1871480 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.305 ms/conn | 1290720 B/op | 9299 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.962 ms/conn | 2445144 B/op | 10805 allocs/op |
| 完整 mTLS 握手 | 5 | 5.976 ms/conn | 1729480 B/op | 14832 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.02 ms/conn | 1932896 B/op | 15355 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.272 ms/conn | 1356320 B/op | 10219 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.923 ms/conn | 1014288 B/op | 7499 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.958 ms/conn | 2459256 B/op | 11366 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.024 ms/conn | 2620664 B/op | 12445 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.52 ms/conn | 3354064 B/op | 22600 allocs/op |
| 会话恢复握手 | 5 | 1008 ms/pair | 3571568 B/op | 20067 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1008 ms/pair | 241376 B/op | 997 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.543 ms/conn | 1592768 B/op | 9799 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.587 ms/conn | 1612768 B/op | 9939 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.955 ms/conn | 1771424 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.896 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.082 ms/conn | 1169496 B/op | 1202 allocs/op |
| 完整 mTLS 握手 | 5 | 8.582 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.528 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.718 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.057 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.986 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.024 ms/conn | 1169408 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.37 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.536 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.899 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 9.797 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 37.52 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 46.1 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 552.6 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.158 us/op | 1898 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.73 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 208.4 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 174.3 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 975.5 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.681 us/op | 1112.77 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 358.8 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 92.05 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 356.5 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 19.5 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 518.7 ns/op | 2313.34 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 467.5 ns/op | 2566.6 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 491.7 ns/op | 2440.58 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 23.701 us/op | 2765.12 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 430 ns/op | 2790.91 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 24.64 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.366 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.591 us/op | 463.19 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 976.8 ns/op | 1228.5 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.88 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.911 us/op | 412.28 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.519 us/op | 789.82 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.6 us/op | 157.9 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.567 us/op | 467.46 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.657 us/op | 451.67 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.382 us/op | 354.77 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.059 us/op | 1132.74 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.264 us/op | 367.69 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.042 us/op | 1151.12 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.113 us/op | 1077.75 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.45 us/op | 827.31 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 3.022 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.546 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.496 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.343 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.447 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.401 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.247 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.407 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.289 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.162 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 9.941 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.104 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 15.807 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 731.3 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.812 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.978 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 731.6 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 750.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.595 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.139 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 8.957 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.618 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.649 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.27 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.606 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.548 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.872 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.73 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.552 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.579 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.4 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 285.4 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 584.2 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 103.1 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 79.19 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 277.1 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 234.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 350.8 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 570.4 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 51.37 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 578.4 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 88.92 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 43.75 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 72.93 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 789.8 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 89.83 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 69.5 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 64.52 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 613.8 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 414.6 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.42 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 55.43 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 91.67 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 61.13 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 29.98 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 294.3 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 155.8 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 78.33 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.223 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 817.1 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 802.7 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 23.06 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 58.87 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 4.922 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 5.926 us/op | 4264 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
