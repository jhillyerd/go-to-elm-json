# go-to-elm-json

A tool to create Elm JSON decoder pipelines from Go struct type definitions.


## Status

- [x] Support basic types: string, int, float, bool
- [x] CamelCase common acronyms
- [x] Generate Elm record
- [x] Generate decoder pipeline
- [x] Generate encoder
- [x] Usage example in README
- [x] Support slice form of basic types
- [ ] Support for optional fields
- [ ] Support for string-keyed basic type maps
- [x] Support nested structs
- [ ] Allow records to be renamed


## Install

```
go get github.com/jhillyerd/go-to-elm-json
```


## Usage

`go-to-elm-json <go source files> -- <package> <type>`

### Example

Given the file `foo/bar.go` containing:

```go
package foo

type User struct {
	Name    string   `json:"name"`
	UserID  int      `json:"userID"`
	Friends []string `json:"friends"`
	Enabled bool     `json:"enabled"`
}
```

Running: `go-to-elm-json foo/*.go -- foo User` will output:

```elm
module User exposing (User, decoder, encode)

import Json.Decode as D
import Json.Decode.Pipeline as P
import Json.Encode as E


type alias User =
    { name : String
    , userId : Int
    , friends : List String
    , enabled : Bool
    }


decoder : D.Decoder User
decoder =
    D.succeed User
        |> P.required "name" D.string
        |> P.required "userID" D.int
        |> P.required "friends" (D.list D.string)
        |> P.required "enabled" D.bool


encode : User -> E.Value
encode r =
    E.object
        [ ("name", E.string r.name)
        , ("userID", E.int r.userId)
        , ("friends", (E.list E.string) r.friends)
        , ("enabled", E.bool r.enabled)
        ]
```
