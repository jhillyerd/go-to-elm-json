package main

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

const examples = "testdata/examples.go"

func TestModuleFromStructErrors(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

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
		structType, err := structFromProg(prog, "main", tt.name)
		if err == nil {
			_, err = moduleFromStruct(structType, tt.name)
		}
		got := err != nil
		if got != tt.errorExpected {
			t.Errorf("moduleFromStruct(%q): got error %v, want error %v\nerror was: %v",
				tt.name, got, tt.errorExpected, err)
		}
	}
}

func TestModuleFromStructNameConversions(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "Strings"
	want := &Module{
		Name: name,
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
	structType, err := structFromProg(prog, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := moduleFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error("Module struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
}

func TestParseStructMultipleNames(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "MultiNames"
	want := &Module{
		Name: name,
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
	structType, err := structFromProg(prog, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := moduleFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error("Module struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
}

func TestModuleFromStructTypeConversions(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "OtherTypes"
	want := &Module{
		Name: name,
		Fields: []Field{
			{
				GoName:     "AnInteger",
				GoType:     "int",
				JSONName:   "AnInteger",
				ElmName:    "anInteger",
				ElmType:    "Int",
				ElmDecoder: "int",
			},
			{
				GoName:     "BigInteger",
				GoType:     "int64",
				JSONName:   "BigInteger",
				ElmName:    "bigInteger",
				ElmType:    "Int",
				ElmDecoder: "int",
			},
			{
				GoName:     "AFloat",
				GoType:     "float32",
				JSONName:   "AFloat",
				ElmName:    "aFloat",
				ElmType:    "Float",
				ElmDecoder: "float",
			},
			{
				GoName:     "BigFloat",
				GoType:     "float64",
				JSONName:   "BigFloat",
				ElmName:    "bigFloat",
				ElmType:    "Float",
				ElmDecoder: "float",
			},
			{
				GoName:     "NoNoNo",
				GoType:     "bool",
				JSONName:   "NoNoNo",
				ElmName:    "noNoNo",
				ElmType:    "Bool",
				ElmDecoder: "bool",
			},
		},
	}
	structType, err := structFromProg(prog, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := moduleFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Error("Module struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
}
