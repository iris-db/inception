package api

import (
	"reflect"
)

// FindByName finds a model by name.
func (m ModelSet) FindByName(name string) *Model {
	for _, m := range m {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

// IsNativeGraphQLType checks if a field is a native graphql type.
func (f Field) IsNativeGraphQLType(m NativeGraphQLTypeMap) bool {
	return reflect.ValueOf(m).FieldByName(f.Type).IsValid()
}

// ModelSet is a collection of models.
type ModelSet []Model

// FieldSet is a collection of fields.
type FieldSet []Field

// Model is a transformable blueprint that can be serialized into a database.
type Model struct {
	Name   string
	Fields FieldSet
}

// Field is a Model attribute.
type Field struct {
	Name     string
	Type     string
	Nullable bool
}

// NativeGraphQLTypeMap maps a native GraphQL type to a native database type.
type NativeGraphQLTypeMap struct {
	Boolean string
	Float   string
	Int     string
	String  string
}
