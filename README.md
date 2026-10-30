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
// lazy default, only evaluated when None
val := a.DefaultFunc(func() string { return expensive() })

// wrap an optional pointer: None if nil, Some(*p) otherwise
o := opt.From(p)
// back to a pointer: nil if None
p := o.Ptr()

// lift a comma-ok result (map lookup, type assertion, channel receive)
v, ok := m[key]
o := opt.FromOk(v, ok)
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

## SQL

`Optional[T]` implements `sql.Scanner` and `driver.Valuer`, so it maps to a
nullable column: None is `NULL`, Some is the value.

```go
var name opt.Optional[string]
row.Scan(&name)                       // NULL -> None
db.Exec(`... VALUES ($1)`, name)      // None -> NULL
```
