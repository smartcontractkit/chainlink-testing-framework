package golang

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
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

// List returns the output of `go list -json ./...`
// Deprecated: Use Packages instead
func List() (*utils.CmdOutput, error) {
	return utils.ExecuteCmd("go", "list", "-json", "./...")
}

// Packages finds all packages in the repository
func Packages(repoPath string) ([]Package, error) {
	var packages []Package

	// Find all go.mod files and run go list -json ./... in the directory of each go.mod file
	// This is necessary because go list -json ./... only returns packages that are associated with the current go.mod file
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "go.mod" {
			cmd := exec.Command("go", "list", "-json", "./...")
			cmd.Dir = filepath.Dir(path)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("error getting packages: %w\nOutput:\n%s", err, string(out))
			}
			scanner := bufio.NewScanner(bytes.NewReader(out))
			var buffer bytes.Buffer

			for scanner.Scan() {
				line := scanner.Text()
				line = strings.TrimSpace(line)
				buffer.WriteString(line)

				if line == "}" {
					var pkg Package
					if err := json.Unmarshal(buffer.Bytes(), &pkg); err != nil {
						return err
					}
					packages = append(packages, pkg)
					buffer.Reset()
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
		}
		return nil
	})

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

// SkipTest is a struct that contains the package and name of the test to skip
type SkipTest struct {
	Package    string
	Name       string
	JiraTicket string

	// Set by SkipTests

	// If the test was already skipped, or there were issues skipping it, this will be false
	Skipped bool

	// The file and line number of the test that was skipped
	File string
	Line int
}

// SkipTests finds all package/test pairs provided and skips them
func SkipTests(repoPath string, testsToSkip []*SkipTest) error {
	packages, err := Packages(repoPath)
	if err != nil {
		return err
	}

	for _, testToSkip := range testsToSkip {
		var (
			packageDir        string
			packageImportPath = testToSkip.Package
		)
		for _, pkg := range packages {
			if pkg.ImportPath == packageImportPath {
				packageDir = pkg.Dir
				break
			}
		}
		if packageDir == "" {
			return fmt.Errorf("directory for package '%s' not found", packageImportPath)
		}

		log.Debug().
			Str("package_dir", packageDir).
			Str("test", testToSkip.Name).
			Str("package", packageImportPath).
			Msg("Skipping test")

		err := filepath.Walk(packageDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), "_test.go") {
				return nil
			}

			fset := token.NewFileSet()
			fileAst, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			found := false
			ast.Inspect(fileAst, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok || fn.Name.Name != testToSkip.Name {
					return true
				}
				// Check if first parameter is *testing.T
				if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
					param := fn.Type.Params.List[0]
					if starExpr, ok := param.Type.(*ast.StarExpr); ok {
						if sel, ok := starExpr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "T" {
							// Check if test is already skipped
							for _, stmt := range fn.Body.List {
								if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
									if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
										if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
											if selExpr.Sel.Name == "Skip" {
												// Test is already skipped, don't modify it
												found = true
												return false
											}
										}
									}
								}
							}

							// Test is not skipped, add skip statement
							skipStmt := &ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent(param.Names[0].Name),
										Sel: ast.NewIdent("Skip"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind: token.STRING,
											Value: fmt.Sprintf(
												"Skipped by flakeguard: https://%s/issues/%s",
												os.Getenv("JIRA_DOMAIN"),
												testToSkip.JiraTicket,
											),
										},
									},
								},
							}
							testToSkip.File = path
							testToSkip.Line = fset.Position(fn.Pos()).Line
							testToSkip.Skipped = true
							fn.Body.List = append([]ast.Stmt{skipStmt}, fn.Body.List...)
							found = true
							return false
						}
					}
				}
				return true
			})

			if found {
				// Write back the file
				var out strings.Builder
				if err := printer.Fprint(&out, fset, fileAst); err != nil {
					return err
				}
				if err := os.WriteFile(path, []byte(out.String()), info.Mode()); err != nil {
					return err
				}
				return filepath.SkipDir // stop walking
			}
			return fmt.Errorf("test function %s not found in package %s", testToSkip.Name, testToSkip.Package)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
