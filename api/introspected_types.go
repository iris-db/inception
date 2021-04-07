package api

import (
	"reflect"
)

// FieldSet is a collection of fields.
type FieldSet []Field

// Field is a Model attribute.
type Field struct {
	Name     string
	Type     string
	Nullable bool
}

// ModelSet is a collection of models.
type ModelSet []Model

// Model is a transformable blueprint that can be serialized into a database.
type Model struct {
	Name   string
	Fields FieldSet
}

// FindByName finds a model by name.
func (m ModelSet) FindByName(name string) *Model {
	for _, m := range m {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

// NativeGraphQLTypeMap maps a native GraphQL type to a native database type.
type NativeGraphQLTypeMap struct {
	Boolean string
	Float   string
	Int     string
	String  string
}

// IsNativeGraphQLType checks if a field is a native graphql type.
func (f Field) IsNativeGraphQLType(m NativeGraphQLTypeMap) bool {
	return reflect.ValueOf(m).FieldByName(f.Type).IsValid()
}
