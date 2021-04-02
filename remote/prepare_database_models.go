package remote

import (
	"sigma-production/interpreter"
)

// DatabaseController initializes and destructs interpreter.Model relational
// database tables. Returns a slice of strings that represent the executed sql
// statements or an error if not all of the sql statements were able to
// complete successfully.
type DatabaseController interface {
	CreateTable(target interpreter.Model, models []interpreter.Model) ([]string, error)
	DeleteTable(model interpreter.Model) error
}

// DatabaseOperation is an operation that can be performed on a supported
// database.
type DatabaseOperation string

const (
	CreateTable DatabaseOperation = "CREATE_TABLE"
	DeleteTable                   = "DELETE_TABLE"
)

func (d *DatabaseExecutionError) Error() string {
	return "error while performing operation " + string(d.Operation) + " on model " + d.Model.Name
}

// DatabaseExecutionError is an error that occurs when a DatabaseController is
// unable to successfully complete a task.
type DatabaseExecutionError struct {
	Operation DatabaseOperation
	Model     interpreter.Model
	Err       error
}
