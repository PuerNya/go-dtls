# go-dtls

[简体中文](README.md) | [English](README.en.md) | [Русский](README.ru.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/puernya/go-dtls.svg)](https://pkg.go.dev/github.com/puernya/go-dtls)
[![CI](https://github.com/puernya/go-dtls/actions/workflows/ci.yml/badge.svg)](https://github.com/puernya/go-dtls/actions/workflows/ci.yml)
[![Benchmarks](https://github.com/puernya/go-dtls/actions/workflows/benchmarks.yml/badge.svg?branch=master)](https://github.com/puernya/go-dtls/blob/benchmark-results/benchmark.ru.md)
[![License: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](../LICENSE)

`go-dtls` — библиотека DTLS 1.3, реализованная на Go. Поведение протокола соответствует [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147). Путь модуля — `github.com/puernya/go-dtls`, имя пакета — `dtls13`.

Реализация охватывает обязательную семантику RFC 9147 и все 11 прямых нормативных ссылок этого RFC. Семантика TLS 1.3 соответствует [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846). Обязательные и рекомендуемые требования для поддерживаемых возможностей реализованы; неподдерживаемые необязательные расширения кратко перечислены в разделе «Границы реализации».

> Статус поддержки RFC не является сертификатом соответствия от третьей стороны. Ограничения безопасности сертификатов из RFC 9325 применяются по умолчанию: конечный RSA-сертификат аутентификации сервера должен содержать ключ не короче 2048 бит, а сертификаты SHA-1/MD5 не могут обойти это правило через самоподпись, корни доверия, `InsecureSkipVerify` или возобновление сессии.

## Основная семантика

DTLS — протокол ненадежных дейтаграмм, а не байтовый поток TLS:

- `Conn` намеренно не реализует ни `net.Conn`, ни `net.PacketConn`.
- Каждый вызов `WriteDatagram` отправляет одну запись DTLS Application Data. Библиотека не фрагментирует, не упорядочивает и не повторяет прикладные дейтаграммы.
- Каждый вызов `ReadDatagram` потребляет одну аутентифицированную запись Application Data. Если буфер слишком мал, остаток отбрасывается, а `DatagramInfo.Truncated` сообщает об усечении.
- По умолчанию прикладная дейтаграмма, превышающая текущий MTU пути или предел записи RFC, завершается ошибкой `ErrDatagramTooLarge` без частичной отправки. `IgnorePathMTU` позволяет пропустить только первую проверку.
- Для сообщений рукопожатия по-прежнему применяются фрагментация RFC 9147, ACK, восстановление после потерь и экспоненциальная задержка. Эти механизмы надежности не меняют семантику прикладных дейтаграмм.
- Гибридный обмен ключами RFC 9954 позволяет явно включить три стандартные группы ECDHE-MLKEM; по умолчанию остаются X25519/P-256.
- ECH по RFC 9849 шифрует настоящий ClientHello; отклонение, HRR, возобновление и 0-RTT работают fail-closed.
- `Listener` принимает аутентифицированные ассоциации DTLS из UDP-сокета и возвращает строго типизированный `*dtls13.Conn`.

## Требования

| Компонент | Требование |
| --- | --- |
| Go | Go 1.26 или новее; точные версии Go, платформа и CPU автоматического запуска указаны в последнем отчете |
| Транспорт | `udp`, `udp4` или `udp6`; TCP не поддерживается |
| Race-тесты Windows | Скрипту репозитория нужны Zig 0.17 и рабочая цепочка CGO |
| Тесты совместимости с wolfSSL | Необязательно; переменная `WOLFSSL_ROOT` должна указывать на совместимый каталог исходников/сборки wolfSSL |

Установка:

```sh
go get github.com/puernya/go-dtls
```

Последняя часть пути импорта содержит дефис, поэтому в исходном коде Go используется объявленное имя пакета `dtls13`:

```go
import dtls13 "github.com/puernya/go-dtls"
```

## Быстрый старт

### Клиент

`Dial` создает UDP-соединение и завершает рукопожатие до возврата. Этот клиент проверяет сертификат сервера, отправляет одну дейтаграмму и читает один ответ:

```go
package main

import (
	"crypto/x509"
	"log"
	"os"
	"time"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	caPEM, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		log.Fatal("ca.pem does not contain a certificate")
	}

	conn, err := dtls13.Dial("udp", "127.0.0.1:4444", &dtls13.Config{
		RootCAs:    roots,
		ServerName: "server.example",
		NextProtos: []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.WriteDatagram([]byte("ping")); err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Fatal(err)
	}
	if info.Truncated {
		log.Fatalf("response needs %d bytes", info.FullLength)
	}
	log.Printf("received %q from %s", buffer[:n], info.Source)
}
```

### Сервер

`Listen` создает UDP-listener. Получив новую ассоциацию, `Accept` возвращает `*Conn`. Приложение может явно вызвать `Handshake` либо позволить первому `ReadDatagram` или `WriteDatagram` запустить рукопожатие.

```go
package main

import (
	"crypto/tls"
	"log"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := dtls13.Listen("udp", ":4444", &dtls13.Config{
		Certificates: []tls.Certificate{certificate},
		NextProtos:   []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go serve(conn)
	}
}

func serve(conn *dtls13.Conn) {
	defer conn.Close()
	if err := conn.Handshake(); err != nil {
		log.Printf("handshake: %v", err)
		return
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Printf("read: %v", err)
		return
	}
	if info.Truncated {
		log.Printf("discarded truncated datagram: need %d bytes", info.FullLength)
		return
	}
	if _, err = conn.WriteDatagram(buffer[:n]); err != nil {
		log.Printf("write: %v", err)
	}
}
```

## Основные API

| API | Назначение |
| --- | --- |
| `Dial` / `DialWithDialer` | Создает клиентское UDP-соединение и завершает рукопожатие; `DialWithDialer` настраивает локальный адрес и тайм-аут подключения |
| `Listen` | Создает UDP-listener и владеет им |
| `NewListener` | Создает DTLS-listener поверх существующего `net.PacketConn`; ошибки конфигурации возвращает `Accept` |
| `Listener.Accept` | Принимает новую ассоциацию и возвращает `*Conn` |
| `Client` / `Server` | Оборачивает существующий подключенный `net.Conn` в клиент или сервер DTLS; рукопожатие можно выполнить явно до первого ввода-вывода |
| `Handshake` / `HandshakeContext` | Явно выполняет рукопожатие; последующие вызовы возвращают тот же результат |
| `ReadDatagram` | Читает одну аутентифицированную дейтаграмму и возвращает источник, полную длину и признак усечения |
| `WriteDatagram` | Отправляет одну дейтаграмму аутентифицированному узлу ассоциации |
| `SetDeadline` / `SetReadDeadline` / `SetWriteDeadline` | Устанавливает крайние сроки для ввода-вывода нижележащих дейтаграмм |
| `ConnectionState` | Возвращает версию, набор шифров, ALPN, сертификаты, состояние возобновления, identity/context внешнего PSK, активные CID, состояние RRC и согласованные ограничения RFC 8449 в обоих направлениях |
| `Close` | Отправляет `close_notify`, очищает секреты трафика/возобновления и закрывает нижележащее соединение |

### Размер дейтаграммы и усечение

Перед записью приложение должно использовать `PathMTU` и `RecordOverhead`, чтобы оценить текущую максимальную полезную нагрузку и учесть возможное уменьшение PMTU в течение жизни соединения:

```go
maximum := conn.PathMTU() - conn.RecordOverhead()
if len(payload) > maximum {
	// Фрагментация относится к прикладному протоколу; каждый фрагмент остается отдельной дейтаграммой.
}
```

Состояние пути может измениться после успешной предварительной проверки, поэтому ошибки записи все равно необходимо обрабатывать:

```go
if _, err := conn.WriteDatagram(payload); errors.Is(err, dtls13.ErrDatagramTooLarge) {
	// Уменьшите или разделите payload по правилам прикладного протокола и повторите как новую дейтаграмму.
}
```

Установите `IgnorePathMTU: true`, если приложение активно определяет PMTU или явно полагается на IP-фрагментацию. Тогда `WriteDatagram` и `WriteEarlyData` пропускают внутреннюю проверку полезной нагрузки по PMTU и передают одну полную запись DTLS нижележащему транспорту. Для flight рукопожатия, ACK и post-handshake по-прежнему применяются фрагментация по PMTU, повторная передача и задержка. Параметр не увеличивает предел содержимого записи `2^14` байт или согласованный `record_size_limit`. Транспорт все еще может фрагментировать, терять или возвращать `ErrDatagramTooLarge`; библиотека не уменьшает PMTU и не повторяет прикладную дейтаграмму автоматически.

Короткий буфер `ReadDatagram` не превращает чтение в потоковое. При `Truncated=true` непрочитанная часть записи уже отброшена, и следующий вызов вернет следующую запись.

### Модель ошибок

| Ошибка | Значение |
| --- | --- |
| `*ConfigError` | Недопустимая локальная конфигурация: например, неверный транспорт, MTU, набор шифров или ограничение ресурсов |
| `*ProtocolError` | Полученные или созданные данные нарушают состояние или формат протокола |
| `AlertError` | Узел вернул фатальное предупреждение TLS; числовое описание предупреждения доступно через `errors.As` |
| `ErrDatagramTooLarge` | Прикладная дейтаграмма превышает текущий PMTU или предел записи; транспорт может вернуть эту ошибку и при `IgnorePathMTU`; проверяйте через `errors.Is` |
| `ErrEarlyDataUnavailable` | Нет подходящего ticket для early data либо соединение не может отправить 0-RTT |
| `ErrEarlyDataRejected` | Рукопожатие завершено, но узел отклонил отправленный 0-RTT из-за HRR, повтора или политики |
| `*ECHRejectionError` | Сервер отклонил ECH после аутентификации соединения `public_name`; ошибка может содержать `RetryConfigList`, допустимый только для того же источника конфигурации и endpoint |
| `io.EOF` | Узел отправил корректный `close_notify`, закрыв направление чтения |

Крайние сроки, закрытие сокета и ошибки нижележащего UDP следуют модели ошибок пакета Go `net`. Вызывающий код не должен зависеть от текста ошибки.

## Расширенные возможности

| Возможность | API / конфигурация | Описание |
| --- | --- | --- |
| Внешний PSK / importer | `ImportExternalPSK`, `NewDirectExternalPSK`, `ExternalPSKs` | Аутентификация без сертификата по RFC 9257/9258; рекомендуется importer, используется только `psk_dhe_ke`, поддерживаются несколько identity, HRR и возобновление через ticket |
| Возобновление сессии | `ClientSessionCache`, `NewLRUClientSessionCache` | Клиент кэширует NewSessionTicket; сервер управляется `SessionTicketKey` и настройками ticket; возобновление mTLS сохраняет состояние аутентификации клиента |
| 0-RTT | `WriteEarlyData`, `MaxEarlyData`, `EarlyDataReplayCache` | Доступно только для возобновленных соединений; вызывающий код обязан обработать `ErrEarlyDataUnavailable` и `ErrEarlyDataRejected`, а early data должны допускать повтор |
| KeyUpdate | `SendKeyUpdate(requestPeer)` | Надежно отправляется, эпоха отправки переключается после ACK; также запускается автоматически при приближении к пределу использования AEAD |
| CID / проверка пути | `ConnectionID`, `GetConnectionID`, `SendNewConnectionIDs`, `RequestConnectionIDs`, `UseNextConnectionID` | Поддерживает согласование и обновление CID по RFC 9146 и по умолчанию согласует RRC по RFC 9853; Listener меняет привязку только после проверки нового пути |
| Сжатие сертификатов | `EnableCertificateCompression` | Явно включает zlib по RFC 8879 для сертификатов сервера и клиентских сертификатов mTLS/PHA; если сжатое сообщение не меньше, отправляется обычный Certificate |
| Динамический выбор сертификата | `GetCertificate`, `GetClientCertificate`, `ServerCertificateAuthorities`, `ClientCertificateOIDFilters` | Подсказки УЦ/OID по RFC 9846 и выбор из нескольких сертификатов на обеих сторонах; начальный mTLS и PHA используют одинаковые правила |
| Encrypted ClientHello | `EncryptedClientHelloConfigList`, `EncryptedClientHelloKeys`, `EncryptedClientHelloGrease` | Inner/Outer ClientHello, HPKE, HRR, подтверждение принятия, retry-конфигурации, возобновление и 0-RTT по RFC 9849 |
| Аутентификация клиента в рукопожатии | `ClientAuth`, `ClientCAs`, `Certificates` | Использует политики клиентских сертификатов из `crypto/tls` |
| Аутентификация клиента после рукопожатия | `PostHandshakeAuth`, `RequestClientCertificate` | Сначала клиент объявляет поддержку, затем сервер инициирует PHA |
| Exporter | `ConnectionState().ExportKeyingMaterial` | Экспортирует материал RFC 8446, раздел 7.5, с меткой DTLS `dtls13` |

### Смена адреса с CID

Клиент, предлагающий CID, по умолчанию также предлагает расширение `rrc` из RFC 9853. Сервер включает проверку пути, только если согласованы и CID, и RRC. Получив от нового источника аутентифицированную запись, маршрутизированную по CID, Listener выполняет enhanced check: сначала посылает challenge старому адресу; если старый путь доступен, сохраняет текущую привязку; после `path_drop` или тайм-аута проверяет адрес-кандидат. `RemoteAddr` и маршрутизация Listener по tuple атомарно обновляются только после правильного возврата случайного cookie кандидатом. Во время проверки прикладные записи продолжают отправляться на старый адрес.

Объем трафика к пути-кандидату не превышает утроенного числа корректных байт, полученных с этого адреса. При наличии измерения RTT таймер равен `3xRTT`, иначе — одной секунде. Каждый challenge использует новый cookie из CSPRNG; неизвестные типы сообщений RRC и недопустимые response/drop молча отбрасываются. Если узел предоставил резервный CID, challenge пути-кандидата временно использует его, чтобы не повторять старый CID на разных путях. Резервный CID активируется только после проверки; во время проверки прикладной трафик старого пути сохраняет исходный CID.

Приложения с эквивалентной проверкой адреса могут установить `DisableReturnRoutabilityCheck: true`. Параметр отключает только RRC, но не CID. `Dial` использует подключенный UDP, и операционные системы обычно не доставляют такому сокету пакеты от сервера с изменившимся адресом источника. Поэтому автоматическая смена привязки в основном относится к ассоциациям Listener, способным маршрутизировать разные источники по CID. Пустой CID может согласовать RRC, но не позволяет однозначно демультиплексировать разные пятиэлементные кортежи и потому не поддерживает миграцию адреса Listener.

## Основная конфигурация

В полях с одинаковой семантикой TLS 1.3 тип `Config` следует `crypto/tls.Config`. Конфигурацию можно использовать повторно, но нельзя изменять после первого использования; для производной конфигурации вызовите `Clone`.

| Параметр | Значение по умолчанию / поведение |
| --- | --- |
| `Certificates` / `GetCertificate` | Сертификаты сервера; по умолчанию выбирается первый сертификат, совместимый с SNI, алгоритмами подписи и подсказкой УЦ клиента; конечный RSA-ключ должен быть не короче 2048 бит, а вся цепочка не должна использовать SHA-1/MD5 |
| `GetClientCertificate` | Необязательный callback выбора клиентского сертификата; при `nil` выбирается первая запись `Certificates`, совместимая с алгоритмами подписи, именами УЦ и распознанными фильтрами OID |
| `ServerCertificateAuthorities` | По умолчанию пусто; DER-кодированные имена X.501 УЦ, отправляемые клиентом в ClientHello для выбора сертификата сервера; `RootCAs` автоматически не раскрывается |
| `RootCAs` / `ServerName` | Проверка сертификата сервера на клиенте; если `ServerName` не задан, `Dial` использует имя целевого узла |
| `ClientCAs` / `ClientAuth` | Политика проверки клиентского сертификата на сервере |
| `ClientCertificateOIDFilters` | По умолчанию пусто; фильтры OID RFC 9846, отправляемые сервером в CertificateRequest; Key Usage/EKU участвуют в выборе клиентом и проверке сервером |
| `VerifyPeerCertificate` | Дополнительная проверка после стандартной обработки сертификатов полного рукопожатия; как и в `crypto/tls`, при возобновлении повторно не вызывается |
| `InsecureSkipVerify` | По умолчанию `false`; производственные приложения не должны полагаться на него для обхода проверки идентичности |
| `NextProtos` | Список протоколов ALPN |
| `CipherSuites` | AES-128-GCM, AES-256-GCM, ChaCha20-Poly1305, AES-128-CCM |
| `CurvePreferences` | По умолчанию X25519 и P-256; явно поддерживаются `tls.X25519MLKEM768`, `tls.SecP256r1MLKEM768` и `tls.SecP384r1MLKEM1024` |
| `ExternalPSKs` | По умолчанию пусто; принимает неизменяемые внешние PSK от `ImportExternalPSK` или `NewDirectExternalPSK`; несовместимо с `ClientAuth` |
| `EnableGREASE` | По умолчанию `false`; отправляет одно случайное пустое расширение RFC 8701 в клиентском CH и серверных CR/NST; оно не согласуется и не меняется после HRR и клонирования ECH Inner/Outer |
| `EncryptedClientHelloConfigList` | По умолчанию `nil`; клиент передает полный ECHConfigList RFC 9849 с двухбайтовой длиной; ненулевое значение делает принятие ECH обязательным |
| `EncryptedClientHelloRejectionVerify` | Необязательная замена встроенной проверки отклонения ECH через `RootCAs` и `public_name` |
| `EncryptedClientHelloKeys` / `GetEncryptedClientHelloKeys` | Серверные ECHConfig и ключи HPKE; хотя бы один ключ должен иметь `SendAsRetry`, callback выполняется до выбора SNI, ALPN и сертификата |
| `EncryptedClientHelloGrease` | По умолчанию `false`; отправляет GREASE ECH без настоящей конфигурации, отклонение не прерывает обычное соединение |
| `MTU` | Полезная нагрузка UDP 1200 байт; минимум 256 |
| `IgnorePathMTU` | По умолчанию `false`; только Application Data пропускает внутреннюю проверку PMTU, рукопожатие не меняется |
| `RecordSizeLimit` | `0` выбирает значение по умолчанию `2^14+1`; диапазон `64..2^14+1` задает максимальный полный `DTLSInnerPlaintext`, принимаемый узлом, и объявляется по RFC 8449 независимо от PMTU |
| `EnableCertificateCompression` | По умолчанию `false`; включает стандартный zlib по RFC 8879 и отправляет `CompressedCertificate`, только если узел предложил zlib и полное сжатое сообщение меньше; выход ограничен `MaxHandshakeMessage` |
| `FlightInterval` | Начальный интервал повторной передачи рукопожатия — одна секунда |
| `MaxFlightInterval` | Максимальная экспоненциальная задержка — 60 секунд |
| `HandshakeTimeout` | 30 секунд |
| `ReplayWindow` | 64 записи на эпоху |
| `MaxHandshakeMessage` | 1 MiB, настраивается до `2^24-1` |
| `MaxBufferedApplicationData` | 1 MiB |
| `MaxBufferedApplicationDatagrams` | 1024 дейтаграммы |
| `MaxPendingConnections` | 128 сессий Listener |
| `MaxSessionQueueDatagrams` | 64 дейтаграммы на сессию Listener |
| `SessionTicketRequest` | По умолчанию отключен; клиент явно включает RFC 9149 и отдельно запрашивает число tickets после полного и возобновленного рукопожатия |
| `MaxSessionTickets` | 4; ограничивает число tickets, выпускаемых сервером по одному запросу RFC 9149 |
| `SessionTicketLifetime` | 24 часа, максимум 7 дней |
| `MaxEarlyData` | 0, то есть 0-RTT по умолчанию выключен |
| `MaxConnectionIDs` | 8 CID в каждом направлении |
| `DisableReturnRoutabilityCheck` | По умолчанию `false`; RRC отключается только при эквивалентной проверке адреса приложением |

### Постквантовый гибридный обмен ключами

Гибридный обмен ключами RFC 9954 использует реализацию ML-KEM из стандартной библиотеки Go и поддерживает три группы IANA с `DTLS-OK=Y`. `tls.X25519MLKEM768` — рекомендуемый универсальный выбор. По умолчанию гибридный обмен не включен:

```go
CurvePreferences: []tls.CurveID{
	tls.X25519MLKEM768,
	tls.X25519,
},
```

Если указана и соответствующая традиционная группа, первый ClientHello повторно использует компонент ECDH гибридной группы для fallback share. После HRR клиент отправляет только запрошенную группу. Большие key share используют обычные механизмы фрагментации, ACK и повторной передачи рукопожатия DTLS.

### Encrypted ClientHello

Клиент получает доверенный параметр `ech` из DNS SVCB/HTTPS и декодирует его представление Base64 по RFC 9848. Библиотека принимает полный ECHConfigList в wire-формате и не выполняет DNS-запрос:

```go
echConfigList, err := base64.StdEncoding.DecodeString(dnsECHParameter)
if err != nil {
	log.Fatal(err)
}

clientConfig := &dtls13.Config{
	RootCAs:                        roots,
	ServerName:                     "private.example",
	EncryptedClientHelloConfigList: echConfigList,
}
```

Сервер устанавливает один или несколько ECHConfig и соответствующие закрытые ключи HPKE, созданные средствами развертывания. `Config` содержит один ECHConfig без внешней длины ECHConfigList; `PrivateKey` использует формат закрытого ключа соответствующего KEM из `crypto/hpke`:

```go
serverConfig := &dtls13.Config{
	Certificates: []tls.Certificate{certificate},
	EncryptedClientHelloKeys: []dtls13.EncryptedClientHelloKey{{
		Config:      echConfig,
		PrivateKey:  echPrivateKey,
		SendAsRetry: true,
	}},
}
```

При настроенном настоящем ECH клиент завершается успешно только после проверки подтверждения принятия в HRR или ServerHello. Аутентифицированное отклонение возвращает `*ECHRejectionError` и проверяет внешнее соединение по ECH `public_name`, даже если задан `InsecureSkipVerify`. Обычный `VerifyPeerCertificate` для отклоненного соединения не вызывается, клиентский сертификат не отправляется. `RetryConfigList` разрешено использовать только с тем же источником DNS-конфигурации и transport endpoint. Для собственной проверки public name задайте `EncryptedClientHelloRejectionVerify`. `ConnectionState().ECHAccepted` сообщает о принятии настоящего ECH; GREASE его не устанавливает.

### Внешние PSK и importer

Рекомендуемая точка входа — importer RFC 9258. Он привязывает EPSK к DTLS 1.3 меткой `dtls13` и выводит отдельные целевые ключи SHA-256 и SHA-384. Возвращаемое значение не хранит исходный EPSK:

```go
psk, err := dtls13.ImportExternalPSK(
	[]byte("device-17"),
	provisionedKey, // Не менее 16 байт, желательно не менее 128 бит энтропии.
	[]byte("client=device-17;server=gateway-2"),
	crypto.SHA256, // Если hash не связан с EPSK, передайте 0; по умолчанию используется SHA-256.
)
if err != nil {
	log.Fatal(err)
}

config := &dtls13.Config{ExternalPSKs: []*dtls13.ExternalPSK{psk}}
```

Существующие развертывания с явно специализированным для TLS ключом могут использовать `NewDirectExternalPSK(identity, key, hash)`. Прямой PSK использует `ext binder`, импортированный — `imp binder`; эти формы невзаимозаменяемы. Клиент может настроить несколько identity. Сервер выбирает первую известную identity, совместимую с hash выбранного набора шифров, или откатывается к сертификатной аутентификации при наличии сертификата. Обе формы предлагают только `psk_dhe_ke`. Для первого рукопожатия с внешним PSK `DidResume` равен `false`; выданный затем ticket возобновляется обычным способом и сохраняет источник аутентификации через `ConnectionState.ExternalPSKIdentity()` и `ExternalPSKContext()`. Удаление или изменение внешнего PSK делает производные tickets недействительными.

Identity и context importer передаются открытым текстом в ClientHello, поэтому повторное использование позволяет связывать соединения, и эти поля не должны содержать секреты. PSK следует выдавать фиксированной паре ролей клиента и сервера. Для группового ключа context должен связывать identity обеих сторон и канал вышестоящего provisioning. Базовый TLS 1.3 не объединяет внешний PSK с сертификатной аутентификацией, поэтому `ClientAuth` нельзя включать вместе с `ExternalPSKs`. Сам внешний PSK не отправляет 0-RTT; только последующее возобновление через ticket может использовать обычную политику `MaxEarlyData` и replay cache.

### Сжатие сертификатов

При `EnableCertificateCompression: true` клиент предлагает zlib в ClientHello, позволяя серверу сжать свой сертификат. Сервер предлагает zlib в CertificateRequest, позволяя включившему параметр клиенту сжать сертификат mTLS или PHA. Для сжатия сертификатов в обоих направлениях включите параметр на обоих узлах.

Реализация использует только zlib из стандартной библиотеки Go. `CompressedCertificate` отправляется лишь тогда, когда полное сообщение меньше обычного Certificate; иначе выполняется безопасный откат. Семантика фрагментации рукопожатия, ACK, повторной передачи, HRR, возобновления и `record_size_limit` не меняется. Заявленная несжатая длина и фактический выход ограничены `MaxHandshakeMessage`.

### Динамический выбор сертификата

Если любая сторона настроила несколько записей в `Certificates`, стандартный выбор берет первый сертификат, совместимый с алгоритмами подписи узла, алгоритмами цепочки и подсказками УЦ; сервер также проверяет SNI. Сервер отправляет subjects из `ClientCAs` в начальном и post-handshake CertificateRequest. Клиент может отдельно задать `ServerCertificateAuthorities` для подсказки серверу в ClientHello; `RootCAs` автоматически не передается.

Пользовательская политика сервера использует `GetCertificate` и `ClientHelloInfo.SupportsCertificate`, политика клиента — `GetClientCertificate` и `CertificateRequestInfo.SupportsCertificate`. Срезы callback действительны только во время вызова: их нельзя изменять или сохранять. Элемент `ClientCertificateOIDFilters` содержит ASN.1 OID и DER-значение расширения без оболочки X.509 extension. Распознаются Key Usage и Extended Key Usage; неизвестные OID остаются в wire и игнорируются согласно RFC 9846. Начальный mTLS и PHA используют одинаковые правила. В возобновленном рукопожатии CertificateRequest отсутствует, поэтому клиентский сертификат не выбирается повторно и callback клиента не вызывается.

### Быстрое возобновление mTLS

Клиент и сервер используют обычные настройки mTLS; дополнительно включите кэш сессий клиента и ticket сервера:

```go
clientConfig := &dtls13.Config{
	Certificates:       []tls.Certificate{clientCertificate},
	RootCAs:            serverRoots,
	ServerName:         "server.example",
	ClientSessionCache: dtls13.NewLRUClientSessionCache(64),
}
serverConfig := &dtls13.Config{
	Certificates:          []tls.Certificate{serverCertificate},
	ClientAuth:            tls.RequireAndVerifyClientCert,
	ClientCAs:             clientRoots,
	SessionTicketKey:      ticketKey,
	SessionTicketLifetime: time.Hour,
}
```

Первое соединение выполняет полную аутентификацию клиента `CertificateRequest -> Certificate -> CertificateVerify`. Последующие соединения возобновляются через PSK-рукопожатие RFC 9147/RFC 9846 без повторной отправки сертификатов. Сервер восстанавливает сертификаты клиента и проверенные цепочки из аутентифицированного и зашифрованного AES-256-GCM ticket, затем повторно оценивает их по активным `ClientAuth`, `ClientCAs` и срокам действия сертификатов. Если политика не выполнена, ticket игнорируется и выполняется полное рукопожатие.

Ticket без идентичности клиента используется для возобновления только при `ClientAuth == tls.NoClientCert`. Любая политика клиентских сертификатов приводит к полному рукопожатию, предотвращая возобновление анонимной сессии.

Обновленный ticket сохраняет время исходного интерактивного `CertificateVerify`. `SessionTicketLifetime` ограничивает и срок ticket, и общий срок этой аутентификации клиента. `VerifyPeerCertificate` не запускается повторно при возобновлении; при изменении политики идентичности приложения смените `SessionTicketKey` или отключите session tickets. Явные ключи ticket необходимо регулярно ротировать. Конфигурация содержит один активный ключ, поэтому его смена немедленно делает старые tickets недействительными. PHA после рукопожатия не выпускает дополнительный ticket автоматически; если идентичность PHA должна использоваться при последующем возобновлении, установите новое полное соединение mTLS.

### Запросы билетов сеанса

RFC 9149 позволяет клиенту отдельно запросить число tickets после полного и возобновленного рукопожатия. По умолчанию расширение отключено. При `Enabled: true` даже два нулевых счетчика отправляют расширение и явно запрашивают отсутствие tickets для этого соединения:

```go
clientConfig.SessionTicketRequest = dtls13.SessionTicketRequest{
	Enabled:         true,
	NewSessionCount: 4,
	ResumptionCount: 1,
}
serverConfig.MaxSessionTickets = 4
```

Сервер возвращает ожидаемое число и использует меньшее из запроса и `MaxSessionTickets`. Нулевой предел сервера означает значение по умолчанию 4, а `SessionTicketsDisabled` сохраняет приоритет. Без расширения сохраняется прежнее поведение с одним ticket. Несколько NewSessionTicket передаются одним надежным flight с отдельными nonce, PSK, tickets и последовательными значениями `message_seq`, используя обычные ACK и повторную передачу DTLS.

Встроенный `NewLRUClientSessionCache` хранит не более 255 запрошенных tickets на один ключ сервера и атомарно расходует разные identity для параллельных соединений. Если возобновление отклонено и аутентифицированное рукопожатие завершено, связанные tickets того же семейства аннулируются. Пользовательские реализации `ClientSessionCache` сохраняют прежнюю семантику одного состояния; необязательный атомарный метод `Take(string)` продолжает использоваться.

Поддерживаемые наборы шифров:

| Константа | ID | Статус |
| --- | --- | --- |
| `TLS_AES_128_GCM_SHA256` | `0x1301` | Поддерживается |
| `TLS_AES_256_GCM_SHA384` | `0x1302` | Поддерживается |
| `TLS_CHACHA20_POLY1305_SHA256` | `0x1303` | Поддерживается |
| `TLS_AES_128_CCM_SHA256` | `0x1304` | Поддерживается |
| `TLS_AES_128_CCM_8_SHA256` | `0x1305` | Явно отключен: библиотека общего назначения не может гарантировать дополнительные меры против подделки на уровне развертывания, требуемые RFC 9147 |

## Полнота RFC 9147

Нормативные ключевые слова интерпретируются по BCP 14. `MUST`, `MUST NOT`, `REQUIRED`, `SHALL` и `SHALL NOT` строго применяются клиентом и сервером при отправке и приеме. Клиент активно реализует требования класса `SHOULD`; сервер допускает отклонения узла только без ослабления аутентификации, конфиденциальности, защиты от повторов, ограничений усиления и согласованности состояния. После согласования возможности `MAY` или `OPTIONAL` все ее условные обязательные требования применяются полностью.

### Общий статус

| Спецификация | Статус | Реализация |
| --- | --- | --- |
| [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147) | Реализовано | Record, Handshake, эпохи, ACK, KeyUpdate, обновление CID, Application Data и применимые требования безопасности; рекомендуемое поведение реализовано для включенных возможностей |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | Реализовано | Согласование CID, направленные CID, обновления, маршрутизация Listener, обработка ошибок и сохранение адреса; детали только для DTLS 1.2 неприменимы |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | Реализовано | Согласование CH/EE по умолчанию, направленные ограничения, минимум 64, фатальный `record_overflow`, HRR, возобновление, 0-RTT, KeyUpdate, ACK и независимость от PMTU |
| [RFC 8879](https://www.rfc-editor.org/rfc/rfc8879) | Реализовано | Явно включаемый zlib; направленное согласование ClientHello/CertificateRequest, сертификаты сервера и клиента mTLS/PHA, transcript CompressedCertificate, безопасный откат и ограничения распаковки |
| [RFC 9149](https://www.rfc-editor.org/rfc/rfc9149) | Реализовано | Явно включаемый `ticket_request(58)`, счетчики полного/возобновленного соединения, инварианты HRR, предел сервера, надежные множественные NST, однократное параллельное расходование cache и аннулирование связанных tickets |
| [RFC 9257](https://www.rfc-editor.org/rfc/rfc9257) | Реализовано | Внешние PSK не короче 128 бит, только DHE, opaque identity, несколько identity, откат к сертификату, рекомендации по приватности и требования развертывания к парам ролей |
| [RFC 9258](https://www.rfc-editor.org/rfc/rfc9258) | Реализовано | Реализованы `ImportedIdentity`, DTLS `0xfefc`, целевые KDF SHA-256/384, исходный hash EPSK, `dtls13derived psk` и `imp binder` |
| [RFC 9848](https://www.rfc-editor.org/rfc/rfc9848) | Реализовано/интеграция приложения | Принимается полный ECHConfigList, декодированный из DNS-параметра `ech`; запрос SVCB/HTTPS и декодирование Base64 выполняет приложение |
| [RFC 9849](https://www.rfc-editor.org/rfc/rfc9849) | Реализовано | Реализованы HPKE, Inner/Outer ClientHello, padding, HRR, подтверждение принятия, retry-конфигурации, аутентифицированное отклонение, GREASE, возобновление и 0-RTT |
| [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954) | Реализовано | Три стандартные группы ECDHE-MLKEM, fallback на традиционный share, HRR, фрагментация, возобновление mTLS, 0-RTT, ECH и строгая семантика ошибок; по умолчанию отключено |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | Реализовано для включенных возможностей | Реализованы подсказки УЦ/OID и выбор из нескольких сертификатов на обеих сторонах; `user_canceled(90)` игнорируется с продолжением ожидания `close_notify` во время рукопожатия, ожидания финального ACK и post-handshake; локальная криптографическая ошибка без более точного предупреждения отправляет `general_error(117)`, а конкретное предупреждение протокола всегда имеет приоритет |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | Частично | Покрыты PFS, AEAD, SNI/ALPN, tickets, 0-RTT, KeyUpdate и ограничения сертификатов; OCSP stapling отсутствует, а модуль намеренно не реализует поддержку DTLS 1.2, требуемую этим BCP |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | Частично | Go X.509 и `ServerName` покрывают DNS-ID/IP-ID; URI-ID, SRV-ID и прикладные service identity делегированы callback-функциям проверки вызывающего кода |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | Реализовано | Расширение 61, защищенный content type 27, все три сообщения, неизвестные типы, enhanced/basic state machine, трехкратное ограничение усиления, таймер одна секунда/3xRTT, NAT rebinding, защита от off-path атак и межпутевая приватность с резервным CID |

### Покрытие разделов

| Раздел RFC 9147 | Статус | Краткое описание реализации |
| --- | --- | --- |
| разделы 1-2 Введение и термины | Неприменимо | Нормативные ключевые слова следуют RFC 2119/8174; отдельной функции wire format нет |
| раздел 3 Цели проектирования | Реализовано | Потери, изменение порядка, дублирование, задержка, фрагментация, защита от повторов, динамический PMTU и освобождение ресурсов |
| раздел 4 Record Layer | Реализовано | DTLSPlaintext, unified header, усеченные номера последовательности, эпохи, AEAD, защита номера последовательности, anti-replay, демультиплексирование CID и пределы использования |
| раздел 5 Handshake | Реализовано | Рукопожатие TLS 1.3, HRR cookie, аутентификация, фрагментация/сборка, flights, ACK, повторная передача, тайм-аут и новая ассоциация на том же пятиэлементном кортеже |
| раздел 6 Epochs | Реализовано | Эпохи 0/1/2/3/4+, 0-RTT, KeyUpdate, ограниченное хранение и очистка старых ключей |
| раздел 7 ACK | Реализовано | Content type 26, пустой ACK, немедленный ACK, частичный ACK, скользящее окно и надежный post-handshake ACK |
| раздел 8 KeyUpdate | Реализовано | Обновления после ACK, повторная передача, `update_requested`, хранение старой эпохи и обработка пределов |
| раздел 9 CID Update | Реализовано | New/RequestConnectionId, динамическая маршрутизация, ACK обновления, ограничения ресурсов и проверка prefix-free |
| раздел 10 Application Data | Реализовано | API подключенных дейтаграмм, границы сообщений, отсутствие упорядочивания/повторов, явное усечение и крайние сроки |
| раздел 11 Security | Реализовано | Ротация cookie, трехкратные ограничения усиления для рукопожатия и пути-кандидата RRC, anti-replay, пределы AEAD, обновление адреса после проверки и ограниченное состояние |
| раздел 12 Отличия DTLS 1.2 | Реализовано | Все применимые отличия включены в record, handshake, epoch, ACK и CID |
| раздел 13 Обновления DTLS 1.2 | Неприменимо | Модуль реализует только DTLS 1.3 |
| раздел 14 IANA | Неприменимо | Используются назначенные значения; библиотека не выполняет операции с реестром |

### Прямые нормативные ссылки

Таблица содержит все 11 `Normative References` из XML RFC 9147, опубликованного RFC Editor. Протоколы нижних уровней реализуются Go и операционной системой; библиотека не переопределяет стек UDP/TCP/IP.

| Спецификация | Назначение в RFC 9147 | Покрытие |
| --- | --- | --- |
| [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439) | Counter, nonce и block function для защиты номера последовательности ChaCha20 | Реализовано |
| [RFC 768](https://www.rfc-editor.org/rfc/rfc768) | Транспорт UDP | Нижний уровень; библиотека ограничена UDP и сохраняет границы дейтаграмм |
| [RFC 793](https://www.rfc-editor.org/rfc/rfc793) | Контекст надежного транспорта и MSL для старых эпох | Применимая семантика реализована; библиотека не принимает транспорт TCP |
| [RFC 1191](https://www.rfc-editor.org/rfc/rfc1191) | Определение IPv4 PMTU | Интеграция нижнего уровня реализована; ошибки MTU платформы приводятся к `ErrDatagramTooLarge`; приложение может активно проверять путь через `IgnorePathMTU` |
| [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) | Нормативные ключевые слова BCP 14 | Применимо |
| [RFC 4443](https://www.rfc-editor.org/rfc/rfc4443) | IPv6 ICMP Packet Too Big | Интеграция нижнего уровня реализована; ОС обрабатывает ICMPv6, библиотека использует обратную связь ошибок записи |
| [RFC 4821](https://www.rfc-editor.org/rfc/rfc4821) | Packetization Layer PMTU Discovery | Реализовано; ошибки записи и последовательные тайм-ауты black hole запускают откат и повторную фрагментацию |
| [RFC 6298](https://www.rfc-editor.org/rfc/rfc6298) | Начальный RTO, экспоненциальная задержка и максимум | Реализовано; одна секунда по умолчанию, максимум 60 секунд и измерения RTT по алгоритму Карна |
| [RFC 8174](https://www.rfc-editor.org/rfc/rfc8174) | Правила регистра BCP 14 | Применимо |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | DTLS Connection ID | Реализовано; согласование, обновления, маршрутизация, изоляция и ограничения ресурсов |
| [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446) | Базовый протокол TLS 1.3 | Поддерживаются рукопожатие, аутентификация, key schedule, PSK/0-RTT, KeyUpdate, PHA и exporter; семантика протокола соответствует RFC 9846 |

### Связанные спецификации и расширения

| Спецификация | Связь с реализацией | Статус |
| --- | --- | --- |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | KeyShare, PSK/HRR, NST, пределы AEAD, KeyUpdate, выбор сертификата, предупреждения и границы векторов TLS 1.3 | Реализовано для включенных возможностей; поддерживаются CH/CR `certificate_authorities`, CR `oid_filters` и выбор из нескольких сертификатов на обеих сторонах; возобновление mTLS сохраняет состояние аутентификации, политику/CA/срок действия и общий срок аутентификации; семантика `user_canceled` и `general_error` описана в общем статусе |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | TLS/DTLS `record_size_limit` | Клиент объявляет расширение по умолчанию; сервер отвечает только на предложение. Отправка следует ограничению узла, прием — локальному ограничению, отсутствие расширения восстанавливает максимум протокола, а PMTU остается независимой нижней границей |
| [RFC 8879](https://www.rfc-editor.org/rfc/rfc8879) | Сжатие сертификатов TLS/DTLS | Явно включаемый стандартный zlib; реализованы согласование CH/CR, сертификаты сервера и клиента, HRR, mTLS, PHA, фрагментация/повторная передача, transcript и ограниченная распаковка; если сжатие не дает выигрыша, используется обычный Certificate |
| [RFC 9149](https://www.rfc-editor.org/rfc/rfc9149) | Запросы tickets TLS/DTLS 1.3 | Клиент отдельно запрашивает число tickets для полного и возобновленного соединения; сервер возвращает ограниченный expected count, надежно отправляет несколько NST и сохраняет один ticket при отсутствии расширения |
| [RFC 9257](https://www.rfc-editor.org/rfc/rfc9257) | Рекомендации TLS 1.3 по внешним PSK | Реализованы DHE-only, несколько identity, откат при неизвестной identity, описание риска открытой identity, привязка происхождения ticket и политика 0-RTT для внешнего PSK |
| [RFC 9258](https://www.rfc-editor.org/rfc/rfc9258) | PSK Importer для TLS/DTLS 1.3 | Реализованы целевой вывод SHA-256/384, метка DTLS, wire-кодирование ImportedIdentity и отдельная метка binder |
| [RFC 9848](https://www.rfc-editor.org/rfc/rfc9848) | Получение конфигурации ECH через DNS | Библиотека разбирает полный ECHConfigList; получение DNS SVCB/HTTPS и декодирование presentation-format Base64 выполняет приложение |
| [RFC 9849](https://www.rfc-editor.org/rfc/rfc9849) | Encrypted ClientHello для TLS/DTLS | Реализованы клиент и сервер, HPKE, реконструкция Inner/Outer, padding, повторное использование контекста HRR, подтверждение принятия, retry-конфигурации, аутентификация отклонения, GREASE, возобновление PSK и 0-RTT |
| [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954) | Гибридный обмен ключами TLS/DTLS | Реализованы конструкция RFC и группы X25519MLKEM768, SecP256r1MLKEM768 и SecP384r1MLKEM1024 из текущего профиля ECDHE-MLKEM; конкретный профиль еще находится в очереди RFC Editor |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | BCP безопасности развертывания TLS/DTLS | Tickets используют AES-256-GCM и имеют срок от одной секунды до семи дней; ограничения RSA 2048 бит и сертификатов SHA-1/MD5 применяются к полному рукопожатию, корням доверия и возобновлению; исключения OCSP и DTLS 1.2 описаны в общем статусе |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | Проверка идентичности службы | DNS-ID/IP-ID строго проверяются по умолчанию; семантику приложения для других reference identifier реализует вызывающий код |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | Return Routability Check при смене адреса CID | Реализовано; enhanced check по умолчанию, basic check после отказа старого пути, смена привязки только после проверки, отдельное ограничение усиления пути-кандидата и проверка с резервным CID при его наличии |
| [RFC 8701](https://www.rfc-editor.org/rfc/rfc8701) | GREASE против окостенения протокола | Стратегия значений расширений реализована: `EnableGREASE` отправляет случайное пустое расширение в CH/CR/NST, получатель игнорирует его без записи в состояние согласования, HRR сохраняет значение; остальные точки MAY активно не используются |

Из необязательных расширений не реализованы RFC 9261 Exported Authenticators и RFC 9345 Delegated Credentials.

### Границы реализации

Следующие ограничения не уменьшают полноту обязательной семантики RFC 9147, но должны учитываться пользователями:

- Модуль реализует только DTLS 1.3 и не откатывается к DTLS 1.2, поэтому не заявляет полного соответствия требованию RFC 9325 о поддержке DTLS 1.2 библиотеками общего назначения.
- Демультиплексирование записей Heartbeat реализовано; полный протокол Heartbeat определен RFC 6520 и находится за пределами RFC 9147.
- Отправитель использует допустимый режим «одна запись на UDP-дейтаграмму» и не предоставляет необязательный API агрегации нескольких записей.
- Параллельные множественные запросы PHA не предоставляются; RFC разрешает, но не требует эту возможность.
- Автоматическая смена привязки RRC требует транспорт, принимающий данные от разных источников и способный отправлять выбранному адресу. Стандартный Listener это поддерживает; подключенные UDP-клиенты ограничены фильтрацией узла операционной системой. Пустой CID не позволяет однозначно маршрутизировать разные пятиэлементные кортежи.
- Сборка wolfSSL master `6502cdd` (строка версии 5.9.2) поддерживает CID, KeyUpdate, PHA, session tickets, 0-RTT, `SESSION_CERTS`, прямые внешние PSK и все три гибридные группы, но не RFC 8449, RFC 8879, RFC 9149, importer RFC 9258 и RRC RFC 9853. Сервер игнорирует `ticket_request` клиента go-dtls и сохраняет обычное возобновление; это доказывает только совместимый откат. Ее сборка ECH/HPKE не завершает DTLS-рукопожатие с принятым ECH, поэтому зафиксировано только успешное обычное рукопожатие с ECH GREASE; совместимость accepted-ECH не заявляется. Для гибридного обмена все три группы проходят с wolfSSL в роли клиента, а сервер wolfSSL проходит группы X25519 и P-256; этот сервер не завершает фрагментированный hybrid ClientHello, а группа P-384 завершается тайм-аутом даже без фрагментации. Другие ограничения узла: HRR сервера отклоняет 0-RTT клиента go-dtls; клиент не разбирает mTLS ticket go-dtls размером 1421 байт; клиент не повторяет Finished после потери финального ACK.

## Тесты производительности

После каждого push в `master` определяется и фиксируется последний commit ветки wolfSSL `master` для данного запуска, затем выполняются все Go benchmark и четырехсторонние benchmark через настоящий UDP. Точный SHA и медианы пяти запусков публикуются в отдельной ветке результатов: [посмотреть последний автоматический отчет](https://github.com/puernya/go-dtls/blob/benchmark-results/benchmark.ru.md). Плановая проверка в 08:00, 16:00 и 00:00 по времени Asia/Shanghai сравнивает SHA веток `master` go-dtls и wolfSSL и запускает тесты повторно, только если один из них изменился и для этого SHA go-dtls нет benchmark в очереди или в процессе выполнения. Исходники и результаты сборки wolfSSL кешируются по SHA. Pull request использует тот же процесс и показывает результаты непосредственно на странице PR, не обновляя отчет master.

Запуск всех тестов производительности:

```sh
go test -run '^$' -bench . -benchmem
```

Отдельный запуск полного соединения и record layer:

```sh
go test -run '^$' -bench '^BenchmarkConnectionHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkECHHandshakeLifecycle/(Direct|HRR)$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkExternalPSKHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkMutualTLSHandshakeLifecycle/(Full|Resumed)$' -benchmem -count=10
go test -run '^$' -bench '^BenchmarkCertificateCompression' -benchmem
go test -run '^$' -bench '^BenchmarkProtectedRecord(Seal|RoundTripInPlace)$' -benchmem -count=5
```

Репозиторий также содержит отдельные тесты производительности для наборов шифров, ACK, records/parsers, transcript, key schedule, KeyUpdate, сообщений рукопожатия, сборки и создания flight. Фиксируйте версию Go, CPU, `-cpu` и `-benchtime` и совместно анализируйте `ns/op`, `B/op`, `allocs/op` и профили полного соединения.

## Покрытие тестами

- Тесты политики сертификатов RFC 9325 охватывают конфигурацию сервера, прием клиентом, самоподписанные сертификаты, неотправленные корни доверия, `InsecureSkipVerify`, обычное и mTLS-возобновление, а также сравнение с поведением `crypto/x509` для корней доверия RSA 1024 бит/SHA-1.
- Тесты RFC 8449 охватывают CH/EE, минимум 64, направленные ограничения, недопустимые значения и сочетания расширений, аутентифицированное превышение, HRR, возобновление, 0-RTT, KeyUpdate, ACK, независимость от PMTU и совместимость со сторонними узлами без согласования расширения.
- Тесты RFC 8879 охватывают согласование ClientHello/CertificateRequest, zlib, CompressedCertificate, transcript, недопустимые алгоритмы/потоки/длины, ограничения распаковки, откат к обычному Certificate, HRR, возобновление, mTLS/PHA, пределы записей, фрагментацию/повторную передачу, слабые сети, жизненный цикл ресурсов и безопасный откат со сторонними узлами без поддержки расширения.
- Тесты RFC 9149 охватывают wire-кодирование CH/EE, строгие позиции расширения и alerts, инварианты HRR, счетчики полного/возобновленного/нулевого запроса, предел сервера, надежные множественные NST, параллельное расходование cache, аннулирование связанных tickets, слабые сети, реальный UDP и откат возобновления при игнорировании расширения сторонним узлом.
- Тесты RFC 8701 охватывают отображение фиксированной случайности, wire CH/CR/NST, требование последнего PSK, HRR, mTLS/PHA, session tickets, отсутствие записи в состояние согласования и допуск всех трех расширений wolfSSL через реальный UDP.
- Тесты RFC 9257/9258 охватывают независимый вывод importer, разделение KDF SHA-256/384 и binder `imp`/`ext`, прямые и импортированные PSK, несколько identity, фильтрацию HRR, ошибки identity/key/context, откат к сертификату, состояние соединения, возобновление и отзыв ticket, политику 0-RTT и слабые сети.
- Тесты RFC 9848/9849 охватывают публичные векторы конфигурации, ECHConfig/ECHConfigList, HPKE, Inner/Outer и padding, реконструкцию внешних расширений, подтверждение принятия HRR и отказ от downgrade, аутентифицированные retry-конфигурации, подавление клиентского сертификата, GREASE, возобновление, 0-RTT, фрагментацию, слабые сети и реальный UDP.
- Тесты слабой сети охватывают двунаправленные потери, задержку, изменение порядка и дублирование, включая комбинации CH/SH/Finished/ACK/HRR/возобновления mTLS.
- Тесты mTLS охватывают полное рукопожатие, возобновление PSK, 0-RTT, откат политики/CA и срок аутентификации обновленного ticket.
- Тесты предупреждений RFC 9846 охватывают рукопожатие, ожидание финального ACK, изменение порядка после рукопожатия, `close_notify` и локальные криптографические ошибки.
- Тесты RFC 9853 охватывают сообщения/state machine RRC, реальный UDP NAT rebinding, обновления CID, комбинации слабой сети и жизненный цикл ресурсов соединения.
- Fuzzing parser/record охватывает дифференциальную проверку копирующего и in-place расшифрования для всех четырех AEAD.
- Двунаправленные тесты через реальный UDP с wolfSSL master `6502cdd` охватывают HRR, рукопожатия с сертификатами RSA-PSS, Finished ACK, прикладные данные, AES-GCM, AES-128-CCM, прямые внешние PSK, CID, KeyUpdate, PHA, обычное возобновление сессии и все три гибридные группы в поддерживаемых узлом направлениях. Дополнительно проверяются отправленные go-dtls GREASE RFC 8701 в CH/CR/NST, откат без согласования RFC 9149, немедленная смена CID, возобновление mTLS и повтор Finished после потери финального ACK, а также 0-RTT клиента wolfSSL. Для ECH проверяется только откат GREASE, поскольку узел пока не завершает DTLS accepted-ECH.

Требования к среде разработки, обязательным проверкам, производительности и commit описаны в [CONTRIBUTING.ru.md](CONTRIBUTING.ru.md).

## Лицензия

Проект распространяется по [GNU General Public License v3.0](../LICENSE) (`GPL-3.0-only`).
