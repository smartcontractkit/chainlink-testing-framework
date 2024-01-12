package osutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// both stop file and target file exist in working directory
func TestFindFile_ExistingInWorkingDir(t *testing.T) {
	_, err := os.Create("stopFile.txt")
	require.NoError(t, err, "error creating temp stop file")

	t.Cleanup(func() {
		_ = os.Remove("stopFile.txt")
	})

	tmpfile, err := os.Create("example")
	require.NoError(t, err, "error creating temp test file")

	t.Cleanup(func() {
		_ = os.Remove("example")
	})

	foundPath, err := FindFile("example", "stopFile.txt", 10)
	require.NoError(t, err, "error calling FindFile")

	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")
	expectedPath := filepath.Join(wd, tmpfile.Name())

	require.Equal(t, expectedPath, expectedPath, "expected %v, got %v", expectedPath, foundPath)
}

// stop file exists in working directory
// target file exists in sub directory
func TestFindFile_ExistsInSubDir(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	err = os.Mkdir("baseDir", 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, "baseDir"))
	})

	baseDir := filepath.Join(wd, "baseDir")
	subDir := filepath.Join(baseDir, "subDir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err, "error creating sub dir")

	// create stop file in base dir
	stopFile := filepath.Join(baseDir, "stopFile.txt")
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	// create nested file in sub dir
	targetFileName := "target.txt"
	targetFile := filepath.Join(subDir, targetFileName)
	_, err = os.Create(targetFile)
	require.NoError(t, err, "error creating temp test file")

	err = os.Chdir(subDir)
	require.NoError(t, err, "error changing working directory")

	foundPath, err := FindFile(targetFileName, "stopFile.txt", 2)
	require.NoError(t, err, "error calling FindFile")
	require.Equal(t, targetFile, foundPath, "expected %v, got %v", targetFile, foundPath)
}

// stop file exists in parent directory relative to working directory
// target file exists in a sub directory
func TestFindFile_ExistsInSameDirAsStop(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	err = os.Mkdir("baseDir", 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, "baseDir"))
	})

	baseDir := filepath.Join(wd, "baseDir")
	// create stop file in base dir
	stopFile := filepath.Join(wd, "stopFile.txt")
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	// create nested file in base dir
	targetFileName := "target.txt"
	targetFile := filepath.Join(baseDir, targetFileName)
	_, err = os.Create(targetFile)
	require.NoError(t, err, "error creating temp test file")

	err = os.Chdir(baseDir)
	require.NoError(t, err, "error changing working directory")

	foundPath, err := FindFile(targetFileName, "stopFile.txt", 2)
	require.NoError(t, err, "error calling FindFile")
	require.Equal(t, targetFile, foundPath, "expected %v, got %v", targetFile, foundPath)
}

// file doesn't exist anywhere
func TestFindFile_FileDoesNotExist(t *testing.T) {
	_, err := os.Create("stopFile.txt")
	require.NoError(t, err, "error creating temp stop file")

	t.Cleanup(func() {
		_ = os.Remove("stopFile.txt")
	})

	_, err = FindFile("nonExistentFile.txt", "stopFile.txt", 2)
	require.Error(t, err, "expected error calling FindFile")
	require.Contains(t, "file does not exist", err.Error(), "got wrong error")
}

// stop file doesn't exist at all
func TestFindFile_StopFileDoesNotExist(t *testing.T) {
	_, err := FindFile("somefile.txt", "nonExistentStopFile.txt", 2)
	require.Error(t, err, "expected error calling FindFile")
	require.Contains(t, err.Error(), ErrStopFileNotFoundWithinLimit, "got wrong error")
}

// stop file doesn't exist within search limit
func TestFindFile_stopFileNotFoundWithinLimit(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	err = os.Mkdir("baseDir", 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, "baseDir"))
	})

	baseDir := filepath.Join(wd, "baseDir")

	currentDir := baseDir
	for i := 0; i <= 3; i++ {
		currentDir = filepath.Join(currentDir, fmt.Sprintf("nested%d", i))
		err = os.Mkdir(currentDir, 0755)
		require.NoError(t, err, "error creating temp nested dir")
	}

	stopFile := filepath.Join(baseDir, "stopFile.txt")
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	err = os.Chdir(currentDir)
	require.NoError(t, err, "error changing working directory")

	_, err = FindFile("target.txt", "stopFile.txt", 2)
	require.Error(t, err, "expected error calling FindFile")
	require.Contains(t, err.Error(), ErrStopFileNotFoundWithinLimit, "got wrong error")
}
