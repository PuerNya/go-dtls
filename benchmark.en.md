# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-18T19:08:04Z`
- Go: `go version go1.26.8 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `f48ee13d9ff88a400af738c1e20f1d41f5678959 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 557.089 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 825.267 us/op | 108516 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 422.675 us/op | 116061 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.063 ms/op | 116320 B/op | 1039 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.443 ms/op | 135167 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 643.112 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 881.901 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 883.26 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 350.966 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 999.856 us/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1 ms/op | 118539 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.664 ms/op | 165287 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.668 ms/op | 153251 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 939.428 us/op | 143818 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 940.905 us/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 858.583 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 803.213 us/op | 145540 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.035 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.879 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.921 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.299 ms/conn | 3084608 B/op | 28756 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.626 ms/conn | 3580472 B/op | 30785 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.928 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.5052 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.898 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.018 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.957 ms/conn | 3878232 B/op | 38887 allocs/op |
| Session resumption handshake | 5 | 0.5441 ms/conn | 4293592 B/op | 37038 allocs/op |
| mTLS session resumption handshake | 5 | 0.6487 ms/conn | 6302816 B/op | 48618 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5652 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.223 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.211 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.419 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.658 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.721 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.04 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.133 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.554 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8286 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.603 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.755 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.102 ms/conn | 1370512 B/op | 12186 allocs/op |
| Session resumption handshake | 5 | 1.035 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.037 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.813 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.857 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.131 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.741 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 5.889 ms/conn | 1666032 B/op | 14425 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.894 ms/conn | 1869376 B/op | 14942 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.215 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.864 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.734 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.889 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.21 ms/conn | 2640016 B/op | 22178 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.326 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.409 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.446 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.571 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.585 ms/conn | 543512 B/op | 1184 allocs/op |
| Full mTLS handshake | 5 | 8.25 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 7.996 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.506 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.932 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.566 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.654 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.29 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558472 B/op | 1186 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.422 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.797 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.42 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 39.11 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 49.98 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 534 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.127 us/op | 1925.98 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.748 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 222.2 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 175 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 988.3 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.621 us/op | 1131.15 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 359.1 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 112.3 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 351.5 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 79.09 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 43.38 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 49.24 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.43 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 537.8 ns/op | 2231.36 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 477.2 ns/op | 2514.92 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 492.7 ns/op | 2435.49 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 25.766 us/op | 2543.46 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 439.8 ns/op | 2728.7 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 27.66 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.362 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.66 us/op | 451.13 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 996.3 ns/op | 1204.45 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.54 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.939 us/op | 408.23 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.518 us/op | 790.46 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.743 us/op | 154.98 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.657 us/op | 451.64 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.752 us/op | 436.1 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.525 us/op | 340.41 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.105 us/op | 1085.79 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.345 us/op | 358.69 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.101 us/op | 1090.09 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.153 us/op | 1040.82 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.473 us/op | 814.6 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.717 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.164 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.492 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.357 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.449 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.291 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9377 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.354 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.084 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 5.919 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 9.747 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 6.903 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.474 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 677 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.769 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.909 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 712.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 706.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.551 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.012 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.897 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.582 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.629 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.177 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.477 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.372 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.775 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.606 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.421 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.448 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.28 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 303.5 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 598.1 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 107.1 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.46 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 283.4 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 235.1 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 333.4 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 552.6 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.03 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 521.9 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 84.66 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 45.55 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.77 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 723.6 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 86.17 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 70.74 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 65.48 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 609.9 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 410.2 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 61.16 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 99.48 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 64.49 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.27 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 292.6 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 153.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 77.57 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.123 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 807.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 803.1 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.79 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 61.23 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.75 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.36 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
