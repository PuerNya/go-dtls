# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-04T17:54:13Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `0fb71c572f1e27e281e0b83a95881ddba0073a73 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 555.571 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 827.821 us/op | 108517 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 424.888 us/op | 116058 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.072 ms/op | 116397 B/op | 1040 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.457 ms/op | 135171 B/op | 1333 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 652.18 us/op | 113730 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 885.311 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 891.775 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 350.614 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.037 ms/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.032 ms/op | 118539 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.697 ms/op | 165288 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.711 ms/op | 153123 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 947.314 us/op | 143818 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 950.52 us/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 863.302 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 818.513 us/op | 145540 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.056 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.908 ms/conn | 2410320 B/op | 20794 allocs/op |
| 1-RTT application-data round trip | 5 | 1.952 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.376 ms/conn | 3085272 B/op | 28764 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.596 ms/conn | 3580472 B/op | 30774 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.951 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.5175 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.94 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.036 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.009 ms/conn | 3878064 B/op | 38884 allocs/op |
| Session resumption handshake | 5 | 0.5584 ms/conn | 4293592 B/op | 37039 allocs/op |
| mTLS session resumption handshake | 5 | 0.6593 ms/conn | 6303232 B/op | 48610 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6126 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.313 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.278 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.455 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.609 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.666 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.13 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.101 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.551 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8398 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.521 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.745 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.031 ms/conn | 1371536 B/op | 12185 allocs/op |
| Session resumption handshake | 5 | 1.049 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.093 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.974 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.911 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.181 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.718 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 5.824 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.873 ms/conn | 1869440 B/op | 14945 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.181 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.751 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.697 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.775 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.25 ms/conn | 2640016 B/op | 22178 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.371 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.345 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.44 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.576 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.402 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 7.904 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.129 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.556 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.917 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.845 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.59 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.08 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.445 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.59 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.505 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 38.36 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 49.81 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 533.5 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.136 us/op | 1917.23 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.916 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 221.9 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 177.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.077 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.61 us/op | 1134.7 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 370.8 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 113.7 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 350.9 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 78.35 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 42.63 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 46.08 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.44 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 555.2 ns/op | 2161.22 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 521.7 ns/op | 2300.18 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 525.8 ns/op | 2282.45 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.597 us/op | 2464.04 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 448.5 ns/op | 2675.49 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 28.02 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.366 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.414 us/op | 351.45 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.259 us/op | 953 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.25 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.118 us/op | 384.87 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.669 us/op | 719.01 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.47 us/op | 126.71 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.122 us/op | 384.42 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.557 us/op | 337.4 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.344 us/op | 276.24 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.123 us/op | 1068.98 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.71 us/op | 323.44 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.215 us/op | 987.32 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.297 us/op | 924.92 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.706 us/op | 703.28 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.765 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.235 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.51 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.393 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.462 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.411 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.248 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.384 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.134 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.313 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.156 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.085 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.715 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 693.2 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.816 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.964 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 724.8 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 722.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.568 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.167 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.953 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.626 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.601 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.234 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.66 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.516 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.84 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.731 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.604 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.603 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.505 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 340.5 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 647.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 115.4 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.41 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 297.5 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 338.5 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 565.1 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.92 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 525.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 86.08 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 46.59 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 70.44 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 729.3 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 87.2 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.54 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 63.72 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 629.2 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 412.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 61.75 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 111.1 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 76.6 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.62 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 349.3 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 170.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.38 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.258 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 838 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 817.3 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.78 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 60.99 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.732 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.599 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
