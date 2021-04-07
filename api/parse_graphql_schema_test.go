package api_test

import (
	"github.com/web-foundation/sigma-production/api"
	"reflect"
	"testing"
)

func TestParseSchema(t *testing.T) {
	schema := `
		type User {
			id: ID!
			username: String
			email: String!
			password: String!
			settings: Settings!
		}
		
		type Settings {
			id: ID!
			theme: String!
		}
`
	expectedResult := []api.Model{
		{
			Name: "Settings",
			Fields: []api.Field{
				{Name: "id", Type: "ID", Nullable: false},
				{Name: "theme", Type: "String", Nullable: false},
			},
		},
		{
			Name: "User",
			Fields: []api.Field{
				{Name: "id", Type: "ID", Nullable: false},
				{Name: "username", Type: "String", Nullable: true},
				{Name: "email", Type: "String", Nullable: false},
				{Name: "password", Type: "String", Nullable: false},
				{Name: "settings", Type: "Settings", Nullable: false},
			},
		},
	}

	res := api.ParseGQLSchema(schema)

	if !reflect.DeepEqual(res, expectedResult) {
		t.Errorf("Not matching. Wanted %+v got %+v", expectedResult, res)
	}
}
