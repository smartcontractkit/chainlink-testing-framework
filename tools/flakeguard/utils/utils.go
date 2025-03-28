package utils

import (
	"path/filepath"
)

func UniqueDirectories(files []string) []string {
	dirSet := make(map[string]struct{})
	for _, file := range files {
		dirname := filepath.Dir(file)
		dirSet[dirname] = struct{}{}
	}
	var dirs []string
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	return dirs
}

func Deduplicate(items []string) []string {
	seen := make(map[string]struct{})
	var uniqueItems []string
	for _, item := range items {
		if _, found := seen[item]; !found {
			seen[item] = struct{}{}
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}

func ResolveFullPath(projectPath string) (string, error) {
	if filepath.IsAbs(projectPath) {
		return filepath.Clean(projectPath), nil
	}
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
