package main

var elmTemplate = `module {{.Name}} exposing ({{.Name}})

import Json.Decode as D
import Json.Decode.Pipeline as P
import Json.Encode as E

type alias {{.Name}} =
{{- range $index, $el := .Fields }}
	{{ if $index }},{{ else }}{{"{"}}{{ end }} {{ .ElmName }} : {{ .ElmType -}}
{{end}}
	{{"}"}}


decoder : D.Decoder {{.Name}}
decoder =
	D.succeed {{.Name}}
{{- range .Fields }}
		|> P.required "{{ .JSONName }}" D.{{ .ElmDecoder -}}
{{end}}

encode : {{.Name}} -> Value
encode r =
	E.object
{{- range $index, $el := .Fields }}
		{{ if $index }},{{ else }}[{{ end }} ("{{ .JSONName }}", E.{{ .ElmDecoder }} r.{{ .ElmName -}}
{{end}}
		]
`
