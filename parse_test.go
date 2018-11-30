package main

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

const examples = "testdata/examples.go"

func TestParseStructTypeErrors(t *testing.T) {
	var tests = []struct {
		name          string
		errorExpected bool
	}{
		{"DoesNotExist", true},
		{"AnInterface", true},
		{"Empty", true},
		{"Strings", false},
	}
	for _, tt := range tests {
		_, err := parseStructType(examples, tt.name)
		got := err != nil
		if got != tt.errorExpected {
			t.Errorf("parseStructType(%q): got error %v, want error %v\nerror was: %v",
				tt.name, got, tt.errorExpected, err)
		}
	}
}

func TestParseStructTypeNameConversions(t *testing.T) {
	want := &Module{
		Name: "Strings",
		Fields: []Field{
			{
				GoName:     "ExportedBareString",
				GoType:     "string",
				JSONName:   "ExportedBareString",
				ElmName:    "exportedBareString",
				ElmType:    "String",
				ElmDecoder: "string",
			},
			{
				GoName:     "ExportedTaggedString",
				GoType:     "string",
				JSONName:   "exported-tagged-string",
				ElmName:    "exportedTaggedString",
				ElmType:    "String",
				ElmDecoder: "string",
			},
			{
				GoName:     "ExportedOptionalString",
				GoType:     "string",
				JSONName:   "exported-optional-string",
				ElmName:    "exportedOptionalString",
				ElmType:    "String",
				ElmDecoder: "string",
			},
		},
	}
	got, err := parseStructType(examples, "Strings")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error("Module struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
}

func TestParseStructMultipleNames(t *testing.T) {
	want := &Module{
		Name: "MultiNames",
		Fields: []Field{
			{
				GoName:     "One",
				GoType:     "string",
				JSONName:   "One",
				ElmName:    "one",
				ElmType:    "String",
				ElmDecoder: "string",
			},
			{
				GoName:     "Two",
				GoType:     "string",
				JSONName:   "Two",
				ElmName:    "two",
				ElmType:    "String",
				ElmDecoder: "string",
			},
			{
				GoName:     "Three",
				GoType:     "string",
				JSONName:   "Three",
				ElmName:    "three",
				ElmType:    "String",
				ElmDecoder: "string",
			},
		},
	}
	got, err := parseStructType(examples, "MultiNames")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error("Module struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
}
