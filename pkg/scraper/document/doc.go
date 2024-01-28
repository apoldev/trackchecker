package document

// Document can be HTML, JSON, XML, etc.
type Document interface {
	Value() interface{}
	FindOne(expr string) Document
	FindAll(expr string) []Document
}
