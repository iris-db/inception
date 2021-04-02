package remote

import (
	"sigma-production/interpreter"
)

var (
	relationContext = DatabaseRelationContext{
		IdType: "SERIAL",
		IdRef:  "INT",
	}
)

// CreateTable creates a new table.
func (p Postgres) CreateTable(target interpreter.Model, models []interpreter.Model) ([]string, error) {
	panic("implement me")
}

// DeleteTable deletes a table.
func (p Postgres) DeleteTable(model interpreter.Model) error {
	panic("implement me")
}

// Postgres is the implementation for the PostgreSQL database.
type Postgres struct{}
