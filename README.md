# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- 生成时间: `2026-08-04T17:54:13Z`
- Go: `go version go1.26.5 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `0fb71c572f1e27e281e0b83a95881ddba0073a73 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 555.571 us/op | 94505 B/op | 688 allocs/op |
| 完整 mTLS 握手 | 5 | 827.821 us/op | 108517 B/op | 865 allocs/op |
| mTLS 会话恢复握手 | 5 | 424.888 us/op | 116058 B/op | 804 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.072 ms/op | 116397 B/op | 1040 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.457 ms/op | 135171 B/op | 1333 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 652.18 us/op | 113730 B/op | 912 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 885.311 us/op | 117769 B/op | 937 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 891.775 us/op | 117769 B/op | 937 allocs/op |
| 直接外部 PSK 握手 | 5 | 350.614 us/op | 98141 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.037 ms/op | 126178 B/op | 970 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.032 ms/op | 118539 B/op | 949 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.697 ms/op | 165288 B/op | 1397 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.711 ms/op | 153123 B/op | 1358 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 947.314 us/op | 143818 B/op | 1188 allocs/op |
| ECH 握手 / 经 HRR | 5 | 950.52 us/op | 146594 B/op | 1209 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 863.302 us/op | 142340 B/op | 720 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 818.513 us/op | 145540 B/op | 750 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.056 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.908 ms/conn | 2410320 B/op | 20794 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.952 ms/conn | 2419040 B/op | 21115 allocs/op |
| 完整 mTLS 握手 | 5 | 3.376 ms/conn | 3085272 B/op | 28764 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.596 ms/conn | 3580472 B/op | 30774 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.951 ms/conn | 2899440 B/op | 24725 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5175 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.94 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.036 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.009 ms/conn | 3878064 B/op | 38884 allocs/op |
| 会话恢复握手 | 5 | 0.5584 ms/conn | 4293592 B/op | 37039 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6593 ms/conn | 6303232 B/op | 48610 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6126 ms/conn | 280144 B/op | 1851 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.313 ms/conn | 3281832 B/op | 21118 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.278 ms/conn | 3320872 B/op | 21438 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.455 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.609 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.666 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 6.13 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.101 ms/conn | 1311024 B/op | 11342 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.551 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.8398 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.521 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.745 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.031 ms/conn | 1371536 B/op | 12185 allocs/op |
| 会话恢复握手 | 5 | 1.049 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.093 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.974 ms/conn | 1829424 B/op | 10762 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 5.911 ms/conn | 1841904 B/op | 10882 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.181 ms/conn | 1229024 B/op | 8975 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.718 ms/conn | 1733464 B/op | 10465 allocs/op |
| 完整 mTLS 握手 | 5 | 5.824 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.873 ms/conn | 1869440 B/op | 14945 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.181 ms/conn | 1293984 B/op | 9895 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.751 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.697 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.775 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.25 ms/conn | 2640016 B/op | 22178 allocs/op |
| 会话恢复握手 | 5 | 1008 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1008 ms/pair | 201928 B/op | 961 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.371 ms/conn | 1476032 B/op | 9455 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.345 ms/conn | 1496192 B/op | 9595 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.44 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.576 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.402 ms/conn | 543488 B/op | 1183 allocs/op |
| 完整 mTLS 握手 | 5 | 7.904 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.129 ms/conn | 34872 B/op | 53 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.556 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.917 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.845 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.59 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.08 ms/conn | 551680 B/op | 1184 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS 会话恢复握手 | 5 | 1015 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.445 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.59 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 9.505 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 38.36 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 49.81 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 533.5 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.136 us/op | 1917.23 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.916 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 221.9 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 177.6 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.077 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.61 us/op | 1134.7 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 370.8 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 113.7 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 350.9 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 78.35 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 42.63 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 46.08 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 40.44 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 555.2 ns/op | 2161.22 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 521.7 ns/op | 2300.18 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 525.8 ns/op | 2282.45 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 26.597 us/op | 2464.04 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 448.5 ns/op | 2675.49 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 28.02 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.366 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.414 us/op | 351.45 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.259 us/op | 953 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 11.25 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.118 us/op | 384.87 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.669 us/op | 719.01 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.47 us/op | 126.71 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.122 us/op | 384.42 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.557 us/op | 337.4 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 4.344 us/op | 276.24 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.123 us/op | 1068.98 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.71 us/op | 323.44 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.215 us/op | 987.32 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.297 us/op | 924.92 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.706 us/op | 703.28 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.765 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.235 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.51 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.393 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.462 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.411 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.384 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.134 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.313 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.156 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.085 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 15.715 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 693.2 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.816 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.964 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 724.8 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 722.7 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.568 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.167 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 8.953 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.626 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.601 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.234 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.66 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.516 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.84 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.731 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.604 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.603 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.505 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 340.5 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 647.5 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 115.4 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 79.41 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 297.5 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 234.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 338.5 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 565.1 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 52.92 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 525.2 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 86.08 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 46.59 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 70.44 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 729.3 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 87.2 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 71.54 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 63.72 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 629.2 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 412.5 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 61.75 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 111.1 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 76.6 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 28.62 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 349.3 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 170.4 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 75.38 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.258 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 838 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 817.3 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.78 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 60.99 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 6.732 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 6.599 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
