# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- Generated: `2026-10-14T20:18:31Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `e9a8f62d6ea0d196e386ce80b7d4df24160835ec (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 441.807 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 670.655 us/op | - | 116336 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 361.224 us/op | - | 123189 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 841.708 us/op | - | 123645 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.165 ms/op | - | 142416 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 539.632 us/op | - | 121345 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 730.423 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 727.562 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 290.189 us/op | - | 103775 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 770.526 us/op | - | 136094 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 770.259 us/op | - | 126017 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.282 ms/op | - | 177158 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.285 ms/op | - | 162355 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 758.667 us/op | - | 151405 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 796.719 us/op | - | 154181 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 655.019 us/op | - | 149313 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 638.739 us/op | - | 152657 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.648 ms/op | - | 178099 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 10.652 us/op | 12.02 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 12.77 us/op | 187.94 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 10.991 us/op | 11.65 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 13.09 us/op | 183.35 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 11.498 us/op | 11.13 MB/s | 593 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 14.634 us/op | 164 MB/s | 2993 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 12.45 us/op | 10.28 MB/s | 593 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 14.944 us/op | 160.6 MB/s | 2993 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 16.556 us/op | 7.73 MB/s | 608 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 20.94 us/op | 114.61 MB/s | 3000 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 24.205 us/op | 5.29 MB/s | 612 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 26.597 us/op | 90.24 MB/s | 3018 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 332.203 us/op | - | 109247 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 338.606 us/op | - | 109120 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 2.881 us/op | 44.43 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 5.586 us/op | 429.66 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.488 ms/conn | 2671704 B/op | 21758 allocs/op |
| 1-RTT application-data round trip | 5 | 1.444 ms/conn | 2682536 B/op | 22097 allocs/op |
| Full mTLS handshake | 5 | 2.579 ms/conn | 3374952 B/op | 30097 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 2.572 ms/conn | 3833552 B/op | 31839 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.365 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.3119 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.502 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.546 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.683 ms/conn | 4180704 B/op | 40355 allocs/op |
| Session resumption handshake | 5 | 0.4187 ms/conn | 4735608 B/op | 39420 allocs/op |
| mTLS session resumption handshake | 5 | 0.5158 ms/conn | 6867128 B/op | 51413 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.45 ms/conn | 299120 B/op | 1962 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.628 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.608 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.62 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.518 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 3.56 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 4.67 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.537 ms/conn | 1495224 B/op | 11925 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.428 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 0.883 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.587 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.712 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.768 ms/conn | 1547088 B/op | 12673 allocs/op |
| Session resumption handshake | 5 | 0.9762 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 0.9259 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.661 ms/conn | 2005320 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.475 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.714 ms/conn | 1326944 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 3.966 ms/conn | 2480888 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 4.428 ms/conn | 1770640 B/op | 15066 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.352 ms/conn | 1974624 B/op | 15593 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.741 ms/conn | 1391488 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.651 ms/conn | 1048624 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.159 ms/conn | 2494744 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.217 ms/conn | 2655384 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.329 ms/conn | 3404600 B/op | 22843 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 3642896 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 244184 B/op | 1014 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.885 ms/conn | 1627968 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.649 ms/conn | 1648288 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.07 ms/conn | 1807040 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 5.72 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63680 B/op | 57 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63760 B/op | 57 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.927 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.032 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 8.291 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1012 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7234 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.7615 ms/conn | 526792 B/op | 3537 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.69 ms/conn | 398008 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.6948 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.655 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7818 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9013 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7181 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 0.9992 ms/conn | 549608 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.696 ms/conn | 462504 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.7153 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.859 ms/conn | 541048 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.4362 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.117 ms/pair | 884744 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 1.593 ms/pair | 1065472 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.182 ms/pair | 882064 B/op | 7299 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.089 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.597 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9971 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.059 ms/conn | 215024 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.035 ms/conn | 146704 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.467 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.863 ms/conn | 186512 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.846 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.005 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 1.61 ms/conn | 222800 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9681 ms/conn | 181712 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.063 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.018 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.6182 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 1.723 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 2.695 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.126 ms/conn | 203312 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.256 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.149 ms/conn | 485672 B/op | 1995 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.252 ms/conn | 517288 B/op | 1995 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.206 ms/conn | 457608 B/op | 1977 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 3.682 ms/conn | 488456 B/op | 2017 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.143 ms/conn | 490824 B/op | 2051 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.034 ms/conn | 518536 B/op | 2187 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.138 ms/conn | 555944 B/op | 2121 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 3.295 ms/conn | 518696 B/op | 2179 allocs/op |
| Full mTLS handshake | 5 | 3.601 ms/conn | 528360 B/op | 2567 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.163 ms/conn | 487112 B/op | 2027 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.296 ms/conn | 520616 B/op | 2321 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.449 ms/conn | 510344 B/op | 2203 allocs/op |
| Direct external PSK handshake | 5 | 2.554 ms/conn | 464552 B/op | 1895 allocs/op |
| Session resumption handshake | 5 | 3.927 ms/pair | 712424 B/op | 3934 allocs/op |
| mTLS session resumption handshake | 5 | 4.839 ms/pair | 811496 B/op | 4734 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.204 ms/pair | 714208 B/op | 3883 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 4.178 ms/conn | 521576 B/op | 2349 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 4.753 ms/pair | 749288 B/op | 4302 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 3.362 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 3.487 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 3.345 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.405 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.56 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.224 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 4.193 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 3.417 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.395 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 3.458 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 2.644 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.412 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 5.969 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 3.98 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 5.695 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7622 ms/conn | 460504 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.786 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7708 ms/conn | 398008 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7105 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.688 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8051 ms/conn | 535520 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9936 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7757 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.015 ms/conn | 549472 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.7531 ms/conn | 462504 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.7179 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.8811 ms/conn | 541944 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6036 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4119 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.185 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 1.585 ms/pair | 1065472 B/op | 8758 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.107 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.103 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.634 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.6009 ms/conn | 179872 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.6841 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7009 ms/conn | 144928 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.5991 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.908 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.703 ms/conn | 203552 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.7812 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7032 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 0.9082 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.6609 ms/conn | 180064 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.6148 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.806 ms/conn | 228896 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.5952 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4204 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.004 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.252 ms/pair | 473440 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.8042 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.6531 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.12 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.116 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.031 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.096 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.131 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.353 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.189 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.302 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.145 ms/conn | 514184 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.216 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.084 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.045 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.222 ms/conn | 506536 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.013 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.761 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 2.488 ms/pair | 701848 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 2.877 ms/pair | 798024 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.411 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.846 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.906 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.84 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.842 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.271 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.942 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.994 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.869 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 1.943 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.789 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.71 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.824 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.854 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.725 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.237 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.371 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7336 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.7482 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.6968 ms/conn | 398008 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.6908 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.6851 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.573 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.7517 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.8927 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7444 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.103 ms/conn | 549096 B/op | 4703 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.6933 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.7266 ms/conn | 493032 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9048 ms/conn | 542072 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.5345 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.4187 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.099 ms/pair | 882248 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 1.429 ms/pair | 1065736 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.121 ms/pair | 882376 B/op | 7303 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.457 ms/conn | 653728 B/op | 6120 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.892 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.868 ms/conn | 244768 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8746 ms/conn | 175776 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.844 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8468 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.648 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8469 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.147 ms/conn | 331456 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.921 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.177 ms/conn | 266560 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8858 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.9099 ms/conn | 246880 B/op | 1987 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.062 ms/conn | 270368 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8145 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.6149 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.462 ms/pair | 471312 B/op | 3881 allocs/op |
| mTLS session resumption handshake | 5 | 1.899 ms/pair | 590240 B/op | 4439 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 0.9986 ms/conn | 270824 B/op | 2313 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.323 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.125 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.172 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.132 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.222 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 5.86 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.39 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.453 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.136 ms/conn | 518632 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 4.453 ms/conn | 525704 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.125 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.252 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.306 ms/conn | 510984 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.542 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 3.894 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 4.791 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 5.226 ms/pair | 810088 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 4.781 ms/pair | 711992 B/op | 3787 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.786 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.287 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.275 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.391 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.326 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.411 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.888 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.438 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.604 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.367 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 4.468 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.23 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.474 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.465 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.626 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.091 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 4.894 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 5.546 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 4.536 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 30.25 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 37.59 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 496.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.944 us/op | 2106.6 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.46 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 168.2 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 136.4 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 809.8 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.998 us/op | 1366.44 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 292.3 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 90.83 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 299.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 24.49 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 21.58 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.67 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.037 us/op | 589.12 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 771.3 ns/op | 1555.75 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.7 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.102 us/op | 571.02 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.094 us/op | 1097.23 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.182 us/op | 194.11 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.068 us/op | 580.15 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.14 us/op | 560.62 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.203 us/op | 374.6 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 813.7 ns/op | 1474.67 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.63 us/op | 456.32 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 754.9 ns/op | 1589.67 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 830 ns/op | 1445.86 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.356 us/op | 884.97 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
