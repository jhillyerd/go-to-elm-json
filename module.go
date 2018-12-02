package main

import (
	"go/types"
	"strings"

	"github.com/pkg/errors"
)

// Module represents an Elm module.
type Module struct {
	Name   string
	Fields []*Field
}

// Equal tests for equality with another Module.
func (m *Module) Equal(o *Module) bool {
	if m.Name != o.Name {
		return false
	}
	if len(m.Fields) != len(o.Fields) {
		return false
	}
	for i, f := range m.Fields {
		if !f.Equal(o.Fields[i]) {
			return false
		}
	}
	return true
}

// Field represents a Go -> JSON -> Elm field.
type Field struct {
	JSONName string
	ElmName  string
	ElmType  ElmType
}

// Equal test for equality with another field.
func (f *Field) Equal(o *Field) bool {
	return f.JSONName == o.JSONName &&
		f.ElmName == o.ElmName &&
		f.ElmType.Equal(o.ElmType)
}

func moduleFromStruct(structDef *types.Struct, moduleName string) (*Module, error) {
	count := structDef.NumFields()
	if count == 0 {
		return nil, errors.Errorf("struct %v had no fields", moduleName)
	}

	// Convert to our field type.
	var fields []*Field
	for i := 0; i < structDef.NumFields(); i++ {
		sfield := structDef.Field(i)
		stag := structDef.Tag(i)
		if !sfield.Exported() {
			continue
		}
		goName := sfield.Name()
		goType := sfield.Type()

		jsonName := goName
		if len(stag) > 2 {
			tagName, _ := parseTag(stag)
			jsonName = tagName
		}

		// Handle abbrevations.
		camelCaseName := camelCase(goName)
		elmName := strings.ToLower(camelCaseName[0:1]) + camelCaseName[1:]
		elmType := elmType(goType)
		fields = append(fields, &Field{
			JSONName: jsonName,
			ElmName:  elmName,
			ElmType:  elmType,
		})
	}

	return &Module{Name: moduleName, Fields: fields}, nil
}
