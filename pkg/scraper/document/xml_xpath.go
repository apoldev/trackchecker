package document

import (
	"bytes"

	"github.com/antchfx/xmlquery"
)

type XMLXpathDoc struct {
	node *xmlquery.Node
}

func NewXMLXpath(data []byte) (*XMLXpathDoc, error) {
	node, err := xmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &XMLXpathDoc{
		node: node,
	}, nil
}

func (d *XMLXpathDoc) Value() interface{} {
	if d.node == nil {
		return nil
	}
	return d.node.InnerText()
}

func (d *XMLXpathDoc) FindOne(expr string) (Document, error) {
	return getFindOne(d.node, expr)
}

func (d *XMLXpathDoc) FindAll(expr string) []Document {
	nodes, err := xmlquery.QueryAll(d.node, expr)
	if err != nil {
		return nil
	}
	docs := make([]Document, 0, len(nodes))
	for _, node := range nodes {
		docs = append(docs, &XMLXpathDoc{
			node: node,
		})
	}
	return docs
}
