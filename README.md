# Automated benchmark results

- Commit: `857762cdec131c22fefbeb378f26a930c83b9160`
- Generated: `2026-07-31T10:44:35Z`
- Go: `go version go1.26.5 linux/amd64`
- Platform: `linux/amd64, AMD EPYC 7763 64-Core Processor`
- wolfSSL: `41b7a0209abbddc579d3d861f36c0f574ae7e907 (Linux Release static)`

Values are medians of the samples emitted by the final benchmark run. Units are reported by Go's benchmark harness.

| Benchmark | Samples | Median metrics |
| --- | ---: | --- |
| `BenchmarkConnectionHandshakeLifecycle` | 5 | 575329 ns/op<br>99289 B/op<br>759 allocs/op |
| `BenchmarkECHHandshakeLifecycle/Direct` | 5 | 943062 ns/op<br>148598 B/op<br>1259 allocs/op |
| `BenchmarkECHHandshakeLifecycle/HRR` | 5 | 948659 ns/op<br>151374 B/op<br>1280 allocs/op |
| `BenchmarkExternalPSKHandshakeLifecycle` | 5 | 337831 ns/op<br>98265 B/op<br>724 allocs/op |
| `BenchmarkMutualTLSHandshakeLifecycle/Full` | 5 | 831169 ns/op<br>115950 B/op<br>973 allocs/op |
| `BenchmarkMutualTLSHandshakeLifecycle/Resumed` | 5 | 412221 ns/op<br>116100 B/op<br>804 allocs/op |
| `BenchmarkProtectedRecordSeal` | 5 | 1103 ns/op<br>1088.4 MB/s<br>1281 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordRoundTrip` | 5 | 3116 ns/op<br>385.11 MB/s<br>3840 B/op<br>3 allocs/op |
| `BenchmarkProtectedRecordRoundTripInPlace` | 5 | 1644 ns/op<br>729.73 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordReceiveErrorUnauthenticated` | 5 | 10.95 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkProtectedRecordSealSuites/AES128GCM` | 5 | 1154 ns/op<br>1040.25 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordSealSuites/AES256GCM` | 5 | 1229 ns/op<br>976.8 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordSealSuites/ChaCha20Poly1305` | 5 | 1572 ns/op<br>763.44 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordSealSuites/AES128CCM` | 5 | 3620 ns/op<br>331.53 MB/s<br>1840 B/op<br>5 allocs/op |
| `BenchmarkProtectedRecordRoundTripSuites/AES128GCM` | 5 | 3037 ns/op<br>395.17 MB/s<br>3840 B/op<br>3 allocs/op |
| `BenchmarkProtectedRecordRoundTripSuites/AES256GCM` | 5 | 2974 ns/op<br>403.48 MB/s<br>3840 B/op<br>3 allocs/op |
| `BenchmarkProtectedRecordRoundTripSuites/ChaCha20Poly1305` | 5 | 3699 ns/op<br>324.42 MB/s<br>3840 B/op<br>3 allocs/op |
| `BenchmarkProtectedRecordRoundTripSuites/AES128CCM` | 5 | 8121 ns/op<br>147.76 MB/s<br>6240 B/op<br>12 allocs/op |
| `BenchmarkProtectedRecordCID/Seal` | 5 | 1074 ns/op<br>1116.97 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkProtectedRecordCID/RoundTrip` | 5 | 2864 ns/op<br>419.04 MB/s<br>3840 B/op<br>3 allocs/op |
| `BenchmarkTranscriptClone/1301` | 5 | 333.5 ns/op<br>288 B/op<br>4 allocs/op |
| `BenchmarkTranscriptClone/1302` | 5 | 619.5 ns/op<br>496 B/op<br>4 allocs/op |
| `BenchmarkTranscriptSum/1301/Owned` | 5 | 114.1 ns/op<br>32 B/op<br>1 allocs/op |
| `BenchmarkTranscriptSum/1301/Reuse` | 5 | 79.43 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkTranscriptSum/1302/Owned` | 5 | 283 ns/op<br>48 B/op<br>1 allocs/op |
| `BenchmarkTranscriptSum/1302/Reuse` | 5 | 239.1 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkBuildProtectedACKRecords/Single` | 5 | 228.4 ns/op<br>72 B/op<br>2 allocs/op |
| `BenchmarkBuildProtectedACKRecords/Sorted64` | 5 | 1036 ns/op<br>1176 B/op<br>2 allocs/op |
| `BenchmarkBuildProtectedACKRecords/Reversed64` | 5 | 1823 ns/op<br>2200 B/op<br>3 allocs/op |
| `BenchmarkBuildProtectedACKRecords/SingleReuse` | 5 | 180.4 ns/op<br>48 B/op<br>1 allocs/op |
| `BenchmarkBuildPlainACKRecords/Empty` | 5 | 39.7 ns/op<br>16 B/op<br>1 allocs/op |
| `BenchmarkBuildPlainACKRecords/Single` | 5 | 51.63 ns/op<br>32 B/op<br>1 allocs/op |
| `BenchmarkBuildPlainACKRecords/Sorted64` | 5 | 566.7 ns/op<br>1152 B/op<br>1 allocs/op |
| `BenchmarkParseACK/Owned` | 5 | 29.18 ns/op<br>16 B/op<br>1 allocs/op |
| `BenchmarkParseACK/ReuseSingle` | 5 | 4.372 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkParseHandshakeFragment/View` | 5 | 63.56 ns/op<br>48 B/op<br>1 allocs/op |
| `BenchmarkParseHandshakeFragment/ReuseSingle` | 5 | 13.4 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkParsePlainRecord/View` | 5 | 63.03 ns/op<br>48 B/op<br>1 allocs/op |
| `BenchmarkParsePlainRecord/ReuseSingle` | 5 | 12.79 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFlightPendingIndices/Allocated` | 5 | 81.66 ns/op<br>80 B/op<br>1 allocs/op |
| `BenchmarkFlightPendingIndices/ReuseWindow` | 5 | 42.22 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFlightWireWindow/Pending` | 5 | 46.38 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFlightWireWindow/Retransmit` | 5 | 40.48 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFlightFirstRefresh` | 5 | 112.7 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFlightInitialHistoryBatch` | 5 | 361.2 ns/op<br>480 B/op<br>1 allocs/op |
| `BenchmarkBuildProtectedFlight` | 5 | 3758 ns/op<br>1090.02 MB/s<br>5616 B/op<br>6 allocs/op |
| `BenchmarkBuildPlainFlight` | 5 | 2199 ns/op<br>1862.65 MB/s<br>5040 B/op<br>9 allocs/op |
| `BenchmarkCombineFlights` | 5 | 373.4 ns/op<br>624 B/op<br>4 allocs/op |
| `BenchmarkParseExtensions/Owned` | 5 | 667.2 ns/op<br>472 B/op<br>8 allocs/op |
| `BenchmarkParseExtensions/View` | 5 | 420.1 ns/op<br>336 B/op<br>2 allocs/op |
| `BenchmarkParseExtensions/OrderedView` | 5 | 63.59 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkParseKeyShares/Client1/Owned` | 5 | 101.7 ns/op<br>64 B/op<br>2 allocs/op |
| `BenchmarkParseKeyShares/Client1/View` | 5 | 67.25 ns/op<br>32 B/op<br>1 allocs/op |
| `BenchmarkParseKeyShares/Client1/ViewInto` | 5 | 28.65 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkParseKeyShares/Client4/Owned` | 5 | 301.3 ns/op<br>256 B/op<br>5 allocs/op |
| `BenchmarkParseKeyShares/Client4/View` | 5 | 153.3 ns/op<br>128 B/op<br>1 allocs/op |
| `BenchmarkParseKeyShares/Client4/ViewInto` | 5 | 75.37 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkParseKeyShares/Client9/Owned` | 5 | 1163 ns/op<br>824 B/op<br>14 allocs/op |
| `BenchmarkParseKeyShares/Client9/View` | 5 | 819.3 ns/op<br>536 B/op<br>5 allocs/op |
| `BenchmarkParseKeyShares/Client9/ViewInto` | 5 | 809.2 ns/op<br>536 B/op<br>5 allocs/op |
| `BenchmarkDeriveTrafficKeys/1301` | 5 | 1506 ns/op<br>976 B/op<br>9 allocs/op |
| `BenchmarkDeriveTrafficKeys/1302` | 5 | 3403 ns/op<br>1520 B/op<br>9 allocs/op |
| `BenchmarkDeriveTrafficKeysInto/1301` | 5 | 1462 ns/op<br>928 B/op<br>8 allocs/op |
| `BenchmarkDeriveTrafficKeysInto/1302` | 5 | 3325 ns/op<br>1440 B/op<br>8 allocs/op |
| `BenchmarkNewRecordCipher/AES128GCM` | 5 | 2532 ns/op<br>3264 B/op<br>13 allocs/op |
| `BenchmarkNewRecordCipher/AES256GCM` | 5 | 4453 ns/op<br>3776 B/op<br>13 allocs/op |
| `BenchmarkNewRecordCipher/ChaCha20Poly1305` | 5 | 1810 ns/op<br>1528 B/op<br>12 allocs/op |
| `BenchmarkNewRecordCipher/AES128CCM` | 5 | 2236 ns/op<br>2520 B/op<br>13 allocs/op |
| `BenchmarkKeyScheduleDerivation/1301` | 5 | 7037 ns/op<br>5184 B/op<br>48 allocs/op |
| `BenchmarkKeyScheduleDerivation/1302` | 5 | 15683 ns/op<br>8224 B/op<br>48 allocs/op |
| `BenchmarkEmptyTranscriptHash/1301` | 5 | 1.249 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkEmptyTranscriptHash/1302` | 5 | 1.247 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkFinishedVerifyData/1301` | 5 | 1386 ns/op<br>992 B/op<br>11 allocs/op |
| `BenchmarkFinishedVerifyData/1302` | 5 | 3128 ns/op<br>1648 B/op<br>11 allocs/op |
| `BenchmarkCalculatePSKBinder/1301` | 5 | 2765 ns/op<br>1952 B/op<br>21 allocs/op |
| `BenchmarkCalculatePSKBinder/1302` | 5 | 6243 ns/op<br>3248 B/op<br>21 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1301/EarlyTraffic` | 5 | 695.7 ns/op<br>480 B/op<br>5 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1301/TrafficUpdate` | 5 | 722.8 ns/op<br>512 B/op<br>6 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1301/ResumptionPSK` | 5 | 720.4 ns/op<br>512 B/op<br>6 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1301/Exporter` | 5 | 1789 ns/op<br>1408 B/op<br>15 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1301/ExporterZero` | 5 | 8.956 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1302/EarlyTraffic` | 5 | 1558 ns/op<br>800 B/op<br>5 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1302/TrafficUpdate` | 5 | 1596 ns/op<br>848 B/op<br>6 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1302/ResumptionPSK` | 5 | 1601 ns/op<br>848 B/op<br>6 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1302/Exporter` | 5 | 4075 ns/op<br>2384 B/op<br>15 allocs/op |
| `BenchmarkKeyScheduleSideDerivations/1302/ExporterZero` | 5 | 8.97 ns/op<br>0 B/op<br>0 allocs/op |
| `BenchmarkReceivingTrafficKeyUpdate/1301` | 5 | 3672 ns/op<br>3776 B/op<br>19 allocs/op |
| `BenchmarkReceivingTrafficKeyUpdate/1302` | 5 | 6527 ns/op<br>4624 B/op<br>19 allocs/op |
| `BenchmarkSendingTrafficKeyUpdate/1301` | 5 | 3517 ns/op<br>3792 B/op<br>19 allocs/op |
| `BenchmarkSendingTrafficKeyUpdate/1302` | 5 | 6360 ns/op<br>4624 B/op<br>19 allocs/op |
| `BenchmarkInstallApplicationKeys/1301` | 5 | 6040 ns/op<br>7488 B/op<br>34 allocs/op |
| `BenchmarkInstallApplicationKeys/1302` | 5 | 9932 ns/op<br>8544 B/op<br>34 allocs/op |
| `BenchmarkMarshalExtensions` | 5 | 332.2 ns/op<br>128 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/ClientHello` | 5 | 524.3 ns/op<br>424 B/op<br>8 allocs/op |
| `BenchmarkMarshalHandshakeMessages/ResumptionClientHello` | 5 | 734.9 ns/op<br>744 B/op<br>9 allocs/op |
| `BenchmarkMarshalHandshakeMessages/ServerHello` | 5 | 88.24 ns/op<br>112 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/HelloRetryRequest` | 5 | 84.94 ns/op<br>128 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/CertificateVerify` | 5 | 52.25 ns/op<br>80 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/Certificate` | 5 | 559.2 ns/op<br>1152 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/NewSessionTicket` | 5 | 69.55 ns/op<br>96 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/SessionTicketState` | 5 | 71.38 ns/op<br>80 B/op<br>1 allocs/op |
| `BenchmarkMarshalHandshakeMessages/NewConnectionID` | 5 | 46.38 ns/op<br>32 B/op<br>1 allocs/op |
| `BenchmarkHandshakeReassembly` | 5 | 25993 ns/op<br>2521.28 MB/s<br>73856 B/op<br>3 allocs/op |
| `BenchmarkHandshakeReassemblySingleFragment` | 5 | 450.9 ns/op<br>2661.36 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkHandshakeInboxSequentialSingleFragment` | 5 | 548.6 ns/op<br>2187.22 MB/s<br>1312 B/op<br>2 allocs/op |
| `BenchmarkHandshakeInboxSequentialSingleFragmentReuse` | 5 | 498.2 ns/op<br>2408.87 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkHandshakeInboxSequentialSingleFragmentBatch` | 5 | 484.7 ns/op<br>2475.54 MB/s<br>1280 B/op<br>1 allocs/op |
| `BenchmarkCertificateCompression/Compress` | 5 | 6762 ns/op<br>336 B/op<br>4 allocs/op |
| `BenchmarkCertificateCompression/Decompress` | 5 | 6512 ns/op<br>4248 B/op<br>6 allocs/op |
| `BenchmarkCertificateCompressionHandshakeLifecycle/ServerCertificate/Plain` | 5 | 1030855 ns/op<br>130938 B/op<br>1041 allocs/op |
| `BenchmarkCertificateCompressionHandshakeLifecycle/ServerCertificate/Zlib` | 5 | 1030414 ns/op<br>123300 B/op<br>1020 allocs/op |
| `BenchmarkCertificateCompressionHandshakeLifecycle/MutualTLS/Plain` | 5 | 1709262 ns/op<br>172671 B/op<br>1505 allocs/op |
| `BenchmarkCertificateCompressionHandshakeLifecycle/MutualTLS/Zlib` | 5 | 1712637 ns/op<br>160572 B/op<br>1466 allocs/op |
| `BenchmarkHybridKeyExchangeHandshakeLifecycle/X25519MLKEM768` | 5 | 879036 ns/op<br>147124 B/op<br>791 allocs/op |
| `BenchmarkHybridKeyExchangeHandshakeLifecycle/SecP256r1MLKEM768` | 5 | 833576 ns/op<br>150324 B/op<br>821 allocs/op |
| `BenchmarkHybridKeyExchangeHandshakeLifecycle/SecP384r1MLKEM1024` | 5 | 2081330 ns/op<br>175766 B/op<br>839 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/GoClient/GoServer` | 5 | 39082924 ns/op<br>1.945 go_ms/conn<br>2621456 B/op<br>24336 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/GoClient/WolfSSLServer` | 5 | 92734525 ns/op<br>4.65 go_ms/conn<br>1157904 B/op<br>10562 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/WolfSSLClient/GoServer` | 5 | 45364911 ns/op<br>2.212 wolfssl_ms/conn<br>1440048 B/op<br>12515 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/CertificateAES128GCM/WolfSSLClient/WolfSSLServer` | 5 | 92258966 ns/op<br>4.536 wolfssl_ms/conn<br>34904 B/op<br>53 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ApplicationDataRoundTrip/GoClient/GoServer` | 5 | 39763252 ns/op<br>1.985 go_ms/conn<br>2630080 B/op<br>24655 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ApplicationDataRoundTrip/GoClient/WolfSSLServer` | 5 | 95282136 ns/op<br>4.761 go_ms/conn<br>1176624 B/op<br>10822 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ApplicationDataRoundTrip/WolfSSLClient/GoServer` | 5 | 95935623 ns/op<br>4.788 wolf_process_ms/conn<br>1944232 B/op<br>14006 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ApplicationDataRoundTrip/WolfSSLClient/WolfSSLServer` | 5 | 151095496 ns/op<br>7.515 wolf_process_ms/conn<br>543488 B/op<br>1183 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLS/GoClient/GoServer` | 5 | 71423942 ns/op<br>3.562 go_ms/conn<br>3407352 B/op<br>33637 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLS/GoClient/WolfSSLServer` | 5 | 125098485 ns/op<br>6.245 go_ms/conn<br>1522544 B/op<br>14882 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLS/WolfSSLClient/GoServer` | 5 | 119169276 ns/op<br>5.893 wolfssl_ms/conn<br>1817784 B/op<br>16822 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLS/WolfSSLClient/WolfSSLServer` | 5 | 161373679 ns/op<br>8.002 wolfssl_ms/conn<br>34896 B/op<br>54 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/AES128CCM/GoClient/GoServer` | 5 | 40014042 ns/op<br>1.994 go_ms/conn<br>3110480 B/op<br>28265 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/AES128CCM/GoClient/WolfSSLServer` | 5 | 92480602 ns/op<br>4.621 go_ms/conn<br>1396304 B/op<br>12122 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/AES128CCM/WolfSSLClient/GoServer` | 5 | 45607956 ns/op<br>2.215 wolfssl_ms/conn<br>1504704 B/op<br>13435 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/AES128CCM/WolfSSLClient/WolfSSLServer` | 5 | 91283935 ns/op<br>4.477 wolfssl_ms/conn<br>34904 B/op<br>53 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ExternalPSK/GoClient/GoServer` | 5 | 10400341 ns/op<br>0.5146 go_ms/conn<br>1659104 B/op<br>14233 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ExternalPSK/GoClient/WolfSSLServer` | 5 | 18878596 ns/op<br>0.9418 go_ms/conn<br>784944 B/op<br>6702 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ExternalPSK/WolfSSLClient/GoServer` | 5 | 16847420 ns/op<br>0.789 wolfssl_ms/conn<br>900752 B/op<br>7195 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ExternalPSK/WolfSSLClient/WolfSSLServer` | 5 | 21910543 ns/op<br>1.019 wolfssl_ms/conn<br>34952 B/op<br>53 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ConnectionID/GoClient/GoServer` | 5 | 38998845 ns/op<br>1.945 go_ms/conn<br>2641048 B/op<br>25452 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ConnectionID/GoClient/WolfSSLServer` | 5 | 92908246 ns/op<br>4.643 go_ms/conn<br>1165104 B/op<br>11122 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ConnectionID/WolfSSLClient/GoServer` | 5 | 96000731 ns/op<br>4.792 wolf_process_ms/conn<br>1956368 B/op<br>14565 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/ConnectionID/WolfSSLClient/WolfSSLServer` | 5 | 151032725 ns/op<br>7.509 wolf_process_ms/conn<br>554584 B/op<br>1185 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/GoServer` | 5 | 42022730 ns/op<br>2.104 go_ms/conn<br>2811584 B/op<br>26007 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/GoClient/WolfSSLServer` | 5 | 96009632 ns/op<br>4.79 go_ms/conn<br>1412624 B/op<br>12162 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/WolfSSLClient/GoServer` | 5 | 97159394 ns/op<br>4.849 wolf_process_ms/conn<br>2119704 B/op<br>15645 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/KeyUpdate/WolfSSLClient/WolfSSLServer` | 5 | 150528180 ns/op<br>7.401 wolf_process_ms/conn<br>543808 B/op<br>1183 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/PostHandshakeAuthentication/GoClient/GoServer` | 5 | 101920779 ns/op<br>5.101 go_ms/conn<br>4315368 B/op<br>45058 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/PostHandshakeAuthentication/GoClient/WolfSSLServer` | 5 | 122363576 ns/op<br>6.111 go_ms/conn<br>1581328 B/op<br>15723 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/PostHandshakeAuthentication/WolfSSLClient/GoServer` | 5 | 247311874 ns/op<br>12.37 wolf_process_ms/conn<br>2734328 B/op<br>23435 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/PostHandshakeAuthentication/WolfSSLClient/WolfSSLServer` | 5 | 222420931 ns/op<br>11.1 wolf_process_ms/conn<br>551680 B/op<br>1184 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/SessionResumption/GoClient/GoServer` | 5 | 75048721 ns/op<br>0.5803 go_ms/conn<br>4500152 B/op<br>40559 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/SessionResumption/GoClient/WolfSSLServer` | 5 | 116815492 ns/op<br>1.075 go_ms/conn<br>2065264 B/op<br>17982 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/SessionResumption/WolfSSLClient/GoServer` | 5 | 20166792467 ns/op<br>1008 wolf_process_ms/pair<br>2984824 B/op<br>22965 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/SessionResumption/WolfSSLClient/WolfSSLServer` | 5 | 20235509981 ns/op<br>1012 wolf_process_ms/pair<br>558472 B/op<br>1186 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLSSessionResumption/GoClient/GoServer` | 5 | 107703513 ns/op<br>0.6349 go_ms/conn<br>6619544 B/op<br>53461 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLSSessionResumption/GoClient/WolfSSLServer` | 5 | 148958705 ns/op<br>1.132 go_ms/conn<br>2435344 B/op<br>22342 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/MutualTLSSessionResumption/WolfSSLClient/WolfSSLServer` | 5 | 20307916930 ns/op<br>1015 wolf_process_ms/pair<br>557832 B/op<br>1186 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/EarlyData/GoClient/GoServer` | 5 | 3725958 ns/op<br>0.5929 go_ms/conn<br>290192 B/op<br>2023 allocs/op |
| `BenchmarkWolfSSLFeatureRealUDP/EarlyData/WolfSSLClient/GoServer` | 5 | 1008301350 ns/op<br>1008 wolf_process_ms/pair<br>212464 B/op<br>1138 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/GoClient/GoServer` | 5 | 46153497 ns/op<br>2.318 go_ms/conn<br>3492872 B/op<br>24658 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/GoClient/WolfSSLServer` | 5 | 97104105 ns/op<br>4.855 go_ms/conn<br>1829424 B/op<br>10762 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/WolfSSLClient/GoServer` | 5 | 48473371 ns/op<br>2.353 wolfssl_ms/conn<br>1686768 B/op<br>12995 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/X25519MLKEM768/WolfSSLClient/WolfSSLServer` | 5 | 94948401 ns/op<br>4.667 wolfssl_ms/conn<br>34888 B/op<br>53 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP256r1MLKEM768/GoClient/GoServer` | 5 | 45591961 ns/op<br>2.271 go_ms/conn<br>3531928 B/op<br>24978 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP256r1MLKEM768/GoClient/WolfSSLServer` | 5 | 116468616 ns/op<br>5.815 go_ms/conn<br>1841904 B/op<br>10882 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP256r1MLKEM768/WolfSSLClient/GoServer` | 5 | 69107294 ns/op<br>3.389 wolfssl_ms/conn<br>1706928 B/op<br>13135 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP256r1MLKEM768/WolfSSLClient/WolfSSLServer` | 5 | 120291078 ns/op<br>5.955 wolfssl_ms/conn<br>34888 B/op<br>53 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/GoClient/GoServer` | 5 | 69894594 ns/op<br>3.483 go_ms/conn<br>4027592 B/op<br>25338 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/WolfSSLClient/GoServer` | 5 | 128947774 ns/op<br>6.381 wolfssl_ms/conn<br>1865808 B/op<br>13315 allocs/op |
| `BenchmarkHybridKeyExchangeRealUDP/SecP384r1MLKEM1024/WolfSSLClient/WolfSSLServer` | 5 | 193016948 ns/op<br>9.565 wolfssl_ms/conn<br>34888 B/op<br>53 allocs/op |

[Raw Go benchmark output](benchmark.txt)
