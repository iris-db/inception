package main

import (
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/compiler/js/rest"
)

func main() {
	rest.CompileAPI(rest.CompilationOpts{
		APIOpts: rest.APIOpts{
			Name:   "Test",
			Prefix: "api",
			Type:   "rest",
			Port:   "4000",
			Models: api.ModelSet{
				{
					Name:   "User",
					Fields: api.FieldSet{},
				},
			},
		},
	})
}
