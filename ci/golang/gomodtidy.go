package golang

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

// GoModTidy run go mod tidy
func GoModTidy(ctx context.Context, container *dagger.Container) error {
	out, err := container.WithExec([]string{"go", "mod", "tidy"}).Stdout(ctx)
	if err != nil {
		// only print output on error
		fmt.Println(out)
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}
	return nil
}

// VerifyTidy verify go.mod and go.sum have changed
func VerifyTidy(ctx context.Context, container *dagger.Container) error {
	out, err := container.WithExec([]string{"git", "diff", "--stat", "--exit-code"}).Stdout(ctx)
	fmt.Println(out)
	if err != nil {
		return fmt.Errorf("go mod tidy: %w Please run `go mod tidy` on your project", err)
	}

	return nil
}
