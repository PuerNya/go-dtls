# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-02T19:27:38Z`
- Go: `go version go1.26.7 linux/amd64`
- Platform: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `9e0d3894d28551cdc5ab97200348aaf1cb8e3a85 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 466.465 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 718.739 us/op | 108519 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 383.293 us/op | 116048 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 932.616 us/op | 116175 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.241 ms/op | 135134 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 559.566 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 759.13 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 760.959 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 314.083 us/op | 98142 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 871.277 us/op | 126177 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 877.702 us/op | 118537 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.444 ms/op | 165285 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.453 ms/op | 153249 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 841.35 us/op | 143819 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 864.619 us/op | 146596 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 760.148 us/op | 142339 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 754.046 us/op | 145539 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.877 ms/op | 170981 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.484 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.512 ms/conn | 2419088 B/op | 21116 allocs/op |
| Full mTLS handshake | 5 | 2.611 ms/conn | 3085072 B/op | 28762 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.884 ms/conn | 3579512 B/op | 30777 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.548 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.3449 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.461 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.57 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.979 ms/conn | 3878728 B/op | 38883 allocs/op |
| Session resumption handshake | 5 | 0.4016 ms/conn | 4293592 B/op | 37022 allocs/op |
| mTLS session resumption handshake | 5 | 0.4973 ms/conn | 6303560 B/op | 48608 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.4226 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.76 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.788 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.881 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.819 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.937 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 5.044 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.023 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.806 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8695 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.792 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.964 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.021 ms/conn | 1380752 B/op | 12198 allocs/op |
| Session resumption handshake | 5 | 0.8944 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.9217 ms/conn | 2228440 B/op | 18824 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.013 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.995 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.803 ms/conn | 1229072 B/op | 8976 allocs/op |
| 1-RTT application-data round trip | 5 | 4.087 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 4.699 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.733 ms/conn | 1869688 B/op | 14948 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.754 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.606 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.094 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.2 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.09 ms/conn | 2640264 B/op | 22181 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.932 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.897 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 5.462 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.68 ms/conn | 34976 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 6.294 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 6.285 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.436 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.512 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.779 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.278 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.444 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 9.322 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1013 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.593 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.973 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.933 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 41.47 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 53 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 638 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.368 us/op | 1729.92 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.991 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 223.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 179.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.104 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.649 us/op | 1122.61 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 377.7 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 108.3 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 359.8 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 91.23 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 46.17 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 50.23 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 50.33 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 563.5 ns/op | 2129.42 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 519.7 ns/op | 2309.2 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 523.8 ns/op | 2291.01 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.851 us/op | 2440.76 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 503.5 ns/op | 2383.42 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 31.59 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.851 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.271 us/op | 366.9 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.179 us/op | 1017.42 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 13.23 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 4.266 us/op | 281.32 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.744 us/op | 688.08 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.038 us/op | 132.77 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.434 us/op | 349.47 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.44 us/op | 348.84 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.958 us/op | 242.04 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.447 us/op | 829.18 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.526 us/op | 340.35 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.353 us/op | 886.85 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.469 us/op | 817.03 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.983 us/op | 605.26 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.73 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.779 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.559 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.673 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.48 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.629 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.037 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.11 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.365 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.402 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.785 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 11.719 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.028 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.963 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 670.1 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.808 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 20.18 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 705.1 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 714.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.78 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.735 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 20.17 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.828 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.824 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.318 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.671 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.772 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.826 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.952 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.441 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.765 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 7.225 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 384.5 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 704 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 134.1 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 84.54 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 339.5 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 269.7 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 360.8 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 578.4 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 53.35 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 535.6 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 87.45 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 46.86 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.68 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 719.1 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 88.48 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.45 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 61.33 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 627.1 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 397.2 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 12.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 64.64 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 100.9 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 62.33 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 23.95 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 301 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 160.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 73.52 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.17 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 811.7 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 807.5 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 11.39 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 64.25 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.915 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.363 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
