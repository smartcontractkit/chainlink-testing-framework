package reports

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// TestFileMap maps test function names to their file paths.
type TestFileMap map[string]string

// ScanTestFiles scans the codebase for test functions and maps them to file paths.
func ScanTestFiles(rootDir string) (TestFileMap, error) {
	testFileMap := make(TestFileMap)

	// Ensure rootDir is absolute
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("error normalizing rootDir: %v", err)
	}

	// Walk through the root directory to find test files
	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip files that are not Go test files
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Normalize path relative to rootDir
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %v", path, err)
		}
		relPath = filepath.ToSlash(relPath) // Ensure Unix-style paths

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", path, err)
		}

		// Traverse the AST to find test or fuzz functions
		ast.Inspect(node, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			// Match both "Test" and "Fuzz" prefixes
			if strings.HasPrefix(funcDecl.Name.Name, "Test") || strings.HasPrefix(funcDecl.Name.Name, "Fuzz") {
				// Add the function to the map with relative path
				testFileMap[funcDecl.Name.Name] = relPath
			}
			return true
		})

		return nil
	})

	return testFileMap, err
}
