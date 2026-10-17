# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-17T02:21:04Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `16cac42a0db8ea6b4df0ba7fe8737e76a5085d6e (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 471.195 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 708.668 us/op | - | 116336 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 361.664 us/op | - | 123196 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 890.901 us/op | - | 123643 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.23 ms/op | - | 142428 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 560.806 us/op | - | 121347 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 764.025 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 766.303 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 294.077 us/op | - | 103775 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 845.528 us/op | - | 136094 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 841.732 us/op | - | 126047 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.413 ms/op | - | 177157 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.416 ms/op | - | 162386 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 823.853 us/op | - | 151404 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 817.124 us/op | - | 154181 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 720.601 us/op | - | 149313 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 677.49 us/op | - | 152657 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.802 ms/op | - | 178099 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 77.497 us/op | 1.65 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 81.318 us/op | 29.51 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 77.459 us/op | 1.65 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 82.003 us/op | 29.27 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 16.088 us/op | 7.96 MB/s | 415 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 16.414 us/op | 146.22 MB/s | 1542 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 16.256 us/op | 7.87 MB/s | 417 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 16.976 us/op | 141.37 MB/s | 1545 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 15.976 us/op | 8.01 MB/s | 493 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 17.013 us/op | 141.07 MB/s | 1687 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 16.582 us/op | 7.72 MB/s | 494 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 15.527 us/op | 154.57 MB/s | 1682 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 15.83 us/op | 151.61 MB/s | 1693 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.03 ms/op | - | 88613 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 461.23 us/op | - | 90126 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 778.156 us/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 308.798 us/op | - | 95276 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 340.907 us/op | - | 109375 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 343.99 us/op | - | 109181 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 2.568 us/op | 49.85 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 5.111 us/op | 469.56 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.643 ms/conn | 2619432 B/op | 21139 allocs/op |
| 1-RTT application-data round trip | 5 | 1.661 ms/conn | 2631496 B/op | 21477 allocs/op |
| Full mTLS handshake | 5 | 2.962 ms/conn | 3274160 B/op | 29157 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.017 ms/conn | 3730624 B/op | 30900 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.677 ms/conn | 3097256 B/op | 24988 allocs/op |
| Direct external PSK handshake | 5 | 0.4263 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.658 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.777 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.385 ms/conn | 4053952 B/op | 39394 allocs/op |
| Session resumption handshake | 5 | 0.5016 ms/conn | 4650808 B/op | 38426 allocs/op |
| mTLS session resumption handshake | 5 | 0.6358 ms/conn | 6671696 B/op | 49952 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5217 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.96 ms/conn | 3432256 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.962 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.014 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.976 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.007 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 5.221 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.213 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.989 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.831 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.96 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.129 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.145 ms/conn | 1547536 B/op | 12670 allocs/op |
| Session resumption handshake | 5 | 0.9195 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 0.9264 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.191 ms/conn | 2005320 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.019 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.004 ms/conn | 1304160 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 4.242 ms/conn | 2458584 B/op | 10765 allocs/op |
| Full mTLS handshake | 5 | 5.115 ms/conn | 1703472 B/op | 14581 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.113 ms/conn | 1906080 B/op | 15102 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.02 ms/conn | 1368896 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.711 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.249 ms/conn | 2472840 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.297 ms/conn | 2633464 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.74 ms/conn | 3309904 B/op | 22332 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 3599376 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 242368 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.199 ms/conn | 1572752 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.996 ms/conn | 1593232 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.716 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 6.483 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.465 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.564 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.603 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1014 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.839 ms/conn | 452696 B/op | 3438 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8485 ms/conn | 518920 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7998 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7606 ms/conn | 459256 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.762 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7993 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.017 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8403 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.135 ms/conn | 537736 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8006 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8264 ms/conn | 487496 B/op | 3664 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9725 ms/conn | 533016 B/op | 4204 allocs/op |
| Direct external PSK handshake | 5 | 0.4896 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.301 ms/pair | 870664 B/op | 7181 allocs/op |
| mTLS session resumption handshake | 5 | 1.704 ms/pair | 1047304 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.326 ms/pair | 868480 B/op | 7130 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.098 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.718 ms/pair | 947768 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.109 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.093 ms/conn | 215664 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.05 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.621 ms/conn | 184752 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.56 ms/conn | 186832 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.736 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.003 ms/conn | 275184 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 1.814 ms/conn | 223120 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.002 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.069 ms/conn | 213648 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.015 ms/conn | 194608 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.637 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 1.84 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 2.626 ms/pair | 459312 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.144 ms/conn | 203632 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.149 ms/pair | 420848 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.471 ms/conn | 481192 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.56 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.463 ms/conn | 453192 B/op | 1923 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4 ms/conn | 484008 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.043 ms/conn | 486472 B/op | 2000 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.063 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.311 ms/conn | 538424 B/op | 2037 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 3.449 ms/conn | 514504 B/op | 2123 allocs/op |
| Full mTLS handshake | 5 | 3.99 ms/conn | 523976 B/op | 2511 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.487 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.596 ms/conn | 516200 B/op | 2267 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.551 ms/conn | 505928 B/op | 2149 allocs/op |
| Direct external PSK handshake | 5 | 2.78 ms/conn | 460104 B/op | 1839 allocs/op |
| Session resumption handshake | 5 | 4.225 ms/pair | 703560 B/op | 3824 allocs/op |
| mTLS session resumption handshake | 5 | 4.801 ms/pair | 802824 B/op | 4628 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.306 ms/pair | 705936 B/op | 3783 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 4.222 ms/conn | 517224 B/op | 2299 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 4.969 ms/pair | 740424 B/op | 4192 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.698 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.752 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.675 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.743 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.724 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.933 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.231 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 4.572 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.674 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.689 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.696 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 2.886 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.328 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 6.187 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 4.264 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 5.923 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7895 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8255 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8114 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.759 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.743 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7956 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.028 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8243 ms/conn | 507104 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.125 ms/conn | 537480 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.7886 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8146 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9669 ms/conn | 533016 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6537 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4807 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.303 ms/pair | 868168 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 1.71 ms/pair | 1047176 B/op | 8531 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.098 ms/conn | 650264 B/op | 5440 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.086 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.711 ms/pair | 947768 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.6427 ms/conn | 180320 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.699 ms/conn | 213472 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.6153 ms/conn | 145376 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.5615 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.026 ms/conn | 185280 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.6827 ms/conn | 203872 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.8236 ms/conn | 270816 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.6915 ms/conn | 180512 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 0.8676 ms/conn | 220320 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.6332 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.6095 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.8151 ms/conn | 228832 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.5515 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4153 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 0.9961 ms/pair | 388512 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.215 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.7971 ms/conn | 241280 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.7223 ms/conn | 201248 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.275 ms/pair | 432128 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.01 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.041 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.994 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.99 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.43 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.03 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.188 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.031 ms/conn | 509672 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 2.215 ms/conn | 516552 B/op | 2341 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.993 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.003 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.072 ms/conn | 501640 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.916 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.825 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 2.55 ms/pair | 693048 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 2.788 ms/pair | 789256 B/op | 4370 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.34 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.775 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.838 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.766 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.752 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.607 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.843 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.909 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.801 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 1.871 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.765 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.772 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.818 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.763 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.645 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.153 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.32 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.918 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7873 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8303 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8197 ms/conn | 390104 B/op | 3400 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8085 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7362 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.744 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7981 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.023 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8331 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.128 ms/conn | 537096 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8144 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8091 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9627 ms/conn | 532888 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6463 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.4847 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.308 ms/pair | 868168 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 1.691 ms/pair | 1047304 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.349 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.546 ms/conn | 641480 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8737 ms/conn | 211040 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9016 ms/conn | 245216 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.23 ms/conn | 176224 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8434 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8239 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.845 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8766 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.109 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8383 ms/conn | 211872 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.118 ms/conn | 266880 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8309 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8062 ms/conn | 247136 B/op | 1986 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.08 ms/conn | 269984 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8352 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.6702 ms/conn | 199328 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.485 ms/pair | 471104 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 1.994 ms/pair | 590608 B/op | 4438 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 0.9461 ms/conn | 271144 B/op | 2313 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.325 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.361 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.324 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.279 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.268 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.31 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.358 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.646 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.32 ms/conn | 514312 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 4.594 ms/conn | 521320 B/op | 2381 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.271 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.35 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.415 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.671 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 3.996 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 5.058 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 5.4 ms/pair | 801320 B/op | 4506 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.879 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.943 ms/conn | 574344 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.334 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.317 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.331 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.341 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.393 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.268 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.429 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.995 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.336 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 4.643 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.321 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.37 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.425 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.736 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.135 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.136 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 5.619 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 4.574 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 11.48 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 13.05 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 12.03 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 13.1 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 12.11 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 9.562 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 19.97 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 26.37 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 14.95 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 9.971 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 24.47 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 26.09 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.126 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 2.141 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 28.43 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 14.09 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 1.641 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 3.025 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 38.96 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 15.27 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 195.7 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 169.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 172.3 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.046 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 172.1 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 304 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 27.98 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 34.45 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 363 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.432 us/op | 2861.24 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.179 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 155.8 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 136 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 671.2 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.692 us/op | 1521.67 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 233.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 70.77 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 268.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 16.19 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 18.3 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.823 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 1.582 us/op | 758.48 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 641.7 ns/op | 1869.97 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.02 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 1.627 us/op | 737.43 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 983.2 ns/op | 1220.51 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 5.627 us/op | 213.26 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 1.566 us/op | 766.11 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 1.652 us/op | 726.28 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.36 us/op | 508.43 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 680.6 ns/op | 1763.24 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.476 us/op | 484.75 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 639.6 ns/op | 1876.04 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 680.9 ns/op | 1762.49 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.028 us/op | 1167.14 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
