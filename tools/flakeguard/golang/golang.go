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
	"strconv"
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
						_ = os.WriteFile(filepath.Join("go_list_output.json"), buffer.Bytes(), 0644)
						return fmt.Errorf("error unmarshalling go list output for dir '%s', see 'go_list_output.json' for output: %w", cmd.Dir, err)
					}
					packages = append(packages, pkg)
					buffer.Reset()
				}
			}

			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error scanning go list output for file %s: %w", path, err)
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
	// These values might be useless info, but for now we'll keep them
	NewlySkipped   bool  // If the test was newly skipped after calling SkipTests
	AlreadySkipped bool  // If the test was already skipped before calling SkipTests
	ErrorSkipping  error // If we failed to skip the test, this will be set
	File           string
	Line           int
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
			found             = false
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
				return fmt.Errorf("error parsing file '%s': %w", path, err)
			}

			found = skipTest(path, fileAst, fset, testToSkip)

			if found {
				log.Debug().
					Str("test", testToSkip.Name).
					Str("file", testToSkip.File).
					Int("line", testToSkip.Line).
					Str("package", testToSkip.Package).
					Msg("Skipped test")

				// Write back the file
				var out strings.Builder
				if err := printer.Fprint(&out, fset, fileAst); err != nil {
					return fmt.Errorf("error printing file '%s': %w", path, err)
				}
				if err := os.WriteFile(path, []byte(out.String()), 0644); err != nil {
					return fmt.Errorf("error writing file '%s': %w", path, err)
				}
				return filepath.SkipDir // stop walking
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error skipping test '%s' in package '%s': %w", testToSkip.Name, testToSkip.Package, err)
		}
		if !found {
			if testToSkip.ErrorSkipping == nil {
				testToSkip.ErrorSkipping = fmt.Errorf("test '%s' not found in package '%s'", testToSkip.Name, testToSkip.Package)
			} else {
				testToSkip.ErrorSkipping = fmt.Errorf("error skipping test '%s' in package '%s': %w", testToSkip.Name, testToSkip.Package, testToSkip.ErrorSkipping)
			}
			log.Warn().
				Str("test", testToSkip.Name).
				Str("package", testToSkip.Package).
				Err(testToSkip.ErrorSkipping).
				Msg("Unable to skip test")
		}
	}
	return nil
}

// skipTest parses through the file AST and skips the test or subtest if it is found
func skipTest(file string, fileAst *ast.File, fset *token.FileSet, testToSkip *SkipTest) (found bool) {
	parentTest := testToSkip.Name
	subTest := ""
	if strings.Contains(testToSkip.Name, "/") {
		subTest = filepath.Base(testToSkip.Name)
		parentTest = filepath.Dir(testToSkip.Name)
	}

	ast.Inspect(fileAst, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Name.Name != parentTest {
			return true
		}
		// Check if first parameter is *testing.T
		if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
			param := fn.Type.Params.List[0]
			if starExpr, ok := param.Type.(*ast.StarExpr); ok {
				if sel, ok := starExpr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "T" {
					if subTest == "" {
						// No subtest: only skip the parent test
						for _, stmt := range fn.Body.List {
							if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
								if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
									if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
										if selExpr.Sel.Name == "Skip" {
											// Test is already skipped, don't modify it
											found = true
											testToSkip.AlreadySkipped = true
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
											"\"Skipped by flakeguard: https://%s/issues/%s\"",
											os.Getenv("JIRA_DOMAIN"),
											testToSkip.JiraTicket,
										),
									},
								},
							},
						}
						testToSkip.File = file
						testToSkip.Line = fset.Position(fn.Pos()).Line
						testToSkip.NewlySkipped = true
						fn.Body.List = append([]ast.Stmt{skipStmt}, fn.Body.List...)
						found = true
						return false
					} else {
						// Subtest: look for t.Run(subTest, ...)
						ast.Inspect(fn.Body, func(n ast.Node) bool {
							call, ok := n.(*ast.CallExpr)
							if !ok {
								return true
							}
							sel, ok := call.Fun.(*ast.SelectorExpr)
							if !ok || sel.Sel.Name != "Run" {
								return true
							}
							// Check if first argument is a string literal matching subTest
							if len(call.Args) < 2 {
								return true
							}
							if nameLit, ok := call.Args[0].(*ast.BasicLit); ok && nameLit.Kind == token.STRING {
								name, _ := strconv.Unquote(nameLit.Value)
								name = strings.ReplaceAll(name, " ", "_")
								if name == subTest {
									// Second argument should be a function literal
									if fnLit, ok := call.Args[1].(*ast.FuncLit); ok {
										// Check if already skipped
										if len(fnLit.Body.List) > 0 {
											if exprStmt, ok := fnLit.Body.List[0].(*ast.ExprStmt); ok {
												if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
													if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
														if selExpr.Sel.Name == "Skip" {
															found = true
															testToSkip.AlreadySkipped = true
															return false
														}
													}
												}
											}
										}
										// Not skipped, insert skip
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
															"\"Skipped by flakeguard: https://%s/issues/%s\"",
															os.Getenv("JIRA_DOMAIN"),
															testToSkip.JiraTicket,
														),
													},
												},
											},
										}
										fnLit.Body.List = append([]ast.Stmt{skipStmt}, fnLit.Body.List...)
										testToSkip.File = file
										testToSkip.Line = fset.Position(fnLit.Pos()).Line
										testToSkip.NewlySkipped = true
										found = true
										return false
									}
								}
							}
							return true
						})
					}
				}
			}
		}
		return true
	})

	return found
}
