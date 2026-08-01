# Automated benchmark results

[简体中文](benchmark.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `b512a578bc6b664ffddb0c20b5ad38882bfa2941`
- Generated: `2026-08-01T01:49:23Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `6502cdd34cab185217b44821d2bcba77383ebebe (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 600.436 us/op | 99289 B/op | 759 allocs/op |
| Full mTLS handshake | 5 | 884.636 us/op | 115949 B/op | 973 allocs/op |
| mTLS session resumption handshake | 5 | 426.08 us/op | 116188 B/op | 805 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 703.055 us/op | 118511 B/op | 983 allocs/op |
| Direct external PSK handshake | 5 | 362.493 us/op | 98268 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.088 ms/op | 130939 B/op | 1041 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.094 ms/op | 123300 B/op | 1020 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.802 ms/op | 172667 B/op | 1505 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.807 ms/op | 160573 B/op | 1466 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.02 ms/op | 148602 B/op | 1259 allocs/op |
| ECH handshake / via HRR | 5 | 1.024 ms/op | 151378 B/op | 1280 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 921.251 us/op | 147125 B/op | 791 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 869.244 us/op | 150325 B/op | 821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.261 ms/op | 175767 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.105 ms/conn | 2621408 B/op | 24335 allocs/op |
| 1-RTT application-data round trip | 5 | 2.13 ms/conn | 2630128 B/op | 24656 allocs/op |
| Full mTLS handshake | 5 | 3.841 ms/conn | 3406808 B/op | 33631 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.132 ms/conn | 3110480 B/op | 28265 allocs/op |
| Direct external PSK handshake | 5 | 0.5237 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.125 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.241 ms/conn | 2811584 B/op | 26008 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.66 ms/conn | 4314360 B/op | 45056 allocs/op |
| Session resumption handshake | 5 | 0.5892 ms/conn | 4503352 B/op | 40579 allocs/op |
| mTLS session resumption handshake | 5 | 0.6665 ms/conn | 6621920 B/op | 53481 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6043 ms/conn | 290696 B/op | 2028 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.471 ms/conn | 3492872 B/op | 24658 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.434 ms/conn | 3531928 B/op | 24978 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.796 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.943 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 5.011 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.587 ms/conn | 1522544 B/op | 14882 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.929 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.9504 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.96 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.125 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.548 ms/conn | 1582352 B/op | 15724 allocs/op |
| Session resumption handshake | 5 | 0.9856 ms/conn | 2068464 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.007 ms/conn | 2438544 B/op | 22362 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.154 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.161 ms/conn | 1841984 B/op | 10883 allocs/op |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.321 ms/conn | 1440048 B/op | 12515 allocs/op |
| 1-RTT application-data round trip | 5 | 5.098 ms/conn | 1944184 B/op | 14005 allocs/op |
| Full mTLS handshake | 5 | 6.405 ms/conn | 1817864 B/op | 16822 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.343 ms/conn | 1504704 B/op | 13435 allocs/op |
| Direct external PSK handshake | 5 | 0.762 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.107 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.196 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.33 ms/conn | 2734576 B/op | 23438 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 2984824 B/op | 22965 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 212712 B/op | 1141 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.546 ms/conn | 1686768 B/op | 12995 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.575 ms/conn | 1706928 B/op | 13135 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.912 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.648 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 8.033 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 8.575 ms/conn | 34896 B/op | 54 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.693 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.955 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.034 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.089 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.14 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 557784 B/op | 1185 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.759 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.596 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.23 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 40.28 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 49.83 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 536.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.228 us/op | 1838.58 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.835 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 221.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 187 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 992.6 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.738 us/op | 1095.74 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 359.4 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 96.93 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 367.7 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 79.52 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 34.81 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 50.16 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 47.32 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 556.6 ns/op | 2155.95 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 493.6 ns/op | 2431.21 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 501.2 ns/op | 2394.42 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 24.857 us/op | 2636.48 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 468.4 ns/op | 2561.73 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 28.26 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.226 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.007 us/op | 399.07 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.172 us/op | 1023.73 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.3 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.453 us/op | 347.5 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.747 us/op | 686.72 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.267 us/op | 129.49 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.087 us/op | 388.73 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.152 us/op | 380.76 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.198 us/op | 285.87 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.187 us/op | 1010.99 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.956 us/op | 303.37 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.289 us/op | 931.31 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.306 us/op | 918.6 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.835 us/op | 654.08 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.783 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.392 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.551 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.508 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.486 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.414 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.403 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.232 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.113 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.093 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.104 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.13 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 702.8 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.792 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.118 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 727 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 723.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.598 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.157 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.121 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.637 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.643 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.221 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.53 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.547 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.822 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.683 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.668 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.558 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.526 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 337.2 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 659.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 119 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.91 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 307.9 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 264 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 358.9 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 560.2 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 53.57 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 499.1 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 85.74 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 47.13 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 70.95 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 683.2 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 86.23 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.55 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 71.38 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 603.7 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 433.7 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.39 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 61.1 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 105.6 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 71.1 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 35.15 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 297.7 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 163.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 88.1 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.125 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 812.5 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 811.3 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.39 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 60.64 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 7.129 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 5.951 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
