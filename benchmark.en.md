# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-24T02:14:23Z`
- Go: `go version go1.26.8 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V45 96-Core Processor`
- wolfSSL: `cb6e4644ae6e29403d79e46795b262b8f90045d9 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 310.481 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 450.177 us/op | 108519 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 245.033 us/op | 116045 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 601.668 us/op | 116156 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 796.48 us/op | 135142 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 364.922 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 502.551 us/op | 117771 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 497.87 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 196.325 us/op | 98145 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 610.429 us/op | 126174 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 600.082 us/op | 118532 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.013 ms/op | 165283 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.018 ms/op | 153245 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 523.958 us/op | 143824 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 517.882 us/op | 146601 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 543.252 us/op | 142335 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 523.268 us/op | 145535 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.279 ms/op | 170977 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.158 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.266 ms/conn | 2419088 B/op | 21116 allocs/op |
| Full mTLS handshake | 5 | 2.072 ms/conn | 3086264 B/op | 28776 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.362 ms/conn | 3579032 B/op | 30776 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.2 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.3002 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.18 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.311 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.179 ms/conn | 3878776 B/op | 38885 allocs/op |
| Session resumption handshake | 5 | 0.3492 ms/conn | 4293592 B/op | 37021 allocs/op |
| mTLS session resumption handshake | 5 | 0.4036 ms/conn | 6302736 B/op | 48621 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3601 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.372 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.354 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.1 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.012 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.02 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 3.855 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.823 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.942 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.7764 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.039 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.185 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.872 ms/conn | 1381776 B/op | 12195 allocs/op |
| Session resumption handshake | 5 | 0.8053 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.7435 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.015 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.489 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.314 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 3.162 ms/conn | 1733512 B/op | 10466 allocs/op |
| Full mTLS handshake | 5 | 3.464 ms/conn | 1666032 B/op | 14425 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.523 ms/conn | 1869936 B/op | 14952 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.301 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.446 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.151 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.433 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 8.054 ms/conn | 2640760 B/op | 22187 allocs/op |
| Session resumption handshake | 5 | 1006 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1006 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.408 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.047 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.893 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.542 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 4.716 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 4.877 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.804 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.558 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.609 ms/conn | 35024 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.155 ms/conn | 554632 B/op | 1186 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.089 ms/conn | 543880 B/op | 1185 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.091 ms/conn | 551752 B/op | 1185 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1010 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.667 ms/conn | 34912 B/op | 54 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.485 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.589 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 27.19 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 36.18 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 547.1 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.903 us/op | 2152.49 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.479 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 160.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 134 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 953.2 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.971 us/op | 1378.58 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 290.1 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 52.84 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 279.2 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 58.48 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 24.19 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 27.66 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 25.4 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 307.1 ns/op | 3907.42 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 327 ns/op | 3669.95 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 296.9 ns/op | 4041.79 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 12.862 us/op | 5095.21 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 272.1 ns/op | 4410.3 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 21.49 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 2.16 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.532 us/op | 474 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 885.4 ns/op | 1355.25 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 5.167 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.757 us/op | 435.19 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.179 us/op | 1017.88 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.123 us/op | 195.98 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.764 us/op | 434.13 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.727 us/op | 439.98 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.341 us/op | 359.2 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.041 us/op | 1152.36 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.747 us/op | 436.84 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 959.7 ns/op | 1250.36 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 944.4 ns/op | 1270.63 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.28 us/op | 937.72 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 1.668 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 3.758 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.029 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.165 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 944.6 ns/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.17 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.4697 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.4693 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 848.3 ns/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 1.857 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 3.751 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 6.923 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 4.523 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 9.91 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 416.6 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.106 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 12.7 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 441.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 432.8 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 987.3 ns/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 2.606 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 12.77 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.085 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.005 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.386 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.833 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 2.85 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.242 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.316 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 3.863 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.185 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 3.994 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 245.3 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 413.9 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 83.22 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 53.39 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 184.4 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 141.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 205.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 313.6 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 28.71 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 308.6 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 49.52 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 24.09 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 35.64 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 402.2 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 50.73 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 35.82 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 34.98 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 440.4 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 325.2 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 7.346 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 49.44 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 59.92 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 40.08 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 14.25 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 195 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 97.23 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 41.84 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 712.5 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 520.7 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 515.4 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 6.825 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 46 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.874 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 4.033 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
