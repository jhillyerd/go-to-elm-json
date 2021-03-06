package main

var elmTemplate = `
{{- with .Record -}}
module {{.Name}} exposing ({{.Name}}, decoder, encode)

import Json.Decode as D
import Json.Decode.Pipeline as P
import Json.Encode as E



-- Generated by https://github.com/jhillyerd/go-to-elm-json


type alias {{.Name}} =
{{- range $index, $el := .Fields }}
    {{ if $index }},{{ else }}{{"{"}}{{ end }} {{ .ElmName }} : {{ .TypeDecl -}}
{{end}}
    {{"}"}}
{{- end}}
{{- range .Nested}}


type alias {{.Name}} =
{{- range $index, $el := .Fields }}
    {{ if $index }},{{ else }}{{"{"}}{{ end }} {{ .ElmName }} : {{ .TypeDecl -}}
{{end}}
    {{"}"}}
{{- end}}


{{with .Record -}}
decoder : D.Decoder {{.Name}}
decoder =
    D.succeed {{.Name}}
{{- range .Fields }}
        |> {{ .Pipeline "P" }} "{{ .JSONName }}" {{ .Decoder "D" }}{{ .Default -}}
{{end}}


encode : {{.Name}} -> E.Value
encode r =
    E.object
{{- range $index, $el := .Fields }}
        {{ if $index }},{{ else }}[{{ end }} ( "{{ .JSONName }}", {{ .Encoder "E" }} r.{{ .ElmName }} )
{{- end}}
        ]
{{- end}}
{{- range .Nested}}


{{.Decoder "D" }} : D.Decoder {{.Name}}
{{.Decoder "D" }} =
    D.succeed {{.Name}}
{{- range .Fields }}
        |> P.required "{{ .JSONName }}" {{ .Decoder "D" -}}
{{end}}


{{.Encoder "E" }} : {{.Name}} -> E.Value
{{.Encoder "E" }} r =
    E.object
{{- range $index, $el := .Fields }}
        {{ if $index }},{{ else }}[{{ end }} ( "{{ .JSONName }}", {{ .Encoder "E" }} r.{{ .ElmName }} )
{{- end}}
        ]
{{- end}}


maybe : (a -> E.Value) -> Maybe a -> E.Value
maybe encoder =
    Maybe.map encoder >> Maybe.withDefault E.null
`
