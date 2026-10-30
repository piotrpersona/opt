package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/piotrpersona/opt"
)

func Test_Optional(t *testing.T) {
	t.Parallel()

	t.Run("Test some", func(t *testing.T) {
		t.Parallel()

		val, err := opt.Some("car").Get()
		noError(t, err)
		equal(t, val, "car")
	})
	t.Run("Test none", func(t *testing.T) {
		t.Parallel()

		val, err := opt.None[string]().Get()
		errorIs(t, err, opt.ErrNone)
		equal(t, val, "")
	})
	t.Run("Test default", func(t *testing.T) {
		t.Parallel()

		equal(t, opt.None[string]().Default("boat"), "boat")
	})
	t.Run("Test IsSome and IsNone", func(t *testing.T) {
		t.Parallel()

		isTrue(t, opt.Some("car").IsSome(), "Some should be IsSome")
		isTrue(t, !opt.Some("car").IsNone(), "Some should not be IsNone")
		isTrue(t, opt.None[string]().IsNone(), "None should be IsNone")
		isTrue(t, !opt.None[string]().IsSome(), "None should not be IsSome")
	})
	t.Run("Test chaining on returned value", func(t *testing.T) {
		t.Parallel()

		val, err := opt.Some("car").Get()
		noError(t, err)
		equal(t, val, "car")
		equal(t, opt.None[string]().Default("boat"), "boat")
	})
	t.Run("Test From nil is None", func(t *testing.T) {
		t.Parallel()

		isTrue(t, opt.From[string](nil).IsNone(), "From(nil) should be None")
	})
	t.Run("Test From value is Some", func(t *testing.T) {
		t.Parallel()

		s := "car"
		val, err := opt.From(&s).Get()
		noError(t, err)
		equal(t, val, "car")
	})
	t.Run("Test Ptr round-trips and copies", func(t *testing.T) {
		t.Parallel()

		isTrue(t, opt.None[string]().Ptr() == nil, "None.Ptr() should be nil")
		s := "car"
		p := opt.From(&s).Ptr()
		isTrue(t, p != nil, "Some.Ptr() should not be nil")
		equal(t, *p, "car")
		isTrue(t, p != &s, "Ptr() should return a copy, not the original pointer")
	})
	t.Run("Test String", func(t *testing.T) {
		t.Parallel()

		equal(t, opt.Some("none").String(), "Some(none)")
		equal(t, opt.None[string]().String(), "None")
	})
	t.Run("Test JSON round-trip", func(t *testing.T) {
		t.Parallel()

		type payload struct {
			Name opt.Optional[string] `json:"name"`
			Age  opt.Optional[int]    `json:"age,omitzero"`
		}

		in := payload{Name: opt.Some("car"), Age: opt.None[int]()}
		data, err := json.Marshal(in)
		noError(t, err)
		equal(t, string(data), `{"name":"car"}`)

		var out payload
		noError(t, json.Unmarshal(data, &out))
		equal(t, out, in)
	})
	t.Run("Test JSON null decodes to None", func(t *testing.T) {
		t.Parallel()

		var o opt.Optional[string]
		noError(t, json.Unmarshal([]byte(`null`), &o))
		isTrue(t, o.IsNone(), "null should decode to None")

		data, err := json.Marshal(opt.None[string]())
		noError(t, err)
		equal(t, string(data), "null")
	})
	t.Run("Test FromOk", func(t *testing.T) {
		t.Parallel()

		m := map[string]int{"car": 1}
		v, ok := m["car"]
		isTrue(t, opt.FromOk(v, ok).IsSome(), "present key should be Some")
		v, ok = m["boat"]
		isTrue(t, opt.FromOk(v, ok).IsNone(), "missing key should be None")
	})
	t.Run("Test DefaultFunc is lazy", func(t *testing.T) {
		t.Parallel()

		calls := 0
		fallback := func() string { calls++; return "boat" }

		equal(t, opt.Some("car").DefaultFunc(fallback), "car")
		equal(t, calls, 0)
		equal(t, opt.None[string]().DefaultFunc(fallback), "boat")
		equal(t, calls, 1)
	})
	t.Run("Test SQL Value", func(t *testing.T) {
		t.Parallel()

		v, err := opt.Some("car").Value()
		noError(t, err)
		isTrue(t, v == "car", "Some.Value() should be the value")

		v, err = opt.None[string]().Value()
		noError(t, err)
		isTrue(t, v == nil, "None.Value() should be nil")
	})
	t.Run("Test SQL Scan", func(t *testing.T) {
		t.Parallel()

		var o opt.Optional[string]
		noError(t, o.Scan("car"))
		val, err := o.Get()
		noError(t, err)
		equal(t, val, "car")

		noError(t, o.Scan(nil))
		isTrue(t, o.IsNone(), "scan nil should be None")
	})
}
