package main

import "go/types"

// ElmType represents a type in Elm.
type ElmType interface {
	Name() string
	Codec(prefix string) string
	Equal(other ElmType) bool
}

// ElmBasicType represents primitive types in Elm.
type ElmBasicType struct {
	name  string
	codec string
}

// Name returns the name of the Elm type.
func (t *ElmBasicType) Name() string {
	return t.name
}

// Codec returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmBasicType) Codec(prefix string) string {
	return prefix + "." + t.codec
}

// Equal tests for equality with another ElmType.
func (t *ElmBasicType) Equal(other ElmType) bool {
	if o, ok := other.(*ElmBasicType); ok {
		return t.name == o.name &&
			t.codec == o.codec
	}
	return false
}

// ElmListType represents a list of another type.
type ElmListType struct {
	elem ElmType
}

// Name returns the name of the Elm type.
func (t *ElmListType) Name() string {
	return "List " + t.elem.Name()
}

// Codec returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmListType) Codec(prefix string) string {
	return "(" + prefix + ".list " + t.elem.Codec(prefix) + ")"
}

// Equal tests for equality with another ElmType.
func (t *ElmListType) Equal(other ElmType) bool {
	if o, ok := other.(*ElmListType); ok {
		return t.elem.Equal(o.elem)
	}
	return false
}

var (
	elmBool   = &ElmBasicType{name: "Bool", codec: "bool"}
	elmFloat  = &ElmBasicType{name: "Float", codec: "float"}
	elmInt    = &ElmBasicType{name: "Int", codec: "int"}
	elmString = &ElmBasicType{name: "String", codec: "string"}
)

// elmType translates a Go type into ana Elm type and JSON decoder pair.
func elmType(goType types.Type) ElmType {
	switch t := goType.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return elmBool
		case types.Float32, types.Float64:
			return elmFloat
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return elmInt
		case types.String:
			return elmString
		}
	case *types.Slice:
		return &ElmListType{elem: elmType(t.Elem())}
	}
	return nil
}
