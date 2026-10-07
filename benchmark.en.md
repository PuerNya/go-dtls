# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- Generated: `2026-10-07T15:17:54Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `e09a10ecd5072c6138504b77fc710607cdcc6771 (Linux Release static)`

178 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (18)](#section-connection-lifecycle)
- [Real UDP interoperability (60)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (15)](#real-udp-wolfssl-client-wolfssl-server)
- [Record layer and reliability (34)](#section-record-layer-and-reliability)
- [Key schedule and cryptography (38)](#section-key-schedule-and-cryptography)
- [Wire encoding and parsing (26)](#section-wire-encoding-and-parsing)
- [Certificate compression (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 545.622 us/op | 97993 B/op | 710 allocs/op |
| Full mTLS handshake | 5 | 835.169 us/op | 112213 B/op | 894 allocs/op |
| mTLS session resumption handshake | 5 | 418.255 us/op | 119463 B/op | 823 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.06 ms/op | 120009 B/op | 1068 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.428 ms/op | 138880 B/op | 1361 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 639.841 us/op | 117569 B/op | 934 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 876.374 us/op | 121529 B/op | 966 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 879.722 us/op | 121529 B/op | 966 allocs/op |
| Direct external PSK handshake | 5 | 354.386 us/op | 101438 B/op | 744 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.04 ms/op | 130098 B/op | 993 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.031 ms/op | 122315 B/op | 971 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.72 ms/op | 169844 B/op | 1428 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.71 ms/op | 157394 B/op | 1387 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 951.206 us/op | 147499 B/op | 1210 allocs/op |
| ECH handshake / via HRR | 5 | 964.75 us/op | 150276 B/op | 1231 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 875.034 us/op | 145972 B/op | 743 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 840.046 us/op | 149316 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.239 ms/op | 174757 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.823 ms/conn | 2487488 B/op | 21255 allocs/op |
| 1-RTT application-data round trip | 5 | 1.837 ms/conn | 2496160 B/op | 21575 allocs/op |
| Full mTLS handshake | 5 | 3.26 ms/conn | 3170176 B/op | 29391 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.508 ms/conn | 3667392 B/op | 31404 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.859 ms/conn | 2976560 B/op | 25185 allocs/op |
| Direct external PSK handshake | 5 | 0.4213 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.823 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.939 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.895 ms/conn | 3965008 B/op | 39501 allocs/op |
| Session resumption handshake | 5 | 0.5073 ms/conn | 4443392 B/op | 37893 allocs/op |
| mTLS session resumption handshake | 5 | 0.5459 ms/conn | 6462392 B/op | 49655 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.4801 ms/conn | 287904 B/op | 1898 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.156 ms/conn | 3358952 B/op | 21578 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.144 ms/conn | 3400872 B/op | 21918 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.452 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.852 ms/conn | 1187344 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 4.871 ms/conn | 1207344 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.638 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.579 ms/conn | 1343344 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.746 ms/conn | 1425744 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9284 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.739 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.839 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.525 ms/conn | 1398928 B/op | 12383 allocs/op |
| Session resumption handshake | 5 | 0.9656 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 0.9488 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.114 ms/conn | 1859080 B/op | 10965 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.251 ms/conn | 1871480 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.192 ms/conn | 1291136 B/op | 9299 allocs/op |
| 1-RTT application-data round trip | 5 | 4.611 ms/conn | 2444504 B/op | 10805 allocs/op |
| Full mTLS handshake | 5 | 5.48 ms/conn | 1729328 B/op | 14829 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.518 ms/conn | 1931776 B/op | 15343 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.164 ms/conn | 1355936 B/op | 10219 allocs/op |
| Direct external PSK handshake | 5 | 0.772 ms/conn | 1014064 B/op | 7499 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.629 ms/conn | 2458680 B/op | 11366 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.748 ms/conn | 2620712 B/op | 12445 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.43 ms/conn | 3354936 B/op | 22599 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3572368 B/op | 20067 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1010 ms/pair | 241376 B/op | 994 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.392 ms/conn | 1592448 B/op | 9799 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.442 ms/conn | 1612768 B/op | 9939 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.486 ms/conn | 1771808 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.355 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.128 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 7.662 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 7.637 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.366 ms/conn | 63712 B/op | 57 allocs/op |
| Direct external PSK handshake | 5 | 0.983 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.014 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.187 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.2 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190232 B/op | 1203 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.402 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.512 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.86 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.68 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 56.52 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 708.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.658 us/op | 1540.9 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.211 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 229.8 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 179.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.191 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.027 us/op | 1017.12 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 416.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 99.18 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 410.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 27.97 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 645.2 ns/op | 1859.81 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 608.2 ns/op | 1972.91 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 615.5 ns/op | 1949.58 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 28.536 us/op | 2296.59 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 591.5 ns/op | 2028.76 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 30.65 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.364 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.709 us/op | 323.57 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.265 us/op | 948.66 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.34 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.889 us/op | 308.58 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.736 us/op | 691.4 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.428 us/op | 127.28 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.67 us/op | 327 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.785 us/op | 317.01 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 5.006 us/op | 239.7 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.356 us/op | 884.67 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.599 us/op | 333.38 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.331 us/op | 901.53 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.377 us/op | 871.3 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.956 us/op | 613.48 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 3.139 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 7.303 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.826 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 4.085 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.743 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.848 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.154 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.564 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.61 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 7.349 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 12.024 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 8.494 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 18.192 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 789.1 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 2.124 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 18.57 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 835.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 812.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.817 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.704 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 18.57 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.862 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.869 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.629 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 3 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 5.259 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 2.111 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 4.325 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.682 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 4.16 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 7.591 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 396.5 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 750.7 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 127.8 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 77.84 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 331.2 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 262.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 344.5 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 641.4 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 58.05 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 611 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 103.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 49.22 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 83.18 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 857.9 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 110 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 78.12 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 66.14 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 704.8 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 471.7 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.69 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 69.33 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 105.1 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 67 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 26.31 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 327.9 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 169.7 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.24 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.21 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 883.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 890.1 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 18.48 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 73.23 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.043 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.105 us/op | 4264 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
