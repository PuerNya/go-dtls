# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12a5008a3c8aa71239d32fd386029316e649de0e`
- Generated: `2026-10-19T15:40:05Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `ecb4697cfba6de50082ecae39c19e413067d4bf0 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 496.805 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 741.504 us/op | - | 116336 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 383.898 us/op | - | 123190 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 943.862 us/op | - | 123646 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.275 ms/op | - | 142422 B/op | 1403 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 576.617 us/op | - | 121353 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 780.351 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 781.233 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 319.828 us/op | - | 103775 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 893.168 us/op | - | 136094 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 891.966 us/op | - | 126049 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.475 ms/op | - | 177158 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.482 ms/op | - | 162534 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 866.324 us/op | - | 151405 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 863.823 us/op | - | 154181 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 774.468 us/op | - | 149314 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 722.435 us/op | - | 152658 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.867 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 76.604 us/op | 1.67 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 81.172 us/op | 29.57 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 77.566 us/op | 1.65 MB/s | 412 B/op | 15 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 82.174 us/op | 29.21 MB/s | 1533 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 16.271 us/op | 7.87 MB/s | 417 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 17.316 us/op | 138.6 MB/s | 1542 B/op | 14 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 16.584 us/op | 7.72 MB/s | 416 B/op | 15 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 16.579 us/op | 144.76 MB/s | 1542 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 16.535 us/op | 7.74 MB/s | 494 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 16.334 us/op | 146.93 MB/s | 1686 B/op | 14 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 16.458 us/op | 7.78 MB/s | 494 B/op | 15 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 16.15 us/op | 148.61 MB/s | 1681 B/op | 14 allocs/op |
| Listener UDP round trip / No server read deadline / 1024 connections / CID / 1200 bytes | 5 | 15.72 us/op | 152.67 MB/s | 1679 B/op | 14 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 1 | 5 | 1.047 ms/op | - | 88613 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Full handshake / Concurrency 64 | 5 | 500.313 us/op | - | 90399 B/op | 767 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 1 | 5 | 797.938 us/op | - | 93941 B/op | 805 allocs/op |
| Listener UDP handshake, ticket and close / Session resumption / Concurrency 64 | 5 | 339.054 us/op | - | 95282 B/op | 805 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 354.798 us/op | - | 109340 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 356.877 us/op | - | 109180 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 2.634 us/op | 48.59 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 5.209 us/op | 460.72 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.678 ms/conn | 2619432 B/op | 21139 allocs/op |
| 1-RTT application-data round trip | 5 | 1.685 ms/conn | 2631544 B/op | 21478 allocs/op |
| Full mTLS handshake | 5 | 3.035 ms/conn | 3274208 B/op | 29158 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.086 ms/conn | 3730672 B/op | 30901 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.723 ms/conn | 3097256 B/op | 24988 allocs/op |
| Direct external PSK handshake | 5 | 0.435 ms/conn | 1752600 B/op | 14596 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.674 ms/conn | 2636808 B/op | 22257 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.817 ms/conn | 2799592 B/op | 22791 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.295 ms/conn | 4053952 B/op | 39387 allocs/op |
| Session resumption handshake | 5 | 0.5294 ms/conn | 4650808 B/op | 38436 allocs/op |
| mTLS session resumption handshake | 5 | 0.6341 ms/conn | 6671712 B/op | 49952 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5343 ms/conn | 294944 B/op | 1912 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.992 ms/conn | 3432256 B/op | 21425 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.983 ms/conn | 3474176 B/op | 21765 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.069 ms/conn | 3946832 B/op | 22125 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.009 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.062 ms/conn | 1354224 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 5.256 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.281 ms/conn | 1495072 B/op | 11923 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.022 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.8508 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.013 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.139 ms/conn | 1590544 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.195 ms/conn | 1547728 B/op | 12673 allocs/op |
| Session resumption handshake | 5 | 0.9371 ms/conn | 2370840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 0.9201 ms/conn | 2538376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.209 ms/conn | 2006104 B/op | 11253 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.063 ms/conn | 2018504 B/op | 11372 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.01 ms/conn | 1304160 B/op | 9259 allocs/op |
| 1-RTT application-data round trip | 5 | 4.307 ms/conn | 2458264 B/op | 10765 allocs/op |
| Full mTLS handshake | 5 | 5.176 ms/conn | 1703472 B/op | 14581 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.168 ms/conn | 1905968 B/op | 15101 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.027 ms/conn | 1368816 B/op | 10179 allocs/op |
| Direct external PSK handshake | 5 | 0.773 ms/conn | 1026736 B/op | 7419 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.5 ms/conn | 2473032 B/op | 11327 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.438 ms/conn | 2633784 B/op | 12405 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.02 ms/conn | 3308912 B/op | 22332 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 3598736 B/op | 19987 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 242240 B/op | 991 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.203 ms/conn | 1572752 B/op | 9681 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.055 ms/conn | 1592752 B/op | 9821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.809 ms/conn | 1751760 B/op | 10001 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 6.726 ms/conn | 1169864 B/op | 1203 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.694 ms/conn | 1183512 B/op | 1203 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.597 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.867 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1013 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8258 ms/conn | 452520 B/op | 3437 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8748 ms/conn | 518888 B/op | 3437 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9314 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7601 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.777 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8436 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.046 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8938 ms/conn | 507104 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.159 ms/conn | 537224 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8155 ms/conn | 454648 B/op | 3485 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8393 ms/conn | 487448 B/op | 3665 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.002 ms/conn | 533016 B/op | 4204 allocs/op |
| Direct external PSK handshake | 5 | 0.4842 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.335 ms/pair | 870664 B/op | 7181 allocs/op |
| mTLS session resumption handshake | 5 | 1.762 ms/pair | 1047560 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.358 ms/pair | 868480 B/op | 7130 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.113 ms/conn | 511576 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.774 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.029 ms/conn | 181744 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.12 ms/conn | 215600 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.068 ms/conn | 147152 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.532 ms/conn | 184880 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.648 ms/conn | 186960 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.711 ms/conn | 215664 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.036 ms/conn | 275184 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 1.803 ms/conn | 222992 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.119 ms/conn | 182160 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.038 ms/conn | 213520 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.075 ms/conn | 194608 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.6658 ms/conn | 159216 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 1.825 ms/pair | 379664 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 2.664 ms/pair | 459568 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.18 ms/conn | 203632 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.198 ms/pair | 420976 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.732 ms/conn | 481192 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.767 ms/conn | 512808 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.618 ms/conn | 453160 B/op | 1921 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.265 ms/conn | 484040 B/op | 1963 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.152 ms/conn | 486456 B/op | 1999 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.268 ms/conn | 514120 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.555 ms/conn | 538424 B/op | 2037 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 3.651 ms/conn | 514472 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 4.17 ms/conn | 523720 B/op | 2509 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.868 ms/conn | 482680 B/op | 1971 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.686 ms/conn | 516216 B/op | 2267 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.843 ms/conn | 505928 B/op | 2149 allocs/op |
| Direct external PSK handshake | 5 | 2.84 ms/conn | 460072 B/op | 1837 allocs/op |
| Session resumption handshake | 5 | 4.515 ms/pair | 703592 B/op | 3826 allocs/op |
| mTLS session resumption handshake | 5 | 5.028 ms/pair | 802984 B/op | 4630 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.556 ms/pair | 705968 B/op | 3786 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 4.599 ms/conn | 517224 B/op | 2299 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 5.379 ms/pair | 740424 B/op | 4192 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.838 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.922 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.859 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.947 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.872 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.135 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.481 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 4.836 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.937 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.856 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.958 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.024 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.604 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 6.519 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 4.519 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.371 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.821 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8784 ms/conn | 518712 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8154 ms/conn | 390104 B/op | 3400 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7786 ms/conn | 459256 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.779 ms/conn | 464312 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8281 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.047 ms/conn | 622264 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8591 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.204 ms/conn | 537224 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8194 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.847 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.033 ms/conn | 532888 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6702 ms/conn | 437400 B/op | 3372 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4923 ms/conn | 419896 B/op | 3184 allocs/op |
| Session resumption handshake | 5 | 1.382 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 1.708 ms/pair | 1047304 B/op | 8531 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.153 ms/conn | 650264 B/op | 5440 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.116 ms/conn | 511448 B/op | 3988 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.763 ms/pair | 947896 B/op | 7911 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.6504 ms/conn | 180320 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.698 ms/conn | 213344 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.6099 ms/conn | 145376 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.5733 ms/conn | 183200 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.141 ms/conn | 185408 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7223 ms/conn | 204000 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.8407 ms/conn | 270816 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.6451 ms/conn | 180640 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 0.8525 ms/conn | 220320 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.605 ms/conn | 180512 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.615 ms/conn | 212096 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.7915 ms/conn | 228704 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.5431 ms/conn | 164384 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.435 ms/conn | 160864 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.137 ms/pair | 388640 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.229 ms/pair | 474080 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.7569 ms/conn | 241280 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.6798 ms/conn | 201120 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.195 ms/pair | 432256 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.225 ms/conn | 476584 B/op | 1861 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.298 ms/conn | 508392 B/op | 1861 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.134 ms/conn | 448872 B/op | 1845 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.083 ms/conn | 479592 B/op | 1885 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.635 ms/conn | 481928 B/op | 1921 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.196 ms/conn | 510024 B/op | 2057 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.36 ms/conn | 530328 B/op | 1875 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.122 ms/conn | 509800 B/op | 2045 allocs/op |
| Full mTLS handshake | 5 | 2.346 ms/conn | 516552 B/op | 2341 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.231 ms/conn | 478120 B/op | 1893 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.204 ms/conn | 492264 B/op | 1973 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.214 ms/conn | 501640 B/op | 2081 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.003 ms/conn | 479688 B/op | 1957 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.917 ms/conn | 457032 B/op | 1761 allocs/op |
| Session resumption handshake | 5 | 2.768 ms/pair | 693064 B/op | 3658 allocs/op |
| mTLS session resumption handshake | 5 | 2.949 ms/pair | 789256 B/op | 4370 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.515 ms/conn | 558280 B/op | 2741 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.91 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.931 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.867 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.832 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.793 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.95 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.995 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.855 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 1.969 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.881 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.875 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.931 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.995 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.789 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.233 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.446 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.972 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7977 ms/conn | 452600 B/op | 3436 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8603 ms/conn | 518840 B/op | 3436 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8071 ms/conn | 389976 B/op | 3400 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8417 ms/conn | 517256 B/op | 4122 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7789 ms/conn | 459384 B/op | 3484 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.808 ms/conn | 464184 B/op | 3556 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8299 ms/conn | 527616 B/op | 3895 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.052 ms/conn | 622392 B/op | 3586 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8403 ms/conn | 507232 B/op | 3815 allocs/op |
| Full mTLS handshake | 5 | 1.187 ms/conn | 537480 B/op | 4556 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8185 ms/conn | 454600 B/op | 3484 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8383 ms/conn | 485128 B/op | 3656 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9888 ms/conn | 533016 B/op | 4204 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.648 ms/conn | 437400 B/op | 3372 allocs/op |
| Direct external PSK handshake | 5 | 0.4764 ms/conn | 417144 B/op | 3176 allocs/op |
| Session resumption handshake | 5 | 1.32 ms/pair | 868296 B/op | 7179 allocs/op |
| mTLS session resumption handshake | 5 | 1.711 ms/pair | 1047432 B/op | 8531 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.392 ms/pair | 868480 B/op | 7130 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.571 ms/conn | 641864 B/op | 5970 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8579 ms/conn | 211168 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9074 ms/conn | 245216 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8704 ms/conn | 176096 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8833 ms/conn | 264800 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8678 ms/conn | 220512 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.873 ms/conn | 220176 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9016 ms/conn | 223840 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.157 ms/conn | 331872 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8958 ms/conn | 211872 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.156 ms/conn | 267136 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8182 ms/conn | 214688 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8728 ms/conn | 247280 B/op | 1987 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.074 ms/conn | 270112 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7136 ms/conn | 194272 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.6523 ms/conn | 199328 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.519 ms/pair | 472112 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 1.968 ms/pair | 591904 B/op | 4440 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 0.9984 ms/conn | 271144 B/op | 2313 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.322 ms/conn | 481032 B/op | 1885 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.53 ms/conn | 512776 B/op | 1885 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.431 ms/conn | 453320 B/op | 1869 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.613 ms/conn | 496008 B/op | 2089 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.374 ms/conn | 483848 B/op | 1909 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.7 ms/conn | 486184 B/op | 1945 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.622 ms/conn | 514216 B/op | 2081 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.136 ms/conn | 533432 B/op | 1899 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.975 ms/conn | 514312 B/op | 2069 allocs/op |
| Full mTLS handshake | 5 | 5.101 ms/conn | 521448 B/op | 2381 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.443 ms/conn | 482440 B/op | 1917 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.423 ms/conn | 513960 B/op | 2137 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.483 ms/conn | 506088 B/op | 2105 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.775 ms/conn | 481096 B/op | 1949 allocs/op |
| Direct external PSK handshake | 5 | 4.21 ms/conn | 459976 B/op | 1773 allocs/op |
| Session resumption handshake | 5 | 4.998 ms/pair | 704680 B/op | 3778 allocs/op |
| mTLS session resumption handshake | 5 | 5.615 ms/pair | 801320 B/op | 4506 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.418 ms/pair | 704464 B/op | 3698 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.235 ms/conn | 574600 B/op | 3085 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.586 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.583 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.437 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.472 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.429 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.384 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.786 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.107 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.83 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 4.674 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.53 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.531 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.631 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.831 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.241 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.198 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 5.853 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 4.856 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Listener CID routing / 1 connection / Fixed length / Hit | 5 | 11.53 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Fixed length / Miss | 5 | 13.1 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Hit | 5 | 12.02 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1 connection / Mixed lengths / Miss | 5 | 13.07 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Hit | 5 | 12.22 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Fixed length / Miss | 5 | 9.594 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Hit | 5 | 20.18 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 64 connections / Mixed lengths / Miss | 5 | 27.36 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Hit | 5 | 15.07 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Fixed length / Miss | 5 | 10.27 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Hit | 5 | 25.02 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Mixed lengths / Miss | 5 | 28.48 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Hit | 5 | 1.142 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / Lengths 4–255 / Miss | 5 | 2.171 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Hit | 5 | 28.7 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 252 connections / 255-byte CIDs / Miss | 5 | 14.06 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Hit | 5 | 1.672 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / Lengths 4–255 / Miss | 5 | 3.007 us/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Hit | 5 | 39.43 ns/op | - | 0 B/op | 0 allocs/op |
| Listener CID routing / 1024 connections / 255-byte CIDs / Miss | 5 | 15.32 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 64 bytes | 5 | 195.9 ns/op | - | 64 B/op | 1 allocs/op |
| Listener packet ingress / 1200 bytes | 5 | 154.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 2048 bytes | 5 | 158.5 ns/op | - | 0 B/op | 0 allocs/op |
| Listener packet ingress / 4096 bytes | 5 | 1.073 us/op | - | 4096 B/op | 1 allocs/op |
| Listener read deadline / Fixed deadline | 5 | 157.8 ns/op | - | 0 B/op | 0 allocs/op |
| Listener read deadline / Updated deadline | 5 | 316.2 ns/op | - | 0 B/op | 0 allocs/op |
| Plain ACK build / Empty | 5 | 28.73 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 35.14 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 366.3 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.47 us/op | 2787.16 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.217 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 161.7 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 140.9 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 680.6 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.717 us/op | 1507.44 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 241.2 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 71.01 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 270.2 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 16.1 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 18.5 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.82 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 1.63 us/op | 736.04 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 654.8 ns/op | 1832.7 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.01 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 1.646 us/op | 729.26 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.021 us/op | 1174.77 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 5.705 us/op | 210.33 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 1.586 us/op | 756.44 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 1.669 us/op | 719.17 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.422 us/op | 495.55 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 860.2 ns/op | 1394.96 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.536 us/op | 473.18 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 649.4 ns/op | 1847.99 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 686.4 ns/op | 1748.29 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.04 us/op | 1154.35 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
