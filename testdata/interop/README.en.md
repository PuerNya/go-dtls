# DTLS 1.3 peer fixtures

[简体中文](README.md) | [English](README.en.md) | [Русский](README.ru.md)

These opt-in tests use native UDP record layers from unmodified libraries.
`openssl.c` serves OpenSSL and BoringSSL; `nss.c` serves NSS. The Go test creates
short-lived credentials, pins trust, compares exporters, and checks application
data and negotiated state. Missing peer environment variables skip that peer.

## Source revisions

The matrix was validated with these exact upstream revisions:

| Library | Repository | Commit |
| --- | --- | --- |
| NSS | https://github.com/nss-dev/nss | `b9283c1e4ee232b4450d36115272af670d38ec81` |
| BoringSSL | https://github.com/google/boringssl | `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978` |
| OpenSSL | https://github.com/openssl/openssl | `41eecc7c7cd4069b1b29bca49b9d97825c6bf648` |

All three negotiate the final RFC 9147 DTLS 1.3 version, `0xfefc`. A successful
TLS 1.3 test or a build alone does not establish DTLS interoperability.

## Linux build and run

Install a native C/C++ toolchain, Git, Make, CMake, Perl, Python, pkg-config and
NSPR development files. Keep sources and build outputs outside this checkout.
From the go-dtls repository, set `repo` and an absolute external `work` directory:

```sh
repo="$PWD"
work="$(mktemp -d)"
git clone https://github.com/nss-dev/nss "$work/nss"
git -C "$work/nss" checkout b9283c1e4ee232b4450d36115272af670d38ec81
git clone https://github.com/google/boringssl "$work/bssl"
git -C "$work/bssl" checkout 4f3b183c5d4a8f8fcaba4c0ab4435df181df6978
git clone https://github.com/openssl/openssl "$work/openssl"
git -C "$work/openssl" checkout 41eecc7c7cd4069b1b29bca49b9d97825c6bf648

(cd "$work/openssl" && ./Configure no-tests && make -j2)
cmake -S "$work/bssl" -B "$work/bssl-build" -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTING=OFF
cmake --build "$work/bssl-build" -j2
make -C "$work/nss" -j2 all USE_64=1 BUILD_OPT=1 NSS_DISABLE_GTESTS=1 \
  NSPR_INCLUDE_DIR="$(pkg-config --variable=includedir nspr)" \
  NSPR_LIB_DIR="$(pkg-config --variable=libdir nspr)"
make -C "$work/nss" latest USE_64=1 BUILD_OPT=1

cc -O2 -Wall -Wextra -Werror -I"$work/openssl/include" \
  "$repo/testdata/interop/openssl.c" -L"$work/openssl" \
  -Wl,-rpath,"$work/openssl" -lssl -lcrypto -o "$work/openssl-peer"
cc -O2 -Wall -Wextra -Werror -I"$work/bssl/include" \
  "$repo/testdata/interop/openssl.c" "$work/bssl-build/libssl.a" \
  "$work/bssl-build/libcrypto.a" -lstdc++ -pthread -o "$work/bssl-peer"
nsslib="$(find "$work/dist" -name libnss3.so -printf '%h\n' | head -n1)"
cc -O2 -Wall -Wextra -Werror -I"$work/dist/public/nss" \
  $(pkg-config --cflags nspr) "$repo/testdata/interop/nss.c" \
  -L"$nsslib" -Wl,-rpath,"$nsslib" -lssl3 -lnss3 -lnssutil3 -lsmime3 \
  $(pkg-config --libs nspr) -o "$work/nss-peer"

export DTLS13_NSS_PEER="$work/nss-peer"
export DTLS13_BORINGSSL_PEER="$work/bssl-peer"
export DTLS13_OPENSSL_PEER="$work/openssl-peer"
export LD_LIBRARY_PATH="$nsslib:$work/openssl"
cd "$repo"
go test . -run '^TestInteropDTLS13Peers$' -count=1 -timeout=3m -v
```

Finish downloads before running tests and clear inherited proxy settings for
the test processes. The test itself neither downloads nor compiles peers.
Existing wolfSSL tests remain under `TestInteropWolfSSL*`; their build is defined
in the [workflow](../../.github/workflows/benchmarks.yml), independently of
whether benchmarks are run.

## Coverage and limits

The matrix covers both peer roles: X.509, resumption, AES-GCM/SHA-256 and
SHA-384, ChaCha20-Poly1305, X25519/P-256/P-384, key-share HRR, fragmentation,
mTLS, OCSP, SNI, ALPN, GREASE, KeyUpdate, and compression/cached-info fallback.
Delegated-credential cases check actual DC state, resumption, fragmentation,
SHA-384, fallback, invalid signatures and expiry. Negative probes are named
separately and are not successful handshakes.

| Peer | Additional tested behavior | Boundary at the pinned revision |
| --- | --- | --- |
| NSS | Server DC presentation; client verification of Go server DC; client 0-RTT including a DC-authenticated session | Rejects DC in CertificateRequest (`tls13con.c`, `KnownExtensions`), contrary to RFC 9345 §4.1.2. Ordinary mTLS disables that Go offer. Native DTLS PHA is unsupported. |
| BoringSSL | Server DC presentation | Client DC verification, PHA and DTLS 0-RTT are unavailable. `tls13_client.cc` explicitly excludes DTLS 0-RTT. |
| OpenSSL | Both PHA directions; client 0-RTT | No DC. Its server rejects a valid initial ACK with `unexpected_message`; positive server cases deliberately drop epoch-2 ACKs and are named `server-ACKLoss`. `server/initial-ACK-rejection` preserves the unmodified failure. |

NSS client 0-RTT sets `SSL_SUPPRESS_END_OF_EARLY_DATA`: RFC 9147 §5 omits
EndOfEarlyData. Without it, NSS sends that forbidden message. NSS documents
that its server needs `SSL_RecordLayerData` to suppress EndOfEarlyData; the
native record-layer server stalls even with the option set, so it is excluded
from positive 0-RTT cases. OpenSSL's listener consumes ClientHello before
returning the connection; `SSL_read_early_data` then rejects the call because
`SSL_in_before` is false. Its server 0-RTT path is also unverified here.

NSS automatic DTLS tickets work and are covered; the separate
`SSL_SendSessionTicket` API rejecting DTLS does not imply that resumption is
unsupported. NSS authenticates an exact per-test DER pin; OpenSSL and BoringSSL
verify a test CA and hostname. None of the fixtures disables peer authentication.
