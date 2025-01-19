package document

type stringDoc struct {
	value string
}

func newStringDoc(value string) Document {
	return &stringDoc{value: value}
}

func (d *stringDoc) Value() interface{} {
	return d.value
}

func (d *stringDoc) FindOne(_ string) (Document, error) {
	return nil, ErrNotExists
}

func (d *stringDoc) FindAll(_ string) []Document {
	return nil
}
