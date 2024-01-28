package document

import (
	"bytes"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

type HtmlXpathDoc struct {
	node *html.Node
}

func NewHtmlXpath(data []byte) *HtmlXpathDoc {
	node, _ := htmlquery.Parse(bytes.NewReader(data))

	return &HtmlXpathDoc{
		node: node,
	}
}

func (d *HtmlXpathDoc) Value() interface{} {
	if d.node == nil {
		return nil
	}

	return htmlquery.InnerText(d.node)
}

func (d *HtmlXpathDoc) FindOne(expr string) (Document, error) {
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

	switch itemNode.(type) {
	case *xpath.NodeIterator:

		iterator := itemNode.(*xpath.NodeIterator)
		iterator.MoveNext()

		if v, ok := iterator.Current().(*htmlquery.NodeNavigator); ok {
			return &HtmlXpathDoc{
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

func (d *HtmlXpathDoc) FindAll(expr string) []Document {
	nodes, err := htmlquery.QueryAll(d.node, expr)

	if err != nil {
		return nil
	}

	docs := make([]Document, 0, len(nodes))

	for _, node := range nodes {
		docs = append(docs, &HtmlXpathDoc{
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

func HtmlNodeToBytes(node html.Node) []byte {
	var b bytes.Buffer

	err := html.Render(&b, &node)

	if err == nil {
		return b.Bytes()
	}

	return nil
}
