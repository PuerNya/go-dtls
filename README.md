# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `12d6cbda9e0eb434c12aaf19a17b774ff6042ec2`
- 生成时间: `2026-10-09T14:07:22Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `bb236b46bc87ef485a381e161c519d3bb9babff1 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 623.161 us/op | 101241 B/op | 741 allocs/op |
| 完整 mTLS 握手 | 5 | 927.9 us/op | 116190 B/op | 938 allocs/op |
| mTLS 会话恢复握手 | 5 | 466.147 us/op | 123014 B/op | 864 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.174 ms/op | 123498 B/op | 1117 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.584 ms/op | 142233 B/op | 1413 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 714.116 us/op | 121281 B/op | 971 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 982.958 us/op | 125505 B/op | 1009 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 984.943 us/op | 125601 B/op | 1009 allocs/op |
| 直接外部 PSK 握手 | 5 | 390.891 us/op | 103597 B/op | 767 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.118 ms/op | 135946 B/op | 1026 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.116 ms/op | 125965 B/op | 1002 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.852 ms/op | 176945 B/op | 1478 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.861 ms/op | 162175 B/op | 1431 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.07 ms/op | 150938 B/op | 1241 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.073 ms/op | 153715 B/op | 1262 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 949.964 us/op | 149230 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 900.332 us/op | 152575 B/op | 805 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.357 ms/op | 178021 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.177 ms/conn | 2655248 B/op | 21796 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.182 ms/conn | 2666080 B/op | 22135 allocs/op |
| 完整 mTLS 握手 | 5 | 3.917 ms/conn | 3357696 B/op | 30137 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.966 ms/conn | 3815376 B/op | 31885 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.183 ms/conn | 3156048 B/op | 25827 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5804 ms/conn | 1767984 B/op | 15014 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.138 ms/conn | 2675216 B/op | 22914 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.311 ms/conn | 2835088 B/op | 23468 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.588 ms/conn | 4165512 B/op | 40431 allocs/op |
| 会话恢复握手 | 5 | 0.6616 ms/conn | 4700912 B/op | 39510 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8547 ms/conn | 6830736 B/op | 51489 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6887 ms/conn | 298024 B/op | 1975 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.571 ms/conn | 3526712 B/op | 22119 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.554 ms/conn | 3568632 B/op | 22459 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 4.006 ms/conn | 4064328 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.128 ms/conn | 1325904 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.19 ms/conn | 1346544 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.745 ms/conn | 1487984 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.761 ms/conn | 1488032 B/op | 11923 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.098 ms/conn | 1564304 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.031 ms/conn | 827824 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.127 ms/conn | 1333104 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.297 ms/conn | 1582864 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.669 ms/conn | 1539600 B/op | 12666 allocs/op |
| 会话恢复握手 | 5 | 1.114 ms/conn | 2354840 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.116 ms/conn | 2522376 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.392 ms/conn | 1998280 B/op | 11245 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.532 ms/conn | 2010680 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.514 ms/conn | 1316064 B/op | 9519 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.691 ms/conn | 2470264 B/op | 11025 allocs/op |
| 完整 mTLS 握手 | 5 | 6.605 ms/conn | 1760120 B/op | 15072 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.642 ms/conn | 1962240 B/op | 15592 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.51 ms/conn | 1381216 B/op | 10439 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.993 ms/conn | 1037328 B/op | 7679 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.541 ms/conn | 2484376 B/op | 11586 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.886 ms/conn | 2646744 B/op | 12685 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 14.24 ms/conn | 3391248 B/op | 22854 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3620272 B/op | 20507 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 243960 B/op | 1019 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.891 ms/conn | 1617888 B/op | 10019 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 4.006 ms/conn | 1637728 B/op | 10159 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.407 ms/conn | 1796928 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.968 ms/conn | 1169544 B/op | 1203 allocs/op |
| 完整 mTLS 握手 | 5 | 8.973 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 9.135 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.311 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.233 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.972 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.914 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.28 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1190872 B/op | 1203 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.361 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 7.115 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 11.86 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 36.33 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 44.77 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 486.4 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.971 us/op | 2078.4 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.524 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 206.3 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 176.5 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 878.5 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.626 us/op | 1129.7 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 306.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 91.24 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 356.9 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 20.77 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 569 ns/op | 2108.82 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 533 ns/op | 2251.47 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 527 ns/op | 2276.84 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 24.341 us/op | 2692.47 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 472.1 ns/op | 2541.75 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 24.07 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.581 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.11 us/op | 568.82 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 834.1 ns/op | 1438.65 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 14.47 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.285 us/op | 525.27 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.324 us/op | 906.65 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.407 us/op | 162.02 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.091 us/op | 574 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.167 us/op | 553.71 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.124 us/op | 384.08 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.016 us/op | 1180.84 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.275 us/op | 366.39 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 853.2 ns/op | 1406.48 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 897.3 ns/op | 1337.41 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.35 us/op | 888.7 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.808 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.577 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.535 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.472 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.526 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.499 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.235 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.412 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.264 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 16.04 ns/op | 0 B/op | 0 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 16.02 ns/op | 0 B/op | 0 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.254 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.458 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 710.7 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.848 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 9.122 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 747 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 753 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.637 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.333 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.121 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.696 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.685 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.313 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.571 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.567 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.835 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.86 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 7.015 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.77 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.799 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 265.1 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 509.6 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 109.6 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 87.62 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 297.4 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 263.5 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 368.7 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 562.1 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 51.89 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 539.8 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 88.49 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 46.04 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 73.52 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 746.6 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 90.41 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 72.2 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 68.28 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 571.5 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 410.6 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.36 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 50.58 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 94.84 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 65.16 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 33.14 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 288.8 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 155.5 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 80.23 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.09 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 798 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 784.9 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 23.33 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 54.87 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.193 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.931 us/op | 4296 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
