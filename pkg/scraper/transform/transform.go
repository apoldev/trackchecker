package transform

import "encoding/json"

const (
	TypeClean         = "clean"
	TypeDate          = "date"
	TypeReplaceString = "replace.string"
)

type Transformer struct {
	Type   string
	Params map[string]string
}

type Transformers []Transformer

func (t *Transformers) UnmarshalJSON(data []byte) error {
	var el []interface{}
	if err := json.Unmarshal(data, &el); err != nil {
		return err
	}
	transformers := make([]Transformer, 0, len(el))

Loop:
	for i := range el {
		var transformer Transformer
		switch v := el[i].(type) {
		case map[string]interface{}:
			transformer = Transformer{
				Type: v["type"].(string),
			}
			if params, ok := v["params"].(map[string]interface{}); ok {
				transformer.Params = make(map[string]string)
				for k, v := range params {
					if str, ok2 := v.(string); ok2 {
						transformer.Params[k] = str
					}
				}
			}
		case string:
			transformer.Type = v
		default:
			continue Loop
		}
		transformers = append(transformers, transformer)
	}

	*t = transformers
	return nil
}
