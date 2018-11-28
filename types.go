package main

// elmType translates a Go type into ana Elm type and JSON decoder pair.
func elmType(goType string) (elmType, elmDecoder string) {
	switch goType {
	case "string":
		return "String", "string"
	}
	return "", ""
}
