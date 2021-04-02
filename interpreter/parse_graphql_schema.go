package interpreter

import (
	"encoding/json"
	"strings"

	"github.com/graph-gophers/graphql-go"
)

// ParseGQLSchema transforms a GraphQL schema into a slice of models.
func ParseGQLSchema(schema string) []Model {
	s, err := graphql.ParseSchema(addTempQuery(schema), nil)
	if err != nil {
		panic(err)
	}

	jsons, err := s.ToJSON()
	if err != nil {
		panic(err)
	}

	var is introspectionQuery
	if err := json.Unmarshal(jsons, &is); err != nil {
		panic(err)
	}

	ft := make([]introspectionType, 0)
	for _, t := range is.Schema.Types {
		if t.Kind == "OBJECT" && t.Name != "Query" && t.Name != "Mutation" && !strings.Contains(t.Name, "__") {
			ft = append(ft, t)
		}
	}

	st := make([]Model, 0)
	for _, t := range ft {
		fs := make([]Field, 0)
		for _, f := range t.Fields {
			var gqlType string

			fType := f.Type
			oType := fType.OfType
			if len(oType.Name) == 0 {
				gqlType = fType.Name
			} else {
				gqlType = oType.Name
			}

			fs = append(fs, Field{
				Name:     f.Name,
				Type:     gqlType,
				Nullable: f.Type.Kind != "NON_NULL",
			})
		}
		st = append(st, Model{
			Name:   t.Name,
			Fields: fs,
		})
	}
	return st
}

func addTempQuery(schema string) string {
	if !strings.Contains(schema, "type Query") {
		return schema + "\ntype Query { _temp: String! }"
	}
	return schema
}

type introspectionQuery struct {
	Schema struct {
		Types []introspectionType `json:"types"`
	} `json:"__schema"`
}

type introspectionType struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Fields []struct {
		Name string `json:"name"`
		Type struct {
			Kind   string `json:"kind"`
			Name   string `json:"name"`
			OfType struct {
				Name string `json:"name"`
			} `json:"ofType"`
		} `json:"type"`
	} `json:"fields"`
}
