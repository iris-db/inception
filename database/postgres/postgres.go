package postgres

import (
	"database/sql"
	"fmt"
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/database"
	"reflect"
	"strings"

	// PostgreSQL driver.
	_ "github.com/lib/pq"
)

var (
	relationContext = database.RelationContext{
		IdType: "SERIAL",
		IdRef:  "INT",
	}
	nativeTypeMap = database.NativeGraphQLTypeMap{
		Boolean: "BOOLEAN",
		Float:   "DECIMAL",
		Int:     "INT",
		String:  "TEXT",
	}
)

func (p Postgres) AddModel(target api.Model, models []api.Model) ([]string, error) {
	fks := make([]database.Relation, 0)
	sqlStmts := make([]string, 0)

	if exists, err := p.checkExistence(target.Name); err != nil {
		return sqlStmts, err
	} else {
		if exists {
			return sqlStmts, nil
		}
	}

	cols := make([]string, 0)
	cols = append(cols, fmt.Sprintf("id %s PRIMARY KEY", relationContext.IdType))

	for _, f := range target.Fields {
		name, ft, nullable := f.Name, f.Type, f.Nullable

		var c string
		if !nullable {
			c = "NOT NULL"
		}

		f := reflect.ValueOf(nativeTypeMap).FieldByName(ft)
		if !f.IsValid() {
			var t *api.Model
			for _, m := range models {
				if m.Name == ft {
					t = &m
					break
				}
			}

			if t == nil {
				continue
			}

			fks = append(fks, database.Relation{
				From:      target,
				To:        *t,
				FieldName: name,
				Nullable:  nullable,
			})
		} else {
			c := fmt.Sprintf("%s %s %s", name, f.String(), c)
			cols = append(cols, strings.TrimSpace(c))
		}
	}

	ctBody := strings.Join(cols, ", ")
	ctStmt := fmt.Sprintf(`CREATE TABLE "%s" (%s)`, target.Name, ctBody)
	if _, err := p.db.Exec(ctStmt); err != nil {
		return sqlStmts, err
	}

	sqlStmts = append(sqlStmts, ctStmt)

	for _, r := range fks {
		addModelStmts, err := p.AddModel(r.To, models)
		if err != nil {
			return sqlStmts, err
		}
		sqlStmts = append(sqlStmts, addModelStmts...)

		relStmts, err := p.CreateRelation(r)
		if err != nil {
			return sqlStmts, err
		}
		sqlStmts = append(sqlStmts, relStmts...)
	}

	return sqlStmts, nil
}

func (p Postgres) RemoveModel(model api.Model) error {
	panic("implement me")
}

func (p Postgres) CreateRelation(r database.Relation) ([]string, error) {
	sqlStmts := make([]string, 0)

	var colC string
	if !r.Nullable {
		colC = "NOT NULL"
	}

	colStmt := fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN "%sId" %s %s`, r.From.Name, r.FieldName, relationContext.IdRef, colC)

	_, err := p.db.Exec(colStmt)
	if err != nil {
		return append(sqlStmts, colStmt), err
	}

	refC := fmt.Sprintf(`ADD CONSTRAINT "fk_%s_id" FOREIGN KEY (id) REFERENCES "%s" (id) ON UPDATE CASCADE ON DELETE CASCADE`, r.FieldName, r.To.Name)
	rStmt := fmt.Sprintf(`ALTER TABLE "%s" %s`, r.From.Name, refC)

	_, err = p.db.Exec(rStmt)
	return append(sqlStmts, colStmt, rStmt), err
}

func (p Postgres) checkExistence(t string) (bool, error) {
	var exists bool
	r := p.db.QueryRow(fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s')`, t))
	if err := r.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// New creates a new Postgres.
func New(db *sql.DB) *Postgres {
	return &Postgres{db: db}
}

// Postgres is the PostgreSQL implementation.
type Postgres struct {
	db *sql.DB
}
