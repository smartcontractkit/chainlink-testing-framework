package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	// golang "github.com/smartcontractkit/chainlink-testing-framework/ci/golang"

	"dagger.io/dagger"
)

func main() {
	if err := verifyGoMod(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("go.mod and go.sum are up to date")
}

func verifyGoMod(ctx context.Context) error {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	container, err := GetGolangContainer(client, fmt.Sprintf("golang:%s", os.Getenv("GO_VERSION")))
	if err != nil {
		return err
	}

	err = GoVersion(ctx, container)
	if err != nil {
		return err
	}

	// err = GoModTidy(ctx, container)
	// if err != nil {
	// 	return err
	// }
	// err = VerifyTidy(ctx, container)
	// if err != nil {
	// 	return err
	// }

	err = RunTests(ctx, container)
	if err != nil {
		return err
	}

	return nil
}

func GetGolangContainer(client *dagger.Client, golangVersion string) (*dagger.Container, error) {
	src := client.Host().Directory(".").WithoutFile("./local_ci_cache")
	container := client.Container()

	// if we want to use a cache to speed up tests locally use a prebuilt tar file
	localCacheImage := os.Getenv("LOCAL_CACHE_IMAGE")
	if len(localCacheImage) > 0 {
		container = container.From(localCacheImage)
	} else {
		container = container.From(golangVersion)
	}
	container = container.WithDirectory("/src", src).
		WithWorkdir("/src")

	return container, nil
}

func MountCache(client *dagger.Client, container *dagger.Container) (*dagger.Container, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	goModCacheDirPath := homeDir + "/go/pkg/mod"
	goModCacheDir := client.Host().Directory(goModCacheDirPath)
	return container.WithMountedDirectory("/go/pkg/mod", goModCacheDir), nil
}

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

// GoVersion print the go version
func GoVersion(ctx context.Context, container *dagger.Container) error {
	out, err := container.WithExec([]string{"go", "version"}).Stdout(ctx)
	fmt.Println(out)
	if err != nil {
		return fmt.Errorf("go version error: %w", err)
	}

	return nil
}

func RunTests(ctx context.Context, container *dagger.Container) error {
	// Define the command to list Go packages
	listCmd := []string{"sh", "-c", "go list ./... | grep -v /k8s/e2e/ | grep -v /k8s/examples/ | grep -v /docker/test_env"}

	// Execute the package list command
	listOutput, err := container.WithExec(listCmd).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error executing package list command: %w", err)
	}

	// Parse the output to get package paths
	packagePaths := strings.Split(listOutput, "\n")

	// Prepare the Go test command
	testCmd := []string{"sh", "-c", fmt.Sprintf("go test -timeout 10m -cover -covermode=count -coverprofile=unit-test-coverage.out %s 2>&1", strings.Join(packagePaths, " "))}

	// Execute the test command
	testOutput, err := container.WithExec(testCmd).Stdout(ctx)
	if err != nil {
		return fmt.Errorf("error executing test command: %w", err)
	}

	// Output test results
	fmt.Println("Test Command Output:", testOutput)
	return nil
}
