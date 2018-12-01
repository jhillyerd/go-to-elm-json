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
				JSONName: "ExportedBareString",
				ElmName:  "exportedBareString",
				ElmType:  elmString,
			},
			{
				JSONName: "exported-tagged-string",
				ElmName:  "exportedTaggedString",
				ElmType:  elmString,
			},
			{
				JSONName: "exported-optional-string",
				ElmName:  "exportedOptionalString",
				ElmType:  elmString,
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
				JSONName: "One",
				ElmName:  "one",
				ElmType:  elmString,
			},
			{
				JSONName: "Two",
				ElmName:  "two",
				ElmType:  elmString,
			},
			{
				JSONName: "Three",
				ElmName:  "three",
				ElmType:  elmString,
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
				JSONName: "AnInteger",
				ElmName:  "anInteger",
				ElmType:  elmInt,
			},
			{
				JSONName: "BigInteger",
				ElmName:  "bigInteger",
				ElmType:  elmInt,
			},
			{
				JSONName: "AFloat",
				ElmName:  "aFloat",
				ElmType:  elmFloat,
			},
			{
				JSONName: "BigFloat",
				ElmName:  "bigFloat",
				ElmType:  elmFloat,
			},
			{
				JSONName: "NoNoNo",
				ElmName:  "noNoNo",
				ElmType:  elmBool,
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
