package document

type StringDoc struct {
	value string
}

func (d *StringDoc) Value() interface{} {
	return d.value
}

func (d *StringDoc) FindOne(_ string) (Document, error) {
	return nil, ErrNotExists
}

func (d *StringDoc) FindAll(_ string) []Document {
	return nil
}
