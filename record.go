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

// CamelCasedName leads with lowercase.
func (r *ElmRecord) CamelCasedName() string {
	return camelCase(r.name)
}

// Decoder for this record type.
func (r *ElmRecord) Decoder(prefix string) string {
	return r.CamelCasedName() + "Decoder"
}

// Encoder for this record type.
func (r *ElmRecord) Encoder(prefix string) string {
	return "encode" + r.name
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

// Nullable indicates whether this type can be nil.
func (r *ElmRecord) Nullable() bool {
	return false
}

// ElmField represents an Elm record field.
type ElmField struct {
	JSONName string
	ElmName  string
	ElmType  ElmType
	Optional bool
}

// Decoder returns the Elm JSON decoder for this field.
func (f *ElmField) Decoder(prefix string) string {
	if f.Optional || f.ElmType.Nullable() {
		return "(" + prefix + ".nullable " + f.ElmType.Decoder(prefix) + ")"
	}
	return f.ElmType.Decoder(prefix)
}

// Default returns a space-prefixed default value, or empty string.
func (f *ElmField) Default() string {
	if f.Optional {
		return " Nothing"
	}
	return ""
}

// Encoder reutrns the Elm JSON encoder for this field.
func (f *ElmField) Encoder(prefix string) string {
	if f.Optional || f.ElmType.Nullable() {
		return "maybe " + f.ElmType.Encoder(prefix)
	}
	return f.ElmType.Encoder(prefix)
}

// Pipeline returns the elm-decode-pipline function for this field.
func (f *ElmField) Pipeline(prefix string) string {
	if f.Optional {
		return prefix + ".optional"
	}
	return prefix + ".required"
}

// TypeDecl returns the type in Elm source format.
func (f *ElmField) TypeDecl() string {
	if f.Optional || f.ElmType.Nullable() {
		return "Maybe " + precedence(f.ElmType.Name())
	}
	return f.ElmType.Name()
}

// Equal test for equality with another field.
func (f *ElmField) Equal(o *ElmField) bool {
	return f.JSONName == o.JSONName &&
		f.ElmName == o.ElmName &&
		f.ElmType != nil &&
		f.ElmType.Equal(o.ElmType) &&
		f.Optional == o.Optional
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
		optional := false
		if len(stag) > 2 {
			tagName, tagOpts := parseTag(stag)
			if tagName != "" {
				jsonName = tagName
			}
			if hasOption("omitempty", tagOpts) {
				optional = true
			}
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
			Optional: optional,
		})
	}

	return &ElmRecord{name: recordName, Fields: fields}, nil
}
