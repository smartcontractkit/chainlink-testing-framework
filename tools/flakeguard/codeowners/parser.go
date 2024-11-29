package codeowners

import (
	"bufio"
	"os"
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
				Pattern: fields[0],
				Owners:  fields[1:],
			})
		}
	}
	return patterns, scanner.Err()
}

// FindOwners determines the owners for a given file path based on patterns
func FindOwners(filePath string, patterns []PatternOwner) []string {
	// Convert filePath to Unix-style for matching
	relFilePath := filepath.ToSlash(filePath)

	var matchedOwners []string
	for _, pattern := range patterns {
		// Ensure the pattern is also converted to Unix-style
		patternPath := strings.TrimPrefix(pattern.Pattern, "/")
		patternPath = strings.TrimSuffix(patternPath, "/")

		// Match if the file is in the directory or is an exact match
		if strings.HasPrefix(relFilePath, patternPath) || relFilePath == patternPath {
			matchedOwners = pattern.Owners
		}
	}
	return matchedOwners
}
