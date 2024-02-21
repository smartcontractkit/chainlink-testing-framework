package github

import "fmt"

func StartGroup(title string) {
	fmt.Printf("::group::%s\n", title)
}
func EndGroup() {
	fmt.Println("::endgroup::")
}
