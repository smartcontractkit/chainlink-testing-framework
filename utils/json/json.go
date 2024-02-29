package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
