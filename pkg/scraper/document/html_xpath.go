package document

import (
	"bytes"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

type HTMLXpathDoc struct {
	node *html.Node
}

func NewHTMLXpath(data []byte) (*HTMLXpathDoc, error) {
	node, err := htmlquery.Parse(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	return &HTMLXpathDoc{
		node: node,
	}, nil
}

func (d *HTMLXpathDoc) Value() interface{} {
	if d.node == nil {
		return nil
	}

	return htmlquery.InnerText(d.node)
}

func (d *HTMLXpathDoc) FindOne(expr string) (Document, error) {
	var err error
	var exp *xpath.Expr

	exp, err = xpath.Compile(expr)

	if err != nil {
		return nil, err
	}

	navigator := htmlquery.CreateXPathNavigator(d.node)

	// h := GetCurrentNodeFromNavigator(navigator)
	// fmt.Println(expr, "navigator", string(HtmlNodeToBytes(*h)))

	itemNode := exp.Evaluate(navigator)

	switch v := itemNode.(type) {
	case *xpath.NodeIterator:

		iterator := v
		iterator.MoveNext()

		if v, ok := iterator.Current().(*htmlquery.NodeNavigator); ok {
			return &HTMLXpathDoc{
				node: v.Current(),
			}, nil
		}

	case string:

		return &StringDoc{
			value: v,
		}, nil

		// todo bool, float64
	}

	return nil, ErrNotExists
}

func (d *HTMLXpathDoc) FindAll(expr string) []Document {
	nodes, err := htmlquery.QueryAll(d.node, expr)

	if err != nil {
		return nil
	}

	docs := make([]Document, 0, len(nodes))

	for _, node := range nodes {
		docs = append(docs, &HTMLXpathDoc{
			node: node,
		})
	}

	return docs
}

func GetCurrentNodeFromNavigator(n *htmlquery.NodeNavigator) *html.Node {
	if n.NodeType() == xpath.AttributeNode {
		childNode := &html.Node{
			Type: html.TextNode,
			Data: n.Value(),
		}
		return &html.Node{
			Type:       html.ElementNode,
			Data:       n.LocalName(),
			FirstChild: childNode,
			LastChild:  childNode,
		}
	}
	return n.Current()
}

func HTMLNodeToBytes(node html.Node) []byte {
	var b bytes.Buffer

	err := html.Render(&b, &node)

	if err == nil {
		return b.Bytes()
	}

	return nil
}
