package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetGoProjectName(path string) (string, error) {
	// Walk up the directory structure to find go.mod
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}
	dir := absPath
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir { // Reached the root without finding go.mod
			return "", fmt.Errorf("go.mod not found in project path, started at %s, ended at %s", path, dir)
		}
		dir = parent
	}

	// Read go.mod to extract the module path
	goModPath := filepath.Join(dir, "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	for _, line := range strings.Split(string(goModData), "\n") {
		if strings.HasPrefix(line, "module ") {
			goProject := strings.TrimSpace(strings.TrimPrefix(line, "module "))
			relativePath := strings.TrimPrefix(path, dir)
			relativePath = strings.TrimLeft(relativePath, string(os.PathSeparator))
			return filepath.Join(goProject, relativePath), nil
		}
	}

	return "", fmt.Errorf("module path not found in go.mod")
}
