# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `f3dfd50e87309e9051b2b839006de3f174b1e328`
- Generated: `2026-10-08T03:38:02Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `799bd17483efd76f30bd0f17cb080a0fc0ea92da (Linux Release static)`

178 results, grouped by workload and ordered by feature, then benchmark name. Values are medians of the samples emitted by the final benchmark run.

Workload-specific connection metrics are preferred over the Go harness time. Memory and allocations remain per Go benchmark operation. Exact raw output remains in the workflow artifact.

## Quick navigation

- [Connection lifecycle (18)](#section-connection-lifecycle)
- [Real UDP interoperability (60)](#section-real-udp-interoperability)
  - [go-dtls client -> go-dtls server (15)](#real-udp-go-dtls-client-go-dtls-server)
  - [go-dtls client -> wolfSSL server (15)](#real-udp-go-dtls-client-wolfssl-server)
  - [wolfSSL client -> go-dtls server (15)](#real-udp-wolfssl-client-go-dtls-server)
  - [wolfSSL client -> wolfSSL server (15)](#real-udp-wolfssl-client-wolfssl-server)
- [Record layer and reliability (34)](#section-record-layer-and-reliability)
- [Key schedule and cryptography (38)](#section-key-schedule-and-cryptography)
- [Wire encoding and parsing (26)](#section-wire-encoding-and-parsing)
- [Certificate compression (2)](#section-certificate-compression)

<a id="section-connection-lifecycle"></a>
## Connection lifecycle

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 560.167 us/op | 97993 B/op | 710 allocs/op |
| Full mTLS handshake | 5 | 843.364 us/op | 112215 B/op | 894 allocs/op |
| mTLS session resumption handshake | 5 | 416.825 us/op | 119419 B/op | 823 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.085 ms/op | 120038 B/op | 1068 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.455 ms/op | 138942 B/op | 1362 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 647.931 us/op | 117567 B/op | 934 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 891.881 us/op | 121529 B/op | 966 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 895.227 us/op | 121529 B/op | 966 allocs/op |
| Direct external PSK handshake | 5 | 355.699 us/op | 101438 B/op | 744 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.014 ms/op | 130098 B/op | 993 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 1.007 ms/op | 122315 B/op | 971 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.683 ms/op | 169843 B/op | 1428 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.687 ms/op | 157393 B/op | 1387 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 957.264 us/op | 147499 B/op | 1210 allocs/op |
| ECH handshake / via HRR | 5 | 958.16 us/op | 150275 B/op | 1231 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 861.323 us/op | 145972 B/op | 743 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 814.283 us/op | 149316 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.08 ms/op | 174758 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.931 ms/conn | 2487488 B/op | 21255 allocs/op |
| 1-RTT application-data round trip | 5 | 1.965 ms/conn | 2496160 B/op | 21575 allocs/op |
| Full mTLS handshake | 5 | 3.351 ms/conn | 3170176 B/op | 29384 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.593 ms/conn | 3666696 B/op | 31403 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.94 ms/conn | 2976560 B/op | 25185 allocs/op |
| Direct external PSK handshake | 5 | 0.5253 ms/conn | 1727584 B/op | 14633 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.93 ms/conn | 2507456 B/op | 22373 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 2.073 ms/conn | 2664848 B/op | 22888 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.017 ms/conn | 3964168 B/op | 39514 allocs/op |
| Session resumption handshake | 5 | 0.5858 ms/conn | 4443080 B/op | 37899 allocs/op |
| mTLS session resumption handshake | 5 | 0.6384 ms/conn | 6461928 B/op | 49652 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.6372 ms/conn | 287904 B/op | 1898 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.295 ms/conn | 3358952 B/op | 21578 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.271 ms/conn | 3400872 B/op | 21918 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.493 ms/conn | 3896552 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.503 ms/conn | 1187344 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 4.578 ms/conn | 1207344 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 5.996 ms/conn | 1343344 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.963 ms/conn | 1343344 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.477 ms/conn | 1425744 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.7904 ms/conn | 810224 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.496 ms/conn | 1194544 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.667 ms/conn | 1443344 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 5.904 ms/conn | 1398992 B/op | 12384 allocs/op |
| Session resumption handshake | 5 | 0.9694 ms/conn | 2141240 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 0.9356 ms/conn | 2302696 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.747 ms/conn | 1859000 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 5.76 ms/conn | 1871480 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.305 ms/conn | 1290720 B/op | 9299 allocs/op |
| 1-RTT application-data round trip | 5 | 4.962 ms/conn | 2445144 B/op | 10805 allocs/op |
| Full mTLS handshake | 5 | 5.976 ms/conn | 1729480 B/op | 14832 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.02 ms/conn | 1932896 B/op | 15355 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.272 ms/conn | 1356320 B/op | 10219 allocs/op |
| Direct external PSK handshake | 5 | 0.923 ms/conn | 1014288 B/op | 7499 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.958 ms/conn | 2459256 B/op | 11366 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.024 ms/conn | 2620664 B/op | 12445 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 12.52 ms/conn | 3354064 B/op | 22600 allocs/op |
| Session resumption handshake | 5 | 1008 ms/pair | 3571568 B/op | 20067 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1008 ms/pair | 241376 B/op | 997 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.543 ms/conn | 1592768 B/op | 9799 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.587 ms/conn | 1612768 B/op | 9939 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.955 ms/conn | 1771424 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.896 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 8.082 ms/conn | 1169496 B/op | 1202 allocs/op |
| Full mTLS handshake | 5 | 8.582 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 8.528 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.718 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.057 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.986 ms/conn | 1183240 B/op | 1204 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 8.024 ms/conn | 1169408 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.37 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1012 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 4.536 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.899 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 9.797 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 37.52 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 46.1 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 552.6 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.158 us/op | 1898 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.73 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 208.4 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 174.3 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 975.5 ns/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.681 us/op | 1112.77 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 358.8 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 92.05 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 356.5 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 19.5 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 518.7 ns/op | 2313.34 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 467.5 ns/op | 2566.6 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 491.7 ns/op | 2440.58 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 23.701 us/op | 2765.12 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 430 ns/op | 2790.91 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 24.64 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.366 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.591 us/op | 463.19 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 976.8 ns/op | 1228.5 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 15.88 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 2.911 us/op | 412.28 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.519 us/op | 789.82 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 7.6 us/op | 157.9 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.567 us/op | 467.46 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.657 us/op | 451.67 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 3.382 us/op | 354.77 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.059 us/op | 1132.74 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.264 us/op | 367.69 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.042 us/op | 1151.12 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.113 us/op | 1077.75 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.45 us/op | 827.31 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 3.022 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.546 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.496 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.343 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.447 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.401 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.247 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.249 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.407 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.289 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.162 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 9.941 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.104 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 15.807 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 731.3 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 1.812 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 8.978 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 731.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 750.4 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.595 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.139 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 8.957 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.618 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.649 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.27 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.606 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.548 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 1.872 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 3.73 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 6.552 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.579 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.4 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 285.4 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 584.2 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 103.1 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 79.19 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 277.1 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 234.3 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 350.8 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 570.4 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 51.37 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 578.4 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 88.92 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 43.75 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 72.93 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 789.8 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 89.83 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 69.5 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 64.52 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 613.8 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 414.6 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 13.42 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 55.43 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 91.67 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 61.13 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 29.98 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 294.3 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 155.8 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 78.33 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.223 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 817.1 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 802.7 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 23.06 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 58.87 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 4.922 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 5.926 us/op | 4264 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
