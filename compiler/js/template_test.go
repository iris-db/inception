package js_test

import (
	"github.com/web-foundation/sigma-production/compiler/js"
	"testing"
)

func TestParseTemplate(t *testing.T) {
	tests := map[string]struct {
		template string
		want     string
		values   js.TemplateValues
	}{
		"const value": {
			template: `function main() { var port = %APP_PORT%; }`,
			want:     `function main() { var port = 4000; }`,
			values:   js.TemplateValues{"APP_PORT": "4000"},
		},
		"inside string": {
			template: `const apiName = "Sigma::%API_NAME%";`,
			want:     `const apiName = "Sigma::AmazingApi";`,
			values:   js.TemplateValues{"API_NAME": "AmazingApi"},
		},
		"reoccurring single var": {
			template: `const apiName = "Sigma::%API_NAME%"; const serviceName = "%API_NAME%-Service"`,
			want:     `const apiName = "Sigma::AmazingApi"; const serviceName = "AmazingApi-Service"`,
			values:   js.TemplateValues{"API_NAME": "AmazingApi"},
		},
		"multiple vars": {
			template: `const apiName = "Sigma::%API_NAME%"; const serviceName = "%SERVICE_PREFIX%-Service"`,
			want:     `const apiName = "Sigma::AmazingApi"; const serviceName = "AmazingPrefix-Service"`,
			values:   js.TemplateValues{"API_NAME": "AmazingApi", "SERVICE_PREFIX": "AmazingPrefix"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := js.ParseTemplate(tc.template, tc.values)
			if res != tc.want {
				t.Fatalf("got %s; wanted %s", res, tc.want)
			}
		})
	}
}
