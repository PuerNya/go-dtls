# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `8010fee127101d213fb208cbfbff51779d6c6954`
- 生成时间: `2026-10-07T03:18:51Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `ab4a14619a4eea21de84603f5e4f581b6296b3b5 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 564.047 us/op | 97817 B/op | 709 allocs/op |
| 完整 mTLS 握手 | 5 | 844.971 us/op | 111974 B/op | 891 allocs/op |
| mTLS 会话恢复握手 | 5 | 422.792 us/op | 119167 B/op | 822 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.109 ms/op | 119873 B/op | 1065 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.484 ms/op | 138667 B/op | 1358 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 665.229 us/op | 117040 B/op | 933 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 913.195 us/op | 121225 B/op | 963 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 916.742 us/op | 121225 B/op | 963 allocs/op |
| 直接外部 PSK 握手 | 5 | 352.549 us/op | 101261 B/op | 743 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.025 ms/op | 129635 B/op | 992 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.027 ms/op | 121852 B/op | 970 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.705 ms/op | 169033 B/op | 1425 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.715 ms/op | 156580 B/op | 1384 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 967.604 us/op | 147130 B/op | 1209 allocs/op |
| ECH 握手 / 经 HRR | 5 | 967.823 us/op | 149906 B/op | 1230 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 891.096 us/op | 145797 B/op | 742 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 829.335 us/op | 149141 B/op | 773 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.122 ms/op | 174582 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.913 ms/conn | 2482048 B/op | 21235 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.954 ms/conn | 2490720 B/op | 21555 allocs/op |
| 完整 mTLS 握手 | 5 | 3.346 ms/conn | 3162960 B/op | 29326 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.601 ms/conn | 3658536 B/op | 31347 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.952 ms/conn | 2971120 B/op | 25165 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5253 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.921 ms/conn | 2502016 B/op | 22353 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.076 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.02 ms/conn | 3957032 B/op | 39444 allocs/op |
| 会话恢复握手 | 5 | 0.5786 ms/conn | 4431648 B/op | 37859 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6418 ms/conn | 6449856 B/op | 49589 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6442 ms/conn | 287200 B/op | 1896 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.319 ms/conn | 3353512 B/op | 21558 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.298 ms/conn | 3395432 B/op | 21898 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.584 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.526 ms/conn | 1186064 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.695 ms/conn | 1204784 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 6.132 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.13 ms/conn | 1341744 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.55 ms/conn | 1424464 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9508 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.624 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.866 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.081 ms/conn | 1396688 B/op | 12383 allocs/op |
| 会话恢复握手 | 5 | 1.069 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.085 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.039 ms/conn | 1857800 B/op | 10965 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.098 ms/conn | 1870200 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.177 ms/conn | 1305760 B/op | 9219 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.858 ms/conn | 2406104 B/op | 10725 allocs/op |
| 完整 mTLS 握手 | 5 | 5.929 ms/conn | 1743336 B/op | 14712 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.974 ms/conn | 1946232 B/op | 15226 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.162 ms/conn | 1370720 B/op | 10139 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.771 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.894 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.953 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.61 ms/conn | 3314392 B/op | 22480 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3492336 B/op | 19887 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 238008 B/op | 987 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.39 ms/conn | 1555648 B/op | 9719 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.502 ms/conn | 1575808 B/op | 9859 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.694 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.542 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.675 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 8.104 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 7.974 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.484 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.025 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.872 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.601 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.42 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.494 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.407 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 11.54 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 38.75 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 48.41 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 555.4 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.229 us/op | 1837.97 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.835 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 212.6 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 183.3 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.017 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.737 us/op | 1096.21 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 368.8 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 92.33 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 359.9 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 63.95 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 22.21 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 42.95 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 22.47 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 527.7 ns/op | 2273.93 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 508.4 ns/op | 2360.41 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 501.7 ns/op | 2391.99 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 24.975 us/op | 2624.08 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 436 ns/op | 2752.06 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 25.36 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.373 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.977 us/op | 403.07 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.042 us/op | 1151.53 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 11.55 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.99 us/op | 401.33 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.617 us/op | 741.9 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 8.222 us/op | 145.95 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.918 us/op | 411.26 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.902 us/op | 413.5 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.549 us/op | 338.17 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.213 us/op | 989.46 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.482 us/op | 344.61 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.202 us/op | 998.45 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.235 us/op | 971.51 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.625 us/op | 738.32 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.831 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.54 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.53 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.526 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.47 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.389 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 0.9399 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.404 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.271 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.128 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.183 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.154 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.015 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 721.8 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.797 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 8.948 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 749.9 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 752 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.552 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.131 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 8.945 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.604 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.6 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.251 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.559 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.628 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.799 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.709 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.659 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.61 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.389 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 309 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 619.2 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 108.3 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 78.93 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 289.1 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 246.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 344.4 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 559.5 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 50.51 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 560.7 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 87.84 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 43.98 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 71.97 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 776.1 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 90.05 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 69.74 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 65.41 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 634.7 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 426.6 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 63.13 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 93.19 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 61.05 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 28.79 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 289.3 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 157.9 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 77.14 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.138 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 804.9 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 816.4 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.48 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 59.44 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 4.952 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.631 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
