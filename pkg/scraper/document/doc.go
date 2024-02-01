package document

import "errors"

var (
	ErrNotExists = errors.New("not exist")
)

// Document can be HTML, JSON, XML, etc.
type Document interface {
	Value() interface{}
	FindOne(expr string) (Document, error)
	FindAll(expr string) []Document
}
