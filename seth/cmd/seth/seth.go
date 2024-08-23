package main

import (
	"os"

	"github.com/smartcontractkit/chainlink-testing-framework/seth/cmd"
)

func main() {
	if err := seth.RunCLI(os.Args); err != nil {
		panic(err)
	}
}
