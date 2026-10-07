# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `8010fee127101d213fb208cbfbff51779d6c6954`
- Generated: `2026-10-07T03:18:51Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `ab4a14619a4eea21de84603f5e4f581b6296b3b5 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 564.047 us/op | 97817 B/op | 709 allocs/op |
| Full mTLS handshake | 5 | 844.971 us/op | 111974 B/op | 891 allocs/op |
| mTLS session resumption handshake | 5 | 422.792 us/op | 119167 B/op | 822 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.109 ms/op | 119873 B/op | 1065 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.484 ms/op | 138667 B/op | 1358 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 665.229 us/op | 117040 B/op | 933 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 913.195 us/op | 121225 B/op | 963 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 916.742 us/op | 121225 B/op | 963 allocs/op |
| Direct external PSK handshake | 5 | 352.549 us/op | 101261 B/op | 743 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.025 ms/op | 129635 B/op | 992 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.027 ms/op | 121852 B/op | 970 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.705 ms/op | 169033 B/op | 1425 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.715 ms/op | 156580 B/op | 1384 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 967.604 us/op | 147130 B/op | 1209 allocs/op |
| ECH handshake / via HRR | 5 | 967.823 us/op | 149906 B/op | 1230 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 891.096 us/op | 145797 B/op | 742 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 829.335 us/op | 149141 B/op | 773 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.122 ms/op | 174582 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.913 ms/conn | 2482048 B/op | 21235 allocs/op |
| 1-RTT application-data round trip | 5 | 1.954 ms/conn | 2490720 B/op | 21555 allocs/op |
| Full mTLS handshake | 5 | 3.346 ms/conn | 3162960 B/op | 29326 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.601 ms/conn | 3658536 B/op | 31347 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.952 ms/conn | 2971120 B/op | 25165 allocs/op |
| Direct external PSK handshake | 5 | 0.5253 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.921 ms/conn | 2502016 B/op | 22353 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.076 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.02 ms/conn | 3957032 B/op | 39444 allocs/op |
| Session resumption handshake | 5 | 0.5786 ms/conn | 4431648 B/op | 37859 allocs/op |
| mTLS session resumption handshake | 5 | 0.6418 ms/conn | 6449856 B/op | 49589 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6442 ms/conn | 287200 B/op | 1896 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.319 ms/conn | 3353512 B/op | 21558 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.298 ms/conn | 3395432 B/op | 21898 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.584 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.526 ms/conn | 1186064 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 4.695 ms/conn | 1204784 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.132 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.13 ms/conn | 1341744 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.55 ms/conn | 1424464 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9508 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.624 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.866 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.081 ms/conn | 1396688 B/op | 12383 allocs/op |
| Session resumption handshake | 5 | 1.069 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 1.085 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.039 ms/conn | 1857800 B/op | 10965 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.098 ms/conn | 1870200 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.177 ms/conn | 1305760 B/op | 9219 allocs/op |
| 1-RTT application-data round trip | 5 | 4.858 ms/conn | 2406104 B/op | 10725 allocs/op |
| Full mTLS handshake | 5 | 5.929 ms/conn | 1743336 B/op | 14712 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.974 ms/conn | 1946232 B/op | 15226 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.162 ms/conn | 1370720 B/op | 10139 allocs/op |
| Direct external PSK handshake | 5 | 0.771 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.894 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.953 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.61 ms/conn | 3314392 B/op | 22480 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3492336 B/op | 19887 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 238008 B/op | 987 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.39 ms/conn | 1555648 B/op | 9719 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.502 ms/conn | 1575808 B/op | 9859 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.694 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.542 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.675 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 8.104 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 7.974 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.484 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.025 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.872 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 7.601 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.42 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190280 B/op | 1204 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.494 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.407 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 11.54 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 38.75 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 48.41 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 555.4 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.229 us/op | 1837.97 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.835 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 212.6 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 183.3 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.017 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.737 us/op | 1096.21 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 368.8 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 92.33 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 359.9 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 63.95 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 22.21 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 42.95 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.47 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 527.7 ns/op | 2273.93 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 508.4 ns/op | 2360.41 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 501.7 ns/op | 2391.99 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 24.975 us/op | 2624.08 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 436 ns/op | 2752.06 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 25.36 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.373 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.977 us/op | 403.07 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.042 us/op | 1151.53 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 11.55 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.99 us/op | 401.33 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.617 us/op | 741.9 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.222 us/op | 145.95 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.918 us/op | 411.26 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.902 us/op | 413.5 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.549 us/op | 338.17 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.213 us/op | 989.46 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.482 us/op | 344.61 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.202 us/op | 998.45 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.235 us/op | 971.51 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.625 us/op | 738.32 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.831 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.54 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.53 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.526 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.47 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.389 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 0.9399 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.404 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.271 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.128 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.183 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.154 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.015 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 721.8 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.797 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.948 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 749.9 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 752 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.552 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.131 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.945 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.604 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.6 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.251 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.559 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.628 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.799 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.709 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.659 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.61 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.389 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 309 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 619.2 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 108.3 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 78.93 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 289.1 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 246.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 344.4 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 559.5 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 50.51 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 560.7 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 87.84 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 43.98 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 71.97 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 776.1 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 90.05 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 69.74 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 65.41 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 634.7 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 426.6 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.4 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 63.13 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 93.19 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 61.05 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 28.79 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 289.3 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 157.9 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 77.14 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.138 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 804.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 816.4 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 12.48 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 59.44 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.952 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.631 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
