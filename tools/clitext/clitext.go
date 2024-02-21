package clitext

import "fmt"

type CliColor string

const (
	ColorGreen CliColor = "\033[32m"
	ColorYelow CliColor = "\033[33m"
	ColorRed   CliColor = "\033[31m"
	ColorReset CliColor = "\033[0m"
)

func Color(color CliColor, s string) string {
	return fmt.Sprintf("%s%s%s", color, s, ColorReset)
}
