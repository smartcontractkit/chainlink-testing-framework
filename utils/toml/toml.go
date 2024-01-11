package toml

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func OpenTomlFileAsStruct(path string, v any) error {
	tomlFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer tomlFile.Close()
	b, _ := io.ReadAll(tomlFile)
	err = toml.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func SaveStructAsToml(v any, dirName, name string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	fileDir := fmt.Sprintf("%s/%s", pwd, dirName)
	if _, err := os.Stat(fileDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fileDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	filePath := fmt.Sprintf("%s/%s.toml", fileDir, name)
	f, err := toml.Marshal(v)
	if err != nil {
		return "", err
	}

	return filePath, os.WriteFile(filePath, f, 0600)
}
