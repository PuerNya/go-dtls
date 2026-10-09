# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `10300127dca519bf9c8f7672c1119d1baaa0893c`
- Generated: `2026-10-09T19:11:06Z`
- Go: `go version go1.27.2 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `23f245c02c4319bb80baf7f3ae0c1de40e3fb68e (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 585.754 us/op | 101978 B/op | 741 allocs/op |
| Full mTLS handshake | 5 | 880.564 us/op | 116990 B/op | 938 allocs/op |
| mTLS session resumption handshake | 5 | 446.123 us/op | 123815 B/op | 864 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.104 ms/op | 124300 B/op | 1117 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.507 ms/op | 143085 B/op | 1413 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 683.684 us/op | 122018 B/op | 971 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 931.296 us/op | 126305 B/op | 1009 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 938.149 us/op | 126401 B/op | 1009 allocs/op |
| Direct external PSK handshake | 5 | 375.985 us/op | 104398 B/op | 767 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.045 ms/op | 136683 B/op | 1026 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.041 ms/op | 126701 B/op | 1002 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.723 ms/op | 177745 B/op | 1478 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.737 ms/op | 162975 B/op | 1431 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 996.297 us/op | 152059 B/op | 1241 allocs/op |
| ECH handshake / via HRR | 5 | 996.978 us/op | 154835 B/op | 1262 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 881.193 us/op | 149966 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 837.207 us/op | 153311 B/op | 805 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.134 ms/op | 178759 B/op | 823 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.003 ms/conn | 2672528 B/op | 21796 allocs/op |
| 1-RTT application-data round trip | 5 | 2.02 ms/conn | 2683360 B/op | 22135 allocs/op |
| Full mTLS handshake | 5 | 3.551 ms/conn | 3375560 B/op | 30130 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.592 ms/conn | 3834128 B/op | 31877 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.127 ms/conn | 3239136 B/op | 25832 allocs/op |
| Direct external PSK handshake | 5 | 0.5599 ms/conn | 1786544 B/op | 15014 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.972 ms/conn | 2692544 B/op | 22915 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.126 ms/conn | 2852368 B/op | 23468 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.129 ms/conn | 4182056 B/op | 40442 allocs/op |
| Session resumption handshake | 5 | 0.6747 ms/conn | 4738512 B/op | 39512 allocs/op |
| mTLS session resumption handshake | 5 | 0.8124 ms/conn | 6870664 B/op | 51513 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6873 ms/conn | 299504 B/op | 1972 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.359 ms/conn | 3543992 B/op | 22119 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.324 ms/conn | 3585912 B/op | 22459 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.627 ms/conn | 4081608 B/op | 22819 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.755 ms/conn | 1332944 B/op | 11042 allocs/op |
| 1-RTT application-data round trip | 5 | 4.792 ms/conn | 1353584 B/op | 11302 allocs/op |
| Full mTLS handshake | 5 | 6.25 ms/conn | 1495024 B/op | 11922 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.278 ms/conn | 1495872 B/op | 11931 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.704 ms/conn | 1573536 B/op | 12619 allocs/op |
| Direct external PSK handshake | 5 | 1.004 ms/conn | 834864 B/op | 7042 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.715 ms/conn | 1340144 B/op | 11602 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.901 ms/conn | 1589904 B/op | 12662 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.116 ms/conn | 1546448 B/op | 12663 allocs/op |
| Session resumption handshake | 5 | 1.152 ms/conn | 2370200 B/op | 19324 allocs/op |
| mTLS session resumption handshake | 5 | 1.148 ms/conn | 2537736 B/op | 20244 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.992 ms/conn | 2006104 B/op | 11253 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.029 ms/conn | 2017800 B/op | 11365 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.365 ms/conn | 1326784 B/op | 9519 allocs/op |
| 1-RTT application-data round trip | 5 | 5.037 ms/conn | 2480184 B/op | 11025 allocs/op |
| Full mTLS handshake | 5 | 6.032 ms/conn | 1770712 B/op | 15063 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.065 ms/conn | 1973920 B/op | 15589 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.358 ms/conn | 1391904 B/op | 10439 allocs/op |
| Direct external PSK handshake | 5 | 0.979 ms/conn | 1048624 B/op | 7679 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.04 ms/conn | 2495000 B/op | 11586 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.064 ms/conn | 2656344 B/op | 12685 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.74 ms/conn | 3405232 B/op | 22866 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3642288 B/op | 20507 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 244984 B/op | 1019 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.538 ms/conn | 1628448 B/op | 10019 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.646 ms/conn | 1648288 B/op | 10159 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.691 ms/conn | 1807296 B/op | 10339 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.893 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.951 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 8.575 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.469 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.885 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.074 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.139 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.017 ms/conn | 1169408 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.58 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190872 B/op | 1203 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.92 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.498 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.832 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 36.05 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 45.18 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 511.8 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.992 us/op | 2055.78 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.623 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 202.7 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 167.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 904.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.488 us/op | 1174.45 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 318.1 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 89.11 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 334.8 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.18 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 513.9 ns/op | 2334.89 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 456.9 ns/op | 2626.26 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 472.3 ns/op | 2540.94 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 23.667 us/op | 2769.03 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 408 ns/op | 2941.3 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 24.04 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.672 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.131 us/op | 563.04 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 836.8 ns/op | 1433.96 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.28 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.256 us/op | 531.92 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.268 us/op | 946.39 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.945 us/op | 172.78 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.09 us/op | 574.05 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.192 us/op | 547.4 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.94 us/op | 408.12 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 930.6 ns/op | 1289.45 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.048 us/op | 393.68 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 831.2 ns/op | 1443.62 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 889 ns/op | 1349.86 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.25 us/op | 960.02 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.721 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.118 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.487 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.303 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.434 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.237 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.354 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.067 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 15.16 ns/op | 0 B/op | 0 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 15.21 ns/op | 0 B/op | 0 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 6.907 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.379 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 687.9 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.802 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.957 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 711.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 755.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.534 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.998 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.044 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.575 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.57 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.228 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.51 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.327 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.78 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.669 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.388 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.533 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.211 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 262.5 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 525.3 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 103 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 78.97 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 266.4 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.5 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 340.3 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 554.6 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 49.94 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 572.1 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 86.28 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 43.33 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 71.62 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 792.8 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 87.46 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.54 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 64.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 601.1 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 402.9 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 53.26 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 90.29 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 59.85 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 29.97 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 278 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 151 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 78.53 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.1 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 797.1 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 806.7 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 23.11 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 55.74 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.909 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.439 us/op | 4296 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
