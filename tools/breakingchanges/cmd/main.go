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
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/mod/modfile"
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

func getRetractedTags(goModPath string) ([]*semver.Constraints, error) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, fmt.Errorf("error reading go.mod file: %w", err)
	}

	modFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing go.mod file: %w", err)
	}

	var retractedTags []*semver.Constraints
	for _, retract := range modFile.Retract {
		lowVersion, err := semver.NewVersion(retract.Low)
		if err != nil {
			return nil, fmt.Errorf("error parsing retracted version: %w", err)
		}
		highVersion, err := semver.NewVersion(retract.High)
		if err != nil {
			return nil, fmt.Errorf("error parsing retracted version: %w", err)
		}
		constraint, err := semver.NewConstraint(fmt.Sprintf(">= %s, <= %s", lowVersion.String(), highVersion.String()))
		if err != nil {
			return nil, fmt.Errorf("error parsing retracted version: %w", err)
		}
		retractedTags = append(retractedTags, constraint)
		fmt.Printf("Retracted version: %s\n", constraint)
	}

	return retractedTags, nil
}

func getLatestTag(pathPrefix string, retractedTags []*semver.Constraints) (string, error) {
	// use regex to find exact matches, as otherwise might include pre-release versions
	// or versions that partially match the path prefix, e.g. when searching for 'lib'
	// we want to make sure we won't include tags like `lib/grafana/v1.0.0`
	grepRegex := fmt.Sprintf("^%s/v[0-9]+\\.[0-9]+\\.[0-9]+$", pathPrefix)

	//nolint
	cmd := exec.Command("sh", "-c", fmt.Sprintf("git tag | grep -E '%s'", grepRegex))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error fetching tags: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found for regex: %s", grepRegex)
	}

	// Parse the tags into semver versions
	var allTags []*semver.Version
	for _, tag := range tags {
		v, err := semver.NewVersion(strings.TrimPrefix(tag, pathPrefix+"/"))
		if err != nil {
			return "", fmt.Errorf("error parsing version tag: %w", err)
		}
		allTags = append(allTags, v)
	}

	// Sort the tags in descending order
	sort.Sort(sort.Reverse(semver.Collection(allTags)))

	if len(retractedTags) == 0 {
		tag := fmt.Sprintf("v%s", allTags[0].String())
		return tag, nil
	}

	// Find the latest tag that doesn't match any of the retracted tags
	for _, tag := range allTags {
		isRetracted := false
		for _, constraint := range retractedTags {
			if constraint.Check(tag) {
				isRetracted = true
				break
			}
		}

		if !isRetracted {
			tag := fmt.Sprintf("v%s", tag.String())
			fmt.Printf("Found non-retracted tag: %s\n", tag)
			return tag, nil
		}
	}

	fmt.Println("No non-retracted tags found")
	fmt.Printf("All tags: %s\n", strings.Join(tags, ", "))
	fmt.Println("Retracted tags:")
	for _, constraint := range retractedTags {
		fmt.Printf("%s\n", constraint)
	}

	return "", fmt.Errorf("failed to find a non-retracted tag got path prefix: %s", pathPrefix)
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
		ignoredDirs = append(ignoredDirs, allDirs...)
	}
	return ignoredDirs
}

func isIgnoredDirRegex(pathPrefix string, ignoredDirs []string) bool {
	for _, d := range ignoredDirs {
		if d == "" {
			continue
		}

		fmt.Printf("Checking regex: %s, path: %s\n", d, pathPrefix)
		re := regexp.MustCompile(d)
		if re.MatchString(pathPrefix) {
			fmt.Printf("Path is ignored, skipping: %s\n", pathPrefix)
			return true
		}
	}
	return false
}

func main() {
	rootFolder := flag.String("root", ".", "The root folder to start scanning from")
	subDir := flag.String("subdir", "", "The subdirectory inside the root folder to scan for modules")
	ignoreDirs := flag.String("ignore", "", "Ignore directory paths matching regex (comma-separated)")
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

		if isIgnoredDirRegex(pathPrefix, ignoredDirs) {
			continue
		}

		retractedVersions, err := getRetractedTags(filepath.Join(dirPath, "go.mod"))
		if err != nil {
			fmt.Printf("Error getting retracted versions: %v\n", err)
			continue
		}

		latestTag, err := getLatestTag(pathPrefix, retractedVersions)
		if err != nil {
			fmt.Printf("Error finding latest tag: %v\n", err)
			continue
		}

		if latestTag != "" {
			fmt.Printf("%sProcessing directory: %s%s\n", Yellow, dirPath, Reset)
			if err := os.Chdir(dirPath); err != nil {
				fmt.Printf("Error changing directory: %v\n", err)
				continue
			}

			stdout, stderr, err := checkBreakingChanges(latestTag)
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
