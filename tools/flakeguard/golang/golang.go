package golang

import (
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
	Deps         []string `json:"Deps,omitempty"`
	TestImports  []string `json:"TestImports,omitempty"`
	XTestImports []string `json:"XTestImports,omitempty"`
	GoFiles      []string `json:"GoFiles,omitempty"`
	TestGoFiles  []string `json:"TestGoFiles,omitempty"`
	XTestGoFiles []string `json:"XTestGoFiles,omitempty"`
	EmbedFiles   []string `json:"EmbedFiles,omitempty"`
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
			// output sometimes contains go downloading things, need to strip that out
			jsonStart := bytes.IndexByte(out, '{')
			if jsonStart == -1 {
				return fmt.Errorf("no JSON found in go list output for dir '%s': go list output:\n%s", cmd.Dir, string(out))
			}
			cleanOutput := out[jsonStart:]

			// go list -json outputs multiple JSON objects (one per package), not a JSON array
			// So we need to decode each JSON object individually
			var pkgs []Package
			decoder := json.NewDecoder(bytes.NewReader(cleanOutput))
			for decoder.More() {
				var pkg Package
				if err := decoder.Decode(&pkg); err != nil {
					return fmt.Errorf("error unmarshalling go list output for dir '%s': %w\ngo list output:\n%s", cmd.Dir, err, string(cleanOutput))
				}
				pkgs = append(pkgs, pkg)
			}
			packages = append(packages, pkgs...)
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

// SkipTest contains info about tests to skip, and the results of skipping them
type SkipTest struct {
	Package    string
	Name       string
	JiraTicket string

	// Set by SkipTests
	// These values might be useless info, but for now we'll keep them
	SimplySkipped  bool   // If the test was newly skipped using simple AST parsing methods
	AlreadySkipped bool   // If the test was already skipped before calling SkipTests
	ErrorSkipping  error  // If we failed to skip the test, this will be set
	File           string // The file that the test was skipped in
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
			testToSkip.ErrorSkipping = fmt.Errorf(
				"unable to find a directory for package '%s' to skip '%s', package may have moved or been deleted",
				packageImportPath, testToSkip.Name,
			)
			return nil
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

			found = skipTestSimple(path, fileAst, fset, testToSkip)
			if found {
				log.Debug().
					Str("test", testToSkip.Name).
					Str("file", testToSkip.File).
					Int("line", testToSkip.Line).
					Str("package", testToSkip.Package).
					Msg("Skipped test using simple AST parsing")

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
				testToSkip.ErrorSkipping = fmt.Errorf(
					"cannot find '%s' in package '%s', the test may be tricky for me to find, or may have moved or been deleted",
					testToSkip.Name, testToSkip.Package,
				)
			} else {
				testToSkip.ErrorSkipping = fmt.Errorf(
					"error skipping '%s' in package '%s': %w",
					testToSkip.Name, testToSkip.Package, testToSkip.ErrorSkipping,
				)
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

// skipTestSimple parses through the file AST and skips the test or subtest if it is found
// This works for most simple cases, is fast and efficient, but can't handle more complex cases like when subtests or helper functions get involved
func skipTestSimple(file string, fileAst *ast.File, fset *token.FileSet, testToSkip *SkipTest) (skipped bool) {
	var (
		parentTest = testToSkip.Name
		subTest    string
		jiraMsg    string
	)
	if strings.Contains(testToSkip.Name, "/") {
		subTest = filepath.Base(testToSkip.Name)
		parentTest = filepath.Dir(testToSkip.Name)
	}

	if os.Getenv("JIRA_DOMAIN") != "" && testToSkip.JiraTicket != "" {
		jiraMsg = fmt.Sprintf("Skipped by flakeguard: https://%s/issues/%s", os.Getenv("JIRA_DOMAIN"), testToSkip.JiraTicket)
	} else {
		jiraMsg = fmt.Sprintf("Skipped by flakeguard: %s", testToSkip.JiraTicket)
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
											skipped = true
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
											"\"%s\"",
											jiraMsg,
										),
									},
								},
							},
						}
						testToSkip.File = file
						testToSkip.Line = fset.Position(fn.Pos()).Line
						testToSkip.SimplySkipped = true
						fn.Body.List = append([]ast.Stmt{skipStmt}, fn.Body.List...)
						skipped = true
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
															skipped = true
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
															"\"%s\"",
															jiraMsg,
														),
													},
												},
											},
										}
										fnLit.Body.List = append([]ast.Stmt{skipStmt}, fnLit.Body.List...)
										testToSkip.File = file
										testToSkip.Line = fset.Position(fnLit.Pos()).Line
										testToSkip.SimplySkipped = true
										skipped = true
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

	return skipped
}

/* WARNING: This was a fun experiment, but it gave massively disappointing results.
 * Might be worth revisiting in the future, but this juice ain't worth the squeeze.

// https://openai.com/api/pricing/
const (
	gpt4_1NanoInputCostPer1MTokens  = 0.100
	gpt4_1NanoOutputCostPer1MTokens = 0.400

	gpt4_1MiniInputCostPer1MTokens  = 0.400
	gpt4_1MiniOutputCostPer1MTokens = 1.600

	gpt4_1InputCostPer1MTokens  = 2.000
	gpt4_1OutputCostPer1MTokens = 8.000
)

var (
	totalOpenAICost float64
)

// skipTestLLM uses an LLM to skip a test or subtest
// It's slower and more expensive, but can handle more complex cases than skipTestSimple
// Use this when skipTestSimple fails to skip a test or subtest

func skipTestLLM(file, openAIKey string, testToSkip *SkipTest) (skipped bool) {
	if openAIKey == "" {
		return false
	}

	fileContent, err := os.ReadFile(file)
	if err != nil {
		testToSkip.ErrorSkipping = fmt.Errorf("failed to read file: %w", err)
		return false
	}

	jiraDomain := os.Getenv("JIRA_DOMAIN")
	jiraMsg := ""
	if jiraDomain != "" && testToSkip.JiraTicket != "" {
		jiraMsg = fmt.Sprintf("Skipped by flakeguard LLM: https://%s/issues/%s", jiraDomain, testToSkip.JiraTicket)
	} else {
		jiraMsg = fmt.Sprintf("Skipped by flakeguard LLM: %s", testToSkip.JiraTicket)
	}

	prompt := fmt.Sprintf(`You are an expert Go developer.
Given the following Go test file, find the exact spot where a t.Skip call would skip this test: %s.
If the test provided is a subtest, we only want to skip that specific subtest.
If you are able to find the spot where the t.Skip call should be inserted, return ONLY the line number, and the necessary code to insert the t.Skip call.
Do so in the following format:

<line number>: <code to insert>

Example:

10: if tt.Name() == "%s" { tt.Skip() }

If you are unable to find the spot where the t.Skip call should be inserted, return ONLY "Unable to skip test: <short reason why>".
If the test is already being skipped, return ONLY "Test is already skipped".
Here is the code:\n\n%s`, testToSkip.Name, testToSkip.Name, string(fileContent))

	client := openai.NewClient(
		option.WithAPIKey(openAIKey),
	)

	completionParams := openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4_1Nano,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	}

	// Estimate costs for OpenAI API calls
	var (
		cost             float64
		inputMultiplier  float64
		outputMultiplier float64
	)
	defer func() {
		totalOpenAICost += cost
	}()
	switch openai.ChatModel(completionParams.Model) {
	case openai.ChatModelGPT4_1Nano:
		inputMultiplier = gpt4_1NanoInputCostPer1MTokens
		outputMultiplier = gpt4_1NanoOutputCostPer1MTokens
	case openai.ChatModelGPT4_1Mini:
		inputMultiplier = gpt4_1MiniInputCostPer1MTokens
		outputMultiplier = gpt4_1MiniOutputCostPer1MTokens
	case openai.ChatModelGPT4_1:
		inputMultiplier = gpt4_1InputCostPer1MTokens
		outputMultiplier = gpt4_1OutputCostPer1MTokens
	default:
		log.Warn().
			Str("model", string(completionParams.Model)).
			Msg("Unknown OpenAI model, unable to accurately calculate cost")
		inputMultiplier = gpt4_1InputCostPer1MTokens
		outputMultiplier = gpt4_1OutputCostPer1MTokens
	}

	start := time.Now()

	resp, err := client.Chat.Completions.New(context.Background(), completionParams)

	duration := time.Since(start)
	cost += inputMultiplier * float64(len(prompt)) / 1_000_000
	log.Trace().
		Str("model", string(completionParams.Model)).
		Str("duration", duration.String()).
		Str("cost", fmt.Sprintf("$%.2f", cost)).
		Msg("OpenAI API call completed")
	if err != nil {
		testToSkip.ErrorSkipping = fmt.Errorf("OpenAI API error: %w", err)
		return false
	}

	if len(resp.Choices) == 0 {
		testToSkip.ErrorSkipping = fmt.Errorf("no response from OpenAI")
		return false
	}
	cost += outputMultiplier * float64(len(resp.Choices[0].Message.Content)) / 1_000_000

	response := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.Contains(response, "Test is already skipped") {
		testToSkip.AlreadySkipped = true
		return true
	}

	if strings.Contains(response, "Unable to skip test") {
		testToSkip.ErrorSkipping = fmt.Errorf(
			"OpenAI could not find a way to skip the test: %s",
			strings.TrimPrefix(response, "Unable to skip test: "),
		)
		return false
	}

	// We should only get back a line number and line of code, so let's parse it and add the skip statement
	// This is cheaper, faster, and more accurate than having the LLM parrot the whole file back to us
	// Clean up the response a bit
	response = strings.TrimPrefix(response, "Line")
	response = strings.TrimPrefix(response, ":")
	response = strings.TrimSpace(response)

	// Split the response into line number and code
	parts := strings.SplitN(response, ":", 2)
	if len(parts) != 2 {
		testToSkip.ErrorSkipping = fmt.Errorf("OpenAI returned invalid response: %s", response)
		return false
	}
	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])

	lineNum, err := strconv.Atoi(parts[0])
	if err != nil {
		testToSkip.ErrorSkipping = fmt.Errorf("OpenAI returned invalid line number '%s': %w", response, err)
		return false
	}
	codeToInsert := parts[1]

	// Add our jira message inside the t.Skip call
	codeToInsert = strings.ReplaceAll(codeToInsert, "t.Skip()", fmt.Sprintf("t.Skip(\"%s\")", jiraMsg))

	// Read the file and split into lines
	lines := strings.Split(string(fileContent), "\n")

	// Validate line number (1-based indexing from OpenAI, convert to 0-based)
	if lineNum < 1 || lineNum > len(lines) {
		testToSkip.ErrorSkipping = fmt.Errorf("line number %d suggested by OpenAI is out of range (file has %d lines)", lineNum, len(lines))
		return false
	}

	// Get the indentation of the line we're inserting before
	targetLineIndex := lineNum - 1 // Convert to 0-based index
	var indentation string
	if targetLineIndex < len(lines) {
		// Extract indentation from the target line
		line := lines[targetLineIndex]
		for _, char := range line {
			if char == ' ' || char == '\t' {
				indentation += string(char)
			} else {
				break
			}
		}
	}

	// Create the skip statement with proper indentation
	skipStatement := fmt.Sprintf("%s%s", indentation, codeToInsert)

	// Insert the skip statement at the specified line
	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:targetLineIndex]...)
	newLines = append(newLines, skipStatement)
	newLines = append(newLines, lines[targetLineIndex:]...)

	// Write the modified content back to the file
	modifiedContent := strings.Join(newLines, "\n")
	err = os.WriteFile(file, []byte(modifiedContent), 0644)
	if err != nil {
		testToSkip.ErrorSkipping = fmt.Errorf("failed to write modified file: %w", err)
		return false
	}

	testToSkip.File = file
	testToSkip.Line = lineNum
	testToSkip.LLMSkipped = true
	return true
}
*/
