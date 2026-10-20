# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-20T02:53:45Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `5f9ab72a4775c4a7ab1a927455178e846f003dc4 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 596.827 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 890.002 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 453.288 us/op | - | 123193 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.125 ms/op | - | 123643 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.524 ms/op | - | 142415 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 691.935 us/op | - | 121345 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 945.511 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 946.798 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 378.608 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.044 ms/op | - | 136092 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.047 ms/op | - | 126047 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.746 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.756 ms/op | - | 162482 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.005 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 1.01 ms/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 913.535 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 851.801 us/op | - | 152656 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.16 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 138.373 us/op | 0.93 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 149.997 us/op | 16 MB/s | 1534 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 139.534 us/op | 0.92 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 153.248 us/op | 15.66 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 21.784 us/op | 5.88 MB/s | 417 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 20.796 us/op | 115.4 MB/s | 1542 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 21.285 us/op | 6.01 MB/s | 415 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 21.718 us/op | 110.51 MB/s | 1543 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 22.48 us/op | 5.69 MB/s | 489 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 20.739 us/op | 115.72 MB/s | 1674 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 19.312 us/op | 6.63 MB/s | 483 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 21.368 us/op | 112.31 MB/s | 1670 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 21.332 us/op | 112.5 MB/s | 1673 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.417 ms/op | - | 88615 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 622.528 us/op | - | 90344 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 1.128 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 427.734 us/op | - | 95364 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 420.459 us/op | - | 109340 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 421.98 us/op | - | 109179 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.499 us/op | 36.58 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.782 us/op | 353.87 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.992 ms/conn | 2619432 B/op | 21139 allocs/op |
| 1-RTT application-data round trip | 5 | 2.001 ms/conn | 2631544 B/op | 21478 allocs/op |
| Full mTLS handshake | 5 | 3.522 ms/conn | 3274208 B/op | 29158 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.64 ms/conn | 3730624 B/op | 30900 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.136 ms/conn | 3097256 B/op | 24988 allocs/op |
| Direct external PSK handshake | 5 | 0.5693 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.996 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.147 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.104 ms/conn | 4053952 B/op | 39396 allocs/op |
| Session resumption handshake | 5 | 0.6665 ms/conn | 4650808 B/op | 38436 allocs/op |
| mTLS session resumption handshake | 5 | 0.767 ms/conn | 6671728 B/op | 49953 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.7029 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.366 ms/conn | 3432272 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.329 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.562 ms/conn | 3946832 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.57 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.72 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.194 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.366 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.785 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.9605 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.749 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.976 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.275 ms/conn | 1547344 B/op | 12667 allocs/op |
| Session resumption handshake | 5 | 1.22 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.192 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.93 ms/conn | 2006888 B/op | 11261 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.95 ms/conn | 2018584 B/op | 11373 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.445 ms/conn | 1303936 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 5.13 ms/conn | 2458264 B/op | 10765 allocs/op |
| Full mTLS handshake | 5 | 6.175 ms/conn | 1703696 B/op | 14581 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.269 ms/conn | 1906016 B/op | 15101 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.459 ms/conn | 1368896 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.981 ms/conn | 1026064 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.132 ms/conn | 2472456 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.285 ms/conn | 2634104 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.93 ms/conn | 3308736 B/op | 22332 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3598832 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 242560 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.629 ms/conn | 1572432 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.695 ms/conn | 1592592 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.758 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63664 B/op | 56 allocs/op |
| 1-RTT application-data round trip | 5 | 7.741 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.806 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.913 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.57 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191192 B/op | 1203 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.057 ms/conn | 452728 B/op | 3438 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.04 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.027 ms/conn | 390152 B/op | 3401 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9442 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.081 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.014 ms/conn | 527664 B/op | 3896 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.386 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.017 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.409 ms/conn | 537352 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.012 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.024 ms/conn | 487496 B/op | 3664 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.242 ms/conn | 534168 B/op | 4204 allocs/op |
| Direct external PSK handshake | 5 | 0.6945 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.655 ms/pair | 868888 B/op | 7181 allocs/op |
| mTLS session resumption handshake | 5 | 2.147 ms/pair | 1047176 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.692 ms/pair | 868480 B/op | 7130 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.324 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.149 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.317 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.409 ms/conn | 215344 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.418 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.874 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.358 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.097 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.461 ms/conn | 275184 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.146 ms/conn | 222992 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.31 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.358 ms/conn | 213648 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.34 ms/conn | 194736 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.8602 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.653 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.274 ms/pair | 459440 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.47 ms/conn | 203632 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.813 ms/pair | 420976 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.312 ms/conn | 481224 B/op | 1939 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.316 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.18 ms/conn | 453224 B/op | 1925 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.804 ms/conn | 484072 B/op | 1965 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.199 ms/conn | 486440 B/op | 1999 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.893 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.401 ms/conn | 538424 B/op | 2036 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.223 ms/conn | 514248 B/op | 2123 allocs/op |
| Full mTLS handshake | 5 | 4.818 ms/conn | 523880 B/op | 2509 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.215 ms/conn | 482664 B/op | 1971 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.333 ms/conn | 516296 B/op | 2273 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.265 ms/conn | 506440 B/op | 2149 allocs/op |
| Direct external PSK handshake | 5 | 3.385 ms/conn | 460072 B/op | 1837 allocs/op |
| Session resumption handshake | 5 | 5.261 ms/pair | 703592 B/op | 3824 allocs/op |
| mTLS session resumption handshake | 5 | 5.925 ms/pair | 802728 B/op | 4630 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.334 ms/pair | 706000 B/op | 3787 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.135 ms/conn | 517192 B/op | 2297 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.139 ms/pair | 740424 B/op | 4192 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.352 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.491 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.359 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.566 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.253 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.813 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.377 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.452 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.386 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.444 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.415 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.479 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.337 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.508 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.011 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.113 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9766 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.066 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9913 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.027 ms/conn | 459256 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.092 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.034 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.26 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.067 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.475 ms/conn | 537352 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9868 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.023 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.23 ms/conn | 532888 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8133 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6106 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.681 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.133 ms/pair | 1047304 B/op | 8531 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.38 ms/conn | 650312 B/op | 5441 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.376 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.147 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8386 ms/conn | 180192 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.899 ms/conn | 213472 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8207 ms/conn | 145376 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8137 ms/conn | 183200 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.421 ms/conn | 185280 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.888 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.033 ms/conn | 270944 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7939 ms/conn | 180640 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 0.9894 ms/conn | 220320 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8415 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.7764 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9666 ms/conn | 228704 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7343 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.52 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.332 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.565 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.029 ms/conn | 241280 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.9315 ms/conn | 201248 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.512 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.453 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.517 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.459 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.407 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.06 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.479 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.652 ms/conn | 530072 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.451 ms/conn | 509800 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 2.688 ms/conn | 516552 B/op | 2341 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.475 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.51 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.723 ms/conn | 501640 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.389 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.177 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 3.074 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 3.477 ms/pair | 789256 B/op | 4370 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.967 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.11 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.185 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.113 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.151 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.218 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.184 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.248 ms/conn | 212512 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.156 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.327 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.176 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.116 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.254 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.143 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.997 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.655 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.85 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.27 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9835 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.05 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.005 ms/conn | 390104 B/op | 3400 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.06 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9162 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.078 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.029 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.252 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.05 ms/conn | 507104 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.443 ms/conn | 537096 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9971 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.045 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.247 ms/conn | 533016 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8107 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.6083 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.668 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.112 ms/pair | 1047304 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.713 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.922 ms/conn | 641864 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.03 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.183 ms/conn | 245088 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.127 ms/conn | 176224 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.067 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.112 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.55 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.155 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.392 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.066 ms/conn | 211744 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.493 ms/conn | 266880 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.116 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.121 ms/conn | 247424 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.401 ms/conn | 270112 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.9821 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.9456 ms/conn | 199328 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 2.005 ms/pair | 471232 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 2.601 ms/pair | 590768 B/op | 4440 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.443 ms/conn | 270400 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.189 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.231 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.294 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.153 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.089 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.45 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.118 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.314 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.124 ms/conn | 514120 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 5.391 ms/conn | 521320 B/op | 2381 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.087 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.292 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.264 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.584 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 4.785 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 5.959 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 6.34 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.12 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.01 ms/conn | 574600 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.229 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.342 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.145 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.142 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.232 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.65 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.222 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.446 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.192 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.349 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.166 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.22 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.367 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.622 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.948 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.186 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 6.817 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.654 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 16.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 16.51 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 17.18 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 16.54 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 19.35 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 12.22 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 26.95 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 29.7 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 22.97 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 13.46 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 36.45 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 36.48 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.781 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 3.445 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 40.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 18.19 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 1.928 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 3.929 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 64.81 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 19.74 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 239 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 195.4 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 204.6 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.456 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 203.7 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 382.3 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 36.76 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 45.22 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 544.7 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.118 us/op | 1933.71 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.698 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 213 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 175.7 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 963 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.689 us/op | 1110.28 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 333.6 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.69 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 352.6 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 19.57 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.6 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.68 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.245 us/op | 534.52 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 887.5 ns/op | 1352.18 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 16.51 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.307 us/op | 520.21 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.31 us/op | 916.24 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.172 us/op | 167.31 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.201 us/op | 545.23 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.306 us/op | 520.44 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.077 us/op | 390.05 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 971.1 ns/op | 1235.71 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.109 us/op | 386 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 872 ns/op | 1376.19 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 927.1 ns/op | 1294.38 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.3 us/op | 922.99 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
