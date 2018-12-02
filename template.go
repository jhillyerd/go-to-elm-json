package main

var elmTemplate = `module {{.Name}} exposing ({{.Name}}, decoder, encode)

import Json.Decode as D
import Json.Decode.Pipeline as P
import Json.Encode as E


type alias {{.Name}} =
{{- range $index, $el := .Fields }}
    {{ if $index }},{{ else }}{{"{"}}{{ end }} {{ .ElmName }} : {{ .ElmType.Name -}}
{{end}}
    {{"}"}}


decoder : D.Decoder {{.Name}}
decoder =
    D.succeed {{.Name}}
{{- range .Fields }}
        |> P.required "{{ .JSONName }}" {{ .ElmType.Codec "D" -}}
{{end}}


encode : {{.Name}} -> E.Value
encode r =
    E.object
{{- range $index, $el := .Fields }}
        {{ if $index }},{{ else }}[{{ end }} ("{{ .JSONName }}", {{ .ElmType.Codec "E" }} r.{{ .ElmName }})
{{- end}}
        ]
`
