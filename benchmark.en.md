# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `8010fee127101d213fb208cbfbff51779d6c6954`
- Generated: `2026-10-06T21:04:43Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `e61f90d0c24a468c11cb0fe09fa87b8e666f0860 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 565.227 us/op | 97817 B/op | 709 allocs/op |
| Full mTLS handshake | 5 | 850.088 us/op | 111974 B/op | 891 allocs/op |
| mTLS session resumption handshake | 5 | 421.799 us/op | 119146 B/op | 822 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.094 ms/op | 119838 B/op | 1065 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.473 ms/op | 138698 B/op | 1359 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 657.858 us/op | 117041 B/op | 933 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 905.351 us/op | 121225 B/op | 963 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 906.836 us/op | 121225 B/op | 963 allocs/op |
| Direct external PSK handshake | 5 | 354.362 us/op | 101261 B/op | 743 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.028 ms/op | 129634 B/op | 992 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.027 ms/op | 121852 B/op | 970 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.713 ms/op | 169030 B/op | 1425 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.717 ms/op | 156580 B/op | 1384 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 962.79 us/op | 147130 B/op | 1209 allocs/op |
| ECH handshake / via HRR | 5 | 967.509 us/op | 149906 B/op | 1230 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 878.489 us/op | 145796 B/op | 742 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 829.026 us/op | 149141 B/op | 773 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.117 ms/op | 174582 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.912 ms/conn | 2482048 B/op | 21235 allocs/op |
| 1-RTT application-data round trip | 5 | 1.968 ms/conn | 2490720 B/op | 21555 allocs/op |
| Full mTLS handshake | 5 | 3.355 ms/conn | 3162712 B/op | 29324 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.617 ms/conn | 3658432 B/op | 31344 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.951 ms/conn | 2971120 B/op | 25165 allocs/op |
| Direct external PSK handshake | 5 | 0.531 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.971 ms/conn | 2502064 B/op | 22354 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.091 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.022 ms/conn | 4023480 B/op | 39433 allocs/op |
| Session resumption handshake | 5 | 0.5855 ms/conn | 4431512 B/op | 37859 allocs/op |
| mTLS session resumption handshake | 5 | 0.7197 ms/conn | 6449872 B/op | 49583 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6281 ms/conn | 287200 B/op | 1896 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.321 ms/conn | 3353512 B/op | 21558 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.268 ms/conn | 3395432 B/op | 21898 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.556 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.781 ms/conn | 1186064 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 4.953 ms/conn | 1204784 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.461 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.529 ms/conn | 1341744 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.789 ms/conn | 1424464 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9753 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.82 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.997 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.282 ms/conn | 1396688 B/op | 12383 allocs/op |
| Session resumption handshake | 5 | 1.123 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 1.163 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.097 ms/conn | 1857720 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.228 ms/conn | 1870200 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.176 ms/conn | 1305760 B/op | 9219 allocs/op |
| 1-RTT application-data round trip | 5 | 4.943 ms/conn | 2406104 B/op | 10725 allocs/op |
| Full mTLS handshake | 5 | 5.917 ms/conn | 1743336 B/op | 14712 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.962 ms/conn | 1947192 B/op | 15236 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.164 ms/conn | 1370720 B/op | 10139 allocs/op |
| Direct external PSK handshake | 5 | 0.813 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.928 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.975 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.57 ms/conn | 3314392 B/op | 22480 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3492320 B/op | 19887 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 238256 B/op | 990 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.38 ms/conn | 1555648 B/op | 9719 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.446 ms/conn | 1575808 B/op | 9859 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.771 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.554 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.886 ms/conn | 1169544 B/op | 1203 allocs/op |
| Full mTLS handshake | 5 | 8.25 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.191 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.326 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.016 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.718 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.903 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.35 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 1190872 B/op | 1203 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.547 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.759 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.795 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.09 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 52.89 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 733 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.526 us/op | 1621.32 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.06 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 227 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 188.8 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.197 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 4.248 us/op | 964.22 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 389.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 94.1 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 392.4 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 74.14 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 22.15 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 43.32 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.52 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 551.5 ns/op | 2175.77 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 500.3 ns/op | 2398.67 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 511.4 ns/op | 2346.56 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 26.264 us/op | 2495.25 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 454.8 ns/op | 2638.78 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 28.18 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.375 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.06 us/op | 392.2 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.076 us/op | 1115.17 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.56 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.307 us/op | 362.88 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.63 us/op | 735.99 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.394 us/op | 142.97 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.139 us/op | 382.25 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.176 us/op | 377.79 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.938 us/op | 304.76 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.239 us/op | 968.87 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.55 us/op | 338.03 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.206 us/op | 994.67 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.23 us/op | 975.34 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.623 us/op | 739.38 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.838 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.367 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.536 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.485 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.459 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.454 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.246 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9371 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.396 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.165 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.959 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 11.128 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.114 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.792 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 702.5 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.837 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.947 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 744.5 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 724.9 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.617 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.475 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.972 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.741 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.722 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.262 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.569 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.686 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.843 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 4.296 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.235 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.864 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.984 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 335 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 616.1 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 111 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 78.96 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 293.8 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 246.2 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 358.6 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 565.1 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 51.17 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 584.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 91.21 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 44.5 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 73.21 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 783.1 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 89.4 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 69.85 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 65.65 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 664.7 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 450.8 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 67.96 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 93.04 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 61.76 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.76 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 310.9 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 163.6 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 77.11 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.141 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 826 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 824.5 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.47 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 67.65 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.933 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.315 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
