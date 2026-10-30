# opt

[![Go Reference](https://pkg.go.dev/badge/github.com/piotrpersona/opt.svg)](https://pkg.go.dev/github.com/piotrpersona/opt)

Simple golang optional package.

## Install

```sh
go get github.com/piotrpersona/opt
```

## Usage

```go
// Create optional variable
a := opt.Some("value")
// create none
n := opt.None[string]()

// check if value exists
val, err := a.Get()
if err != nil {
    // it is None (err is opt.ErrNone)
}

// or get default
val := a.Default("default")

// wrap an optional pointer: None if nil, Some(*p) otherwise
o := opt.From(p)

// back to a pointer: nil if None
p := o.Ptr()
```

## JSON

`Optional[T]` marshals Some as the value and None as `null`. Tag a field
`,omitzero` to drop it entirely when None.

```go
type User struct {
    Name opt.Optional[string] `json:"name"`
    Age  opt.Optional[int]    `json:"age,omitzero"`
}
```
