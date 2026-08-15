# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-08-15T08:40:17Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `889de4784d81cb4346db2f537313e7e699b730e6 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 590.991 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 875.356 us/op | 108517 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 429.555 us/op | 116007 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.143 ms/op | 116448 B/op | 1040 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.543 ms/op | 135201 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 692.719 us/op | 113729 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 942.704 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 947.328 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 363.326 us/op | 98140 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.078 ms/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.084 ms/op | 118539 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.803 ms/op | 165291 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.803 ms/op | 153123 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.007 ms/op | 143817 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 1.017 ms/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 905.258 us/op | 142340 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 850.371 us/op | 145540 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.24 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.058 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 2.091 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.649 ms/conn | 3086560 B/op | 28780 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.917 ms/conn | 3580272 B/op | 30783 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.098 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.5266 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.067 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.2 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.449 ms/conn | 3945008 B/op | 38879 allocs/op |
| Session resumption handshake | 5 | 0.575 ms/conn | 4293640 B/op | 37039 allocs/op |
| mTLS session resumption handshake | 5 | 0.6602 ms/conn | 6302984 B/op | 48624 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6219 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.448 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.412 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.772 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.92 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 5.053 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.549 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.556 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.911 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.9679 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.969 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.174 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.462 ms/conn | 1374960 B/op | 12190 allocs/op |
| Session resumption handshake | 5 | 0.9845 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.9968 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.196 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.227 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.293 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 5.108 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 6.331 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.343 ms/conn | 1869160 B/op | 14942 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.297 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.758 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.148 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.126 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.37 ms/conn | 2640416 B/op | 22183 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.506 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.586 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.922 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.671 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 8.097 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 8.617 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.611 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.87 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 1.042 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.994 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.017 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.99 ms/conn | 551704 B/op | 1185 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.782 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.379 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.94 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 39.06 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 50.04 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 572.2 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.229 us/op | 1837.46 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.899 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 226.5 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 190.8 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.074 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.886 us/op | 1053.93 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 367.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 99.5 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 380.1 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 80.51 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 34.66 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 49.16 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 46.09 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 567.6 ns/op | 2114.23 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 525.4 ns/op | 2283.91 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 537.7 ns/op | 2231.83 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 25.425 us/op | 2577.63 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 477.9 ns/op | 2510.77 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 28.56 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.577 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.154 us/op | 380.45 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.199 us/op | 1000.55 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.28 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.587 us/op | 334.54 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.742 us/op | 688.97 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.332 us/op | 128.59 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.159 us/op | 379.87 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.261 us/op | 367.96 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.197 us/op | 285.89 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.336 us/op | 898.17 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.871 us/op | 309.99 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.282 us/op | 936.35 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.345 us/op | 891.88 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.78 us/op | 674.11 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.76 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.399 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.558 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.508 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.497 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.433 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.408 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.4 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.199 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.122 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.139 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.152 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.151 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 698.3 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.794 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.133 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 723 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 720.1 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.59 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.133 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.151 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.632 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.629 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.229 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.56 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.561 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.827 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.684 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.706 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.572 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.548 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 344.3 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 622 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 117.7 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 88.09 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 311.7 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 263.9 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 360 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 559.2 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 53.64 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 502.8 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 86.1 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 47.71 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 71.88 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 688.9 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 86.24 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 71.84 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 71.3 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 604.6 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 434 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 63.24 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 106.5 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 72.47 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 35.15 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 300.3 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 166.1 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 87.65 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.15 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 823 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 819.8 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.39 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 62.93 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 7.095 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 5.968 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
