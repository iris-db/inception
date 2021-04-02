package interpreter

// Model is a transformable blueprint that can be serialized into a database.
type Model struct {
	Name   string
	Fields []Field
}

// Field is a Model attribute.
type Field struct {
	Name     string
	Type     string
	Nullable bool
}
