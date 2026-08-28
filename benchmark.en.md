# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-28T07:40:13Z`
- Go: `go version go1.26.7 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
- wolfSSL: `e68ef1a047ef0a9b26130bc37553275f58d0ba71 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 530.221 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 788.74 us/op | 108517 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 408.746 us/op | 116063 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.029 ms/op | 116364 B/op | 1039 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.392 ms/op | 135122 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 622.625 us/op | 113729 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 853.115 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 854.73 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 338.803 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 979.92 us/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 982.578 us/op | 118539 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.64 ms/op | 165288 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.639 ms/op | 153122 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 914.317 us/op | 143818 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 922.517 us/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 822.62 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 789.765 us/op | 145540 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.135 ms/op | 170981 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.809 ms/conn | 2410464 B/op | 20797 allocs/op |
| 1-RTT application-data round trip | 5 | 1.826 ms/conn | 2419088 B/op | 21116 allocs/op |
| Full mTLS handshake | 5 | 3.214 ms/conn | 3085768 B/op | 28764 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.434 ms/conn | 3579552 B/op | 30780 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.839 ms/conn | 2899488 B/op | 24726 allocs/op |
| Direct external PSK handshake | 5 | 0.4033 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.801 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.909 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.804 ms/conn | 3878480 B/op | 38882 allocs/op |
| Session resumption handshake | 5 | 0.4553 ms/conn | 4293592 B/op | 37030 allocs/op |
| mTLS session resumption handshake | 5 | 0.5322 ms/conn | 6303312 B/op | 48618 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.49 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.119 ms/conn | 3281848 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.096 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.406 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.974 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 5.079 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.746 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.68 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.974 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8888 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.955 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.047 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.692 ms/conn | 1375632 B/op | 12189 allocs/op |
| Session resumption handshake | 5 | 0.9232 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.942 ms/conn | 2228440 B/op | 18824 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.21 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.168 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.142 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.321 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 5.245 ms/conn | 1666032 B/op | 14425 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.322 ms/conn | 1869408 B/op | 14945 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.033 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.593 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.338 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.399 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.06 ms/conn | 2639272 B/op | 22169 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.245 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.202 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.112 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.862 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 6.497 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 7.312 ms/conn | 34944 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 7.091 ms/conn | 34944 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.781 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.879 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 6.604 ms/conn | 554632 B/op | 1186 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.987 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.09 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558424 B/op | 1185 allocs/op |
| mTLS session resumption handshake | 5 | 1015 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.902 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.486 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 8.702 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 40.99 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 54.24 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 609.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.302 us/op | 1779.36 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.927 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 237.2 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 182.9 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.073 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.559 us/op | 1150.96 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 381.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 100.3 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 363.6 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 86.71 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 40.43 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 47.09 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 47.3 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 536.1 ns/op | 2238.5 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 493.1 ns/op | 2433.63 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 493.8 ns/op | 2430.34 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 24.509 us/op | 2673.99 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 460.3 ns/op | 2607.22 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 31.26 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.082 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.156 us/op | 380.19 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.105 us/op | 1086.37 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 9.005 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.8 us/op | 315.8 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.703 us/op | 704.45 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.715 us/op | 137.7 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.141 us/op | 382.03 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.374 us/op | 355.61 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.488 us/op | 267.36 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.303 us/op | 921.14 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.467 us/op | 346.12 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.251 us/op | 959.14 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.313 us/op | 913.84 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.89 us/op | 634.82 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.844 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.773 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.615 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.666 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.57 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.63 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.222 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9621 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.413 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.41 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.425 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.665 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.303 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 17.301 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 707 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.935 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 18.45 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 741 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 733.9 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.691 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.387 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 18.47 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.738 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.732 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.393 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.713 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.791 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.917 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.886 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.001 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.717 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.718 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 362.8 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 668.4 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 123.8 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 78.01 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 333.9 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 260.4 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 322.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 550.6 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.12 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 524.9 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 89.77 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 45.68 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 69.3 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 726 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 93.57 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 69.96 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 66.21 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 654.5 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 427.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.66 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 69.83 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 105.2 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 64.8 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 26.1 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 301.7 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 159 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 74.23 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.127 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 778.3 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 766.7 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 10.09 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 67.96 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.886 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 6.022 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
