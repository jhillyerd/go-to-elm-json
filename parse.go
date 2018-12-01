package main

import (
	"go/types"
	"strings"

	"github.com/pkg/errors"
)

// Module represents an Elm module.
type Module struct {
	Name   string
	Fields []Field
}

// Field represents a Go -> JSON -> Elm field.
type Field struct {
	GoName     string
	GoType     string
	JSONName   string
	ElmName    string
	ElmType    string
	ElmDecoder string
}

func parseStructType(typeName string, root *types.Struct) (*Module, error) {
	count := root.NumFields()
	if count == 0 {
		return nil, errors.Errorf("struct %v had no fields", typeName)
	}

	// Convert to our field type.
	var fields []Field
	for i := 0; i < root.NumFields(); i++ {
		sfield := root.Field(i)
		stag := root.Tag(i)
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
		elmType, elmDecoder := elmType(goType)
		fields = append(fields, Field{
			GoName:     goName,
			GoType:     goType.String(),
			JSONName:   jsonName,
			ElmName:    elmName,
			ElmType:    elmType,
			ElmDecoder: elmDecoder,
		})
	}

	return &Module{Name: typeName, Fields: fields}, nil
}
