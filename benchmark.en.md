# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-31T21:36:19Z`
- Go: `go version go1.26.7 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) 6973P-C`
- wolfSSL: `0827c4bca614cfa46c50d3cdef2bf79edccf67f9 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 336.865 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 515.588 us/op | 108520 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 306.094 us/op | 116104 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 692.934 us/op | 116090 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 912.995 us/op | 135070 B/op | 1333 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 423.798 us/op | 113727 B/op | 911 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 557.652 us/op | 117770 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 560.135 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 232.055 us/op | 98144 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 626.021 us/op | 126175 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 622.334 us/op | 118534 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.039 ms/op | 165284 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.12 ms/op | 153246 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 596.402 us/op | 143822 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 617.05 us/op | 146599 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 547.545 us/op | 142337 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 513.361 us/op | 145537 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.331 ms/op | 170978 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.141 ms/conn | 2410416 B/op | 20796 allocs/op |
| 1-RTT application-data round trip | 5 | 1.366 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 2.022 ms/conn | 3085520 B/op | 28761 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.231 ms/conn | 3578736 B/op | 30777 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.169 ms/conn | 2899488 B/op | 24726 allocs/op |
| Direct external PSK handshake | 5 | 0.2436 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.125 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.218 ms/conn | 2587776 B/op | 22429 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.096 ms/conn | 3876744 B/op | 38870 allocs/op |
| Session resumption handshake | 5 | 0.2934 ms/conn | 4293592 B/op | 37021 allocs/op |
| mTLS session resumption handshake | 5 | 0.3441 ms/conn | 6302040 B/op | 48612 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3295 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.41 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.382 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.269 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.084 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.116 ms/conn | 1181816 B/op | 10829 allocs/op |
| Full mTLS handshake | 5 | 4.008 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.932 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.987 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.9956 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.148 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.161 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.97 ms/conn | 1378704 B/op | 12192 allocs/op |
| Session resumption handshake | 5 | 1.219 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.369 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.293 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.87 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.339 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 3.08 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 3.54 ms/conn | 1665784 B/op | 14422 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.538 ms/conn | 1869672 B/op | 14948 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.354 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.467 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.086 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.122 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.507 ms/conn | 2640760 B/op | 22187 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 202128 B/op | 963 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.52 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.141 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 4.069 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.559 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 4.444 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 4.649 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.661 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.523 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.586 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.53 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.626 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.615 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1012 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.62 ms/conn | 34912 B/op | 54 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.776 ms/conn | 34960 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.501 ms/conn | 34960 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 28.2 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 35.33 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 489.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.711 us/op | 2394.51 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.367 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 154.6 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 127.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 789.2 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.617 us/op | 1564.91 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 270.6 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 81.12 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 271.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 66.88 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 32.77 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 36.98 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 35.84 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 423.6 ns/op | 2832.62 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 391.6 ns/op | 3064.71 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 390 ns/op | 3076.55 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 21.08 us/op | 3108.94 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 377.7 ns/op | 3176.91 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 20.82 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 2.649 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.604 us/op | 460.75 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 910.4 ns/op | 1318.06 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 7.214 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 4.011 us/op | 299.2 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.549 us/op | 774.57 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.317 us/op | 163.99 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.457 us/op | 488.47 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.883 us/op | 416.28 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.038 us/op | 297.15 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.339 us/op | 895.94 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.886 us/op | 415.75 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.24 us/op | 967.75 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.28 us/op | 937.74 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.519 us/op | 789.89 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 1.92 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 4.753 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.08 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.578 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.047 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.557 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.6214 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.7315 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 968.7 ns/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.355 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 4.505 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 7.552 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 4.951 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 11.934 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 477.4 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.272 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 12.76 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 499.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 500.5 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.199 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.03 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 12.76 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.2 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.187 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.653 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.847 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 3.371 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.288 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.664 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 4.896 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.533 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 4.752 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 277 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 477.8 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 103.7 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 61.95 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 251.1 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 196.2 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 237.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 410.8 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 36.69 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 365.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 58.62 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 30.22 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 47.57 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 497.5 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 58.62 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 47.39 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 44.44 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 432.7 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 282 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 9.187 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 43.46 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 69.61 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 42.99 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 17.18 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 205.6 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 108.2 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 52.01 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 787.2 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 553.4 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 554.7 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 8.391 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 42.7 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.036 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 4.919 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
