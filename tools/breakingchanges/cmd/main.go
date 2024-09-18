package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Reset  = "\033[0m"
)

func findGoModDirs(rootFolder, subDir string) ([]string, error) {
	var goModDirs []string

	err := filepath.WalkDir(filepath.Join(rootFolder, subDir), func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			goModPath := filepath.Join(path, "go.mod")
			if _, err := os.Stat(goModPath); !os.IsNotExist(err) {
				// Ensure we store absolute paths
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				goModDirs = append(goModDirs, absPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return goModDirs, nil
}

func getLastTag(pathPrefix string) (string, error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("git tag | grep '%s' | tail -1", pathPrefix))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error fetching tags: %w", err)
	}

	tag := strings.TrimSpace(out.String())
	if tag == "" {
		return "", nil
	}

	// Use regex to find the version tag starting with 'v'
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 0 {
		tag = matches[0]
	} else {
		return "", fmt.Errorf("no valid version tag found in '%s'", tag)
	}

	return tag, nil
}

func checkBreakingChanges(tag string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("gorelease", "-base", tag)
	fmt.Printf("%sExecuting command: %s %s\n", Yellow, cmd, Reset)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func getIgnoredDirs(flag *string) []string {
	ignoredDirs := make([]string, 0)
	if flag != nil {
		allDirs := strings.Split(*flag, ",")
		for _, d := range allDirs {
			ignoredDirs = append(ignoredDirs, d)
		}
	}
	return ignoredDirs
}

func isIgnoredDirPrefix(pathPrefix string, ignoredDirs []string) bool {
	for _, d := range ignoredDirs {
		if d == "" {
			continue
		}
		fmt.Printf("Checking prefix: %s, path: %s\n", d, pathPrefix)
		if strings.HasPrefix(pathPrefix, d) {
			fmt.Printf("Path is ignored, skipping: %s\n", pathPrefix)
			return true
		}
	}
	return false
}

func main() {
	rootFolder := flag.String("root", ".", "The root folder to start scanning from")
	subDir := flag.String("subdir", "", "The subdirectory inside the root folder to scan for modules")
	ignoreDirs := flag.String("ignore", "", "Ignore directory paths starting with prefix")
	flag.Parse()

	absRootFolder, err := filepath.Abs(*rootFolder)
	if err != nil {
		fmt.Printf("Error getting absolute path of root folder: %v\n", err)
		return
	}

	goModDirs, err := findGoModDirs(absRootFolder, *subDir)
	if err != nil {
		fmt.Printf("Error finding directories: %v\n", err)
		return
	}

	ignoredDirs := getIgnoredDirs(ignoreDirs)

	breakingChanges := false
	for _, dirPath := range goModDirs {
		// Convert the stripped path back to absolute
		pathPrefix := strings.TrimPrefix(dirPath, absRootFolder+string(os.PathSeparator))

		if isIgnoredDirPrefix(pathPrefix, ignoredDirs) {
			continue
		}

		lastTag, err := getLastTag(pathPrefix)
		if err != nil {
			fmt.Printf("Error finding last tag: %v\n", err)
			continue
		}

		if lastTag != "" {
			fmt.Printf("%sProcessing directory: %s%s\n", Yellow, dirPath, Reset)
			if err := os.Chdir(dirPath); err != nil {
				fmt.Printf("Error changing directory: %v\n", err)
				continue
			}

			stdout, stderr, err := checkBreakingChanges(lastTag)
			if err != nil {
				fmt.Printf("Error running gorelease: %v\n", err)
				breakingChanges = true
			}
			fmt.Printf("%sgorelease output:\n%s%s\n", Green, stdout, Reset)
			if stderr != "" {
				fmt.Printf("%sgorelease errors:\n%s%s\n", Red, stderr, Reset)
			}

			if err := os.Chdir(absRootFolder); err != nil {
				fmt.Printf("Error changing back to root directory: %v\n", err)
			}
		} else {
			fmt.Printf("No valid tags found for path prefix: %s\n", pathPrefix)
		}
	}
	if breakingChanges {
		log.Fatalf("breaking changes detected!")
	}
}
