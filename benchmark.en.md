# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-19T02:35:16Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `7be1181fa997901e3c07fe2a56e151d6d333e054 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 588.855 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 880.862 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 440.453 us/op | - | 123190 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.096 ms/op | - | 123642 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.505 ms/op | - | 142423 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 689.63 us/op | - | 121348 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 938.636 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 936.126 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 370.315 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.044 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.055 ms/op | - | 126014 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.746 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.751 ms/op | - | 162352 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 993.882 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 993.269 us/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 896.724 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 839.335 us/op | - | 152656 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.151 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 136.274 us/op | 0.94 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 152.415 us/op | 15.75 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 137.579 us/op | 0.93 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 155.402 us/op | 15.44 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 21.951 us/op | 5.83 MB/s | 418 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 21.796 us/op | 110.11 MB/s | 1544 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 20.441 us/op | 6.26 MB/s | 421 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 21.212 us/op | 113.14 MB/s | 1548 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 19.035 us/op | 6.72 MB/s | 460 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 21.392 us/op | 112.19 MB/s | 1672 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 18.746 us/op | 6.83 MB/s | 454 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 22.113 us/op | 108.53 MB/s | 1671 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 20.701 us/op | 115.94 MB/s | 1683 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.444 ms/op | - | 88615 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 659.196 us/op | - | 90609 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 1.156 ms/op | - | 93941 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 444.069 us/op | - | 96095 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 425.899 us/op | - | 109342 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 428.748 us/op | - | 109179 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.48 us/op | 36.78 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.682 us/op | 359.19 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.015 ms/conn | 2619432 B/op | 21139 allocs/op |
| 1-RTT application-data round trip | 5 | 2.021 ms/conn | 2631496 B/op | 21477 allocs/op |
| Full mTLS handshake | 5 | 3.587 ms/conn | 3274160 B/op | 29158 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.656 ms/conn | 3730624 B/op | 30900 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.004 ms/conn | 3097304 B/op | 24989 allocs/op |
| Direct external PSK handshake | 5 | 0.5743 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.004 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.151 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.195 ms/conn | 4053952 B/op | 39396 allocs/op |
| Session resumption handshake | 5 | 0.6719 ms/conn | 4651288 B/op | 38434 allocs/op |
| mTLS session resumption handshake | 5 | 0.7744 ms/conn | 6671712 B/op | 49952 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6816 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.381 ms/conn | 3432256 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.353 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.572 ms/conn | 3946832 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.708 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.854 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.299 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.386 ms/conn | 1495176 B/op | 11924 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.711 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 1.045 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.747 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.95 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.206 ms/conn | 1547088 B/op | 12663 allocs/op |
| Session resumption handshake | 5 | 1.167 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.151 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.985 ms/conn | 2006888 B/op | 11261 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.955 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.437 ms/conn | 1303968 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 5.218 ms/conn | 2458264 B/op | 10765 allocs/op |
| Full mTLS handshake | 5 | 6.13 ms/conn | 1703280 B/op | 14581 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.214 ms/conn | 1906640 B/op | 15101 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.45 ms/conn | 1368672 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.949 ms/conn | 1026512 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.207 ms/conn | 2473224 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.373 ms/conn | 2633784 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.28 ms/conn | 3308976 B/op | 22331 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3599712 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 242368 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.676 ms/conn | 1572272 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.641 ms/conn | 1592752 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.017 ms/conn | 1751504 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.934 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.919 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.107 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.81 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64464 B/op | 57 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9934 ms/conn | 452568 B/op | 3437 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.047 ms/conn | 518760 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.015 ms/conn | 389976 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9718 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.092 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.021 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.264 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.045 ms/conn | 507312 B/op | 3816 allocs/op |
| Full mTLS handshake | 5 | 1.425 ms/conn | 537352 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.018 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.025 ms/conn | 487496 B/op | 3664 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.219 ms/conn | 534040 B/op | 4204 allocs/op |
| Direct external PSK handshake | 5 | 0.6593 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.685 ms/pair | 869320 B/op | 7180 allocs/op |
| mTLS session resumption handshake | 5 | 2.156 ms/pair | 1047048 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.705 ms/pair | 868480 B/op | 7130 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.369 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.176 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.264 ms/conn | 181744 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.349 ms/conn | 215344 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.348 ms/conn | 147024 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.899 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.315 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.077 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.456 ms/conn | 275184 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.171 ms/conn | 222992 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.311 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.273 ms/conn | 213648 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.305 ms/conn | 194608 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.9309 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.417 ms/pair | 379536 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.265 ms/pair | 459248 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.478 ms/conn | 203632 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.728 ms/pair | 420976 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.398 ms/conn | 481240 B/op | 1938 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.544 ms/conn | 512840 B/op | 1939 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.335 ms/conn | 453224 B/op | 1925 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.914 ms/conn | 484072 B/op | 1965 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.265 ms/conn | 486408 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.905 ms/conn | 514184 B/op | 2137 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.369 ms/conn | 538424 B/op | 2037 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.254 ms/conn | 514280 B/op | 2125 allocs/op |
| Full mTLS handshake | 5 | 4.896 ms/conn | 523912 B/op | 2511 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.214 ms/conn | 482696 B/op | 1973 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.387 ms/conn | 516232 B/op | 2269 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.353 ms/conn | 506440 B/op | 2149 allocs/op |
| Direct external PSK handshake | 5 | 3.422 ms/conn | 460136 B/op | 1841 allocs/op |
| Session resumption handshake | 5 | 5.329 ms/pair | 703688 B/op | 3831 allocs/op |
| mTLS session resumption handshake | 5 | 5.869 ms/pair | 802824 B/op | 4636 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.569 ms/pair | 706064 B/op | 3791 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.14 ms/conn | 517192 B/op | 2297 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.084 ms/pair | 740136 B/op | 4190 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.496 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.638 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.43 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.678 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.168 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.777 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.388 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.547 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.378 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.423 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.433 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.491 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.516 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.611 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.08 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.127 ms/pair | 207264 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.033 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.062 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9908 ms/conn | 389976 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9623 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.126 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.012 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.271 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.073 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.45 ms/conn | 537224 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.048 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.057 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.208 ms/conn | 534168 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8489 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6658 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.691 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.11 ms/pair | 1047304 B/op | 8531 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.375 ms/conn | 650184 B/op | 5441 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.385 ms/conn | 511448 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.174 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7761 ms/conn | 180320 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.7965 ms/conn | 213344 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8408 ms/conn | 145248 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8118 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.454 ms/conn | 185408 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9569 ms/conn | 203872 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9988 ms/conn | 270944 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9165 ms/conn | 180640 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.129 ms/conn | 220320 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.7743 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.87 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.945 ms/conn | 229472 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7072 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.5736 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.378 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.568 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.9276 ms/conn | 241152 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.9271 ms/conn | 201120 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.706 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.509 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.523 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.531 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.455 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.364 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.617 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.817 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.63 ms/conn | 509864 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 2.958 ms/conn | 516552 B/op | 2341 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.64 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.503 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.592 ms/conn | 502152 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.396 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.284 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 3.186 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 3.49 ms/pair | 789272 B/op | 4370 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.814 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.12 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.171 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.218 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.142 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.365 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.319 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.392 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.204 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.397 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.21 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.131 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.236 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.197 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.072 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.664 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.891 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.298 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9947 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.082 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.026 ms/conn | 389976 B/op | 3400 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.031 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.932 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.108 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.013 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.246 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.032 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.433 ms/conn | 537224 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.001 ms/conn | 454472 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.03 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.212 ms/conn | 533016 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8386 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.6705 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.703 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.161 ms/pair | 1047048 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.709 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.945 ms/conn | 641608 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.103 ms/conn | 211040 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.296 ms/conn | 245216 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.218 ms/conn | 176096 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.105 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.118 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.465 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.102 ms/conn | 223712 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.435 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.077 ms/conn | 211872 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.565 ms/conn | 267008 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.22 ms/conn | 214560 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.068 ms/conn | 247216 B/op | 1985 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.422 ms/conn | 270112 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.143 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.863 ms/conn | 199328 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 2.042 ms/pair | 471232 B/op | 3878 allocs/op |
| mTLS session resumption handshake | 5 | 2.862 ms/pair | 590912 B/op | 4440 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.403 ms/conn | 270400 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.224 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.184 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.2 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.165 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.155 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.544 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.264 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.412 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.172 ms/conn | 514248 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 5.479 ms/conn | 521320 B/op | 2381 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.119 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.251 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.376 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.25 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 5.043 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 6.487 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 7.207 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.222 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.288 ms/conn | 574344 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.234 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.287 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.18 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.197 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.254 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.741 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.303 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.498 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.237 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.392 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.189 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.607 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.633 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.9 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 5.192 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.643 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 7.4 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.865 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 16.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 16.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 17.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 16.55 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 17.27 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 12.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 26.97 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 35.36 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 22.34 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 13.42 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 34.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 33.11 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.717 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 3.282 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 41.26 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 18.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 1.996 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 3.979 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 67.06 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 19.72 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 234.4 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 190.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 200.4 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.42 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 197.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 395.9 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 36.6 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.75 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 528 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.066 us/op | 1982.9 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.668 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 210.8 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 173.7 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 942.4 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.599 us/op | 1138.1 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 331.3 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.56 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 346.2 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 19.56 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.46 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.672 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.187 us/op | 548.66 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 873.1 ns/op | 1374.46 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 16.5 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.284 us/op | 525.42 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.321 us/op | 908.47 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.097 us/op | 169.09 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.158 us/op | 556.15 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.262 us/op | 530.42 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.001 us/op | 399.89 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 925.1 ns/op | 1297.15 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.064 us/op | 391.62 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 859.3 ns/op | 1396.45 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 915.4 ns/op | 1310.86 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.276 us/op | 940.27 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
