/* Test-only DTLS 1.3 peer for OpenSSL and BoringSSL. */
#include <openssl/ssl.h>
#include <openssl/pem.h>
#include <openssl/err.h>
#ifndef OPENSSL_IS_BORINGSSL
#include <openssl/ech.h>
#include <openssl/sslerr.h>
#endif
#ifdef OPENSSL_IS_BORINGSSL
#include <openssl/pool.h>
#include <openssl/hpke.h>
#endif
#include <arpa/inet.h>
#include <sys/socket.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void check(int ok, const char *what) {
    if (!ok) { fprintf(stderr, "%s\n", what); ERR_print_errors_fp(stderr); exit(1); }
}

static SSL_SESSION *session;
static unsigned char ocsp_response[4096];
static size_t ocsp_length;
static EVP_PKEY *expected_rpk;
static int remember_session(SSL *ssl, SSL_SESSION *value) {
    (void)ssl;
    SSL_SESSION_free(session);
    check(SSL_SESSION_up_ref(value) == 1, "session reference");
    session = value;
    return 0;
}

#ifndef OPENSSL_IS_BORINGSSL
static const unsigned char psk_identity[] = "interop-psk";
static SSL_SESSION *external_psk_session(SSL *ssl) {
    static const unsigned char cipher_id[] = { 0x13, 0x01 };
    unsigned char key[32]; memset(key, 0x42, sizeof(key));
    SSL_SESSION *value = SSL_SESSION_new();
    const SSL_CIPHER *cipher = SSL_CIPHER_find(ssl, cipher_id);
    if (value == NULL || cipher == NULL || !SSL_SESSION_set1_master_key(value, key, sizeof(key)) ||
        !SSL_SESSION_set_cipher(value, cipher) || !SSL_SESSION_set_protocol_version(value, DTLS1_3_VERSION)) {
        SSL_SESSION_free(value);
        return NULL;
    }
    return value;
}
static int use_external_psk(SSL *ssl, const EVP_MD *md, const unsigned char **id,
                            size_t *idlen, SSL_SESSION **value) {
    if (md != NULL && EVP_MD_type(md) != NID_sha256) return 0;
    *value = external_psk_session(ssl);
    *id = psk_identity; *idlen = sizeof(psk_identity) - 1;
    return *value != NULL;
}
static int find_external_psk(SSL *ssl, const unsigned char *id, size_t idlen,
                             SSL_SESSION **value) {
    *value = NULL;
    if (idlen != sizeof(psk_identity) - 1 || memcmp(id, psk_identity, idlen) != 0) return 1;
    *value = external_psk_session(ssl);
    return *value != NULL;
}

static int staple(SSL *ssl, void *arg) {
    (void)arg;
    unsigned char *data = OPENSSL_memdup(ocsp_response, ocsp_length);
    check(data != NULL && SSL_set_tlsext_status_ocsp_resp(ssl, data, ocsp_length), "OCSP staple");
    return SSL_TLSEXT_ERR_OK;
}
#endif

static int alpn(SSL *ssl, const unsigned char **out, unsigned char *len,
                const unsigned char *in, unsigned int inlen, void *arg) {
    static const unsigned char protocols[] = {4, 'c', 'o', 'a', 'p'};
    (void)ssl; (void)arg;
    return SSL_select_next_proto((unsigned char **)out, len, protocols,
                                 sizeof(protocols), in, inlen) == OPENSSL_NPN_NEGOTIATED
           ? SSL_TLSEXT_ERR_OK : SSL_TLSEXT_ERR_ALERT_FATAL;
}

static void trace(int writing, int version, int type, const void *data, size_t len, SSL *ssl, void *arg) {
    (void)ssl; (void)arg;
    if (getenv("DTLS13_PEER_TRACE") && len != 0) {
        const unsigned char *b = data;
        fprintf(stderr, "%s version=%04x type=%d len=%zu first=%02x\n", writing ? "write" : "read", version, type, len, b[0]);
    }
}

static int read_datagram(SSL *ssl, unsigned char *data, int len) {
    int n;
    do { n = SSL_read(ssl, data, len); }
    while (n <= 0 && (SSL_get_error(ssl, n) == SSL_ERROR_WANT_READ || SSL_get_error(ssl, n) == SSL_ERROR_WANT_WRITE));
    if (n <= 0) { fprintf(stderr, "SSL_read=%d error=%d\n", n, SSL_get_error(ssl, n)); ERR_print_errors_fp(stderr); }
    return n;
}

#ifdef OPENSSL_IS_BORINGSSL
static enum ssl_verify_result_t verify_rpk(SSL *ssl, uint8_t *alert) {
    EVP_PKEY *peer = SSL_get0_peer_rpk(ssl);
    if (peer != NULL && expected_rpk != NULL && EVP_PKEY_cmp(peer, expected_rpk) == 1)
        return ssl_verify_ok;
    *alert = SSL_AD_BAD_CERTIFICATE;
    return ssl_verify_invalid;
}

static EVP_PKEY *read_key(const char *path) {
    FILE *f = fopen(path, "rb"); check(f != NULL, "open key");
    EVP_PKEY *key = PEM_read_PrivateKey(f, NULL, NULL, NULL);
    fclose(f); check(key != NULL, "parse key"); return key;
}

static CRYPTO_BUFFER *read_buffer(const char *path) {
    FILE *f = fopen(path, "rb"); check(f != NULL, "open credential");
    unsigned char data[65536]; size_t n = fread(data, 1, sizeof(data), f);
    check(!ferror(f) && feof(f), "read credential"); fclose(f);
    CRYPTO_BUFFER *b = CRYPTO_BUFFER_new(data, n, NULL); check(b != NULL, "credential buffer"); return b;
}
#endif

static EVP_PKEY *read_trusted_public_key(const char *cert_path) {
    char path[4096]; snprintf(path, sizeof(path), "%s.trust", cert_path);
    FILE *f = fopen(path, "rb"); check(f != NULL, "open RPK pin");
    X509 *cert = PEM_read_X509(f, NULL, NULL, NULL);
    fclose(f); check(cert != NULL, "parse RPK pin");
    EVP_PKEY *key = X509_get_pubkey(cert); X509_free(cert);
    check(key != NULL, "extract RPK pin");
    return key;
}

static size_t read_bytes(const char *path, unsigned char *data, size_t capacity) {
    FILE *f = fopen(path, "rb"); check(f != NULL, "open fixture");
    size_t n = fread(data, 1, capacity, f);
    check(!ferror(f) && feof(f) && n != 0, "read fixture"); fclose(f);
    return n;
}

int main(int argc, char **argv) {
    check(argc == 11 || argc == 12, "role port mode cert.pem key.pem cert.der key.der dc.bin delegated.pem delegated.der [benchmark-count]");
    int benchmark = argc == 12;
    int repeats = benchmark ? atoi(argv[11]) : 1;
    check(repeats > 0 && repeats <= 10000, "benchmark count");
    int server = strcmp(argv[1], "server") == 0;
    const char *mode = argv[3];
    int earlyIO = strcmp(mode, "early-io") == 0;
    int early = strcmp(mode, "early") == 0 || earlyIO;
    check(!early || !server, "OpenSSL listener cannot enter SSL_read_early_data");
    if (strcmp(mode, "ocsp") == 0) {
        char path[4096]; snprintf(path, sizeof(path), "%s.ocsp", argv[4]);
        FILE *f = fopen(path, "rb"); check(f != NULL, "OCSP file");
        ocsp_length = fread(ocsp_response, 1, sizeof(ocsp_response), f);
        check(!ferror(f) && feof(f) && ocsp_length != 0, "read OCSP"); fclose(f);
    }
    SSL_CTX *ctx = SSL_CTX_new(DTLS_method()); check(ctx != NULL, "context");
#ifdef OPENSSL_IS_BORINGSSL
    if (strcmp(mode, "ccm-fallback") == 0)
        check(SSL_get_cipher_by_value(0x1304) == NULL, "AES-128-CCM unexpectedly implemented");
    if (strcmp(mode, "imported-psk") == 0) {
        unsigned char key[32]; memset(key, 0x42, sizeof(key));
        static const unsigned char identity[] = "interop-psk", context[] = "interop-context";
        SSL_CREDENTIAL *credential = SSL_CREDENTIAL_new_pre_shared_key(key, sizeof(key), identity,
            sizeof(identity) - 1, EVP_sha256(), context, sizeof(context) - 1);
        check(credential != NULL && SSL_CTX_add1_credential(ctx, credential), "imported PSK credential");
        SSL_CREDENTIAL_free(credential);
    }
#endif
    SSL_CTX_set_msg_callback(ctx, trace);
    SSL_CTX_set_session_cache_mode(ctx, SSL_SESS_CACHE_BOTH);
    if (!server) SSL_CTX_sess_set_new_cb(ctx, remember_session);
    check(SSL_CTX_set_min_proto_version(ctx, DTLS1_3_VERSION) &&
          SSL_CTX_set_max_proto_version(ctx, DTLS1_3_VERSION), "DTLS 1.3 version");
    check(SSL_CTX_use_certificate_chain_file(ctx, argv[4]) == 1 &&
          SSL_CTX_use_PrivateKey_file(ctx, argv[5], SSL_FILETYPE_PEM) == 1, "X.509 credential");
    SSL_CTX_set_verify(ctx, !server ? SSL_VERIFY_PEER : strcmp(mode, "mutual") == 0 || strcmp(mode, "mutual-resume") == 0 || strcmp(mode, "pha") == 0 ? SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT : SSL_VERIFY_NONE, NULL);
#ifndef OPENSSL_IS_BORINGSSL
    if (!server && strcmp(mode, "ech-client-unavailable") == 0) {
        unsigned char config[512];
        size_t n = read_bytes(argv[8], config, sizeof(config));
        SSL *probe = SSL_new(ctx);
        check(probe != NULL, "ECH probe");
        BIO *input = BIO_new(BIO_s_mem()), *output = BIO_new(BIO_s_mem());
        check(input != NULL && output != NULL, "ECH probe BIOs");
        SSL_set_bio(probe, input, output);
        SSL_set_connect_state(probe);
        check(SSL_set1_ech_config_list(probe, config, n), "ECH probe config");
        check(SSL_set_tlsext_host_name(probe, "server.test"), "ECH probe SNI");
        ERR_clear_error();
        int result = SSL_do_handshake(probe);
        check(result != 1 && ERR_GET_REASON(ERR_peek_error()) == SSL_R_UNSUPPORTED_SSL_VERSION, "DTLS ECH version rejection");
        SSL_free(probe); ERR_clear_error();
    }
#endif
    if (strcmp(mode, "rpk") == 0) {
        static const unsigned char rpk[] = { TLSEXT_cert_type_rpk };
#ifdef OPENSSL_IS_BORINGSSL
        EVP_PKEY *key = read_key(argv[5]);
        SSL_CREDENTIAL *credential = SSL_CREDENTIAL_new_raw_public_key(key);
        check(credential != NULL && SSL_CTX_add1_credential(ctx, credential), "RPK credential");
        SSL_CREDENTIAL_free(credential); EVP_PKEY_free(key);
        if (!server) {
            expected_rpk = read_trusted_public_key(argv[4]);
            check(SSL_CTX_set1_accepted_peer_cert_types(ctx, rpk, sizeof(rpk)), "accept RPK");
            SSL_CTX_set_custom_verify(ctx, SSL_VERIFY_PEER, verify_rpk);
        }
#else
        check(SSL_CTX_set1_server_cert_type(ctx, rpk, sizeof(rpk)), "offer RPK");
        if (!server) check(SSL_CTX_dane_enable(ctx) > 0, "enable RPK pinning");
        if (!server) expected_rpk = read_trusted_public_key(argv[4]);
#endif
    }
    if (!server || strcmp(mode, "mutual") == 0 || strcmp(mode, "mutual-resume") == 0 || strcmp(mode, "pha") == 0) {
        char trust[4096]; snprintf(trust, sizeof(trust), "%s.trust", argv[4]);
        check(SSL_CTX_load_verify_locations(ctx, trust, NULL) == 1, "mutual trust");
    }
#ifndef OPENSSL_IS_BORINGSSL
    if (strcmp(mode, "external-psk") == 0) {
        static const unsigned char context[] = "interop";
        if (server) {
            check(SSL_CTX_set_session_id_context(ctx, context, sizeof(context)), "PSK context");
            SSL_CTX_set_psk_find_session_callback(ctx, find_external_psk);
        } else {
            SSL_CTX_set_psk_use_session_callback(ctx, use_external_psk);
        }
    }
    if (server && strcmp(mode, "ech-server-rejection") == 0) {
        BIO *pem = BIO_new_file(argv[9], "r");
        OSSL_ECHSTORE *store = OSSL_ECHSTORE_new(NULL, NULL);
        check(pem != NULL && store != NULL && OSSL_ECHSTORE_read_pem(store, pem, OSSL_ECH_FOR_RETRY) &&
              SSL_CTX_set1_echstore(ctx, store), "ECH server keys");
        OSSL_ECHSTORE_free(store); BIO_free(pem);
    }
    if (early) check(SSL_CTX_set_max_early_data(ctx, 256), "early data limit");
    SSL_CTX_set_post_handshake_auth(ctx, 1);
    if (server && strcmp(mode, "pha") == 0) SSL_CTX_set_verify(ctx, SSL_VERIFY_PEER | SSL_VERIFY_POST_HANDSHAKE, NULL);
    if (server && strcmp(mode, "ocsp") == 0) SSL_CTX_set_tlsext_status_cb(ctx, staple);
#else
    if (server && strcmp(mode, "ocsp") == 0) check(SSL_CTX_set_ocsp_response(ctx, ocsp_response, ocsp_length), "OCSP staple");
    if (server && strcmp(mode, "ech") == 0) {
        unsigned char config[512], private_key[64];
        size_t config_len = read_bytes(argv[8], config, sizeof(config));
        size_t key_len = read_bytes(argv[9], private_key, sizeof(private_key));
        EVP_HPKE_KEY *key = EVP_HPKE_KEY_new();
        SSL_ECH_KEYS *keys = SSL_ECH_KEYS_new();
        check(key != NULL && keys != NULL && EVP_HPKE_KEY_init(key, EVP_hpke_x25519_hkdf_sha256(), private_key, key_len) &&
              SSL_ECH_KEYS_add(keys, 1, config, config_len, key) && SSL_CTX_set1_ech_keys(ctx, keys), "ECH server keys");
        SSL_ECH_KEYS_free(keys); EVP_HPKE_KEY_free(key);
    }
#endif
    if (strcmp(mode, "alpn") == 0) {
        if (server) SSL_CTX_set_alpn_select_cb(ctx, alpn, NULL);
        else { static const unsigned char p[] = {4, 'c','o','a','p'}; check(SSL_CTX_set_alpn_protos(ctx, p, sizeof(p)) == 0, "ALPN"); }
    }
    const char *group = strcmp(mode, "p256") == 0 || strcmp(mode, "hrr") == 0 ? "P-256" : strcmp(mode, "p384") == 0 ? "P-384" : strcmp(mode, "hybrid") == 0 ? "X25519MLKEM768" : "X25519";
    if (!server && strcmp(mode, "hrr") == 0) group = "X25519:P-256";
    check(SSL_CTX_set1_groups_list(ctx, group) == 1, "groups");
    const char *cipher = (strcmp(mode, "sha384") == 0 || strcmp(mode, "dc-sha384") == 0) ? "TLS_AES_256_GCM_SHA384" : strcmp(mode, "chacha") == 0 ? "TLS_CHACHA20_POLY1305_SHA256" : strcmp(mode, "ccm") == 0 ? "TLS_AES_128_CCM_SHA256" : "TLS_AES_128_GCM_SHA256";
    uint16_t cipher_id = (strcmp(mode, "sha384") == 0 || strcmp(mode, "dc-sha384") == 0) ? 0x1302 : strcmp(mode, "chacha") == 0 ? 0x1303 : strcmp(mode, "ccm") == 0 ? 0x1304 : 0x1301;
#ifndef OPENSSL_IS_BORINGSSL
    check(SSL_CTX_set_ciphersuites(ctx, cipher) == 1, "cipher suites");
#else
    check(SSL_CTX_set1_tls13_ciphers(ctx, &cipher_id, NULL, 1), "TLS 1.3 cipher suite");
    (void)cipher;
    if (server && strncmp(mode, "dc", 2) == 0) {
        SSL_CREDENTIAL *cred = SSL_CREDENTIAL_new_delegated(); check(cred != NULL, "new DC");
        CRYPTO_BUFFER *cert = read_buffer(argv[6]), *dc = read_buffer(argv[8]);
        EVP_PKEY *key = read_key(argv[9]);
        check(SSL_CREDENTIAL_set1_cert_chain(cred, &cert, 1) &&
              SSL_CREDENTIAL_set1_private_key(cred, key) &&
              SSL_CREDENTIAL_set1_delegated_credential(cred, dc) &&
              SSL_CTX_add1_credential(ctx, cred), "configure DC");
        SSL_CREDENTIAL_free(cred); CRYPTO_BUFFER_free(cert); CRYPTO_BUFFER_free(dc); EVP_PKEY_free(key);
    }
#endif
    int pair = strcmp(mode, "resume") == 0 || strcmp(mode, "mutual-resume") == 0 || strcmp(mode, "dc-resume") == 0 || strcmp(mode, "ticket-request-fallback") == 0 || early;
    if (server && strcmp(mode, "mutual-resume") == 0) {
        static const unsigned char context[] = "interop-mtls";
        check(SSL_CTX_set_session_id_context(ctx, context, sizeof(context)), "mTLS session context");
    }
    if (benchmark && !pair) SSL_CTX_set_options(ctx, SSL_OP_NO_TICKET);
    int rounds = (pair ? 2 : 1) * repeats;
    for (int round = 0; round < rounds; round++) {
        if (benchmark && !server && round > 0 && getenv("DTLS13_BENCH_SYNC")) {
            unsigned char ready;
            printf("WAIT %d\n", round); fflush(stdout);
            check(read(STDIN_FILENO, &ready, 1) == 1 && ready == 1, "benchmark server ready");
        }
        int resumed = pair && round % 2 != 0;
        if (benchmark && pair && !resumed) { SSL_SESSION_free(session); session = NULL; }
        int fd = socket(AF_INET, SOCK_DGRAM, 0); check(fd >= 0, "socket");
        struct sockaddr_in address = {0}; address.sin_family = AF_INET;
        address.sin_addr.s_addr = htonl(INADDR_LOOPBACK); address.sin_port = htons((unsigned short)atoi(argv[2]));
        struct timeval timeout = {8, 0}; setsockopt(fd, SOL_SOCKET, SO_RCVTIMEO, &timeout, sizeof(timeout));
        SSL *listener = NULL, *ssl = NULL;
        if (server) {
            check(bind(fd, (struct sockaddr *)&address, sizeof(address)) == 0, "bind");
            printf("READY %d\n", round); fflush(stdout);
#ifndef OPENSSL_IS_BORINGSSL
            listener = SSL_new_listener(ctx, SSL_LISTENER_FLAG_SINGLE_THREAD);
            check(listener != NULL, "DTLS listener");
            BIO *bio = BIO_new_dgram(fd, BIO_NOCLOSE); check(bio != NULL, "datagram BIO");
            SSL_set_bio(listener, bio, bio);
            check(SSL_listen(listener) == 1, "listen");
            ssl = SSL_accept_connection(listener, 0); check(ssl != NULL, "accept connection");
#else
            unsigned char b; socklen_t size = sizeof(address);
            do {
                check(recvfrom(fd, &b, 1, MSG_PEEK, (struct sockaddr *)&address, &size) > 0, "first datagram");
                if (!benchmark || b == 22) break;
                unsigned char discard[2048];
                check(recvfrom(fd, discard, sizeof(discard), 0, (struct sockaddr *)&address, &size) > 0, "discard previous association");
            } while (1);
            check(connect(fd, (struct sockaddr *)&address, size) == 0, "connect server UDP");
#endif
        } else {
            check(connect(fd, (struct sockaddr *)&address, sizeof(address)) == 0, "connect client UDP");
        }
        if (ssl == NULL) {
            ssl = SSL_new(ctx); check(ssl != NULL, "connection");
#ifdef OPENSSL_IS_BORINGSSL
            check(SSL_set_fd(ssl, fd) == 1, "set fd");
#else
            BIO *bio = BIO_new_dgram(fd, BIO_NOCLOSE); check(bio != NULL, "datagram BIO");
            BIO_ctrl(bio, BIO_CTRL_DGRAM_SET_CONNECTED, 0, &address);
            SSL_set_bio(ssl, bio, bio);
#endif
        }
        if (strcmp(mode, "fragmented") == 0 || strcmp(mode, "dc-fragmented") == 0) check(SSL_set_mtu(ssl, 256) > 0, "MTU");
#ifdef OPENSSL_IS_BORINGSSL
        if (!server && strcmp(mode, "hrr") == 0) {
            uint16_t share = SSL_GROUP_X25519;
            check(SSL_set1_client_key_shares(ssl, &share, 1), "initial key share");
        }
#endif
        if (!server) {
            SSL_set_connect_state(ssl);
            if (strcmp(mode, "ech") == 0) {
                unsigned char config[512];
                size_t n = read_bytes(argv[8], config, sizeof(config));
                check(SSL_set1_ech_config_list(ssl, config, n), "ECH client config");
            }
            check(SSL_set_tlsext_host_name(ssl, "server.test") == 1, "SNI");
            check(X509_VERIFY_PARAM_set1_host(SSL_get0_param(ssl), "server.test", 11) == 1, "expected hostname");
#ifndef OPENSSL_IS_BORINGSSL
            if (strcmp(mode, "rpk") == 0) {
                check(SSL_dane_enable(ssl, "server.test") > 0, "enable RPK verification");
                check(SSL_add_expected_rpk(ssl, expected_rpk), "RPK pin");
            }
#endif
            if (strcmp(mode, "ocsp") == 0) {
#ifdef OPENSSL_IS_BORINGSSL
                SSL_enable_ocsp_stapling(ssl);
#else
                check(SSL_set_tlsext_status_type(ssl, TLSEXT_STATUSTYPE_ocsp), "OCSP request");
#endif
            }
        }
        if (!server && resumed) check(session != NULL && SSL_set_session(ssl, session), "restore session");
#ifndef OPENSSL_IS_BORINGSSL
        if (early && !server && resumed) {
            size_t written; check(SSL_write_early_data(ssl, "early", 5, &written) && written == 5, "write early data");
        }
#endif
        check((server ? SSL_accept(ssl) : SSL_connect(ssl)) == 1, "handshake");
#ifndef OPENSSL_IS_BORINGSSL
        if (early && resumed) check(SSL_get_early_data_status(ssl) == SSL_EARLY_DATA_ACCEPTED, "early data accepted");
#endif
        check(SSL_session_reused(ssl) == (resumed || strcmp(mode, "external-psk") == 0), "session reuse");
        if (!server && (strcmp(mode, "external-psk") == 0 || strcmp(mode, "imported-psk") == 0)) {
            X509 *peer_cert = SSL_get_peer_certificate(ssl);
            check(peer_cert == NULL, "external PSK selected");
            X509_free(peer_cert);
        }
        check(SSL_version(ssl) == DTLS1_3_VERSION, "negotiated version");
        check(SSL_CIPHER_get_protocol_id(SSL_get_current_cipher(ssl)) == cipher_id, "negotiated cipher suite");
#ifdef OPENSSL_IS_BORINGSSL
        if (strcmp(mode, "ech") == 0) check(SSL_ech_accepted(ssl), "ECH not accepted");
#else
        if (strcmp(mode, "ech") == 0) check(SSL_ech_get1_status(ssl, NULL, NULL) == SSL_ECH_STATUS_SUCCESS, "ECH not accepted");
#endif
        if (strcmp(mode, "rpk") == 0) {
#ifdef OPENSSL_IS_BORINGSSL
            check(server || SSL_get_peer_cert_type(ssl) == TLSEXT_cert_type_rpk, "negotiated server RPK");
#else
            check(SSL_get_negotiated_server_cert_type(ssl) == TLSEXT_cert_type_rpk, "negotiated server RPK");
#endif
        }
        if (server) { const char *name = SSL_get_servername(ssl, TLSEXT_NAMETYPE_host_name); check(name && strcmp(name, "server.test") == 0, "received SNI"); }
        if (!server && strcmp(mode, "ocsp") == 0) {
            const unsigned char *response; size_t length;
#ifdef OPENSSL_IS_BORINGSSL
            SSL_get0_ocsp_response(ssl, &response, &length);
#else
            length = SSL_get_tlsext_status_ocsp_resp(ssl, &response);
#endif
            check(length == ocsp_length && memcmp(response, ocsp_response, length) == 0, "received OCSP");
        }
        if (strcmp(mode, "alpn") == 0) { const unsigned char *p; unsigned n; SSL_get0_alpn_selected(ssl, &p, &n); check(n == 4 && memcmp(p, "coap", 4) == 0, "negotiated ALPN"); }
        unsigned char exporter[32]; check(SSL_export_keying_material(ssl, exporter, sizeof(exporter), "interop", 7, NULL, 0, 0) == 1, "exporter");
        printf("ROUND=%d VERSION=%04x CIPHER=%s EXPORTER=", round, SSL_version(ssl), SSL_get_cipher(ssl));
        for (unsigned i = 0; i < sizeof(exporter); i++) printf("%02x", exporter[i]);
        puts(""); fflush(stdout);
        unsigned char buf[64]; int n;
        if (earlyIO && resumed) {
            n = read_datagram(ssl, buf, sizeof(buf));
            check(n == 11 && memcmp(buf, "early reply", 11) == 0, "early response");
        }
        if (strcmp(mode, "keyupdate") == 0 && (!benchmark || !server || getenv("DTLS13_BENCH_SERVER_KEYUPDATE"))) {
            check(SSL_key_update(ssl, SSL_KEY_UPDATE_REQUESTED) == 1, "KeyUpdate");
#ifdef OPENSSL_IS_BORINGSSL
            if (benchmark) check(SSL_do_handshake(ssl) == 1, "flush KeyUpdate");
#endif
        }
        if (!server) check(SSL_write(ssl, "interop", 7) == 7, "write datagram");
        n = read_datagram(ssl, buf, sizeof(buf)); check(n == 7 && memcmp(buf, "interop", 7) == 0, "read datagram");
#ifndef OPENSSL_IS_BORINGSSL
        if (server && strcmp(mode, "pha") == 0) check(SSL_verify_client_post_handshake(ssl), "request PHA");
#endif
        if (server) check(SSL_write(ssl, buf, n) == n, "echo datagram");
        if (benchmark) check(SSL_write(ssl, "done", 4) == 4, "benchmark completion");
        /* Keep the association alive until the other endpoint completes. */
        n = read_datagram(ssl, buf, sizeof(buf)); check(n == 4 && memcmp(buf, "done", 4) == 0, "completion datagram");
        if (server && strcmp(mode, "pha") == 0) { X509 *cert = SSL_get_peer_certificate(ssl); check(cert != NULL && SSL_get_verify_result(ssl) == X509_V_OK, "PHA identity"); X509_free(cert); }
        /* The authenticated completion datagram closes the test exchange. */
        SSL_set_shutdown(ssl, SSL_SENT_SHUTDOWN | SSL_RECEIVED_SHUTDOWN);
        SSL_free(ssl); SSL_free(listener); close(fd);
    }
    SSL_SESSION_free(session); SSL_CTX_free(ctx); EVP_PKEY_free(expected_rpk);
    return 0;
}
