package main

import "testing"

func TestCamelCase(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", ""},
		{"abc", "abc"},
		{"MyWWW", "MyWww"},
		{"HTMLBody", "HtmlBody"},
		{"UserID", "UserId"},
		{"totalMBUploaded", "totalMbUploaded"},
	}
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			got := camelCase(tc.input)
			if got != tc.want {
				t.Errorf("camelCase(%q) got %q, want %q", tc.input, got, tc.want)
			}
		})
	}

}
