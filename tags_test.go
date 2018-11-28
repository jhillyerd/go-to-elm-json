package main

import "testing"

func TestParseTag(t *testing.T) {
	testCases := []struct {
		input, name, options string
	}{
		{"bad-prefix\"`", "", ""},
		{"`json:\"bad-suffix", "", ""},
		{"`json:\"name\"`", "name", ""},
		{"`json:\"name,option\"`", "name", "option"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			gname, goptions := parseTag(tc.input)
			if gname != tc.name {
				t.Errorf("name got %q, want %q", gname, tc.name)
			}
			if goptions != tc.options {
				t.Errorf("options got %q, want %q", goptions, tc.options)
			}
		})
	}
}

func TestHasOption(t *testing.T) {
	testCases := []struct {
		input  string
		option string
		want   bool
	}{
		{"foo", "bar", false},
		{"bar", "bar", true},
		{"foo,bar", "foo", true},
		{"foo,bar", "bar", true},
		{"foo,bar", "baz", false},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := hasOption(tc.option, tc.input)
			if got != tc.want {
				t.Errorf("option %q in %q got %v, want %v", tc.option, tc.input, got, tc.want)
			}
		})
	}
}
