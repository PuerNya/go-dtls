# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- Generated: `2026-10-13T21:29:36Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `460292b036b8e22516b7a160af8966af90951951 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 478.588 us/op | - | 101321 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 724.809 us/op | - | 116336 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 371.107 us/op | - | 123189 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 935.315 us/op | - | 123646 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.25 ms/op | - | 142415 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 564.533 us/op | - | 121347 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 768.236 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 766.966 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 302.938 us/op | - | 103775 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 841.473 us/op | - | 136094 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 843.844 us/op | - | 126049 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.414 ms/op | - | 177157 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.429 ms/op | - | 162388 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 838.991 us/op | - | 151404 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 838.947 us/op | - | 154181 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 737.991 us/op | - | 149314 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 703.491 us/op | - | 152658 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.82 ms/op | - | 178101 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 18.49 us/op | 6.92 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 21.012 us/op | 114.22 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 18.709 us/op | 6.84 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 21.14 us/op | 113.53 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 19.233 us/op | 6.66 MB/s | 595 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 21.309 us/op | 112.63 MB/s | 2993 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 19.566 us/op | 6.54 MB/s | 596 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 21.759 us/op | 110.3 MB/s | 2995 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 24.584 us/op | 5.21 MB/s | 611 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 27.463 us/op | 87.39 MB/s | 3020 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 29.465 us/op | 4.34 MB/s | 628 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 31.844 us/op | 75.37 MB/s | 3034 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 337.119 us/op | - | 109246 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 338.71 us/op | - | 109117 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 2.49 us/op | 51.4 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 5.097 us/op | 470.9 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.657 ms/conn | 2671704 B/op | 21758 allocs/op |
| 1-RTT application-data round trip | 5 | 1.659 ms/conn | 2682536 B/op | 22097 allocs/op |
| Full mTLS handshake | 5 | 2.972 ms/conn | 3374704 B/op | 30091 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 2.995 ms/conn | 3833056 B/op | 31836 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.676 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.4321 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.669 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.833 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.445 ms/conn | 4180456 B/op | 40360 allocs/op |
| Session resumption handshake | 5 | 0.5209 ms/conn | 4735624 B/op | 39430 allocs/op |
| mTLS session resumption handshake | 5 | 0.6241 ms/conn | 6867224 B/op | 51423 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5425 ms/conn | 299120 B/op | 1962 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.993 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.955 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.022 ms/conn | 4080784 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.931 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.007 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 5.186 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.147 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.942 ms/conn | 1571392 B/op | 12603 allocs/op |
| Direct external PSK handshake | 5 | 0.8604 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.023 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.156 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.242 ms/conn | 1546832 B/op | 12669 allocs/op |
| Session resumption handshake | 5 | 0.9458 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 0.9153 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.201 ms/conn | 2005320 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.019 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.023 ms/conn | 1327168 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 4.238 ms/conn | 2481000 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 5.042 ms/conn | 1771080 B/op | 15070 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.118 ms/conn | 1974400 B/op | 15589 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.993 ms/conn | 1392128 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.765 ms/conn | 1048400 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.585 ms/conn | 2494744 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.775 ms/conn | 2655704 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.21 ms/conn | 3403016 B/op | 22834 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 3643184 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 244432 B/op | 1017 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.197 ms/conn | 1627808 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.048 ms/conn | 1647968 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.65 ms/conn | 1806912 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 6.454 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.911 ms/conn | 1183512 B/op | 1203 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.997 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.995 ms/conn | 1172800 B/op | 1206 allocs/op |
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
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8259 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8426 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8605 ms/conn | 398056 B/op | 3501 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7852 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.787 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8192 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.072 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.846 ms/conn | 517120 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.162 ms/conn | 549976 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8202 ms/conn | 462680 B/op | 3585 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8817 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.03 ms/conn | 542200 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.5575 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.426 ms/pair | 884584 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 1.779 ms/pair | 1065224 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.412 ms/pair | 881984 B/op | 7301 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.201 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.709 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.214 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.22 ms/conn | 215344 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.194 ms/conn | 146704 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.651 ms/conn | 184432 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.868 ms/conn | 186640 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.782 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.192 ms/conn | 274736 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 1.96 ms/conn | 222800 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.173 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.256 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.153 ms/conn | 194416 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.716 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.102 ms/pair | 378896 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 2.899 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.27 ms/conn | 203312 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.334 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.57 ms/conn | 485576 B/op | 1989 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.37 ms/conn | 517224 B/op | 1991 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.156 ms/conn | 457608 B/op | 1977 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.149 ms/conn | 488392 B/op | 2013 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 6.395 ms/conn | 490856 B/op | 2053 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.192 ms/conn | 518504 B/op | 2185 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.517 ms/conn | 555880 B/op | 2117 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.719 ms/conn | 518632 B/op | 2175 allocs/op |
| Full mTLS handshake | 5 | 4.339 ms/conn | 528328 B/op | 2566 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.245 ms/conn | 487080 B/op | 2025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.833 ms/conn | 520584 B/op | 2319 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.528 ms/conn | 510888 B/op | 2205 allocs/op |
| Direct external PSK handshake | 5 | 3.286 ms/conn | 464520 B/op | 1893 allocs/op |
| Session resumption handshake | 5 | 5.282 ms/pair | 712392 B/op | 3932 allocs/op |
| mTLS session resumption handshake | 5 | 5.358 ms/pair | 811528 B/op | 4736 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.5 ms/pair | 714336 B/op | 3885 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 4.965 ms/conn | 521656 B/op | 2352 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 5.276 ms/pair | 749160 B/op | 4294 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.547 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.424 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.286 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.398 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.421 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.635 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.379 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.646 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.408 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.654 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.581 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.633 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.378 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.18 ms/pair | 207520 B/op | 185 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 4.889 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.756 ms/pair | 207776 B/op | 185 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8031 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.875 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8664 ms/conn | 398008 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7633 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.77 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8029 ms/conn | 535520 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.02 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.8382 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.147 ms/conn | 549096 B/op | 4703 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8095 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8197 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9877 ms/conn | 541048 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6575 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4798 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.323 ms/pair | 882968 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 1.708 ms/pair | 1065224 B/op | 8755 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.117 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.114 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.827 ms/pair | 961848 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.6603 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.7705 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.6602 ms/conn | 144928 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.6329 ms/conn | 182880 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.11 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.6981 ms/conn | 203552 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.8678 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.6908 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 0.8643 ms/conn | 219872 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.6976 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.6827 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.939 ms/conn | 228512 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.5725 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.4636 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.198 ms/pair | 388016 B/op | 3280 allocs/op |
| mTLS session resumption handshake | 5 | 1.445 ms/pair | 473440 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.9758 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.7537 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.398 ms/pair | 431488 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.159 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.756 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.231 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.126 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.597 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.157 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.327 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.545 ms/conn | 514120 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.367 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.753 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.054 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.333 ms/conn | 506024 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.039 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.931 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 2.773 ms/pair | 701848 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 3.699 ms/pair | 798040 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.778 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.156 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.615 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.954 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.911 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.678 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.95 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.954 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.328 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.512 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.608 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.235 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.025 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.944 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 1.85 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.315 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.928 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.552 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7966 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8991 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8556 ms/conn | 398136 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.8221 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7549 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 1.891 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.8315 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.025 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.909 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.139 ms/conn | 548976 B/op | 4701 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8381 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.827 ms/conn | 493032 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.96 ms/conn | 542072 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6636 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.5071 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.31 ms/pair | 882248 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 1.706 ms/pair | 1065472 B/op | 8758 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.328 ms/pair | 882064 B/op | 7300 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.547 ms/conn | 652984 B/op | 6112 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.966 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.029 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.133 ms/conn | 175904 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.9036 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9942 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.062 ms/conn | 219872 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9951 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.24 ms/conn | 331552 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.285 ms/conn | 266688 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.003 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.9708 ms/conn | 246976 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.22 ms/conn | 270368 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7729 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.7222 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.58 ms/pair | 470464 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 2.048 ms/pair | 590256 B/op | 4440 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.021 ms/conn | 270824 B/op | 2313 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.952 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.612 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.374 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.369 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.432 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.822 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.499 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.62 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.363 ms/conn | 518696 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 4.68 ms/conn | 525704 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.457 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.603 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.232 ms/conn | 510984 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.674 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 4.078 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 5.167 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 5.744 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.076 ms/pair | 711992 B/op | 3787 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.265 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.81 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.699 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.524 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.513 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.612 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.668 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.3 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 4.933 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.974 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 4.979 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.062 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.897 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.72 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 4.931 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.644 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 5.472 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 5.926 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 4.87 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 27.67 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 34.39 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 366.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.454 us/op | 2817.26 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.21 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 159.3 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 136.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 682.6 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.714 us/op | 1509.41 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 234.7 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 68.7 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 266.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 17.77 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 18.22 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.547 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 1.714 us/op | 700.04 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 684.7 ns/op | 1752.64 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.19 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 1.682 us/op | 713.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.032 us/op | 1162.38 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 5.735 us/op | 209.23 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 1.645 us/op | 729.32 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 1.744 us/op | 688.18 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.453 us/op | 489.23 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 734 ns/op | 1634.87 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.526 us/op | 475.12 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 661.7 ns/op | 1813.49 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 699.2 ns/op | 1716.36 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.061 us/op | 1131.31 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
