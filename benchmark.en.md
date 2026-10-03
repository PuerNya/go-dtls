# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `cbf2e3ffed88e32107fe483dc45b3833bd409f48`
- Generated: `2026-10-03T02:56:19Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V45 96-Core Processor`
- wolfSSL: `25c706099dbf04d0ad2a98146ec210388d0424cc (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 311.229 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 469.701 us/op | 108520 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 267.1 us/op | 116002 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 622.514 us/op | 116226 B/op | 1039 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 836.257 us/op | 135141 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 371.793 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 502.259 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 500.965 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 206.287 us/op | 98144 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 551.808 us/op | 126174 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 542.174 us/op | 118532 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 896.77 us/op | 165285 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 900.507 us/op | 153116 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 548.925 us/op | 143823 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 541.375 us/op | 146600 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 468.594 us/op | 142335 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 471.825 us/op | 145536 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.187 ms/op | 170977 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.129 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.111 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 1.964 ms/conn | 3086016 B/op | 28773 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.181 ms/conn | 3580272 B/op | 30783 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.166 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.2918 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.164 ms/conn | 2430336 B/op | 21913 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.222 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 2.978 ms/conn | 3879224 B/op | 38880 allocs/op |
| Session resumption handshake | 5 | 0.3465 ms/conn | 4294384 B/op | 37021 allocs/op |
| mTLS session resumption handshake | 5 | 0.3961 ms/conn | 6302488 B/op | 48615 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3841 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.4 ms/conn | 3281848 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.384 ms/conn | 3320888 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.164 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.558 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 2.476 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 3.369 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.462 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.697 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.4578 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.567 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.681 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.384 ms/conn | 1369488 B/op | 12183 allocs/op |
| Session resumption handshake | 5 | 0.4803 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.4902 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.733 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.374 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.265 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 2.879 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 3.296 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.427 ms/conn | 1869408 B/op | 14943 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.305 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.451 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.061 ms/conn | 1745336 B/op | 11026 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.111 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.342 ms/conn | 2639712 B/op | 22172 allocs/op |
| Session resumption handshake | 5 | 1005 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1005 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.411 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.149 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.934 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.447 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 4.293 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 4.565 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.731 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.486 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.563 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.642 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.593 ms/conn | 543880 B/op | 1185 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.731 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1010 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.681 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.831 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.654 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 23.47 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 27.55 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 379.8 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.213 us/op | 3375.55 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.093 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 139.6 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 118.6 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 680.1 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 2.225 us/op | 1840.6 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 194.7 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 59.47 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 207.7 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 41.67 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 17.99 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 27.42 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 17.88 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 274.4 ns/op | 4373.49 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 265.4 ns/op | 4520.76 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 251.5 ns/op | 4770.54 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 12.422 us/op | 5275.87 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 250.9 ns/op | 4782.18 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 16.6 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 2.956 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 1.984 us/op | 604.86 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 756 ns/op | 1587.2 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 5.688 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.358 us/op | 509.01 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.116 us/op | 1075.37 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 5.76 us/op | 208.33 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 1.896 us/op | 632.94 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 1.989 us/op | 603.32 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 2.867 us/op | 418.53 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 916.8 ns/op | 1308.88 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.251 us/op | 533.19 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 821.8 ns/op | 1460.19 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 835.1 ns/op | 1436.87 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.114 us/op | 1077.52 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 1.713 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 3.767 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 904.2 ns/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 1.966 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 879.1 ns/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 1.938 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.6899 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.6735 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 839.5 ns/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 1.861 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 3.456 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 5.605 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 4.207 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 9.098 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 418 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.063 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 12.97 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 426.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 441.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 885.9 ns/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 2.444 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 13.38 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 946.4 ns/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 915.7 ns/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.266 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 1.462 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 2.548 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.063 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 2.13 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 3.748 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 2.051 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 3.766 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 190.6 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 354.9 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 73.22 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 55.9 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 172.5 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 148.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 191.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 289.3 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 26.69 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 250.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 42.83 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 24.8 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 36.49 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 349.2 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 43.3 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 35.31 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 32.1 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 346 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 247.4 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 7.551 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 33.5 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 51.27 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 31.93 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 13.65 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 160.1 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 79.6 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 38.19 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 620.2 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 456.4 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 454.2 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 7.181 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 31.65 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 3.004 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 3.431 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
