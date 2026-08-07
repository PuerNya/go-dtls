# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-07T02:23:04Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `80d296934876a0134a476322da8fb3595e307ece (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 455.45 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 705.712 us/op | 108518 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 357.416 us/op | 116072 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 908.826 us/op | 116177 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.229 ms/op | 135124 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 540.987 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 741.994 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 746.506 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 299.049 us/op | 98142 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 836.465 us/op | 126177 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 834.336 us/op | 118536 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.395 ms/op | 165284 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.4 ms/op | 153250 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 818.361 us/op | 143851 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 821.296 us/op | 146596 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 701.945 us/op | 142338 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 657.882 us/op | 145538 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.732 ms/op | 170980 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.611 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.626 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 2.839 ms/conn | 3085816 B/op | 28764 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.051 ms/conn | 3579528 B/op | 30782 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.628 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.4003 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.606 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.696 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.236 ms/conn | 3879200 B/op | 38883 allocs/op |
| Session resumption handshake | 5 | 0.4607 ms/conn | 4293592 B/op | 37023 allocs/op |
| mTLS session resumption handshake | 5 | 0.52 ms/conn | 6303312 B/op | 48630 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.4637 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.905 ms/conn | 3281848 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.88 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.921 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.879 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.995 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 5.135 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.136 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.869 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.7981 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.885 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.047 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.104 ms/conn | 1372560 B/op | 12186 allocs/op |
| Session resumption handshake | 5 | 0.817 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.8224 ms/conn | 2228440 B/op | 18824 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.119 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.891 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.805 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 3.981 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 4.887 ms/conn | 1665784 B/op | 14422 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.902 ms/conn | 1870416 B/op | 14957 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.801 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.619 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.964 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.987 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.28 ms/conn | 2640512 B/op | 22184 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.909 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.708 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.279 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.797 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 6.237 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 6.671 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.658 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.607 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.736 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.232 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.273 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.275 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1013 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.777 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.911 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.925 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 30.97 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 39.17 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 435.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.721 us/op | 2379.34 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.516 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 185.3 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 148.8 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 883.9 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.057 us/op | 1339.96 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 277.5 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 92.77 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 288.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 60.7 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 26.87 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 38.13 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 35.75 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 363 ns/op | 3305.44 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 334.6 ns/op | 3586.91 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 335.6 ns/op | 3575.17 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 17.527 us/op | 3739.12 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 321.7 ns/op | 3730.54 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 22.54 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.55 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.79 us/op | 430.18 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.034 us/op | 1160.36 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 8.741 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.567 us/op | 336.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.588 us/op | 755.81 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.799 us/op | 153.87 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.111 us/op | 385.73 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.17 us/op | 378.6 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.702 us/op | 324.15 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.229 us/op | 976.34 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.257 us/op | 368.49 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.31 us/op | 915.83 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.313 us/op | 913.73 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.676 us/op | 716.06 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.118 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 4.934 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.188 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.713 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.155 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.67 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.093 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.094 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.063 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.466 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 4.526 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 7.629 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 5.377 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 12.329 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 532.5 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.503 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 7.081 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 551.1 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 549.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.228 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.165 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 7.086 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.255 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.258 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.679 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.939 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 3.439 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.391 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.802 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 5.211 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.743 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 4.959 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 281 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 523 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 92.1 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 68.46 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 245.3 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 204.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 276.8 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 407.2 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 39.95 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 386.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 62.94 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 36.24 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 53.84 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 520 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 64.12 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 53.82 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 55.26 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 455.7 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 330.2 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.39 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 50.1 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 82.33 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 55.56 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 27.27 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 235.8 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 127.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 67.89 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 859 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 627.6 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 622 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 9.619 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 48.27 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.557 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 4.616 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
