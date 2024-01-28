package document

import (
	"github.com/tidwall/gjson"
)

type JsonDoc struct {
	data *gjson.Result
}

func NewJson(data []byte) Document {

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

func (d *JsonDoc) FindOne(path string) Document {

	result := d.data.Get(path)

	return &JsonDoc{
		data: &result,
	}

}

func (d *JsonDoc) FindAll(path string) []Document {

	array := d.data.Get(path).Array()
	result := make([]Document, 0, len(array))

	for i := range array {
		result = append(result, &JsonDoc{
			data: &array[i],
		})
	}

	return result
}
