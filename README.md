# 自动化基准测试结果

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- 提交: `30146f0ccfc956cc3f7d27a13ca1eda152f15c02`
- 生成时间: `2026-10-08T12:51:21Z`
- Go: `go version go1.27.1 linux/amd64`
- 平台: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `799bd17483efd76f30bd0f17cb080a0fc0ea92da (Linux Release static)`

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
| 证书认证完整握手 / AES-128-GCM | 5 | 536.731 us/op | 98185 B/op | 710 allocs/op |
| 完整 mTLS 握手 | 5 | 812.849 us/op | 112583 B/op | 895 allocs/op |
| mTLS 会话恢复握手 | 5 | 412.561 us/op | 119639 B/op | 823 allocs/op |
| 按 CA 与 OID filters 选择多证书的 mTLS 握手 | 5 | 1.056 ms/op | 120444 B/op | 1069 allocs/op |
| 握手后认证的多证书选择 | 5 | 1.42 ms/op | 139532 B/op | 1363 allocs/op |
| 完整握手 + 4 个已确认会话票据 | 5 | 633.43 us/op | 117760 B/op | 934 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 关闭 | 5 | 875.002 us/op | 121897 B/op | 967 allocs/op |
| 完整 mTLS 握手 + 会话票据 / GREASE 启用 | 5 | 875.273 us/op | 121897 B/op | 967 allocs/op |
| 直接外部 PSK 握手 | 5 | 347.554 us/op | 101630 B/op | 744 allocs/op |
| 服务器证书完整握手 / 证书未压缩 | 5 | 1.004 ms/op | 130290 B/op | 993 allocs/op |
| zlib 服务器证书压缩握手 | 5 | 994.883 us/op | 122507 B/op | 971 allocs/op |
| 完整 mTLS 握手 / 证书未压缩 | 5 | 1.675 ms/op | 170216 B/op | 1429 allocs/op |
| zlib mTLS 证书压缩握手 | 5 | 1.684 ms/op | 157779 B/op | 1388 allocs/op |
| ECH 握手 / 直接（无 HRR） | 5 | 928.09 us/op | 147691 B/op | 1210 allocs/op |
| ECH 握手 / 经 HRR | 5 | 931.618 us/op | 150467 B/op | 1231 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 842.436 us/op | 146164 B/op | 743 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 815.192 us/op | 149508 B/op | 774 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 2.246 ms/op | 174949 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## 真实 UDP 互通

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls 客户端 -> go-dtls 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 1.832 ms/conn | 2493840 B/op | 21254 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 1.845 ms/conn | 2502608 B/op | 21576 allocs/op |
| 完整 mTLS 握手 | 5 | 3.266 ms/conn | 3179920 B/op | 29404 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 4.553 ms/conn | 3678912 B/op | 31426 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 1.863 ms/conn | 2982960 B/op | 25185 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.4213 ms/conn | 1733984 B/op | 14633 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 1.844 ms/conn | 2513856 B/op | 22373 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 1.944 ms/conn | 2671248 B/op | 22888 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 4.906 ms/conn | 3978536 B/op | 39548 allocs/op |
| 会话恢复握手 | 5 | 0.4825 ms/conn | 4457888 B/op | 37893 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.6133 ms/conn | 6481712 B/op | 49695 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 0.5302 ms/conn | 288328 B/op | 1895 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.135 ms/conn | 3365352 B/op | 21578 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 2.131 ms/conn | 3407272 B/op | 21918 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 3.466 ms/conn | 3902952 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls 客户端 -> wolfSSL 服务端

中位耗时由 go-dtls 客户端计时；`ms/conn` 表示一次完整连接工作负载。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 4.944 ms/conn | 1191184 B/op | 10762 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 5.013 ms/conn | 1211184 B/op | 11022 allocs/op |
| 完整 mTLS 握手 | 5 | 6.696 ms/conn | 1347184 B/op | 11602 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 6.658 ms/conn | 1347184 B/op | 11602 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.867 ms/conn | 1429584 B/op | 12322 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.9189 ms/conn | 814064 B/op | 6882 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.908 ms/conn | 1198384 B/op | 11322 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 5.22 ms/conn | 1447184 B/op | 12362 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 6.713 ms/conn | 1403024 B/op | 12387 allocs/op |
| 会话恢复握手 | 5 | 0.9636 ms/conn | 2150200 B/op | 18384 allocs/op |
| mTLS 会话恢复握手 | 5 | 0.9575 ms/conn | 2311656 B/op | 19264 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 go-dtls 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 5.224 ms/conn | 1862840 B/op | 10964 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.4 ms/conn | 1875320 B/op | 11084 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | - | 跳过: MTU 4096 下 P384 ClientHello 超出 wolfSSL 服务端接收缓冲；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL 客户端 -> go-dtls 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 2.18 ms/conn | 1293024 B/op | 9299 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 4.595 ms/conn | 2447768 B/op | 10805 allocs/op |
| 完整 mTLS 握手 | 5 | 5.434 ms/conn | 1736128 B/op | 14855 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 5.444 ms/conn | 1938944 B/op | 15369 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 2.134 ms/conn | 1359104 B/op | 10219 allocs/op |
| 直接外部 PSK 握手 | 5 | 0.761 ms/conn | 1017296 B/op | 7499 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 4.641 ms/conn | 2461624 B/op | 11366 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 4.672 ms/conn | 2623544 B/op | 12445 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 11.4 ms/conn | 3365800 B/op | 22634 allocs/op |
| 会话恢复握手 | 5 | 1010 ms/pair | 3576336 B/op | 20067 allocs/op |
| mTLS 会话恢复握手 | - | 跳过: wolfSSL 客户端拒绝分片发送携带 mTLS ticket 的首个 ClientHello；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 应用数据 1-RTT 往返 | 5 | 1009 ms/pair | 241472 B/op | 997 allocs/op |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 2.404 ms/conn | 1595168 B/op | 9799 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 3.454 ms/conn | 1615488 B/op | 9939 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 6.43 ms/conn | 1774240 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL 客户端 -> wolfSSL 服务端

中位耗时由 wolfSSL 客户端计时；`ms/conn` 表示单个连接，`ms/pair` 表示由完整连接和恢复连接组成的一组，并包含客户端内置等待。

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 证书认证完整握手 / AES-128-GCM | 5 | 3.967 ms/conn | 63640 B/op | 55 allocs/op |
| 应用数据 1-RTT 往返 | 5 | 7.03 ms/conn | 1169472 B/op | 1201 allocs/op |
| 完整 mTLS 握手 | 5 | 7.328 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE 兼容性 / 完整 mTLS 握手 + 会话票据 | 5 | 7.649 ms/conn | 63608 B/op | 55 allocs/op |
| 证书认证完整握手 / AES-128-CCM | 5 | 4.057 ms/conn | 63640 B/op | 55 allocs/op |
| 直接外部 PSK 握手 | 5 | 1.129 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 应用数据 1-RTT 往返 | 5 | 7.079 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 应用数据 1-RTT 往返 | 5 | 6.889 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 应用数据 1-RTT 往返 | 5 | 10.13 ms/conn | 1172480 B/op | 1206 allocs/op |
| 会话恢复握手 | 5 | 1013 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS 会话恢复握手 | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 应用数据 1-RTT 往返 | - | 跳过: 该场景启用了 cookie HRR，按 RFC 必须拒绝 wolfSSL 客户端 0-RTT；该限制最后验证于 wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 后量子混合密钥交换 / X25519MLKEM768 | 5 | 60.16 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP256r1MLKEM768 | 5 | 6.504 ms/conn | 64392 B/op | 55 allocs/op |
| 后量子混合密钥交换 / SecP384r1MLKEM1024 | 5 | 65.42 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## 记录层与可靠性

| 基准测试 | 样本数 | 中位耗时 | 吞吐量 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: | :---: |
| 明文 ACK 构建 / 空 | 5 | 37.97 ns/op | - | 16 B/op | 1 allocs/op |
| 明文 ACK 构建 / 单条 | 5 | 48.91 ns/op | - | 32 B/op | 1 allocs/op |
| 明文 ACK 构建 / 已排序 64 | 5 | 600.3 ns/op | - | 1152 B/op | 1 allocs/op |
| 构建明文握手报文组 | 5 | 2.432 us/op | 1684.15 MB/s | 5008 B/op | 9 allocs/op |
| 受保护 ACK 构建 / 逆序 64 | 5 | 1.826 us/op | - | 2200 B/op | 3 allocs/op |
| 受保护 ACK 构建 / 单条 | 5 | 211.6 ns/op | - | 72 B/op | 2 allocs/op |
| 受保护 ACK 构建 / 单条复用 | 5 | 166.1 ns/op | - | 48 B/op | 1 allocs/op |
| 受保护 ACK 构建 / 已排序 64 | 5 | 1.003 us/op | - | 1176 B/op | 2 allocs/op |
| 构建受保护握手报文组 | 5 | 3.713 us/op | 1103.27 MB/s | 5584 B/op | 6 allocs/op |
| 合并握手报文组 | 5 | 383.2 ns/op | - | 592 B/op | 4 allocs/op |
| 握手报文组首次刷新 | 5 | 105.3 ns/op | - | 0 B/op | 0 allocs/op |
| 握手报文组初始历史批次 | 5 | 381.5 ns/op | - | 480 B/op | 1 allocs/op |
| 握手报文组传输窗口 / 重传 | 5 | 27.97 ns/op | - | 0 B/op | 0 allocs/op |
| 接收缓存 / 单分片 | 5 | 605.8 ns/op | 1980.78 MB/s | 1312 B/op | 2 allocs/op |
| 接收缓存 / 分片批次 | 5 | 569.7 ns/op | 2106.41 MB/s | 1280 B/op | 1 allocs/op |
| 接收缓存 / 分片复用 | 5 | 579.5 ns/op | 2070.58 MB/s | 1280 B/op | 1 allocs/op |
| 握手重组 | 5 | 28.401 us/op | 2307.55 MB/s | 73856 B/op | 3 allocs/op |
| 握手重组单分片 | 5 | 552.8 ns/op | 2170.89 MB/s | 1280 B/op | 1 allocs/op |
| 解析 ACK / 独占 | 5 | 27.27 ns/op | - | 16 B/op | 1 allocs/op |
| 解析 ACK / 单条复用 | 5 | 4.659 ns/op | - | 0 B/op | 0 allocs/op |
| 受保护记录 CID / 往返 | 5 | 2.799 us/op | 428.74 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录 CID / 封装 | 5 | 1.02 us/op | 1176.06 MB/s | 1280 B/op | 1 allocs/op |
| 拒绝未认证记录 | 5 | 12.31 ns/op | - | 0 B/op | 0 allocs/op |
| 记录往返 | 5 | 3.439 us/op | 348.92 MB/s | 3840 B/op | 3 allocs/op |
| 受保护记录原地往返 | 5 | 1.599 us/op | 750.49 MB/s | 1280 B/op | 1 allocs/op |
| 记录往返 / AES-128-CCM | 5 | 8.001 us/op | 149.97 MB/s | 6240 B/op | 12 allocs/op |
| 记录往返 / AES-128-GCM | 5 | 2.782 us/op | 431.31 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / AES-256-GCM | 5 | 2.91 us/op | 412.34 MB/s | 3840 B/op | 3 allocs/op |
| 记录往返 / ChaCha20-Poly1305 | 5 | 4.057 us/op | 295.82 MB/s | 3840 B/op | 3 allocs/op |
| 记录封装 | 5 | 1.248 us/op | 961.69 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-128-CCM | 5 | 3.306 us/op | 362.98 MB/s | 1840 B/op | 5 allocs/op |
| 记录封装 / AES-128-GCM | 5 | 1.155 us/op | 1039.4 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / AES-256-GCM | 5 | 1.2 us/op | 1000.37 MB/s | 1280 B/op | 1 allocs/op |
| 记录封装 / ChaCha20-Poly1305 | 5 | 1.801 us/op | 666.15 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## 密钥调度与密码学

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 计算 PSK 绑定值 / AES-128-GCM | 5 | 3.022 us/op | 1952 B/op | 21 allocs/op |
| 计算 PSK 绑定值 / AES-256-GCM | 5 | 6.936 us/op | 3248 B/op | 21 allocs/op |
| 派生流量密钥 / AES-128-GCM | 5 | 1.683 us/op | 976 B/op | 9 allocs/op |
| 派生流量密钥 / AES-256-GCM | 5 | 3.721 us/op | 1520 B/op | 9 allocs/op |
| 派生流量密钥并写入 / AES-128-GCM | 5 | 1.63 us/op | 928 B/op | 8 allocs/op |
| 派生流量密钥并写入 / AES-256-GCM | 5 | 3.656 us/op | 1440 B/op | 8 allocs/op |
| 空握手转录哈希 / AES-128-GCM | 5 | 1.154 ns/op | 0 B/op | 0 allocs/op |
| 空握手转录哈希 / AES-256-GCM | 5 | 1.153 ns/op | 0 B/op | 0 allocs/op |
| Finished 验证数据 / AES-128-GCM | 5 | 1.513 us/op | 992 B/op | 11 allocs/op |
| Finished 验证数据 / AES-256-GCM | 5 | 3.471 us/op | 1648 B/op | 11 allocs/op |
| 安装应用密钥 / AES-128-GCM | 5 | 6.904 us/op | 7488 B/op | 34 allocs/op |
| 安装应用密钥 / AES-256-GCM | 5 | 11.396 us/op | 8544 B/op | 34 allocs/op |
| 密钥调度派生 / AES-128-GCM | 5 | 7.807 us/op | 5184 B/op | 48 allocs/op |
| 密钥调度派生 / AES-256-GCM | 5 | 17.177 us/op | 8224 B/op | 48 allocs/op |
| 密钥派生 / AES-128-GCM / 早期流量 | 5 | 768.2 ns/op | 480 B/op | 5 allocs/op |
| 密钥派生 / AES-128-GCM / 导出器 | 5 | 2.033 us/op | 1408 B/op | 15 allocs/op |
| 密钥派生 / AES-128-GCM / 零值导出器 | 5 | 18.54 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-128-GCM / 恢复 PSK | 5 | 793 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-128-GCM / 流量更新 | 5 | 789.6 ns/op | 512 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 早期流量 | 5 | 1.734 us/op | 800 B/op | 5 allocs/op |
| 密钥派生 / AES-256-GCM / 导出器 | 5 | 4.479 us/op | 2384 B/op | 15 allocs/op |
| 密钥派生 / AES-256-GCM / 零值导出器 | 5 | 18.54 ns/op | 0 B/op | 0 allocs/op |
| 密钥派生 / AES-256-GCM / 恢复 PSK | 5 | 1.787 us/op | 848 B/op | 6 allocs/op |
| 密钥派生 / AES-256-GCM / 流量更新 | 5 | 1.78 us/op | 848 B/op | 6 allocs/op |
| 新建记录密码器 / AES-128-CCM | 5 | 2.538 us/op | 2520 B/op | 13 allocs/op |
| 新建记录密码器 / AES-128-GCM | 5 | 2.872 us/op | 3264 B/op | 13 allocs/op |
| 新建记录密码器 / AES-256-GCM | 5 | 4.917 us/op | 3776 B/op | 13 allocs/op |
| 新建记录密码器 / ChaCha20-Poly1305 | 5 | 2.003 us/op | 1528 B/op | 12 allocs/op |
| 接收 KeyUpdate / AES-128-GCM | 5 | 4.162 us/op | 3776 B/op | 19 allocs/op |
| 接收 KeyUpdate / AES-256-GCM | 5 | 7.347 us/op | 4624 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-128-GCM | 5 | 3.993 us/op | 3792 B/op | 19 allocs/op |
| 发送 KeyUpdate / AES-256-GCM | 5 | 6.981 us/op | 4624 B/op | 19 allocs/op |
| 握手转录克隆 / AES-128-GCM | 5 | 329.9 ns/op | 288 B/op | 4 allocs/op |
| 握手转录克隆 / AES-256-GCM | 5 | 633.5 ns/op | 496 B/op | 4 allocs/op |
| 握手转录求和 / AES-128-GCM / 独占 | 5 | 116.7 ns/op | 32 B/op | 1 allocs/op |
| 握手转录求和 / AES-128-GCM / 复用 | 5 | 77.94 ns/op | 0 B/op | 0 allocs/op |
| 握手转录求和 / AES-256-GCM / 独占 | 5 | 315.6 ns/op | 48 B/op | 1 allocs/op |
| 握手转录求和 / AES-256-GCM / 复用 | 5 | 262.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## 报文编码与解析

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 编码扩展 | 5 | 329.7 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 证书 | 5 | 626.9 ns/op | 1152 B/op | 1 allocs/op |
| 编码握手 / 证书验证 | 5 | 56.44 ns/op | 80 B/op | 1 allocs/op |
| 编码握手 / 客户端 Hello | 5 | 601.3 ns/op | 424 B/op | 8 allocs/op |
| 编码握手 / Hello 重试请求 | 5 | 101.1 ns/op | 128 B/op | 1 allocs/op |
| 编码握手 / 新连接 ID | 5 | 47.65 ns/op | 32 B/op | 1 allocs/op |
| 编码握手 / 新会话票据 | 5 | 79 ns/op | 96 B/op | 1 allocs/op |
| 编码握手 / 恢复 Client Hello | 5 | 832.5 ns/op | 744 B/op | 9 allocs/op |
| 编码握手 / 服务端 Hello | 5 | 105.5 ns/op | 112 B/op | 1 allocs/op |
| 编码握手 / 会话票据状态 | 5 | 74.67 ns/op | 80 B/op | 1 allocs/op |
| 解析扩展 / 有序视图 | 5 | 63.46 ns/op | 0 B/op | 0 allocs/op |
| 解析扩展 / 独占 | 5 | 652.3 ns/op | 472 B/op | 8 allocs/op |
| 解析扩展 / 视图 | 5 | 444.4 ns/op | 336 B/op | 2 allocs/op |
| 解析握手分片 / 单条复用 | 5 | 10.66 ns/op | 0 B/op | 0 allocs/op |
| 解析握手分片 / 视图 | 5 | 58.63 ns/op | 48 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 独占 | 5 | 96.56 ns/op | 64 B/op | 2 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 视图 | 5 | 61.19 ns/op | 32 B/op | 1 allocs/op |
| 解析密钥份额 / 1 个密钥份额 / 写入视图 | 5 | 27.05 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 独占 | 5 | 304.7 ns/op | 256 B/op | 5 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 视图 | 5 | 168.2 ns/op | 128 B/op | 1 allocs/op |
| 解析密钥份额 / 4 个密钥份额 / 写入视图 | 5 | 75.99 ns/op | 0 B/op | 0 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 独占 | 5 | 1.149 us/op | 824 B/op | 14 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 视图 | 5 | 825.9 ns/op | 536 B/op | 5 allocs/op |
| 解析密钥份额 / 9 个密钥份额 / 写入视图 | 5 | 819.1 ns/op | 536 B/op | 5 allocs/op |
| 解析明文记录 / 单条复用 | 5 | 18.46 ns/op | 0 B/op | 0 allocs/op |
| 解析明文记录 / 视图 | 5 | 64.52 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## 证书压缩

| 基准测试 | 样本数 | 中位耗时 | 测试框架内存 | 测试框架分配次数 |
| --- | :---: | :---: | :---: | :---: |
| 压缩 | 5 | 5.043 us/op | 176 B/op | 3 allocs/op |
| 解压 | 5 | 5.743 us/op | 4264 B/op | 6 allocs/op |

[Go benchmark 原始输出](benchmark.txt)
