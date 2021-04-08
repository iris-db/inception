package compiler_test

import (
	"github.com/web-foundation/sigma-production/compiler"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	tests := map[string]struct {
		template string
		want     string
		values   compiler.TemplateValues
	}{
		"const value": {
			template: `function main() { var port = %APP_PORT%; }`,
			want:     `function main() { var port = 4000; }`,
			values:   compiler.TemplateValues{"APP_PORT": compiler.StrPtr("4000")},
		},
		"inside string": {
			template: `const apiName = "Sigma::%API_NAME%";`,
			want:     `const apiName = "Sigma::AmazingApi";`,
			values:   compiler.TemplateValues{"API_NAME": compiler.StrPtr("AmazingApi")},
		},
		"reoccurring single var": {
			template: `const apiName = "Sigma::%API_NAME%"; const serviceName = "%API_NAME%-Service"`,
			want:     `const apiName = "Sigma::AmazingApi"; const serviceName = "AmazingApi-Service"`,
			values:   compiler.TemplateValues{"API_NAME": compiler.StrPtr("AmazingApi")},
		},
		"multiple vars": {
			template: `const apiName = "Sigma::%API_NAME%"; const serviceName = "%SERVICE_PREFIX%-Service"`,
			want:     `const apiName = "Sigma::AmazingApi"; const serviceName = "AmazingPrefix-Service"`,
			values:   compiler.TemplateValues{"API_NAME": compiler.StrPtr("AmazingApi"), "SERVICE_PREFIX": compiler.StrPtr("AmazingPrefix")},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := compiler.ParseTemplate(tc.template, tc.values)
			if res != tc.want {
				t.Fatalf("got %s; wanted %s", res, tc.want)
			}
		})
	}
}
