package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

// Module represents an Elm module.
type Module struct {
	Name   string
	Fields []Field
}

// Field represents a Go -> JSON -> Elm field.
type Field struct {
	GoName     string
	GoType     string
	JSONName   string
	ElmName    string
	ElmType    string
	ElmDecoder string
}

func parseStructType(filePath string, typeName string) (*Module, error) {
	// Parse source file
	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "parse of %q failed", filePath)
	}

	// Remove nodes that are not named typeName.
	found := ast.FilterFile(fileNode, func(s string) bool {
		return s == typeName
	})
	if !found {
		return nil, errors.Errorf("no declaration found for name %q", typeName)
	}

	// Locate the struct definition.
	var astFields *ast.FieldList
	ast.Inspect(fileNode, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.StructType:
			astFields = x.Fields
			return false
		}
		return true
	})

	if astFields == nil {
		return nil, errors.Errorf("no struct found for type %q", typeName)
	}
	count := astFields.NumFields()
	if count == 0 {
		return nil, errors.Errorf("struct %q had no fields", typeName)
	}

	// Convert to our field type.
	var fields []Field
	for _, f := range astFields.List {
		for _, fnam := range f.Names {
			if !fnam.IsExported() {
				continue
			}
			goName := fnam.Name

			typeIdent, ok := f.Type.(*ast.Ident)
			if !ok {
				return nil, errors.Errorf("expected %q to have an Ident in Type, got %T", goName, f.Type)
			}

			jsonName := goName
			if f.Tag != nil && len(f.Tag.Value) > 2 {
				tagName, _ := parseTag(f.Tag.Value)
				jsonName = tagName
			}

			elmName := strings.ToLower(goName[0:1]) + goName[1:]
			elmType, elmDecoder := elmType(typeIdent.Name)
			fields = append(fields, Field{
				GoName:     goName,
				GoType:     typeIdent.Name,
				JSONName:   jsonName,
				ElmName:    elmName,
				ElmType:    elmType,
				ElmDecoder: elmDecoder,
			})
		}
	}

	return &Module{Name: typeName, Fields: fields}, nil
}
