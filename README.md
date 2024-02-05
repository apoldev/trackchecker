# TrackChecker 

![deploy action](https://github.com/apoldev/trackchecker/actions/workflows/deploy.yml/badge.svg?branch=develop_apoldev)

TrackChecker - приложение для отслеживания посылок из различных почтовых и курьерских служб.

Это приложение вляется примером работы с GRPC, NATS, Swagger, Docker, Docker Swarm, CI/CD, и другими технологиями в golang.

### Содержание
___
#### [Использовано в разработке](#h1)
#### [Как посмотреть в бою](#h3)
#### [Как запустить локально](#h4)
#### [Как это работает](#h5)
#### [Добавленные парсеры](#h6)

<h3 id="h1">
Использовал в разработке
</h3>

___

* [Swagger, go-swagger](https://github.com/go-swagger/go-swagger) - для генерации http сервера из swagger.yml
* [GRPC](https://github.com/grpc/grpc-go)
* [NATS (JetStream)](https://github.com/nats-io/nats.go) - библиотека для работы с брокером сообщений NATS
* [Watermill](https://watermill.io/) - библиотека для работы с Kafka, RabbitMQ, etc...
* [Redis](https://github.com/redis/go-redis) - redis клиент
* [Testify](https://github.com/stretchr/testify) - тестирование
* [Mockery](https://github.com/vektra/mockery) - для генерации моков
* [Xpath](https://github.com/antchfx/xpath) - для парсера html, json, xml документов
* GitHub Action
* Docker, Docker Swarm кластер

<h3 id="h3">
Как посмотреть
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

<h3>Демонстрация</h3>

_Требуется предварительно запустить основное приложение командой выше_

#### Для запуска демонстрации запустите демо-клиент с помощью команды:

```bash
go run ./cmd/client/main.go
```


> Важно! Многие иностранные сайты не позволяют делать запросы с российских ip адресов. Поэтому, если вы запустите демо-клиент с российского ip, то вам частично могут вернуться пустые результаты или отмененные по таймауту.

<h3 id="h5">
Как это работает
</h3>

___

1. TrackChecker получает запрос со списком номеров отслеживания в виде массива строк.
2. Каждый трек-код отправляется в очередь отдельным сообщением. В текущей версии в качестве брокера сообщений используется NATS (JetStream).
3. Другая часть приложения забирает из очереди по одному трек-коду.
4. Трек-код проверяется в каждом парсере, у которого совпал по регулярному выражению.
5. Результаты складываются в HSET Redis.

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

<h3 id="h6">Добавленные парсеры</h3>

___

- [x] Почта России
- [x] Почта США
- [x] Почта Новой Зеландии
- [x] Почта Южной Кореи
- [x] Почта Малайзии
- [x] DPD Польши
- [x] Global Track&Trace
- [x] Почта Швеции

