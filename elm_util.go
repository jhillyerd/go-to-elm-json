package main

import "strings"

func precedence(s string) string {
	if strings.ContainsRune(s, ' ') {
		return "(" + s + ")"
	}
	return s
}
