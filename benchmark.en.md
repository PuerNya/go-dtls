# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `bdf1290d8aa1edc71ac136266046d72816d5ae14`
- Generated: `2026-09-12T02:10:36Z`
- Go: `go version go1.26.8 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `c74d6c669e9622dec03eff50931d66141b896433 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 613.245 us/op | 94505 B/op | 688 allocs/op |
| Full mTLS handshake | 5 | 899.261 us/op | 108517 B/op | 865 allocs/op |
| mTLS session resumption handshake | 5 | 448.435 us/op | 116037 B/op | 804 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.169 ms/op | 116394 B/op | 1040 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.573 ms/op | 135155 B/op | 1334 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 701.44 us/op | 113729 B/op | 912 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 956.952 us/op | 117769 B/op | 937 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 956.833 us/op | 117769 B/op | 937 allocs/op |
| Direct external PSK handshake | 5 | 379.425 us/op | 98140 B/op | 724 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.106 ms/op | 126178 B/op | 970 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.108 ms/op | 118540 B/op | 949 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.837 ms/op | 165287 B/op | 1397 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.837 ms/op | 153253 B/op | 1358 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.042 ms/op | 143817 B/op | 1188 allocs/op |
| ECH handshake / via HRR | 5 | 1.045 ms/op | 146594 B/op | 1209 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 943.89 us/op | 142341 B/op | 720 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 895.098 us/op | 145541 B/op | 750 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.28 ms/op | 170982 B/op | 768 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.091 ms/conn | 2410416 B/op | 20796 allocs/op |
| 1-RTT application-data round trip | 5 | 2.126 ms/conn | 2419040 B/op | 21115 allocs/op |
| Full mTLS handshake | 5 | 3.726 ms/conn | 3086064 B/op | 28774 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.024 ms/conn | 3579816 B/op | 30792 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.147 ms/conn | 2899440 B/op | 24725 allocs/op |
| Direct external PSK handshake | 5 | 0.5295 ms/conn | 1659104 B/op | 14233 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.084 ms/conn | 2430328 B/op | 21912 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.207 ms/conn | 2587728 B/op | 22428 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.428 ms/conn | 3879720 B/op | 38875 allocs/op |
| Session resumption handshake | 5 | 0.5824 ms/conn | 4293592 B/op | 37039 allocs/op |
| mTLS session resumption handshake | 5 | 0.7092 ms/conn | 6302320 B/op | 48617 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6254 ms/conn | 280144 B/op | 1851 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.45 ms/conn | 3281832 B/op | 21118 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.445 ms/conn | 3320872 B/op | 21438 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.792 ms/conn | 3816552 B/op | 21798 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.028 ms/conn | 1157904 B/op | 10562 allocs/op |
| 1-RTT application-data round trip | 5 | 5.142 ms/conn | 1176624 B/op | 10822 allocs/op |
| Full mTLS handshake | 5 | 6.628 ms/conn | 1311024 B/op | 11342 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.642 ms/conn | 1311024 B/op | 11342 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 5.047 ms/conn | 1396304 B/op | 12122 allocs/op |
| Direct external PSK handshake | 5 | 1.007 ms/conn | 784944 B/op | 6702 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.016 ms/conn | 1165104 B/op | 11122 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.188 ms/conn | 1412624 B/op | 12162 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.567 ms/conn | 1373584 B/op | 12187 allocs/op |
| Session resumption handshake | 5 | 1.013 ms/conn | 2069744 B/op | 18002 allocs/op |
| mTLS session resumption handshake | 5 | 1.028 ms/conn | 2228304 B/op | 18822 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.284 ms/conn | 1829440 B/op | 10762 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.383 ms/conn | 1841904 B/op | 10882 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.298 ms/conn | 1229024 B/op | 8975 allocs/op |
| 1-RTT application-data round trip | 5 | 5.226 ms/conn | 1733464 B/op | 10465 allocs/op |
| Full mTLS handshake | 5 | 6.421 ms/conn | 1666280 B/op | 14428 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.355 ms/conn | 1869392 B/op | 14945 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.442 ms/conn | 1293984 B/op | 9895 allocs/op |
| Direct external PSK handshake | 5 | 0.79 ms/conn | 900752 B/op | 7195 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.193 ms/conn | 1745328 B/op | 11025 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.263 ms/conn | 1908984 B/op | 12105 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.73 ms/conn | 2639672 B/op | 22174 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 2774104 B/op | 19425 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 202176 B/op | 964 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.525 ms/conn | 1476032 B/op | 9455 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.686 ms/conn | 1496192 B/op | 9595 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.976 ms/conn | 1655072 B/op | 9775 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.901 ms/conn | 34904 B/op | 53 allocs/op |
| 1-RTT application-data round trip | 5 | 8.628 ms/conn | 543560 B/op | 1185 allocs/op |
| Full mTLS handshake | 5 | 8.761 ms/conn | 34872 B/op | 53 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.79 ms/conn | 34872 B/op | 53 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.738 ms/conn | 34904 B/op | 53 allocs/op |
| Direct external PSK handshake | 5 | 1.044 ms/conn | 34952 B/op | 53 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.359 ms/conn | 554560 B/op | 1184 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.352 ms/conn | 543880 B/op | 1185 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.21 ms/conn | 551704 B/op | 1185 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 558400 B/op | 1184 allocs/op |
| mTLS session resumption handshake | 5 | 1017 ms/pair | 557760 B/op | 1184 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.805 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.767 ms/conn | 34888 B/op | 53 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 12.03 ms/conn | 34888 B/op | 53 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 43.31 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 57.7 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 704.7 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.793 us/op | 1466.33 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.225 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 243.1 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 197.5 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.257 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.611 us/op | 888.4 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 418.6 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 98.52 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 445.6 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 88.27 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 35.09 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 49.25 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 46.12 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 609.1 ns/op | 1970.16 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 547.7 ns/op | 2190.94 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 595.1 ns/op | 2016.55 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.272 us/op | 2494.55 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 567.9 ns/op | 2112.9 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 30.97 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.233 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 4.02 us/op | 298.48 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.48 us/op | 810.88 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 10.93 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 4.243 us/op | 282.8 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.968 us/op | 609.89 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 10.741 us/op | 111.72 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 4.065 us/op | 295.22 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 4.17 us/op | 287.78 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 5.078 us/op | 236.31 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.505 us/op | 797.55 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 4.315 us/op | 278.11 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.545 us/op | 776.73 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.599 us/op | 750.63 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 2.04 us/op | 588.25 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.968 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.7 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.731 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.871 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.678 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.712 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.234 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.235 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.557 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.485 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.844 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 11.388 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 8.077 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 18.019 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 740.7 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.856 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.175 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 759.2 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 754.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.662 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.294 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.159 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.688 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.691 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.695 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 3.112 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 5.177 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 2.067 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.928 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.097 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.903 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.862 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 404.9 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 742 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 124.4 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.9 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 318 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 264 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 378.3 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 599.2 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 55.35 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 507.6 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 87.72 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 49.39 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 74.24 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 709.4 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 89.32 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 74.11 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 77.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 661.8 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 473 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.38 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 68.57 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 119.7 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 80.24 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 35.56 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 341.1 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 193.2 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 85.68 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.251 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 915.7 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 928.9 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.32 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 68.93 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 7.41 us/op | 336 B/op | 4 allocs/op |
| Decompress | 5 | 7.085 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
