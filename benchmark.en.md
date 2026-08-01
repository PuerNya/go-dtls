# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `6db340ebc19d2834b8e162d5ac719f5748de2dab`
- Generated: `2026-08-01T07:45:16Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `6502cdd34cab185217b44821d2bcba77383ebebe (Linux Release static)`

173 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (14)](#section-connection-lifecycle)
- [Real UDP interoperability (56)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (14)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (14)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (14)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (14)](#real-udp-wolfssl-client-wolfssl-server)
- [Record layer and reliability (37)](#section-record-layer-and-reliability)
- [Key schedule and cryptography (38)](#section-key-schedule-and-cryptography)
- [Wire encoding and parsing (26)](#section-wire-encoding-and-parsing)
- [Certificate compression (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 587.139 us/op | 99289 B/op | 759 allocs/op |
| Full mTLS handshake | 5 | 865.63 us/op | 115951 B/op | 973 allocs/op |
| mTLS session resumption handshake | 5 | 438.955 us/op | 116163 B/op | 804 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 691.413 us/op | 118513 B/op | 983 allocs/op |
| Direct external PSK handshake | 5 | 355.703 us/op | 98268 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.067 ms/op | 130939 B/op | 1041 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.065 ms/op | 123300 B/op | 1020 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.753 ms/op | 172668 B/op | 1505 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.758 ms/op | 160571 B/op | 1466 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 993.193 us/op | 148602 B/op | 1259 allocs/op |
| ECH handshake / via HRR | 5 | 988.989 us/op | 151378 B/op | 1280 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 915.081 us/op | 147125 B/op | 791 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 864.318 us/op | 150325 B/op | 821 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.118 ms/op | 175766 B/op | 839 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.968 ms/conn | 2621408 B/op | 24335 allocs/op |
| 1-RTT application-data round trip | 5 | 1.989 ms/conn | 2630080 B/op | 24655 allocs/op |
| Full mTLS handshake | 5 | 3.59 ms/conn | 3470856 B/op | 33637 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.989 ms/conn | 3110480 B/op | 28265 allocs/op |
| Direct external PSK handshake | 5 | 0.5325 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.99 ms/conn | 2641048 B/op | 25452 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.18 ms/conn | 2811584 B/op | 26007 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.148 ms/conn | 4315648 B/op | 45057 allocs/op |
| Session resumption handshake | 5 | 0.5974 ms/conn | 4503352 B/op | 40579 allocs/op |
| mTLS session resumption handshake | 5 | 0.6459 ms/conn | 6621936 B/op | 53467 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5952 ms/conn | 290696 B/op | 2028 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.331 ms/conn | 3492872 B/op | 24658 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.292 ms/conn | 3531928 B/op | 24978 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.485 ms/conn | 4027592 B/op | 25338 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.673 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.802 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.262 ms/conn | 1522544 B/op | 14882 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.643 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.9988 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.692 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.87 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.164 ms/conn | 1580304 B/op | 15722 allocs/op |
| Session resumption handshake | 5 | 1.088 ms/conn | 2068464 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.089 ms/conn | 2438544 B/op | 22362 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.886 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.893 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.218 ms/conn | 1440048 B/op | 12515 allocs/op |
| 1-RTT application-data round trip | 5 | 4.846 ms/conn | 1944184 B/op | 14005 allocs/op |
| Full mTLS handshake | 5 | 5.881 ms/conn | 1818360 B/op | 16829 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.232 ms/conn | 1504704 B/op | 13435 allocs/op |
| Direct external PSK handshake | 5 | 0.793 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.837 ms/conn | 1956368 B/op | 14565 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.904 ms/conn | 2119704 B/op | 15645 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.57 ms/conn | 2734728 B/op | 23440 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2984824 B/op | 22965 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 212464 B/op | 1138 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.365 ms/conn | 1686768 B/op | 12995 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.443 ms/conn | 1706928 B/op | 13135 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.41 ms/conn | 1865808 B/op | 13315 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.293 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.514 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 8.245 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.284 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 1.016 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.541 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.482 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.3 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.399 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.545 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.602 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.55 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 54.82 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 614.5 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.41 us/op | 1699.73 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.025 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 235.2 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 185 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.113 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.036 us/op | 1014.76 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 391.8 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 110.2 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 380.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 87.2 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 42.51 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 46.33 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.45 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 591.9 ns/op | 2027.33 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 533 ns/op | 2251.2 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 544.6 ns/op | 2203.29 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 27.093 us/op | 2418.89 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 483.2 ns/op | 2483.41 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 30.41 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.363 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.087 us/op | 388.75 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.137 us/op | 1055.54 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.28 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.552 us/op | 337.86 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.768 us/op | 678.78 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.783 us/op | 136.63 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.147 us/op | 381.32 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.158 us/op | 379.94 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.957 us/op | 303.24 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.196 us/op | 1003.03 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.796 us/op | 316.16 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.259 us/op | 953.45 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.318 us/op | 910.33 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.724 us/op | 696.25 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.882 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.429 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.571 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.582 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.519 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.48 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.431 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.226 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.439 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.583 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.259 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.507 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 715.4 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.847 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.963 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 747 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 739.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.595 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.284 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.95 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.638 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.643 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.35 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.68 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.664 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.899 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.845 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.813 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.701 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.641 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 350.2 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 673.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 118.6 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.4 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 297.7 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.2 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 341.9 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 582.2 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 53.57 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 530.7 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 87.02 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 48.8 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 71.17 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 744 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 89.57 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 72.48 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 66.29 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 672.1 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 439.4 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 67.1 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 107.5 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 69.49 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.59 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 313.2 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 158.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 76.13 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.209 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 847.7 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 839.9 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.78 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 67.12 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.763 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.754 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
