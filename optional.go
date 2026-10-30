package gopt

import (
	"errors"
	"fmt"
)

// NoneErr will be returned if the optional value is none.
var NoneErr = errors.New("the value is none")

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

// FromPtr wraps a pointer into an Optional: None if p is nil, otherwise
// Some(*p). Use it to lift a genuinely optional pointer. It copies the
// pointee, so prefer a plain pointer for references to large objects.
func FromPtr[T any](p *T) Optional[T] {
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

// Get returns value if Some, and error if None.
func (o Optional[T]) Get() (T, error) {
	if !o.some {
		return o.value, NoneErr
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

// Ptr returns a pointer to a copy of the value if Some, and nil if None.
// It is the inverse of FromPtr.
func (o Optional[T]) Ptr() *T {
	if !o.some {
		return nil
	}
	v := o.value
	return &v
}

// String returns string representation of value.
func (o Optional[T]) String() string {
	if o.some {
		return fmt.Sprintf("%v", o.value)
	}
	return "none"
}
