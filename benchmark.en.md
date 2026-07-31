# Automated benchmark results

[简体中文](benchmark.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `b512a578bc6b664ffddb0c20b5ad38882bfa2941`
- Generated: `2026-07-31T16:52:59Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `41b7a0209abbddc579d3d861f36c0f574ae7e907 (Linux Release static)`

169 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (14)](#section-connection-lifecycle)
- [Real UDP interoperability (52)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (14)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (12)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (13)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (13)](#real-udp-wolfssl-client-wolfssl-server)
- [Record layer and reliability (37)](#section-record-layer-and-reliability)
- [Key schedule and cryptography (38)](#section-key-schedule-and-cryptography)
- [Wire encoding and parsing (26)](#section-wire-encoding-and-parsing)
- [Certificate compression (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 586.707 us/op | 99289 B/op | 759 allocs/op |
| Full mTLS handshake | 5 | 870.204 us/op | 115950 B/op | 973 allocs/op |
| mTLS session resumption handshake | 5 | 431.089 us/op | 116187 B/op | 804 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 682.04 us/op | 118512 B/op | 983 allocs/op |
| Direct external PSK handshake | 5 | 366.787 us/op | 98268 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.034 ms/op | 130938 B/op | 1041 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.047 ms/op | 123300 B/op | 1020 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.725 ms/op | 172667 B/op | 1505 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.73 ms/op | 160571 B/op | 1466 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 993.077 us/op | 148602 B/op | 1259 allocs/op |
| ECH handshake / via HRR | 5 | 996.413 us/op | 151378 B/op | 1280 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 905.079 us/op | 147125 B/op | 791 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 860.261 us/op | 150325 B/op | 821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.088 ms/op | 175766 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.965 ms/conn | 2621456 B/op | 24336 allocs/op |
| 1-RTT application-data round trip | 5 | 1.996 ms/conn | 2630128 B/op | 24656 allocs/op |
| Full mTLS handshake | 5 | 3.481 ms/conn | 3471400 B/op | 33631 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.001 ms/conn | 3110528 B/op | 28266 allocs/op |
| Direct external PSK handshake | 5 | 0.5129 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.926 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.04 ms/conn | 2811760 B/op | 26008 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.031 ms/conn | 4316408 B/op | 45071 allocs/op |
| Session resumption handshake | 5 | 0.5705 ms/conn | 4503400 B/op | 40579 allocs/op |
| mTLS session resumption handshake | 5 | 0.6273 ms/conn | 6622496 B/op | 53481 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5905 ms/conn | 290696 B/op | 2028 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.288 ms/conn | 3492888 B/op | 24658 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.309 ms/conn | 3531928 B/op | 24978 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.456 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.736 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.826 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.455 ms/conn | 1522544 B/op | 14882 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.773 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.9095 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.636 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.837 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.35 ms/conn | 1584400 B/op | 15727 allocs/op |
| Session resumption handshake | 5 | 1.06 ms/conn | 2068464 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.063 ms/conn | 2438544 B/op | 22362 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.908 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.891 ms/conn | 1841904 B/op | 10882 allocs/op |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.186 ms/conn | 1440048 B/op | 12515 allocs/op |
| 1-RTT application-data round trip | 5 | 4.805 ms/conn | 1944184 B/op | 14005 allocs/op |
| Full mTLS handshake | 5 | 5.851 ms/conn | 1817536 B/op | 16819 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.186 ms/conn | 1504704 B/op | 13435 allocs/op |
| Direct external PSK handshake | 5 | 0.791 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.745 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.802 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.32 ms/conn | 2734592 B/op | 23438 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2984824 B/op | 22965 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 212712 B/op | 1141 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.375 ms/conn | 1686768 B/op | 12995 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.375 ms/conn | 1706928 B/op | 13135 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.441 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.518 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.477 ms/conn | 543560 B/op | 1185 allocs/op |
| Full mTLS handshake | 5 | 7.965 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.498 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.919 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.46 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.418 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.11 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1011 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557760 B/op | 1184 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.401 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.88 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.634 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.16 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 56.52 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 688.5 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.376 us/op | 1724.24 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.899 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 238.2 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 188.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.047 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.18 us/op | 979.94 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 382.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 110.2 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 378.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 90.3 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 42.19 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 46.44 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.54 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 568.2 ns/op | 2112.11 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 525 ns/op | 2285.89 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 517.1 ns/op | 2320.45 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 27.012 us/op | 2426.18 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 489.2 ns/op | 2452.74 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 31.54 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.373 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.043 us/op | 394.41 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.158 us/op | 1035.97 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.29 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.359 us/op | 357.29 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.605 us/op | 747.53 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.798 us/op | 136.39 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.056 us/op | 392.67 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.246 us/op | 369.66 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.828 us/op | 313.45 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.162 us/op | 1032.53 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.708 us/op | 323.64 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.158 us/op | 1036.44 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.243 us/op | 965.43 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.57 us/op | 764.52 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.759 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.251 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.505 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.369 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.459 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.309 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.375 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.156 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.371 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.788 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.077 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.718 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 691.2 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.842 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.993 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 732.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 730.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.633 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.084 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.956 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.662 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.645 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.221 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.514 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.424 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.835 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.713 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.684 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.717 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.549 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 330.3 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 660.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 114.9 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.39 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 288.8 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 239.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 340.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 594.1 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.21 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 523.9 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 85.73 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 47.4 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 70.1 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 725.1 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 88.41 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.29 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 66.36 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 639.6 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 451 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 70.16 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 109.3 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 68.96 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.59 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 313.5 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 152.9 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.91 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.157 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 821.3 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 823.1 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.79 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 71.54 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.721 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.652 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
