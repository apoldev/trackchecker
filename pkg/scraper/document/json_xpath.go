package document

import (
	"bytes"

	"github.com/antchfx/jsonquery"
)

type JsonXpathDoc struct {
	node *jsonquery.Node
}

func NewJsonXpath(data []byte) *JsonXpathDoc {
	node, _ := jsonquery.Parse(bytes.NewReader(data))

	return &JsonXpathDoc{
		node: node,
	}
}

func (d *JsonXpathDoc) Value() interface{} {
	if d.node == nil {
		return nil
	}

	return d.node.Value()
}

func (d *JsonXpathDoc) FindOne(expr string) (Document, error) {
	node, err := jsonquery.Query(d.node, expr)

	if err != nil {
		return nil, err
	}

	return &JsonXpathDoc{
		node: node,
	}, nil
}

func (d *JsonXpathDoc) FindAll(expr string) []Document {
	nodes, err := jsonquery.QueryAll(d.node, expr)

	if err != nil {
		return nil
	}

	docs := make([]Document, 0, len(nodes))

	for _, node := range nodes {
		docs = append(docs, &JsonXpathDoc{
			node: node,
		})
	}

	return docs
}
