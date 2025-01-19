package document

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

type HTMLDoc struct {
	selection *goquery.Selection
}

func NewHTML(data []byte) (*HTMLDoc, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	return &HTMLDoc{
		selection: doc.Selection,
	}, nil
}

func (d *HTMLDoc) Value() interface{} {
	if d.selection == nil {
		return nil
	}

	return d.selection.First().Text()
}

func (d *HTMLDoc) FindOne(expr string) (Document, error) {
	selection := d.selection.Find(expr)

	return &HTMLDoc{
		selection: selection,
	}, nil
}

func (d *HTMLDoc) FindAll(expr string) []Document {
	selection := d.selection.Find(expr)

	docs := make([]Document, 0, selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		docs = append(docs, &HTMLDoc{
			selection: selection,
		})
	})

	return docs
}
