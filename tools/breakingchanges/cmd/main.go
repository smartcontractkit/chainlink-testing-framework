package main

import (
	"flag"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/breakingchanges"
)

func main() {
	pathFlag := flag.String("path", ".", "Path to start searching for go.mod files")
	flag.Parse()

	breakingchanges.DetectBreakingChanges(*pathFlag)
}
