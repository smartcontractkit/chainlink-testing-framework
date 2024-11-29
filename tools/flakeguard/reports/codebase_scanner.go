package reports

import (
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

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
		if err != nil {
			return err
		}

		// Traverse the AST to find test functions
		ast.Inspect(node, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if strings.HasPrefix(funcDecl.Name.Name, "Test") {
				// Add both the package and test function to the map
				testFileMap[funcDecl.Name.Name] = path
			}
			return true
		})
		return nil
	})

	return testFileMap, err
}
