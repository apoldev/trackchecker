package document

import (
	"bytes"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
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
	var err error
	var exp *xpath.Expr

	exp, err = xpath.Compile(expr)
	if err != nil {
		return nil, err
	}

	navigator := xmlquery.CreateXPathNavigator(d.node)
	itemNode := exp.Evaluate(navigator)

	switch itemNode.(type) {
	case *xpath.NodeIterator:
		iterator := itemNode.(*xpath.NodeIterator)
		iterator.MoveNext()
		if v, ok := iterator.Current().(*xmlquery.NodeNavigator); ok {
			return &XMLXpathDoc{
				node: v.Current(),
			}, nil
		}
	case string:
		return &StringDoc{
			value: itemNode.(string),
		}, nil
		// todo bool, float64
	}

	return nil, ErrorNotexist
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
