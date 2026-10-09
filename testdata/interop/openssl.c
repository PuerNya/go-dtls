/* Test-only DTLS 1.3 peer for OpenSSL and BoringSSL. */
#include <openssl/ssl.h>
#include <openssl/pem.h>
#include <openssl/err.h>
#ifdef OPENSSL_IS_BORINGSSL
#include <openssl/pool.h>
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
static int remember_session(SSL *ssl, SSL_SESSION *value) {
    (void)ssl;
    SSL_SESSION_free(session);
    check(SSL_SESSION_up_ref(value) == 1, "session reference");
    session = value;
    return 0;
}

#ifndef OPENSSL_IS_BORINGSSL
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

int main(int argc, char **argv) {
    check(argc == 11, "role port mode cert.pem key.pem cert.der key.der dc.bin delegated.pem delegated.der");
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
    SSL_CTX_set_msg_callback(ctx, trace);
    SSL_CTX_set_session_cache_mode(ctx, SSL_SESS_CACHE_BOTH);
    if (!server) SSL_CTX_sess_set_new_cb(ctx, remember_session);
    check(SSL_CTX_set_min_proto_version(ctx, DTLS1_3_VERSION) &&
          SSL_CTX_set_max_proto_version(ctx, DTLS1_3_VERSION), "DTLS 1.3 version");
    check(SSL_CTX_use_certificate_chain_file(ctx, argv[4]) == 1 &&
          SSL_CTX_use_PrivateKey_file(ctx, argv[5], SSL_FILETYPE_PEM) == 1, "X.509 credential");
    SSL_CTX_set_verify(ctx, !server ? SSL_VERIFY_PEER : strcmp(mode, "mutual") == 0 || strcmp(mode, "pha") == 0 ? SSL_VERIFY_PEER | SSL_VERIFY_FAIL_IF_NO_PEER_CERT : SSL_VERIFY_NONE, NULL);
    if (!server || strcmp(mode, "mutual") == 0 || strcmp(mode, "pha") == 0) {
        char trust[4096]; snprintf(trust, sizeof(trust), "%s.trust", argv[4]);
        check(SSL_CTX_load_verify_locations(ctx, trust, NULL) == 1, "mutual trust");
    }
#ifndef OPENSSL_IS_BORINGSSL
    if (early) check(SSL_CTX_set_max_early_data(ctx, 256), "early data limit");
    SSL_CTX_set_post_handshake_auth(ctx, 1);
    if (server && strcmp(mode, "pha") == 0) SSL_CTX_set_verify(ctx, SSL_VERIFY_PEER | SSL_VERIFY_POST_HANDSHAKE, NULL);
    if (server && strcmp(mode, "ocsp") == 0) SSL_CTX_set_tlsext_status_cb(ctx, staple);
#else
    if (server && strcmp(mode, "ocsp") == 0) check(SSL_CTX_set_ocsp_response(ctx, ocsp_response, ocsp_length), "OCSP staple");
#endif
    if (strcmp(mode, "alpn") == 0) {
        if (server) SSL_CTX_set_alpn_select_cb(ctx, alpn, NULL);
        else { static const unsigned char p[] = {4, 'c','o','a','p'}; check(SSL_CTX_set_alpn_protos(ctx, p, sizeof(p)) == 0, "ALPN"); }
    }
    const char *group = strcmp(mode, "p256") == 0 || strcmp(mode, "hrr") == 0 ? "P-256" : strcmp(mode, "p384") == 0 ? "P-384" : "X25519";
    if (!server && strcmp(mode, "hrr") == 0) group = "X25519:P-256";
    check(SSL_CTX_set1_groups_list(ctx, group) == 1, "groups");
    const char *cipher = (strcmp(mode, "sha384") == 0 || strcmp(mode, "dc-sha384") == 0) ? "TLS_AES_256_GCM_SHA384" : strcmp(mode, "chacha") == 0 ? "TLS_CHACHA20_POLY1305_SHA256" : "TLS_AES_128_GCM_SHA256";
#ifndef OPENSSL_IS_BORINGSSL
    check(SSL_CTX_set_ciphersuites(ctx, cipher) == 1, "cipher suites");
#else
    /* BoringSSL's TLS 1.3 cipher list is fixed; the Go peer restricts selection. */
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
    int rounds = (strcmp(mode, "resume") == 0 || strcmp(mode, "dc-resume") == 0 || early) ? 2 : 1;
    for (int round = 0; round < rounds; round++) {
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
            check(recvfrom(fd, &b, 1, MSG_PEEK, (struct sockaddr *)&address, &size) > 0, "first datagram");
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
            check(SSL_set_tlsext_host_name(ssl, "server.test") == 1, "SNI");
            check(X509_VERIFY_PARAM_set1_host(SSL_get0_param(ssl), "server.test", 11) == 1, "expected hostname");
            if (strcmp(mode, "ocsp") == 0) {
#ifdef OPENSSL_IS_BORINGSSL
                SSL_enable_ocsp_stapling(ssl);
#else
                check(SSL_set_tlsext_status_type(ssl, TLSEXT_STATUSTYPE_ocsp), "OCSP request");
#endif
            }
        }
        if (!server && round != 0) check(session != NULL && SSL_set_session(ssl, session), "restore session");
#ifndef OPENSSL_IS_BORINGSSL
        if (early && !server && round != 0) {
            size_t written; check(SSL_write_early_data(ssl, "early", 5, &written) && written == 5, "write early data");
        }
#endif
        check((server ? SSL_accept(ssl) : SSL_connect(ssl)) == 1, "handshake");
#ifndef OPENSSL_IS_BORINGSSL
        if (early && round != 0) check(SSL_get_early_data_status(ssl) == SSL_EARLY_DATA_ACCEPTED, "early data accepted");
#endif
        check(SSL_session_reused(ssl) == (round != 0), "session reuse");
        check(SSL_version(ssl) == DTLS1_3_VERSION, "negotiated version");
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
        printf("VERSION=%04x CIPHER=%s EXPORTER=", SSL_version(ssl), SSL_get_cipher(ssl));
        for (unsigned i = 0; i < sizeof(exporter); i++) printf("%02x", exporter[i]);
        puts(""); fflush(stdout);
        unsigned char buf[64]; int n;
        if (earlyIO && round != 0) {
            n = read_datagram(ssl, buf, sizeof(buf));
            check(n == 11 && memcmp(buf, "early reply", 11) == 0, "early response");
        }
        if (strcmp(mode, "keyupdate") == 0) check(SSL_key_update(ssl, SSL_KEY_UPDATE_REQUESTED) == 1, "KeyUpdate");
        if (!server) check(SSL_write(ssl, "interop", 7) == 7, "write datagram");
        n = read_datagram(ssl, buf, sizeof(buf)); check(n == 7 && memcmp(buf, "interop", 7) == 0, "read datagram");
#ifndef OPENSSL_IS_BORINGSSL
        if (server && strcmp(mode, "pha") == 0) check(SSL_verify_client_post_handshake(ssl), "request PHA");
#endif
        if (server) check(SSL_write(ssl, buf, n) == n, "echo datagram");
        /* Keep the association alive until Go has received the echo and ACKed KU. */
        n = read_datagram(ssl, buf, sizeof(buf)); check(n == 4 && memcmp(buf, "done", 4) == 0, "completion datagram");
        if (server && strcmp(mode, "pha") == 0) { X509 *cert = SSL_get_peer_certificate(ssl); check(cert != NULL && SSL_get_verify_result(ssl) == X509_V_OK, "PHA identity"); X509_free(cert); }
        /* The authenticated completion datagram closes the test exchange. */
        SSL_set_shutdown(ssl, SSL_SENT_SHUTDOWN | SSL_RECEIVED_SHUTDOWN);
        SSL_free(ssl); SSL_free(listener); close(fd);
    }
    SSL_SESSION_free(session); SSL_CTX_free(ctx);
    return 0;
}
