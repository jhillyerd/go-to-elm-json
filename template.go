package main

var elmTemplate = `module {{.Name}} exposing ({{.Name}})

import Json.Decode as D
import Json.Decode.Pipeline as P

type alias {{.Name}} =
{{- range $index, $el := .Fields }}
	{{ if $index }},{{ else }}{{"{"}}{{ end }} {{ $el.ElmName }} : {{ $el.ElmType -}}
{{end}}
	{{"}"}}


decoder : D.Decoder {{.Name}}
decoder =
	D.succeed {{.Name}}
{{- range .Fields }}
		|> P.required "{{ .JSONName }}" D.{{ .ElmDecoder -}}
{{end}}
`
