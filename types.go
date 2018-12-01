package main

import "go/types"

// elmType translates a Go type into ana Elm type and JSON decoder pair.
func elmType(goType types.Type) (elmType, elmDecoder string) {
	switch t := goType.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return "Bool", "bool"
		case types.Float32, types.Float64:
			return "Float", "float"
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return "Int", "int"
		case types.String:
			return "String", "string"
		}
		return "", ""
	}
	return "", ""
}
