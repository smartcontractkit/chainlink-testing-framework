package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Setenv("SAVE_FILE", "kendrick.json")
	t.Setenv("LOG_LEVEL", "butterfly")

	config := ReadConfig()
	assert.Equal(t, config.SaveFile, "kendrick.json")
	assert.Equal(t, config.LogLevel, "butterfly")
}

func TestConfigDefaults(t *testing.T) {
	config := ReadConfig()
	assert.Equal(t, config.SaveFile, "save.json")
	assert.Equal(t, config.LogLevel, "debug")
}
