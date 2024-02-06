package main

import (
	"fmt"
	"github.com/apoldev/trackchecker/internal/app/models"
	"github.com/apoldev/trackchecker/pkg/scraper"
	"github.com/apoldev/trackchecker/pkg/scraper/transform"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

func main() {

	trackingNumbers := []string{
		"1515005006",
		"1340837686",
		"1517180126",
		"CN0008610350RU9",
		"1490020945",
		"1516134320",
		"1517556629",
		"1516621181",
		"1515939433",
		"1516783111",
		"1516835272",
		"1517550130",
		"CNV0000100342RU2",
	}

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

	results := make(map[string]string, len(trackingNumbers))
	mu := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(trackingNumbers))
	for _, track := range trackingNumbers {
		go func(track string) {
			defer wg.Done()

			args := scraper.NewArgs(scraper.Variables{
				"[track]": track,
			}, http.DefaultClient)
			err := cdek.Scrape(args)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				results[track] = err.Error()
				return
			}
			results[track] = args.ResultBuilder.GetString()
		}(track)
	}
	wg.Wait()

	// Show table
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Tracking Number", "Goods", "Weight, kg", "To", "Office coords", "Event"})
	for i := range trackingNumbers {
		data := results[trackingNumbers[i]]
		valid := gjson.Valid(data)
		if !valid {
			t.AppendRow(table.Row{
				i + 1,
				trackingNumbers[i],
				data,
			})
			continue
		}

		addressTo := gjson.Get(data, "AddressTo")

		coords := gjson.Get(data, "delivery_office.latitude").String() + ", " +
			gjson.Get(data, "delivery_office.longitude").String()

		var event gjson.Result
		if events := gjson.Get(data, "events").Array(); len(events) > 0 {
			event = events[len(events)-1]
		}

		status := event.Get("date").String() + ", " + event.Get("status").String()
		weight := gjson.Get(data, "Weight")
		goods := gjson.Get(data, "Goods").Array()
		goodsStr := make([]string, 0, len(goods))
		for j := range goods {
			runes := []rune(goods[j].String())
			if len(runes) > 21 {
				runes = runes[:20]
				runes = append(runes, '.', '.', '.')
			}
			goodsStr = append(goodsStr, string(runes))
		}

		t.AppendRow(table.Row{
			i + 1,
			trackingNumbers[i],
			strings.Join(goodsStr, "\n"),
			weight,
			addressTo,
			coords,
			status,
		})

		t.AppendRow(table.Row{""})
		t.AppendRow(table.Row{""})
	}
	fmt.Println(t.Render())

}
