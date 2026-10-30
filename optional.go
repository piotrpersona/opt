package gopt

import (
	"encoding/json"
	"fmt"
)

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

// Get returns the value and true if Some, and the zero value and false if None.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.some
}

// Default returns provided value if None.
func (o Optional[T]) Default(d T) T {
	if !o.some {
		return d
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

// String returns string representation of value.
func (o Optional[T]) String() string {
	if o.some {
		return fmt.Sprintf("Some(%v)", o.value)
	}
	return "None"
}
