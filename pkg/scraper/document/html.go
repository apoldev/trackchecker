package document

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

type HtmlDoc struct {
	selection *goquery.Selection
}

func NewHtml(data []byte) *HtmlDoc {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(data))

	return &HtmlDoc{
		selection: doc.Selection,
	}
}

func (d *HtmlDoc) Value() interface{} {
	if d.selection == nil {
		return nil
	}

	return d.selection.First().Text()
}

func (d *HtmlDoc) FindOne(expr string) (Document, error) {
	selection := d.selection.Find(expr)

	return &HtmlDoc{
		selection: selection,
	}, nil
}

func (d *HtmlDoc) FindAll(expr string) []Document {
	selection := d.selection.Find(expr)

	docs := make([]Document, 0, selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		docs = append(docs, &HtmlDoc{
			selection: selection,
		})
	})

	return docs
}
