# gopt

[![Go Reference](https://pkg.go.dev/badge/github.com/piotrpersona/gopt.svg)](https://pkg.go.dev/github.com/piotrpersona/gopt)

Simple golang optional package.

## Install

```sh
go get github.com/piotrpersona/gopt
```

## Usage

```go
// Create optional variable
a := gopt.Some("value")
// create none
n := gopt.None[string]()

// check if value exists
val, ok := a.Get()
if !ok {
    // it is None
}

// or get default
val := a.Default("default")

// wrap an optional pointer: None if nil, Some(*p) otherwise
opt := gopt.From(p)

// back to a pointer: nil if None
p := opt.Ptr()
```

## JSON

`Optional[T]` marshals Some as the value and None as `null`. Tag a field
`,omitzero` to drop it entirely when None.

```go
type User struct {
    Name gopt.Optional[string] `json:"name"`
    Age  gopt.Optional[int]    `json:"age,omitzero"`
}
```

