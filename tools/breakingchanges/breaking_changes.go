package breakingchanges

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	green   = "\033[0;32m"
	yellow  = "\033[0;33m"
	noColor = "\033[0m"
)

func DetectBreakingChanges(rootPath string) {
	// Check if gorelease is installed
	if _, err := exec.LookPath("gorelease"); err != nil {
		log.Fatalf("%sgorelease could not be found. Please install it with 'go install golang.org/x/exp/cmd/gorelease@latest'.%s\n", green, noColor)
	}

	var eg errgroup.Group

	// Function to process each directory
	processDirectory := func(path string) error {
		return runGorelease(path)
	}

	// Check root directory for go.mod
	if _, err := os.Stat(filepath.Join(rootPath, "go.mod")); err == nil {
		eg.Go(func() error {
			return processDirectory(rootPath)
		})
	}

	// Walk through directories starting from rootPath
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory or if it's the root directory (already processed)
		if !info.IsDir() || path == rootPath {
			return nil
		}

		// Check if go.mod exists in the directory
		goModPath := filepath.Join(path, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			eg.Go(func() error {
				return processDirectory(path)
			})
		}

		return nil
	})

	if err != nil {
		log.Fatalf("%sError walking through directories: %v%s\n", green, err, noColor)
	}

	// Wait for all goroutines to finish and collect errors
	if err := eg.Wait(); err != nil {
		log.Fatalf("%sErrors occurred while running gorelease: %v%s\n", green, err, noColor)
	}

	fmt.Printf("%sAll checks completed successfully.%s\n", green, noColor)
}

func runGorelease(path string) error {
	packageFolder := filepath.Base(path)
	var output bytes.Buffer

	// Find the second latest tag for the package
	cmd := exec.Command("git", "tag", "--sort=-creatordate")
	cmd.Dir = path
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sFailed to retrieve git tags for package %s: %v%s\n", green, packageFolder, err, noColor)
		return err
	}

	tags := strings.Split(output.String(), "\n")
	var previousTag string
	for _, tag := range tags {
		if strings.HasPrefix(tag, fmt.Sprintf("%s/v", packageFolder)) {
			if previousTag != "" {
				break
			}
			previousTag = tag
		}
	}

	if previousTag == "" {
		fmt.Printf("%sNo previous tag found for package %s. Skipping.%s\n", green, packageFolder, noColor)
		return nil
	}

	versionTag := strings.Split(previousTag, "/")[1]
	fmt.Printf("%sRunning gorelease for package %s with base tag %s%s\n", green, packageFolder, versionTag, noColor)

	// Run gorelease
	cmd = exec.Command("gorelease", "-base", versionTag)
	cmd.Dir = path
	output.Reset()
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()

	if err != nil {
		fmt.Printf("%sgorelease command failed for package %s with error: %v%s\n%sOutput:%s\n%s%s%s", green, packageFolder, err, noColor, yellow, noColor, yellow, output.String(), noColor)
		return err
	}

	if strings.Contains(output.String(), "Breaking changes") {
		fmt.Printf("%sBreaking changes found for package %s:%s\n%s%s%s", green, packageFolder, noColor, yellow, output.String(), noColor)
	} else {
		fmt.Printf("%sNo breaking changes found for package %s.%s\n", green, packageFolder, noColor)
	}

	return nil
}
