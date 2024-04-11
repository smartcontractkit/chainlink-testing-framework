package main

import (
	"testing"
)

func TestTidy(t *testing.T) {
	project := "."

	Main(project, false, false)
}
