# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- Generated: `2026-10-07T05:04:55Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `ab4a14619a4eea21de84603f5e4f581b6296b3b5 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 551.426 us/op | 97993 B/op | 710 allocs/op |
| Full mTLS handshake | 5 | 855.542 us/op | 112214 B/op | 894 allocs/op |
| mTLS session resumption handshake | 5 | 427.607 us/op | 119411 B/op | 823 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.076 ms/op | 119968 B/op | 1067 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.448 ms/op | 138902 B/op | 1361 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 643.856 us/op | 117568 B/op | 934 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 883.894 us/op | 121529 B/op | 966 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 883.92 us/op | 121529 B/op | 966 allocs/op |
| Direct external PSK handshake | 5 | 376.59 us/op | 101438 B/op | 744 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.073 ms/op | 130098 B/op | 993 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.059 ms/op | 122315 B/op | 971 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.754 ms/op | 169846 B/op | 1428 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.749 ms/op | 157392 B/op | 1387 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 984.343 us/op | 147499 B/op | 1210 allocs/op |
| ECH handshake / via HRR | 5 | 1.002 ms/op | 150276 B/op | 1231 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 903.97 us/op | 145972 B/op | 743 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 870.754 us/op | 149316 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.271 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.839 ms/conn | 2487536 B/op | 21256 allocs/op |
| 1-RTT application-data round trip | 5 | 1.863 ms/conn | 2496208 B/op | 21576 allocs/op |
| Full mTLS handshake | 5 | 3.28 ms/conn | 3170272 B/op | 29390 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.581 ms/conn | 3668136 B/op | 31406 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.877 ms/conn | 2976560 B/op | 25185 allocs/op |
| Direct external PSK handshake | 5 | 0.4283 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.839 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.955 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.924 ms/conn | 3964168 B/op | 39510 allocs/op |
| Session resumption handshake | 5 | 0.5558 ms/conn | 4443344 B/op | 37896 allocs/op |
| mTLS session resumption handshake | 5 | 0.5617 ms/conn | 6527384 B/op | 49652 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5266 ms/conn | 287656 B/op | 1895 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.176 ms/conn | 3358952 B/op | 21578 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.157 ms/conn | 3400872 B/op | 21918 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.529 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.849 ms/conn | 1187344 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 5.059 ms/conn | 1207344 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.676 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.733 ms/conn | 1343344 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.935 ms/conn | 1425744 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9358 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.05 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.062 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.707 ms/conn | 1398992 B/op | 12384 allocs/op |
| Session resumption handshake | 5 | 0.9579 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 0.9522 ms/conn | 2302680 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.356 ms/conn | 1859000 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.454 ms/conn | 1871480 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.177 ms/conn | 1291360 B/op | 9299 allocs/op |
| 1-RTT application-data round trip | 5 | 4.77 ms/conn | 2444824 B/op | 10805 allocs/op |
| Full mTLS handshake | 5 | 5.544 ms/conn | 1729400 B/op | 14826 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.538 ms/conn | 1932512 B/op | 15349 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.166 ms/conn | 1356576 B/op | 10219 allocs/op |
| Direct external PSK handshake | 5 | 0.75 ms/conn | 1013840 B/op | 7499 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.793 ms/conn | 2459064 B/op | 11366 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.74 ms/conn | 2620024 B/op | 12445 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.77 ms/conn | 3354928 B/op | 22598 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 3571248 B/op | 20067 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1010 ms/pair | 241184 B/op | 994 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.364 ms/conn | 1592928 B/op | 9799 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.513 ms/conn | 1612768 B/op | 9939 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.525 ms/conn | 1771936 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.39 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.351 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 7.805 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 63.67 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.381 ms/conn | 63664 B/op | 56 allocs/op |
| Direct external PSK handshake | 5 | 57.67 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.449 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.456 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.46 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 1190848 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.409 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.991 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.137 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 45.47 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 63.04 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 911.8 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 3.2 us/op | 1279.91 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.437 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 235.9 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 182.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.334 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.825 us/op | 848.91 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 519.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 100.9 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 492.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 28.01 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 715.6 ns/op | 1676.85 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 732.7 ns/op | 1637.72 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 700.3 ns/op | 1713.48 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 33.244 us/op | 1971.35 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 658.6 ns/op | 1822.12 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 34.28 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.349 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 4.007 us/op | 299.46 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.429 us/op | 839.79 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.17 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 4.086 us/op | 293.66 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.84 us/op | 652.29 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.99 us/op | 120.12 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 4.059 us/op | 295.61 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 4.204 us/op | 285.44 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 5.273 us/op | 227.56 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.388 us/op | 864.82 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.782 us/op | 317.32 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.466 us/op | 818.42 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.527 us/op | 785.71 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 2.092 us/op | 573.72 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 3.552 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 8.095 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 2.015 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 4.259 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.916 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 4.18 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.155 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.786 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 4.062 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 7.734 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 13.048 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 9.424 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 20.486 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 894.2 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 2.28 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 18.58 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 903.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 902.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.955 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 5.089 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 18.58 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 2.006 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.98 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 3.082 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 3.453 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 5.99 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 2.453 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 4.711 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 8.183 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 4.334 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 7.968 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 425.7 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 798.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 130.9 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 78.08 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 335.6 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 262.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 359.2 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 671 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 59.4 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 643.8 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 109.1 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 58.2 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 90.27 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 910 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 119.2 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 85.07 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 67.64 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 808.9 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 543.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.68 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 78.96 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 122.7 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 73.86 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 26.69 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 393 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 204.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 76.88 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.467 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 1.017 us/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 1.013 us/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 18.51 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 82.37 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.085 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.626 us/op | 4264 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
