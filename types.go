package main

// elmType translates a Go type into ana Elm type and JSON decoder pair.
func elmType(goType string) (elmType, elmDecoder string) {
	switch goType {
	case "bool":
		return "Bool", "bool"
	case "float32":
		return "Float", "float"
	case "float64":
		return "Float", "float"
	case "int":
		return "Int", "int"
	case "int64":
		return "Int", "int"
	case "string":
		return "String", "string"
	}
	return "", ""
}
