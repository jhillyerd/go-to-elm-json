package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		log.Fatal("Expecting two args, go source and struct type name")
	}
	tmpl, err := template.New("elm").Parse(elmTemplate)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	module, err := parseStructType(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalf("%+v", err)
	}
	err = tmpl.Execute(os.Stdout, module)
}
