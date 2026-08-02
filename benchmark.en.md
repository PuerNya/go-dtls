# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-02T03:41:29Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 578.025 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 857.737 us/op | 108518 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 438.926 us/op | 116092 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.105 ms/op | 116314 B/op | 1039 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.508 ms/op | 135130 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 675.626 us/op | 113727 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 910.612 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 909.41 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 370.501 us/op | 98141 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.046 ms/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.044 ms/op | 118540 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.736 ms/op | 165288 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.735 ms/op | 153126 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 988.705 us/op | 143818 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 985.24 us/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 893.579 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 850.076 us/op | 145541 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.11 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.95 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.988 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.433 ms/conn | 3086560 B/op | 28773 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.692 ms/conn | 3579816 B/op | 30785 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.994 ms/conn | 2964976 B/op | 24726 allocs/op |
| Direct external PSK handshake | 5 | 0.5353 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.953 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.08 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.184 ms/conn | 3879968 B/op | 38894 allocs/op |
| Session resumption handshake | 5 | 0.5891 ms/conn | 4294384 B/op | 37039 allocs/op |
| mTLS session resumption handshake | 5 | 0.7059 ms/conn | 6303064 B/op | 48614 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6675 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.337 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.329 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.507 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.735 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 4.843 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.232 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.257 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.748 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 1.016 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.776 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.949 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.304 ms/conn | 1371112 B/op | 12185 allocs/op |
| Session resumption handshake | 5 | 1.116 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.112 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.982 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.992 ms/conn | 1841984 B/op | 10883 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.18 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 4.989 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 5.967 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.952 ms/conn | 1869640 B/op | 14948 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.204 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.805 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.992 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.057 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.93 ms/conn | 2639488 B/op | 22171 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.384 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.371 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.472 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.295 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 7.795 ms/conn | 543560 B/op | 1185 allocs/op |
| Full mTLS handshake | 5 | 7.93 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.298 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.313 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 1.019 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.757 ms/conn | 554632 B/op | 1186 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.708 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.38 ms/conn | 551704 B/op | 1185 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.684 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.94 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.756 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.93 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 56.31 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 691.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.855 us/op | 1434.9 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.274 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 238.3 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 189.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.238 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.437 us/op | 923.23 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 436.8 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 113.5 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 415.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 89.74 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 42.64 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 46.28 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 40.5 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 653.3 ns/op | 1836.76 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 612.4 ns/op | 1959.36 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 622.8 ns/op | 1926.72 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 30.827 us/op | 2125.9 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 610.8 ns/op | 1964.76 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 31.92 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.372 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.557 us/op | 337.36 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.311 us/op | 915.39 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.25 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.95 us/op | 303.78 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.874 us/op | 640.47 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.708 us/op | 123.6 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.535 us/op | 339.48 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.712 us/op | 323.28 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.468 us/op | 268.55 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.424 us/op | 842.95 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.807 us/op | 315.24 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.42 us/op | 845.03 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.468 us/op | 817.27 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.8 us/op | 666.67 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 3.034 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.707 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.67 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.707 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.621 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.596 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.253 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.252 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.512 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.378 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 7.427 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 11.586 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.825 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 17.294 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 764.9 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.989 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.951 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 792.1 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 780.5 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.68 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.433 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.967 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.718 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.723 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.572 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 3.011 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 5.049 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 2.043 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 4.243 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.259 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 4.109 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 7.158 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 383.2 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 737 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 119.8 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.3 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 300.5 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.7 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 351.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 640.5 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 58.54 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 575.6 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 94.73 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 49.9 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 77.52 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 806.3 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 95.63 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 80.08 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 63.81 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 708 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 461 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 73.54 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 125.8 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 83.44 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.63 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 385.9 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 186.3 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.54 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.385 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 908.7 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 904 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.8 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 72.12 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 6.835 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 7.007 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
