package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

type Package struct {
	ImportPath  string   `json:"ImportPath"`
	Deps        []string `json:"Deps"`
	TestImports []string `json:"TestImports"`
}

type Config struct {
	Branch      string
	ProjectPath string
	Excludes    []string
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		stop() // restore default exit behavior
		log.Println("Cancelling... interrupt again to exit")
	}()

	branch := flag.String("b", "", "The base git branch to compare current changes with. Required.")
	projectPath := flag.String("p", "", "The path to the project. Default is the current directory. Useful for subprojects.")
	excludes := flag.String("e", "", "The comma separated list of paths to exclude. Useful for repositories with multiple go projects within.")
	flag.Parse()

	config := SetConfig(branch, projectPath, excludes)

	Run(config)
}

func SetConfig(branch, projectPath, excludes *string) *Config {
	if *branch == "" {
		log.Fatalf("Branch is required")
	}

	parsedExcludes := []string{}
	if *excludes != "" {
		parsedExcludes = strings.Split(*excludes, ",")
		for i, e := range parsedExcludes {
			parsedExcludes[i] = strings.TrimSpace(e)
		}
	}
	return &Config{
		Branch:      *branch,
		ProjectPath: *projectPath,
		Excludes:    parsedExcludes,
	}
}

func Run(config *Config) {
	var changedPackages []string
	gitDiff := gitDiff(config.Branch)
	changedPackages, err := getChangedPackages(gitDiff, config.ProjectPath, config.Excludes)
	if err != nil {
		log.Fatalf("Error getting changed packages: %v", err)
	}

	gitModDiff := gitModDiff(config.Branch, config.ProjectPath)
	modChangedPackages, err := getGoModChanges(gitModDiff)
	if err != nil {
		log.Fatalf("Error getting go.mod changes: %v", err)
	}

	goList := goList()
	depGraph := getGoDepMap(goList)

	// Find affected packages
	affectedPkgs := map[string]bool{}
	for _, pkg := range changedPackages {
		p := findAffectedPackages(pkg, depGraph, false)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	for _, pkg := range modChangedPackages {
		p := findAffectedPackages(pkg, depGraph, true)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	o := ""
	for k := range affectedPkgs {
		o = fmt.Sprintf("%s ./%s ", o, k)
	}
	fmt.Println(o)
}

func executeCommand(name string, args ...string) bytes.Buffer {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		joined := strings.Join(args, " ")
		log.Fatalf("Error running command: %s %s\n Err: %v", name, joined, err)
	}
	return out
}

func goList() bytes.Buffer {
	return executeCommand("go", "list", "-json", "./...")
}

// getGoDepMap returns a map of Go packages to their dependencies.
func getGoDepMap(goList bytes.Buffer) map[string][]string {
	var packages []Package
	scanner := bufio.NewScanner(&goList)
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		if line == "}" {
			buffer.WriteString(line)
			var pkg Package
			if err := json.Unmarshal(buffer.Bytes(), &pkg); err != nil {
				panic(err)
			}
			packages = append(packages, pkg)
			buffer.Reset()
		} else {
			buffer.WriteString(line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Build dependency graph
	depGraph := make(map[string][]string)
	for _, pkg := range packages {
		for _, dep := range pkg.Deps {
			depGraph[dep] = append(depGraph[dep], pkg.ImportPath)
		}
	}
	return depGraph
}

// findAffectedPackages takes a package name and a dependency graph and returns the list of packages that are affected by the change.
func findAffectedPackages(pkg string, depGraph map[string][]string, externalPackage bool) []string {
	visited := make(map[string]bool)
	var affected []string

	var dfs func(string)
	dfs = func(p string) {
		if visited[p] {
			return
		}

		visited[p] = true
		// exclude the package itself if it is an external package
		if !(externalPackage && p == pkg) {
			affected = append(affected, p)
		}
		for _, dep := range depGraph[p] {
			dfs(dep)
		}
	}

	dfs(pkg)
	return affected
}

func gitDiff(baseBranch string) bytes.Buffer {
	return executeCommand("git", "diff", "--name-only", baseBranch)
}

// SliceContains checks if a slice contains a given item
func shouldExclude(excludes []string, item string) bool {
	for _, v := range excludes {
		if strings.HasPrefix(item, v) {
			return true
		}
	}
	return false
}

// getChangedPackages takes a base branch and returns the list of Go packages that have changed in the current branch.
func getChangedPackages(out bytes.Buffer, projectPath string, excludes []string) ([]string, error) {
	// Get the list of changed files
	changedFiles := strings.Split(out.String(), "\n")

	// Filter out non-Go files and directories
	changedPackages := make(map[string]struct{})
	for _, file := range changedFiles {
		if strings.HasSuffix(file, ".go") && !shouldExclude(excludes, file) && strings.HasPrefix(file, projectPath) {
			changedPackages[file] = struct{}{}
		}
	}

	// Convert map keys to slice
	var packages []string
	for pkg := range changedPackages {
		packages = append(packages, pkg)
	}

	return packages, nil
}

func gitModDiff(baseBranch, projectPath string) bytes.Buffer {
	return executeCommand("git", "diff", baseBranch, "--unified=0", "--", filepath.Join(projectPath, "go.mod"))
}

// getGoModChanges takes a base branch and returns the list of packages that have changed in the go.mod file.
func getGoModChanges(lines bytes.Buffer) ([]string, error) {
	// Get the list of changed lines
	changedLines := strings.Split(lines.String(), "\n")

	// Filter out lines that do not indicate package changes
	var packages []string
	for _, line := range changedLines {
		if strings.HasPrefix(line, "+") {
			// Ignore lines that are not relevant (e.g., comments or empty lines)
			if strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "+ ") {
				continue
			}
			// Split the line into fields
			fields := strings.Fields(line)
			if len(fields) > 1 {
				// The second field should contains the module path
				packages = append(packages, fields[1])
			}
		}
	}

	return packages, nil
}
