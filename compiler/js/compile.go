package js

import (
	"fmt"
	"github.com/web-foundation/sigma-production/api"
)

// MakeAPI makes the api.
func MakeAPI() {
	models := api.ModelSet{
		api.Model{
			Name: "User",
			Fields: api.FieldSet{
				{Name: "username", Type: "String", Nullable: false},
			},
		},
	}

	fmt.Print(models)

	initProject()
}
