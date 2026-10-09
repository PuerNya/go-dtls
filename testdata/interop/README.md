# DTLS 1.3 对端测试程序

[简体中文](README.md) | [English](README.en.md) | [Русский](README.ru.md)

这些按需启用的测试使用未经修改的第三方库的原生 UDP 记录层。
`openssl.c` 用于 OpenSSL 和 BoringSSL，`nss.c` 用于 NSS。Go 测试创建
短期凭据、固定信任材料、比较 exporter，并检查应用数据和协商状态。
未设置某个对端的环境变量时，该对端的测试会被跳过。

## 源码版本

互通矩阵使用以下确切的上游版本完成验证：

| 库 | 仓库 | Commit |
| --- | --- | --- |
| NSS | https://github.com/nss-dev/nss | `b9283c1e4ee232b4450d36115272af670d38ec81` |
| BoringSSL | https://github.com/google/boringssl | `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978` |
| OpenSSL | https://github.com/openssl/openssl | `41eecc7c7cd4069b1b29bca49b9d97825c6bf648` |

三个对端均协商 RFC 9147 最终版 DTLS 1.3，版本号为 `0xfefc`。
仅 TLS 1.3 测试成功或构建成功，不能证明 DTLS 互通。

## Linux 构建与运行

安装原生 C/C++ 工具链、Git、Make、CMake、Perl、Python、pkg-config 和
NSPR 开发文件。源码及构建产物放在本仓库之外。从 go-dtls 仓库目录开始，
设置 `repo` 和指向仓库外目录的绝对路径 `work`：

```sh
repo="$PWD"
work="$(mktemp -d)"
git clone https://github.com/nss-dev/nss "$work/nss"
git -C "$work/nss" checkout b9283c1e4ee232b4450d36115272af670d38ec81
git clone https://github.com/google/boringssl "$work/bssl"
git -C "$work/bssl" checkout 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978
git clone https://github.com/openssl/openssl "$work/openssl"
git -C "$work/openssl" checkout 41eecc7c7cd4069b1b29bca49b9d97825c6bf648

(cd "$work/openssl" && ./Configure no-tests && make -j2)
cmake -S "$work/bssl" -B "$work/bssl-build" -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTING=OFF
cmake --build "$work/bssl-build" -j2
make -C "$work/nss" -j2 all USE_64=1 BUILD_OPT=1 NSS_DISABLE_GTESTS=1 \
  NSPR_INCLUDE_DIR="$(pkg-config --variable=includedir nspr)" \
  NSPR_LIB_DIR="$(pkg-config --variable=libdir nspr)"
make -C "$work/nss" latest USE_64=1 BUILD_OPT=1

cc -O2 -Wall -Wextra -Werror -I"$work/openssl/include" \
  "$repo/testdata/interop/openssl.c" -L"$work/openssl" \
  -Wl,-rpath,"$work/openssl" -lssl -lcrypto -o "$work/openssl-peer"
cc -O2 -Wall -Wextra -Werror -I"$work/bssl/include" \
  "$repo/testdata/interop/openssl.c" "$work/bssl-build/libssl.a" \
  "$work/bssl-build/libcrypto.a" -lstdc++ -pthread -o "$work/bssl-peer"
nsslib="$(find "$work/dist" -name libnss3.so -printf '%h\n' | head -n1)"
cc -O2 -Wall -Wextra -Werror -I"$work/dist/public/nss" \
  $(pkg-config --cflags nspr) "$repo/testdata/interop/nss.c" \
  -L"$nsslib" -Wl,-rpath,"$nsslib" -lssl3 -lnss3 -lnssutil3 -lsmime3 \
  $(pkg-config --libs nspr) -o "$work/nss-peer"

export DTLS13_NSS_PEER="$work/nss-peer"
export DTLS13_BORINGSSL_PEER="$work/bssl-peer"
export DTLS13_OPENSSL_PEER="$work/openssl-peer"
export LD_LIBRARY_PATH="$nsslib:$work/openssl"
cd "$repo"
go test . -run '^TestInteropDTLS13Peers$' -count=1 -timeout=3m -v
```

运行测试前先完成下载，并清除测试进程继承的代理设置。测试本身不会下载或
编译对端。已有 wolfSSL 测试仍使用 `TestInteropWolfSSL*`；其构建方式见
[workflow](../../.github/workflows/benchmarks.yml)，不依赖是否运行 benchmark。

## 覆盖范围与限制

矩阵覆盖对端作为客户端和服务端的两种角色：X.509、会话恢复、
AES-GCM/SHA-256 和 SHA-384、ChaCha20-Poly1305、X25519/P-256/P-384、
key-share HRR、分片、mTLS、OCSP、SNI、ALPN、GREASE、KeyUpdate，
以及证书压缩和缓存信息的回退行为。委托凭据场景检查实际 DC 状态、恢复、
分片、SHA-384、回退、无效签名和过期。负向探测单独命名，不算成功握手。

NSS 和 OpenSSL 的 `client/early-io` 场景扣留客户端 Finished，要求释放前
已收到 epoch-1 请求并发出 epoch-3 响应。NSS 首次发送使用非阻塞模式，避免
socket API 在写入请求前完成握手。普通 early-data 场景也检查
`DatagramInfo.EarlyData`。wolfSSL RPK 矩阵在两个角色覆盖 `early-io`，使用
未修改的上游 `23f245c02c4319bb80baf7f3ae0c1de40e3fb68e`，分别经 Windows
MSVC 和 Linux 原生工具链构建；客户端正确处理连接期间的 `APP_DATA_READY`。
BoringSSL 不支持原生 DTLS 0-RTT，以及下述服务端限制仍然适用。
wolfSSL 服务端的 HRR 用例只验证拒绝 0-RTT 提议，不发送 epoch-1 数据：
此前验证的 `c6e4286c0269a111799ff059badd94a272e37574` 在 Go 客户端报告
HRR 已拒绝早期数据后，仍会回显该数据；当前用例不重新探测该 payload 行为。

| 对端 | 额外验证的行为 | 固定版本下的能力边界 |
| --- | --- | --- |
| NSS | 服务端发送 DC；客户端验证 Go 服务端的 DC；客户端 0-RTT，包括经 DC 认证的会话 | 拒绝 CertificateRequest 中的 DC（`tls13con.c`、`KnownExtensions`），违反 RFC 9345 §4.1.2。普通 mTLS 场景关闭 Go 端的这项扩展声明。原生 DTLS 不支持 PHA。 |
| BoringSSL | 服务端发送 DC | 不支持客户端 DC 验证、PHA 和 DTLS 0-RTT。`tls13_client.cc` 明确排除 DTLS 0-RTT。 |
| OpenSSL | 双向 PHA；客户端 0-RTT | 不支持 DC。其服务端以 `unexpected_message` 拒绝合法的初始 ACK；服务端正向场景有意丢弃 epoch-2 ACK，并命名为 `server-ACKLoss`。`server/initial-ACK-rejection` 保留未修改报文时的失败探测。 |

NSS 客户端 0-RTT 设置 `SSL_SUPPRESS_END_OF_EARLY_DATA`：RFC 9147 §5
省略 EndOfEarlyData；未设置该选项时，NSS 会发送这一被禁止的消息。
NSS 文档说明其服务端需要 `SSL_RecordLayerData` 才能抑制 EndOfEarlyData；
即使设置该选项，使用原生记录层的服务端仍会停滞，因此不纳入正向 0-RTT
场景。OpenSSL 的 listener 在返回连接前已经消费 ClientHello，随后调用
`SSL_read_early_data` 会因 `SSL_in_before` 为 false 而遭到拒绝。
这里也未验证其服务端 0-RTT 路径。

NSS 自动发送 DTLS 票据的功能正常，并已纳入测试；独立的
`SSL_SendSessionTicket` API 拒绝 DTLS，不代表不支持会话恢复。
NSS 使用每次测试固定的 DER 证书进行精确匹配认证；OpenSSL 和 BoringSSL
验证测试 CA 和主机名。这些测试程序均未关闭对端认证。

握手元数据扩展（`0xff02`，版本 1）是私有约定，不是这些对端的标准功能。
`metadata-unoffered` 验证启用元数据支持的 Go 服务端仍接受普通客户端；
`metadata-ech-rejected` 验证对端未配置 ECH 时客户端收到经过认证的 ECH 拒绝。
两者都不代表成功交换加密元数据；后者在 OpenSSL 上沿用 `server-ACKLoss` 夹具。
wolfSSL client-without-metadata 测试同样确认未提供扩展时不会调用接受回调。
完整/恢复连接的元数据交换、分片、重传、保密性和 UDP 由自互通覆盖。
