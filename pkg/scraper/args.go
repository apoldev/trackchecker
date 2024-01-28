package scraper

import (
	"github.com/tidwall/sjson"
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
	Document      Document
	Variables     Variables
	ResultBuilder *ResultBuilder
}

// Document can be HTML, JSON, XML, etc.
type Document interface {
	Value() interface{}
	FindOne(expr string) Document
	FindAll(expr string) []Document
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
