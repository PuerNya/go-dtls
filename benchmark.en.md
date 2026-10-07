# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- Generated: `2026-10-07T21:20:11Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `0656a210be0269fa08a9a15874931fedd32eab0f (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 594.567 us/op | 97993 B/op | 710 allocs/op |
| Full mTLS handshake | 5 | 888.294 us/op | 112214 B/op | 894 allocs/op |
| mTLS session resumption handshake | 5 | 429.41 us/op | 119443 B/op | 823 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.157 ms/op | 119991 B/op | 1068 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.55 ms/op | 138947 B/op | 1362 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 690.613 us/op | 117569 B/op | 934 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 954.081 us/op | 121529 B/op | 966 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 951.755 us/op | 121529 B/op | 966 allocs/op |
| Direct external PSK handshake | 5 | 368.096 us/op | 101437 B/op | 744 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.076 ms/op | 130099 B/op | 993 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.077 ms/op | 122316 B/op | 971 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.816 ms/op | 169848 B/op | 1428 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.819 ms/op | 157395 B/op | 1387 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.02 ms/op | 147498 B/op | 1210 allocs/op |
| ECH handshake / via HRR | 5 | 1.022 ms/op | 150274 B/op | 1231 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 926.312 us/op | 145973 B/op | 743 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 869.604 us/op | 149317 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.324 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.093 ms/conn | 2487488 B/op | 21255 allocs/op |
| 1-RTT application-data round trip | 5 | 2.105 ms/conn | 2496208 B/op | 21576 allocs/op |
| Full mTLS handshake | 5 | 3.667 ms/conn | 3169928 B/op | 29390 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.928 ms/conn | 3666400 B/op | 31392 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.102 ms/conn | 2976560 B/op | 25185 allocs/op |
| Direct external PSK handshake | 5 | 0.5266 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.08 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.2 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.514 ms/conn | 3964168 B/op | 39514 allocs/op |
| Session resumption handshake | 5 | 0.6104 ms/conn | 4443032 B/op | 37899 allocs/op |
| mTLS session resumption handshake | 5 | 0.6562 ms/conn | 6462176 B/op | 49660 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6215 ms/conn | 287904 B/op | 1898 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.49 ms/conn | 3358952 B/op | 21578 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.46 ms/conn | 3400872 B/op | 21918 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.835 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.027 ms/conn | 1187344 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 5.125 ms/conn | 1207344 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.57 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.623 ms/conn | 1343344 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.996 ms/conn | 1425744 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9836 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.967 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.22 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.601 ms/conn | 1398928 B/op | 12383 allocs/op |
| Session resumption handshake | 5 | 1.024 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 1.04 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.278 ms/conn | 1859000 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.363 ms/conn | 1871480 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.477 ms/conn | 1291136 B/op | 9299 allocs/op |
| 1-RTT application-data round trip | 5 | 5.394 ms/conn | 2445144 B/op | 10805 allocs/op |
| Full mTLS handshake | 5 | 6.496 ms/conn | 1730392 B/op | 14829 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.508 ms/conn | 1932616 B/op | 15352 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.464 ms/conn | 1355680 B/op | 10219 allocs/op |
| Direct external PSK handshake | 5 | 0.958 ms/conn | 1014288 B/op | 7499 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.376 ms/conn | 2458728 B/op | 11366 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.431 ms/conn | 2620344 B/op | 12445 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.59 ms/conn | 3354712 B/op | 22597 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3571472 B/op | 20067 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 241256 B/op | 994 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.725 ms/conn | 1592768 B/op | 9799 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.836 ms/conn | 1612928 B/op | 9939 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.249 ms/conn | 1771552 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.93 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.523 ms/conn | 1169496 B/op | 1202 allocs/op |
| Full mTLS handshake | 5 | 9.326 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.869 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.263 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.309 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.527 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.528 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.42 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1017 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.014 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 7.389 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.6 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 38 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 48.27 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 596.8 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.216 us/op | 1847.96 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.839 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 214.8 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 184.3 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.008 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.848 us/op | 1064.4 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 351.7 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 95.37 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 379.6 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 20.79 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 581.8 ns/op | 2062.6 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 549.1 ns/op | 2185.36 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 551.5 ns/op | 2176 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 25.361 us/op | 2584.13 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 500.4 ns/op | 2397.93 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 26.4 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.579 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.161 us/op | 379.67 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.234 us/op | 972.12 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.15 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.75 us/op | 320.03 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.905 us/op | 630 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.598 us/op | 125.02 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.156 us/op | 380.27 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.316 us/op | 361.84 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.146 us/op | 289.45 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.399 us/op | 857.91 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.889 us/op | 308.57 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.467 us/op | 817.79 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.518 us/op | 790.66 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.945 us/op | 616.92 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.787 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.491 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.601 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.591 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.52 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.468 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.407 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.405 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.233 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.333 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.388 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.135 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.325 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 708.3 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.822 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.085 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 725.1 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 727.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.609 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.223 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.091 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.648 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.646 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.306 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.629 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.622 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.872 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.804 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.794 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.673 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.684 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 327.1 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 655.9 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 121.2 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.66 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 308.7 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 263.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 369.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 582.5 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 53.7 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 542.4 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 91.69 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 46.43 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 75.44 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 742.6 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 91.74 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 72.58 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 66.87 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 597.5 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 436.6 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.37 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 55.48 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 94.72 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 66.1 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 33.64 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 288.1 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 162 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 81.55 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.114 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 813.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 829.4 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 23.35 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 58.69 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.308 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 5.843 us/op | 4264 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
