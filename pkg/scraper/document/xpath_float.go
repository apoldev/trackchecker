package document

type floatDoc struct {
	value float64
}

func newFloatDoc(value float64) Document {
	return &floatDoc{value: value}
}

func (d *floatDoc) Value() interface{} {
	return d.value
}

func (d *floatDoc) FindOne(_ string) (Document, error) {
	return nil, ErrNotExists
}

func (d *floatDoc) FindAll(_ string) []Document {
	return nil
}
