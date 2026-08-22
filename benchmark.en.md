# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-22T00:54:42Z`
- Go: `go version go1.26.6 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `ac45965828714783de56045d839196afb3234250 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 553.928 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 817.045 us/op | 108519 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 413.748 us/op | 116053 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.071 ms/op | 116410 B/op | 1040 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.453 ms/op | 135173 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 646.302 us/op | 113727 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 878.478 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 881.318 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 346.712 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 996.521 us/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 997.312 us/op | 118538 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.667 ms/op | 165283 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.67 ms/op | 153252 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 939.876 us/op | 143818 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 944.009 us/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 845.8 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 799.271 us/op | 145540 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.023 ms/op | 170981 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.909 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.933 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.401 ms/conn | 3085816 B/op | 28771 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.566 ms/conn | 3579232 B/op | 30780 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.934 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.5146 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.902 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.044 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.94 ms/conn | 3878976 B/op | 38889 allocs/op |
| Session resumption handshake | 5 | 0.5626 ms/conn | 4293592 B/op | 37039 allocs/op |
| mTLS session resumption handshake | 5 | 0.6577 ms/conn | 6302072 B/op | 48608 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5896 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.234 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.205 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.418 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.468 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.573 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 5.937 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.957 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.445 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.7942 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.459 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.644 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.867 ms/conn | 1369488 B/op | 12183 allocs/op |
| Session resumption handshake | 5 | 1.061 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.067 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.673 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.658 ms/conn | 1841984 B/op | 10883 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.142 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.718 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 5.897 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.886 ms/conn | 1869376 B/op | 14942 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.188 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.777 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.731 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.767 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.45 ms/conn | 2640264 B/op | 22181 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.365 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.298 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.335 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.464 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.413 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 8.273 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.272 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.51 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.916 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.449 ms/conn | 554632 B/op | 1186 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.44 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.16 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558424 B/op | 1185 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.375 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.537 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.27 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 39.05 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 49.63 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 534.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.152 us/op | 1902.94 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.715 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 216.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 174.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 966.5 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.615 us/op | 1132.93 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 363.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 112.5 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 348.9 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 78.41 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 43.1 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 50.08 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.43 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 550 ns/op | 2181.83 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 480.5 ns/op | 2497.31 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 502 ns/op | 2390.39 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.113 us/op | 2509.68 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 439.9 ns/op | 2727.96 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 27.38 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.361 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.571 us/op | 466.81 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 966.9 ns/op | 1241.04 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.79 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.88 us/op | 416.7 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.51 us/op | 794.72 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.605 us/op | 157.8 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.603 us/op | 461.07 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.733 us/op | 439.11 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.396 us/op | 353.31 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.082 us/op | 1109.07 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.292 us/op | 364.55 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.041 us/op | 1152.29 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.129 us/op | 1062.52 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.468 us/op | 817.2 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.687 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.155 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.481 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.36 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.442 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.295 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9368 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.342 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.077 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 5.976 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 9.79 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 6.908 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.483 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 672.5 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.758 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.892 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 701.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 696.8 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.548 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.995 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.885 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.58 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.58 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.179 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.49 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.385 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.789 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.614 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.4 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.467 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.324 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 295.1 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 578 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 105.1 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.48 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 277.7 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 331.9 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 552.9 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.16 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 523.7 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 85.16 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 45.84 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.88 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 731.6 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 86.53 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 70.89 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 65.61 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 605.6 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 409.3 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 61.48 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 99.26 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 64.96 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.29 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 290.5 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 153.3 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 77.59 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.124 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 808.6 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 808.6 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.78 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 61.22 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.75 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.177 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
