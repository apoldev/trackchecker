package document

import (
	"github.com/tidwall/gjson"
)

type JSONDoc struct {
	data *gjson.Result
}

func NewJSON(data []byte) (*JSONDoc, error) {
	result := gjson.ParseBytes(data)

	return &JSONDoc{
		data: &result,
	}, nil
}

func (d *JSONDoc) Value() interface{} {
	if d.data == nil {
		return nil
	}

	return d.data.Value()
}

func (d *JSONDoc) FindOne(path string) (Document, error) {
	result := d.data.Get(path)

	if !result.Exists() {
		return nil, ErrNotExists
	}

	return &JSONDoc{
		data: &result,
	}, nil
}

func (d *JSONDoc) FindAll(path string) []Document {
	array := d.data.Get(path).Array()
	result := make([]Document, 0, len(array))

	for i := range array {
		result = append(result, &JSONDoc{
			data: &array[i],
		})
	}

	return result
}
