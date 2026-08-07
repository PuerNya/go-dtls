# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- 生成时间: `2026-08-07T02:23:04Z`
- Go: `go version go1.26.5 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `80d296934876a0134a476322da8fb3595e307ece (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 455.45 us/op | 94505 B/op | 688 allocs/op |
| 完整 mTLS 握手 | 5 | 705.712 us/op | 108518 B/op | 865 allocs/op |
| mTLS 会话恢复握手 | 5 | 357.416 us/op | 116072 B/op | 805 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 908.826 us/op | 116177 B/op | 1038 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.229 ms/op | 135124 B/op | 1334 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 540.987 us/op | 113728 B/op | 912 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 741.994 us/op | 117769 B/op | 937 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 746.506 us/op | 117769 B/op | 937 allocs/op |
| 直接外部 PSK 握手 | 5 | 299.049 us/op | 98142 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 836.465 us/op | 126177 B/op | 970 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 834.336 us/op | 118536 B/op | 949 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.395 ms/op | 165284 B/op | 1397 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.4 ms/op | 153250 B/op | 1358 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 818.361 us/op | 143851 B/op | 1188 allocs/op |
| ECH 握手 / 经 HRR | 5 | 821.296 us/op | 146596 B/op | 1209 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 701.945 us/op | 142338 B/op | 720 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 657.882 us/op | 145538 B/op | 750 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 1.732 ms/op | 170980 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.611 ms/conn | 2410368 B/op | 20795 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.626 ms/conn | 2419040 B/op | 21115 allocs/op |
| 完整 mTLS 握手 | 5 | 2.839 ms/conn | 3085816 B/op | 28764 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.051 ms/conn | 3579528 B/op | 30782 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.628 ms/conn | 2899440 B/op | 24725 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4003 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.606 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.696 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.236 ms/conn | 3879200 B/op | 38883 allocs/op |
| 会话恢复握手 | 5 | 0.4607 ms/conn | 4293592 B/op | 37023 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.52 ms/conn | 6303312 B/op | 48630 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.4637 ms/conn | 279896 B/op | 1848 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.905 ms/conn | 3281848 B/op | 21118 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 1.88 ms/conn | 3320872 B/op | 21438 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.921 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 3.879 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 3.995 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 5.135 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.136 ms/conn | 1311024 B/op | 11342 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 3.869 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.7981 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 3.885 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.047 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.104 ms/conn | 1372560 B/op | 12186 allocs/op |
| 会话恢复握手 | 5 | 0.817 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8224 ms/conn | 2228440 B/op | 18824 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.119 ms/conn | 1829504 B/op | 10763 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 4.891 ms/conn | 1841904 B/op | 10882 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.805 ms/conn | 1229024 B/op | 8975 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 3.981 ms/conn | 1733464 B/op | 10465 allocs/op |
| 完整 mTLS 握手 | 5 | 4.887 ms/conn | 1665784 B/op | 14422 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.902 ms/conn | 1870416 B/op | 14957 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.801 ms/conn | 1293984 B/op | 9895 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.619 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 3.964 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 3.987 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 10.28 ms/conn | 2640512 B/op | 22184 allocs/op |
| 会话恢复握手 | 5 | 1007 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1007 ms/pair | 201928 B/op | 961 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 1.909 ms/conn | 1476032 B/op | 9455 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.708 ms/conn | 1496192 B/op | 9595 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 5.279 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 3.797 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 6.237 ms/conn | 543488 B/op | 1183 allocs/op |
| 完整 mTLS 握手 | 5 | 6.671 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.658 ms/conn | 34944 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 3.607 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.736 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 6.232 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 6.273 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 9.275 ms/conn | 551680 B/op | 1184 allocs/op |
| 会话恢复握手 | 5 | 1010 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS 会话恢复握手 | 5 | 1013 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 3.777 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 4.911 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.925 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 30.97 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 39.17 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 435.4 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.721 us/op | 2379.34 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.516 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 185.3 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 148.8 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 883.9 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.057 us/op | 1339.96 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 277.5 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 92.77 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 288.3 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 60.7 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 26.87 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 38.13 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 35.75 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 363 ns/op | 3305.44 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 334.6 ns/op | 3586.91 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 335.6 ns/op | 3575.17 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 17.527 us/op | 3739.12 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 321.7 ns/op | 3730.54 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 22.54 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 3.55 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.79 us/op | 430.18 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.034 us/op | 1160.36 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 8.741 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.567 us/op | 336.41 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.588 us/op | 755.81 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.799 us/op | 153.87 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.111 us/op | 385.73 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.17 us/op | 378.6 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.702 us/op | 324.15 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.229 us/op | 976.34 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.257 us/op | 368.49 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.31 us/op | 915.83 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.313 us/op | 913.73 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.676 us/op | 716.06 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.118 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 4.934 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.188 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 2.713 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.155 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 2.67 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.093 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.094 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.063 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 2.466 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 4.526 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 7.629 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 5.377 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 12.329 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 532.5 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.503 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 7.081 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 551.1 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 549.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.228 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 3.165 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 7.086 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.255 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.258 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 1.679 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 1.939 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 3.439 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.391 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 2.802 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 5.211 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 2.743 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 4.959 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 281 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 523 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 92.1 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 68.46 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 245.3 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 204.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 276.8 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 407.2 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 39.95 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 386.2 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 62.94 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 36.24 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 53.84 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 520 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 64.12 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 53.82 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 55.26 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 455.7 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 330.2 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 10.39 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 50.1 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 82.33 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 55.56 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 27.27 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 235.8 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 127.4 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 67.89 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 859 ns/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 627.6 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 622 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 9.619 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 48.27 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.557 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 4.616 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
