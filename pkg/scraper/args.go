package scraper

import (
	"github.com/apoldev/trackchecker/pkg/scraper/document"
	"github.com/tidwall/sjson"
	"net/http"
)

const (
	HTML      = "html"
	JSON      = "json"
	JSONXpath = "json.xpath"
	XML       = "xml"
	XPATH     = "xpath"
)

// Args is a struct for passing arguments to scraper
type Args struct {
	ResultBuilder *ResultBuilder

	document   document.Document
	variables  Variables
	httpClient *http.Client
}

func NewArgs(variables Variables, httpClient *http.Client) *Args {
	cl := http.DefaultClient
	if httpClient != nil {
		cl = httpClient
	}

	return &Args{
		ResultBuilder: NewResultBuilder(),
		variables:     variables,
		httpClient:    cl,
	}
}

type ResultBuilder struct {
	data []byte
}

func NewResultBuilder() *ResultBuilder {
	return &ResultBuilder{
		data: []byte(`{}`),
	}
}

func (b *ResultBuilder) GetString() string {
	return string(b.data)
}

func (b *ResultBuilder) GetData() []byte {
	return b.data
}

func (b *ResultBuilder) Set(path string, value interface{}) {
	var err error
	var newdata []byte

	newdata, err = sjson.SetBytes(b.data, path, value)

	if err != nil {
		return
	}

	b.data = newdata
}
