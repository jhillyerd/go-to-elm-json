package main

import "fmt"

// AnInterface is a boring interface.
type AnInterface interface {
	boring()
}

// Empty is a struct with no fields.
type Empty struct{}

// JSONObject contains an abbreviation in its name.
type JSONObject struct {
	Value string
}

// Strings is a struct of strings.
type Strings struct {
	ExportedBareString     string
	ExportedTaggedString   string `json:"exported-tagged-string"`
	ExportedOptionalString string `json:"exported-optional-string,omitempty"`
	internalBareString     string
}

// OtherTypes is a struct with types other than string.
type OtherTypes struct {
	AnInteger  int
	BigInteger int64
	AFloat     float32
	BigFloat   float64
	NoNoNo     bool
}

// MultiNames defines multiple fields per line.
type MultiNames struct {
	One, Two, Three string
}

// SliceTypes defines some list-like fields.
type SliceTypes struct {
	Bools   []bool
	Floats  []float32
	Strings []string
}

func main() {
	fmt.Printf("Hello world!\n")
}
