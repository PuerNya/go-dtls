# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `8010fee127101d213fb208cbfbff51779d6c6954`
- Generated: `2026-10-06T03:54:10Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, INTEL(R) XEON(R) PLATINUM 8573C`
- wolfSSL: `b874f7167fc592b331c79a33674bc51c2f0286b0 (Linux Release static)`

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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 394.977 us/op | 97817 B/op | 709 allocs/op |
| Full mTLS handshake | 5 | 607.17 us/op | 111975 B/op | 891 allocs/op |
| mTLS session resumption handshake | 5 | 322.684 us/op | 119186 B/op | 822 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 798.085 us/op | 119589 B/op | 1064 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.07 ms/op | 138634 B/op | 1358 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 472.564 us/op | 117039 B/op | 932 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 652.52 us/op | 121225 B/op | 963 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 652.658 us/op | 121225 B/op | 963 allocs/op |
| Direct external PSK handshake | 5 | 271.591 us/op | 101264 B/op | 743 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 711.875 us/op | 129632 B/op | 992 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 710.437 us/op | 121751 B/op | 970 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.196 ms/op | 169029 B/op | 1425 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.191 ms/op | 156479 B/op | 1384 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 703.72 us/op | 147133 B/op | 1209 allocs/op |
| ECH handshake / via HRR | 5 | 721.386 us/op | 149910 B/op | 1230 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 612.875 us/op | 145793 B/op | 742 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 595.158 us/op | 149138 B/op | 773 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 1.537 ms/op | 174579 B/op | 791 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.226 ms/conn | 2482048 B/op | 21235 allocs/op |
| 1-RTT application-data round trip | 5 | 1.254 ms/conn | 2490720 B/op | 21555 allocs/op |
| Full mTLS handshake | 5 | 2.173 ms/conn | 3163456 B/op | 29333 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.369 ms/conn | 3658944 B/op | 31346 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.235 ms/conn | 2971168 B/op | 25166 allocs/op |
| Direct external PSK handshake | 5 | 0.2951 ms/conn | 1724064 B/op | 14613 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.235 ms/conn | 2502016 B/op | 22353 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.337 ms/conn | 2659408 B/op | 22868 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 3.251 ms/conn | 3956952 B/op | 39455 allocs/op |
| Session resumption handshake | 5 | 0.3481 ms/conn | 4432304 B/op | 37841 allocs/op |
| mTLS session resumption handshake | 5 | 0.4032 ms/conn | 6449872 B/op | 49583 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.3905 ms/conn | 286952 B/op | 1893 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.482 ms/conn | 3353512 B/op | 21558 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 1.477 ms/conn | 3395432 B/op | 21898 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.425 ms/conn | 3891128 B/op | 22258 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.099 ms/conn | 1186064 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 3.144 ms/conn | 1204784 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 4.071 ms/conn | 1341744 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.044 ms/conn | 1341744 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 3.086 ms/conn | 1424464 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.624 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.074 ms/conn | 1193264 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.198 ms/conn | 1440784 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.008 ms/conn | 1396816 B/op | 12385 allocs/op |
| Session resumption handshake | 5 | 0.8169 ms/conn | 2124600 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 0.8138 ms/conn | 2285736 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 3.234 ms/conn | 1857800 B/op | 10965 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.985 ms/conn | 1870200 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Unsupported: wolfSSL server does not complete this DTLS 1.3 hybrid handshake; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.53 ms/conn | 1305808 B/op | 9220 allocs/op |
| 1-RTT application-data round trip | 5 | 3.476 ms/conn | 2406104 B/op | 10725 allocs/op |
| Full mTLS handshake | 5 | 3.938 ms/conn | 1742840 B/op | 14706 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 3.941 ms/conn | 1946448 B/op | 15229 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.419 ms/conn | 1370720 B/op | 10139 allocs/op |
| Direct external PSK handshake | 5 | 0.516 ms/conn | 973648 B/op | 7399 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 3.457 ms/conn | 2418936 B/op | 11286 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 3.522 ms/conn | 2581624 B/op | 12365 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 8.498 ms/conn | 3314640 B/op | 22483 allocs/op |
| Session resumption handshake | 5 | 1007 ms/pair | 3492336 B/op | 19887 allocs/op |
| mTLS session resumption handshake | - | Unsupported: wolfSSL client cannot parse the go-dtls mTLS session ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1007 ms/pair | 238256 B/op | 990 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 1.596 ms/conn | 1555648 B/op | 9719 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.471 ms/conn | 1575808 B/op | 9859 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 4.727 ms/conn | 1734688 B/op | 10039 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.956 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 5.197 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 5.22 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.279 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.849 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 0.684 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 5.24 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.352 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 7.796 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1009 ms/pair | 1190848 B/op | 1202 allocs/op |
| mTLS session resumption handshake | 5 | 1011 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Unsupported: wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.918 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 4.652 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.703 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 29.95 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 37.73 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 524 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 1.934 us/op | 2117.67 MB/s | 5040 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.423 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 157.1 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 131.7 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 788.7 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.019 us/op | 1356.74 MB/s | 5616 B/op | 6 allocs/op |
| Combine Flights | 5 | 299.4 ns/op | - | 624 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 105 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 301.8 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Pending Indices / Allocated | 5 | 62.38 ns/op | - | 80 B/op | 1 allocs/op |
| Flight Pending Indices / Reuse Window | 5 | 24.97 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Pending | 5 | 41.81 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Wire Window / Retransmit | 5 | 23.35 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 459.7 ns/op | 2610.42 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 433.8 ns/op | 2766.3 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 436 ns/op | 2752.08 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 22.06 us/op | 2970.83 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 405.4 ns/op | 2960.25 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 21.43 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 3.494 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.038 us/op | 588.73 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 767.8 ns/op | 1562.9 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 9.2 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.232 us/op | 537.66 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.07 us/op | 1121.42 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 5.91 us/op | 203.03 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.026 us/op | 592.29 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.162 us/op | 555.07 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.086 us/op | 388.82 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 868 ns/op | 1382.56 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 2.466 us/op | 486.69 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 782.4 ns/op | 1533.82 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 836.8 ns/op | 1434.09 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.307 us/op | 918.14 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 2.256 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 5.501 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.261 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 2.959 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.229 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 2.926 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 0.9168 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.139 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 2.783 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 5.268 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 8.756 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 5.788 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 13.935 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 559.6 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.499 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 16.87 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 588.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 589.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.379 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 3.555 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 16.86 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.423 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.421 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 1.92 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.204 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 3.928 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.509 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.142 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 5.689 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.002 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 5.554 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 243.8 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 485.4 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 97.96 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 70.34 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 262.1 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 225.1 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 277.5 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 475.7 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 43.11 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 422.2 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 72.36 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 36.04 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 59.18 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 601.7 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 73.34 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 55.9 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 52.77 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 512 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 327.5 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.37 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 45.38 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 73.87 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 45.51 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 22.68 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 240.5 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 128.7 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 62.62 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 905.3 ns/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 647.5 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 647.6 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 9.252 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 45.18 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.173 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 5.129 us/op | 4248 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
