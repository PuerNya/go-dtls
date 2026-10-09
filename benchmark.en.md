# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `12d6cbda9e0eb434c12aaf19a17b774ff6042ec2`
- Generated: `2026-10-09T14:07:22Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `bb236b46bc87ef485a381e161c519d3bb9babff1 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 623.161 us/op | 101241 B/op | 741 allocs/op |
| Full mTLS handshake | 5 | 927.9 us/op | 116190 B/op | 938 allocs/op |
| mTLS session resumption handshake | 5 | 466.147 us/op | 123014 B/op | 864 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.174 ms/op | 123498 B/op | 1117 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.584 ms/op | 142233 B/op | 1413 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 714.116 us/op | 121281 B/op | 971 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 982.958 us/op | 125505 B/op | 1009 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 984.943 us/op | 125601 B/op | 1009 allocs/op |
| Direct external PSK handshake | 5 | 390.891 us/op | 103597 B/op | 767 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.118 ms/op | 135946 B/op | 1026 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.116 ms/op | 125965 B/op | 1002 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.852 ms/op | 176945 B/op | 1478 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.861 ms/op | 162175 B/op | 1431 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.07 ms/op | 150938 B/op | 1241 allocs/op |
| ECH handshake / via HRR | 5 | 1.073 ms/op | 153715 B/op | 1262 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 949.964 us/op | 149230 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 900.332 us/op | 152575 B/op | 805 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.357 ms/op | 178021 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.177 ms/conn | 2655248 B/op | 21796 allocs/op |
| 1-RTT application-data round trip | 5 | 2.182 ms/conn | 2666080 B/op | 22135 allocs/op |
| Full mTLS handshake | 5 | 3.917 ms/conn | 3357696 B/op | 30137 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.966 ms/conn | 3815376 B/op | 31885 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.183 ms/conn | 3156048 B/op | 25827 allocs/op |
| Direct external PSK handshake | 5 | 0.5804 ms/conn | 1767984 B/op | 15014 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.138 ms/conn | 2675216 B/op | 22914 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.311 ms/conn | 2835088 B/op | 23468 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.588 ms/conn | 4165512 B/op | 40431 allocs/op |
| Session resumption handshake | 5 | 0.6616 ms/conn | 4700912 B/op | 39510 allocs/op |
| mTLS session resumption handshake | 5 | 0.8547 ms/conn | 6830736 B/op | 51489 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6887 ms/conn | 298024 B/op | 1975 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.571 ms/conn | 3526712 B/op | 22119 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.554 ms/conn | 3568632 B/op | 22459 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 4.006 ms/conn | 4064328 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.128 ms/conn | 1325904 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 5.19 ms/conn | 1346544 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.745 ms/conn | 1487984 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.761 ms/conn | 1488032 B/op | 11923 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.098 ms/conn | 1564304 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 1.031 ms/conn | 827824 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.127 ms/conn | 1333104 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.297 ms/conn | 1582864 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.669 ms/conn | 1539600 B/op | 12666 allocs/op |
| Session resumption handshake | 5 | 1.114 ms/conn | 2354840 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.116 ms/conn | 2522376 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.392 ms/conn | 1998280 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.532 ms/conn | 2010680 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.514 ms/conn | 1316064 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.691 ms/conn | 2470264 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.605 ms/conn | 1760120 B/op | 15072 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.642 ms/conn | 1962240 B/op | 15592 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.51 ms/conn | 1381216 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.993 ms/conn | 1037328 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.541 ms/conn | 2484376 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.886 ms/conn | 2646744 B/op | 12685 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 14.24 ms/conn | 3391248 B/op | 22854 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3620272 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 243960 B/op | 1019 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.891 ms/conn | 1617888 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.006 ms/conn | 1637728 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.407 ms/conn | 1796928 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.968 ms/conn | 1169544 B/op | 1203 allocs/op |
| Full mTLS handshake | 5 | 8.973 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 9.135 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.311 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.233 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.972 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.914 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.28 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 1190872 B/op | 1203 allocs/op |
| mTLS session resumption handshake | 5 | 1017 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.361 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 7.115 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.86 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.33 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.77 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 486.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.971 us/op | 2078.4 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.524 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 206.3 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 176.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 878.5 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.626 us/op | 1129.7 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 306.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 91.24 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 356.9 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 20.77 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 569 ns/op | 2108.82 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 533 ns/op | 2251.47 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 527 ns/op | 2276.84 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 24.341 us/op | 2692.47 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 472.1 ns/op | 2541.75 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 24.07 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.581 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.11 us/op | 568.82 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 834.1 ns/op | 1438.65 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 14.47 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.285 us/op | 525.27 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.324 us/op | 906.65 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.407 us/op | 162.02 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.091 us/op | 574 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.167 us/op | 553.71 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.124 us/op | 384.08 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.016 us/op | 1180.84 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.275 us/op | 366.39 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 853.2 ns/op | 1406.48 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 897.3 ns/op | 1337.41 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.35 us/op | 888.7 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.808 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.577 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.535 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.472 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.526 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.499 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.235 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.412 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.264 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 16.04 ns/op | 0 B/op | 0 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 16.02 ns/op | 0 B/op | 0 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.254 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.458 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 710.7 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.848 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.122 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 747 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 753 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.637 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.333 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.121 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.696 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.685 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.313 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.571 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.567 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.835 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.86 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.015 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.77 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.799 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 265.1 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 509.6 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 109.6 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.62 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 297.4 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 263.5 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 368.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 562.1 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 51.89 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 539.8 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 88.49 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 46.04 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 73.52 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 746.6 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 90.41 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 72.2 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 68.28 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 571.5 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 410.6 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.36 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 50.58 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 94.84 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 65.16 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 33.14 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 288.8 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 155.5 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 80.23 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.09 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 798 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 784.9 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 23.33 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 54.87 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.193 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.931 us/op | 4296 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
