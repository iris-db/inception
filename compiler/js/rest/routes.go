package rest

import (
	_ "embed"
	"github.com/web-foundation/sigma-production/api"
	"github.com/web-foundation/sigma-production/compiler"
)

// CrudGenerator generates crud rest endpoints from a model.
type CrudGenerator struct {
	Model   api.Model
	Methods [4]string
}

//go:embed templates/crud/get.txt
var getTemplate string

func (t CrudGenerator) Get() string {
	tv := compiler.TemplateValues{
		"MODEL_NAME": compiler.StrPtr(t.Model.Name),
	}
	return compiler.ParseTemplate(getTemplate, tv)
}
