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
	"golang.org/x/tools/go/loader"
)

// logger used by unit tests.
var logger = zerolog.New(
	zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}).With().Timestamp().Logger()

// TemplateData holds the context for the template.
type TemplateData struct {
	Record *ElmRecord
	Nested []*ElmRecord
}

func main() {
	// Flags.
	verbose := flag.Bool("v", false, "verbose (debug) output")
	color := flag.Bool("color", runtime.GOOS != "windows", "colorize debug output")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage of %s [opts] <args> -- <pkg name> <type name>:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(),
			"  Ex: 'go-to-elm-json *.go -- main MyStruct > MyStruct.elm'\n\n")
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), loader.FromArgsUsage)
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

	// Parse Go.
	prog, rest, err := progFromArgs(flag.Args())
	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't parse Go")
	}
	if len(rest) != 2 {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Wanted a package and a struct type to convert, got: %v\n\n", rest)
		flag.Usage()
		os.Exit(1)
	}
	packageName := rest[0]
	objectName := rest[1]

	// Output Elm.
	err = generateElm(os.Stdout, prog, packageName, objectName)
	if err != nil {
		logger.Fatal().Err(err).Msg("Generation failed")
	}
}

// generateElm processes the provided program and outputs Elm code to the provider writer.
func generateElm(w io.Writer, prog *loader.Program, packageName string, objectName string) error {
	// Load output template.
	tmpl, err := template.New("elm").Parse(elmTemplate)
	if err != nil {
		return errors.Wrap(err, "Couldn't parse template")
	}

	// Process definition.
	structType, err := structFromProg(prog, packageName, objectName)
	if err != nil {
		return errors.Wrap(err, "Couldn't find struct")
	}
	resolver := NewResolver()
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
