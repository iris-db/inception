package remote

import (
	"sigma-production/interpreter"
)

// NativeGraphQLTypeMap maps a native GraphQL type to a native database type.
type NativeGraphQLTypeMap struct {
	Boolean string
	Float   string
	Int     string
	String  string
}

// DatabaseRelationContext describes the type names that are required to form
// relations across tables.
type DatabaseRelationContext struct {
	IdType string
	IdRef  string
}

// RelationalDatabaseController initializes and destructs interpreter.Model relational
// database tables. Returns a slice of strings that represent the executed sql
// statements or an error if not all of the sql statements were able to
// complete successfully.
type RelationalDatabaseController interface {
	AddModel(target interpreter.Model, models []interpreter.Model) ([]string, error)
	CreateRelation(r Relation) ([]string, error)
	RemoveModel(model interpreter.Model) error
	CheckExistence(t string) (bool, error)
}

// DatabaseOperation is an operation that can be performed on a supported
// database.
type DatabaseOperation string

const (
	SyncModel   DatabaseOperation = "SYNC_MODEL"
	DeleteModel                   = "DELETE_MODEL"
)

func (d *DatabaseExecutionError) Error() string {
	return "error while performing operation " + string(d.Operation) + " on model " + d.Model.Name
}

// DatabaseExecutionError is an error that occurs when a RelationalDatabaseController is
// unable to successfully complete a task.
type DatabaseExecutionError struct {
	Operation DatabaseOperation
	Model     interpreter.Model
	Err       error
}
