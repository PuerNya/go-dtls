# Тестовые программы для узлов DTLS 1.3

[简体中文](README.md) | [English](README.en.md) | [Русский](README.ru.md)

Эти явно включаемые тесты используют собственные слои записей UDP неизменённых
сторонних библиотек. `openssl.c` служит для OpenSSL и BoringSSL, `nss.c` — для
NSS. Тест Go создаёт краткосрочные учётные данные, задаёт доверенные сертификаты,
сравнивает exporter и проверяет прикладные данные и согласованное состояние.
Если переменная окружения для узла не задана, его тесты пропускаются.

## Версии исходного кода

Матрица проверена на следующих точных ревизиях upstream:

| Библиотека | Репозиторий | Commit |
| --- | --- | --- |
| NSS | https://github.com/nss-dev/nss | `b9283c1e4ee232b4450d36115272af670d38ec81` |
| BoringSSL | https://github.com/google/boringssl | `4f3b183c5d4a8f8fcaba4c0ab4435df181df6978` |
| OpenSSL | https://github.com/openssl/openssl | `41eecc7c7cd4069b1b29bca49b9d97825c6bf648` |

Все три узла согласуют окончательную версию DTLS 1.3 из RFC 9147, `0xfefc`.
Успешный тест TLS 1.3 или одна лишь сборка не доказывают совместимость DTLS.

## Сборка и запуск в Linux

Установите нативный инструментарий C/C++, Git, Make, CMake, Perl, Python,
pkg-config и файлы разработки NSPR. Храните исходники и результаты сборки вне
этого репозитория. Из каталога go-dtls задайте `repo` и абсолютный путь `work`
к внешнему каталогу:

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

Завершите загрузки до запуска тестов и удалите унаследованные настройки прокси
из окружения тестовых процессов. Сам тест не загружает и не собирает узлы.
Существующие тесты wolfSSL по-прежнему называются `TestInteropWolfSSL*`;
сборка определена в [workflow](../../.github/workflows/benchmarks.yml)
независимо от того, запускаются ли benchmark.

## Покрытие и ограничения

Матрица охватывает обе роли узла — клиента и сервера: X.509, возобновление,
AES-GCM/SHA-256 и SHA-384, ChaCha20-Poly1305, X25519/P-256/P-384, HRR для
key share, фрагментацию, mTLS, OCSP, SNI, ALPN, GREASE, KeyUpdate и откат при
отсутствии сжатия сертификатов или cached information. Сценарии делегированных
учётных данных проверяют фактическое использование DC, возобновление,
фрагментацию, SHA-384, откат, неверные подписи и истечение срока. Отрицательные
проверки названы отдельно и не считаются успешными рукопожатиями.

Сценарии `client/early-io` для NSS и OpenSSL задерживают клиентский Finished
и требуют запрос эпохи 1 и ответ эпохи 3 до его освобождения. Первый вызов
отправки NSS неблокирующий, чтобы socket API не завершил рукопожатие до записи
запроса. Обычные проверки early data также проверяют `DatagramInfo.EarlyData`.
Матрица wolfSSL RPK проверяет `early-io` в обеих ролях на неизмененной ревизии
`c6e4286c0269a111799ff059badd94a272e37574`, собранной MSVC в Windows и штатным
компилятором в Linux. Клиент обрабатывает `APP_DATA_READY` при подключении.
Отсутствие DTLS 0-RTT в BoringSSL и серверные ограничения ниже сохраняются.
Сценарий HRR сервера wolfSSL проверяет отклонение предложения 0-RTT без отправки
данных эпохи 1: эта ревизия возвращала их даже после того, как клиент Go
сообщил об отклонении ранних данных вследствие HRR.

| Узел | Дополнительно проверенное поведение | Граница возможностей закреплённой ревизии |
| --- | --- | --- |
| NSS | Отправка DC сервером; проверка DC сервера Go клиентом; клиентский 0-RTT, включая сессию с аутентификацией DC | Отклоняет DC в CertificateRequest (`tls13con.c`, `KnownExtensions`) вопреки RFC 9345 §4.1.2. В обычных сценариях mTLS это предложение Go отключено. PHA для DTLS с собственным слоем записей не поддерживается. |
| BoringSSL | Отправка DC сервером | Проверка DC клиентом, PHA и DTLS 0-RTT недоступны. `tls13_client.cc` явно исключает DTLS 0-RTT. |
| OpenSSL | Оба направления PHA; клиентский 0-RTT | DC отсутствует. Сервер отклоняет допустимый начальный ACK с `unexpected_message`; положительные серверные сценарии намеренно отбрасывают ACK эпохи 2 и называются `server-ACKLoss`. `server/initial-ACK-rejection` сохраняет проверку отказа без изменения обмена. |

Для клиентского 0-RTT NSS задаётся `SSL_SUPPRESS_END_OF_EARLY_DATA`: RFC 9147 §5
исключает EndOfEarlyData. Без этой настройки NSS отправляет запрещённое
сообщение. Документация NSS указывает, что серверу нужен `SSL_RecordLayerData`
для подавления EndOfEarlyData; сервер с собственным слоем записей зависает
даже с этой настройкой, поэтому исключён из положительных сценариев 0-RTT.
Listener OpenSSL обрабатывает ClientHello до возврата соединения;
`SSL_read_early_data` затем отклоняет вызов, поскольку `SSL_in_before` имеет
значение false. Серверный путь 0-RTT OpenSSL здесь также не подтверждён.

Автоматическая отправка tickets DTLS в NSS работает и покрыта тестами;
отказ отдельного API `SSL_SendSessionTicket` для DTLS не означает отсутствия
возобновления. NSS аутентифицирует узел по точному совпадению с DER-сертификатом,
закреплённым для каждого теста; OpenSSL и BoringSSL проверяют тестовый УЦ и имя
хоста. Ни одна из тестовых программ не отключает аутентификацию узла.
