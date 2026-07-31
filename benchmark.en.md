# Automated benchmark results

[简体中文](benchmark.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `196ad5c06d35e7c73aae7fbcd1ce0be1f4e2be5d`
- Generated: `2026-07-31T15:43:55Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `41b7a0209abbddc579d3d861f36c0f574ae7e907 (Linux Release static)`

168 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (13)](#section-connection-lifecycle)
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 583.338 us/op | 99289 B/op | 759 allocs/op |
| Full mTLS handshake | 5 | 836.809 us/op | 115950 B/op | 973 allocs/op |
| mTLS session resumption handshake | 5 | 426.15 us/op | 116082 B/op | 803 allocs/op |
| Direct external PSK handshake | 5 | 340.342 us/op | 98265 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.025 ms/op | 130938 B/op | 1041 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.032 ms/op | 123299 B/op | 1020 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.704 ms/op | 172668 B/op | 1505 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.708 ms/op | 160570 B/op | 1466 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 947.842 us/op | 148598 B/op | 1259 allocs/op |
| ECH handshake / via HRR | 5 | 953.056 us/op | 151374 B/op | 1280 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 874.918 us/op | 147124 B/op | 791 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 827.7 us/op | 150324 B/op | 821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.071 ms/op | 175766 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.965 ms/conn | 2621408 B/op | 24335 allocs/op |
| 1-RTT application-data round trip | 5 | 1.982 ms/conn | 2630128 B/op | 24656 allocs/op |
| Full mTLS handshake | 5 | 3.538 ms/conn | 3406328 B/op | 33631 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.974 ms/conn | 3110480 B/op | 28265 allocs/op |
| Direct external PSK handshake | 5 | 0.5162 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.941 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.101 ms/conn | 2811808 B/op | 26008 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.122 ms/conn | 4315104 B/op | 45058 allocs/op |
| Session resumption handshake | 5 | 0.5727 ms/conn | 4500152 B/op | 40559 allocs/op |
| mTLS session resumption handshake | 5 | 0.626 ms/conn | 6619792 B/op | 53470 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.578 ms/conn | 290440 B/op | 2026 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.281 ms/conn | 3492872 B/op | 24658 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.243 ms/conn | 3531928 B/op | 24978 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.451 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.566 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.671 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.115 ms/conn | 1522544 B/op | 14882 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.579 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8768 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.592 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.726 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.075 ms/conn | 1582352 B/op | 15724 allocs/op |
| Session resumption handshake | 5 | 1.06 ms/conn | 2065264 B/op | 17982 allocs/op |
| mTLS session resumption handshake | 5 | 1.131 ms/conn | 2435480 B/op | 22344 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.753 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.729 ms/conn | 1841904 B/op | 10882 allocs/op |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.237 ms/conn | 1440048 B/op | 12515 allocs/op |
| 1-RTT application-data round trip | 5 | 4.808 ms/conn | 1944184 B/op | 14005 allocs/op |
| Full mTLS handshake | 5 | 5.884 ms/conn | 1818312 B/op | 16828 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.239 ms/conn | 1504704 B/op | 13435 allocs/op |
| Direct external PSK handshake | 5 | 0.812 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.799 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.859 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.32 ms/conn | 2734080 B/op | 23432 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2984824 B/op | 22965 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 212464 B/op | 1138 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.393 ms/conn | 1686768 B/op | 12995 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.383 ms/conn | 1706928 B/op | 13135 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.448 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.532 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.473 ms/conn | 543560 B/op | 1185 allocs/op |
| Full mTLS handshake | 5 | 7.879 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.269 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.941 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.439 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.429 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.2 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557784 B/op | 1185 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.393 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.253 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.32 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 40.95 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 52.42 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 575.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.208 us/op | 1855.12 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.992 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 230 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 181.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.119 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.839 us/op | 1066.85 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 372.6 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 114.7 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 362.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 83.14 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 42.89 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 46.39 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.48 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 564 ns/op | 2127.84 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 496.8 ns/op | 2415.39 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 495.5 ns/op | 2421.75 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.273 us/op | 2494.45 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 442.2 ns/op | 2713.46 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 29.18 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.37 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.282 us/op | 365.63 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.186 us/op | 1011.64 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 10.94 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.029 us/op | 396.2 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.644 us/op | 730.09 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.227 us/op | 130.05 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.474 us/op | 345.4 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.41 us/op | 351.91 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.926 us/op | 305.62 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.126 us/op | 1066.08 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.631 us/op | 330.52 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.176 us/op | 1020 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.229 us/op | 976.15 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.616 us/op | 742.51 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.81 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.384 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.576 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.431 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.493 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.373 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.397 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.151 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.004 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.007 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.196 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.248 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 698.4 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.786 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.954 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 724.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 721.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.559 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.105 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.967 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.602 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.597 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.309 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.565 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.623 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.892 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.713 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.464 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.515 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.383 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 352.9 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 681.3 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 114.3 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.47 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 286 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 239.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 335.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 555.8 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.49 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 528.6 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 85.69 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 46.02 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.78 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 743.4 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 88.03 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.05 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 63.91 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 674.8 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 422.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 65.4 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 101.8 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 66.93 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.61 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 306 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 153.7 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.45 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.17 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 836.3 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 839.8 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.79 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 65.2 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.703 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.585 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
