package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/clihelper"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

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

func Main(project string, commit, subProjects, verbose bool) {
	// clean up the project path
	projectPath, err := filepath.Abs(filepath.Dir(project))
	if err != nil {
		log.Fatal(err)
	}

	projectsToCheck := []string{projectPath}
	if subProjects {
		projectsToCheck, err = FindSubProjects(projectPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// stash changes if the user wants to commit tidy fixes
	if commit {
		err = StashChanges()
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
		err = PopStash()
		if err != nil {
			ErrorString("Error un-stashing changes\n")
			log.Fatal(err)
		}
	}
}

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

func ErrorString(message string) {
	fmt.Print(clihelper.Color(clihelper.ColorRed, message))
}

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

func CommitChanges(dir string) (err error) {
	fmt.Println("Committing changes...")
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, fmt.Sprintf("git add %s/go.mod", dir), func(m string) {
		fmt.Println(m)
	})
	if err != nil {
		return
	}
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, fmt.Sprintf("git add %s/go.sum", dir), func(m string) {
		fmt.Println(m)
	})
	if err != nil {
		return
	}
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "git commit -m \"go_mod_tidy_cleanup\"", func(m string) {
		fmt.Println(m)
	})

	return
}

func GoModTidy() (err error) {
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "go mod tidy", func(m string) {
		fmt.Println(m)
	})

	return
}

func StashChanges() (err error) {
	fmt.Println("Doing a git stash before tidying...")
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "git stash", func(m string) {
		fmt.Println(m)
	})

	return
}

func PopStash() (err error) {
	fmt.Println("Popping the stash to return previous changes...")
	err = osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "git stash pop", func(m string) {
		fmt.Println(m)
	})

	return
}

func CleanOnError(commit bool, err error) {
	if err != nil && commit {
		pserr := PopStash()
		if pserr != nil {
			ErrorString("Error popping stash on cleanup\n")
			log.Fatal(pserr)
		}
		log.Fatal(err)
	}
}

func FindSubProjects(dir string) ([]string, error) {
	subprojectPaths := []string{}
	// Walk through all Go files in the directory and its sub-directories
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".mod") {
			// Found a go.mod file, add the directory to the list of subprojects
			subprojectPaths = append(subprojectPaths, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
		return nil, err
	}
	return subprojectPaths, nil
}
