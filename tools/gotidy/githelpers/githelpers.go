package githelpers

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/osutil"
)

func outputFunc(m string) {
	fmt.Println(m)
}

func StashChanges() (err error) {
	fmt.Println("Doing a git stash before tidying...")
	return osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "git stash", outputFunc)
}

func PopStash() error {
	fmt.Println("Popping the stash to return previous changes...")
	return osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, "git stash pop", outputFunc)
}

func AddFile(file string) error {
	fmt.Printf("Adding %s to commit...\n", file)
	return osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, fmt.Sprintf("git add %s", file), outputFunc)
}

func CommitChanges(message string) error {
	fmt.Println("Committing changes...")
	return osutil.ExecCmdWithOptions(context.Background(), zerolog.Logger{}, fmt.Sprintf("git commit -m \"%s\"", message), outputFunc)
}
