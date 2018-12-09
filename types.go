package main

import (
	"go/types"

	"github.com/pkg/errors"
)

// ElmType represents a type in Elm.
type ElmType interface {
	Name() string
	Codec(prefix string) string
	Equal(other ElmType) bool
}

func elmTypeName(t ElmType) string {
	if t == nil {
		return "<undefined>"
	}
	return t.Name()
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

// ElmList represents a list of another type.
type ElmList struct {
	elem ElmType
}

// Name returns the name of the Elm type.
func (t *ElmList) Name() string {
	return "List " + t.elem.Name()
}

// Codec returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmList) Codec(prefix string) string {
	return "(" + prefix + ".list " + t.elem.Codec(prefix) + ")"
}

// Equal tests for equality with another ElmType.
func (t *ElmList) Equal(other ElmType) bool {
	if o, ok := other.(*ElmList); ok {
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

// ElmTypeResolver maintains a cache of Go to Elm type conversions.
type ElmTypeResolver struct {
	resolved map[string]ElmType
}

// NewResolver creates an empty resolver.
func NewResolver() *ElmTypeResolver {
	return &ElmTypeResolver{
		make(map[string]ElmType),
	}
}

// Convert translates a Go type into an Elm type and JSON decoder pair.
func (r *ElmTypeResolver) Convert(goType types.Type) (ElmType, error) {
	switch t := goType.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return elmBool, nil
		case types.Float32, types.Float64:
			return elmFloat, nil
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return elmInt, nil
		case types.String:
			return elmString, nil
		}
	case *types.Slice:
		elemType, err := r.Convert(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ElmList{elem: elemType}, nil
	case *types.Named:
		goName := t.Obj().Name()
		switch u := t.Underlying().(type) {
		case *types.Struct:
			if record := r.resolved[goName]; record != nil {
				return record, nil
			}
			record, err := recordFromStruct(r, u, goName)
			if err != nil {
				return nil, err
			}
			logger.Debug().
				Str("name", goName).
				Str("type", elmTypeName(record)).
				Msg("Caching resolved type")
			r.resolved[goName] = record
			return record, nil
		}
	}
	return nil, errors.Errorf("don't know how to handle Go type %s (%T)", goType, goType)
}
