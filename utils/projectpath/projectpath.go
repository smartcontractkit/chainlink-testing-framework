// Package utils contains some common paths used in configuration and tests
package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// ProjectRoot Root folder of this project
	ProjectRoot = filepath.Join(filepath.Dir(b), "/../..")
	// ChartsRoot test suite root
	ChartsRoot = filepath.Join(ProjectRoot, "charts")
	// K8sRoot test suite soak root
	K8sRoot = filepath.Join(ProjectRoot, "soak")
	// PresetRoot root folder for environments preset
	PresetRoot = filepath.Join(ProjectRoot, "preset")
	// ContractsDir path to our contracts
	ContractsDir = filepath.Join(ProjectRoot, "contracts")
	// EthereumContractsDir path to our ethereum contracts
	EthereumContractsDir = filepath.Join(ContractsDir, "ethereum")
	// MirrorDir path to our ecr mirror helpers
	MirrorDir = filepath.Join(ProjectRoot, "mirror")
)
