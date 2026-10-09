# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `10300127dca519bf9c8f7672c1119d1baaa0893c`
- 生成时间: `2026-10-09T19:11:06Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `23f245c02c4319bb80baf7f3ae0c1de40e3fb68e (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 585.754 us/op | 101978 B/op | 741 allocs/op |
| 完整 mTLS 握手 | 5 | 880.564 us/op | 116990 B/op | 938 allocs/op |
| mTLS 会话恢复握手 | 5 | 446.123 us/op | 123815 B/op | 864 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.104 ms/op | 124300 B/op | 1117 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.507 ms/op | 143085 B/op | 1413 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 683.684 us/op | 122018 B/op | 971 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 931.296 us/op | 126305 B/op | 1009 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 938.149 us/op | 126401 B/op | 1009 allocs/op |
| 直接外部 PSK 握手 | 5 | 375.985 us/op | 104398 B/op | 767 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.045 ms/op | 136683 B/op | 1026 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.041 ms/op | 126701 B/op | 1002 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.723 ms/op | 177745 B/op | 1478 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.737 ms/op | 162975 B/op | 1431 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 996.297 us/op | 152059 B/op | 1241 allocs/op |
| ECH 握手 / 经 HRR | 5 | 996.978 us/op | 154835 B/op | 1262 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 881.193 us/op | 149966 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 837.207 us/op | 153311 B/op | 805 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.134 ms/op | 178759 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.003 ms/conn | 2672528 B/op | 21796 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.02 ms/conn | 2683360 B/op | 22135 allocs/op |
| 完整 mTLS 握手 | 5 | 3.551 ms/conn | 3375560 B/op | 30130 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.592 ms/conn | 3834128 B/op | 31877 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.127 ms/conn | 3239136 B/op | 25832 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5599 ms/conn | 1786544 B/op | 15014 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.972 ms/conn | 2692544 B/op | 22915 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.126 ms/conn | 2852368 B/op | 23468 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.129 ms/conn | 4182056 B/op | 40442 allocs/op |
| 会话恢复握手 | 5 | 0.6747 ms/conn | 4738512 B/op | 39512 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8124 ms/conn | 6870664 B/op | 51513 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6873 ms/conn | 299504 B/op | 1972 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.359 ms/conn | 3543992 B/op | 22119 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.324 ms/conn | 3585912 B/op | 22459 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.627 ms/conn | 4081608 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.755 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.792 ms/conn | 1353584 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.25 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.278 ms/conn | 1495872 B/op | 11931 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.704 ms/conn | 1573536 B/op | 12619 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.004 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.715 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.901 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.116 ms/conn | 1546448 B/op | 12663 allocs/op |
| 会话恢复握手 | 5 | 1.152 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.148 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.992 ms/conn | 2006104 B/op | 11253 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.029 ms/conn | 2017800 B/op | 11365 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.365 ms/conn | 1326784 B/op | 9519 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.037 ms/conn | 2480184 B/op | 11025 allocs/op |
| 完整 mTLS 握手 | 5 | 6.032 ms/conn | 1770712 B/op | 15063 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.065 ms/conn | 1973920 B/op | 15589 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.358 ms/conn | 1391904 B/op | 10439 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.979 ms/conn | 1048624 B/op | 7679 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.04 ms/conn | 2495000 B/op | 11586 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.064 ms/conn | 2656344 B/op | 12685 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.74 ms/conn | 3405232 B/op | 22866 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3642288 B/op | 20507 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 244984 B/op | 1019 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.538 ms/conn | 1628448 B/op | 10019 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.646 ms/conn | 1648288 B/op | 10159 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.691 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.893 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.951 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 8.575 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.469 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.885 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.074 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.139 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.017 ms/conn | 1169408 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.58 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1190872 B/op | 1203 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.92 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.498 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 9.832 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 36.05 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 45.18 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 511.8 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.992 us/op | 2055.78 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.623 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 202.7 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 167.6 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 904.7 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.488 us/op | 1174.45 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 318.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 89.11 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 334.8 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 22.18 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 513.9 ns/op | 2334.89 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 456.9 ns/op | 2626.26 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 472.3 ns/op | 2540.94 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 23.667 us/op | 2769.03 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 408 ns/op | 2941.3 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 24.04 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.672 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.131 us/op | 563.04 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 836.8 ns/op | 1433.96 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.28 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.256 us/op | 531.92 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.268 us/op | 946.39 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 6.945 us/op | 172.78 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.09 us/op | 574.05 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.192 us/op | 547.4 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 2.94 us/op | 408.12 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 930.6 ns/op | 1289.45 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.048 us/op | 393.68 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 831.2 ns/op | 1443.62 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 889 ns/op | 1349.86 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.25 us/op | 960.02 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.721 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.118 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.487 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.303 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.434 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.237 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.354 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.067 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 15.16 ns/op | 0 B/op | 0 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 15.21 ns/op | 0 B/op | 0 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 6.907 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 15.379 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 687.9 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.802 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.957 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 711.3 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 755.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.534 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 3.998 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.044 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.575 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.57 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.228 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.51 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.327 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.78 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.669 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.388 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.533 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.211 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 262.5 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 525.3 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 103 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 78.97 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 266.4 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 234.5 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 340.3 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 554.6 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 49.94 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 572.1 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 86.28 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 43.33 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 71.62 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 792.8 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 87.46 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 71.54 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 64.4 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 601.1 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 402.9 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 53.26 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 90.29 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 59.85 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 29.97 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 278 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 151 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 78.53 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.1 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 797.1 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 806.7 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 23.11 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 55.74 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 4.909 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.439 us/op | 4296 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
