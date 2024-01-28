package document

import "errors"

var (
	ErrorNotexist = errors.New("not exist")
)

// Document can be HTML, JSON, XML, etc.
type Document interface {
	Value() interface{}
	FindOne(expr string) (Document, error)
	FindAll(expr string) []Document
}
