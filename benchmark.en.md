# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- Generated: `2026-10-14T02:50:16Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `28bc7540dd420114dc33371e181b5617a4a66611 (Linux Release static)`
- NSS: `b9283c1e4ee232b4450d36115272af670d38ec81 (Linux Release)`
- BoringSSL: `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 (Linux Release)`
- OpenSSL: `41eecc7c7cd4069b1b29bca49b9d97825c6bf648 (Linux Release)`

347 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (34)](#section-connection-lifecycle)
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
- [Record layer and reliability (29)](#section-record-layer-and-reliability)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 554.561 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 830.806 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 429.473 us/op | - | 123190 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.057 ms/op | - | 123645 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.449 ms/op | - | 142424 B/op | 1403 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 653.118 us/op | - | 121345 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 894.479 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 895.344 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 356.07 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 999.569 us/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 991.952 us/op | - | 126045 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.651 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.667 ms/op | - | 162384 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 948.15 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 961.804 us/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 828.221 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 786.788 us/op | - | 152655 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.19 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 13.177 us/op | 9.71 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 15.876 us/op | 151.17 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 13.294 us/op | 9.63 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 15.94 us/op | 150.56 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 14.15 us/op | 9.05 MB/s | 594 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 16.713 us/op | 143.6 MB/s | 2994 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 14.98 us/op | 8.54 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 17.502 us/op | 137.13 MB/s | 2997 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 16.829 us/op | 7.61 MB/s | 645 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 20.848 us/op | 115.12 MB/s | 2999 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 24.342 us/op | 5.26 MB/s | 613 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 27.959 us/op | 85.84 MB/s | 3021 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 398.86 us/op | - | 109279 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 399.825 us/op | - | 109145 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.236 us/op | 39.55 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.497 us/op | 369.4 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.882 ms/conn | 2671704 B/op | 21758 allocs/op |
| 1-RTT application-data round trip | 5 | 1.883 ms/conn | 2682536 B/op | 22097 allocs/op |
| Full mTLS handshake | 5 | 3.421 ms/conn | 3375200 B/op | 30093 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.441 ms/conn | 3833384 B/op | 31838 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.899 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.4314 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.864 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.986 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.945 ms/conn | 4181448 B/op | 40363 allocs/op |
| Session resumption handshake | 5 | 0.5099 ms/conn | 4735608 B/op | 39431 allocs/op |
| mTLS session resumption handshake | 5 | 0.6313 ms/conn | 6867472 B/op | 51429 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.55 ms/conn | 299120 B/op | 1962 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.19 ms/conn | 3543216 B/op | 22082 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.151 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.543 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.982 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 5.153 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.821 ms/conn | 1495824 B/op | 11930 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.84 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.967 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.911 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.964 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.222 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.719 ms/conn | 1546704 B/op | 12667 allocs/op |
| Session resumption handshake | 5 | 1.063 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.031 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.398 ms/conn | 2006104 B/op | 11253 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.564 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.244 ms/conn | 1326768 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 4.675 ms/conn | 2480504 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 5.527 ms/conn | 1772032 B/op | 15063 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.427 ms/conn | 1974880 B/op | 15595 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.228 ms/conn | 1391904 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.78 ms/conn | 1048848 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.677 ms/conn | 2494984 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.759 ms/conn | 2656024 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.34 ms/conn | 3404504 B/op | 22848 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 3644096 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244184 B/op | 1014 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.469 ms/conn | 1628128 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.52 ms/conn | 1648288 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.557 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 6.743 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63632 B/op | 56 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.704 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.815 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.09 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64416 B/op | 56 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8744 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9631 ms/conn | 526920 B/op | 3537 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8759 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8165 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.107 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8805 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.091 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8881 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.261 ms/conn | 549352 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8658 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8767 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.054 ms/conn | 542072 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.5097 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.395 ms/pair | 884744 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 1.816 ms/pair | 1065472 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.416 ms/pair | 882312 B/op | 7303 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.235 ms/conn | 519480 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.849 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.131 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.233 ms/conn | 215344 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.125 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.702 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.273 ms/conn | 186640 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.796 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.235 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 1.885 ms/conn | 222672 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.152 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.12 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.109 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.704 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.095 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 2.907 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.287 ms/conn | 203184 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.394 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.451 ms/conn | 485640 B/op | 1993 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.514 ms/conn | 517288 B/op | 1995 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.417 ms/conn | 457576 B/op | 1975 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 3.968 ms/conn | 488392 B/op | 2013 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.527 ms/conn | 490792 B/op | 2049 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.097 ms/conn | 518568 B/op | 2189 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.301 ms/conn | 555912 B/op | 2119 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 3.432 ms/conn | 518888 B/op | 2175 allocs/op |
| Full mTLS handshake | 5 | 3.949 ms/conn | 528264 B/op | 2565 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.39 ms/conn | 487048 B/op | 2023 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.457 ms/conn | 520616 B/op | 2321 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.473 ms/conn | 510888 B/op | 2205 allocs/op |
| Direct external PSK handshake | 5 | 2.626 ms/conn | 464552 B/op | 1895 allocs/op |
| Session resumption handshake | 5 | 4.189 ms/pair | 712392 B/op | 3932 allocs/op |
| mTLS session resumption handshake | 5 | 4.848 ms/pair | 811496 B/op | 4734 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.363 ms/pair | 714240 B/op | 3883 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 4.235 ms/conn | 521640 B/op | 2353 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 5.069 ms/pair | 749224 B/op | 4298 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.753 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.758 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.659 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.688 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.491 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.894 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.215 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 4.528 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.622 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.625 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.627 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 2.789 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.478 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 6.35 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 4.292 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.094 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8982 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9363 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8806 ms/conn | 398008 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8085 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.101 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8762 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.1 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9033 ms/conn | 517120 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.24 ms/conn | 549352 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8667 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.9032 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.07 ms/conn | 542072 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.698 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.5086 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.399 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 1.817 ms/pair | 1065344 B/op | 8758 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.215 ms/conn | 658344 B/op | 5541 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.24 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.845 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7471 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.714 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.6968 ms/conn | 144928 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7201 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.38 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.724 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9679 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.6824 ms/conn | 180192 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.003 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.6556 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.6853 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9501 ms/conn | 229024 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6016 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4849 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.125 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.351 ms/pair | 473312 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.8808 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.7543 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.32 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.114 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.142 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.084 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.042 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.699 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.149 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.245 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.103 ms/conn | 514056 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.352 ms/conn | 521064 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.097 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.105 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.203 ms/conn | 506536 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.009 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.816 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 2.632 ms/pair | 701832 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 2.987 ms/pair | 798024 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.451 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.849 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.86 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.81 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.814 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.757 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.897 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.988 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.808 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 1.969 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.833 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.81 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.891 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.801 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.701 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.252 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.462 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.938 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8706 ms/conn | 460504 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9443 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8746 ms/conn | 398008 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8632 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.821 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.108 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8902 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.099 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8777 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.248 ms/conn | 548976 B/op | 4701 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8729 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8912 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.068 ms/conn | 542072 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.711 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.5092 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.383 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 1.826 ms/pair | 1065104 B/op | 8752 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.405 ms/pair | 882064 B/op | 7300 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.662 ms/conn | 653976 B/op | 6118 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9472 ms/conn | 210720 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.006 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9583 ms/conn | 175776 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.9895 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9965 ms/conn | 220064 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.474 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.014 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.216 ms/conn | 331552 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9416 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.27 ms/conn | 266688 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.948 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8851 ms/conn | 246816 B/op | 1985 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.144 ms/conn | 270368 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8423 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.7338 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.487 ms/pair | 470496 B/op | 3878 allocs/op |
| mTLS session resumption handshake | 5 | 2.166 ms/pair | 591280 B/op | 4441 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.186 ms/conn | 270080 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.35 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.452 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.292 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.3 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.305 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.732 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.416 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.704 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.325 ms/conn | 518632 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 4.652 ms/conn | 525704 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.431 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.474 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.57 ms/conn | 510984 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.786 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 4.077 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 5.198 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 5.538 ms/pair | 810088 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.168 ms/pair | 711992 B/op | 3781 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.055 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.499 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.477 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.444 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.393 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.508 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.933 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.576 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.762 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.482 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 4.684 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.476 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.516 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.655 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.942 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.231 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.313 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 5.907 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 4.747 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.23 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 45.4 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 533.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.047 us/op | 2000.83 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.656 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 198.9 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 160.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 893.3 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.441 us/op | 1190.28 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 330.7 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.74 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 352.4 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 27.96 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.91 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.331 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.281 us/op | 526.13 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 822 ns/op | 1459.92 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.24 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.404 us/op | 499.22 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.187 us/op | 1011.15 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.76 us/op | 177.51 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.268 us/op | 529.17 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.359 us/op | 508.68 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.473 us/op | 345.55 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 897.7 ns/op | 1336.71 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.877 us/op | 417.13 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 812.9 ns/op | 1476.11 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 876.2 ns/op | 1369.52 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.452 us/op | 826.57 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
