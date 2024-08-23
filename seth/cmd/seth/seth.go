package main

import (
	"os"

	"github.com/smartcontractkit/seth/cmd"
)

func main() {
	if err := seth.RunCLI(os.Args); err != nil {
		panic(err)
	}
}
