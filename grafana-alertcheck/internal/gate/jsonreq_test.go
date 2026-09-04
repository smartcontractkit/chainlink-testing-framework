package gate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func rawMap(t *testing.T, jsonObj string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(jsonObj), &m))
	return m
}

func TestReq(t *testing.T) {
	m := rawMap(t, `{"present":"hello","wrongtype":123,"nullval":null}`)

	t.Run("present key decodes", func(t *testing.T) {
		var s string
		require.NoError(t, req(m, "present", &s))
		require.Equal(t, "hello", s)
	})

	t.Run("absent key errors", func(t *testing.T) {
		var s string
		require.Error(t, req(m, "missing", &s))
	})

	t.Run("wrong type errors", func(t *testing.T) {
		var s string
		require.Error(t, req(m, "wrongtype", &s))
	})

	t.Run("explicit JSON null errors, never a zero value", func(t *testing.T) {
		var s string
		err := req(m, "nullval", &s)
		require.Error(t, err, "a null required field must not silently become a zero value")
	})
}

func TestOpt(t *testing.T) {
	m := rawMap(t, `{"present":"hello","wrongtype":123,"nullval":null}`)

	t.Run("present key decodes", func(t *testing.T) {
		var s string
		require.NoError(t, opt(m, "present", &s))
		require.Equal(t, "hello", s)
	})

	t.Run("absent key leaves dst untouched", func(t *testing.T) {
		s := "unchanged"
		require.NoError(t, opt(m, "missing", &s))
		require.Equal(t, "unchanged", s)
	})

	t.Run("wrong type errors", func(t *testing.T) {
		var s string
		require.Error(t, opt(m, "wrongtype", &s))
	})

	t.Run("explicit JSON null leaves dst at its zero value", func(t *testing.T) {
		var s string
		require.NoError(t, opt(m, "nullval", &s))
		require.Equal(t, "", s)
	})
}
