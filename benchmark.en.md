# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c5e4bf75e5d3c6f6367cc7ebd1b55d11e81ab5a3`
- Generated: `2026-10-12T15:14:42Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `e65388ddcfe321fd6a54ab5bc5eb079753780aad (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 595.048 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 899.051 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 461.616 us/op | - | 123193 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.115 ms/op | - | 123644 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.528 ms/op | - | 142420 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 691.607 us/op | - | 121344 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 945.655 us/op | - | 125649 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 939.02 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 386.137 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.049 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.057 ms/op | - | 126045 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.74 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.752 ms/op | - | 162447 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.007 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 1.024 ms/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 893.857 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 842.816 us/op | - | 152655 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.16 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 22.287 us/op | 5.74 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 25.699 us/op | 93.39 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 22.614 us/op | 5.66 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 26.172 us/op | 91.7 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 22.891 us/op | 5.59 MB/s | 598 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 25.851 us/op | 92.84 MB/s | 2995 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 23.584 us/op | 5.43 MB/s | 600 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 27.015 us/op | 88.84 MB/s | 2999 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 28.848 us/op | 4.44 MB/s | 627 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 32.195 us/op | 74.55 MB/s | 3035 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 37.321 us/op | 3.43 MB/s | 624 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 41.869 us/op | 57.32 MB/s | 2997 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 425.644 us/op | - | 109308 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 429.208 us/op | - | 109114 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.491 us/op | 36.66 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.818 us/op | 352.02 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.987 ms/conn | 2671656 B/op | 21757 allocs/op |
| 1-RTT application-data round trip | 5 | 2.029 ms/conn | 2682584 B/op | 22098 allocs/op |
| Full mTLS handshake | 5 | 3.592 ms/conn | 3375792 B/op | 30108 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.617 ms/conn | 3833304 B/op | 31832 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.106 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.5571 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.015 ms/conn | 2689160 B/op | 22877 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.142 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.113 ms/conn | 4181200 B/op | 40348 allocs/op |
| Session resumption handshake | 5 | 0.6388 ms/conn | 4735608 B/op | 39436 allocs/op |
| mTLS session resumption handshake | 5 | 0.8112 ms/conn | 6867472 B/op | 51423 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6463 ms/conn | 299120 B/op | 1961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.375 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.324 ms/conn | 3585104 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.559 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.78 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.86 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.355 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.395 ms/conn | 1495824 B/op | 11930 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.756 ms/conn | 1573536 B/op | 12619 allocs/op |
| Direct external PSK handshake | 5 | 1.026 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.773 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.92 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.132 ms/conn | 1546640 B/op | 12666 allocs/op |
| Session resumption handshake | 5 | 1.19 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.175 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.003 ms/conn | 2006024 B/op | 11252 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.04 ms/conn | 2018504 B/op | 11372 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.372 ms/conn | 1326528 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.062 ms/conn | 2480824 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.033 ms/conn | 1771184 B/op | 15066 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.002 ms/conn | 1973984 B/op | 15592 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.378 ms/conn | 1392176 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.921 ms/conn | 1049296 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.232 ms/conn | 2495128 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.084 ms/conn | 2655384 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.73 ms/conn | 3404240 B/op | 22837 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3642608 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244504 B/op | 1015 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.586 ms/conn | 1628176 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.559 ms/conn | 1648128 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.76 ms/conn | 1807424 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.918 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.193 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 8.609 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.644 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.655 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.172 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.334 ms/conn | 1183560 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.997 ms/conn | 1169728 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.46 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191168 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190528 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.571 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.436 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.71 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.003 ms/conn | 460760 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.085 ms/conn | 526920 B/op | 3537 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.9929 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9629 ms/conn | 467336 B/op | 3585 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.115 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.007 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.241 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.042 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.447 ms/conn | 549096 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.027 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.034 ms/conn | 495528 B/op | 3764 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.189 ms/conn | 542072 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.6386 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.689 ms/pair | 882840 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 2.151 ms/pair | 1065224 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.704 ms/pair | 881936 B/op | 7300 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.357 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.173 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.315 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.319 ms/conn | 215344 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.421 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.957 ms/conn | 184432 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.505 ms/conn | 186512 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.021 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.446 ms/conn | 274864 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.287 ms/conn | 222480 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.264 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.277 ms/conn | 213328 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.285 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.7751 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.356 ms/pair | 378896 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.269 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.481 ms/conn | 203312 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.781 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.234 ms/conn | 485640 B/op | 1992 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.35 ms/conn | 517224 B/op | 1991 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.158 ms/conn | 457608 B/op | 1977 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 4.749 ms/conn | 488456 B/op | 2017 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.334 ms/conn | 490856 B/op | 2053 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 4.891 ms/conn | 518568 B/op | 2189 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.318 ms/conn | 555944 B/op | 2121 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.206 ms/conn | 518664 B/op | 2177 allocs/op |
| Full mTLS handshake | 5 | 4.776 ms/conn | 528168 B/op | 2565 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.307 ms/conn | 487016 B/op | 2021 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.335 ms/conn | 520648 B/op | 2323 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.273 ms/conn | 510888 B/op | 2205 allocs/op |
| Direct external PSK handshake | 5 | 3.419 ms/conn | 464488 B/op | 1891 allocs/op |
| Session resumption handshake | 5 | 5.286 ms/pair | 712360 B/op | 3930 allocs/op |
| mTLS session resumption handshake | 5 | 5.822 ms/pair | 811640 B/op | 4742 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.44 ms/pair | 714240 B/op | 3884 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.152 ms/conn | 521576 B/op | 2349 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.179 ms/pair | 749096 B/op | 4290 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.358 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.48 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.418 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.601 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.49 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.806 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.411 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.445 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.375 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.421 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.46 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.45 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.492 ms/pair | 207008 B/op | 184 allocs/op |
| mTLS session resumption handshake | 5 | 7.597 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.116 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.244 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.008 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.043 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.015 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9329 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.074 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.006 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.284 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.08 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.418 ms/conn | 549480 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.006 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.014 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.192 ms/conn | 541048 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8111 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6438 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.689 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.184 ms/pair | 1064848 B/op | 8752 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.369 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.353 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.134 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.8601 ms/conn | 180000 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.8923 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.8406 ms/conn | 145056 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.8063 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.498 ms/conn | 184960 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9451 ms/conn | 203680 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.061 ms/conn | 270624 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.7853 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.084 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.7841 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8039 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9569 ms/conn | 228512 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.6902 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.5381 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.307 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.781 ms/pair | 473328 B/op | 3656 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.9513 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.9386 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.594 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.492 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.481 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.498 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.395 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.134 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.484 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.691 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.452 ms/conn | 514184 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.74 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.479 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.529 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.732 ms/conn | 506024 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.396 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.234 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 3.175 ms/pair | 701848 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 3.512 ms/pair | 798040 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.847 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.131 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.225 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.113 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.117 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.296 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.203 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.248 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.143 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.35 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.143 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.125 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.222 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.142 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.654 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 2.84 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.296 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9962 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.028 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.017 ms/conn | 398136 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 0.9935 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9408 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.097 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.015 ms/conn | 535520 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.269 ms/conn | 642488 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.033 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.392 ms/conn | 549728 B/op | 4707 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9932 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.016 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.198 ms/conn | 540920 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8418 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.6319 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.667 ms/pair | 882248 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.141 ms/pair | 1064848 B/op | 8752 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.674 ms/pair | 881816 B/op | 7297 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.927 ms/conn | 654368 B/op | 6121 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.122 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.126 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.225 ms/conn | 175904 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.137 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.135 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.371 ms/conn | 219968 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.195 ms/conn | 223392 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.423 ms/conn | 331552 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.121 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.403 ms/conn | 266816 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.099 ms/conn | 214368 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.158 ms/conn | 247104 B/op | 1988 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.374 ms/conn | 269664 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.007 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.8443 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 2.041 ms/pair | 470272 B/op | 3877 allocs/op |
| mTLS session resumption handshake | 5 | 2.618 ms/pair | 590496 B/op | 4438 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.334 ms/conn | 270208 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.147 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.323 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.198 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.152 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.088 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.395 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.166 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.418 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.123 ms/conn | 518632 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 5.483 ms/conn | 525832 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.107 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.178 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.303 ms/conn | 510472 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.673 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 4.816 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 6.014 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 6.593 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.906 ms/pair | 711992 B/op | 3784 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.886 ms/conn | 578984 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.178 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.308 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.122 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.145 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.207 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 8.621 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.235 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.455 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.072 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.471 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.166 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.227 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.425 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.676 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 4.915 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.196 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 6.954 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.489 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.82 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.76 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 513.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.999 us/op | 2049.09 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.664 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 202.4 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 168.8 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 922.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.575 us/op | 1145.82 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 321.9 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.69 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 344 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.5 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 23.79 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.362 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.18 us/op | 550.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 847.2 ns/op | 1416.4 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.31 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.233 us/op | 537.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.303 us/op | 921.18 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.995 us/op | 171.56 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.152 us/op | 557.75 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.252 us/op | 532.89 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.986 us/op | 401.9 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 944.5 ns/op | 1270.5 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.018 us/op | 397.61 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 854.8 ns/op | 1403.8 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 897.3 ns/op | 1337.34 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.275 us/op | 941 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
