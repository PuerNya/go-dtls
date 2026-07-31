# 自动化基准测试结果

[简体中文](benchmark.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `b512a578bc6b664ffddb0c20b5ad38882bfa2941`
- 生成时间: `2026-07-31T16:52:59Z`
- Go: `go version go1.26.5 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `41b7a0209abbddc579d3d861f36c0f574ae7e907 (Linux Release static)`

共 169 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (14)](#section-connection-lifecycle)
- [真实 UDP 互通 (52)](#section-real-udp-interoperability)
  - [go-dtls 客户端 -> go-dtls 服务端 (14)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls 客户端 -> wolfSSL 服务端 (12)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL 客户端 -> go-dtls 服务端 (13)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL 客户端 -> wolfSSL 服务端 (13)](#real-udp-wolfssl-client-wolfssl-server)
- [记录层与可靠性 (37)](#section-record-layer-and-reliability)
- [密钥调度与密码学 (38)](#section-key-schedule-and-cryptography)
- [报文编码与解析 (26)](#section-wire-encoding-and-parsing)
- [证书压缩 (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## 连接生命周期

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 586.707 us/op | 99289 B/op | 759 allocs/op |
| 完整 mTLS 握手 | 5 | 870.204 us/op | 115950 B/op | 973 allocs/op |
| mTLS 会话恢复握手 | 5 | 431.089 us/op | 116187 B/op | 804 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 682.04 us/op | 118512 B/op | 983 allocs/op |
| 直接外部 PSK 握手 | 5 | 366.787 us/op | 98268 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.034 ms/op | 130938 B/op | 1041 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.047 ms/op | 123300 B/op | 1020 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.725 ms/op | 172667 B/op | 1505 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.73 ms/op | 160571 B/op | 1466 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 993.077 us/op | 148602 B/op | 1259 allocs/op |
| ECH 握手 / 经 HRR | 5 | 996.413 us/op | 151378 B/op | 1280 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 905.079 us/op | 147125 B/op | 791 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 860.261 us/op | 150325 B/op | 821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.088 ms/op | 175766 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.965 ms/conn | 2621456 B/op | 24336 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.996 ms/conn | 2630128 B/op | 24656 allocs/op |
| 完整 mTLS 握手 | 5 | 3.481 ms/conn | 3471400 B/op | 33631 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.001 ms/conn | 3110528 B/op | 28266 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5129 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.926 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.04 ms/conn | 2811760 B/op | 26008 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.031 ms/conn | 4316408 B/op | 45071 allocs/op |
| 会话恢复握手 | 5 | 0.5705 ms/conn | 4503400 B/op | 40579 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6273 ms/conn | 6622496 B/op | 53481 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.5905 ms/conn | 290696 B/op | 2028 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.288 ms/conn | 3492888 B/op | 24658 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.309 ms/conn | 3531928 B/op | 24978 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.456 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.736 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.826 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 6.455 ms/conn | 1522544 B/op | 14882 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.773 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9095 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.636 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.837 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.35 ms/conn | 1584400 B/op | 15727 allocs/op |
| 会话恢复握手 | 5 | 1.06 ms/conn | 2068464 B/op | 18002 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.063 ms/conn | 2438544 B/op | 22362 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.908 ms/conn | 1829504 B/op | 10763 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.891 ms/conn | 1841904 B/op | 10882 allocs/op |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.186 ms/conn | 1440048 B/op | 12515 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.805 ms/conn | 1944184 B/op | 14005 allocs/op |
| 完整 mTLS 握手 | 5 | 5.851 ms/conn | 1817536 B/op | 16819 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.186 ms/conn | 1504704 B/op | 13435 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.791 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.745 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.802 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.32 ms/conn | 2734592 B/op | 23438 allocs/op |
| 会话恢复握手 | 5 | 1008 ms/pair | 2984824 B/op | 22965 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1008 ms/pair | 212712 B/op | 1141 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.375 ms/conn | 1686768 B/op | 12995 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.375 ms/conn | 1706928 B/op | 13135 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.441 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.518 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.477 ms/conn | 543560 B/op | 1185 allocs/op |
| 完整 mTLS 握手 | 5 | 7.965 ms/conn | 34872 B/op | 53 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.498 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.919 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.46 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.418 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.11 ms/conn | 551680 B/op | 1184 allocs/op |
| 会话恢复握手 | 5 | 1011 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS 会话恢复握手 | 5 | 1015 ms/pair | 557760 B/op | 1184 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.401 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.88 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 9.634 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 42.16 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 56.52 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 688.5 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.376 us/op | 1724.24 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.899 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 238.2 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 188.1 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.047 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 4.18 us/op | 979.94 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 382.9 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 110.2 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 378.1 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 90.3 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 42.19 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 46.44 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 40.54 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 568.2 ns/op | 2112.11 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 525 ns/op | 2285.89 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 517.1 ns/op | 2320.45 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 27.012 us/op | 2426.18 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 489.2 ns/op | 2452.74 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 31.54 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.373 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.043 us/op | 394.41 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.158 us/op | 1035.97 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 11.29 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.359 us/op | 357.29 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.605 us/op | 747.53 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 8.798 us/op | 136.39 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.056 us/op | 392.67 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.246 us/op | 369.66 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.828 us/op | 313.45 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.162 us/op | 1032.53 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.708 us/op | 323.64 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.158 us/op | 1036.44 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.243 us/op | 965.43 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.57 us/op | 764.52 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.759 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.251 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.505 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.369 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.459 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.309 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.375 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.156 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.371 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.788 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.077 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 15.718 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 691.2 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.842 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.993 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 732.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 730.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.633 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.084 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 8.956 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.662 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.645 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.221 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.514 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.424 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.835 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.713 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.684 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.717 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.549 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 330.3 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 660.5 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 114.9 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 79.39 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 288.8 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 239.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 340.6 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 594.1 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 52.21 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 523.9 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 85.73 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 47.4 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 70.1 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 725.1 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 88.41 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 71.29 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 66.36 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 639.6 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 451 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 70.16 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 109.3 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 68.96 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 28.59 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 313.5 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 152.9 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 75.91 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.157 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 821.3 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 823.1 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.79 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 71.54 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 6.721 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 6.652 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
