package main

import (
	"bytes"
	"testing"

	"github.com/jhillyerd/goldiff"
)

func TestMainOutput(t *testing.T) {
	prog, err := programs.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name, goldenFile string
	}{
		{"Strings", "strings.golden"},
		{"OtherTypes", "othertypes.golden"},
		{"SliceTypes", "slicetypes.golden"},
	}

	buf := &bytes.Buffer{}
	for _, tt := range tests {
		buf.Reset()
		err = generateElm(buf, prog, "main", tt.name)
		goldiff.File(t, buf.Bytes(), "testdata", "examples", tt.goldenFile)
	}
}
