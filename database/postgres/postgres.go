package postgres

import (
	"database/sql"
	"errors"
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
	nativeTypeMap = api.NativeGraphQLTypeMap{
		Boolean: "BOOLEAN",
		Float:   "DECIMAL",
		Int:     "INT",
		String:  "TEXT",
	}
)

func (p *Postgres) AddModel(model api.Model) (sqlStmts []string, err error) {
	var relations []database.Relation
	var cols []string

	if exists, err := p.checkExistence(model.Name); err != nil {
		return sqlStmts, err
	} else {
		if exists {
			return sqlStmts, nil
		}
	}

	cols = append(cols, fmt.Sprintf("id %s PRIMARY KEY", relationContext.IdType))

	for _, f := range model.Fields {
		name, ft, nullable := f.Name, f.Type, f.Nullable

		var c string
		if !nullable {
			c = "NOT NULL"
		}

		f := reflect.ValueOf(nativeTypeMap).FieldByName(ft)
		// Not a native GraphQL type. Add it to the relations slice
		if !f.IsValid() {
			m := p.models.FindByName(ft)
			if m == nil {
				continue
			}

			relations = append(relations, database.Relation{
				From:      model,
				To:        *m,
				FieldName: name,
				Nullable:  nullable,
			})
		} else {
			c := fmt.Sprintf("%s %s %s", name, f.String(), c)
			cols = append(cols, strings.TrimSpace(c))
		}
	}

	ctBody := strings.Join(cols, ", ")
	ctStmt := fmt.Sprintf(`CREATE TABLE "%s" (%s)`, model.Name, ctBody)
	if _, err := p.db.Exec(ctStmt); err != nil {
		return sqlStmts, err
	}

	sqlStmts = append(sqlStmts, ctStmt)

	for _, r := range relations {
		// Attempt to create if model does not exist
		addModelStmts, err := p.AddModel(r.To)
		if err != nil {
			return sqlStmts, err
		}
		sqlStmts = append(sqlStmts, addModelStmts...)

		relStmts, err := p.createRelation(r)
		if err != nil {
			return sqlStmts, err
		}
		sqlStmts = append(sqlStmts, relStmts...)
	}

	return sqlStmts, nil
}

func (p *Postgres) RemoveModel(model api.Model) (sqlStmts []string, err error) {
	var deps api.ModelSet

	for _, f := range model.Fields {
		if f.IsNativeGraphQLType(nativeTypeMap) {
			continue
		}

		m := p.models.FindByName(f.Name)
		if m == nil {
			return sqlStmts, errors.New("no model found for foreign key")
		}
		deps = append(deps, *m)
	}

	for _, d := range deps {
		fk := p.toForeignKey(d.Name)
		dropStmt := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", d.Name, fk)

		sqlStmts = append(sqlStmts, dropStmt)
		if _, err := p.db.Exec(dropStmt); err != nil {
			return sqlStmts, err
		}
	}

	dropStmt := fmt.Sprintf("DROP TABLE %s", model.Name)
	sqlStmts = append(sqlStmts, dropStmt)

	if _, err := p.db.Exec(dropStmt); err != nil {
		return sqlStmts, err
	}

	return sqlStmts, nil
}

func (p *Postgres) createRelation(r database.Relation) ([]string, error) {
	sqlStmts := make([]string, 0)

	var constraint string
	if !r.Nullable {
		constraint = "NOT NULL"
	}

	colStmt := fmt.Sprintf(`ALTER TABLE "%s" ADD COLUMN "%sId" %s %s`, r.From.Name, r.FieldName, relationContext.IdRef, constraint)

	_, err := p.db.Exec(colStmt)
	if err != nil {
		return append(sqlStmts, colStmt), err
	}

	refC := fmt.Sprintf(`ADD CONSTRAINT "fk_%s_id" FOREIGN KEY (id) REFERENCES "%s" (id) ON UPDATE CASCADE ON DELETE CASCADE`, r.FieldName, r.To.Name)
	rStmt := fmt.Sprintf(`ALTER TABLE "%s" %s`, r.From.Name, refC)

	_, err = p.db.Exec(rStmt)
	return append(sqlStmts, colStmt, rStmt), err
}

func (p *Postgres) checkExistence(t string) (bool, error) {
	var exists bool
	r := p.db.QueryRow(fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = "public"" AND table_name = "%s")`, t))
	if err := r.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (p *Postgres) toForeignKey(col string) string {
	return fmt.Sprintf("fk_%s_id", col)
}

// New creates a new Postgres.
func New(db *sql.DB, models []api.Model) *Postgres {
	return &Postgres{db: db, models: models}
}

// Postgres is the PostgreSQL implementation.
type Postgres struct {
	db     *sql.DB
	models api.ModelSet
}
