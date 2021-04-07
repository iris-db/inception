package api_test

import (
	"github.com/web-foundation/sigma-production/api"
	"testing"
)

func TestField_IsNativeGraphQLType(t *testing.T) {
	tests := map[string]struct {
		typeName string
		want     bool
	}{
		"String is native":       {typeName: "String", want: true},
		"Int is native":          {typeName: "Int", want: true},
		"Boolean is native":      {typeName: "Boolean", want: true},
		"Float is native":        {typeName: "Float", want: true},
		"Settings is not native": {typeName: "Settings", want: false},
		"User is not native":     {typeName: "User", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := api.Field{Name: "testField", Type: tc.typeName, Nullable: false}
			if f.IsNativeGraphQLType(api.NativeGraphQLTypeMap{}) != tc.want {
				t.Fatalf("%s: expected true; got false", name)
			}
		})
	}
}
