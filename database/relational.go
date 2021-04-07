package database

import (
	"github.com/web-foundation/sigma-production/api"
)

// Relation describes a dependency between two models.
type Relation struct {
	From      api.Model
	To        api.Model
	FieldName string
	Nullable  bool
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
	AddModel(model api.Model) (sqlStmts []string, err error)
	RemoveModel(model api.Model) (sqlStmts []string, err error)
}

// Operation is an operation that can be performed on a supported
// database.
type Operation string

const (
	SyncModel   Operation = "SYNC_MODEL"
	DeleteModel           = "DELETE_MODEL"
)

// ExecutionError is an error that occurs when a RelationalDatabaseController is
// unable to successfully complete a task.
type ExecutionError struct {
	Operation Operation
	Model     api.Model
	Err       error
}

func (d *ExecutionError) Error() string {
	return "error while performing operation " + string(d.Operation) + " on model " + d.Model.Name
}
