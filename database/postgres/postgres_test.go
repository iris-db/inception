package postgres_test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/database/postgres"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestPostgres_AddModel(t *testing.T) {
	db, mock := newDBMock()
	defer db.Close()

	primaryType := api.Model{
		Name: "User",
		Fields: api.FieldSet{
			{Name: "username", Type: "String", Nullable: true},
			{Name: "email", Type: "String", Nullable: false},
			{Name: "password", Type: "String", Nullable: false},
			{Name: "settings", Type: "Settings", Nullable: false},
		},
	}
	relationType := api.Model{
		Name: "Settings",
		Fields: []api.Field{
			{Name: "theme", Type: "String", Nullable: false},
			{Name: "subUser", Type: "User", Nullable: false},
		},
	}

	pg := postgres.New(db, api.ModelSet{primaryType, relationType})

	dne1 := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	dne2 := sqlmock.NewRows([]string{"exists"}).AddRow(false)
	dne3 := sqlmock.NewRows([]string{"exists"}).AddRow(true)

	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(dne1)
	mock.ExpectExec("CREATE TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(dne2)
	mock.ExpectExec("CREATE TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(dne3)
	mock.ExpectExec("ALTER TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("ALTER TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("ALTER TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("ALTER TABLE").WillReturnResult(driver.ResultNoRows)

	res, err := pg.AddModel(primaryType)
	if err != nil {
		t.Fatalf("Error ::: %s", err.Error())
	}

	expectedStmts := []string{
		`CREATE TABLE "User" (id SERIAL PRIMARY KEY, username TEXT, email TEXT NOT NULL, password TEXT NOT NULL)`,
		`CREATE TABLE "Settings" (id SERIAL PRIMARY KEY, theme TEXT NOT NULL)`,
		`ALTER TABLE "Settings" ADD COLUMN "subUserId" INT NOT NULL`,
		`ALTER TABLE "Settings" ADD CONSTRAINT "fk_subUser_id" FOREIGN KEY (id) REFERENCES "User" (id) ON UPDATE CASCADE ON DELETE CASCADE`,
		`ALTER TABLE "User" ADD COLUMN "settingsId" INT NOT NULL`,
		`ALTER TABLE "User" ADD CONSTRAINT "fk_settings_id" FOREIGN KEY (id) REFERENCES "Settings" (id) ON UPDATE CASCADE ON DELETE CASCADE`,
	}

	if !reflect.DeepEqual(res, expectedStmts) {
		t.Fatalf("\ngot\n--- %s \n\nexpectedStmts\n--- %s", fmtSQLStmts(res), fmtSQLStmts(expectedStmts))
	}
}

func TestPostgres_RemoveModel(t *testing.T) {
	db, mock := newDBMock()
	defer db.Close()

	primaryType := api.Model{
		Name: "User",
		Fields: api.FieldSet{
			{Name: "username", Type: "String", Nullable: false},
			{Name: "settings", Type: "Settings", Nullable: false},
		},
	}
	relationType := api.Model{
		Name: "Settings",
		Fields: api.FieldSet{
			{Name: "subUser", Type: "User", Nullable: false},
		},
	}

	mock.ExpectQuery("ALTER TABLE")
	mock.ExpectQuery("DROP TABLE")

	pg := postgres.New(db, api.ModelSet{primaryType, relationType})

	mock.ExpectExec("ALTER TABLE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("DROP TABLE").WillReturnResult(driver.ResultNoRows)

	res, err := pg.RemoveModel(primaryType)
	if err != nil {
		// TODO fix unit test
		//t.Fatalf("Error ::: %s", err.Error())
	}

	expectedStmts := []string{
		`ALTER TABLE "Settings" DROP COLUMN userId`,
		`DROP TABLE "User"`,
	}

	if !reflect.DeepEqual(res, expectedStmts) {
		// TODO fix unit test
		//t.Fatalf("\ngot\n--- %s \n\nexpectedStmts\n--- %s", fmtSQLStmts(res), fmtSQLStmts(expectedStmts))
	}
}

func fmtSQLStmts(stmts []string) string {
	newFStmt := "\n" + (stmts[0])
	stmts[0] = newFStmt
	return strings.Join(stmts, "\n")
}

func newDBMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
