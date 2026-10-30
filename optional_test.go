package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/piotrpersona/opt"
	"github.com/stretchr/testify/require"
)

func Test_Optional(t *testing.T) {
	t.Parallel()

	t.Run("Test some", func(t *testing.T) {
		t.Parallel()

		o := opt.Some("car")
		val, err := o.Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
	})
	t.Run("Test none", func(t *testing.T) {
		t.Parallel()

		o := opt.None[string]()
		val, err := o.Get()
		require.Equal(t, "", val)
		require.ErrorIs(t, err, opt.ErrNone)
	})
	t.Run("Test default", func(t *testing.T) {
		t.Parallel()

		o := opt.None[string]()
		val := o.Default("boat")
		require.Equal(t, "boat", val)
	})
	t.Run("Test IsSome and IsNone", func(t *testing.T) {
		t.Parallel()

		require.True(t, opt.Some("car").IsSome())
		require.False(t, opt.Some("car").IsNone())
		require.True(t, opt.None[string]().IsNone())
		require.False(t, opt.None[string]().IsSome())
	})
	t.Run("Test chaining on returned value", func(t *testing.T) {
		t.Parallel()

		val, err := opt.Some("car").Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
		require.Equal(t, "boat", opt.None[string]().Default("boat"))
	})
	t.Run("Test From nil is None", func(t *testing.T) {
		t.Parallel()

		o := opt.From[string](nil)
		require.True(t, o.IsNone())
	})
	t.Run("Test From value is Some", func(t *testing.T) {
		t.Parallel()

		s := "car"
		o := opt.From(&s)
		require.True(t, o.IsSome())
		val, err := o.Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)
	})
	t.Run("Test Ptr round-trips and copies", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, opt.None[string]().Ptr())

		s := "car"
		p := opt.From(&s).Ptr()
		require.NotNil(t, p)
		require.Equal(t, "car", *p)
		require.NotSame(t, &s, p)
	})
	t.Run("Test String", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, "Some(none)", opt.Some("none").String())
		require.Equal(t, "None", opt.None[string]().String())
	})
	t.Run("Test JSON round-trip", func(t *testing.T) {
		t.Parallel()

		type payload struct {
			Name opt.Optional[string] `json:"name"`
			Age  opt.Optional[int]    `json:"age,omitzero"`
		}

		in := payload{Name: opt.Some("car"), Age: opt.None[int]()}
		data, err := json.Marshal(in)
		require.NoError(t, err)
		require.JSONEq(t, `{"name":"car"}`, string(data))

		var out payload
		require.NoError(t, json.Unmarshal(data, &out))
		require.Equal(t, in, out)
	})
	t.Run("Test JSON null decodes to None", func(t *testing.T) {
		t.Parallel()

		var o opt.Optional[string]
		require.NoError(t, json.Unmarshal([]byte(`null`), &o))
		require.True(t, o.IsNone())

		data, err := json.Marshal(opt.None[string]())
		require.NoError(t, err)
		require.Equal(t, "null", string(data))
	})
	t.Run("Test FromOk", func(t *testing.T) {
		t.Parallel()

		m := map[string]int{"car": 1}
		v, ok := m["car"]
		require.True(t, opt.FromOk(v, ok).IsSome())
		v, ok = m["boat"]
		require.True(t, opt.FromOk(v, ok).IsNone())
	})
	t.Run("Test DefaultFunc is lazy", func(t *testing.T) {
		t.Parallel()

		calls := 0
		fallback := func() string { calls++; return "boat" }

		require.Equal(t, "car", opt.Some("car").DefaultFunc(fallback))
		require.Equal(t, 0, calls)
		require.Equal(t, "boat", opt.None[string]().DefaultFunc(fallback))
		require.Equal(t, 1, calls)
	})
	t.Run("Test SQL Value", func(t *testing.T) {
		t.Parallel()

		v, err := opt.Some("car").Value()
		require.NoError(t, err)
		require.Equal(t, "car", v)

		v, err = opt.None[string]().Value()
		require.NoError(t, err)
		require.Nil(t, v)
	})
	t.Run("Test SQL Scan", func(t *testing.T) {
		t.Parallel()

		var o opt.Optional[string]
		require.NoError(t, o.Scan("car"))
		require.True(t, o.IsSome())
		val, err := o.Get()
		require.NoError(t, err)
		require.Equal(t, "car", val)

		require.NoError(t, o.Scan(nil))
		require.True(t, o.IsNone())
	})
}
