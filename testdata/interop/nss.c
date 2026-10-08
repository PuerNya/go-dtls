/* Test-only NSS peer. NoDB credentials are imported from per-test DER files. */
#include <nss.h>
#include <ssl.h>
#include <sslexp.h>
#include <sslproto.h>
#include <pk11pub.h>
#include <cert.h>
#include <keyhi.h>
#include <prerror.h>
#include <prnetdb.h>
#include <prio.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void check(int ok, const char *what) {
    if (!ok) { PRErrorCode code = PR_GetError(); fprintf(stderr, "%s: %s (%d)\n", what, PR_ErrorToName(code), code); exit(1); }
}
static SECItem read_item(const char *path) {
    FILE *f = fopen(path, "rb"); check(f != NULL, "open DER");
    SECItem item = {siBuffer, NULL, 0}; unsigned char data[65536];
    size_t n = fread(data, 1, sizeof(data), f); check(!ferror(f) && feof(f), "read DER"); fclose(f);
    check(SECITEM_AllocItem(NULL, &item, n) != NULL, "allocate DER"); memcpy(item.data, data, n); return item;
}
static SECKEYPrivateKey *read_key(PK11SlotInfo *slot, const char *path) {
    SECItem der = read_item(path); SECKEYPrivateKey *key = NULL;
    check(PK11_ImportDERPrivateKeyInfoAndReturnKey(slot, &der, NULL, NULL, PR_FALSE, PR_FALSE, KU_ALL, &key, NULL) == SECSuccess && key != NULL, "import key");
    SECITEM_FreeItem(&der, PR_FALSE); return key;
}
static SECStatus auth(void *arg, PRFileDesc *fd, PRBool checksig, PRBool server) {
    (void)checksig; (void)server;
    CERTCertificate *cert = SSL_PeerCertificate(fd);
    SECStatus result = cert != NULL && SECITEM_ItemsAreEqual(&cert->derCert, arg) ? SECSuccess : SECFailure;
    if (cert) CERT_DestroyCertificate(cert);
    return result;
}
static PRInt32 sni(PRFileDesc *fd, const SECItem *names, PRUint32 count, void *arg) {
    (void)fd; (void)arg;
    check(count == 1 && names[0].len == 11 && memcmp(names[0].data, "server.test", 11) == 0, "SNI");
    return 0;
}
struct credential { CERTCertificate *cert; SECKEYPrivateKey *key; };
static SECStatus client_auth(void *arg, PRFileDesc *fd, CERTDistNames *names, CERTCertificate **cert, SECKEYPrivateKey **key) {
    struct credential *c = arg; (void)fd; (void)names;
    *cert = CERT_DupCertificate(c->cert); *key = SECKEY_CopyPrivateKey(c->key); return SECSuccess;
}
int main(int argc, char **argv) {
    check(argc == 11, "role port mode cert.pem key.pem cert.der key.der dc.bin delegated.pem delegated.der");
    int server = strcmp(argv[1], "server") == 0; const char *mode = argv[3];
    int early = strcmp(mode, "early") == 0 || strcmp(mode, "dc-early") == 0;
    check(NSS_NoDB_Init(NULL) == SECSuccess, "initialize NSS");
    PK11SlotInfo *slot = PK11_GetInternalKeySlot(); check(slot != NULL, "key slot");
    SECItem der = read_item(argv[6]);
    struct credential c = {CERT_NewTempCertificate(CERT_GetDefaultCertDB(), &der, NULL, PR_FALSE, PR_TRUE), read_key(slot, argv[7])};
    SECITEM_FreeItem(&der, PR_FALSE); check(c.cert != NULL, "parse certificate");
    char trustPath[4096]; snprintf(trustPath, sizeof(trustPath), "%s.trust", argv[6]);
    SECItem trust = read_item(trustPath);
    SECItem response = {siBuffer, NULL, 0};
    if (strcmp(mode, "ocsp") == 0) {
        snprintf(trustPath, sizeof(trustPath), "%s.ocsp", argv[4]);
        response = read_item(trustPath);
    }
    SECItemArray responses = {&response, 1};
    if (server) check(SSL_ConfigServerSessionIDCache(256, 0, 0, NULL) == SECSuccess, "session cache");
    check(!early || !server, "NSS native server requires EndOfEarlyData");
    int rounds = (strcmp(mode, "resume") == 0 || strcmp(mode, "dc-resume") == 0 || early) ? 2 : 1;
    for (int round = 0; round < rounds; round++) {
        PRFileDesc *udp = PR_OpenUDPSocket(PR_AF_INET); check(udp != NULL, "socket");
        PRNetAddr address; check(PR_InitializeNetAddr(PR_IpAddrLoopback, (PRUint16)atoi(argv[2]), &address) == PR_SUCCESS, "address");
        PRIntervalTime timeout = PR_SecondsToInterval(8);
        if (server) {
            check(PR_Bind(udp, &address) == PR_SUCCESS, "bind"); printf("READY %d\n", round); fflush(stdout);
            unsigned char b; check(PR_RecvFrom(udp, &b, 1, PR_MSG_PEEK, &address, timeout) > 0, "first datagram");
        }
        check(PR_Connect(udp, &address, timeout) == PR_SUCCESS, "connect UDP");
        PRFileDesc *fd = DTLS_ImportFD(NULL, udp); check(fd != NULL, "DTLS fd");
        SSLVersionRange versions = {SSL_LIBRARY_VERSION_TLS_1_3, SSL_LIBRARY_VERSION_TLS_1_3};
        check(SSL_VersionRangeSet(fd, &versions) == SECSuccess, "DTLS 1.3 version");
        check(SSL_OptionSet(fd, SSL_ENABLE_SESSION_TICKETS, PR_TRUE) == SECSuccess, "session tickets");
        if (early) {
            check(SSL_OptionSet(fd, SSL_ENABLE_0RTT_DATA, PR_TRUE) == SECSuccess, "enable early data");
            check(SSL_OptionSet(fd, SSL_SUPPRESS_END_OF_EARLY_DATA, PR_TRUE) == SECSuccess, "DTLS omits EndOfEarlyData");
        }
        check(SSL_OptionSet(fd, SSL_ENABLE_POST_HANDSHAKE_AUTH, PR_TRUE) == SECSuccess, "enable PHA");
        check(SSL_OptionSet(fd, SSL_ENABLE_OCSP_STAPLING, PR_TRUE) == SECSuccess, "enable OCSP");
        check(SSL_OptionSet(fd, SSL_ENABLE_DELEGATED_CREDENTIALS, !server && strcmp(mode, "dc-fallback") == 0 ? PR_FALSE : PR_TRUE) == SECSuccess, "enable DC");
        check(SSL_AuthCertificateHook(fd, auth, &trust) == SECSuccess, "pinned certificate hook");
        SSLNamedGroup group = strcmp(mode, "p256") == 0 || strcmp(mode, "hrr") == 0 ? ssl_grp_ec_secp256r1 : strcmp(mode, "p384") == 0 ? ssl_grp_ec_secp384r1 : ssl_grp_ec_curve25519;
        SSLNamedGroup groups[] = {ssl_grp_ec_curve25519, ssl_grp_ec_secp256r1};
        check(SSL_NamedGroupConfig(fd, !server && strcmp(mode, "hrr") == 0 ? groups : &group, !server && strcmp(mode, "hrr") == 0 ? 2 : 1) == SECSuccess, "group");
        PRUint16 suite = (strcmp(mode, "sha384") == 0 || strcmp(mode, "dc-sha384") == 0) ? TLS_AES_256_GCM_SHA384 : strcmp(mode, "chacha") == 0 ? TLS_CHACHA20_POLY1305_SHA256 : TLS_AES_128_GCM_SHA256;
        for (unsigned i = 0; i < SSL_NumImplementedCiphers; i++) check(SSL_CipherPrefSet(fd, SSL_ImplementedCiphers[i], SSL_ImplementedCiphers[i] == suite) == SECSuccess, "cipher");
        if (strcmp(mode, "alpn") == 0) { static const unsigned char p[] = {4, 'c','o','a','p'}; check(SSL_SetNextProtoNego(fd, p, sizeof(p)) == SECSuccess, "ALPN"); }
        SECKEYPrivateKey *delegatedKey = NULL; SECItem dc = {siBuffer, NULL, 0};
        if (server) {
            check(SSL_SNISocketConfigHook(fd, sni, NULL) == SECSuccess, "SNI hook");
            SSLExtraServerCertData extra = {0};
            if (response.data) extra.stapledOCSPResponses = &responses;
            if (strncmp(mode, "dc", 2) == 0) {
                dc = read_item(argv[8]); delegatedKey = read_key(slot, argv[10]);
                extra.delegCred = &dc; extra.delegCredPrivKey = delegatedKey;
            }
            check(SSL_ConfigServerCert(fd, c.cert, c.key, &extra, sizeof(extra)) == SECSuccess, "server credential");
            if (strcmp(mode, "mutual") == 0) {
                check(SSL_OptionSet(fd, SSL_REQUEST_CERTIFICATE, PR_TRUE) == SECSuccess && SSL_OptionSet(fd, SSL_REQUIRE_CERTIFICATE, SSL_REQUIRE_ALWAYS) == SECSuccess, "client authentication");
            }
        } else {
            check(SSL_SetURL(fd, "server.test") == SECSuccess, "SNI");
            check(SSL_GetClientAuthDataHook(fd, client_auth, &c) == SECSuccess, "client credential");
        }
        check(SSL_ResetHandshake(fd, server) == SECSuccess, "reset handshake");
        if (early && round != 0) {
            check(PR_Send(fd, "early", 5, 0, timeout) == 5, "write early data");
        }
        check(SSL_ForceHandshakeWithTimeout(fd, timeout) == SECSuccess, "handshake");
        SSLChannelInfo info = {0}; check(SSL_GetChannelInfo(fd, &info, sizeof(info)) == SECSuccess && info.protocolVersion == SSL_LIBRARY_VERSION_TLS_1_3, "negotiated version");
        check(info.resumed == (round != 0), "session reuse");
        if (early && round != 0) check(info.earlyDataAccepted, "early data accepted");
        if (!server && strcmp(mode, "ocsp") == 0) {
            const SECItemArray *responses = SSL_PeerStapledOCSPResponses(fd);
            check(responses != NULL && responses->len == 1 && SECITEM_ItemsAreEqual(&responses->items[0], &response), "received OCSP");
        }
        if (strcmp(mode, "keyupdate") == 0) check(SSL_KeyUpdate(fd, PR_TRUE) == SECSuccess, "KeyUpdate");
        unsigned char exporter[32]; check(SSL_ExportKeyingMaterial(fd, "interop", 7, PR_FALSE, NULL, 0, exporter, sizeof(exporter)) == SECSuccess, "exporter");
        printf("VERSION=fefc CIPHER=%04x PEER_DC=%d EXPORTER=", info.cipherSuite, info.peerDelegCred);
        for (unsigned i = 0; i < sizeof(exporter); i++) printf("%02x", exporter[i]);
        puts(""); fflush(stdout);
        unsigned char buf[64]; PRInt32 n;
        if (!server) check(PR_Send(fd, "interop", 7, 0, timeout) == 7, "write datagram");
        n = PR_Recv(fd, buf, sizeof(buf), 0, timeout); check(n == 7 && memcmp(buf, "interop", 7) == 0, "read datagram");
        if (server) check(PR_Send(fd, buf, n, 0, timeout) == n, "echo datagram");
        n = PR_Recv(fd, buf, sizeof(buf), 0, timeout); check(n == 4 && memcmp(buf, "done", 4) == 0, "completion datagram");
        PR_Close(fd);
        if (delegatedKey) SECKEY_DestroyPrivateKey(delegatedKey);
        if (dc.data) SECITEM_FreeItem(&dc, PR_FALSE);
    }
    SSL_ClearSessionCache();
    if (server) check(SSL_ShutdownServerSessionIDCache() == SECSuccess, "shutdown session cache");
    CERT_DestroyCertificate(c.cert); SECKEY_DestroyPrivateKey(c.key);
    SECITEM_FreeItem(&trust, PR_FALSE);
    if (response.data) SECITEM_FreeItem(&response, PR_FALSE);
    PK11_FreeSlot(slot); check(NSS_Shutdown() == SECSuccess, "shutdown NSS"); return 0;
}
