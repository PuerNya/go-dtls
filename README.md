# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- 生成时间: `2026-09-12T02:10:36Z`
- Go: `go version go1.26.8 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `c74d6c669e9622dec03eff50931d66141b896433 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 613.245 us/op | 94505 B/op | 688 allocs/op |
| 完整 mTLS 握手 | 5 | 899.261 us/op | 108517 B/op | 865 allocs/op |
| mTLS 会话恢复握手 | 5 | 448.435 us/op | 116037 B/op | 804 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.169 ms/op | 116394 B/op | 1040 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.573 ms/op | 135155 B/op | 1334 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 701.44 us/op | 113729 B/op | 912 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 956.952 us/op | 117769 B/op | 937 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 956.833 us/op | 117769 B/op | 937 allocs/op |
| 直接外部 PSK 握手 | 5 | 379.425 us/op | 98140 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.106 ms/op | 126178 B/op | 970 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.108 ms/op | 118540 B/op | 949 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.837 ms/op | 165287 B/op | 1397 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.837 ms/op | 153253 B/op | 1358 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.042 ms/op | 143817 B/op | 1188 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.045 ms/op | 146594 B/op | 1209 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 943.89 us/op | 142341 B/op | 720 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 895.098 us/op | 145541 B/op | 750 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.28 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.091 ms/conn | 2410416 B/op | 20796 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.126 ms/conn | 2419040 B/op | 21115 allocs/op |
| 完整 mTLS 握手 | 5 | 3.726 ms/conn | 3086064 B/op | 28774 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.024 ms/conn | 3579816 B/op | 30792 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.147 ms/conn | 2899440 B/op | 24725 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5295 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.084 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.207 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.428 ms/conn | 3879720 B/op | 38875 allocs/op |
| 会话恢复握手 | 5 | 0.5824 ms/conn | 4293592 B/op | 37039 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.7092 ms/conn | 6302320 B/op | 48617 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6254 ms/conn | 280144 B/op | 1851 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.45 ms/conn | 3281832 B/op | 21118 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.445 ms/conn | 3320872 B/op | 21438 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.792 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.028 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.142 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 6.628 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.642 ms/conn | 1311024 B/op | 11342 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.047 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.007 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.016 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.188 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.567 ms/conn | 1373584 B/op | 12187 allocs/op |
| 会话恢复握手 | 5 | 1.013 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.028 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.284 ms/conn | 1829440 B/op | 10762 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.383 ms/conn | 1841904 B/op | 10882 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.298 ms/conn | 1229024 B/op | 8975 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.226 ms/conn | 1733464 B/op | 10465 allocs/op |
| 完整 mTLS 握手 | 5 | 6.421 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.355 ms/conn | 1869392 B/op | 14945 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.442 ms/conn | 1293984 B/op | 9895 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.79 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.193 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.263 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.73 ms/conn | 2639672 B/op | 22174 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 202176 B/op | 964 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.525 ms/conn | 1476032 B/op | 9455 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.686 ms/conn | 1496192 B/op | 9595 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.976 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.901 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.628 ms/conn | 543560 B/op | 1185 allocs/op |
| 完整 mTLS 握手 | 5 | 8.761 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.79 ms/conn | 34872 B/op | 53 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.738 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.044 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.359 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.352 ms/conn | 543880 B/op | 1185 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.21 ms/conn | 551704 B/op | 1185 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.805 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.767 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 12.03 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 43.31 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 57.7 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 704.7 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.793 us/op | 1466.33 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 2.225 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 243.1 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 197.5 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.257 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 4.611 us/op | 888.4 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 418.6 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 98.52 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 445.6 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 88.27 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 35.09 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 49.25 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 46.12 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 609.1 ns/op | 1970.16 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 547.7 ns/op | 2190.94 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 595.1 ns/op | 2016.55 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 26.272 us/op | 2494.55 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 567.9 ns/op | 2112.9 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 30.97 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.233 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 4.02 us/op | 298.48 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.48 us/op | 810.88 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 10.93 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 4.243 us/op | 282.8 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.968 us/op | 609.89 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 10.741 us/op | 111.72 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 4.065 us/op | 295.22 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 4.17 us/op | 287.78 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 5.078 us/op | 236.31 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.505 us/op | 797.55 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 4.315 us/op | 278.11 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.545 us/op | 776.73 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.599 us/op | 750.63 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 2.04 us/op | 588.25 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.968 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.7 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.731 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.871 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.678 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.712 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.234 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.235 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.557 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.485 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.844 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 11.388 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 8.077 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 18.019 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 740.7 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.856 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 9.175 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 759.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 754.7 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.662 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.294 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.159 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.688 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.691 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.695 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 3.112 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 5.177 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 2.067 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.928 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 7.097 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.903 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.862 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 404.9 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 742 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 124.4 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 87.9 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 318 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 264 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 378.3 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 599.2 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 55.35 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 507.6 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 87.72 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 49.39 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 74.24 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 709.4 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 89.32 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 74.11 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 77.4 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 661.8 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 473 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.38 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 68.57 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 119.7 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 80.24 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 35.56 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 341.1 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 193.2 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 85.68 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.251 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 915.7 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 928.9 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.32 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 68.93 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 7.41 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 7.085 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
