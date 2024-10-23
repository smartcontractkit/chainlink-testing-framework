package time

import (
	"encoding/json"
	"errors"
	"time"
)

// StrDuration is JSON/TOML friendly duration that can be parsed from "1h2m0s" Go format
type StrDuration struct {
	time.Duration
}

func (d *StrDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *StrDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// MarshalText implements the text.Marshaler interface (used by toml)
func (d StrDuration) MarshalText() ([]byte, error) {
	return []byte(d.Duration.String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface (used by toml)
func (d *StrDuration) UnmarshalText(b []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(b))
	if err != nil {
		return err
	}
	return nil

}
