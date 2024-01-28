package scraper

const (
	FieldTypeCommon = ""
	FieldTypeObject = "object"
	FieldTypeArray  = "array"
)

// Field describes a field in result.
type Field struct {
	Path string `json:"path,omitempty"`
	Type string `json:"type,omitempty"`

	Query string `json:"query,omitempty"` // XPath or JSONPath or CSS selector

	// Element is used for FieldTypeArray
	Element *Field `json:"element,omitempty"`

	// Object is used for FieldTypeObject
	Object []*Field `json:"object,omitempty"`
}
