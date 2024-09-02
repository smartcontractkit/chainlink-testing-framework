package main

import (
	"os"

	seth "github.com/smartcontractkit/chainlink-testing-framework/lib/seth/cmd"
)

func main() {
	if err := seth.RunCLI(os.Args); err != nil {
		panic(err)
	}
}
