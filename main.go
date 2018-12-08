package main

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

func main() {
	// Load output template.
	tmpl, err := template.New("elm").Parse(elmTemplate)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// Parse Go.
	prog, rest, err := progFromArgs(os.Args[1:])
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if len(rest) != 2 {
		fmt.Fprintf(os.Stderr, "Want package and a single type to find, got: %v\n", rest)
		fmt.Fprintln(os.Stderr, loader.FromArgsUsage)
		os.Exit(1)
	}
	packageName := rest[0]
	objectName := rest[1]

	// Process definition.
	structType, err := structFromProg(prog, packageName, objectName)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	record, err := recordFromStruct(structType, objectName)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = tmpl.Execute(os.Stdout, record)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

// progFromArgs takes an x/tools/go/loader argument string and parses the specified Go files,
// leaving remaining arguments in rest.
func progFromArgs(args []string) (prog *loader.Program, rest []string, err error) {
	// Configure type checker.
	var conf loader.Config
	rest, err = conf.FromArgs(args, false)
	if err != nil {
		return nil, nil, err
	}

	// Load Go package.
	prog, err = conf.Load()
	if err != nil {
		return nil, nil, err
	}
	return prog, rest, nil
}

// structFromProg finds the requested object and confirms it's a struct type definition.
func structFromProg(prog *loader.Program, packageName, objectName string) (*types.Struct, error) {
	// Lookup package.
	pkg := prog.Package(packageName)
	if pkg == nil {
		return nil, errors.Errorf("Package %s not found", packageName)
	}

	// Lookup struct object.
	obj := pkg.Pkg.Scope().Lookup(objectName)
	if obj == nil {
		return nil, errors.Errorf("Definition %s.%s not found", packageName, objectName)
	}
	objType := obj.Type().Underlying()
	structType, ok := objType.(*types.Struct)
	if !ok {
		return nil, errors.Errorf("%s type is %T, want *types.Struct", obj.Id(), objType)
	}
	return structType, nil
}
