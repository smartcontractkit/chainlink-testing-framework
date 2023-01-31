package utils

import (
	"encoding/json"
	"errors"
	"time"
)

// JSONStrDuration is JSON friendly duration that can be parsed from "1h2m0s" Go format
type JSONStrDuration struct {
	time.Duration
}

func (d *JSONStrDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *JSONStrDuration) UnmarshalJSON(b []byte) error {
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
