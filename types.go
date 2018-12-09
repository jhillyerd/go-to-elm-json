package main

import (
	"go/types"

	"github.com/pkg/errors"
)

var (
	elmBool   = &ElmBasicType{name: "Bool", codec: "bool"}
	elmFloat  = &ElmBasicType{name: "Float", codec: "float"}
	elmInt    = &ElmBasicType{name: "Int", codec: "int"}
	elmString = &ElmBasicType{name: "String", codec: "string"}
)

// ElmType represents a type in Elm.
type ElmType interface {
	Name() string
	Decoder(prefix string) string
	Encoder(prefix string) string
	Equal(other ElmType) bool
	Nullable() bool
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

// Decoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmBasicType) Decoder(prefix string) string {
	return prefix + "." + t.codec
}

// Encoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmBasicType) Encoder(prefix string) string {
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

// Nullable indicates whether this type can be nil.
func (t *ElmBasicType) Nullable() bool {
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

// Decoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmList) Decoder(prefix string) string {
	return "(" + prefix + ".list " + t.elem.Decoder(prefix) + ")"
}

// Encoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmList) Encoder(prefix string) string {
	return "(" + prefix + ".list " + t.elem.Encoder(prefix) + ")"
}

// Equal tests for equality with another ElmType.
func (t *ElmList) Equal(other ElmType) bool {
	if o, ok := other.(*ElmList); ok {
		return t.elem.Equal(o.elem)
	}
	return false
}

// Nullable indicates whether this type can be nil.
func (t *ElmList) Nullable() bool {
	return true
}

// ElmPointer represents a pointer to an instance of another type.
type ElmPointer struct {
	elem ElmType
}

// Name returns the name of the Elm type.
func (t *ElmPointer) Name() string {
	return t.elem.Name()
}

// Decoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmPointer) Decoder(prefix string) string {
	return t.elem.Decoder(prefix)
}

// Encoder returns the name of the Elm JSON encoder/decoder for this type.
func (t *ElmPointer) Encoder(prefix string) string {
	return t.elem.Encoder(prefix)
}

// Equal tests for equality with another ElmType.
func (t *ElmPointer) Equal(other ElmType) bool {
	if o, ok := other.(*ElmPointer); ok {
		return t.elem.Equal(o.elem)
	}
	return false
}

// Nullable indicates whether this type can be nil.
func (t *ElmPointer) Nullable() bool {
	return true
}

// ElmTypeResolver maintains a cache of Go to Elm type conversions.
type ElmTypeResolver struct {
	resolved map[string]*ElmRecord
	ordered  []*ElmRecord
}

// NewResolver creates an empty resolver.
func NewResolver() *ElmTypeResolver {
	return &ElmTypeResolver{
		resolved: make(map[string]*ElmRecord),
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
	case *types.Pointer:
		elemType, err := r.Convert(t.Elem())
		if err != nil {
			return nil, err
		}
		return &ElmPointer{elem: elemType}, nil
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
			return r.resolveRecord(goName, u)
		}
	}
	return nil, errors.Errorf("don't know how to handle Go type %s (%T)", goType, goType)
}

// CachedRecords returns slice of resolved Elm records.
func (r *ElmTypeResolver) CachedRecords() []*ElmRecord {
	return r.ordered
}

// resolveRecord converts the struct to an Elm record, or returns the cached version.
func (r *ElmTypeResolver) resolveRecord(goName string, stype *types.Struct) (*ElmRecord, error) {
	if record := r.resolved[goName]; record != nil {
		return record, nil
	}
	record, err := recordFromStruct(r, stype, goName)
	if err != nil {
		return nil, err
	}
	logger.Debug().
		Str("name", goName).
		Str("type", elmTypeName(record)).
		Msg("Caching resolved type")
	r.resolved[goName] = record
	r.ordered = append(r.ordered, record)
	return record, nil
}
