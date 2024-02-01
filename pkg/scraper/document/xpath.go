package document

import (
	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

func getFindOne(node interface{}, expr string) (Document, error) {
	var err error
	var exp *xpath.Expr

	exp, err = xpath.Compile(expr)
	if err != nil {
		return nil, ErrInvalidQuery
	}
	navigator := getXpathNavigator(node)
	if navigator == nil {
		return nil, ErrNotExists
	}

	itemNode := exp.Evaluate(navigator)
	switch v := itemNode.(type) {
	case *xpath.NodeIterator:
		if v.MoveNext() {
			return getCurrentNode(v.Current())
		}
	case string:
		return newStringDoc(v), nil
	case float64:
		return newFloatDoc(v), nil
	}
	return nil, ErrNotExists
}

func getCurrentNode(iterator interface{}) (Document, error) {
	var doc Document
	switch n := iterator.(type) {
	case *htmlquery.NodeNavigator:
		return &HTMLXpathDoc{
			node: n.Current(),
		}, nil
	case *jsonquery.NodeNavigator:
		return &JSONXpathDoc{
			node: n.Current(),
		}, nil
	case *xmlquery.NodeNavigator:
		return &XMLXpathDoc{
			node: n.Current(),
		}, nil
	}
	return nil, ErrNotExists
}

func getXpathNavigator(node interface{}) xpath.NodeNavigator {
	switch n := node.(type) {
	case *xmlquery.Node:
		return xmlquery.CreateXPathNavigator(n)
	case *jsonquery.Node:
		return jsonquery.CreateXPathNavigator(n)
	case *html.Node:
		return htmlquery.CreateXPathNavigator(n)
	}

	return nil
}
