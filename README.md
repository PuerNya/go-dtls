# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- 生成时间: `2026-08-15T08:40:17Z`
- Go: `go version go1.26.5 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `889de4784d81cb4346db2f537313e7e699b730e6 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 590.991 us/op | 94505 B/op | 688 allocs/op |
| 完整 mTLS 握手 | 5 | 875.356 us/op | 108517 B/op | 865 allocs/op |
| mTLS 会话恢复握手 | 5 | 429.555 us/op | 116007 B/op | 804 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.143 ms/op | 116448 B/op | 1040 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.543 ms/op | 135201 B/op | 1334 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 692.719 us/op | 113729 B/op | 912 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 942.704 us/op | 117769 B/op | 937 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 947.328 us/op | 117769 B/op | 937 allocs/op |
| 直接外部 PSK 握手 | 5 | 363.326 us/op | 98140 B/op | 724 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.078 ms/op | 126178 B/op | 970 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.084 ms/op | 118539 B/op | 949 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.803 ms/op | 165291 B/op | 1397 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.803 ms/op | 153123 B/op | 1358 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.007 ms/op | 143817 B/op | 1188 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.017 ms/op | 146594 B/op | 1209 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 905.258 us/op | 142340 B/op | 720 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 850.371 us/op | 145540 B/op | 750 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.24 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.058 ms/conn | 2410368 B/op | 20795 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.091 ms/conn | 2419040 B/op | 21115 allocs/op |
| 完整 mTLS 握手 | 5 | 3.649 ms/conn | 3086560 B/op | 28780 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.917 ms/conn | 3580272 B/op | 30783 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.098 ms/conn | 2899440 B/op | 24725 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5266 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.067 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.2 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.449 ms/conn | 3945008 B/op | 38879 allocs/op |
| 会话恢复握手 | 5 | 0.575 ms/conn | 4293640 B/op | 37039 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6602 ms/conn | 6302984 B/op | 48624 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6219 ms/conn | 280144 B/op | 1851 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.448 ms/conn | 3281832 B/op | 21118 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.412 ms/conn | 3320872 B/op | 21438 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.772 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.92 ms/conn | 1157904 B/op | 10562 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.053 ms/conn | 1176624 B/op | 10822 allocs/op |
| 完整 mTLS 握手 | 5 | 6.549 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.556 ms/conn | 1311024 B/op | 11342 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.911 ms/conn | 1396304 B/op | 12122 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9679 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.969 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.174 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.462 ms/conn | 1374960 B/op | 12190 allocs/op |
| 会话恢复握手 | 5 | 0.9845 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.9968 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.196 ms/conn | 1829424 B/op | 10762 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.227 ms/conn | 1841904 B/op | 10882 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 不支持: wolfSSL 服务端无法完成该 DTLS 1.3 hybrid 握手；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.293 ms/conn | 1229024 B/op | 8975 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.108 ms/conn | 1733464 B/op | 10465 allocs/op |
| 完整 mTLS 握手 | 5 | 6.331 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.343 ms/conn | 1869160 B/op | 14942 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.297 ms/conn | 1293984 B/op | 9895 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.758 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.148 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.126 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.37 ms/conn | 2640416 B/op | 22183 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS 会话恢复握手 | - | 不支持: wolfSSL 客户端无法解析 go-dtls 的 mTLS session ticket；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 202176 B/op | 964 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.506 ms/conn | 1476032 B/op | 9455 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.586 ms/conn | 1496192 B/op | 9595 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.922 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.671 ms/conn | 34904 B/op | 53 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.097 ms/conn | 543488 B/op | 1183 allocs/op |
| 完整 mTLS 握手 | 5 | 8.617 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.611 ms/conn | 34872 B/op | 53 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.87 ms/conn | 34904 B/op | 53 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.042 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.994 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.017 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.99 ms/conn | 551704 B/op | 1185 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 不支持: wolfSSL 服务端在 HelloRetryRequest 后拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.782 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.379 ms/conn | 34888 B/op | 53 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 11.94 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 39.06 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 50.04 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 572.2 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.229 us/op | 1837.46 MB/s | 5040 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.899 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 226.5 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 190.8 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.074 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.886 us/op | 1053.93 MB/s | 5616 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 367.9 ns/op | - | 624 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 99.5 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 380.1 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 已分配 | 5 | 80.51 ns/op | - | 80 B/op | 1 allocs/op |
| 握手报文组待处理索引 / 复用窗口 | 5 | 34.66 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 待处理 | 5 | 49.16 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 46.09 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 567.6 ns/op | 2114.23 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 525.4 ns/op | 2283.91 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 537.7 ns/op | 2231.83 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 25.425 us/op | 2577.63 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 477.9 ns/op | 2510.77 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 28.56 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.577 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.154 us/op | 380.45 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.199 us/op | 1000.55 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 11.28 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.587 us/op | 334.54 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.742 us/op | 688.97 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.332 us/op | 128.59 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.159 us/op | 379.87 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.261 us/op | 367.96 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 4.197 us/op | 285.89 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.336 us/op | 898.17 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.871 us/op | 309.99 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.282 us/op | 936.35 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.345 us/op | 891.88 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.78 us/op | 674.11 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.76 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.399 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.558 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.508 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.497 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.433 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.4 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.199 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.122 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.139 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.152 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.151 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 698.3 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.794 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 9.133 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 723 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 720.1 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.59 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.133 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.151 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.632 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.629 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.229 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.56 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.561 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.827 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.684 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.706 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.572 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.548 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 344.3 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 622 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 117.7 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 88.09 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 311.7 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 263.9 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 360 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 559.2 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 53.64 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 502.8 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 86.1 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 47.71 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 71.88 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 688.9 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 86.24 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 71.84 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 71.3 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 604.6 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 434 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 63.24 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 106.5 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 72.47 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 35.15 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 300.3 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 166.1 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 87.65 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.15 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 823 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 819.8 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 12.39 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 62.93 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 7.095 us/op | 336 B/op | 4 allocs/op |
| 解压 | 5 | 5.968 us/op | 4248 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
