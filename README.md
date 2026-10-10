# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `10300127dca519bf9c8f7672c1119d1baaa0893c`
- 生成时间: `2026-10-10T03:22:53Z`
- Go: `go version go1.27.2 linux/amd64`
- 平台: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `7499fc5b6c99eb39055a59f538bb3709971341c9 (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 606.054 us/op | 101977 B/op | 741 allocs/op |
| 完整 mTLS 握手 | 5 | 910.285 us/op | 116990 B/op | 938 allocs/op |
| mTLS 会话恢复握手 | 5 | 449.353 us/op | 123819 B/op | 864 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.144 ms/op | 124299 B/op | 1117 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.567 ms/op | 143077 B/op | 1413 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 707.838 us/op | 122021 B/op | 971 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 971.487 us/op | 126305 B/op | 1009 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 971.657 us/op | 126401 B/op | 1009 allocs/op |
| 直接外部 PSK 握手 | 5 | 378.445 us/op | 104397 B/op | 767 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.084 ms/op | 136683 B/op | 1026 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 1.087 ms/op | 126701 B/op | 1002 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.809 ms/op | 177745 B/op | 1478 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.813 ms/op | 162975 B/op | 1431 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 1.043 ms/op | 152058 B/op | 1241 allocs/op |
| ECH 握手 / 经 HRR | 5 | 1.043 ms/op | 154835 B/op | 1262 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 912.322 us/op | 149966 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 861.686 us/op | 153312 B/op | 805 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.306 ms/op | 178759 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.146 ms/conn | 2672528 B/op | 21796 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 2.175 ms/conn | 2683408 B/op | 22136 allocs/op |
| 完整 mTLS 握手 | 5 | 3.881 ms/conn | 3375512 B/op | 30121 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 3.876 ms/conn | 3834672 B/op | 31876 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.282 ms/conn | 3239088 B/op | 25831 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.5498 ms/conn | 1786544 B/op | 15014 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 2.143 ms/conn | 2692496 B/op | 22914 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 2.266 ms/conn | 2852368 B/op | 23468 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 5.535 ms/conn | 4182552 B/op | 40426 allocs/op |
| 会话恢复握手 | 5 | 0.6609 ms/conn | 4738032 B/op | 39510 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.8238 ms/conn | 6871104 B/op | 51505 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.6858 ms/conn | 299752 B/op | 1975 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.57 ms/conn | 3543992 B/op | 22119 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.524 ms/conn | 3585912 B/op | 22459 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.923 ms/conn | 4081608 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.056 ms/conn | 1332944 B/op | 11042 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.076 ms/conn | 1353584 B/op | 11302 allocs/op |
| 完整 mTLS 握手 | 5 | 6.665 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.67 ms/conn | 1495024 B/op | 11922 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 5.047 ms/conn | 1571344 B/op | 12602 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.001 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.043 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.219 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.574 ms/conn | 1546448 B/op | 12663 allocs/op |
| 会话恢复握手 | 5 | 1.094 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS 会话恢复握手 | 5 | 1.11 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.242 ms/conn | 2005320 B/op | 11245 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.447 ms/conn | 2017720 B/op | 11364 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.498 ms/conn | 1326800 B/op | 9519 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.395 ms/conn | 2479672 B/op | 11025 allocs/op |
| 完整 mTLS 握手 | 5 | 6.543 ms/conn | 1771928 B/op | 15078 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.57 ms/conn | 1974256 B/op | 15589 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.507 ms/conn | 1391232 B/op | 10439 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.915 ms/conn | 1048848 B/op | 7679 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 5.401 ms/conn | 2495000 B/op | 11586 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.465 ms/conn | 2656024 B/op | 12685 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 13.71 ms/conn | 3403088 B/op | 22857 allocs/op |
| 会话恢复握手 | 5 | 1009 ms/pair | 3641872 B/op | 20507 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 244600 B/op | 1019 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.8 ms/conn | 1627968 B/op | 10019 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.843 ms/conn | 1648608 B/op | 10159 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 7.415 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 5.273 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 8.669 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 8.881 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 9.326 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.938 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.281 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 8.471 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 8.347 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 12.36 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1017 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.109 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.7 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 10.61 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 35.41 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 44.18 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 477.5 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 1.945 us/op | 2105.38 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.527 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 200 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 174.4 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 866.6 ns/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.483 us/op | 1175.85 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 297.1 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 87.78 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 343.2 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 22.58 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 497.1 ns/op | 2413.86 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 462.4 ns/op | 2595.18 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 472.2 ns/op | 2541.04 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 22.456 us/op | 2918.39 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 408.3 ns/op | 2939.37 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 23.59 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.928 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.051 us/op | 585.06 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 825.7 ns/op | 1453.31 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 14.44 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 2.176 us/op | 551.55 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.276 us/op | 940.38 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 7.208 us/op | 166.48 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.022 us/op | 593.42 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.123 us/op | 565.12 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 3.034 us/op | 395.54 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 917.4 ns/op | 1308.06 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.197 us/op | 375.3 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 821.7 ns/op | 1460.37 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 869.9 ns/op | 1379.51 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.316 us/op | 911.9 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 2.988 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.326 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.501 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.379 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.472 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.386 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.409 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.394 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.218 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 16.04 ns/op | 0 B/op | 0 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 16.05 ns/op | 0 B/op | 0 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.118 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 16.016 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 692 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 1.81 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 9.127 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 743.4 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 722.7 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.585 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.095 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 9.162 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.617 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.611 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.148 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.456 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.389 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 1.782 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 3.629 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 6.578 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.465 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.441 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 254.3 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 508 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 107 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 87.57 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 291.3 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 264.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 354.4 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 539.6 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 49.79 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 521.3 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 85.38 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 44.35 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 69.95 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 719.5 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 86.01 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 70 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 67.27 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 619.6 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 402.8 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 49.84 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 101.2 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 62.55 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 33.5 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 268.9 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 152 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 80.48 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.048 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 769.8 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 776.3 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 23.31 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 53.43 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.166 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 6.318 us/op | 4296 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
