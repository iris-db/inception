package remote

import "sigma-production/interpreter"

type createTableOutput struct {
	SQLStmts    []string
	ForeignKeys []Relation
}

// Relation describes a dependency between two models.
type Relation struct {
	From   			interpreter.Model
	To 				interpreter.Model
	FieldName       string
	Nullable        bool
}
