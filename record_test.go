package main

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
)

const examples = "testdata/examples.go"

func TestRecordFromStructErrors(t *testing.T) {
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
			_, err = recordFromStruct(structType, tt.name)
		}
		got := err != nil
		if got != tt.errorExpected {
			t.Errorf("recordFromStruct(%q): got error %v, want error %v\nerror was: %v",
				tt.name, got, tt.errorExpected, err)
		}
	}
}

func TestRecordFromStructAbbrev(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	input := "JSONObject"
	want := "JsonObject"
	structType, err := structFromProg(prog, "main", input)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	record, err := recordFromStruct(structType, input)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got := record.Name()
	if got != want {
		t.Errorf("Got record name %q, want %q", got, want)
	}
}

func TestRecordFromStructNameConversions(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "Strings"
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
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
	got, err := recordFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Fatal("ElmRecord struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
	if !got.Equal(want) {
		t.Error("ElmRecord struct did not match expectations, likely in an ElmType field.")
	}
}

func TestParseStructMultipleNames(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "MultiNames"
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
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
	got, err := recordFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Fatal("ElmRecord struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
	if !got.Equal(want) {
		t.Error("ElmRecord struct did not match expectations, likely in an ElmType field.")
	}
}

func TestRecordFromStructTypeConversions(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "OtherTypes"
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
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
	got, err := recordFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Fatal("ElmRecord struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
	if !got.Equal(want) {
		t.Error("ElmRecord struct did not match expectations, likely in an ElmType field.")
	}
}

func TestRecordFromStructSlices(t *testing.T) {
	prog, _, err := progFromArgs([]string{examples})
	if err != nil {
		t.Fatal(err)
	}

	name := "SliceTypes"
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
			{
				JSONName: "Bools",
				ElmName:  "bools",
				ElmType:  &ElmList{elem: elmBool},
			},
			{
				JSONName: "Floats",
				ElmName:  "floats",
				ElmType:  &ElmList{elem: elmFloat},
			},
			{
				JSONName: "Strings",
				ElmName:  "strings",
				ElmType:  &ElmList{elem: elmString},
			},
		},
	}
	structType, err := structFromProg(prog, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Fatal("ElmRecord struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
	if !got.Equal(want) {
		t.Error("ElmRecord struct did not match expectations, likely in an ElmType field.")
	}
}
