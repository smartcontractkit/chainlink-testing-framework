package internal

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	ProjectRoot = filepath.Join(filepath.Dir(b), "/..")
)
