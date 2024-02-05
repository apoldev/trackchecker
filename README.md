# TrackChecker 

![deploy action](https://github.com/apoldev/trackchecker/actions/workflows/deploy.yml/badge.svg?branch=develop_apoldev)

TrackChecker - приложение для отслеживания посылок из различных почтовых и курьерских служб.

Это приложение является примером работы с GRPC, NATS, Kafka, RabbitMQ, Swagger, Docker, Docker Swarm, CI/CD, и другими технологиями в golang.

Для примера здесь также реализовано мини-ядро парсера, которое принимает конфигурацию для парсера простых сайтов/api, отдающих xml, json, html.

#### [Подробнее в разделе Как это работает](#h6)

### Содержание
___
- #### [Использовано в разработке](#h1)
- #### [Как посмотреть в бою](#h3)
- #### [Как запустить локально](#h4)
  - #### [Демо клиент](#h5)
  - #### [Результат работы демо клиента](#h51)
  - #### [Пример запроса](#h52)
  - #### [Пример ответа](#h53)
- #### [Как это работает](#h6)
- #### [Добавленные парсеры](#h7)

<h3 id="h1">
Использовал в разработке
</h3>

___

* [Swagger, go-swagger](https://github.com/go-swagger/go-swagger) - для генерации http сервера из swagger.yml
* [GRPC](https://github.com/grpc/grpc-go)
* [NATS (JetStream)](https://github.com/nats-io/nats.go) - библиотека для работы с брокером сообщений NATS
* [Kafka, Watermill](https://watermill.io/) - библиотека для работы с Kafka, RabbitMQ, etc...
* [Redis](https://github.com/redis/go-redis) - redis клиент для golang
* [Testify](https://github.com/stretchr/testify) - тестирование
* [Mockery](https://github.com/vektra/mockery) - для генерации моков
* [Xpath](https://github.com/antchfx/xpath) - для парсера html, json, xml документов
* [Goquery](https://github.com/PuerkitoBio/goquery) - для парсера html документов
* [Gjson](https://github.com/tidwall/gjson) - для парсера json документов
* [Sjson](https://github.com/tidwall/sjson) - для записи в json документ
* GitHub Action
* Docker
* Docker Swarm кластер

<h3 id="h3">
Как посмотреть в бою
</h3>

___

Приложение деплоится на сервер с помощью github actions. 
Документация swagger доступна по адресу [trackchecker.1trackapp.com/docs](https://trackchecker.1trackapp.com/docs)

<h3 id="h4">
Как запустить локально
</h3>

```bash
docker-compose up
```

Если требуется изменить [/api/swagger.yaml](./api/swagger.yaml) и перестроить http сервер, то выполнить команду:

```bash
swagger generate server --exclude-main -f ./api/swagger.yaml -t ./internal/app/restapi --exclude-main
```

<h3 id="h5">Демо клиент</h3>

_Требуется предварительно запустить основное приложение командой выше_

#### Для запуска демонстрации запустите демо-клиент с помощью команды:

```bash
go run ./cmd/client/main.go
```

> Важно! Многие иностранные сайты не позволяют делать запросы с российских ip адресов. Поэтому, если вы запустите демо-клиент с российского ip, то вам частично могут вернуться пустые результаты или отмененные по таймауту.

<h4 id="h51">Результат работы демо клиента</h4>

___

```bash
LE704280574SE found at sweden-post: 2023-12-06T17:06:00Z, The shipment item has been dropped off by sender
LE704280574SE found at global-track-trace: 2023-12-07T16:22:00Z, Departure from outward office of exchange, SGSINN
RK166520145LV found at global-track-trace: 2024-01-19T11:28:00Z, Departure from outward office of exchange, LVRIXF
EH036261918US found at global-track-trace: 2024-01-21T00:43:00Z, Departure from outward office of exchange, USLAXA
EH036261918US found at usps: 2024-01-31T16:22:00Z, Delivered
LE704180823SE found at global-track-trace: 2023-12-05T17:19:00Z, Departure from outward office of exchange, SGSINN
LE704180823SE found at sweden-post: 2023-12-04T17:00:00Z, The shipment item has been dropped off by sender
LE703771444SE found at global-track-trace: 2023-11-16T14:02:00Z, Departure from outward office of exchange, SGSINN
LE703771444SE found at sweden-post: 2023-11-16T14:52:00Z, The shipment item has been dropped off by sender
LE704046934SE found at global-track-trace: 2023-11-29T16:36:00Z, Departure from outward office of exchange, SGSINN
LE704046934SE found at sweden-post: 2023-12-08T13:16:00Z, The shipment item has arrived at the country of destination
92748999997762543411393406 found at usps: 2024-01-12T12:38:00Z, Departed Shipping Partner Facility, USPS Awaiting Item, SAN LEANDRO, CA 94578
LM219469900SE not found at global-track-trace: unexpected end of JSON input
LM219469900SE found at sweden-post: 2024-02-05T10:17:00Z, The shipment item has been dropped off by sender
420751699214490233605881208508 found at usps: 2024-02-03T16:38:00Z, Arrived at USPS Regional Origin Facility, LOS ANGELES CA DISTRIBUTION CENTER
UA132229530HU found at global-track-trace: 2024-01-05T06:36:00Z, Departure from outward office of exchange, HUBUDB
9300110555800007963985 found at usps: 2024-02-04T18:32:00Z, Arrived at USPS Regional Origin Facility, CHICAGO IL NETWORK DISTRIBUTION CENTER
9400108205498532578275 found at usps: 2024-01-09T13:39:00Z, Shipping Label Created, USPS Awaiting Item, NIAGARA FALLS, NY 14305
LH256986182AU not found at global-track-trace: unexpected end of JSON input
LH256986182AU found at new-zealand-post: 2023-07-24T11:44:00Z, Picked up/Collected, , Your item has been collected by the overseas postal service and is en route to their depot
UM908307556US not found at global-track-trace: unexpected end of JSON input
UM908307556US found at usps: 2023-12-30T07:08:00Z, Departed, SAO PAULO
UE400083227US not found at global-track-trace: unexpected end of JSON input
UE400083227US found at usps: 2023-12-14T10:38:00Z, Departed, BRUSSELS
9200190348376028555454 found at usps: 2023-12-14T11:31:00Z, Arrived at Post Office, MIAMI, FL 33166
4203316692748927005455000598734364 found at usps: 2024-01-28T11:08:00Z, Delivered to Agent for Final Delivery, MIAMI, FL 33166
RC211515121MY found at global-track-trace: 2024-01-10T13:31:00Z, Posting/Collection, MYPENB
RC211515121MY found at malaysia-post: 2024-01-10T17:17:38Z, Departed from International Hub to Overseas Destination
EG014132620KR not found at global-track-trace: unexpected end of JSON input
EG014132620KR found at south-korea: 2024-02-01T16:33:00Z, 발송, 고양일산우체국
EX407171527KR not found at global-track-trace: unexpected end of JSON input
EX407171527KR found at south-korea: 2024-01-16T20:34:00Z, 발송, 인천해상교환우체국
EG012385263KR found at global-track-trace: 2023-12-23T14:04:00Z, Final delivery, MNUB19
EG012385263KR found at south-korea: 2023-12-13T19:01:00Z, 발송, 고양덕양우체국
EG013492697KR not found at global-track-trace: unexpected end of JSON input
EG013492697KR found at south-korea: 2024-01-17T17:48:00Z, 발송, 안산우체국
UD657184909MY found at global-track-trace: 2024-01-02T16:56:00Z, Departure from outward office of exchange, MYPENB
UD657184909MY found at malaysia-post: 2024-01-02T18:11:54Z, Item Sent to Burma, In Transit
10881378073059 not found at russian-post: Get "http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru&track-numbers=10881378073059": context deadline exceeded
80513392264272 not found at russian-post: Get "http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru&track-numbers=80513392264272": context deadline exceeded
12907591218904 not found at russian-post: Get "http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru&track-numbers=12907591218904": dial tcp 212.164.138.79:80: i/o timeout
66013290014431 not found at russian-post: Get "http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru&track-numbers=66013290014431": dial tcp 212.164.138.79:80: i/o timeout
CL123084655RU not found at global-track-trace: unexpected end of JSON input
CL123084655RU not found at russian-post: Get "http://www.pochta.ru/api/tracking/api/v1/trackings/by-barcodes?language=ru&track-numbers=CL123084655RU": context deadline exceeded
UD660079300MY found at global-track-trace: 2024-01-29T01:35:00Z, Posting/Collection, MYKULC
UD660079300MY found at malaysia-post: 2024-01-29T10:33:56Z, Item Sent to Uzbekistan, In Transit
UD656337373MY found at global-track-trace: 2024-01-03T14:18:00Z, Departure from outward office of exchange, MYKULC
UD656337373MY found at malaysia-post: 2023-12-29T10:21:31Z, Item Sent to Namibia, In Transit
LP610391713MY found at global-track-trace: 2024-01-15T15:03:00Z, Departure from outward office of exchange, MYJHBB
LP610391713MY found at malaysia-post: 2024-02-05T10:59:00Z, Departed from International Hub to domestic location
```

> Почта России отвалилась по таймауту, так как она заблокировала мой ip адрес за злоупотребление запросами :)


<h4 id="h52">Пример запроса</h4>
___

```json
{
    "tracking_numbers": [
        "LH256986182AU",
        "UD656337373MY"
    ]
}
```

<h4 id="h53">Пример ответа</h4>
___

```json
{
  "data": [
    {
      "code": "UD656337373MY",
      "id": "50047d80-1880-4071-953f-b3bec70c3a91",
      "results": [
        {
          "execute_time": 0.017348659,
          "result": {
            "CountryTo": "NA",
            "CountryFrom": "MY",
            "events": [
              {
                "status": "Departure from outward office of exchange",
                "date": "2024-01-03T14:18:00Z",
                "place": "MYKULC"
              }
            ]
          },
          "spider": "global-track-trace",
          "tracking_number": "UD656337373MY"
        },
        {
          "execute_time": 0.843198383,
          "result": {
            "events": [
              {
                "status": "Item Sent to Namibia",
                "date": "2023-12-29T10:21:31Z",
                "place": "In Transit"
              },
              {
                "status": "Item Posted Over The Counter to Namibia",
                "date": "2023-12-29T09:45:49Z",
                "place": "In Transit"
              },
              {
                "status": "Dispatch PreAlert to Namibia",
                "date": "2023-12-26T15:02:53Z",
                "place": "In Transit"
              }
            ]
          },
          "spider": "malaysia-post",
          "tracking_number": "UD656337373MY"
        }
      ],
      "status": "finish",
      "uuid": "ffd8c81b-9804-4622-9e21-48f15ae69e55"
    },
    {
      "code": "LH256986182AU",
      "id": "50047d80-1880-4071-953f-b3bec70c3a91",
      "results": [
        {
          "error": "unexpected end of JSON input",
          "execute_time": 0.246106831,
          "result": null,
          "spider": "global-track-trace",
          "tracking_number": "LH256986182AU"
        },
        {
          "execute_time": 1.242484113,
          "result": {
            "SignedBy": "Acp602 Mailroom",
            "events": [
              {
                "status": "Picked up/Collected",
                "date": "2023-07-24T11:44:00Z",
                "details": "Your item has been collected by the overseas postal service and is en route to their depot"
              },
              {
                "status": "International departure",
                "date": "2023-07-28T04:21:00Z",
                "place": "BRISBANE",
                "details": "Departure from country of origin. Your item is in transit to New Zealand"
              },
              {
                "status": "International arrival",
                "date": "2023-07-31T19:17:04Z",
                "place": "AUCKLAND",
                "details": "Your item has arrived in New Zealand"
              },
              {
                "status": "In transit to local depot",
                "date": "2023-07-31T23:21:16Z",
                "place": "AUCKLAND",
                "details": "Your item has left our International Mail Centre in Auckland and is on its way to a local/regional delivery depot"
              },
              {
                "status": "At local/regional depot",
                "date": "2023-08-01T11:45:22Z",
                "place": "Auckland (Ak Central/East Depot)",
                "details": "Your item has been sorted at a parcel depot"
              },
              {
                "status": "With courier for delivery",
                "date": "2023-08-01T17:53:44Z",
                "place": "Auckland (Ak City CP Depot)",
                "details": "Your item is with a courier for delivery"
              },
              {
                "status": "Delivery Complete",
                "date": "2023-08-01T18:48:26Z",
                "place": "Auckland (Ak City CP Depot)",
                "details": "Your item has been successfully delivered and was signed for by \"Acp602 Mailroom\""
              }
            ]
          },
          "spider": "new-zealand-post",
          "tracking_number": "LH256986182AU"
        }
      ],
      "status": "finish",
      "uuid": "a105efb8-e409-4ebb-a796-92bd5d073e42"
    }
  ],
  "status": true
}
```

<h3 id="h6">
Как это работает
</h3>

___

1. TrackChecker получает запрос со списком номеров отслеживания в виде массива строк.
2. Каждый трек-код отправляется в очередь отдельным сообщением. В текущей версии в качестве брокера сообщений используется NATS (JetStream).
3. Другая часть приложения забирает из очереди по одному трек-коду.
4. Трек-код проверяется в каждом парсере, у которого совпал по регулярному выражению.
5. Результаты складываются в HSET Redis.
6. Клиент через некоторое время запрашивает результаты запроса.

Сами парсеры представляют из себя структуру, деклараттивно описывающую "как парсить", в которых указана последовательность действий для выполнения http запросов и дальнейшего парсинга этого документа с помощью:
* [xpath](https://github.com/antchfx/xpath) для html, xml, json
* [goquery](https://github.com/PuerkitoBio/goquery) для html
* [gjson](https://github.com/tidwall/gjson) запросы для json

Пример декларативного описания парсера для Почты США (USPS):
```json
{
  "code":"usps",
  "masks": [
    "[A-Z]{2}[0-9]{9}US"
  ],
  "examples": [
    "EH036261918US"
  ],
  "tasks":[
    {
      "type":"request",
      "payload":"http://production.shippingapis.com/ShippingApi.dll?API=TrackV2&XML=%3CTrackFieldRequest%20USERID=%22707HGUPS0501%22%3E%3CTrackID%20ID=%22[track]%22/%3E%3C/TrackFieldRequest%3E",
      "params":{
        "method":"GET",
        "type":"xml"
      }
    },
    {
      "type":"query",
      "payload":"//TrackSummary|//TrackDetail",
      "field":{
        "path":"events",
        "type":"array",
        "element":{
          "type":"object",
          "object":[
            {
              "path":"status",
              "query":".//Event"
            },
            {
              "path":"date",
              "query":"concat(./EventDate,' ', ./EventTime)"
            },
            {
              "path":"place",
              "query":"concat(./EventCity,', ', ./EventState, ' ', ./EventZIPCode)"
            },
            {
              "path":"country",
              "query":"./EventCountry"
            }
          ]
        }
      }
    }
  ]
}
```

<h3 id="h7">Добавленные почтовые службы</h3>

___

- [x] Почта России
- [x] Почта США
- [x] Почта Новой Зеландии
- [x] Почта Южной Кореи
- [x] Почта Малайзии
- [x] DPD Польши
- [x] Global Track&Trace
- [x] Почта Швеции

