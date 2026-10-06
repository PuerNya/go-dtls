# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `8010fee127101d213fb208cbfbff51779d6c6954`
- Generated: `2026-10-06T15:00:24Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 9V74 80-Core Processor`
- wolfSSL: `6aff4edc1776049c713b41d3685fce43406a0601 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 600.391 us/op | 97817 B/op | 709 allocs/op |
| Full mTLS handshake | 5 | 889.599 us/op | 111973 B/op | 891 allocs/op |
| mTLS session resumption handshake | 5 | 442.852 us/op | 119166 B/op | 822 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.15 ms/op | 119816 B/op | 1065 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.552 ms/op | 138670 B/op | 1358 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 697.295 us/op | 117041 B/op | 933 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 954.889 us/op | 121225 B/op | 963 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 956.052 us/op | 121225 B/op | 963 allocs/op |
| Direct external PSK handshake | 5 | 377.867 us/op | 101261 B/op | 743 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.091 ms/op | 129634 B/op | 992 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.088 ms/op | 121852 B/op | 970 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.826 ms/op | 169029 B/op | 1425 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.824 ms/op | 156578 B/op | 1384 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 1.048 ms/op | 147130 B/op | 1209 allocs/op |
| ECH handshake / via HRR | 5 | 1.04 ms/op | 149906 B/op | 1230 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 922.517 us/op | 145796 B/op | 742 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 864.548 us/op | 149141 B/op | 773 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.329 ms/op | 174582 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.104 ms/conn | 2482048 B/op | 21235 allocs/op |
| 1-RTT application-data round trip | 5 | 2.132 ms/conn | 2490720 B/op | 21555 allocs/op |
| Full mTLS handshake | 5 | 3.661 ms/conn | 3163456 B/op | 29333 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.912 ms/conn | 3658200 B/op | 31333 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.136 ms/conn | 2971120 B/op | 25165 allocs/op |
| Direct external PSK handshake | 5 | 0.5346 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 2.128 ms/conn | 2502016 B/op | 22353 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.234 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.407 ms/conn | 3957200 B/op | 39447 allocs/op |
| Session resumption handshake | 5 | 0.6102 ms/conn | 4431824 B/op | 37859 allocs/op |
| mTLS session resumption handshake | 5 | 0.6997 ms/conn | 6448880 B/op | 49577 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5978 ms/conn | 287200 B/op | 1896 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.474 ms/conn | 3353512 B/op | 21558 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.446 ms/conn | 3395432 B/op | 21898 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.907 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 5.005 ms/conn | 1186064 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 5.137 ms/conn | 1204784 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.556 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.596 ms/conn | 1341744 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.999 ms/conn | 1424464 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9799 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.979 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.202 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.539 ms/conn | 1396688 B/op | 12383 allocs/op |
| Session resumption handshake | 5 | 1.05 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 1.006 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.302 ms/conn | 1857720 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.408 ms/conn | 1870200 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.387 ms/conn | 1305760 B/op | 9219 allocs/op |
| 1-RTT application-data round trip | 5 | 5.263 ms/conn | 2406104 B/op | 10725 allocs/op |
| Full mTLS handshake | 5 | 6.349 ms/conn | 1743336 B/op | 14712 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.334 ms/conn | 1946448 B/op | 15229 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.334 ms/conn | 1370720 B/op | 10139 allocs/op |
| Direct external PSK handshake | 5 | 0.814 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.182 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.3 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 13.72 ms/conn | 3314408 B/op | 22480 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 3492320 B/op | 19887 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 238008 B/op | 987 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.549 ms/conn | 1555648 B/op | 9719 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.67 ms/conn | 1575808 B/op | 9859 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 7.164 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.686 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.172 ms/conn | 1169544 B/op | 1203 allocs/op |
| Full mTLS handshake | 5 | 8.613 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.619 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.898 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 0.969 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 8.247 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.259 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.32 ms/conn | 1172552 B/op | 1208 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190848 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.062 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.518 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 10.51 ms/conn | 64416 B/op | 56 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 42.18 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 56.79 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 745.1 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.28 us/op | 1796.75 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 2.399 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 236.8 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 196 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.333 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.792 us/op | 1080.1 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 349.9 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 93.93 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 367.7 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 68.42 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 23.26 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 48.88 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 22.92 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 568.7 ns/op | 2110.1 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 521.3 ns/op | 2302.14 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 529.9 ns/op | 2264.67 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 23.909 us/op | 2741.03 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 471.6 ns/op | 2544.47 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 28.59 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.579 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 3.465 us/op | 346.3 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.306 us/op | 918.93 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 10.95 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.836 us/op | 312.84 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.879 us/op | 638.5 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 9.945 us/op | 120.67 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 3.429 us/op | 349.99 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 3.532 us/op | 339.77 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.585 us/op | 261.72 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.426 us/op | 841.35 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 4.037 us/op | 297.25 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.431 us/op | 838.62 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.495 us/op | 802.85 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.899 us/op | 631.9 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.936 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.616 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.547 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.536 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.534 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.511 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.407 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.407 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.467 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.515 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.454 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 10.831 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.3 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 16.369 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 722.2 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.908 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 9.148 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 744.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 748.7 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.679 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.18 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 9.142 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.692 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.715 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.252 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.59 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.756 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.882 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.79 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.854 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.708 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.724 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 362.1 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 721.9 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 121.8 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 87.75 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 316.4 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 264.7 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 370.2 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 563.8 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 52.7 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 514.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 89.53 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 45.62 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 75.29 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 713.5 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 90.16 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 70.51 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 76.28 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 581.1 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 417.4 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.38 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 66.99 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 94.01 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 65.25 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 35.23 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 281 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 163.4 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 84.5 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.067 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 784.2 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 791.2 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 11.99 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 63.92 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.182 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 6.214 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
