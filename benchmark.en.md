# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- Generated: `2026-10-14T13:41:07Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `782235be7cbb83f9551612a54cfc5f220d7b7d5f (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 587.766 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 887.29 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 455.714 us/op | - | 123196 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.133 ms/op | - | 123643 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.533 ms/op | - | 142426 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 686.541 us/op | - | 121345 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 939.481 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 937.336 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 375.011 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.055 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.048 ms/op | - | 126046 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.748 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.759 ms/op | - | 162384 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 994.885 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 1.004 ms/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 910.062 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 862.134 us/op | - | 152656 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.171 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 22.229 us/op | 5.76 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 25.789 us/op | 93.06 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 22.445 us/op | 5.7 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 26.276 us/op | 91.34 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 22.701 us/op | 5.64 MB/s | 598 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 25.947 us/op | 92.5 MB/s | 2997 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 23.515 us/op | 5.44 MB/s | 601 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 26.606 us/op | 90.2 MB/s | 2998 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 27.587 us/op | 4.64 MB/s | 623 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 31.152 us/op | 77.04 MB/s | 3032 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 36.74 us/op | 3.48 MB/s | 644 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 41.207 us/op | 58.24 MB/s | 2995 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 431.18 us/op | - | 109275 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 433.512 us/op | - | 109114 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.469 us/op | 36.9 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.984 us/op | 343.63 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.01 ms/conn | 2671704 B/op | 21758 allocs/op |
| 1-RTT application-data round trip | 5 | 2.035 ms/conn | 2682536 B/op | 22097 allocs/op |
| Full mTLS handshake | 5 | 3.592 ms/conn | 3375048 B/op | 30099 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.627 ms/conn | 3833600 B/op | 31833 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.999 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.55 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.977 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.125 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.112 ms/conn | 4181448 B/op | 40336 allocs/op |
| Session resumption handshake | 5 | 0.6689 ms/conn | 4735608 B/op | 39436 allocs/op |
| mTLS session resumption handshake | 5 | 0.8408 ms/conn | 6866976 B/op | 51418 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6852 ms/conn | 298872 B/op | 1959 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.356 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.334 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.569 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.75 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.89 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.417 ms/conn | 1495072 B/op | 11923 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.388 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.76 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.9872 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.701 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.942 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.243 ms/conn | 1546576 B/op | 12665 allocs/op |
| Session resumption handshake | 5 | 1.182 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.174 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.003 ms/conn | 2005400 B/op | 11246 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.026 ms/conn | 2018584 B/op | 11373 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.41 ms/conn | 1327216 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.359 ms/conn | 2480824 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.222 ms/conn | 1771680 B/op | 15069 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.134 ms/conn | 1974600 B/op | 15592 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.438 ms/conn | 1391712 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.914 ms/conn | 1048400 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.138 ms/conn | 2494360 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.275 ms/conn | 2655064 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.88 ms/conn | 3404232 B/op | 22846 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3643344 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244184 B/op | 1014 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.635 ms/conn | 1627808 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.67 ms/conn | 1648128 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.765 ms/conn | 1807552 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.998 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.721 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.824 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.78 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190552 B/op | 1203 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64416 B/op | 56 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9946 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.517 ms/conn | 526920 B/op | 3537 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.005 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9342 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.12 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.021 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.285 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.07 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.439 ms/conn | 549736 B/op | 4706 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.015 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.038 ms/conn | 496040 B/op | 3765 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.202 ms/conn | 542152 B/op | 4305 allocs/op |
| Direct external PSK handshake | 5 | 0.6491 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.723 ms/pair | 882968 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 2.128 ms/pair | 1065856 B/op | 8754 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.731 ms/pair | 882184 B/op | 7303 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.383 ms/conn | 519480 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.17 ms/pair | 961848 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.313 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.667 ms/conn | 215280 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.44 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.928 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.339 ms/conn | 186512 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.111 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.627 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.264 ms/conn | 222800 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.365 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.339 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.369 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.8979 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.426 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.314 ms/pair | 458928 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.504 ms/conn | 203184 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.608 ms/pair | 420208 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.384 ms/conn | 485640 B/op | 1993 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.309 ms/conn | 517256 B/op | 1993 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.351 ms/conn | 457576 B/op | 1975 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.844 ms/conn | 488424 B/op | 2015 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.313 ms/conn | 490856 B/op | 2053 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.046 ms/conn | 518568 B/op | 2189 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.871 ms/conn | 555880 B/op | 2117 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.547 ms/conn | 518856 B/op | 2173 allocs/op |
| Full mTLS handshake | 5 | 5.03 ms/conn | 528392 B/op | 2565 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.48 ms/conn | 487080 B/op | 2025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.36 ms/conn | 520648 B/op | 2323 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.345 ms/conn | 510600 B/op | 2203 allocs/op |
| Direct external PSK handshake | 5 | 3.464 ms/conn | 464520 B/op | 1893 allocs/op |
| Session resumption handshake | 5 | 5.503 ms/pair | 712392 B/op | 3932 allocs/op |
| mTLS session resumption handshake | 5 | 6.069 ms/pair | 811816 B/op | 4734 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.485 ms/pair | 714456 B/op | 3885 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.104 ms/conn | 521640 B/op | 2353 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.268 ms/pair | 749224 B/op | 4298 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.518 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.555 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.425 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.737 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.396 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.019 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.818 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.663 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.533 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.581 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.506 ms/conn | 212192 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.476 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.646 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.665 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.189 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.259 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.009 ms/conn | 460504 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.046 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.003 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.935 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.098 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.044 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.24 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.025 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.4 ms/conn | 549224 B/op | 4701 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.037 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.029 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.193 ms/conn | 542072 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8401 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6396 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.684 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.179 ms/pair | 1065728 B/op | 8758 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.385 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.354 ms/conn | 519480 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.139 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.763 ms/conn | 179872 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8288 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7565 ms/conn | 145056 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7424 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.421 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8545 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.006 ms/conn | 270496 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9049 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.119 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8608 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.7844 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.035 ms/conn | 229024 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7745 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.5191 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.552 ms/pair | 388016 B/op | 3280 allocs/op |
| mTLS session resumption handshake | 5 | 1.786 ms/pair | 473568 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.005 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.94 ms/conn | 200800 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.556 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.55 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.562 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.481 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.476 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.204 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.468 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.642 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.47 ms/conn | 514184 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.977 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.506 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.515 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.536 ms/conn | 506536 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.386 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.184 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 3.176 ms/pair | 701848 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 3.548 ms/pair | 798168 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.901 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.163 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.278 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.131 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.161 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.224 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.259 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.298 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.153 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.593 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.159 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.126 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.23 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.191 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.019 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.732 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.895 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.281 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.005 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.047 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.02 ms/conn | 398136 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.016 ms/conn | 525160 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.965 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.121 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9876 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.256 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.031 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.461 ms/conn | 549224 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.024 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.032 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.198 ms/conn | 541048 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8546 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.6375 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.648 ms/pair | 882248 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.193 ms/pair | 1065096 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.711 ms/pair | 882184 B/op | 7303 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.919 ms/conn | 653728 B/op | 6120 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.063 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.167 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.129 ms/conn | 175904 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.216 ms/conn | 264352 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.107 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.502 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.277 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.465 ms/conn | 331440 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.123 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.488 ms/conn | 266688 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.191 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.18 ms/conn | 247104 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.496 ms/conn | 269792 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.053 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.9226 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.987 ms/pair | 470464 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 2.827 ms/pair | 590128 B/op | 4440 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.342 ms/conn | 270080 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.195 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.206 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.273 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.187 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.178 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.566 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.204 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.414 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.521 ms/conn | 518696 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 5.769 ms/conn | 525704 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.856 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.435 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.498 ms/conn | 510472 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.799 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 5.092 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 6.214 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 6.65 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.993 ms/pair | 711744 B/op | 3784 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.943 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.227 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.462 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.189 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.167 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.252 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.724 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.289 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.614 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.399 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.553 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.617 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.435 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.488 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.918 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.996 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.296 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 6.936 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.514 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.86 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 45.03 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 523.7 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.057 us/op | 1991.47 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.699 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 202.4 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 170.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 928 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.605 us/op | 1136.07 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 329.8 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.86 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 345.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.53 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.04 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.372 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.222 us/op | 539.94 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 870.7 ns/op | 1378.28 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.3 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.32 us/op | 517.31 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.349 us/op | 889.84 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.066 us/op | 169.84 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.188 us/op | 548.51 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.307 us/op | 520.26 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.037 us/op | 395.19 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 929.1 ns/op | 1291.58 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.037 us/op | 395.14 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 868.3 ns/op | 1381.95 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 922.7 ns/op | 1300.58 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.294 us/op | 927.16 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
