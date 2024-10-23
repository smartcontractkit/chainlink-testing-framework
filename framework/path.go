package framework

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	PathRoot   = filepath.Join(filepath.Dir(b), ".")
)
