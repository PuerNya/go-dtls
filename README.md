# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- 生成时间: `2026-10-07T21:20:11Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `0656a210be0269fa08a9a15874931fedd32eab0f (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 594.567 us/op | 97993 B/op | 710 allocs/op |
| 完整 mTLS 握手 | 5 | 888.294 us/op | 112214 B/op | 894 allocs/op |
| mTLS 会话恢复握手 | 5 | 429.41 us/op | 119443 B/op | 823 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.157 ms/op | 119991 B/op | 1068 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.55 ms/op | 138947 B/op | 1362 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 690.613 us/op | 117569 B/op | 934 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 954.081 us/op | 121529 B/op | 966 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 951.755 us/op | 121529 B/op | 966 allocs/op |
| 直接外部 PSK 握手 | 5 | 368.096 us/op | 101437 B/op | 744 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.076 ms/op | 130099 B/op | 993 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.077 ms/op | 122316 B/op | 971 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.816 ms/op | 169848 B/op | 1428 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.819 ms/op | 157395 B/op | 1387 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.02 ms/op | 147498 B/op | 1210 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.022 ms/op | 150274 B/op | 1231 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 926.312 us/op | 145973 B/op | 743 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 869.604 us/op | 149317 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.324 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.093 ms/conn | 2487488 B/op | 21255 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.105 ms/conn | 2496208 B/op | 21576 allocs/op |
| 完整 mTLS 握手 | 5 | 3.667 ms/conn | 3169928 B/op | 29390 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.928 ms/conn | 3666400 B/op | 31392 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.102 ms/conn | 2976560 B/op | 25185 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5266 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.08 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.2 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.514 ms/conn | 3964168 B/op | 39514 allocs/op |
| 会话恢复握手 | 5 | 0.6104 ms/conn | 4443032 B/op | 37899 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6562 ms/conn | 6462176 B/op | 49660 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6215 ms/conn | 287904 B/op | 1898 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.49 ms/conn | 3358952 B/op | 21578 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.46 ms/conn | 3400872 B/op | 21918 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.835 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.027 ms/conn | 1187344 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.125 ms/conn | 1207344 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 6.57 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.623 ms/conn | 1343344 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.996 ms/conn | 1425744 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9836 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.967 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.22 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.601 ms/conn | 1398928 B/op | 12383 allocs/op |
| 会话恢复握手 | 5 | 1.024 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.04 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.278 ms/conn | 1859000 B/op | 10964 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.363 ms/conn | 1871480 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.477 ms/conn | 1291136 B/op | 9299 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.394 ms/conn | 2445144 B/op | 10805 allocs/op |
| 完整 mTLS 握手 | 5 | 6.496 ms/conn | 1730392 B/op | 14829 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.508 ms/conn | 1932616 B/op | 15352 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.464 ms/conn | 1355680 B/op | 10219 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.958 ms/conn | 1014288 B/op | 7499 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.376 ms/conn | 2458728 B/op | 11366 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.431 ms/conn | 2620344 B/op | 12445 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.59 ms/conn | 3354712 B/op | 22597 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3571472 B/op | 20067 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 241256 B/op | 994 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.725 ms/conn | 1592768 B/op | 9799 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.836 ms/conn | 1612928 B/op | 9939 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.249 ms/conn | 1771552 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.93 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.523 ms/conn | 1169496 B/op | 1202 allocs/op |
| 完整 mTLS 握手 | 5 | 9.326 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 8.869 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.263 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.309 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.527 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.528 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.42 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.014 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 7.389 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 10.6 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 38 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 48.27 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 596.8 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.216 us/op | 1847.96 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.839 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 214.8 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 184.3 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.008 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.848 us/op | 1064.4 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 351.7 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 95.37 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 379.6 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 20.79 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 581.8 ns/op | 2062.6 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 549.1 ns/op | 2185.36 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 551.5 ns/op | 2176 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 25.361 us/op | 2584.13 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 500.4 ns/op | 2397.93 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 26.4 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.579 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 3.161 us/op | 379.67 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.234 us/op | 972.12 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 15.15 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.75 us/op | 320.03 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.905 us/op | 630 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 9.598 us/op | 125.02 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 3.156 us/op | 380.27 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 3.316 us/op | 361.84 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 4.146 us/op | 289.45 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.399 us/op | 857.91 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.889 us/op | 308.57 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.467 us/op | 817.79 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.518 us/op | 790.66 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.945 us/op | 616.92 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.787 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.491 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.601 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.591 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.52 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.468 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.407 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.405 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.233 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.333 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 10.388 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.135 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.325 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 708.3 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.822 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 9.085 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 725.1 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 727.3 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.609 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.223 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.091 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.648 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.646 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.306 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.629 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.622 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.872 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.804 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.794 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.673 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.684 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 327.1 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 655.9 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 121.2 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 87.66 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 308.7 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 263.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 369.6 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 582.5 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 53.7 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 542.4 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 91.69 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 46.43 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 75.44 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 742.6 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 91.74 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 72.58 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 66.87 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 597.5 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 436.6 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.37 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 55.48 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 94.72 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 66.1 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 33.64 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 288.1 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 162 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 81.55 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.114 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 813.9 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 829.4 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 23.35 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 58.69 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.308 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 5.843 us/op | 4264 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
