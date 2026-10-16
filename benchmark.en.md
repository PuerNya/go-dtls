# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-16T13:24:06Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `bc9c0790253382e73198689dbef57b5c5c4792a5 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 586.776 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 867.68 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 432.105 us/op | - | 123193 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.09 ms/op | - | 123643 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.497 ms/op | - | 142418 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 681.272 us/op | - | 121347 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 919.177 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 919.884 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 368.562 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.039 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.036 ms/op | - | 126045 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.723 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.725 ms/op | - | 162447 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 985.938 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 986.57 us/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 875.528 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 828.247 us/op | - | 152656 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.129 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 134.456 us/op | 0.95 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 145.969 us/op | 16.44 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 135.645 us/op | 0.94 MB/s | 413 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 150.414 us/op | 15.96 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 22.081 us/op | 5.8 MB/s | 422 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 20.399 us/op | 117.65 MB/s | 1541 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 18.702 us/op | 6.84 MB/s | 420 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 21.359 us/op | 112.36 MB/s | 1548 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 20.16 us/op | 6.35 MB/s | 452 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 21.158 us/op | 113.43 MB/s | 1666 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 19.34 us/op | 6.62 MB/s | 459 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 21.15 us/op | 113.48 MB/s | 1676 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 20.848 us/op | 115.12 MB/s | 1670 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.416 ms/op | - | 88615 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 614.436 us/op | - | 90550 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 1.128 ms/op | - | 93942 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 426.153 us/op | - | 95503 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 413.51 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 418.675 us/op | - | 109177 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.479 us/op | 36.79 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.613 us/op | 362.95 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.985 ms/conn | 2619384 B/op | 21138 allocs/op |
| 1-RTT application-data round trip | 5 | 1.996 ms/conn | 2631496 B/op | 21477 allocs/op |
| Full mTLS handshake | 5 | 3.528 ms/conn | 3274160 B/op | 29157 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.573 ms/conn | 3730672 B/op | 30901 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2 ms/conn | 3097256 B/op | 24988 allocs/op |
| Direct external PSK handshake | 5 | 0.5454 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.972 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.112 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.061 ms/conn | 4054000 B/op | 39397 allocs/op |
| Session resumption handshake | 5 | 0.6342 ms/conn | 4650808 B/op | 38432 allocs/op |
| mTLS session resumption handshake | 5 | 0.7982 ms/conn | 6671728 B/op | 49952 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6799 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.346 ms/conn | 3432256 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.342 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.583 ms/conn | 3946816 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.611 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.695 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.176 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.023 ms/conn | 1495872 B/op | 11931 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.508 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.8067 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.519 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.677 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.96 ms/conn | 1547152 B/op | 12664 allocs/op |
| Session resumption handshake | 5 | 1.13 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.158 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.84 ms/conn | 2006024 B/op | 11252 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.063 ms/conn | 2017880 B/op | 11366 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.416 ms/conn | 1303936 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 5.087 ms/conn | 2458632 B/op | 10766 allocs/op |
| Full mTLS handshake | 5 | 6.117 ms/conn | 1703344 B/op | 14581 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.07 ms/conn | 1906240 B/op | 15101 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.339 ms/conn | 1369344 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.964 ms/conn | 1026736 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.089 ms/conn | 2472840 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.141 ms/conn | 2634104 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.69 ms/conn | 3308512 B/op | 22331 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3599360 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 242240 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.615 ms/conn | 1572432 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.687 ms/conn | 1592912 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.973 ms/conn | 1751632 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.72 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63712 B/op | 57 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.719 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.804 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.4 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
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
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9967 ms/conn | 452648 B/op | 3437 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.065 ms/conn | 518920 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.013 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9459 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.124 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.026 ms/conn | 527488 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.284 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.1 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.463 ms/conn | 537096 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.045 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.042 ms/conn | 488008 B/op | 3665 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.233 ms/conn | 534040 B/op | 4204 allocs/op |
| Direct external PSK handshake | 5 | 0.6686 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.71 ms/pair | 868936 B/op | 7181 allocs/op |
| mTLS session resumption handshake | 5 | 2.158 ms/pair | 1047048 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.695 ms/pair | 868480 B/op | 7130 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.396 ms/conn | 511624 B/op | 3989 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.183 ms/pair | 947768 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.372 ms/conn | 181872 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.397 ms/conn | 215664 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.481 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.096 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.408 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.144 ms/conn | 215536 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.471 ms/conn | 275056 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.33 ms/conn | 222800 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.345 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.304 ms/conn | 213648 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.382 ms/conn | 194608 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.8039 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.484 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.398 ms/pair | 459248 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.527 ms/conn | 203632 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.785 ms/pair | 420848 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.264 ms/conn | 481288 B/op | 1943 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.482 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.731 ms/conn | 453192 B/op | 1923 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.148 ms/conn | 484072 B/op | 1965 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.558 ms/conn | 486408 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.327 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.788 ms/conn | 538456 B/op | 2039 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.654 ms/conn | 514504 B/op | 2123 allocs/op |
| Full mTLS handshake | 5 | 5.03 ms/conn | 523720 B/op | 2509 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.369 ms/conn | 482632 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.473 ms/conn | 516232 B/op | 2269 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.508 ms/conn | 506440 B/op | 2149 allocs/op |
| Direct external PSK handshake | 5 | 3.528 ms/conn | 460072 B/op | 1837 allocs/op |
| Session resumption handshake | 5 | 5.72 ms/pair | 703560 B/op | 3824 allocs/op |
| mTLS session resumption handshake | 5 | 6.148 ms/pair | 802760 B/op | 4632 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.348 ms/pair | 706000 B/op | 3787 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.441 ms/conn | 517224 B/op | 2299 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.38 ms/pair | 740488 B/op | 4196 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.524 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.826 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.771 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.892 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.529 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.179 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.713 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.704 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.475 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.55 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.44 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.503 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.626 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.628 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.251 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.451 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.034 ms/conn | 452472 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.093 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.061 ms/conn | 389976 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9457 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.115 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.035 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.276 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.036 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.464 ms/conn | 537224 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.997 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.067 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.231 ms/conn | 534168 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.841 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6236 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.707 ms/pair | 868168 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.123 ms/pair | 1047560 B/op | 8531 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.403 ms/conn | 650264 B/op | 5440 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.389 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.186 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8244 ms/conn | 180192 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9267 ms/conn | 213344 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8007 ms/conn | 145248 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8225 ms/conn | 183328 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.383 ms/conn | 185280 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9656 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9663 ms/conn | 270944 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7902 ms/conn | 180640 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.11 ms/conn | 220320 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8243 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.856 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.003 ms/conn | 229472 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.699 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.609 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.355 ms/pair | 388512 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.753 ms/pair | 474208 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.058 ms/conn | 241280 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.8869 ms/conn | 201248 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.563 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.523 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.568 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.574 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.51 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.105 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.614 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.748 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.521 ms/conn | 509864 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 2.841 ms/conn | 516552 B/op | 2341 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.71 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.612 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.817 ms/conn | 502152 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.471 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.259 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 3.237 ms/pair | 693080 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 3.486 ms/pair | 789384 B/op | 4370 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.916 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.16 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.248 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.196 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.184 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.389 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.221 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.239 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.204 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.39 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.178 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.232 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.318 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.15 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.031 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.692 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.901 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.309 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.001 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.046 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.026 ms/conn | 390104 B/op | 3400 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.019 ms/conn | 517128 B/op | 4122 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9384 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.132 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.009 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.267 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.048 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.408 ms/conn | 537352 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.01 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.041 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.209 ms/conn | 532888 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8338 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.6809 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.724 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 2.171 ms/pair | 1047176 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.708 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.924 ms/conn | 641992 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.114 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.098 ms/conn | 245088 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.097 ms/conn | 176224 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.142 ms/conn | 264672 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.171 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.412 ms/conn | 220288 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.201 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.463 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.138 ms/conn | 211872 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.454 ms/conn | 267008 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.146 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.209 ms/conn | 247424 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.484 ms/conn | 269984 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.021 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.9282 ms/conn | 199328 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 2.127 ms/pair | 470944 B/op | 3877 allocs/op |
| mTLS session resumption handshake | 5 | 2.731 ms/pair | 592752 B/op | 4441 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.266 ms/conn | 270528 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.303 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.508 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.41 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.234 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.774 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.135 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.596 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.653 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.246 ms/conn | 514248 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 5.6 ms/conn | 521320 B/op | 2381 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.583 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.732 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.487 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.981 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 5.115 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 6.589 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 6.672 ms/pair | 801192 B/op | 4506 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.272 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.287 ms/conn | 574600 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.504 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.743 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.33 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.4 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.768 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.251 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.478 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.971 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.367 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.516 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.589 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.592 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.552 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.179 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 5.195 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.708 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 7.008 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.648 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 16.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 16.55 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 17.17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 16.52 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 17 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 12.19 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 27.77 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 29.62 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 23.03 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 13.41 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 35.21 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 33.13 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.734 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 3.335 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 39.62 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 18.14 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 2.029 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 4.041 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 61.25 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 28.49 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 232.7 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 188.9 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 198 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.379 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 198.2 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 363.8 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 36.28 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.35 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 521.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.007 us/op | 2041.17 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.668 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 211 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 173.4 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 939.6 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.57 us/op | 1147.29 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 319.7 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.2 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 343.8 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 19.58 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.39 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.681 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.154 us/op | 557.05 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 863.4 ns/op | 1389.83 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 16.49 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.225 us/op | 539.21 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.27 us/op | 944.75 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.026 us/op | 170.79 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.114 us/op | 567.69 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.222 us/op | 540.13 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.96 us/op | 405.41 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 920 ns/op | 1304.36 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.044 us/op | 394.27 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 844.6 ns/op | 1420.82 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 899.8 ns/op | 1333.6 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.26 us/op | 952.49 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
