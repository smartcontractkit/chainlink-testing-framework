package tools

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// ProjectRoot Root folder of this project
	ProjectRoot          = filepath.Join(filepath.Dir(b), "/..")
	ContractsDir         = filepath.Join(ProjectRoot, "contracts")
	EthereumContractsDir = filepath.Join(ContractsDir, "ethereum")
)
