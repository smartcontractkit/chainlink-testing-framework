package golang

import (
	"fmt"
	"os"

	"dagger.io/dagger"
)

func GolangImage(client *dagger.Client, golangVersion string) (*dagger.Container, error) {
	// Use a Docker image with Go and Git installed
	src := client.Host().Directory(".")
	container := client.Container().
		From(golangVersion).
		WithDirectory("/src", src).
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
