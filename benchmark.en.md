# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-09T02:09:52Z`
- Go: `go version go1.26.8 linux/amd64`
- Platform: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `bd753e2f2e6eda55dd85344aad64c33654d0950c (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 389.885 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 590.482 us/op | 108517 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 306.89 us/op | 116065 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 760.574 us/op | 116190 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.027 ms/op | 135134 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 452.439 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 613.346 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 614.838 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 266.061 us/op | 98143 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 716.118 us/op | 126175 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 713.681 us/op | 118534 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.183 ms/op | 165283 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.189 ms/op | 153119 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 687.255 us/op | 143821 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 707.274 us/op | 146597 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 613.287 us/op | 142337 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 601.714 us/op | 145537 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.537 ms/op | 170979 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.205 ms/conn | 2410368 B/op | 20795 allocs/op |
| 1-RTT application-data round trip | 5 | 1.208 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 2.208 ms/conn | 3086016 B/op | 28773 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.435 ms/conn | 3579072 B/op | 30773 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.229 ms/conn | 2899488 B/op | 24726 allocs/op |
| Direct external PSK handshake | 5 | 0.283 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.208 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.308 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.341 ms/conn | 3881728 B/op | 38866 allocs/op |
| Session resumption handshake | 5 | 0.3389 ms/conn | 4294384 B/op | 37022 allocs/op |
| mTLS session resumption handshake | 5 | 0.4094 ms/conn | 6368352 B/op | 48624 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3962 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.479 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.467 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.427 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.056 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.132 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 4.055 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.101 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.172 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.6115 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.047 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.16 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.058 ms/conn | 1376232 B/op | 12192 allocs/op |
| Session resumption handshake | 5 | 0.7196 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.7226 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.249 ms/conn | 1829424 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.008 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.492 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 3.541 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 3.994 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.029 ms/conn | 1869376 B/op | 14940 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.5 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.521 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.584 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.582 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 8.655 ms/conn | 2640264 B/op | 22181 allocs/op |
| Session resumption handshake | 5 | 1006 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1006 ms/pair | 201928 B/op | 961 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.603 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.48 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 4.685 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.909 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 5.491 ms/conn | 543512 B/op | 1184 allocs/op |
| Full mTLS handshake | 5 | 5.345 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.615 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.964 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 0.613 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.45 ms/conn | 554584 B/op | 1185 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.51 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.91 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1011 ms/pair | 557784 B/op | 1185 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.967 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.16 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 8.096 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 30.96 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 40.64 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 538 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.994 us/op | 2054.26 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.473 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 165.3 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 136.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 821.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.116 us/op | 1314.68 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 325.3 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 92.32 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 307 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 74.83 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 39.23 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 42.34 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 42.68 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 490.1 ns/op | 2448.34 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 458.8 ns/op | 2615.68 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 460.1 ns/op | 2608.1 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 23.512 us/op | 2787.32 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 435.2 ns/op | 2757.44 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 23.72 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.275 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.129 us/op | 563.72 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 794 ns/op | 1511.32 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.4 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.216 us/op | 541.45 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.1 us/op | 1090.64 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 6.053 us/op | 198.23 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.127 us/op | 564.16 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.218 us/op | 541.05 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.207 us/op | 374.16 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 827.3 ns/op | 1450.52 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.556 us/op | 469.44 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 798.4 ns/op | 1502.98 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 864.5 ns/op | 1388.06 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.356 us/op | 884.63 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.237 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 5.519 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.27 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.991 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.226 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.916 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.8822 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9632 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.129 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.741 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 5.364 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 8.874 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 5.801 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 13.79 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 556.8 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.499 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 17.15 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 582.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 587 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.364 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.568 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 17.12 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.399 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.401 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.942 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.216 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 3.97 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.516 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.186 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 5.748 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.06 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 5.629 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 257.4 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 505.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 99.89 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 71.56 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 267 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 229.1 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 278.3 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 486.9 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 45.21 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 420.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 73.25 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 37.7 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 58.51 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 590.5 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 73.67 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 58.81 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 51.96 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 519.4 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 333.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.51 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 50.86 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 83.1 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 49.89 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 20.43 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 251.6 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 129.6 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 62.46 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 935.9 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 654.6 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 654.1 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 9.665 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 50.53 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.805 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 5.296 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
