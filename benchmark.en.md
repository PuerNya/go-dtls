# Automated benchmark results

[简体中文](README.md) | [English](benchmark.en.md) | [Русский](benchmark.ru.md)

- Commit: `30146f0ccfc956cc3f7d27a13ca1eda152f15c02`
- Generated: `2026-10-08T12:51:21Z`
- Go: `go version go1.27.1 linux/amd64`
- Platform: `linux/amd64, Intel(R) Xeon(R) Platinum 8370C CPU @ 2.80GHz`
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
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 536.731 us/op | 98185 B/op | 710 allocs/op |
| Full mTLS handshake | 5 | 812.849 us/op | 112583 B/op | 895 allocs/op |
| mTLS session resumption handshake | 5 | 412.561 us/op | 119639 B/op | 823 allocs/op |
| Multi-certificate mTLS selection by CA and OID filters | 5 | 1.056 ms/op | 120444 B/op | 1069 allocs/op |
| Multi-certificate post-handshake authentication selection | 5 | 1.42 ms/op | 139532 B/op | 1363 allocs/op |
| Full handshake + 4 acknowledged session tickets | 5 | 633.43 us/op | 117760 B/op | 934 allocs/op |
| Full mTLS handshake + session ticket / GREASE disabled | 5 | 875.002 us/op | 121897 B/op | 967 allocs/op |
| Full mTLS handshake + session ticket / GREASE enabled | 5 | 875.273 us/op | 121897 B/op | 967 allocs/op |
| Direct external PSK handshake | 5 | 347.554 us/op | 101630 B/op | 744 allocs/op |
| Full server-certificate handshake / uncompressed certificate | 5 | 1.004 ms/op | 130290 B/op | 993 allocs/op |
| zlib-compressed server-certificate handshake | 5 | 994.883 us/op | 122507 B/op | 971 allocs/op |
| Full mTLS handshake / uncompressed certificates | 5 | 1.675 ms/op | 170216 B/op | 1429 allocs/op |
| zlib-compressed mTLS handshake | 5 | 1.684 ms/op | 157779 B/op | 1388 allocs/op |
| ECH handshake / direct (no HRR) | 5 | 928.09 us/op | 147691 B/op | 1210 allocs/op |
| ECH handshake / via HRR | 5 | 931.618 us/op | 150467 B/op | 1231 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 842.436 us/op | 146164 B/op | 743 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 815.192 us/op | 149508 B/op | 774 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 2.246 ms/op | 174949 B/op | 792 allocs/op |

<a id="section-real-udp-interoperability"></a>
## Real UDP interoperability

<a id="real-udp-go-dtls-client-go-dtls-server"></a>
### go-dtls client -> go-dtls server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 1.832 ms/conn | 2493840 B/op | 21254 allocs/op |
| 1-RTT application-data round trip | 5 | 1.845 ms/conn | 2502608 B/op | 21576 allocs/op |
| Full mTLS handshake | 5 | 3.266 ms/conn | 3179920 B/op | 29404 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 4.553 ms/conn | 3678912 B/op | 31426 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 1.863 ms/conn | 2982960 B/op | 25185 allocs/op |
| Direct external PSK handshake | 5 | 0.4213 ms/conn | 1733984 B/op | 14633 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 1.844 ms/conn | 2513856 B/op | 22373 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 1.944 ms/conn | 2671248 B/op | 22888 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 4.906 ms/conn | 3978536 B/op | 39548 allocs/op |
| Session resumption handshake | 5 | 0.4825 ms/conn | 4457888 B/op | 37893 allocs/op |
| mTLS session resumption handshake | 5 | 0.6133 ms/conn | 6481712 B/op | 49695 allocs/op |
| 0-RTT + 1-RTT application-data round trip | 5 | 0.5302 ms/conn | 288328 B/op | 1895 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.135 ms/conn | 3365352 B/op | 21578 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 2.131 ms/conn | 3407272 B/op | 21918 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 3.466 ms/conn | 3902952 B/op | 22278 allocs/op |

<a id="real-udp-go-dtls-client-wolfssl-server"></a>
### go-dtls client -> wolfSSL server

Median time is measured by the go-dtls client; `ms/conn` means one complete connection workload.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 4.944 ms/conn | 1191184 B/op | 10762 allocs/op |
| 1-RTT application-data round trip | 5 | 5.013 ms/conn | 1211184 B/op | 11022 allocs/op |
| Full mTLS handshake | 5 | 6.696 ms/conn | 1347184 B/op | 11602 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 6.658 ms/conn | 1347184 B/op | 11602 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.867 ms/conn | 1429584 B/op | 12322 allocs/op |
| Direct external PSK handshake | 5 | 0.9189 ms/conn | 814064 B/op | 6882 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.908 ms/conn | 1198384 B/op | 11322 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 5.22 ms/conn | 1447184 B/op | 12362 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 6.713 ms/conn | 1403024 B/op | 12387 allocs/op |
| Session resumption handshake | 5 | 0.9636 ms/conn | 2150200 B/op | 18384 allocs/op |
| mTLS session resumption handshake | 5 | 0.9575 ms/conn | 2311656 B/op | 19264 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting go-dtls 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 5.224 ms/conn | 1862840 B/op | 10964 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.4 ms/conn | 1875320 B/op | 11084 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | - | Skipped: wolfSSL server receive buffers cannot hold the P384 ClientHello at MTU 4096; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |

<a id="real-udp-wolfssl-client-go-dtls-server"></a>
### wolfSSL client -> go-dtls server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 2.18 ms/conn | 1293024 B/op | 9299 allocs/op |
| 1-RTT application-data round trip | 5 | 4.595 ms/conn | 2447768 B/op | 10805 allocs/op |
| Full mTLS handshake | 5 | 5.434 ms/conn | 1736128 B/op | 14855 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 5.444 ms/conn | 1938944 B/op | 15369 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 2.134 ms/conn | 1359104 B/op | 10219 allocs/op |
| Direct external PSK handshake | 5 | 0.761 ms/conn | 1017296 B/op | 7499 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 4.641 ms/conn | 2461624 B/op | 11366 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 4.672 ms/conn | 2623544 B/op | 12445 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 11.4 ms/conn | 3365800 B/op | 22634 allocs/op |
| Session resumption handshake | 5 | 1010 ms/pair | 3576336 B/op | 20067 allocs/op |
| mTLS session resumption handshake | - | Skipped: wolfSSL client refuses to fragment the first ClientHello containing the mTLS ticket; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| 0-RTT + 1-RTT application-data round trip | 5 | 1009 ms/pair | 241472 B/op | 997 allocs/op |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 2.404 ms/conn | 1595168 B/op | 9799 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 3.454 ms/conn | 1615488 B/op | 9939 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 6.43 ms/conn | 1774240 B/op | 10119 allocs/op |

<a id="real-udp-wolfssl-client-wolfssl-server"></a>
### wolfSSL client -> wolfSSL server

Median time is measured by the wolfSSL client; `ms/conn` means one connection, while `ms/pair` means a full-plus-resumed connection pair and includes the client's built-in wait.

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Certificate-authenticated full handshake / AES-128-GCM | 5 | 3.967 ms/conn | 63640 B/op | 55 allocs/op |
| 1-RTT application-data round trip | 5 | 7.03 ms/conn | 1169472 B/op | 1201 allocs/op |
| Full mTLS handshake | 5 | 7.328 ms/conn | 63608 B/op | 55 allocs/op |
| GREASE compatibility / full mTLS handshake + session ticket | 5 | 7.649 ms/conn | 63608 B/op | 55 allocs/op |
| Certificate-authenticated full handshake / AES-128-CCM | 5 | 4.057 ms/conn | 63640 B/op | 55 allocs/op |
| Direct external PSK handshake | 5 | 1.129 ms/conn | 63688 B/op | 55 allocs/op |
| CID + 1-RTT application-data round trip | 5 | 7.079 ms/conn | 1183168 B/op | 1202 allocs/op |
| KeyUpdate + 1-RTT application-data round trip | 5 | 6.889 ms/conn | 1169792 B/op | 1201 allocs/op |
| PHA + 1-RTT application-data round trip | 5 | 10.13 ms/conn | 1172480 B/op | 1206 allocs/op |
| Session resumption handshake | 5 | 1013 ms/pair | 1190920 B/op | 1204 allocs/op |
| mTLS session resumption handshake | 5 | 1016 ms/pair | 1190208 B/op | 1202 allocs/op |
| 0-RTT + 1-RTT application-data round trip | - | Skipped: this fixture enables cookie HRR, which requires rejecting wolfSSL client 0-RTT; last verified against wolfSSL commit 7a8aae3e40138d19c640ae5bc0bc4e8f2998c22d | - | - |
| Post-quantum hybrid key exchange / X25519MLKEM768 | 5 | 60.16 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP256r1MLKEM768 | 5 | 6.504 ms/conn | 64392 B/op | 55 allocs/op |
| Post-quantum hybrid key exchange / SecP384r1MLKEM1024 | 5 | 65.42 ms/conn | 64392 B/op | 55 allocs/op |

<a id="section-record-layer-and-reliability"></a>
## Record layer and reliability

| Benchmark | Samples | Median time | Throughput | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: | :---: |
| Plain ACK build / Empty | 5 | 37.97 ns/op | - | 16 B/op | 1 allocs/op |
| Plain ACK build / Single | 5 | 48.91 ns/op | - | 32 B/op | 1 allocs/op |
| Plain ACK build / Sorted64 | 5 | 600.3 ns/op | - | 1152 B/op | 1 allocs/op |
| Build Plain Flight | 5 | 2.432 us/op | 1684.15 MB/s | 5008 B/op | 9 allocs/op |
| Protected ACK build / Reversed64 | 5 | 1.826 us/op | - | 2200 B/op | 3 allocs/op |
| Protected ACK build / Single | 5 | 211.6 ns/op | - | 72 B/op | 2 allocs/op |
| Protected ACK build / Single Reuse | 5 | 166.1 ns/op | - | 48 B/op | 1 allocs/op |
| Protected ACK build / Sorted64 | 5 | 1.003 us/op | - | 1176 B/op | 2 allocs/op |
| Build Protected Flight | 5 | 3.713 us/op | 1103.27 MB/s | 5584 B/op | 6 allocs/op |
| Combine Flights | 5 | 383.2 ns/op | - | 592 B/op | 4 allocs/op |
| Flight First Refresh | 5 | 105.3 ns/op | - | 0 B/op | 0 allocs/op |
| Flight Initial History Batch | 5 | 381.5 ns/op | - | 480 B/op | 1 allocs/op |
| Flight Wire Window / Retransmit | 5 | 27.97 ns/op | - | 0 B/op | 0 allocs/op |
| Inbox / Single fragment | 5 | 605.8 ns/op | 1980.78 MB/s | 1312 B/op | 2 allocs/op |
| Inbox / Fragment batch | 5 | 569.7 ns/op | 2106.41 MB/s | 1280 B/op | 1 allocs/op |
| Inbox / Fragment reuse | 5 | 579.5 ns/op | 2070.58 MB/s | 1280 B/op | 1 allocs/op |
| Handshake Reassembly | 5 | 28.401 us/op | 2307.55 MB/s | 73856 B/op | 3 allocs/op |
| Handshake Reassembly Single Fragment | 5 | 552.8 ns/op | 2170.89 MB/s | 1280 B/op | 1 allocs/op |
| Parse ACK / Owned | 5 | 27.27 ns/op | - | 16 B/op | 1 allocs/op |
| Parse ACK / Reuse Single | 5 | 4.659 ns/op | - | 0 B/op | 0 allocs/op |
| Protected Record CID / Round Trip | 5 | 2.799 us/op | 428.74 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record CID / Seal | 5 | 1.02 us/op | 1176.06 MB/s | 1280 B/op | 1 allocs/op |
| Reject unauthenticated record | 5 | 12.31 ns/op | - | 0 B/op | 0 allocs/op |
| Record round trip | 5 | 3.439 us/op | 348.92 MB/s | 3840 B/op | 3 allocs/op |
| Protected Record Round Trip In Place | 5 | 1.599 us/op | 750.49 MB/s | 1280 B/op | 1 allocs/op |
| Record round trip / AES-128-CCM | 5 | 8.001 us/op | 149.97 MB/s | 6240 B/op | 12 allocs/op |
| Record round trip / AES-128-GCM | 5 | 2.782 us/op | 431.31 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / AES-256-GCM | 5 | 2.91 us/op | 412.34 MB/s | 3840 B/op | 3 allocs/op |
| Record round trip / ChaCha20-Poly1305 | 5 | 4.057 us/op | 295.82 MB/s | 3840 B/op | 3 allocs/op |
| Record seal | 5 | 1.248 us/op | 961.69 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-128-CCM | 5 | 3.306 us/op | 362.98 MB/s | 1840 B/op | 5 allocs/op |
| Record seal / AES-128-GCM | 5 | 1.155 us/op | 1039.4 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / AES-256-GCM | 5 | 1.2 us/op | 1000.37 MB/s | 1280 B/op | 1 allocs/op |
| Record seal / ChaCha20-Poly1305 | 5 | 1.801 us/op | 666.15 MB/s | 1280 B/op | 1 allocs/op |

<a id="section-key-schedule-and-cryptography"></a>
## Key schedule and cryptography

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Calculate PSK Binder / AES-128-GCM | 5 | 3.022 us/op | 1952 B/op | 21 allocs/op |
| Calculate PSK Binder / AES-256-GCM | 5 | 6.936 us/op | 3248 B/op | 21 allocs/op |
| Derive Traffic Keys / AES-128-GCM | 5 | 1.683 us/op | 976 B/op | 9 allocs/op |
| Derive Traffic Keys / AES-256-GCM | 5 | 3.721 us/op | 1520 B/op | 9 allocs/op |
| Derive Traffic Keys Into / AES-128-GCM | 5 | 1.63 us/op | 928 B/op | 8 allocs/op |
| Derive Traffic Keys Into / AES-256-GCM | 5 | 3.656 us/op | 1440 B/op | 8 allocs/op |
| Empty Transcript Hash / AES-128-GCM | 5 | 1.154 ns/op | 0 B/op | 0 allocs/op |
| Empty Transcript Hash / AES-256-GCM | 5 | 1.153 ns/op | 0 B/op | 0 allocs/op |
| Finished Verify Data / AES-128-GCM | 5 | 1.513 us/op | 992 B/op | 11 allocs/op |
| Finished Verify Data / AES-256-GCM | 5 | 3.471 us/op | 1648 B/op | 11 allocs/op |
| Install Application Keys / AES-128-GCM | 5 | 6.904 us/op | 7488 B/op | 34 allocs/op |
| Install Application Keys / AES-256-GCM | 5 | 11.396 us/op | 8544 B/op | 34 allocs/op |
| Key Schedule Derivation / AES-128-GCM | 5 | 7.807 us/op | 5184 B/op | 48 allocs/op |
| Key Schedule Derivation / AES-256-GCM | 5 | 17.177 us/op | 8224 B/op | 48 allocs/op |
| Key derivation / AES-128-GCM / Early Traffic | 5 | 768.2 ns/op | 480 B/op | 5 allocs/op |
| Key derivation / AES-128-GCM / Exporter | 5 | 2.033 us/op | 1408 B/op | 15 allocs/op |
| Key derivation / AES-128-GCM / Exporter Zero | 5 | 18.54 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-128-GCM / Resumption PSK | 5 | 793 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-128-GCM / Traffic Update | 5 | 789.6 ns/op | 512 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Early Traffic | 5 | 1.734 us/op | 800 B/op | 5 allocs/op |
| Key derivation / AES-256-GCM / Exporter | 5 | 4.479 us/op | 2384 B/op | 15 allocs/op |
| Key derivation / AES-256-GCM / Exporter Zero | 5 | 18.54 ns/op | 0 B/op | 0 allocs/op |
| Key derivation / AES-256-GCM / Resumption PSK | 5 | 1.787 us/op | 848 B/op | 6 allocs/op |
| Key derivation / AES-256-GCM / Traffic Update | 5 | 1.78 us/op | 848 B/op | 6 allocs/op |
| New Record Cipher / AES-128-CCM | 5 | 2.538 us/op | 2520 B/op | 13 allocs/op |
| New Record Cipher / AES-128-GCM | 5 | 2.872 us/op | 3264 B/op | 13 allocs/op |
| New Record Cipher / AES-256-GCM | 5 | 4.917 us/op | 3776 B/op | 13 allocs/op |
| New Record Cipher / ChaCha20-Poly1305 | 5 | 2.003 us/op | 1528 B/op | 12 allocs/op |
| Receive KeyUpdate / AES-128-GCM | 5 | 4.162 us/op | 3776 B/op | 19 allocs/op |
| Receive KeyUpdate / AES-256-GCM | 5 | 7.347 us/op | 4624 B/op | 19 allocs/op |
| Send KeyUpdate / AES-128-GCM | 5 | 3.993 us/op | 3792 B/op | 19 allocs/op |
| Send KeyUpdate / AES-256-GCM | 5 | 6.981 us/op | 4624 B/op | 19 allocs/op |
| Transcript Clone / AES-128-GCM | 5 | 329.9 ns/op | 288 B/op | 4 allocs/op |
| Transcript Clone / AES-256-GCM | 5 | 633.5 ns/op | 496 B/op | 4 allocs/op |
| Transcript Sum / AES-128-GCM / Owned | 5 | 116.7 ns/op | 32 B/op | 1 allocs/op |
| Transcript Sum / AES-128-GCM / Reuse | 5 | 77.94 ns/op | 0 B/op | 0 allocs/op |
| Transcript Sum / AES-256-GCM / Owned | 5 | 315.6 ns/op | 48 B/op | 1 allocs/op |
| Transcript Sum / AES-256-GCM / Reuse | 5 | 262.6 ns/op | 0 B/op | 0 allocs/op |

<a id="section-wire-encoding-and-parsing"></a>
## Wire encoding and parsing

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Marshal Extensions | 5 | 329.7 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / Certificate | 5 | 626.9 ns/op | 1152 B/op | 1 allocs/op |
| Handshake marshal / Certificate Verify | 5 | 56.44 ns/op | 80 B/op | 1 allocs/op |
| Handshake marshal / Client Hello | 5 | 601.3 ns/op | 424 B/op | 8 allocs/op |
| Handshake marshal / Hello Retry Request | 5 | 101.1 ns/op | 128 B/op | 1 allocs/op |
| Handshake marshal / New Connection ID | 5 | 47.65 ns/op | 32 B/op | 1 allocs/op |
| Handshake marshal / New Session Ticket | 5 | 79 ns/op | 96 B/op | 1 allocs/op |
| Handshake marshal / Resumption Client Hello | 5 | 832.5 ns/op | 744 B/op | 9 allocs/op |
| Handshake marshal / Server Hello | 5 | 105.5 ns/op | 112 B/op | 1 allocs/op |
| Handshake marshal / Session Ticket State | 5 | 74.67 ns/op | 80 B/op | 1 allocs/op |
| Parse Extensions / Ordered View | 5 | 63.46 ns/op | 0 B/op | 0 allocs/op |
| Parse Extensions / Owned | 5 | 652.3 ns/op | 472 B/op | 8 allocs/op |
| Parse Extensions / View | 5 | 444.4 ns/op | 336 B/op | 2 allocs/op |
| Parse Handshake Fragment / Reuse Single | 5 | 10.66 ns/op | 0 B/op | 0 allocs/op |
| Parse Handshake Fragment / View | 5 | 58.63 ns/op | 48 B/op | 1 allocs/op |
| Key share parse / 1 key share / Owned | 5 | 96.56 ns/op | 64 B/op | 2 allocs/op |
| Key share parse / 1 key share / View | 5 | 61.19 ns/op | 32 B/op | 1 allocs/op |
| Key share parse / 1 key share / View Into | 5 | 27.05 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 4 key shares / Owned | 5 | 304.7 ns/op | 256 B/op | 5 allocs/op |
| Key share parse / 4 key shares / View | 5 | 168.2 ns/op | 128 B/op | 1 allocs/op |
| Key share parse / 4 key shares / View Into | 5 | 75.99 ns/op | 0 B/op | 0 allocs/op |
| Key share parse / 9 key shares / Owned | 5 | 1.149 us/op | 824 B/op | 14 allocs/op |
| Key share parse / 9 key shares / View | 5 | 825.9 ns/op | 536 B/op | 5 allocs/op |
| Key share parse / 9 key shares / View Into | 5 | 819.1 ns/op | 536 B/op | 5 allocs/op |
| Parse Plain Record / Reuse Single | 5 | 18.46 ns/op | 0 B/op | 0 allocs/op |
| Parse Plain Record / View | 5 | 64.52 ns/op | 48 B/op | 1 allocs/op |

<a id="section-certificate-compression"></a>
## Certificate compression

| Benchmark | Samples | Median time | Harness memory | Harness allocations |
| --- | :---: | :---: | :---: | :---: |
| Compress | 5 | 5.043 us/op | 176 B/op | 3 allocs/op |
| Decompress | 5 | 5.743 us/op | 4264 B/op | 6 allocs/op |

[Raw Go benchmark output](benchmark.txt)
