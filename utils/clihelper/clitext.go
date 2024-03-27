package clihelper

import (
	"fmt"
	"strings"
)

type CliColor string

const (
	ColorGreen  CliColor = "\033[0;32m"
	ColorYellow CliColor = "\033[0;33m"
	ColorRed    CliColor = "\033[0;31m"
	ColorReset  CliColor = "\033[0m"
)

func Color(color CliColor, s string) string {
	return fmt.Sprintf("%s%s %s\n", color, strings.TrimRight(s, "\n"), ColorReset)
}
