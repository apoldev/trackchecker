# TrackChecker 

![deploy action](https://github.com/apoldev/trackchecker/actions/workflows/deploy.yml/badge.svg?branch=develop_apoldev)

TrackChecker is an application for tracking parcels from various postal and courier services.

This application demonstrates working with GRPC, NATS, Kafka, RabbitMQ, Swagger, Docker, Docker Swarm, GitHub Actions, and other technologies in Go.

As an example, it also implements a mini-parser engine that accepts configurations for parsing simple websites/APIs returning XML, JSON, or HTML.

#### [Learn more in the How It Works section](#h6)

### Table of Contents
___
- #### [Technologies Used](#h1)
- #### [How to View in Action](#h3)
- #### [How to Run Locally](#h4)
  - #### [Demo Client](#h5)
  - #### [Demo Client Results](#h51)
  - #### [Sample Request](#h52)
  - #### [Sample Response](#h53)

- #### [How It Works](#h6)
  - #### [Parsers](#h7)
    - #### [Sample Declarative Description of USPS Parser in JSON](#h71)
    - #### [Sample CDEK Parser in Go](#h72)
  - #### [Added Parsers](#h99)

<h3 id="h1">
Technologies Used
</h3>

___
* Technologies
  * [Golang](https://golang.org/)
  * [Docker](https://www.docker.com/)
  * [NATS (JetStream)](https://docs.nats.io/nats-concepts/jetstream)
  * [Kafka](https://kafka.apache.org/)
  * [Swagger](https://swagger.io/)
  * [GRPC](https://grpc.io/)
  * [Redis](https://redis.io/)

* Development
  * [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - for automatic application rebuild in a container
  * [go-swagger](https://github.com/go-swagger/go-swagger) - for generating an HTTP server from swagger.yml
  * [xpath](https://github.com/antchfx/xpath) - for parsing HTML, JSON, XML documents
  * [goquery](https://github.com/PuerkitoBio/goquery) - for parsing HTML documents
  * [gjson](https://github.com/tidwall/gjson) - for parsing JSON documents
  * [sjson](https://github.com/tidwall/sjson) - for writing to JSON documents
  * [grpc-go](https://github.com/grpc/grpc-go)
  * [nats.go](https://github.com/nats-io/nats.go) - library for working with the NATS message broker
  * [watermill](https://watermill.io/) - library for working with Kafka, RabbitMQ, etc.
  * [go-redis](https://github.com/redis/go-redis) - Redis client for Golang

  * Testing
    * [Testify](https://github.com/stretchr/testify) - testing
    * [Mockery](https://github.com/vektra/mockery) - for generating mocks

* Deployment
  * [Traefik](https://doc.traefik.io/traefik/)
  * GitHub Actions
  * Docker
  * Docker Swarm


<h3 id="h3">
How to View in Action
</h3>

___

The application is deployed to a server using GitHub Actions. 
Swagger documentation is available at [trackchecker.1trackapp.com/docs](https://trackchecker.1trackapp.com/docs)

<h3 id="h4">
How to Run Locally
</h3>

```bash
docker-compose up
```

If you need to modify [/api/swagger.yaml](./api/swagger.yaml) and rebuild the HTTP server, install [go-swagger](https://goswagger.io/) and run the command to generate code:

```bash
swagger generate server --exclude-main -f ./api/swagger.yaml -t ./internal/app/restapi --exclude-main
```

<h3 id="h5">Demo Client</h3>

_The main application must be running before proceeding._

#### To run the demo, start the demo client with the following command:

```bash
go run ./cmd/client/main.go
```

> Important! Many foreign websites restrict requests from Russian IP addresses. Therefore, if you run the demo client from a Russian IP, some results may be partially empty or canceled due to timeouts.

<h4 id="h51">Demo Client Results</h4>

___

```bash
LE704280574SE found at sweden-post: 2023-12-06T17:06:00Z, The shipment item has been dropped off by sender
LE704280574SE found at global-track-trace: 2023-12-07T16:22:00Z, Departure from outward office of exchange, SGSINN
RK166520145LV found at global-track-trace: 2024-01-19T11:28:00Z, Departure from outward office of exchange, LVRIXF
EH036261918US found at global-track-trace: 2024-01-21T00:43:00Z, Departure from outward office of exchange, USLAXA
EH036261918US found at usps: 2024-01-31T16:22:00Z, Delivered
LE704180823SE found at global-track-trace: 2023-12-05T17:19:00Z, Departure from outward office of exchange, SGSINN
LE704180823SE found at sweden-post: 2023-12-04T17:00:00Z, The shipment item has been dropped off by sender
...
```

> Russian Post timed out because it blocked my IP address for excessive requests. :)


<h4 id="h52">Sample POST Request</h4>
___

```bash
curl --location 'http://localhost:7777/track' --header 'Content-Type: application/json' --data '{
    "tracking_numbers": [
        "LH256986182AU",
        "UD656337373MY"
    ]
}'
```


<h4>Sample Response</h4>
___

```json
{
  "tracking_numbers": [
    {
      "code": "LH256986182AU",
      "uuid": "50047d80-1880-4071-953f-b3bec70c3a91"
    },
    {
      "code": "UD656337373MY",
      "uuid": "50047d80-1880-4071-953f-b3bec70c3a91"
    }
  ],
  "tracking_id": "4a89c077-20b7-447c-86b7-dc93582af6b2"
}
```


<h4>Sample GET Request for Results</h4>

___

```bash
curl --location 'http://localhost:7777/track?id=4a89c077-20b7-447c-86b7-dc93582af6b2'
```

<h4 id="h53">Sample Response for Results</h4>
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
    }
  ],
  "status": true
}
```

<h3 id="h6">
How It Works
</h3>

___

1. TrackChecker receives a request with a list of tracking numbers as an array of strings.
2. Each tracking code is sent to the queue as a separate message. The current version uses NATS (JetStream) as the message broker.
3. Another part of the application retrieves one tracking code at a time from the queue.
4. The tracking code is checked with each parser matching a regular expression.
5. Results are stored in HSET Redis.
6. The client queries the results after some time.


<h3 id="h7">Parsers</h3>

___

The parsers are structures that declaratively describe "how to parse," specifying a sequence of actions to perform HTTP requests and subsequently parse the document using:
* [xpath](https://github.com/antchfx/xpath) for HTML, XML, JSON
* [goquery](https://github.com/PuerkitoBio/goquery) for HTML
* [gjson](https://github.com/tidwall/gjson) queries for JSON

<h3 id="h71">Sample Declarative Description of USPS Parser</h3>

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

<h3 id="h72">Sample Declarative Description of CDEK Parser (Go)</h3>

```go
cdek := models.Spider{
		Scraper: scraper.Scraper{
			Code: "cdek",
			Tasks: []scraper.Task{
				{
					Type:    scraper.TaskTypeRequest,
					Payload: `https://mobile-apps.cdek.ru/api/v2/order/[track]`,
					Params: map[string]interface{}{
						"type":   scraper.JSONXpath,
						"method": "GET",
						"headers": map[string]string{
							"User-Agent":      "CDEK/2.5 (com.cdek.cdekapp; build:1; iOS 13.3.1) Alamofire/4.9.1",
							"Accept-Language": "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7",
							"X-User-Lang":     "ru",
						},
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `concat(//office/city/name,', ', //office/address)`,
					Field: scraper.Field{
						Path: "AddressTo",
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `//office/latitude`,
					Field: scraper.Field{
						Path: "delivery_office.latitude",
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `//office/longitude`,
					Field: scraper.Field{
						Path: "delivery_office.longitude",
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `//realWeight`,
					Field: scraper.Field{
						Path: "Weight",
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `//additionalInfo//goods/*/name`,
					Field: scraper.Field{
						Path: "Goods",
						Type: scraper.FieldTypeArray,
					},
				},
				{
					Type:    scraper.TaskTypeQuery,
					Payload: `(//orderStatusGroups/* | //orderStatusGroups//statuses/*)`,
					Field: scraper.Field{
						Path: "events",
						Type: scraper.FieldTypeArray,
						Element: &scraper.Field{
							Type: scraper.FieldTypeObject,
							Object: []*scraper.Field{
								{
									Path:  "status",
									Query: "./title",
								},
								{
									Path:  "date",
									Query: "./date",
									Transformers: []transform.Transformer{
										{
											Type: transform.TypeReplaceRegexp,
											Params: map[string]string{
												"regexp": "(\\d+)\\.(\\d+)\\.(\\d+)",
												"new":    "$3-$2-$1",
											},
										},
										{
											Type: transform.TypeDate,
										},
									},
								},
								{
									Path:  "place",
									Query: "../../city",
								},
							},
						},
					},
				},
			},
		},
		RegexpMasks: []*regexp.Regexp{
			regexp.MustCompile(`^(\d{10})$`),
			regexp.MustCompile(`^CN[A-Z0-9]+RU[0-9]{1}$`),
		},
	}
```

<h3 id="h99">Added parsers (postal services)</h3>

___

- [x] Russian post
- [x] USPS (USA post)
- [x] New Zealand post
- [x] South Korea
- [x] Malaysia post
- [x] DPD Poland
- [x] Global Track&Trace
- [x] Sweden Post
- [x] CDEK (Express courier)

