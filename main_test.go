package main

import (
	"bytes"
	"testing"

	"github.com/jhillyerd/goldiff"
)

func TestMainOutput(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name, goldenFile string
	}{
		{"Strings", "strings.golden"},
		{"OtherTypes", "othertypes.golden"},
		{"SliceTypes", "slicetypes.golden"},
		{"NestedStructs", "nestedstructs.golden"},
		{"OptionalValues", "optionalvalues.golden"},
		{"NullableValues", "nullablevalues.golden"},
	}

	buf := &bytes.Buffer{}
	for _, tt := range tests {
		buf.Reset()
		err = generateElm(buf, pkgs, "main", tt.name, make(TypeNamePairs))
		if err != nil {
			t.Error(err)
			continue
		}
		goldiff.File(t, buf.Bytes(), "testdata", "examples", tt.goldenFile)
	}
}
