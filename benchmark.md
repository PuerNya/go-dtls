# 自动化基准测试结果

[简体中文](benchmark.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `196ad5c06d35e7c73aae7fbcd1ce0be1f4e2be5d`
- 生成时间: `2026-07-31T15:43:55Z`
- Go: `go version go1.26.5 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `41b7a0209abbddc579d3d861f36c0f574ae7e907 (Linux Release static)`

共 168 项结果，按工作负载分组，并按功能、基准测试名称排序。数值为最终测试运行所输出样本的中位数。

工作负载专用的连接指标优先于 Go 基准测试框架耗时。内存和分配次数仍按每次 Go 基准测试操作统计。精确原始输出保留在 Workflow Artifact 中。

## 快速跳转

- [连接生命周期 (13)](#section-connection-lifecycle)
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
| 证书认证完整握手 / AES-128-GCM | 5 | 583.338 us/op | 99289 B/op | 759 allocs/op |
| 完整 mTLS 握手 | 5 | 836.809 us/op | 115950 B/op | 973 allocs/op |
| mTLS 会话恢复握手 | 5 | 426.15 us/op | 116082 B/op | 803 allocs/op |
| 直接外部 PSK 握手 | 5 | 340.342 us/op | 98265 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.025 ms/op | 130938 B/op | 1041 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.032 ms/op | 123299 B/op | 1020 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.704 ms/op | 172668 B/op | 1505 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.708 ms/op | 160570 B/op | 1466 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 947.842 us/op | 148598 B/op | 1259 allocs/op |
| ECH 握手 / 经 HRR | 5 | 953.056 us/op | 151374 B/op | 1280 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 874.918 us/op | 147124 B/op | 791 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 827.7 us/op | 150324 B/op | 821 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.071 ms/op | 175766 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.965 ms/conn | 2621408 B/op | 24335 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.982 ms/conn | 2630128 B/op | 24656 allocs/op |
| 完整 mTLS 握手 | 5 | 3.538 ms/conn | 3406328 B/op | 33631 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.974 ms/conn | 3110480 B/op | 28265 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5162 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.941 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.101 ms/conn | 2811808 B/op | 26008 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.122 ms/conn | 4315104 B/op | 45058 allocs/op |
| 会话恢复握手 | 5 | 0.5727 ms/conn | 4500152 B/op | 40559 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.626 ms/conn | 6619792 B/op | 53470 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.578 ms/conn | 290440 B/op | 2026 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.281 ms/conn | 3492872 B/op | 24658 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.243 ms/conn | 3531928 B/op | 24978 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.451 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.566 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.671 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 6.115 ms/conn | 1522544 B/op | 14882 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.579 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8768 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.592 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.726 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.075 ms/conn | 1582352 B/op | 15724 allocs/op |
| 会话恢复握手 | 5 | 1.06 ms/conn | 2065264 B/op | 17982 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.131 ms/conn | 2435480 B/op | 22344 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.753 ms/conn | 1829504 B/op | 10763 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.729 ms/conn | 1841904 B/op | 10882 allocs/op |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.237 ms/conn | 1440048 B/op | 12515 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.808 ms/conn | 1944184 B/op | 14005 allocs/op |
| 完整 mTLS 握手 | 5 | 5.884 ms/conn | 1818312 B/op | 16828 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.239 ms/conn | 1504704 B/op | 13435 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.812 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.799 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.859 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.32 ms/conn | 2734080 B/op | 23432 allocs/op |
| 会话恢复握手 | 5 | 1008 ms/pair | 2984824 B/op | 22965 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1008 ms/pair | 212464 B/op | 1138 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.393 ms/conn | 1686768 B/op | 12995 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.383 ms/conn | 1706928 B/op | 13135 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.448 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.532 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.473 ms/conn | 543560 B/op | 1185 allocs/op |
| 完整 mTLS 握手 | 5 | 7.879 ms/conn | 34872 B/op | 53 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.269 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.941 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.439 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.429 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.2 ms/conn | 551680 B/op | 1184 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS 会话恢复握手 | 5 | 1015 ms/pair | 557784 B/op | 1185 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.393 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.253 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 11.32 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 40.95 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 52.42 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 575.2 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.208 us/op | 1855.12 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.992 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 230 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 181.6 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.119 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.839 us/op | 1066.85 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 372.6 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 114.7 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 362.1 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 83.14 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 42.89 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 46.39 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 40.48 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 564 ns/op | 2127.84 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 496.8 ns/op | 2415.39 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 495.5 ns/op | 2421.75 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 26.273 us/op | 2494.45 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 442.2 ns/op | 2713.46 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 29.18 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.37 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.282 us/op | 365.63 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.186 us/op | 1011.64 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 10.94 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.029 us/op | 396.2 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.644 us/op | 730.09 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.227 us/op | 130.05 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.474 us/op | 345.4 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.41 us/op | 351.91 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.926 us/op | 305.62 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.126 us/op | 1066.08 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.631 us/op | 330.52 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.176 us/op | 1020 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.229 us/op | 976.15 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.616 us/op | 742.51 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.81 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.384 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.576 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.431 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.493 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.373 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.397 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.151 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.004 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.007 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.196 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.248 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 698.4 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.786 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.954 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 724.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 721.3 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.559 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.105 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 8.967 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.602 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.597 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.309 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.565 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.623 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.892 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.713 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.464 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.515 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.383 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 352.9 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 681.3 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 114.3 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 79.47 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 286 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 239.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 335.7 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 555.8 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 52.49 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 528.6 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 85.69 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 46.02 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 69.78 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 743.4 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 88.03 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 71.05 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 63.91 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 674.8 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 422.5 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 65.4 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 101.8 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 66.93 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 28.61 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 306 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 153.7 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 75.45 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.17 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 836.3 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 839.8 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.79 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 65.2 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 6.703 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 6.585 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
