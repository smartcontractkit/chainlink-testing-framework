package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

func OpenJsonFileAsStruct(path string, v any) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	b, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func SaveStructAsJson(v any, dirName, name string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	confDir := fmt.Sprintf("%s/%s", pwd, dirName)
	if _, err := os.Stat(confDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(confDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	confPath := fmt.Sprintf("%s/%s.json", confDir, name)
	f, _ := json.MarshalIndent(v, "", "   ")
	err = os.WriteFile(confPath, f, 0600)

	return confPath, err
}

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
