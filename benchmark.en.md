# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-05T11:52:30Z`
- Go: `go version go1.26.7 linux/amd64`
- Platform: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `d72f6d9e4e85ffcadfa0c737959dc26b8717947a (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 385.45 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 586.357 us/op | 108518 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 334.194 us/op | 116026 B/op | 805 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 782.838 us/op | 116043 B/op | 1038 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.049 ms/op | 135103 B/op | 1333 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 485.231 us/op | 113728 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 653.034 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 651.102 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 266.628 us/op | 98143 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 729.335 us/op | 126176 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 730.013 us/op | 118535 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.209 ms/op | 165285 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.218 ms/op | 153245 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 685.076 us/op | 143821 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 702.948 us/op | 146597 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 627.473 us/op | 142337 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 615.497 us/op | 145538 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.585 ms/op | 170979 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.225 ms/conn | 2410416 B/op | 20796 allocs/op |
| 1-RTT application-data round trip | 5 | 1.225 ms/conn | 2419088 B/op | 21116 allocs/op |
| Full mTLS handshake | 5 | 2.206 ms/conn | 3086016 B/op | 28767 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.435 ms/conn | 3577976 B/op | 30769 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.32 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.2869 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.216 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.302 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.295 ms/conn | 3877736 B/op | 38875 allocs/op |
| Session resumption handshake | 5 | 0.3387 ms/conn | 4293592 B/op | 37021 allocs/op |
| mTLS session resumption handshake | 5 | 0.4346 ms/conn | 6302368 B/op | 48611 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3715 ms/conn | 279896 B/op | 1848 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.481 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.495 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.438 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.261 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 3.357 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 4.345 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.318 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.297 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 0.8161 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.291 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.427 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.319 ms/conn | 1375632 B/op | 12189 allocs/op |
| Session resumption handshake | 5 | 0.8469 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 0.8359 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.523 ms/conn | 1829504 B/op | 10763 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.287 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.437 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 3.499 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 3.944 ms/conn | 1665784 B/op | 14422 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.972 ms/conn | 1869688 B/op | 14948 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.468 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.504 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.484 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.503 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 8.626 ms/conn | 2640016 B/op | 22178 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.67 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.486 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 4.778 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.876 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 5.44 ms/conn | 543488 B/op | 1183 allocs/op |
| Full mTLS handshake | 5 | 5.28 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.324 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.855 ms/conn | 34976 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 0.616 ms/conn | 34976 B/op | 54 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.398 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.218 ms/conn | 543808 B/op | 1183 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.791 ms/conn | 551680 B/op | 1184 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1011 ms/pair | 557832 B/op | 1186 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.963 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.163 ms/conn | 34912 B/op | 54 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.836 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 34.62 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 43.54 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 558.3 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.018 us/op | 2029.52 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.547 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 196.9 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 149.2 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 905.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.164 us/op | 1294.63 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 324.8 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 91.69 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 310.9 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 77.58 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 39.35 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 42.54 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 43.08 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 494.4 ns/op | 2427.28 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 464.8 ns/op | 2581.9 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 466.9 ns/op | 2570.19 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 23.383 us/op | 2802.73 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 439.7 ns/op | 2729.08 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 26.24 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.226 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.943 us/op | 407.7 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 950.4 ns/op | 1262.67 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.91 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.684 us/op | 325.75 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.441 us/op | 832.93 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.926 us/op | 151.4 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.509 us/op | 478.31 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.608 us/op | 460.07 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.754 us/op | 319.69 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.218 us/op | 984.87 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.06 us/op | 392.13 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.146 us/op | 1047.06 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.185 us/op | 1012.63 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.758 us/op | 682.46 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.243 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 5.519 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.278 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.021 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.235 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.966 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.8666 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9552 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.134 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.785 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 5.393 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 8.968 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 5.789 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 13.936 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 555.9 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.5 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 16.89 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 585.3 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 587.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.378 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.567 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 16.9 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.416 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.415 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.952 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.236 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.016 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.539 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.188 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 5.75 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.043 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 5.662 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 354.6 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 615.3 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 117.3 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 71.29 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 293.6 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 227.8 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 283.1 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 491.9 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 45.26 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 423 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 73.48 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 37.6 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 58.59 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 595 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 73.71 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 58.8 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 51.99 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 530.5 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 338.8 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 53.61 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 84.58 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 51.13 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 20.59 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 254.6 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 132.9 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 63.14 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 956.8 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 668.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 668 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 9.727 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 52.79 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.865 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 5.333 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
