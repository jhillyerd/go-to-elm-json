package main

import "strings"

// abbreviations are converted from ABC in Go to abc or Abc in Elm.  Sort by length (long first),
// then alphabetically for ease of updating.
var abbreviations = []string{
	"HTML",
	"HTTP",
	"IMAP",
	"JSON",
	"POP3",
	"REST",
	"SMTP",
	"URI",
	"URL",
	"WWW",
	"XML",
	"GB",
	"ID",
	"KB",
	"MB",
	"MS",
	"TB",
}

var replacements = abbrevReplacements(abbreviations)

func camelCase(s string) string {
	if s == "" {
		return ""
	}
	r := strings.NewReplacer(replacements...)
	s = r.Replace(s)
	return strings.ToLower(s[:1]) + s[1:]
}

// abbrevReplacements initial-uppercases each string in a, then interleaves the replacements for use
// with strings.Replacer.
func abbrevReplacements(a []string) []string {
	r := make([]string, 0, len(a)*2)
	for _, s := range a {
		r = append(r, s)
		r = append(r, s[:1]+strings.ToLower(s[1:]))
	}
	return r
}
