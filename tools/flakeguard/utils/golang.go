package utils

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func GetPackageNames(dirs []string, repoPath string) []string {
	var packageNames []string
	for _, dir := range dirs {
		cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", ".")
		cmd.Dir = repoPath + dir
		out, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error getting package name for directory %s: %s\n", dir, err)
			continue
		}
		packageName := strings.TrimSpace(string(out))
		if packageName != "" {
			packageNames = append(packageNames, packageName)
		}
	}
	return packageNames
}

// FindDependentPackages returns a list of packages that depend on any of the specified packages.
func FindDependentPackages(repoPath string, targetPackages []string) ([]string, error) {
	dependentPackages := make([]string, 0)

	// Execute 'go list' to find all Go test packages with their dependencies in the current module
	cmd := exec.Command("go", "list", "-f", `{{if .TestGoFiles}}{{.ImportPath}} {{join .Imports " "}} {{join .TestImports " "}}{{end}}`, "./...")
	cmd.Dir = repoPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running 'go list': %w", err)
	}

	// Scan each line to determine if it imports any of the target packages
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue // Skip if there aren't enough parts to include imports
		}
		packageName := parts[0]
		imports := strings.Fields(parts[1] + " " + parts[2])

		for _, target := range targetPackages {
			for _, imp := range imports {
				if imp == target {
					dependentPackages = append(dependentPackages, packageName)
					break
				}
			}
		}
	}

	return dependentPackages, nil
}
