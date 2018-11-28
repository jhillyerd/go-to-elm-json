package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
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
