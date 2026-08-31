package gate

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// req decodes m[key] into *dst. It returns an error when key is absent from m
// or explicitly JSON null, so a caller can never mistake absence for a zero
// value (H1) — json.Unmarshal treats "null" as a documented no-op for
// non-pointer targets (string, bool, int, ...), so without this check a
// required field sent as null would silently pass through as its zero value.
func req[T any](m map[string]json.RawMessage, key string, dst *T) error {
	raw, ok := m[key]
	if !ok || isJSONNull(raw) {
		return fmt.Errorf("required field %q is absent", key)
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("field %q: %w", key, err)
	}
	return nil
}

func isJSONNull(raw json.RawMessage) bool {
	return string(bytes.TrimSpace(raw)) == "null"
}

// opt decodes m[key] into *dst when present, leaving *dst untouched when key is absent.
func opt[T any](m map[string]json.RawMessage, key string, dst *T) error {
	raw, ok := m[key]
	if !ok {
		return nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return fmt.Errorf("field %q: %w", key, err)
	}
	return nil
}
