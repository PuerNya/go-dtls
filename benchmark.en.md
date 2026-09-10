# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-10T02:04:59Z`
- Go: `go version go1.26.7 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `765d2165e3b7e5fc84cc41ba16d8077f886cef17 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 457.13 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 682.289 us/op | 108518 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 343.418 us/op | 116064 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 907.33 us/op | 116330 B/op | 1039 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.214 ms/op | 135179 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 535.586 us/op | 113727 B/op | 911 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 733.158 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 735.252 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 283.282 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 824.492 us/op | 126176 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 824.59 us/op | 118536 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.376 ms/op | 165285 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.383 ms/op | 153119 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 786.101 us/op | 143819 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 789.301 us/op | 146595 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 695.5 us/op | 142338 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 654.622 us/op | 145538 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.729 ms/op | 170980 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.616 ms/conn | 2410464 B/op | 20797 allocs/op |
| 1-RTT application-data round trip | 5 | 1.631 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 2.823 ms/conn | 3085768 B/op | 28770 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.031 ms/conn | 3579064 B/op | 30776 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.646 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.3996 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.608 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.709 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.24 ms/conn | 3880736 B/op | 38900 allocs/op |
| Session resumption handshake | 5 | 0.4492 ms/conn | 4293592 B/op | 37025 allocs/op |
| mTLS session resumption handshake | 5 | 0.5178 ms/conn | 6303480 B/op | 48624 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.4917 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.894 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.865 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.922 ms/conn | 3816568 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.933 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.042 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 5.156 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.179 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.93 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8143 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.955 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.12 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.174 ms/conn | 1372664 B/op | 12187 allocs/op |
| Session resumption handshake | 5 | 0.8348 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.8441 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.074 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.912 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.855 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.054 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 4.895 ms/conn | 1666032 B/op | 14425 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.953 ms/conn | 1869640 B/op | 14948 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.868 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.628 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.066 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.076 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.52 ms/conn | 2640440 B/op | 22181 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.047 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.803 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.485 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.823 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 6.662 ms/conn | 543512 B/op | 1184 allocs/op |
| Full mTLS handshake | 5 | 6.993 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.844 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.666 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.837 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.386 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.515 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.599 ms/conn | 551752 B/op | 1186 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 558424 B/op | 1185 allocs/op |
| mTLS session resumption handshake | 5 | 1013 ms/pair | 557784 B/op | 1185 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.735 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.522 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 8.088 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 30.79 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 39.4 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 425.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.539 us/op | 2661.39 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.433 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 176.6 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 145.9 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 832.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.815 us/op | 1455.16 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 262.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 76.39 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 273.3 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 65.79 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 28.63 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 38.21 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 35.8 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 365.7 ns/op | 3281.33 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 333.2 ns/op | 3601.73 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 333.7 ns/op | 3595.81 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 17.799 us/op | 3681.93 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 332.7 ns/op | 3606.9 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 21.92 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.276 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.784 us/op | 431.04 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 977.2 ns/op | 1227.94 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 8.469 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.358 us/op | 357.35 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.537 us/op | 780.78 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.216 us/op | 146.05 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.262 us/op | 367.83 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.39 us/op | 353.97 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.887 us/op | 308.73 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.212 us/op | 990.38 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.301 us/op | 363.51 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.203 us/op | 997.55 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.225 us/op | 979.49 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.564 us/op | 767.03 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.113 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 4.955 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.168 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.666 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.139 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.638 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.091 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9583 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.075 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.492 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 4.552 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 7.644 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 5.467 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 12.489 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 535.3 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.363 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 7.101 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 553.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 553 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.228 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.191 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 7.097 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.253 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.248 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.704 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.94 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 3.44 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.416 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.794 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 5.072 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.704 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 4.956 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 276.6 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 492.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 91.71 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 68.37 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 241 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 204.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 274.8 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 401.4 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 39.65 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 372.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 63.2 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 36.59 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 54.45 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 505 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 64.27 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 54.01 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 59.98 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 442.3 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 324.2 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.37 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 48.4 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 79.7 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 55.92 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 27.65 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 222.4 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 128.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 67.17 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 842.1 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 620 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 615 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 9.566 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 47.91 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.697 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 4.808 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
