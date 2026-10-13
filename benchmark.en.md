# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c5e4bf75e5d3c6f6367cc7ebd1b55d11e81ab5a3`
- Generated: `2026-10-13T04:24:38Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `6d48df25bca7efef49eb62816ae8080bbf6b96e3 (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 587.48 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 874.395 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 453.986 us/op | - | 123190 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.113 ms/op | - | 123644 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.525 ms/op | - | 142428 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 684.988 us/op | - | 121350 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 929.922 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 936.17 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 369.961 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.034 ms/op | - | 136090 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.033 ms/op | - | 126045 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.712 ms/op | - | 177152 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.721 ms/op | - | 162383 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 988.517 us/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 989.554 us/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 871.885 us/op | - | 149310 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 823.457 us/op | - | 152655 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.111 ms/op | - | 178102 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 22.386 us/op | 5.72 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 25.111 us/op | 95.58 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 22.636 us/op | 5.65 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 25.844 us/op | 92.86 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 22.728 us/op | 5.63 MB/s | 598 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 25.513 us/op | 94.07 MB/s | 2996 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 23.421 us/op | 5.47 MB/s | 599 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 26.362 us/op | 91.04 MB/s | 2997 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 27.002 us/op | 4.74 MB/s | 621 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 30.492 us/op | 78.71 MB/s | 3031 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 35.667 us/op | 3.59 MB/s | 649 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 40.861 us/op | 58.74 MB/s | 2995 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 425.4 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 422.855 us/op | - | 109114 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.337 us/op | 38.36 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.61 us/op | 363.1 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.982 ms/conn | 2671752 B/op | 21759 allocs/op |
| 1-RTT application-data round trip | 5 | 2.016 ms/conn | 2682536 B/op | 22097 allocs/op |
| Full mTLS handshake | 5 | 3.571 ms/conn | 3373712 B/op | 30082 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.587 ms/conn | 3832808 B/op | 31829 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.011 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.5547 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.994 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.143 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.104 ms/conn | 4181496 B/op | 40357 allocs/op |
| Session resumption handshake | 5 | 0.6542 ms/conn | 4735608 B/op | 39436 allocs/op |
| mTLS session resumption handshake | 5 | 0.8265 ms/conn | 6867968 B/op | 51429 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.7676 ms/conn | 299120 B/op | 1961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.424 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.369 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.683 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.837 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.954 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.374 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.433 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.762 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 1.02 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.825 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.011 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.36 ms/conn | 1546512 B/op | 12664 allocs/op |
| Session resumption handshake | 5 | 1.227 ms/conn | 2367584 B/op | 19299 allocs/op |
| mTLS session resumption handshake | 5 | 1.227 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.103 ms/conn | 2006888 B/op | 11261 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.093 ms/conn | 2018584 B/op | 11373 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.349 ms/conn | 1326496 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.08 ms/conn | 2480248 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.048 ms/conn | 1771432 B/op | 15072 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.997 ms/conn | 1974152 B/op | 15592 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.349 ms/conn | 1391504 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.936 ms/conn | 1048624 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.087 ms/conn | 2494744 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.125 ms/conn | 2656024 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.79 ms/conn | 3404224 B/op | 22837 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3643024 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244312 B/op | 1015 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.605 ms/conn | 1628288 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.678 ms/conn | 1648288 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.961 ms/conn | 1807168 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.557 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.168 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 8.629 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.224 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.585 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.154 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.952 ms/conn | 1183560 B/op | 1203 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.131 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.81 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190552 B/op | 1203 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.971 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.25 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.07 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.03 ms/conn | 460552 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.055 ms/conn | 526792 B/op | 3537 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.033 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9969 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.166 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.062 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.291 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.089 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.468 ms/conn | 549848 B/op | 4706 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9984 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.081 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.238 ms/conn | 541048 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.6192 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.739 ms/pair | 882968 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 2.216 ms/pair | 1065096 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.739 ms/pair | 881984 B/op | 7301 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.392 ms/conn | 519480 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.165 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.521 ms/conn | 181472 B/op | 1535 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.54 ms/conn | 215024 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.377 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.924 ms/conn | 184432 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.491 ms/conn | 186640 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.23 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.616 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.377 ms/conn | 222672 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.3 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.424 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.506 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.8707 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.485 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.72 ms/pair | 458608 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.565 ms/conn | 203184 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.709 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.484 ms/conn | 485640 B/op | 1993 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.617 ms/conn | 517256 B/op | 1993 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.525 ms/conn | 457576 B/op | 1975 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.155 ms/conn | 488424 B/op | 2015 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.859 ms/conn | 490856 B/op | 2051 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.306 ms/conn | 518504 B/op | 2185 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.986 ms/conn | 555912 B/op | 2119 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.633 ms/conn | 518888 B/op | 2175 allocs/op |
| Full mTLS handshake | 5 | 5.312 ms/conn | 528264 B/op | 2563 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.608 ms/conn | 487080 B/op | 2025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.724 ms/conn | 520616 B/op | 2321 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.607 ms/conn | 510376 B/op | 2205 allocs/op |
| Direct external PSK handshake | 5 | 3.753 ms/conn | 464520 B/op | 1893 allocs/op |
| Session resumption handshake | 5 | 5.941 ms/pair | 712360 B/op | 3930 allocs/op |
| mTLS session resumption handshake | 5 | 6.362 ms/pair | 811528 B/op | 4736 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.88 ms/pair | 714304 B/op | 3888 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.663 ms/conn | 521608 B/op | 2351 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.439 ms/pair | 749224 B/op | 4298 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.648 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.704 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.802 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.984 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.826 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.242 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.883 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.842 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.791 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.877 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.707 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.816 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.995 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.942 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.532 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.728 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.051 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.075 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.027 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9903 ms/conn | 467288 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.159 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.079 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.296 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.079 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.483 ms/conn | 548976 B/op | 4701 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.058 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.045 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.23 ms/conn | 542200 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.864 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6334 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.73 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.172 ms/pair | 1065480 B/op | 8755 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.396 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.403 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.255 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8604 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9268 ms/conn | 213024 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7844 ms/conn | 145056 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8384 ms/conn | 182880 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.633 ms/conn | 184960 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9956 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 0.9997 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9558 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.036 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8556 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8345 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.06 ms/conn | 229152 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7511 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.5962 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.557 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.833 ms/pair | 473440 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.9577 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.9147 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.848 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.778 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.834 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.656 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.68 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.581 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.801 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.805 ms/conn | 546088 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.675 ms/conn | 514120 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.871 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.834 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.663 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.71 ms/conn | 506536 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.594 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.347 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 3.384 ms/pair | 701864 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 3.96 ms/pair | 798040 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 3.017 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.219 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.275 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.22 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.231 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.68 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.328 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.432 ms/conn | 212512 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.235 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.376 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.215 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.194 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.32 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.189 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.143 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.758 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 3.038 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.375 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9921 ms/conn | 460504 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.105 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9854 ms/conn | 398136 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.068 ms/conn | 525160 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9871 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.161 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.075 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.26 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.104 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.437 ms/conn | 549472 B/op | 4706 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.055 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.11 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.222 ms/conn | 542200 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8379 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.7092 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.7 ms/pair | 882248 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.195 ms/pair | 1065096 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.751 ms/pair | 882312 B/op | 7303 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.966 ms/conn | 653480 B/op | 6118 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.176 ms/conn | 210720 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.231 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.09 ms/conn | 175904 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.318 ms/conn | 264352 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.325 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.581 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.276 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.455 ms/conn | 331424 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.296 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.612 ms/conn | 266688 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.113 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.333 ms/conn | 247104 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.425 ms/conn | 270496 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.161 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.949 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 1.975 ms/pair | 470464 B/op | 3877 allocs/op |
| mTLS session resumption handshake | 5 | 2.847 ms/pair | 590032 B/op | 4439 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.41 ms/conn | 270080 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.795 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 6.132 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.954 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.866 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.835 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.321 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.135 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.007 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.745 ms/conn | 518632 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 6.402 ms/conn | 525704 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.918 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.054 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.539 ms/conn | 510984 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.422 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 5.485 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 6.827 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 6.998 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.73 ms/pair | 711744 B/op | 3784 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.758 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.552 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.883 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.859 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.765 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.811 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.478 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.757 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.315 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.422 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 6.046 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.725 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.569 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.982 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.364 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 5.656 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.913 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 7.24 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 6.24 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.52 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.48 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 503.9 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.954 us/op | 2096.38 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.662 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 202.1 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 166.8 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 910.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.505 us/op | 1168.75 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 313.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.49 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 336.2 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.53 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 24.2 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.376 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.164 us/op | 554.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 847.5 ns/op | 1415.96 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.28 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.331 us/op | 514.7 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.357 us/op | 884.27 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.92 us/op | 173.4 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.114 us/op | 567.62 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.217 us/op | 541.25 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.953 us/op | 406.37 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.038 us/op | 1156.6 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.976 us/op | 403.18 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 861.6 ns/op | 1392.75 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 914.8 ns/op | 1311.79 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.262 us/op | 950.56 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
