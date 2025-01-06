package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Setenv("SAVE_FILE", "kendrick.json")
	t.Setenv("LOG_LEVEL", "butterfly")

	config := readConfig()
	assert.Equal(t, "kendrick.json", config.SaveFile)
	assert.Equal(t, "butterfly", config.LogLevel)
}

func TestConfigDefaults(t *testing.T) {
	config := readConfig()
	assert.Equal(t, "save.json", config.SaveFile)
	assert.Equal(t, "debug", config.LogLevel)
}
