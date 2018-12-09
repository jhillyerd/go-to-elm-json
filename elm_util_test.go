package main

import "testing"

func TestPrecedence(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", ""},
		{"String", "String"},
		{"List String", "(List String)"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := precedence(tc.input)
			if got != tc.want {
				t.Errorf("%q got %q, want %q", tc.input, got, tc.want)
			}
		})
	}

}
