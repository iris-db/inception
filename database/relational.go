package database

import (
	"github.com/web-foundation/sigma-production/api"
)

// NativeGraphQLTypeMap maps a native GraphQL type to a native database type.
type NativeGraphQLTypeMap struct {
	Boolean string
	Float   string
	Int     string
	String  string
}

// RelationContext describes the type names that are required to form
// relations across tables.
type RelationContext struct {
	IdType string
	IdRef  string
}

// RelationalDatabaseController initializes and destructs api.Model relational
// database tables. Returns a slice of strings that represent the executed sql
// statements or an error if not all of the sql statements were able to
// complete successfully.
type RelationalDatabaseController interface {
	AddModel(target api.Model, models []api.Model) ([]string, error)
	CreateRelation(r Relation) ([]string, error)
	RemoveModel(model api.Model) error
	CheckExistence(t string) (bool, error)
}

// Operation is an operation that can be performed on a supported
// database.
type Operation string

const (
	SyncModel   Operation = "SYNC_MODEL"
	DeleteModel           = "DELETE_MODEL"
)

func (d *ExecutionError) Error() string {
	return "error while performing operation " + string(d.Operation) + " on model " + d.Model.Name
}

// ExecutionError is an error that occurs when a RelationalDatabaseController is
// unable to successfully complete a task.
type ExecutionError struct {
	Operation Operation
	Model     api.Model
	Err       error
}

// Relation describes a dependency between two models.
type Relation struct {
	From      api.Model
	To        api.Model
	FieldName string
	Nullable  bool
}
