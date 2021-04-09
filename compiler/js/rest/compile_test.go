package rest_test

import (
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/compiler/js/rest"
	"testing"
)

func TestTemplateGenerator_Get(t *testing.T) {
	tests := map[string]struct {
		model api.Model
		want  string
	}{
		"query params": {
			model: api.Model{
				Name: "User",
				Fields: api.FieldSet{
					{Name: "username", Type: "String", Nullable: true},
					{Name: "email", Type: "String", Nullable: false},
					{Name: "password", Type: "String", Nullable: false},
				},
			},
			want: `
router.get("/", async (req, res) => {
	const users = await UserController.findAll(req.query);
	return res.json(users);
})
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			g := rest.CrudGenerator{Model: tc.model}

			template := g.Get()

			if template != tc.want[1:] {
				t.Fatalf("got \n%s; wanted \n%s", template, tc.want)
			}
		})
	}
}

func BenchmarkCompileAPI(b *testing.B) {
	opts := rest.CompilationOpts{
		APIOpts: rest.APIOpts{
			Name:   "MyApi",
			Prefix: " ",
			Type:   "rest",
			Port:   "4000",
		},
	}

	for n := 0; n < b.N; n++ {
		rest.CompileAPI(opts)
	}
}
