package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

func stringToBytesBuffer(s string) bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString(s)
	return buffer
}

func fileToString(t *testing.T, filepath string) string {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	return string(content)
}

func verifyAllItemsPresent(t *testing.T, expected, actual []string) {
	require.Equal(t, len(expected), len(actual))
	for _, item := range expected {
		require.Contains(t, actual, item)
	}
}

func TestGetChangedPackages(t *testing.T) {
	input := fileToString(t, "testdata/gitdiff.txt")
	listInput := fileToString(t, "testdata/golist.txt")
	config := SetConfig(ptr.Ptr("main"), ptr.Ptr(""), ptr.Ptr(""), ptr.Ptr(2))
	parsedPackages, err := parsePackages(stringToBytesBuffer(listInput))
	require.NoError(t, err)
	packages, err := getChangedPackages(stringToBytesBuffer(input), config.ProjectPath, config.Excludes, getGoFileMap(parsedPackages))
	fmt.Printf("packages: %++v\n", packages)
	require.NoError(t, err)
	expected := []string{
		"github.com/friendsofgo/killgrave/internal",
		"github.com/friendsofgo/killgrave/internal/server/http",
		"github.com/friendsofgo/killgrave/internal/app/cmd",
	}
	verifyAllItemsPresent(t, expected, packages)
}

func TestGetChangedPackagesWithExcludes(t *testing.T) {
	input := fileToString(t, "testdata/gitdiff.txt")
	listInput := fileToString(t, "testdata/golist.txt")
	config := SetConfig(ptr.Ptr("main"), ptr.Ptr(""), ptr.Ptr("internal/server/http/,internal/app/cmd/"), ptr.Ptr(2))
	parsedPackages, err := parsePackages(stringToBytesBuffer(listInput))
	require.NoError(t, err)
	packages, err := getChangedPackages(stringToBytesBuffer(input), config.ProjectPath, config.Excludes, getGoFileMap(parsedPackages))
	fmt.Printf("packages: %++v\n", packages)
	require.NoError(t, err)
	expected := []string{
		"github.com/friendsofgo/killgrave/internal",
	}
	verifyAllItemsPresent(t, expected, packages)
}

func TestGetChangedPackagesWithProjectPath(t *testing.T) {
	input := fileToString(t, "testdata/gitdiff.txt")
	listInput := fileToString(t, "testdata/golist.txt")
	config := SetConfig(ptr.Ptr("main"), ptr.Ptr("internal/server/http"), ptr.Ptr(""), ptr.Ptr(2))
	parsedPackages, err := parsePackages(stringToBytesBuffer(listInput))
	require.NoError(t, err)
	packages, err := getChangedPackages(stringToBytesBuffer(input), config.ProjectPath, config.Excludes, getGoFileMap(parsedPackages))
	fmt.Printf("packages: %++v\n", packages)
	require.NoError(t, err)
	expected := []string{
		"github.com/friendsofgo/killgrave/internal/server/http",
	}
	verifyAllItemsPresent(t, expected, packages)
}

func TestGetGoModChanges(t *testing.T) {
	input := fileToString(t, "testdata/gitdiffmod.txt")
	packages, err := getGoModChanges(stringToBytesBuffer(input))
	fmt.Printf("packages: %++v\n", packages)
	require.NoError(t, err)
	require.Equal(t, 1, len(packages))
	require.Equal(t, "github.com/stretchr/testify", packages[0])
}

func TestGetGoDepMap(t *testing.T) {
	input := fileToString(t, "testdata/golist.txt")
	parsedPackages, err := parsePackages(stringToBytesBuffer(input))
	require.NoError(t, err)
	packages := getGoDepMap(parsedPackages)
	require.Equal(t, 160, len(packages))
	p := packages["github.com/friendsofgo/killgrave/internal/server/http"]
	fmt.Printf("packages: %++v\n", p)
	require.Equal(t, 3, len(p))
	require.Equal(t, "github.com/friendsofgo/killgrave/cmd/killgrave", p[0])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app", p[1])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app/cmd", p[2])
}

func TestGetGoFileMap(t *testing.T) {
	input := fileToString(t, "testdata/golist.txt")
	parsedPackages, err := parsePackages(stringToBytesBuffer(input))
	require.NoError(t, err)
	packages := getGoFileMap(parsedPackages)
	require.Equal(t, 20, len(packages))
	require.Equal(t, "github.com/friendsofgo/killgrave/internal", packages["internal/config.go"])
}

func TestFindAffectedPackagesExternalDep(t *testing.T) {
	input := fileToString(t, "testdata/golist.txt")
	parsedPackages, err := parsePackages(stringToBytesBuffer(input))
	require.NoError(t, err)
	packages := getGoDepMap(parsedPackages)
	affected := findAffectedPackages("github.com/spf13/cobra", packages, true, 2)
	require.Equal(t, 3, len(affected))
	fmt.Printf("%++v\n", affected)
	require.Equal(t, "github.com/friendsofgo/killgrave/cmd/killgrave", affected[0])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app", affected[1])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app/cmd", affected[2])
}

func TestFindAffectedPackagesInternalDep(t *testing.T) {
	input := fileToString(t, "testdata/golist.txt")
	parsedPackages, err := parsePackages(stringToBytesBuffer(input))
	require.NoError(t, err)
	packages := getGoDepMap(parsedPackages)
	affected := findAffectedPackages("github.com/friendsofgo/killgrave/internal/server/http", packages, false, 2)
	require.Equal(t, 4, len(affected))
	fmt.Printf("%++v\n", affected)
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/server/http", affected[0])
	require.Equal(t, "github.com/friendsofgo/killgrave/cmd/killgrave", affected[1])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app", affected[2])
	require.Equal(t, "github.com/friendsofgo/killgrave/internal/app/cmd", affected[3])
}
