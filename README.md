# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- 生成时间: `2026-10-07T05:04:55Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `ab4a14619a4eea21de84603f5e4f581b6296b3b5 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 551.426 us/op | 97993 B/op | 710 allocs/op |
| 完整 mTLS 握手 | 5 | 855.542 us/op | 112214 B/op | 894 allocs/op |
| mTLS 会话恢复握手 | 5 | 427.607 us/op | 119411 B/op | 823 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.076 ms/op | 119968 B/op | 1067 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.448 ms/op | 138902 B/op | 1361 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 643.856 us/op | 117568 B/op | 934 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 883.894 us/op | 121529 B/op | 966 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 883.92 us/op | 121529 B/op | 966 allocs/op |
| 直接外部 PSK 握手 | 5 | 376.59 us/op | 101438 B/op | 744 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.073 ms/op | 130098 B/op | 993 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.059 ms/op | 122315 B/op | 971 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.754 ms/op | 169846 B/op | 1428 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.749 ms/op | 157392 B/op | 1387 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 984.343 us/op | 147499 B/op | 1210 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.002 ms/op | 150276 B/op | 1231 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 903.97 us/op | 145972 B/op | 743 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 870.754 us/op | 149316 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.271 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.839 ms/conn | 2487536 B/op | 21256 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.863 ms/conn | 2496208 B/op | 21576 allocs/op |
| 完整 mTLS 握手 | 5 | 3.28 ms/conn | 3170272 B/op | 29390 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.581 ms/conn | 3668136 B/op | 31406 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.877 ms/conn | 2976560 B/op | 25185 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4283 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.839 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.955 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.924 ms/conn | 3964168 B/op | 39510 allocs/op |
| 会话恢复握手 | 5 | 0.5558 ms/conn | 4443344 B/op | 37896 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.5617 ms/conn | 6527384 B/op | 49652 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.5266 ms/conn | 287656 B/op | 1895 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.176 ms/conn | 3358952 B/op | 21578 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.157 ms/conn | 3400872 B/op | 21918 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.529 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.849 ms/conn | 1187344 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.059 ms/conn | 1207344 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 6.676 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.733 ms/conn | 1343344 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.935 ms/conn | 1425744 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9358 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.05 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.062 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.707 ms/conn | 1398992 B/op | 12384 allocs/op |
| 会话恢复握手 | 5 | 0.9579 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.9522 ms/conn | 2302680 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.356 ms/conn | 1859000 B/op | 10964 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.454 ms/conn | 1871480 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.177 ms/conn | 1291360 B/op | 9299 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.77 ms/conn | 2444824 B/op | 10805 allocs/op |
| 完整 mTLS 握手 | 5 | 5.544 ms/conn | 1729400 B/op | 14826 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.538 ms/conn | 1932512 B/op | 15349 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.166 ms/conn | 1356576 B/op | 10219 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.75 ms/conn | 1013840 B/op | 7499 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.793 ms/conn | 2459064 B/op | 11366 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.74 ms/conn | 2620024 B/op | 12445 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.77 ms/conn | 3354928 B/op | 22598 allocs/op |
| 会话恢复握手 | 5 | 1010 ms/pair | 3571248 B/op | 20067 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1010 ms/pair | 241184 B/op | 994 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.364 ms/conn | 1592928 B/op | 9799 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.513 ms/conn | 1612768 B/op | 9939 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.525 ms/conn | 1771936 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.39 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.351 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 7.805 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 63.67 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.381 ms/conn | 63664 B/op | 56 allocs/op |
| 直接外部 PSK 握手 | 5 | 57.67 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.449 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 7.456 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 10.46 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1190848 B/op | 1202 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 4.409 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.991 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 9.137 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 45.47 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 63.04 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 911.8 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 3.2 us/op | 1279.91 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 2.437 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 235.9 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 182.5 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.334 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 4.825 us/op | 848.91 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 519.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 100.9 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 492.3 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 28.01 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 715.6 ns/op | 1676.85 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 732.7 ns/op | 1637.72 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 700.3 ns/op | 1713.48 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 33.244 us/op | 1971.35 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 658.6 ns/op | 1822.12 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 34.28 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.349 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 4.007 us/op | 299.46 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.429 us/op | 839.79 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 12.17 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 4.086 us/op | 293.66 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.84 us/op | 652.29 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.99 us/op | 120.12 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 4.059 us/op | 295.61 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 4.204 us/op | 285.44 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 5.273 us/op | 227.56 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.388 us/op | 864.82 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.782 us/op | 317.32 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.466 us/op | 818.42 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.527 us/op | 785.71 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 2.092 us/op | 573.72 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 3.552 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 8.095 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 2.015 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 4.259 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.916 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 4.18 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.786 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 4.062 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 7.734 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 13.048 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 9.424 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 20.486 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 894.2 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 2.28 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 18.58 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 903.2 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 902.6 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.955 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 5.089 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 18.58 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 2.006 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.98 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 3.082 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 3.453 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 5.99 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 2.453 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 4.711 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 8.183 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 4.334 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 7.968 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 425.7 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 798.5 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 130.9 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 78.08 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 335.6 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 262.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 359.2 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 671 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 59.4 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 643.8 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 109.1 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 58.2 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 90.27 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 910 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 119.2 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 85.07 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 67.64 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 808.9 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 543.5 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 10.68 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 78.96 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 122.7 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 73.86 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 26.69 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 393 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 204.4 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 76.88 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.467 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 1.017 us/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 1.013 us/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 18.51 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 82.37 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.085 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.626 us/op | 4264 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
