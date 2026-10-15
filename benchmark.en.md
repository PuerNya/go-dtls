# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `c7cb9ac2cd73dcebe3ab9e32545048ee026e9bd8`
- Generated: `2026-10-15T02:31:12Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `360f93b322addb95af2dadbb19a7eba3a06725bc (Linux Release static)`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 602.046 us/op | - | 101322 B/op | 732 allocs/op |
| Full mTLS handshake | 5 | 909.355 us/op | - | 116334 B/op | 929 allocs/op |
| mTLS session resumption handshake | 5 | 460.019 us/op | - | 123190 B/op | 855 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.14 ms/op | - | 123643 B/op | 1108 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.527 ms/op | - | 142436 B/op | 1404 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 703.826 us/op | - | 121344 B/op | 962 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 958.36 us/op | - | 125650 B/op | 1000 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 959.425 us/op | - | 125745 B/op | 1000 allocs/op |
| Direct external PSK handshake | 5 | 392.854 us/op | - | 103774 B/op | 758 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.031 ms/op | - | 136091 B/op | 1017 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.056 ms/op | - | 126045 B/op | 993 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.754 ms/op | - | 177153 B/op | 1469 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.753 ms/op | - | 162384 B/op | 1422 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.031 ms/op | - | 151403 B/op | 1232 allocs/op |
| ECH handshake / via HRR | 5 | 1.033 ms/op | - | 154179 B/op | 1253 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 892.236 us/op | - | 149311 B/op | 765 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 840.72 us/op | - | 152655 B/op | 796 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.161 ms/op | - | 178103 B/op | 814 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 64 bytes | 5 | 22.35 us/op | 5.73 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / No CID / 1200 bytes | 5 | 25.71 us/op | 93.35 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 64 bytes | 5 | 22.494 us/op | 5.69 MB/s | 592 B/op | 12 allocs/op |
| Listener UDP round trip / 1 connection / CID / 1200 bytes | 5 | 26.046 us/op | 92.14 MB/s | 2992 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 64 bytes | 5 | 22.669 us/op | 5.65 MB/s | 597 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / No CID / 1200 bytes | 5 | 25.589 us/op | 93.79 MB/s | 2996 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 64 bytes | 5 | 23.582 us/op | 5.43 MB/s | 600 B/op | 12 allocs/op |
| Listener UDP round trip / 64 connections / CID / 1200 bytes | 5 | 26.825 us/op | 89.47 MB/s | 2999 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 64 bytes | 5 | 28.046 us/op | 4.56 MB/s | 624 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / No CID / 1200 bytes | 5 | 31.351 us/op | 76.55 MB/s | 3033 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 64 bytes | 5 | 37.207 us/op | 3.44 MB/s | 619 B/op | 12 allocs/op |
| Listener UDP round trip / 1024 connections / CID / 1200 bytes | 5 | 41.28 us/op | 58.14 MB/s | 2996 B/op | 12 allocs/op |
| 0-RTT / 64 initial replay entries | 5 | 430.442 us/op | - | 109276 B/op | 796 allocs/op |
| 0-RTT / 4096 initial replay entries | 5 | 434.216 us/op | - | 109115 B/op | 796 allocs/op |
| In-memory datagram round trip / 64 bytes | 5 | 3.404 us/op | 37.6 MB/s | 384 B/op | 4 allocs/op |
| In-memory datagram round trip / 1200 bytes | 5 | 6.734 us/op | 356.39 MB/s | 5184 B/op | 4 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.004 ms/conn | 2671752 B/op | 21759 allocs/op |
| 1-RTT application-data round trip | 5 | 2.004 ms/conn | 2682584 B/op | 22098 allocs/op |
| Full mTLS handshake | 5 | 3.614 ms/conn | 3374688 B/op | 30093 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.607 ms/conn | 3834872 B/op | 31836 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.036 ms/conn | 3149576 B/op | 25608 allocs/op |
| Direct external PSK handshake | 5 | 0.5736 ms/conn | 1785080 B/op | 14976 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.994 ms/conn | 2689112 B/op | 22876 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.143 ms/conn | 2850176 B/op | 23392 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.254 ms/conn | 4181696 B/op | 40362 allocs/op |
| Session resumption handshake | 5 | 0.6796 ms/conn | 4735608 B/op | 39436 allocs/op |
| mTLS session resumption handshake | 5 | 0.8264 ms/conn | 6867144 B/op | 51418 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6773 ms/conn | 299120 B/op | 1962 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.366 ms/conn | 3543168 B/op | 22081 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.353 ms/conn | 3585088 B/op | 22421 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.592 ms/conn | 4080768 B/op | 22781 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.831 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.976 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.456 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.485 ms/conn | 1495824 B/op | 11930 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.89 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 1.055 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.84 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.079 ms/conn | 1590040 B/op | 12664 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.307 ms/conn | 1546512 B/op | 12664 allocs/op |
| Session resumption handshake | 5 | 1.204 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.236 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.118 ms/conn | 2005320 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.144 ms/conn | 2018504 B/op | 11372 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.493 ms/conn | 1326496 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.277 ms/conn | 2480552 B/op | 11026 allocs/op |
| Full mTLS handshake | 5 | 6.219 ms/conn | 1771560 B/op | 15069 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.197 ms/conn | 1973880 B/op | 15586 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.461 ms/conn | 1391680 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.997 ms/conn | 1048848 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.419 ms/conn | 2494936 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.583 ms/conn | 2655384 B/op | 12665 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.9 ms/conn | 3403576 B/op | 22837 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3642416 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244240 B/op | 1017 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.688 ms/conn | 1627968 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.723 ms/conn | 1648128 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.803 ms/conn | 1807040 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.997 ms/conn | 1169792 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 1024 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1024 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1024 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.171 ms/conn | 1183488 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.203 ms/conn | 1170112 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.8 ms/conn | 1172800 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1191240 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190600 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7499fc5b6c99eb39055a59f538bb3709971341c9 | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1024 ms/conn | 64464 B/op | 57 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1024 ms/conn | 64392 B/op | 55 allocs/op |

<a id="real-udp-nss-go-go"></a>
### NSS / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.053 ms/conn | 460680 B/op | 3537 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.097 ms/conn | 526744 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.024 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.163 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.102 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.003 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.262 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.039 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.457 ms/conn | 549400 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.001 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.047 ms/conn | 495480 B/op | 3765 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.198 ms/conn | 540920 B/op | 4304 allocs/op |
| Direct external PSK handshake | 5 | 0.6209 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.681 ms/pair | 882968 B/op | 7357 allocs/op |
| mTLS session resumption handshake | 5 | 2.176 ms/pair | 1065720 B/op | 8761 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.679 ms/pair | 882312 B/op | 7302 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.372 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.188 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-nss-go-peer"></a>
### NSS / go-dtls client -> NSS server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.454 ms/conn | 181552 B/op | 1534 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.555 ms/conn | 214960 B/op | 1534 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.467 ms/conn | 146832 B/op | 1514 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.965 ms/conn | 184560 B/op | 1558 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.371 ms/conn | 186640 B/op | 1594 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.052 ms/conn | 215344 B/op | 1810 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.559 ms/conn | 274736 B/op | 1646 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation with an NSS server; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 2.29 ms/conn | 222672 B/op | 1910 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.334 ms/conn | 181840 B/op | 1546 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.378 ms/conn | 213200 B/op | 1718 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.382 ms/conn | 194288 B/op | 1638 allocs/op |
| Direct external PSK handshake | 5 | 0.9056 ms/conn | 158896 B/op | 1306 allocs/op |
| Session resumption handshake | 5 | 2.495 ms/pair | 379024 B/op | 3141 allocs/op |
| mTLS session resumption handshake | 5 | 3.458 ms/pair | 458800 B/op | 3517 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 1.607 ms/conn | 203312 B/op | 1714 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.8 ms/pair | 420336 B/op | 3477 allocs/op |

<a id="real-udp-nss-peer-go"></a>
### NSS / NSS client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.416 ms/conn | 485640 B/op | 1991 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.56 ms/conn | 517256 B/op | 1993 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.626 ms/conn | 457576 B/op | 1975 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.065 ms/conn | 488472 B/op | 2016 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.531 ms/conn | 490840 B/op | 2050 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.142 ms/conn | 518600 B/op | 2191 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 6.596 ms/conn | 555912 B/op | 2119 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 4.409 ms/conn | 518888 B/op | 2175 allocs/op |
| Full mTLS handshake | 5 | 5.06 ms/conn | 528296 B/op | 2567 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.321 ms/conn | 487112 B/op | 2027 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.552 ms/conn | 520616 B/op | 2321 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.337 ms/conn | 510376 B/op | 2205 allocs/op |
| Direct external PSK handshake | 5 | 3.567 ms/conn | 464488 B/op | 1891 allocs/op |
| Session resumption handshake | 5 | 5.383 ms/pair | 712392 B/op | 3932 allocs/op |
| mTLS session resumption handshake | 5 | 6.144 ms/pair | 811432 B/op | 4730 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 5.514 ms/pair | 714456 B/op | 3881 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 5.346 ms/conn | 521608 B/op | 2351 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 6.284 ms/pair | 749256 B/op | 4300 allocs/op |

<a id="real-udp-nss-peer-peer"></a>
### NSS / NSS client -> NSS server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 4.629 ms/conn | 212704 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 4.706 ms/conn | 212192 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 4.586 ms/conn | 212448 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.842 ms/conn | 212192 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.434 ms/conn | 212192 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 6.037 ms/conn | 212448 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 7.617 ms/conn | 212448 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | - | Skipped: NSS fixture cannot force handshake fragmentation for two NSS peers; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Full mTLS handshake | 5 | 5.639 ms/conn | 212192 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 4.857 ms/conn | 212192 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.594 ms/conn | 212192 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 4.537 ms/conn | 212448 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 3.676 ms/conn | 212448 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.742 ms/pair | 207264 B/op | 185 allocs/op |
| mTLS session resumption handshake | 5 | 7.704 ms/pair | 207264 B/op | 184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: NSS server cannot accept DTLS 1.3 early data in this fixture; last verified against NSS commit b9283c1e4ee232b4450d36115272af670d38ec81 | - | - |
| Delegated-credential handshake + datagram exchange | 5 | 5.267 ms/conn | 212448 B/op | 176 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 7.348 ms/pair | 207520 B/op | 184 allocs/op |

<a id="real-udp-boringssl-go-go"></a>
### BoringSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.003 ms/conn | 460504 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.048 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.002 ms/conn | 398136 B/op | 3500 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9549 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.093 ms/conn | 472344 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.073 ms/conn | 535520 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.301 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.088 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.412 ms/conn | 549224 B/op | 4704 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.999 ms/conn | 462632 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.043 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.197 ms/conn | 540920 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8162 ms/conn | 445432 B/op | 3472 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.6186 ms/conn | 425944 B/op | 3260 allocs/op |
| Session resumption handshake | 5 | 1.694 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.15 ms/pair | 1065728 B/op | 8755 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 1.423 ms/conn | 658296 B/op | 5540 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 1.357 ms/conn | 519608 B/op | 4088 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 2.219 ms/pair | 961976 B/op | 8087 allocs/op |

<a id="real-udp-boringssl-go-peer"></a>
### BoringSSL / go-dtls client -> BoringSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.7682 ms/conn | 179872 B/op | 1524 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 0.9243 ms/conn | 213152 B/op | 1524 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 0.7942 ms/conn | 145056 B/op | 1503 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.7589 ms/conn | 183008 B/op | 1548 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.488 ms/conn | 185088 B/op | 1584 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 0.9559 ms/conn | 203552 B/op | 1756 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.074 ms/conn | 270496 B/op | 1592 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 0.9072 ms/conn | 180320 B/op | 1556 allocs/op |
| Full mTLS handshake | 5 | 1.033 ms/conn | 220000 B/op | 1900 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.8789 ms/conn | 180192 B/op | 1536 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 0.8007 ms/conn | 211776 B/op | 1708 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 0.9791 ms/conn | 228384 B/op | 2068 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.7677 ms/conn | 164064 B/op | 1364 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 0.582 ms/conn | 160544 B/op | 1308 allocs/op |
| Session resumption handshake | 5 | 1.345 ms/pair | 388000 B/op | 3279 allocs/op |
| mTLS session resumption handshake | 5 | 1.703 ms/pair | 473568 B/op | 3655 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 0.9792 ms/conn | 240960 B/op | 2252 allocs/op |
| Delegated-credential handshake + datagram exchange | 5 | 0.973 ms/conn | 200928 B/op | 1704 allocs/op |
| Delegated-credential handshake + session resumption | 5 | 1.775 ms/pair | 431616 B/op | 3631 allocs/op |

<a id="real-udp-boringssl-peer-go"></a>
### BoringSSL / BoringSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.552 ms/conn | 480968 B/op | 1913 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.574 ms/conn | 512776 B/op | 1913 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.622 ms/conn | 453256 B/op | 1897 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.483 ms/conn | 483976 B/op | 1937 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.311 ms/conn | 486312 B/op | 1973 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.624 ms/conn | 514408 B/op | 2109 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.773 ms/conn | 546344 B/op | 1933 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.556 ms/conn | 514248 B/op | 2097 allocs/op |
| Full mTLS handshake | 5 | 2.798 ms/conn | 520936 B/op | 2393 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.562 ms/conn | 482504 B/op | 1945 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.585 ms/conn | 496648 B/op | 2025 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.632 ms/conn | 506024 B/op | 2133 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.47 ms/conn | 484072 B/op | 2009 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.303 ms/conn | 461416 B/op | 1813 allocs/op |
| Session resumption handshake | 5 | 3.406 ms/pair | 701832 B/op | 3762 allocs/op |
| mTLS session resumption handshake | 5 | 3.687 ms/pair | 798168 B/op | 4474 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.988 ms/conn | 562664 B/op | 2793 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-boringssl-peer-peer"></a>
### BoringSSL / BoringSSL client -> BoringSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 2.261 ms/conn | 212768 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 2.243 ms/conn | 212512 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 2.165 ms/conn | 212832 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 2.156 ms/conn | 212256 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 4.274 ms/conn | 212256 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 2.293 ms/conn | 212768 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 2.33 ms/conn | 212768 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 2.159 ms/conn | 212768 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 2.334 ms/conn | 212512 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 2.198 ms/conn | 212256 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.159 ms/conn | 212512 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 2.277 ms/conn | 212512 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 2.236 ms/conn | 212512 B/op | 176 allocs/op |
| Imported external PSK handshake + datagram exchange | 5 | 2.056 ms/conn | 212512 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 2.747 ms/pair | 210592 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 3.016 ms/pair | 210848 B/op | 192 allocs/op |
| ECH-accepted handshake + datagram exchange | 5 | 2.365 ms/conn | 212256 B/op | 176 allocs/op |
| Delegated-credential handshake + datagram exchange | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |
| Delegated-credential handshake + session resumption | - | Skipped: BoringSSL client does not authenticate delegated credentials in DTLS 1.3; last verified against BoringSSL commit 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978 | - | - |

<a id="real-udp-openssl-go-go"></a>
### OpenSSL / go-dtls client -> go-dtls server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 0.9938 ms/conn | 460632 B/op | 3536 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.123 ms/conn | 526872 B/op | 3536 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.008 ms/conn | 398008 B/op | 3500 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.099 ms/conn | 525288 B/op | 4222 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 0.9414 ms/conn | 467416 B/op | 3584 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 2.089 ms/conn | 472216 B/op | 3656 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.055 ms/conn | 535648 B/op | 3995 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.266 ms/conn | 642360 B/op | 3716 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.124 ms/conn | 517248 B/op | 3939 allocs/op |
| Full mTLS handshake | 5 | 1.48 ms/conn | 549600 B/op | 4707 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 0.9982 ms/conn | 462504 B/op | 3584 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.069 ms/conn | 493160 B/op | 3756 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.249 ms/conn | 541944 B/op | 4304 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 0.8146 ms/conn | 445432 B/op | 3472 allocs/op |
| Direct external PSK handshake | 5 | 0.6301 ms/conn | 423192 B/op | 3252 allocs/op |
| Session resumption handshake | 5 | 1.72 ms/pair | 882376 B/op | 7355 allocs/op |
| mTLS session resumption handshake | 5 | 2.186 ms/pair | 1065224 B/op | 8755 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1.704 ms/pair | 881936 B/op | 7300 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 1.955 ms/conn | 653608 B/op | 6118 allocs/op |

<a id="real-udp-openssl-go-peer"></a>
### OpenSSL / go-dtls client -> OpenSSL server

Median time is measured by the go-dtls client and includes the authenticated datagram exchange; `ms/pair` includes the full connection that obtains the ticket. Runs with an OpenSSL server deliberately drop the client's epoch-2 ACK.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 1.095 ms/conn | 210848 B/op | 1776 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 1.221 ms/conn | 244896 B/op | 1776 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 1.187 ms/conn | 175776 B/op | 1755 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.209 ms/conn | 264480 B/op | 2316 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 1.157 ms/conn | 220192 B/op | 1844 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 3.581 ms/conn | 219856 B/op | 1860 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 1.213 ms/conn | 223520 B/op | 1892 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 1.502 ms/conn | 331424 B/op | 1880 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 1.141 ms/conn | 211552 B/op | 1788 allocs/op |
| Full mTLS handshake | 5 | 1.616 ms/conn | 266688 B/op | 2240 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 1.204 ms/conn | 214240 B/op | 1812 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.106 ms/conn | 247040 B/op | 1986 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 1.515 ms/conn | 270240 B/op | 2416 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 1.042 ms/conn | 193952 B/op | 1600 allocs/op |
| Direct external PSK handshake | 5 | 0.8293 ms/conn | 199008 B/op | 1680 allocs/op |
| Session resumption handshake | 5 | 2.106 ms/pair | 470592 B/op | 3879 allocs/op |
| mTLS session resumption handshake | 5 | 2.792 ms/pair | 590720 B/op | 4439 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 1.295 ms/conn | 270080 B/op | 2304 allocs/op |

<a id="real-udp-openssl-peer-go"></a>
### OpenSSL / OpenSSL client -> go-dtls server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.499 ms/conn | 485416 B/op | 1937 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.518 ms/conn | 517160 B/op | 1937 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.522 ms/conn | 457704 B/op | 1921 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.388 ms/conn | 500392 B/op | 2141 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.392 ms/conn | 488232 B/op | 1961 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 7.99 ms/conn | 490568 B/op | 1997 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.438 ms/conn | 518600 B/op | 2133 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.789 ms/conn | 549960 B/op | 1957 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.399 ms/conn | 518696 B/op | 2121 allocs/op |
| Full mTLS handshake | 5 | 5.838 ms/conn | 525832 B/op | 2433 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.436 ms/conn | 486824 B/op | 1969 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.742 ms/conn | 518344 B/op | 2189 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.945 ms/conn | 510984 B/op | 2157 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 6.26 ms/conn | 485480 B/op | 2001 allocs/op |
| Direct external PSK handshake | 5 | 5.114 ms/conn | 464360 B/op | 1825 allocs/op |
| Session resumption handshake | 5 | 6.251 ms/pair | 713448 B/op | 3882 allocs/op |
| mTLS session resumption handshake | 5 | 6.791 ms/pair | 809960 B/op | 4610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 6.354 ms/pair | 711992 B/op | 3785 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.342 ms/conn | 578728 B/op | 3137 allocs/op |

<a id="real-udp-openssl-peer-peer"></a>
### OpenSSL / OpenSSL client -> OpenSSL server

Median time includes the peer client process startup and authenticated datagram exchange; `ms/pair` includes both full and resumed connections. These process metrics are not directly comparable with client-only timings.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated AES-128-GCM connection + datagram exchange | 5 | 5.465 ms/conn | 212832 B/op | 176 allocs/op |
| AES-256-GCM/SHA-384 connection + datagram exchange | 5 | 5.567 ms/conn | 212576 B/op | 176 allocs/op |
| ChaCha20-Poly1305 connection + datagram exchange | 5 | 5.356 ms/conn | 212896 B/op | 176 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.501 ms/conn | 212576 B/op | 176 allocs/op |
| P-256 key exchange + datagram exchange | 5 | 5.534 ms/conn | 212320 B/op | 176 allocs/op |
| P-384 key exchange + datagram exchange | 5 | 9.114 ms/conn | 212320 B/op | 176 allocs/op |
| HelloRetryRequest handshake + datagram exchange | 5 | 5.553 ms/conn | 212832 B/op | 176 allocs/op |
| X25519MLKEM768 hybrid handshake + datagram exchange | 5 | 5.759 ms/conn | 212576 B/op | 176 allocs/op |
| MTU-256 fragmented handshake + datagram exchange | 5 | 5.398 ms/conn | 212832 B/op | 176 allocs/op |
| Full mTLS handshake | 5 | 5.714 ms/conn | 212576 B/op | 176 allocs/op |
| ALPN-negotiated connection + datagram exchange | 5 | 5.523 ms/conn | 212320 B/op | 176 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.614 ms/conn | 212576 B/op | 176 allocs/op |
| OCSP-stapled handshake + datagram exchange | 5 | 5.719 ms/conn | 212576 B/op | 176 allocs/op |
| Raw-public-key handshake + datagram exchange | 5 | 5.981 ms/conn | 212576 B/op | 176 allocs/op |
| Direct external PSK handshake | 5 | 5.338 ms/conn | 212576 B/op | 176 allocs/op |
| Session resumption handshake | 5 | 6.609 ms/pair | 210656 B/op | 192 allocs/op |
| mTLS session resumption handshake | 5 | 7.085 ms/pair | 210912 B/op | 192 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: OpenSSL listener cannot enter the DTLS early-data read state; last verified against OpenSSL commit 41eecc7c7cd4069b1b29bca49b9d97825c6bf648 | - | - |
| PHA + 1-RTT application-data round trip | 5 | 5.784 ms/conn | 213088 B/op | 176 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.89 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 45.54 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 525.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.006 us/op | 2041.98 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.696 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 202.7 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 169.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 935.5 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.561 us/op | 1150.22 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 321.4 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 93.25 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 343.4 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.52 ns/op | - | 0 B/op | 0 allocs/op |
| Parse ACK / Owned | 5 | 23.98 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.362 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.197 us/op | 546.2 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 853.4 ns/op | 1406.09 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.3 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.282 us/op | 525.82 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.311 us/op | 915.37 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.93 us/op | 173.17 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.166 us/op | 554.02 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.262 us/op | 530.4 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.963 us/op | 404.95 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.017 us/op | 1179.75 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.021 us/op | 397.23 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 843.7 ns/op | 1422.23 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 898.3 ns/op | 1335.84 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.274 us/op | 942.05 MB/s | 1280 B/op | 1 allocs/op |

[Raw Go benchmark output](benchmark.txt)
