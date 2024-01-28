package document

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/apoldev/trackchecker/pkg/scraper"
)

type HtmlDoc struct {
	selection *goquery.Selection
}

func NewHtml(data []byte) scraper.Document {

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

func (d *HtmlDoc) FindOne(expr string) scraper.Document {

	selection := d.selection.Find(expr)

	return &HtmlDoc{
		selection: selection,
	}
}

func (d *HtmlDoc) FindAll(expr string) []scraper.Document {

	selection := d.selection.Find(expr)

	docs := make([]scraper.Document, 0, selection.Length())

	selection.Each(func(i int, selection *goquery.Selection) {
		docs = append(docs, &HtmlDoc{
			selection: selection,
		})
	})

	return docs
}
