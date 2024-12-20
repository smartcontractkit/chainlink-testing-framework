package golang

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/utils"
)

type Package struct {
	Dir          string   `json:"Dir"`
	ImportPath   string   `json:"ImportPath"`
	Root         string   `json:"Root"`
	Deps         []string `json:"Deps"`
	TestImports  []string `json:"TestImports"`
	XTestImports []string `json:"XTestImports"`
	GoFiles      []string `json:"GoFiles"`
	TestGoFiles  []string `json:"TestGoFiles"`
	XTestGoFiles []string `json:"XTestGoFiles"`
	EmbedFiles   []string `json:"EmbedFiles"`
}

type DepMap map[string][]string

func GoList() (*utils.CmdOutput, error) {
	return utils.ExecuteCmd("go", "list", "-json", "./...")
}

// ParsePackages parses the output of `go list -json ./...` and returns a slice of Package structs
func ParsePackages(goList bytes.Buffer) ([]Package, error) {
	var packages []Package
	scanner := bufio.NewScanner(&goList)
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		buffer.WriteString(line)

		if line == "}" {
			var pkg Package
			if err := json.Unmarshal(buffer.Bytes(), &pkg); err != nil {
				return nil, err
			}
			packages = append(packages, pkg)
			buffer.Reset()
		}
	}

	err := scanner.Err()
	return packages, err
}

func GetGoDepMap(packages []Package) DepMap {
	depMap := make(map[string][]string)
	for _, pkg := range packages {
		for _, dep := range pkg.Deps {
			depMap[dep] = append(depMap[dep], pkg.ImportPath)
		}
		for _, dep := range pkg.TestImports {
			depMap[dep] = append(depMap[dep], pkg.ImportPath)
		}
		for _, dep := range pkg.XTestImports {
			depMap[dep] = append(depMap[dep], pkg.ImportPath)
		}
	}
	return depMap
}

// GetGoFileMap returns a map of go files to the packages that embed them
func GetGoFileMap(packages []Package, includeTestFiles bool) map[string][]string {
	// Build dependency map
	fileMap := make(map[string][]string)
	for _, pkg := range packages {
		addToMap(pkg, pkg.GoFiles, fileMap)
		addToMap(pkg, pkg.EmbedFiles, fileMap)
		if includeTestFiles {
			addToMap(pkg, pkg.TestGoFiles, fileMap)
			addToMap(pkg, pkg.XTestGoFiles, fileMap)
		}

	}
	return fileMap
}

func GetPackageNames(dirs []string) []string {
	var packageNames []string
	for _, dir := range dirs {
		cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", ".")
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			log.Error().Str("directory", dir).Err(err).Msg("Error getting package name")
			continue
		}
		packageName := strings.TrimSpace(string(out))
		if packageName != "" {
			packageNames = append(packageNames, packageName)
		}
	}
	return packageNames
}

func addToMap(pkg Package, files []string, fileMap map[string][]string) {
	for _, file := range files {
		separator := "/"
		if strings.Contains(pkg.Root, "\\") {
			// windows path replace
			separator = "\\"
		}
		path := strings.Replace(pkg.Dir, fmt.Sprintf("%s%s", pkg.Root, separator), "", 1)
		key := fmt.Sprintf("%s/%s", path, file)
		if pkg.Dir == pkg.Root {
			key = file
		}
		// Multiple packages can embed the same file so we need to take that into account
		if keys, exists := fileMap[key]; exists {
			fileMap[key] = append(keys, pkg.ImportPath)
		} else {
			keys := []string{
				pkg.ImportPath,
			}
			fileMap[key] = keys
		}
	}
}

//nolint:revive
func FindAffectedPackages(pkg string, depMap DepMap, externalPackage bool, maxDepth int) []string {
	visited := make(map[string]bool)
	var affected []string

	var dfs func(string, int)
	dfs = func(p string, depthLeft int) {
		if visited[p] {
			return
		}

		visited[p] = true
		// exclude the package itself if it is an external package
		if !(externalPackage && p == pkg) {
			affected = append(affected, p)
		}
		d := depthLeft - 1
		if d != 0 {
			for _, dep := range depMap[p] {
				dfs(dep, d)
			}
		}
	}

	depth := maxDepth
	// depth is zero then we want infinite recursion, set this to -1 to enable this
	if depth <= 0 {
		depth = -1
	}
	dfs(pkg, depth)
	return affected
}

func GetFilePackages(files []string) ([]string, error) {
	uniqueDirs := utils.UniqueDirectories(files)
	return GetPackageNames(uniqueDirs), nil
}

// Function to check if a package contains any test functions
var hasTests = func(pkgName string) (bool, error) {
	cmd := exec.Command("go", "test", pkgName, "-run=^$", "-list", ".")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false, err
	}

	// Parse the output to find test names, excluding the "no test files" message
	output := strings.TrimSpace(out.String())
	return output != "" && !strings.Contains(output, "no test files"), nil
}

// Filter out test packages with no actual test functions
func FilterPackagesWithTests(pkgs []string) []string {
	testPkgs := []string{}
	for _, pkg := range pkgs {
		hasT, err := hasTests(pkg)
		if err != nil {
			log.Error().Err(err).Str("package", pkg).Msg("Error checking for tests")
		}
		if hasT {
			testPkgs = append(testPkgs, pkg)
		}
	}
	return testPkgs
}
