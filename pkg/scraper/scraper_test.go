package scraper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/tidwall/gjson"
)

var htmlData = `<div xmlns:xs="http://www.w3.org/2001/XMLSchema" class="single-package"><input type="hidden" value="13559300030524" class="js-waybill"><input type="hidden" value="13559300030524" class="js-waybill-paczki"><fieldset class="compact"> <div class="form-group"><span class="label">Przesyłka</span><span class="input"><span class="input-text">13559300030524</span></span></div> <div class="form-group"><label>Paczki w przesyłce</label><span class="input"><select name="parcel" class="custom-select"> <option value="13559300030524">13559300030524</option></select></span></div> <div class="form-group"><span class="label">Paczka</span><span class="input"><span class="input-text">13559300030524  </span></span></div> <div class="form-group"><span class="label">Jesteś odbiorcą? Możesz samodzielnie zarządzać przesyłką</span><span class="input"><a href="https://mojapaczka.dpd.com.pl/login?parcel=13559300030524"> <div class="column arrow"><span class="btn--arrow-right btn--no-text"></span></div></a></span><span class="label l-info"><span class="l-info">Portal pozwala na przekierowanie, zmianę daty doręczenia, rezygnację z paczki i przekierowanie do punktu Pickup</span></span></div> <div class="js-package-details subform" style="display: none;"></div> </fieldset> <fieldset class="compact"> <h3>Historia przesyłki</h3> <div class="table-wrapper-400"> <table class="table-track"> <thead> <th>Data</th> <th>Godzina</th> <th>Opis</th> <th>Oddział</th> </thead> <tbody> <tr> <td>2023-12-08</td> <td>11:45:04</td> <td>Przesyłka doręczona</td> <td></td> </tr> <tr> <td>2023-12-08</td> <td>09:22:53</td> <td>Przyjęcie przesyłki do punktu Pickup</td> <td></td> </tr> <tr> <td>2023-12-08</td> <td>09:21:32</td> <td>Przyjęcie przesyłki do punktu Pickup</td> <td></td> </tr> <tr> <td>2023-12-08</td> <td>08:16:57</td> <td>Wydanie do doręczenia za granicą</td> <td></td> </tr> <tr> <td>2023-12-08</td> <td>05:24:10</td> <td>Przyjęcie przesyłki w oddziale doręczenia za granicą</td> <td></td> </tr> <tr> <td>2023-12-07</td> <td>18:35:44</td> <td>Przeładunek w sortowni za granicą</td> <td></td> </tr> <tr> <td>2023-12-07</td> <td>18:35:40</td> <td>Przeładunek w sortowni za granicą</td> <td></td> </tr> <tr> <td>2023-12-07</td> <td>01:02:02</td> <td>Przekazano za granicę<br><a href="https://www.dpdgroup.com/pl/mydpd/my-parcels/search?lang=pl&parcelNumber=13559300030524" class="underline" target="_blank">13559300030524</a></td> <td>NCS</td> </tr> <tr> <td>2023-12-06</td> <td>18:10:18</td> <td>Przyjęcie przesyłki w oddziale DPD </td> <td>LEG</td> </tr> <tr> <td>2023-12-06</td> <td>15:54:08</td> <td>Powiadomienie mail</td> <td>WA1</td> </tr> <tr> <td>2023-12-06</td> <td>15:14:57</td> <td>Przesyłka odebrana przez Kuriera</td> <td>MDZ</td> </tr> <tr> <td>2023-12-06</td> <td>10:57:52</td> <td>Zarejestrowano dane przesyłki, przesyłka jeszcze nienadana</td> <td></td> </tr> </tbody> </table> </div> </fieldset> </div>`

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type BrokenReader struct{}

func (br *BrokenReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed reading")
}

func (br *BrokenReader) Close() error {
	return fmt.Errorf("failed closing")
}

func TestScraper_Request(t *testing.T) {
	httpClient := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			assert.Equal(t, req.URL.String(), "https://tracktrace.dpd.com.pl/EN/findPackage")
			assert.Equal(t, req.Method, "POST")

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(htmlData)),
			}
		},
		)}

	dpdPolandScraperXpath := Scraper{
		Code: "dpd-poland",
		Tasks: []Task{
			{
				Type:    "request",
				Payload: "https://tracktrace.dpd.com.pl/EN/findPackage",
				Params: map[string]interface{}{
					"type":   "xpath",
					"method": "POST",
					"body":   "q=[track]&typ=1",
				},
			},
			{
				Type:    "query",
				Payload: `//table[@class="table-track"]/tbody/tr`,
				Field: Field{
					Path: "events",
					Type: FieldTypeArray,
					Element: &Field{
						Type: FieldTypeObject,
						Object: []*Field{
							{
								Path:  "status",
								Query: "./td[3]/text()[1]",
							},
							{
								Path:  "fulldate",
								Query: "concat(//td[1],' ', //td[2])",
							},
							{
								Path:  "date",
								Query: "td[1]",
							},
							{
								Path:  "time",
								Query: "td[2]",
							},
							{
								Path:  "place",
								Query: "./td[4]",
							},
							{
								Path:  "details",
								Query: "./td[3]/a",
							},
						},
					},
				},
			},
		},
	}

	args := NewArgs(Variables{
		"[track]": "13559300030524",
	}, httpClient)

	dpdPolandScraperXpath.Scrape(args)

	data := args.ResultBuilder.GetData()
	eventCount := gjson.GetBytes(data, "events.#").Int()

	assert.Equal(t, eventCount, int64(12))
	assert.Equal(t, gjson.GetBytes(data, "events.0.date").String(), "2023-12-08")

	httpClient2 := &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       &BrokenReader{},
			}
		},
		)}

	args2 := NewArgs(nil, httpClient2)
	err := dpdPolandScraperXpath.Scrape(args2)

	assert.Equal(t, err.Error(), "failed reading")
	assert.Equal(t, args2.ResultBuilder.GetString(), "{}")
}

func TestScraper_Parsing(t *testing.T) {
	dpdPolandScraperXpath := Scraper{
		Code: "dpd-poland",
		Tasks: []Task{
			{
				Type:    "source",
				Payload: htmlData,
				Params: map[string]interface{}{
					"type": "xpath",
				},
			},
			{
				Type:    "query",
				Payload: `//table[@class="table-track"]/tbody/tr`,
				Field: Field{
					Path: "events",
					Type: FieldTypeArray,
					Element: &Field{
						Type: FieldTypeObject,
						Object: []*Field{
							{
								Path:  "status",
								Query: "./td[3]/text()[1]",
							},
							{
								Path:  "fulldate",
								Query: "concat(//td[1],' ', //td[2])",
							},
							{
								Path:  "date",
								Query: "td[1]",
							},
							{
								Path:  "time",
								Query: "td[2]",
							},
							{
								Path:  "place",
								Query: "./td[4]",
							},
							{
								Path:  "details",
								Query: "./td[3]/a",
							},
						},
					},
				},
			},
		},
	}

	args := NewArgs(nil, http.DefaultClient)

	dpdPolandScraperXpath.Scrape(args)

	data := args.ResultBuilder.GetData()
	eventCount := gjson.GetBytes(data, "events.#").Int()

	assert.Equal(t, eventCount, int64(12))
	assert.Equal(t, gjson.GetBytes(data, "events.0.date").String(), "2023-12-08")

	args2 := NewArgs(nil, http.DefaultClient)
	dpdPolandScraperXpath.Tasks[0].Params["type"] = "unknown"
	dpdPolandScraperXpath.Scrape(args2)

	assert.Equal(t, args2.ResultBuilder.GetString(), "{}")
}
