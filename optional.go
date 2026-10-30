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

// String returns string representation of value.
func (o Optional[T]) String() string {
	if o.some {
		return fmt.Sprintf("%v", o.value)
	}
	return "none"
}
