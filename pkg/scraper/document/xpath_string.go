package document

import "github.com/apoldev/trackchecker/pkg/scraper"

type StringDoc struct {
	value string
}

func (d *StringDoc) Value() interface{} {
	return d.value
}

func (d *StringDoc) FindOne(expr string) scraper.Document {
	return nil
}

func (d *StringDoc) FindAll(expr string) []scraper.Document {
	return nil
}
