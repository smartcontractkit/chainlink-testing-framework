package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/smartcontractkit/chainlink-testing-framework/tools/gotidy/git"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/clihelper"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

var gomodRegex = regexp.MustCompile("^go.mod$")

type GoProject struct {
	BeforeMod string
	BeforeSum string
	AfterMod  string
	AfterSum  string
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		stop() // restore default exit behavior
		log.Println("Cancelling... interrupt again to exit")
	}()
	project := flag.String("path", ".", "Path to the go project to check for tidy.")
	commit := flag.Bool("commit", false, "Commit the changes if there are any.")
	subProjects := flag.Bool("subprojects", false, "Check subprojects for tidy.")
	verbose := flag.Bool("v", false, "Print verbose output.")
	flag.Parse()
	Main(*project, *commit, *subProjects, *verbose)
}

// Main is the entrypoint for the gotidy tool
func Main(project string, commit, subProjects, verbose bool) {
	// clean up the project path
	projectPath, err := filepath.Abs(filepath.Dir(project))
	if err != nil {
		log.Fatal(err)
	}

	projectsToCheck := []string{projectPath}
	if subProjects {
		projectsToCheck, err = osutil.FindDirectoriesContainingFile(projectPath, gomodRegex)
		if err != nil {
			log.Fatal(err)
		}
	}

	// stash changes if the user wants to commit tidy fixes
	if commit {
		err = git.StashChanges()
		if err != nil {
			ErrorString("Error stashing changes\n")
			log.Fatal(err)
		}
	}

	foundChanges, err := TidyProjects(projectsToCheck, verbose)
	CleanOnError(commit, err)

	if foundChanges {
		if commit {
			err = CommitChanges(project)
			CleanOnError(commit, err)
		} else {
			ErrorString(fmt.Sprintf("Some projects were not tidy. Please run `gotidy -path=%s -subprojects=%t -commit=true` to commit the changes.\n", project, subProjects))
			os.Exit(1)
		}
	}

	if commit {
		err = git.PopStash()
		if err != nil {
			ErrorString("Error un-stashing changes\n")
			log.Fatal(err)
		}
	}
}

// TidyProjects runs `go mod tidy` on a list of projects and compares the changes
func TidyProjects(projects []string, verbose bool) (bool, error) {
	var err error
	foundChanges := false
	for _, projectPath := range projects {
		mod := GoProject{}

		// read the go.mod and go.sum files before tidying
		mod.BeforeMod, mod.BeforeSum, err = ReadModFiles(projectPath)
		if err != nil {
			return true, err
		}

		// change to the project directory
		err = os.Chdir(projectPath)
		if err != nil {
			return true, err
		}

		// run go mod tidy
		fmt.Println("Running go mod tidy on project: ", projectPath)
		err = GoModTidy()
		if err != nil {
			return true, err
		}

		mod.AfterMod, mod.AfterSum, err = ReadModFiles(projectPath)
		if err != nil {
			return true, err
		}

		modDiff := CompareFiles(mod.BeforeMod, mod.AfterMod)
		sumDiff := CompareFiles(mod.BeforeSum, mod.AfterSum)

		if modDiff != "" {
			ErrorString(fmt.Sprintf("Found changes in %s/go.mod\n", projectPath))
			if verbose {
				fmt.Println(modDiff)
			}
			foundChanges = true
		} else if sumDiff != "" {
			ErrorString(fmt.Sprintf("Found changes in %s/go.sum\n", projectPath))
			if verbose {
				fmt.Println(modDiff)
			}
			foundChanges = true
		}
	}
	return foundChanges, nil
}

// ErrorString prints a message in red
func ErrorString(message string) {
	fmt.Print(clihelper.Color(clihelper.ColorRed, message))
}

// ReadModFiles reads the go.mod and go.sum files in a directory
func ReadModFiles(dir string) (gomod string, gosum string, err error) {
	var gomodB, gomodS []byte
	gomodB, err = os.ReadFile(fmt.Sprintf("%s/go.mod", dir))
	if err != nil {
		return
	}
	gomodS, err = os.ReadFile(fmt.Sprintf("%s/go.sum", dir))
	if err != nil {
		return
	}
	gomod = string(gomodB)
	gosum = string(gomodS)
	return
}

// CompareFiles compares two strings and returns a human-readable diff, or empty string if there are no changes
func CompareFiles(before, after string) string {
	// Create a diff object
	dmp := diffmatchpatch.New()

	// Find differences
	rawDiffs := dmp.DiffMain(before, after, true)

	// Check if there are any changes
	hasChanges := false
	for _, diff := range rawDiffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			hasChanges = true
			break
		}
	}

	if hasChanges {
		// Process the diff to make it human-readable
		return dmp.DiffPrettyText(rawDiffs)
	}

	return ""
}

// CommitChanges adds the go.mod and go.sum files to the git index and commits them
func CommitChanges(dir string) (err error) {
	err = git.AddFile("**/go.mod")
	if err != nil {
		return
	}
	err = git.AddFile("**/go.sum")
	if err != nil {
		return
	}

	err = git.CommitChanges("go_mod_tidy_cleanup")
	return
}

// GoModTidy runs `go mod tidy` in the current directory
func GoModTidy() (err error) {
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "go mod tidy", func(m string) {
		fmt.Println(m)
	})

	return
}

// CleanOnError pops the stash if there was an error to return the branch to the original state
func CleanOnError(commit bool, err error) {
	if err != nil && commit {
		pserr := git.PopStash()
		if pserr != nil {
			ErrorString("Error popping stash on cleanup\n")
			log.Fatal(pserr)
		}
		log.Fatal(err)
	}
}
