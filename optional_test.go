package gopt_test

import (
	"testing"

	"github.com/piotrpersona/gopt"
	"github.com/stretchr/testify/require"
)

func Test_Optional(t *testing.T) {
	t.Parallel()

	t.Run("Test some", func(t *testing.T) {
		t.Parallel()

		opt := gopt.Some("car")
		val, err := opt.Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
	})
	t.Run("Test none", func(t *testing.T) {
		t.Parallel()

		opt := gopt.None[string]()
		val, err := opt.Get()
		require.Equal(t, "", val)
		require.Error(t, err)
		require.ErrorIs(t, err, gopt.NoneErr)
	})
	t.Run("Test default", func(t *testing.T) {
		t.Parallel()

		opt := gopt.None[string]()
		val := opt.Default("boat")
		require.Equal(t, "boat", val)
	})
	t.Run("Test IsSome and IsNone", func(t *testing.T) {
		t.Parallel()

		require.True(t, gopt.Some("car").IsSome())
		require.False(t, gopt.Some("car").IsNone())
		require.True(t, gopt.None[string]().IsNone())
		require.False(t, gopt.None[string]().IsSome())
	})
	t.Run("Test chaining on returned value", func(t *testing.T) {
		t.Parallel()

		val, err := gopt.Some("car").Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
		require.Equal(t, "boat", gopt.None[string]().Default("boat"))
	})
	t.Run("Test From nil is None", func(t *testing.T) {
		t.Parallel()

		opt := gopt.From[string](nil)
		require.True(t, opt.IsNone())
		_, err := opt.Get()
		require.ErrorIs(t, err, gopt.NoneErr)
	})
	t.Run("Test From value is Some", func(t *testing.T) {
		t.Parallel()

		s := "car"
		opt := gopt.From(&s)
		require.True(t, opt.IsSome())
		val, err := opt.Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
	})
	t.Run("Test Ptr round-trips and copies", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, gopt.None[string]().Ptr())

		s := "car"
		p := gopt.From(&s).Ptr()
		require.NotNil(t, p)
		require.Equal(t, "car", *p)
		require.NotSame(t, &s, p)
	})
}
