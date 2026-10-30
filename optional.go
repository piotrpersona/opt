package opt

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNone is returned by Get when the optional is None.
var ErrNone = errors.New("the value is none")

// Optional is type wrapper for any value
type Optional[T any] struct {
	value T
	some  bool
}

// Some returns non-empty value using v.
func Some[T any](v T) Optional[T] {
	return Optional[T]{value: v, some: true}
}

// None returns empty value.
func None[T any]() Optional[T] {
	return Optional[T]{some: false}
}

// From wraps a pointer into an Optional: None if p is nil, otherwise
// Some(*p). Use it to lift a genuinely optional pointer. It copies the
// pointee, so prefer a plain pointer for references to large objects.
func From[T any](p *T) Optional[T] {
	if p == nil {
		return None[T]()
	}
	return Some(*p)
}

// FromOk wraps a comma-ok result: Some(v) if ok, otherwise None. Handy to
// lift a map lookup, type assertion or channel receive into an Optional.
func FromOk[T any](v T, ok bool) Optional[T] {
	if !ok {
		return None[T]()
	}
	return Some(v)
}

// IsSome reports whether the optional holds a value.
func (o Optional[T]) IsSome() bool {
	return o.some
}

// IsNone reports whether the optional is empty.
func (o Optional[T]) IsNone() bool {
	return !o.some
}

// IsZero reports whether the optional is empty, so a struct field tagged
// json:",omitzero" is omitted when None.
func (o Optional[T]) IsZero() bool {
	return !o.some
}

// Get returns the value if Some, and ErrNone if None.
func (o Optional[T]) Get() (T, error) {
	if !o.some {
		return o.value, ErrNone
	}
	return o.value, nil
}

// Default returns provided value if None.
func (o Optional[T]) Default(d T) T {
	if !o.some {
		return d
	}
	return o.value
}

// DefaultFunc returns the value if Some, otherwise the result of f. Unlike
// Default it does not evaluate the fallback unless the optional is None.
func (o Optional[T]) DefaultFunc(f func() T) T {
	if !o.some {
		return f()
	}
	return o.value
}

// Ptr returns a pointer to a copy of the value if Some, and nil if None.
// It is the inverse of From.
func (o Optional[T]) Ptr() *T {
	if !o.some {
		return nil
	}
	v := o.value
	return &v
}

// MarshalJSON encodes Some as the wrapped value and None as null.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.some {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON decodes null as None and any other value as Some.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = None[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}

// Value implements driver.Valuer: None becomes SQL NULL, Some the value.
func (o Optional[T]) Value() (driver.Value, error) {
	if !o.some {
		return nil, nil
	}
	return sql.Null[T]{V: o.value, Valid: true}.Value()
}

// Scan implements sql.Scanner: SQL NULL becomes None, anything else Some.
func (o *Optional[T]) Scan(src any) error {
	var n sql.Null[T]
	if err := n.Scan(src); err != nil {
		return err
	}
	if n.Valid {
		*o = Some(n.V)
	} else {
		*o = None[T]()
	}
	return nil
}

// String returns string representation of value.
func (o Optional[T]) String() string {
	if o.some {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}
