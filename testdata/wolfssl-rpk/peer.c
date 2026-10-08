/* Authenticated RFC 7250 DTLS 1.3 interoperability peer. */
#ifdef _WIN32
#include <winsock2.h>
#include <ws2tcpip.h>
typedef int socklen_t;
#else
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <unistd.h>
#define closesocket close
#define SOCKET int
#define INVALID_SOCKET (-1)
#endif
#include <wolfssl/options.h>
#include <wolfssl/ssl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int read_file(const char* name, unsigned char* data, size_t capacity)
{
    FILE* file = fopen(name, "rb");
    size_t size;
    if (file == NULL) return -1;
    size = fread(data, 1, capacity, file);
    if (ferror(file) || size == capacity) size = 0;
    fclose(file);
    return size == 0 ? -1 : (int)size;
}
#define CHECK(call) do { if (!(call)) { \
    fprintf(stderr, "%s failed (wolfSSL %d)\n", #call, \
        ssl == NULL ? 0 : wolfSSL_get_error(ssl, -1)); goto cleanup; } } while (0)

int main(int argc, char** argv)
{
    WOLFSSL_CTX* ctx = NULL;
    WOLFSSL* ssl = NULL;
    WOLFSSL_SESSION* session = NULL;
    SOCKET fd = INVALID_SOCKET;
    struct sockaddr_in address;
    unsigned char data[8192];
    int server, result, status = 1, port, resume, early, pha, ech, serverOnly;
    int rawLocal, rawPeer, round, step, count;
    char serverType, clientType;
    const char* mode;
#ifdef _WIN32
    WSADATA wsa;
    DWORD timeout = 5000;
    if (WSAStartup(MAKEWORD(2, 2), &wsa) != 0) return 1;
#else
    struct timeval timeout = {5, 0};
#endif
    setvbuf(stdout, NULL, _IONBF, 0);
    if (argc != 8 || (strcmp(argv[1], "server") != 0 && strcmp(argv[1], "client") != 0)) {
        fprintf(stderr, "usage: peer server|client port credential key trust scenario ech-file\n");
        return 1;
    }
    server = strcmp(argv[1], "server") == 0;
    port = atoi(argv[2]);
    mode = argv[6];
    early = strcmp(mode, "early") == 0;
    pha = strcmp(mode, "pha") == 0 || strcmp(mode, "resume-pha") == 0;
    resume = early || strcmp(mode, "resume") == 0 || strcmp(mode, "resume-pha") == 0;
    ech = strcmp(mode, "ech") == 0;
    serverOnly = strcmp(mode, "server-only") == 0;
    serverType = strcmp(mode, "raw-client") == 0 ? WOLFSSL_CERT_TYPE_X509 : WOLFSSL_CERT_TYPE_RPK;
    clientType = strcmp(mode, "raw-server") == 0 ? WOLFSSL_CERT_TYPE_X509 : WOLFSSL_CERT_TYPE_RPK;
    rawLocal = (server ? serverType : clientType) == WOLFSSL_CERT_TYPE_RPK;
    rawPeer = (server ? clientType : serverType) == WOLFSSL_CERT_TYPE_RPK;
    CHECK(server || (port > 0 && port <= 65535));
    CHECK(wolfSSL_Init() == WOLFSSL_SUCCESS);
    ctx = wolfSSL_CTX_new(server ? wolfDTLSv1_3_server_method() : wolfDTLSv1_3_client_method());
    CHECK(ctx != NULL);
    CHECK(wolfSSL_CTX_set_server_cert_type(ctx, &serverType, 1) == WOLFSSL_SUCCESS);
    CHECK(wolfSSL_CTX_set_client_cert_type(ctx, &clientType, 1) == WOLFSSL_SUCCESS);
    if (server || !serverOnly) {
        CHECK(wolfSSL_CTX_use_certificate_file(ctx, argv[3], rawLocal ? WOLFSSL_FILETYPE_ASN1 : WOLFSSL_FILETYPE_PEM) == WOLFSSL_SUCCESS);
        CHECK(wolfSSL_CTX_use_PrivateKey_file(ctx, argv[4], WOLFSSL_FILETYPE_PEM) == WOLFSSL_SUCCESS);
    }
    if (rawPeer) {
        result = read_file(argv[5], data, sizeof(data));
        CHECK(result > 0);
        CHECK(wolfSSL_CTX_set_expected_rpk(ctx, data, result) == WOLFSSL_SUCCESS);
    }
    else CHECK(wolfSSL_CTX_load_verify_locations(ctx, argv[5], NULL) == WOLFSSL_SUCCESS);
    wolfSSL_CTX_set_verify(ctx, server && serverOnly ? WOLFSSL_VERIFY_NONE :
        WOLFSSL_VERIFY_PEER | (server ? WOLFSSL_VERIFY_FAIL_IF_NO_PEER_CERT : 0), NULL);
    if (!server && pha) CHECK(wolfSSL_CTX_allow_post_handshake_auth(ctx) == 0);
    if (!server && resume) CHECK(wolfSSL_CTX_UseSessionTicket(ctx) == WOLFSSL_SUCCESS);
    if (early && server) CHECK(wolfSSL_CTX_set_max_early_data(ctx, 256) >= 0);
    if (strcmp(mode, "sha384") == 0) CHECK(wolfSSL_CTX_set_cipher_list(ctx, "TLS13-AES256-GCM-SHA384") == WOLFSSL_SUCCESS);
    if (strcmp(mode, "ccm") == 0) CHECK(wolfSSL_CTX_set_cipher_list(ctx, "TLS13-AES128-CCM-SHA256") == WOLFSSL_SUCCESS);
    if (strcmp(mode, "chacha") == 0) CHECK(wolfSSL_CTX_set_cipher_list(ctx, "TLS13-CHACHA20-POLY1305-SHA256") == WOLFSSL_SUCCESS);
    if (ech) {
        if (server) {
            word32 size = sizeof(data);
            FILE* file;
            CHECK(wolfSSL_CTX_GenerateEchConfig(ctx, "public.test", 0x20, 1, 1) == WOLFSSL_SUCCESS);
            CHECK(wolfSSL_CTX_GetEchConfigs(ctx, data, &size) == WOLFSSL_SUCCESS);
            file = fopen(argv[7], "wb");
            CHECK(file != NULL);
            result = (int)fwrite(data, 1, size, file);
            fclose(file);
            CHECK(result == (int)size);
        }
        else {
            CHECK(wolfSSL_CTX_UseSNI(ctx, WOLFSSL_SNI_HOST_NAME, "server.test", 11) == WOLFSSL_SUCCESS);
        }
    }
    for (round = 0; round < (resume ? 2 : 1); round++) {
        fd = socket(AF_INET, SOCK_DGRAM, 0);
        CHECK(fd != INVALID_SOCKET);
        CHECK(setsockopt(fd, SOL_SOCKET, SO_RCVTIMEO, (const char*)&timeout, sizeof(timeout)) == 0);
        memset(&address, 0, sizeof(address));
        address.sin_family = AF_INET;
        address.sin_port = htons((unsigned short)port);
        address.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
        if (server) {
            socklen_t length = sizeof(address);
            address.sin_port = 0;
            CHECK(bind(fd, (struct sockaddr*)&address, sizeof(address)) == 0);
            CHECK(getsockname(fd, (struct sockaddr*)&address, &length) == 0);
            printf("READY %d %u\n", round, (unsigned int)ntohs(address.sin_port));
            CHECK(recvfrom(fd, (char*)data, sizeof(data), MSG_PEEK, (struct sockaddr*)&address, &length) > 0);
        }
        CHECK(connect(fd, (struct sockaddr*)&address, sizeof(address)) == 0);
        ssl = wolfSSL_new(ctx);
        CHECK(ssl != NULL);
        /* A resumed PSK needs no client certificate. Without PHA, suppress
         * wolfSSL's client_certificate_type response (RFC 7250 section 4.2). */
        if (server && round && !pha)
            wolfSSL_set_verify(ssl, WOLFSSL_VERIFY_NONE, NULL);
        if (ech && !server) {
            result = read_file(argv[7], data, sizeof(data));
            CHECK(result > 0);
            CHECK(wolfSSL_SetEchConfigs(ssl, data, result) == WOLFSSL_SUCCESS);
        }
        CHECK(wolfSSL_set_fd(ssl, (int)fd) == WOLFSSL_SUCCESS);
        CHECK(wolfSSL_dtls_set_peer(ssl, &address, sizeof(address)) == WOLFSSL_SUCCESS);
        CHECK(wolfSSL_dtls_set_mtu(ssl, strcmp(mode, "fragmented") == 0 ? 512 : 4096) == WOLFSSL_SUCCESS);
        if (server && (early || ech)) CHECK(wolfSSL_disable_hrr_cookie(ssl) == WOLFSSL_SUCCESS);
        if (early || strcmp(mode, "fragmented") == 0) {
            int group = WOLFSSL_ECC_SECP256R1;
            CHECK(wolfSSL_set_groups(ssl, &group, 1) == WOLFSSL_SUCCESS);
            if (!server) CHECK(wolfSSL_UseKeyShare(ssl, group) == WOLFSSL_SUCCESS);
        }
        if (strcmp(mode, "hrr") == 0) {
            int group = WOLFSSL_ECC_SECP256R1;
            if (server) CHECK(wolfSSL_set_groups(ssl, &group, 1) == WOLFSSL_SUCCESS);
            else CHECK(wolfSSL_UseKeyShare(ssl, WOLFSSL_ECC_X25519) == WOLFSSL_SUCCESS);
        }
        if (!server && round) CHECK(wolfSSL_set_session(ssl, session) == WOLFSSL_SUCCESS);
        if (early && round) {
            if (server) {
                CHECK(wolfSSL_read_early_data(ssl, data, sizeof(data), &count) > 0);
                CHECK(count == 5 && memcmp(data, "early", 5) == 0);
            }
            else CHECK(wolfSSL_write_early_data(ssl, "early", 5, &count) > 0 && count == 5);
        }
        result = server ? wolfSSL_accept(ssl) : wolfSSL_connect(ssl);
        if (ech) fprintf(stderr, "ECH status=%d\n", wolfSSL_GetEchStatus(ssl));
        CHECK(result == WOLFSSL_SUCCESS);
        CHECK(wolfSSL_session_reused(ssl) == (round != 0));
        if (ech) CHECK(wolfSSL_GetEchStatus(ssl) == WOLFSSL_ECH_STATUS_ACCEPTED);
        if (strcmp(mode, "keyupdate") == 0 && !server) CHECK(wolfSSL_update_keys(ssl) == WOLFSSL_SUCCESS);
        for (step = 0; step < 2; step++) {
            if (!server) CHECK(wolfSSL_write(ssl, "ping", 4) == 4);
            result = wolfSSL_read(ssl, data, sizeof(data));
            CHECK(result == 4 && memcmp(data, "ping", 4) == 0);
            if (server) {
                if (pha && step == 0) CHECK(wolfSSL_request_certificate(ssl) == WOLFSSL_SUCCESS);
                CHECK(wolfSSL_write(ssl, data, result) == result);
            }
        }
        if (early && round) CHECK(wolfSSL_get_early_data_status(ssl) == WOLFSSL_EARLY_DATA_ACCEPTED);
        if (!server && resume && round == 0) {
            session = wolfSSL_get1_session(ssl);
            CHECK(session != NULL);
        }
        printf("DONE %d resumed=%d\n", round, wolfSSL_session_reused(ssl));
        wolfSSL_free(ssl);
        ssl = NULL;
        closesocket(fd);
        fd = INVALID_SOCKET;
    }
    status = 0;
cleanup:
    if (session != NULL) wolfSSL_SESSION_free(session);
    if (ssl != NULL) wolfSSL_free(ssl);
    if (fd != INVALID_SOCKET) closesocket(fd);
    if (ctx != NULL) wolfSSL_CTX_free(ctx);
    wolfSSL_Cleanup();
#ifdef _WIN32
    WSACleanup();
#endif
    return status;
}
