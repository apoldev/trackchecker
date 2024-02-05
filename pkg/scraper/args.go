package scraper

import (
	"net/http"
	"time"

	"github.com/apoldev/trackchecker/pkg/scraper/transform"

	"github.com/apoldev/trackchecker/pkg/scraper/document"
	"github.com/tidwall/sjson"
)

const (
	HTML      = "html"
	JSON      = "json"
	JSONXpath = "json.xpath"
	XML       = "xml"
	XPATH     = "xpath"
)

// Args is a struct for passing arguments to scraper.
type Args struct {
	ResultBuilder *ResultBuilder
	ExecuteTime   time.Duration

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

func (b *ResultBuilder) Set(path string, value interface{}, transformers ...transform.Transformer) {
	var err error
	var newdata []byte

	if v, ok := value.(string); ok {
		value = b.applyTransformers(v, transformers)
	}
	newdata, err = sjson.SetBytes(b.data, path, value)

	if err != nil {
		return
	}

	b.data = newdata
}

func (b *ResultBuilder) applyTransformers(data string, transformers []transform.Transformer) string {
	str := data
	for _, transformer := range transformers {
		switch transformer.Type {
		case transform.TypeClean:
			str = transform.Clean(str)
		case transform.TypeDate:
			str = transform.Date(str)
		case transform.TypeReplaceString:
			str = transform.ReplaceStr(str, transformer.Params["old"], transformer.Params["new"])
		}
	}
	return str
}
