package main

import (
	"strings"
	"sync"
	"testing"

	"github.com/go-test/deep"
	"golang.org/x/tools/go/packages"
)

// packageCache prevents each test from having to re-parse the Go example code.
type packageCache struct {
	sync.Mutex
	loadedPkgs map[string][]*packages.Package
}

func (p *packageCache) load(path string) (pkgs []*packages.Package, err error) {
	p.Lock()
	defer p.Unlock()

	if p.loadedPkgs == nil {
		p.loadedPkgs = make(map[string][]*packages.Package)
	}

	// Check cache first.
	if pkgs, ok := p.loadedPkgs[path]; ok {
		return pkgs, nil
	}

	pkgs, err = loadPackages([]string{path})
	if err != nil {
		return nil, err
	}

	p.loadedPkgs[path] = pkgs
	return pkgs, nil
}

var pkgCache = &packageCache{}

const examples = "testdata/examples.go"

func TestRecordFromStructErrors(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
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
		{"OptionalValues", false},
	}
	for _, tt := range tests {
		structType, err := getStructDef(pkgs, "main", tt.name)
		if err == nil {
			_, err = recordFromStruct(NewResolver(make(TypeNamePairs)), structType, tt.name)
		}
		got := err != nil
		if got != tt.errorExpected {
			t.Errorf("recordFromStruct(%q): got error %v, want error %v\nerror was: %v",
				tt.name, got, tt.errorExpected, err)
		}
	}
}

func TestRecordFromStructAbbrev(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	input := "JSONObject"
	want := "JsonObject"
	structType, err := getStructDef(pkgs, "main", input)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	record, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, input)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got := record.Name()
	if got != want {
		t.Errorf("Got record name %q, want %q", got, want)
	}
}

func TestRecordFromStructNameConversions(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
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
				Optional: true,
			},
			{
				JSONName: "AnotherOptionalString",
				ElmName:  "anotherOptionalString",
				ElmType:  elmString,
				Optional: true,
			},
		},
	}
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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
	pkgs, err := pkgCache.load(examples)
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
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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
	pkgs, err := pkgCache.load(examples)
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
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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
	pkgs, err := pkgCache.load(examples)
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
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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

func TestRecordFromStructOptionals(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	name := "OptionalValues"
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
			{
				JSONName: "opt-string",
				ElmName:  "optString",
				ElmType:  elmString,
				Optional: true,
			},
			{
				JSONName: "OptInt",
				ElmName:  "optInt",
				ElmType:  elmInt,
				Optional: true,
			},
			{
				JSONName: "OptBool",
				ElmName:  "optBool",
				ElmType:  elmBool,
				Optional: true,
			},
		},
	}
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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

func TestRecordFromStructNullables(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	name := "NullableValues"
	innerType := &ElmRecord{
		name: "InnerStruct",
		Fields: []*ElmField{
			{
				JSONName: "Value",
				ElmName:  "value",
				ElmType:  elmString,
			},
		},
	}
	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
			{
				JSONName: "NullString",
				ElmName:  "nullString",
				ElmType:  &ElmPointer{elem: elmString},
			},
			{
				JSONName: "OptNullString",
				ElmName:  "optNullString",
				ElmType:  &ElmPointer{elem: elmString},
				Optional: true,
			},
			{
				JSONName: "NullInt",
				ElmName:  "nullInt",
				ElmType:  &ElmPointer{elem: elmInt},
			},
			{
				JSONName: "NullStruct",
				ElmName:  "nullStruct",
				ElmType:  &ElmPointer{elem: innerType},
			},
		},
	}
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
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

func TestRecordFromStructNested(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	name := "NestedStructs"
	innerType := &ElmRecord{
		name: "InnerStruct",
		Fields: []*ElmField{
			{
				JSONName: "Value",
				ElmName:  "value",
				ElmType:  elmString,
			},
		},
	}

	want := &ElmRecord{
		name: name,
		Fields: []*ElmField{
			{
				JSONName: "OuterName",
				ElmName:  "outerName",
				ElmType:  elmString,
			},
			{
				JSONName: "InnerValue1",
				ElmName:  "innerValue1",
				ElmType:  innerType,
			},
			{
				JSONName: "InnerValue2",
				ElmName:  "innerValue2",
				ElmType:  innerType,
			},
		},
	}
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(make(TypeNamePairs)), structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(got, want); diff != nil {
		t.Fatal("ElmRecord struct did not match expectations:\n" + strings.Join(diff, "\n"))
	}
	if innerRecord, ok := got.Fields[1].ElmType.(*ElmRecord); ok {
		t.Logf("innerRecord: %#v", innerRecord)
		for i, f := range innerRecord.Fields {
			t.Logf("innerRecord[%v]: %#v", i, f)
			t.Logf("innerRecord[%v].ElmType: %s", i, f.ElmType.Name())
		}
	} else {
		t.Errorf("Fields[1].ElmType was %T, not *ElmRecord", got.Fields[1].ElmType)
	}

	if !got.Equal(want) {
		t.Error("ElmRecord struct did not match expectations, likely in an ElmType field.")
	}
}

func TestRecordFromStructNestedRenames(t *testing.T) {
	pkgs, err := pkgCache.load(examples)
	if err != nil {
		t.Fatal(err)
	}

	renames := make(TypeNamePairs)
	renames.Add("NestedStructs:NewOuter")
	renames.Add("innerStruct:NewInner")

	name := "NestedStructs"
	structType, err := getStructDef(pkgs, "main", name)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	got, err := recordFromStruct(NewResolver(renames), structType, name)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	wantName := "NewOuter"
	gotName := got.name
	if gotName != wantName {
		t.Errorf("got name %q, want %q", gotName, wantName)
	}

	wantName = "NewInner"
	gotType, ok := got.Fields[1].ElmType.(*ElmRecord)
	if !ok {
		t.Fatalf("want type *ElmRecord, got %T", got.Fields[1].ElmType)
	}
	gotName = gotType.name
	if gotName != wantName {
		t.Errorf("got name %q, want %q", gotName, wantName)
	}
}
