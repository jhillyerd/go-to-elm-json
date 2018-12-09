package main

import (
	"go/types"
	"strings"

	"github.com/pkg/errors"
)

// ElmRecord represents an Elm record.
type ElmRecord struct {
	name   string
	Fields []*ElmField
}

// Name of this record type.
func (r *ElmRecord) Name() string {
	return r.name
}

// Codec for this record type.
func (r *ElmRecord) Codec(prefix string) string {
	return "dunno"
}

// Equal tests for equality with another ElmType.
func (r *ElmRecord) Equal(other ElmType) bool {
	if o, ok := other.(*ElmRecord); ok {
		if r.name != o.name {
			return false
		}
		if len(r.Fields) != len(o.Fields) {
			return false
		}
		for i, f := range r.Fields {
			if !f.Equal(o.Fields[i]) {
				return false
			}
		}
		return true
	}
	return false
}

// ElmField represents an Elm record field.
type ElmField struct {
	JSONName string
	ElmName  string
	ElmType  ElmType
}

// Equal test for equality with another field.
func (f *ElmField) Equal(o *ElmField) bool {
	return f.JSONName == o.JSONName &&
		f.ElmName == o.ElmName &&
		f.ElmType != nil &&
		f.ElmType.Equal(o.ElmType)
}

func recordFromStruct(resolver *ElmTypeResolver, structDef *types.Struct, typeName string) (*ElmRecord, error) {
	count := structDef.NumFields()
	if count == 0 {
		return nil, errors.Errorf("struct %v had no fields", typeName)
	}
	camelCaseName := camelCase(typeName)
	recordName := strings.ToUpper(camelCaseName[:1]) + camelCaseName[1:]

	// Convert to our field type.
	var fields []*ElmField
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
		elmName := strings.ToLower(camelCaseName[:1]) + camelCaseName[1:]
		elmType, err := resolver.Convert(goType)
		if err != nil {
			return nil, err
		}
		logger.Debug().
			Str("field", recordName+":"+jsonName).
			Str("goType", goType.String()).
			Str("elmType", elmTypeName(elmType)).
			Msg("Type conversion")
		fields = append(fields, &ElmField{
			JSONName: jsonName,
			ElmName:  elmName,
			ElmType:  elmType,
		})
	}

	return &ElmRecord{name: recordName, Fields: fields}, nil
}
