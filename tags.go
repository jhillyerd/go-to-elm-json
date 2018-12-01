package main

import (
	"strings"
)

const (
	tagPrefix = "json:\""
	tagSuffix = "\""
)

// parseTag splits a struct field's json tag into its name and comma-separated options.
func parseTag(tag string) (string, string) {
	if !strings.HasPrefix(tag, tagPrefix) {
		return "", ""
	}
	if !strings.HasSuffix(tag, tagSuffix) {
		return "", ""
	}
	tag = tag[len(tagPrefix) : len(tag)-len(tagSuffix)]
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return tag, ""
}

// hasOption tests for presence of option name in option string s.
func hasOption(name string, s string) bool {
	if len(s) == 0 {
		return false
	}
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == name {
			return true
		}
		s = next
	}
	return false
}
