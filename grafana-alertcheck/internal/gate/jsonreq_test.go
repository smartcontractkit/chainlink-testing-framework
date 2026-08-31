package gate

import (
	"encoding/json"
	"testing"
)

func rawMap(t *testing.T, jsonObj string) map[string]json.RawMessage {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonObj), &m); err != nil {
		t.Fatalf("rawMap: %v", err)
	}
	return m
}

func TestReq(t *testing.T) {
	m := rawMap(t, `{"present":"hello","wrongtype":123,"nullval":null}`)

	t.Run("present key decodes", func(t *testing.T) {
		var s string
		if err := req(m, "present", &s); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s != "hello" {
			t.Errorf("got %q, want hello", s)
		}
	})

	t.Run("absent key errors", func(t *testing.T) {
		var s string
		if err := req(m, "missing", &s); err == nil {
			t.Fatalf("expected an error, got none")
		}
	})

	t.Run("wrong type errors", func(t *testing.T) {
		var s string
		if err := req(m, "wrongtype", &s); err == nil {
			t.Fatalf("expected an error, got none")
		}
	})

	t.Run("explicit JSON null errors, never a zero value", func(t *testing.T) {
		var s string
		err := req(m, "nullval", &s)
		if err == nil {
			t.Fatalf("expected an error, got none (s=%q) — a null required field must not silently become a zero value", s)
		}
	})
}

func TestOpt(t *testing.T) {
	m := rawMap(t, `{"present":"hello","wrongtype":123,"nullval":null}`)

	t.Run("present key decodes", func(t *testing.T) {
		var s string
		if err := opt(m, "present", &s); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s != "hello" {
			t.Errorf("got %q, want hello", s)
		}
	})

	t.Run("absent key leaves dst untouched", func(t *testing.T) {
		s := "unchanged"
		if err := opt(m, "missing", &s); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s != "unchanged" {
			t.Errorf("got %q, want unchanged", s)
		}
	})

	t.Run("wrong type errors", func(t *testing.T) {
		var s string
		if err := opt(m, "wrongtype", &s); err == nil {
			t.Fatalf("expected an error, got none")
		}
	})

	t.Run("explicit JSON null leaves dst at its zero value", func(t *testing.T) {
		var s string
		if err := opt(m, "nullval", &s); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s != "" {
			t.Errorf("got %q, want empty string", s)
		}
	})
}
