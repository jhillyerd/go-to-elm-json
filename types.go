package main

import "go/types"

// ElmType represents a type in Elm.
type ElmType interface {
	Name() string
	Codec(prefix string) string
}

// BasicElmType represents primitive types in Elm.
type BasicElmType struct {
	name  string
	codec string
}

// Name returns the name of the Elm type.
func (t *BasicElmType) Name() string {
	return t.name
}

// Codec returns the name of the Elm JSON encoder/decoder for this type.
func (t *BasicElmType) Codec(prefix string) string {
	return prefix + "." + t.codec
}

var (
	elmBool   = &BasicElmType{name: "Bool", codec: "bool"}
	elmFloat  = &BasicElmType{name: "Float", codec: "float"}
	elmInt    = &BasicElmType{name: "Int", codec: "int"}
	elmString = &BasicElmType{name: "String", codec: "string"}
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
	}
	return nil
}
