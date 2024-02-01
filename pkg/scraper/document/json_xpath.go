package document

import (
	"bytes"

	"github.com/antchfx/jsonquery"
)

type JSONXpathDoc struct {
	node *jsonquery.Node
}

func NewJSONXpath(data []byte) (*JSONXpathDoc, error) {
	node, err := jsonquery.Parse(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	return &JSONXpathDoc{
		node: node,
	}, nil
}

func (d *JSONXpathDoc) Value() interface{} {
	if d.node == nil {
		return nil
	}

	return d.node.Value()
}

func (d *JSONXpathDoc) FindOne(expr string) (Document, error) {
	node, err := jsonquery.Query(d.node, expr)

	if err != nil {
		return nil, err
	}

	return &JSONXpathDoc{
		node: node,
	}, nil
}

func (d *JSONXpathDoc) FindAll(expr string) []Document {
	nodes, err := jsonquery.QueryAll(d.node, expr)

	if err != nil {
		return nil
	}

	docs := make([]Document, 0, len(nodes))

	for _, node := range nodes {
		docs = append(docs, &JSONXpathDoc{
			node: node,
		})
	}

	return docs
}
