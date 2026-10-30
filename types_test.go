package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/piotrpersona/opt"
)

// jsonRoundTrip checks that Some(v) marshals to want, round-trips back to
// Some(v), and that the zero value of T is still encoded as a value, not null.
func jsonRoundTrip[T comparable](t *testing.T, v T, want string) {
	t.Helper()

	data, err := json.Marshal(opt.Some(v))
	noError(t, err)
	equal(t, string(data), want)

	var o opt.Optional[T]
	noError(t, json.Unmarshal(data, &o))
	isTrue(t, o.IsSome(), "decoded value should be Some")
	got, err := o.Get()
	noError(t, err)
	equal(t, got, v)
}

// sqlRoundTrip checks that Some(v).Value() scans back to Some(v).
func sqlRoundTrip[T comparable](t *testing.T, v T) {
	t.Helper()

	dv, err := opt.Some(v).Value()
	noError(t, err)

	var o opt.Optional[T]
	noError(t, o.Scan(dv))
	got, err := o.Get()
	noError(t, err)
	equal(t, got, v)
}

func Test_JSON_Types(t *testing.T) {
	t.Parallel()

	type point struct{ X, Y int }

	t.Run("int", func(t *testing.T) { jsonRoundTrip(t, 42, "42") })
	t.Run("int zero", func(t *testing.T) { jsonRoundTrip(t, 0, "0") })
	t.Run("bool true", func(t *testing.T) { jsonRoundTrip(t, true, "true") })
	t.Run("bool zero", func(t *testing.T) { jsonRoundTrip(t, false, "false") })
	t.Run("string", func(t *testing.T) { jsonRoundTrip(t, "car", `"car"`) })
	t.Run("string zero", func(t *testing.T) { jsonRoundTrip(t, "", `""`) })
	t.Run("float", func(t *testing.T) { jsonRoundTrip(t, 3.14, "3.14") })
	t.Run("struct", func(t *testing.T) { jsonRoundTrip(t, point{1, 2}, `{"X":1,"Y":2}`) })
}

func Test_SQL_Types(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) { sqlRoundTrip(t, 42) })
	t.Run("int zero", func(t *testing.T) { sqlRoundTrip(t, 0) })
	t.Run("int64", func(t *testing.T) { sqlRoundTrip(t, int64(7)) })
	t.Run("bool", func(t *testing.T) { sqlRoundTrip(t, true) })
	t.Run("float", func(t *testing.T) { sqlRoundTrip(t, 2.5) })
	t.Run("string", func(t *testing.T) { sqlRoundTrip(t, "car") })
	t.Run("string zero", func(t *testing.T) { sqlRoundTrip(t, "") })
}

// Some(zero value) must stay present under json:",omitzero"; only None is dropped.
func Test_JSON_Omitzero(t *testing.T) {
	t.Parallel()

	type payload struct {
		N opt.Optional[int] `json:"n,omitzero"`
	}

	data, err := json.Marshal(payload{N: opt.Some(0)})
	noError(t, err)
	equal(t, string(data), `{"n":0}`)

	data, err = json.Marshal(payload{N: opt.None[int]()})
	noError(t, err)
	equal(t, string(data), `{}`)
}
