# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `10300127dca519bf9c8f7672c1119d1baaa0893c`
- Generated: `2026-10-10T03:22:53Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `7499fc5b6c99eb39055a59f538bb3709971341c9 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 606.054 us/op | 101977 B/op | 741 allocs/op |
| Full mTLS handshake | 5 | 910.285 us/op | 116990 B/op | 938 allocs/op |
| mTLS session resumption handshake | 5 | 449.353 us/op | 123819 B/op | 864 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.144 ms/op | 124299 B/op | 1117 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.567 ms/op | 143077 B/op | 1413 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 707.838 us/op | 122021 B/op | 971 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 971.487 us/op | 126305 B/op | 1009 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 971.657 us/op | 126401 B/op | 1009 allocs/op |
| Direct external PSK handshake | 5 | 378.445 us/op | 104397 B/op | 767 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.084 ms/op | 136683 B/op | 1026 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.087 ms/op | 126701 B/op | 1002 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.809 ms/op | 177745 B/op | 1478 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.813 ms/op | 162975 B/op | 1431 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.043 ms/op | 152058 B/op | 1241 allocs/op |
| ECH handshake / via HRR | 5 | 1.043 ms/op | 154835 B/op | 1262 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 912.322 us/op | 149966 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 861.686 us/op | 153312 B/op | 805 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.306 ms/op | 178759 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.146 ms/conn | 2672528 B/op | 21796 allocs/op |
| 1-RTT application-data round trip | 5 | 2.175 ms/conn | 2683408 B/op | 22136 allocs/op |
| Full mTLS handshake | 5 | 3.881 ms/conn | 3375512 B/op | 30121 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.876 ms/conn | 3834672 B/op | 31876 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.282 ms/conn | 3239088 B/op | 25831 allocs/op |
| Direct external PSK handshake | 5 | 0.5498 ms/conn | 1786544 B/op | 15014 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.143 ms/conn | 2692496 B/op | 22914 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.266 ms/conn | 2852368 B/op | 23468 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.535 ms/conn | 4182552 B/op | 40426 allocs/op |
| Session resumption handshake | 5 | 0.6609 ms/conn | 4738032 B/op | 39510 allocs/op |
| mTLS session resumption handshake | 5 | 0.8238 ms/conn | 6871104 B/op | 51505 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6858 ms/conn | 299752 B/op | 1975 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.57 ms/conn | 3543992 B/op | 22119 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.524 ms/conn | 3585912 B/op | 22459 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.923 ms/conn | 4081608 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.056 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 5.076 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.665 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.67 ms/conn | 1495024 B/op | 11922 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.047 ms/conn | 1571344 B/op | 12602 allocs/op |
| Direct external PSK handshake | 5 | 1.001 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.043 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.219 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.574 ms/conn | 1546448 B/op | 12663 allocs/op |
| Session resumption handshake | 5 | 1.094 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.11 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.242 ms/conn | 2005320 B/op | 11245 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.447 ms/conn | 2017720 B/op | 11364 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.498 ms/conn | 1326800 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.395 ms/conn | 2479672 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.543 ms/conn | 1771928 B/op | 15078 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.57 ms/conn | 1974256 B/op | 15589 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.507 ms/conn | 1391232 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.915 ms/conn | 1048848 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.401 ms/conn | 2495000 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.465 ms/conn | 2656024 B/op | 12685 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.71 ms/conn | 3403088 B/op | 22857 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3641872 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244600 B/op | 1019 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.8 ms/conn | 1627968 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.843 ms/conn | 1648608 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.415 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.273 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.669 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 8.881 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 9.326 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.938 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.281 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.471 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.347 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.36 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1017 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.109 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.7 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.61 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 35.41 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 44.18 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 477.5 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.945 us/op | 2105.38 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.527 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 200 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 174.4 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 866.6 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.483 us/op | 1175.85 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 297.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 87.78 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 343.2 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.58 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 497.1 ns/op | 2413.86 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 462.4 ns/op | 2595.18 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 472.2 ns/op | 2541.04 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 22.456 us/op | 2918.39 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 408.3 ns/op | 2939.37 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 23.59 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.928 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.051 us/op | 585.06 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 825.7 ns/op | 1453.31 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 14.44 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.176 us/op | 551.55 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.276 us/op | 940.38 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.208 us/op | 166.48 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.022 us/op | 593.42 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.123 us/op | 565.12 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.034 us/op | 395.54 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 917.4 ns/op | 1308.06 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.197 us/op | 375.3 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 821.7 ns/op | 1460.37 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 869.9 ns/op | 1379.51 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.316 us/op | 911.9 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.988 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.326 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.501 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.379 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.472 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.386 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.409 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.394 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.218 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 16.04 ns/op | 0 B/op | 0 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 16.05 ns/op | 0 B/op | 0 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.118 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.016 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 692 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.81 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.127 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 743.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 722.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.585 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.095 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.162 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.617 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.611 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.148 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.456 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.389 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.782 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.629 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.578 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.465 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.441 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 254.3 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 508 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 107 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.57 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 291.3 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 264.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 354.4 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 539.6 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 49.79 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 521.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 85.38 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 44.35 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.95 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 719.5 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 86.01 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 70 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 67.27 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 619.6 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 402.8 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 49.84 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 101.2 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 62.55 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 33.5 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 268.9 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 152 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 80.48 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.048 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 769.8 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 776.3 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 23.31 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 53.43 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.166 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.318 us/op | 4296 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
