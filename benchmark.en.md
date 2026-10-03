# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `cbf2e3ffed88e32107fe483dc45b3833bd409f48`
- Generated: `2026-10-03T13:14:53Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V45 96-Core Processor`
- wolfSSL: `2411aae3f74d0fc6ccb09d3d6dfdc69e7b32c632 (Linux Release static)`

181 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (18)](#section-connection-lifecycle)
- [Real UDP interoperability (60)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (15)](#real-udp-wolfssl-client-wolfssl-server)
- [Record layer and reliability (37)](#section-record-layer-and-reliability)
- [Key schedule and cryptography (38)](#section-key-schedule-and-cryptography)
- [Wire encoding and parsing (26)](#section-wire-encoding-and-parsing)
- [Certificate compression (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 292.261 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 448.462 us/op | 108520 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 247.202 us/op | 116121 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 621.504 us/op | 116136 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 820.859 us/op | 135141 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 366.645 us/op | 113727 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 500.162 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 490.353 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 197.638 us/op | 98144 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 579.651 us/op | 126174 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 549.262 us/op | 118532 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 938.375 us/op | 165286 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 939.673 us/op | 153116 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 522.407 us/op | 143823 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 524.385 us/op | 146600 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 509.457 us/op | 142335 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 474.291 us/op | 145535 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.226 ms/op | 170977 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.089 ms/conn | 2410416 B/op | 20796 allocs/op |
| 1-RTT application-data round trip | 5 | 1.119 ms/conn | 2419088 B/op | 21116 allocs/op |
| Full mTLS handshake | 5 | 1.971 ms/conn | 3085816 B/op | 28767 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.139 ms/conn | 3579552 B/op | 30780 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.121 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.2936 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.155 ms/conn | 2430336 B/op | 21913 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.182 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.091 ms/conn | 3878232 B/op | 38880 allocs/op |
| Session resumption handshake | 5 | 0.351 ms/conn | 4293672 B/op | 37021 allocs/op |
| mTLS session resumption handshake | 5 | 0.3765 ms/conn | 6302072 B/op | 48606 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3553 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.424 ms/conn | 3281848 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.308 ms/conn | 3320888 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.068 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.705 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 2.819 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 3.523 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.567 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.7 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.7221 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.714 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.86 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.754 ms/conn | 1374656 B/op | 12189 allocs/op |
| Session resumption handshake | 5 | 0.7358 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.7527 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.634 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.298 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.213 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 2.873 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 3.136 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.107 ms/conn | 1869392 B/op | 14945 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.225 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.391 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.956 ms/conn | 1745336 B/op | 11026 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.117 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.446 ms/conn | 2640016 B/op | 22178 allocs/op |
| Session resumption handshake | 5 | 1006 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1005 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.404 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.975 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.833 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.303 ms/conn | 34976 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 4.35 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 4.248 ms/conn | 34944 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.38 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.425 ms/conn | 34976 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 0.552 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.412 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.562 ms/conn | 543832 B/op | 1184 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.618 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 558424 B/op | 1185 allocs/op |
| mTLS session resumption handshake | 5 | 1010 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.47 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.724 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.688 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 24.69 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 29.63 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 426.5 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.626 us/op | 2519.49 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.318 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 150.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 124.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 760.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.505 us/op | 1635.15 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 230 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 45.07 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 243 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 44.11 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 16.93 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 26.12 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 16.9 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 283.5 ns/op | 4232.16 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 254.1 ns/op | 4721.74 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 273.9 ns/op | 4381.72 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 12.591 us/op | 5204.99 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 310.9 ns/op | 3859.26 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 17.13 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 2.726 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.645 us/op | 453.76 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 970.8 ns/op | 1236.08 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 5.796 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.513 us/op | 477.47 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.252 us/op | 958.81 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.715 us/op | 178.71 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.236 us/op | 536.67 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.037 us/op | 589.16 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.213 us/op | 373.51 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 993.9 ns/op | 1207.31 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.703 us/op | 444 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.019 us/op | 1177.78 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.044 us/op | 1149.32 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.271 us/op | 944.39 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 1.646 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 3.597 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 919.8 ns/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.069 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 919.7 ns/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 1.993 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.6744 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.6855 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 847.9 ns/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 1.851 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 3.623 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 6.207 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 4.146 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 9.105 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 412.8 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.14 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 12.89 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 430.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 431.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 925.6 ns/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 2.476 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 13.53 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 968.8 ns/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.265 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.478 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 2.607 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.084 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.272 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 4.031 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.037 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 3.807 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 227.9 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 401.6 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 77.88 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 54.31 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 182 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 144.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 199.9 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 323.3 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 28.54 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 249.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 46.93 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 26.04 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 38.43 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 360.7 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 46.45 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 39.09 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 29.24 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 370.3 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 276.1 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 7.002 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 34.7 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 53.14 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 33.24 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 13.04 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 156 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 78.88 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 37.04 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 619.4 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 448.3 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 448.2 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 6.66 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 32.18 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 3.038 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 3.796 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
