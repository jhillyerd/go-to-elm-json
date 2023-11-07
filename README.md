# go-to-elm-json
[![Build and Test](https://github.com/jhillyerd/go-to-elm-json/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/jhillyerd/go-to-elm-json/actions/workflows/build-and-test.yml)
[![Coverage Status](https://coveralls.io/repos/github/jhillyerd/go-to-elm-json/badge.svg?branch=main)](https://coveralls.io/github/jhillyerd/go-to-elm-json?branch=main)

A tool to create Elm JSON decoder pipelines from Go struct type definitions.


## Status

Useful, but not feature complete.

- [x] Support basic types: string, int, float, bool
- [x] CamelCase common acronyms
- [x] Generate Elm record
- [x] Generate decoder pipeline
- [x] Generate encoder
- [x] Usage example in README
- [x] Support slice form of basic types
- [x] Support for optional fields
- [x] Support nested structs
- [x] Support nullable pointer types
- [x] Allow records to be renamed
- [ ] Handle `json:"-"` correctly
- [ ] Specify module name
- [ ] Support for string-keyed basic type maps


## Install

```
go get github.com/jhillyerd/go-to-elm-json
```


## Usage

`go-to-elm-json <go source files> -- <package> <go type:elm name>`

### Example

Given the file `foo/bar.go` containing:

```go
package foo

type UserJSON struct {
	Name    string   `json:"name"`
	UserID  int      `json:"userID"`
	Friends []string `json:"friends"`
	Enabled bool     `json:"enabled"`
}
```

Running: `go-to-elm-json foo/*.go -- foo UserJSON:User` will output:

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


## Contributing

PRs welcome, please:

- Base your work off of the `development` branch, and target pull requests to
  the same.
- Run the unit tests before filing a PR.  `make` will run tests and lint.
- Include unit tests for your changes.


[Build Status]:    https://travis-ci.org/jhillyerd/go-to-elm-json
[Coverage Status]: https://coveralls.io/github/jhillyerd/go-to-elm-json?branch=master
