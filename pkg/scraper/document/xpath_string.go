package document

type StringDoc struct {
	value string
}

func (d *StringDoc) Value() interface{} {
	return d.value
}

func (d *StringDoc) FindOne(_ string) Document {
	return nil
}

func (d *StringDoc) FindAll(_ string) []Document {
	return nil
}
