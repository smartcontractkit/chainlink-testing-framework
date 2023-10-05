package flake

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TestFlake struct {
	Test    string `json:"test"`
	Package string `json:"package"`
	Error   string `json:"error"`
}

func GetDuplicateTestNames(goProjectPath string) (map[string]int, error) {
	// Set up a map to track test names and their occurrence count
	testNames := map[string]int{}

	// Set up the Go parser
	fset := token.NewFileSet()

	// Walk through all Go files in the directory and its sub-directories
	err := filepath.Walk(goProjectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			node, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}

			// Check each function in the file
			for _, decl := range node.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv == nil && strings.HasPrefix(fn.Name.Name, "Test") && fn.Name.Name != "TestMain" {
					// Increment the count for this test name
					testNames[fn.Name.Name]++
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
		return nil, err
	}

	// Print out test names that appear more than once
	for testName, count := range testNames {
		if count > 1 {
			fmt.Printf("Duplicate test name: %s (Count: %d)\n", testName, count)
		}
	}

	return testNames, nil
}

func ReadFlakyTests(flakyTestFile string) ([]TestFlake, error) {
	// Open the JSON file
	file, err := os.Open(flakyTestFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the content of the file
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a slice of TestError structs
	var flakyTests []TestFlake
	err = json.Unmarshal(bytes, &flakyTests)
	if err != nil {
		return nil, err
	}
	return flakyTests, nil
}

// CompareDuplicateTestNamesToFlakeTestNames checks for the test names in the flaky test file that are duplicates
func CompareDuplicateTestNamesToFlakeTestNames(flakyTests []TestFlake, allTestNames map[string]int) error {
	for _, ft := range flakyTests {
		if _, ok := allTestNames[ft.Test]; !ok {
			return fmt.Errorf("test name %s in flaky test file does not exist in project", ft.Test)
		} else if allTestNames[ft.Test] > 1 {
			return fmt.Errorf("test name %s is a duplicate test name", ft.Test)
		}
	}

	return nil
}
