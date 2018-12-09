package main

import "strings"

func precedence(s string) string {
	if strings.ContainsRune(s, ' ') {
		return "(" + s + ")"
	}
	return s
}

func splitTypeNamePair(s string) (string, string) {
	els := strings.Split(s, ":")
	goName := els[0]
	elmName := goName
	if len(els) > 1 {
		elmName = els[1]
	}
	return goName, elmName
}

// TypeNamePairs maps the source Go type name to the target Elm type name.
type TypeNamePairs map[string]string

// Add splits the input string on : and updates the map.
func (m TypeNamePairs) Add(s string) {
	goName, elmName := splitTypeNamePair(s)
	m[goName] = elmName
}

// ElmName returns the Elm record name for a Go struct type.
func (m TypeNamePairs) ElmName(typeName string) string {
	recordName := m[typeName]
	if recordName == "" {
		camelCaseName := camelCase(typeName)
		recordName = strings.ToUpper(camelCaseName[:1]) + camelCaseName[1:]
	}
	return recordName
}
