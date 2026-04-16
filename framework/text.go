package framework

import "fmt"

func RedText(text string, args ...any) string {
	return fmt.Sprintf("\033[31m%s\033[0m", fmt.Sprintf(text, args...))
}
