package github

import (
	"fmt"
)

func StartGroup(title string) {
	fmt.Printf("::group:: %s", title)
}
func EndGroup() {
	fmt.Println("::endgroup::")
}
