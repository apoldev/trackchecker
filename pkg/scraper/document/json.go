package document

import (
	"github.com/apoldev/trackchecker/pkg/scraper"
	"github.com/tidwall/gjson"
)

type JsonDoc struct {
	data *gjson.Result
}

func NewJson(data []byte) scraper.Document {

	result := gjson.ParseBytes(data)

	return &JsonDoc{
		data: &result,
	}

}

func (d *JsonDoc) Value() interface{} {

	if d.data == nil {
		return nil

	}

	return d.data.Value()
}

func (d *JsonDoc) FindOne(path string) scraper.Document {

	result := d.data.Get(path)

	return &JsonDoc{
		data: &result,
	}

}

func (d *JsonDoc) FindAll(path string) []scraper.Document {

	array := d.data.Get(path).Array()
	result := make([]scraper.Document, 0, len(array))

	for i := range array {
		result = append(result, &JsonDoc{
			data: &array[i],
		})
	}

	return result
}
