package codeowners

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// PatternOwner maps a file pattern to its owners
type PatternOwner struct {
	Pattern string
	Owners  []string
}

// Parse reads the CODEOWNERS file and returns a list of PatternOwner
func Parse(filePath string) ([]PatternOwner, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []PatternOwner
	scanner := bufio.NewScanner(file)
	commentPattern := regexp.MustCompile(`^\s*#`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || commentPattern.MatchString(line) {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 1 {
			patterns = append(patterns, PatternOwner{
				Pattern: filepath.ToSlash(fields[0]), // Normalize to Unix-style
				Owners:  fields[1:],
			})
		}
	}
	return patterns, scanner.Err()
}

func IsWildcardPattern(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

// FindOwners finds the owners of a file based on the CODEOWNERS patterns
func FindOwners(filePath string, patterns []PatternOwner) []string {
	// Normalize the file path to Unix-style
	relFilePath := filepath.ToSlash(filePath)

	var matchedOwners []string
	for _, pattern := range patterns {
		// Normalize the pattern to Unix-style and remove leading and trailing slashes
		normalizedPattern := filepath.ToSlash(strings.TrimPrefix(pattern.Pattern, "/"))
		normalizedPattern = strings.TrimSuffix(normalizedPattern, "/")

		if IsWildcardPattern(normalizedPattern) {
			matched, err := path.Match(normalizedPattern, relFilePath)
			if err != nil {
				fmt.Printf("Error matching pattern: %s to file: %s, error: %v\n", normalizedPattern, relFilePath, err)
				continue
			}

			if matched {
				matchedOwners = pattern.Owners
			}
		} else {
			if relFilePath == normalizedPattern {
				// Exact file or directory match
				matchedOwners = pattern.Owners
			} else if strings.HasPrefix(relFilePath, normalizedPattern+"/") {
				// File is under the directory pattern
				matchedOwners = pattern.Owners
			}
		}
	}

	return matchedOwners
}
