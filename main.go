package main

import (
	"flag"
	"fmt"
	"go/types"
	"io"
	"os"
	"runtime"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/tools/go/packages"
)

// logger used by unit tests.
var logger = zerolog.New(
	zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}).With().Timestamp().Logger()

// TemplateData holds the context for the template.
type TemplateData struct {
	Record *ElmRecord
	Nested []*ElmRecord
}

const help = `
Usage Example:
  go-to-elm-json *.go -- main MyThingJSON:MyThing > MyThing.elm

<go files> syntax:
  This list is passed to packages.Load() unmodified.  It can be a literal list
  of files, or a list of packages.  Examples:
    path/main.go                        - single file
    path/*.go                           - shell glob
    ./path                              - local package path
    github.com/jhillyerd/go-to-elm-json - full package path

  See https://pkg.go.dev/golang.org/x/tools/go/packages for more.

`

func main() {
	// Flags.
	verbose := flag.Bool("v", false, "verbose (debug) output")
	color := flag.Bool("color", runtime.GOOS != "windows", "colorize debug output")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [opts] <go files> -- <pkg name> \\\n"+
			"  <root go type:elm name> [<go type:elm name> ...]:\n\n", os.Args[0])
		fmt.Fprint(flag.CommandLine.Output(), help)
		flag.PrintDefaults()
	}
	flag.Parse()

	// Logging.
	if *verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, NoColor: !*color}).
		With().Timestamp().Logger()

	// Split files and args at `--`.
	var args []string
	files := flag.Args()
	for i, arg := range files {
		if arg == "--" {
			args = files[i+1:]
			files = files[:i]
			break
		}
	}
	if len(args) < 2 {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Wanted a package and a struct type to convert, got: %v\n\n", args)
		flag.Usage()
		os.Exit(1)
	}

	// Parse Go.
	pkgs, err := loadPackages(files)
	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't load Go package")
	}
	packageName := args[0]
	objectName, _ := splitTypeNamePair(args[1])
	renames := make(TypeNamePairs)
	for _, arg := range args[1:] {
		renames.Add(arg)
	}

	// Output Elm.
	err = generateElm(os.Stdout, pkgs, packageName, objectName, renames)
	if err != nil {
		logger.Fatal().Err(err).Msg("Generation failed")
	}
}

// generateElm processes the provided program and outputs Elm code to the provider writer.
func generateElm(
	w io.Writer,
	pkgs []*packages.Package,
	packageName string,
	objectName string,
	renames TypeNamePairs) error {
	// Load output template.
	tmpl, err := template.New("elm").Parse(elmTemplate)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse template")
	}

	// Process definition.
	structType, err := structFromPackage(pkgs, packageName, objectName)
	if err != nil {
		return errors.Wrap(err, "Couldn't find struct")
	}
	resolver := NewResolver(renames)
	record, err := recordFromStruct(resolver, structType, objectName)
	if err != nil {
		return errors.Wrap(err, "Couldn't convert struct")
	}

	// Render Elm.
	data := &TemplateData{
		Record: record,
		Nested: resolver.CachedRecords(),
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return errors.Wrap(err, "Couldn't render template")
	}

	return nil
}

// loadPackages takes an x/tools/go/packages argument list and parses the specified Go files.
func loadPackages(args []string) (pkgs []*packages.Package, err error) {
	// Configure package loader, load packages.
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes,
	}
	pkgs, err = packages.Load(cfg, args...)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		logger.Warn().Msg("There were non-fatal errors loading package")
	}

	// Dump package type info for troubleshooting loading problems.
	logger.Debug().Func(func(e *zerolog.Event) {
		for _, p := range pkgs {
			e.Str("ID", p.ID).Str("Name", p.Name).Str("PkgPath", p.PkgPath).Msg("Go package loaded")
			for _, sn := range p.Types.Scope().Names() {
				logger.Debug().
					Str("Name", sn).
					Str("PkgName", p.Name).
					Msg("Type found")
			}
		}
	})

	return pkgs, nil
}

// structFromPackage finds the requested object and confirms it's a struct type definition.
func structFromPackage(pkgs []*packages.Package, packageName, typeName string) (*types.Struct, error) {
	// Lookup package.
	var pkg *packages.Package
	for _, p := range pkgs {
		if p.Name == packageName {
			pkg = p
			break
		}
	}
	if pkg == nil {
		return nil, errors.Errorf("Package %s not found", packageName)
	}

	// Lookup type definition.
	obj := pkg.Types.Scope().Lookup(typeName)
	if obj == nil {
		return nil, errors.Errorf("Definition %s.%s not found", packageName, typeName)
	}
	objType := obj.Type().Underlying()
	structType, ok := objType.(*types.Struct)
	if !ok {
		return nil, errors.Errorf("%s type is %T, want *types.Struct", obj.Id(), objType)
	}

	return structType, nil
}
