package document

type StringDoc struct {
	value string
}

func (d *StringDoc) Value() interface{} {
	return d.value
}

func (d *StringDoc) FindOne(expr string) Document {
	return nil
}

func (d *StringDoc) FindAll(expr string) []Document {
	return nil
}
