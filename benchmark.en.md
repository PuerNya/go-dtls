# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-21T02:46:58Z`
- Go: `go version go1.27.3 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `aa3f7cb6e3c32f2b409055f75ea52de3539e5283 (Linux Release static)`
- NSS: `b9283c1e4ee232b4450d36115272af670d38ec81 (Linux Release)`
- BoringSSL: `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 (Linux Release)`
- OpenSSL: `41eecc7c7cd4069b1b29bca49b9d97825c6bf648 (Linux Release)`

378 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per operation; Listener UDP sums the separate client/server processes and excludes the relay. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (39)](#section-connection-lifecycle)
- [Real UDP interoperability (284)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (15)](#real-udp-wolfssl-client-wolfssl-server)
  - [NSS / go-dtls client -> go-dtls server (18)](#real-udp-nss-go-go)
  - [NSS / go-dtls client -> NSS server (18)](#real-udp-nss-go-peer)
  - [NSS / NSS client -> go-dtls server (18)](#real-udp-nss-peer-go)
  - [NSS / NSS client -> NSS server (18)](#real-udp-nss-peer-peer)
  - [BoringSSL / go-dtls client -> go-dtls server (19)](#real-udp-boringssl-go-go)
  - [BoringSSL / go-dtls client -> BoringSSL server (19)](#real-udp-boringssl-go-peer)
  - [BoringSSL / BoringSSL client -> go-dtls server (19)](#real-udp-boringssl-peer-go)
  - [BoringSSL / BoringSSL client -> BoringSSL server (19)](#real-udp-boringssl-peer-peer)
  - [OpenSSL / go-dtls client -> go-dtls server (19)](#real-udp-openssl-go-go)
  - [OpenSSL / go-dtls client -> OpenSSL server (19)](#real-udp-openssl-go-peer)
  - [OpenSSL / OpenSSL client -> go-dtls server (19)](#real-udp-openssl-peer-go)
  - [OpenSSL / OpenSSL client -> OpenSSL server (19)](#real-udp-openssl-peer-peer)
- [Record layer and reliability (55)](#section-record-layer-and-reliability)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 601.905 us/op | - | 101386 B/op | 733 allocs/op |
| Full mTLS handshake | 5 | 894.05 us/op | - | 116462 B/op | 931 allocs/op |
| mTLS session resumption handshake | 5 | 451.665 us/op | - | 123254 B/op | 856 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.106 ms/op | - | 123771 B/op | 1110 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.535 ms/op | - | 142635 B/op | 1407 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 706.698 us/op | - | 121415 B/op | 963 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 950.302 us/op | - | 125777 B/op | 1002 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 953.902 us/op | - | 125873 B/op | 1002 allocs/op |
| Direct external PSK handshake | 5 | 377.682 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.037 ms/op | - | 136154 B/op | 1018 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.052 ms/op | - | 126076 B/op | 994 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.748 ms/op | - | 177280 B/op | 1471 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.76 ms/op | - | 162479 B/op | 1424 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.017 ms/op | - | 151467 B/op | 1233 allocs/op |
| ECH handshake / via HRR | 5 | 1.005 ms/op | - | 154243 B/op | 1254 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 892.094 us/op | - | 149374 B/op | 766 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 833.138 us/op | - | 152719 B/op | 797 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.134 ms/op | - | 178165 B/op | 815 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 118.675 us/op | 1.08 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 123.465 us/op | 19.44 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 119.717 us/op | 1.07 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 126.073 us/op | 19.04 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 21.507 us/op | 5.95 MB/s | 420 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 20.286 us/op | 118.31 MB/s | 1546 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 20.548 us/op | 6.23 MB/s | 418 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 20.471 us/op | 117.24 MB/s | 1546 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 21.882 us/op | 5.85 MB/s | 466 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 21.87 us/op | 109.74 MB/s | 1686 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 21.407 us/op | 5.98 MB/s | 468 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 22.59 us/op | 106.24 MB/s | 1697 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 22.183 us/op | 108.19 MB/s | 1680 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.433 ms/op | - | 88680 B/op | 768 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 652.561 us/op | - | 90447 B/op | 768 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 1.145 ms/op | - | 93943 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 438.042 us/op | - | 95068 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 444.285 us/op | - | 109340 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 434.536 us/op | - | 109180 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.407 us/op | 37.57 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.753 us/op | 355.39 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.995 ms/conn | 2619416 B/op | 21138 allocs/op |
| 1-RTT application-data round trip | 5 | 2.026 ms/conn | 2631496 B/op | 21477 allocs/op |
| Full mTLS handshake | 5 | 3.591 ms/conn | 3275440 B/op | 29177 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.629 ms/conn | 3731904 B/op | 30920 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.012 ms/conn | 3097256 B/op | 24988 allocs/op |
| Direct external PSK handshake | 5 | 0.5546 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.162 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.237 ms/conn | 4056560 B/op | 39430 allocs/op |
| Session resumption handshake | 5 | 0.6799 ms/conn | 4651288 B/op | 38434 allocs/op |
| mTLS session resumption handshake | 5 | 0.8239 ms/conn | 6674256 B/op | 49992 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6801 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.368 ms/conn | 3432256 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.323 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.627 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.708 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.813 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.382 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.427 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.607 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.9992 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.817 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.015 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.271 ms/conn | 1547216 B/op | 12665 allocs/op |
| Session resumption handshake | 5 | 1.195 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.181 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.956 ms/conn | 2006104 B/op | 11253 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.04 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.394 ms/conn | 1304160 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 5.232 ms/conn | 2458264 B/op | 10765 allocs/op |
| Full mTLS handshake | 5 | 6.122 ms/conn | 1704528 B/op | 14601 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.134 ms/conn | 1907040 B/op | 15121 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.427 ms/conn | 1368704 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.977 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.322 ms/conn | 2473224 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.418 ms/conn | 2633784 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.35 ms/conn | 3311152 B/op | 22372 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3598352 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.594 ms/conn | 1572752 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.698 ms/conn | 1592592 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.849 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.829 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.074 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.039 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.81 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.012 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.089 ms/conn | 519144 B/op | 3441 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.004 ms/conn | 390232 B/op | 3404 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9479 ms/conn | 459688 B/op | 3489 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.11 ms/conn | 464568 B/op | 3560 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.014 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.251 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.094 ms/conn | 507360 B/op | 3819 allocs/op |
| Full mTLS handshake | 5 | 1.443 ms/conn | 537736 B/op | 4564 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.026 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.083 ms/conn | 487752 B/op | 3668 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.221 ms/conn | 533400 B/op | 4212 allocs/op |
| Direct external PSK handshake | 5 | 0.6437 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.696 ms/pair | 869144 B/op | 7185 allocs/op |
| mTLS session resumption handshake | 5 | 2.168 ms/pair | 1048200 B/op | 8543 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.695 ms/pair | 868736 B/op | 7134 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.359 ms/conn | 511704 B/op | 3992 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.168 ms/pair | 948152 B/op | 7915 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.43 ms/conn | 182128 B/op | 1538 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.335 ms/conn | 215920 B/op | 1538 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.422 ms/conn | 147280 B/op | 1518 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.879 ms/conn | 185136 B/op | 1562 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.37 ms/conn | 187216 B/op | 1598 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.08 ms/conn | 215920 B/op | 1814 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.446 ms/conn | 275440 B/op | 1650 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.12 ms/conn | 223248 B/op | 1914 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.313 ms/conn | 182416 B/op | 1550 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.275 ms/conn | 213904 B/op | 1722 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.281 ms/conn | 195120 B/op | 1646 allocs/op |
| Direct external PSK handshake | 5 | 0.8094 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.452 ms/pair | 379920 B/op | 3145 allocs/op |
| mTLS session resumption handshake | 5 | 3.562 ms/pair | 459696 B/op | 3521 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.572 ms/conn | 203760 B/op | 1718 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.745 ms/pair | 421232 B/op | 3481 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.325 ms/conn | 481240 B/op | 1939 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.424 ms/conn | 512840 B/op | 1939 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.33 ms/conn | 453160 B/op | 1921 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.89 ms/conn | 484008 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.326 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.01 ms/conn | 514152 B/op | 2135 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.467 ms/conn | 538456 B/op | 2037 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.332 ms/conn | 514504 B/op | 2123 allocs/op |
| Full mTLS handshake | 5 | 4.908 ms/conn | 524008 B/op | 2514 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.328 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.443 ms/conn | 516232 B/op | 2269 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.331 ms/conn | 505992 B/op | 2153 allocs/op |
| Direct external PSK handshake | 5 | 3.446 ms/conn | 460104 B/op | 1839 allocs/op |
| Session resumption handshake | 5 | 5.444 ms/pair | 703560 B/op | 3824 allocs/op |
| mTLS session resumption handshake | 5 | 5.878 ms/pair | 803336 B/op | 4644 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.436 ms/pair | 706000 B/op | 3786 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.225 ms/conn | 517192 B/op | 2297 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.385 ms/pair | 740456 B/op | 4194 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.501 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.591 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.542 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.702 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.273 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.903 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.491 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.57 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.479 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.547 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.472 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.535 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.596 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.67 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.222 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.203 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.01 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.051 ms/conn | 519096 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.019 ms/conn | 390232 B/op | 3404 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.965 ms/conn | 459640 B/op | 3488 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.16 ms/conn | 464440 B/op | 3560 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.072 ms/conn | 527872 B/op | 3899 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.252 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.034 ms/conn | 507488 B/op | 3819 allocs/op |
| Full mTLS handshake | 5 | 1.499 ms/conn | 537736 B/op | 4564 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.003 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.041 ms/conn | 485384 B/op | 3660 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.231 ms/conn | 533528 B/op | 4212 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8528 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6597 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.767 ms/pair | 868424 B/op | 7183 allocs/op |
| mTLS session resumption handshake | 5 | 2.172 ms/pair | 1047944 B/op | 8543 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.37 ms/conn | 650520 B/op | 5444 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.368 ms/conn | 511832 B/op | 3992 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.165 ms/pair | 948024 B/op | 7915 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8435 ms/conn | 180576 B/op | 1528 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8211 ms/conn | 213728 B/op | 1528 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7905 ms/conn | 145504 B/op | 1507 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.841 ms/conn | 183584 B/op | 1552 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.585 ms/conn | 185536 B/op | 1588 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9588 ms/conn | 204256 B/op | 1760 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.096 ms/conn | 271200 B/op | 1596 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7903 ms/conn | 180896 B/op | 1560 allocs/op |
| Full mTLS handshake | 5 | 1.131 ms/conn | 220576 B/op | 1904 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.7693 ms/conn | 180768 B/op | 1540 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.9356 ms/conn | 212352 B/op | 1712 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.071 ms/conn | 229344 B/op | 2076 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7582 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6364 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.614 ms/pair | 388768 B/op | 3283 allocs/op |
| mTLS session resumption handshake | 5 | 1.904 ms/pair | 474336 B/op | 3659 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.041 ms/conn | 241536 B/op | 2256 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.9783 ms/conn | 201504 B/op | 1708 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.82 ms/pair | 432384 B/op | 3636 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.54 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.773 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.694 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.625 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.855 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.629 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.772 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.757 ms/conn | 509800 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 3.174 ms/conn | 516808 B/op | 2345 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.119 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.07 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.909 ms/conn | 501640 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.6 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.477 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 3.512 ms/pair | 693064 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 4.397 ms/pair | 789752 B/op | 4378 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.972 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.237 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.39 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.395 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.428 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.577 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.34 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.376 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.608 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.72 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.491 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.394 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.56 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.441 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.236 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 3.124 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 3.081 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.395 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.041 ms/conn | 452856 B/op | 3440 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.089 ms/conn | 519096 B/op | 3440 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.053 ms/conn | 390360 B/op | 3404 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.073 ms/conn | 517384 B/op | 4126 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.091 ms/conn | 459512 B/op | 3488 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.112 ms/conn | 464568 B/op | 3560 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.046 ms/conn | 527744 B/op | 3899 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.302 ms/conn | 622648 B/op | 3590 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.059 ms/conn | 507360 B/op | 3819 allocs/op |
| Full mTLS handshake | 5 | 1.479 ms/conn | 537864 B/op | 4564 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.062 ms/conn | 454856 B/op | 3488 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.066 ms/conn | 485384 B/op | 3660 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.207 ms/conn | 533528 B/op | 4212 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8378 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.6488 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.686 ms/pair | 868552 B/op | 7183 allocs/op |
| mTLS session resumption handshake | 5 | 2.178 ms/pair | 1047944 B/op | 8543 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.739 ms/pair | 868736 B/op | 7134 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.906 ms/conn | 642120 B/op | 5982 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.191 ms/conn | 211424 B/op | 1780 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.398 ms/conn | 245472 B/op | 1780 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.154 ms/conn | 176480 B/op | 1759 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.191 ms/conn | 264928 B/op | 2320 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.2 ms/conn | 220640 B/op | 1848 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.676 ms/conn | 220544 B/op | 1864 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.296 ms/conn | 223968 B/op | 1896 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.491 ms/conn | 332128 B/op | 1884 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.216 ms/conn | 212000 B/op | 1792 allocs/op |
| Full mTLS handshake | 5 | 1.52 ms/conn | 267264 B/op | 2244 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.303 ms/conn | 214944 B/op | 1816 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.31 ms/conn | 247536 B/op | 1991 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.542 ms/conn | 270624 B/op | 2424 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.129 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.9883 ms/conn | 199184 B/op | 1679 allocs/op |
| Session resumption handshake | 5 | 2.216 ms/pair | 471488 B/op | 3883 allocs/op |
| mTLS session resumption handshake | 5 | 3.079 ms/pair | 590864 B/op | 4442 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.293 ms/conn | 270528 B/op | 2308 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 6.954 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 6.805 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 6.489 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 6.812 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.565 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.547 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.201 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.457 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.352 ms/conn | 514120 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 6.39 ms/conn | 521576 B/op | 2385 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 6.888 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.857 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 7.128 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.945 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 5.585 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 7.239 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 7.594 ms/pair | 801704 B/op | 4514 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.192 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.36 ms/conn | 574856 B/op | 3093 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 6.567 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 6.742 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 6.7 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 6.501 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.569 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.994 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.386 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.591 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 6.367 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 6.882 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 6.802 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.877 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 7.021 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.181 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 5.78 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.836 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 7.826 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.838 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 16.59 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 15.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 17.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 15.94 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 17.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 11.85 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 28.31 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 29.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 22.59 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 13.33 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 35.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 33.12 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.768 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 3.434 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 40.79 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 18 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 2.121 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 4.13 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 66.79 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 19.66 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 237.9 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 199 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 205.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.532 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 200 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 359.5 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 36.72 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.88 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 518.7 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.04 us/op | 2007.42 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.656 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 210.7 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 167.2 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 926.1 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.511 us/op | 1166.57 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 314.2 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.9 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 338.8 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.5 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.03 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.367 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.182 us/op | 550.05 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 863.3 ns/op | 1389.99 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.62 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.215 us/op | 541.7 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.311 us/op | 915.12 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.958 us/op | 172.47 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.15 us/op | 558.17 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.244 us/op | 534.81 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.98 us/op | 402.7 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 938.7 ns/op | 1278.3 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.005 us/op | 399.38 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 854.6 ns/op | 1404.11 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 903.1 ns/op | 1328.79 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.273 us/op | 942.72 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
