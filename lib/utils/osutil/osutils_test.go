package osutil

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

// both stop file and target file exist in working directory
func TestFindFile_ExistingInWorkingDir(t *testing.T) {
	st := "stopFile1.txt"
	_, err := os.Create(st)
	require.NoError(t, err, "error creating temp stop file")

	t.Cleanup(func() {
		_ = os.Remove(st)
	})

	tmpfile, err := os.Create("example")
	require.NoError(t, err, "error creating temp test file")

	t.Cleanup(func() {
		_ = os.Remove("example")
	})

	foundPath, err := FindFile("example", st, 10)
	require.NoError(t, err, "error calling FindFile")

	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")
	expectedPath := filepath.Join(wd, tmpfile.Name())

	require.Equal(t, expectedPath, expectedPath, "expected %v, got %v", expectedPath, foundPath)
}

// stop file exists in working directory
// target file exists in sub directory
func TestFindFile_ExistsInSubDir(t *testing.T) {
	st := "stopFile2.txt"
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	b := "baseDir2"
	s := "subDir2"
	err = os.Mkdir(b, 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, b))
		_ = os.Chdir(wd)
	})

	baseDir := filepath.Join(wd, b)
	subDir := filepath.Join(baseDir, s)
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err, "error creating sub dir")

	// create stop file in base dir
	stopFile := filepath.Join(baseDir, st)
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	// create nested file in sub dir
	targetFileName := "target2.txt"
	targetFile := filepath.Join(subDir, targetFileName)
	_, err = os.Create(targetFile)
	require.NoError(t, err, "error creating temp test file")

	err = os.Chdir(subDir)
	require.NoError(t, err, "error changing working directory")

	foundPath, err := FindFile(targetFileName, st, 2)
	require.NoError(t, err, "error calling FindFile")
	require.Equal(t, targetFile, foundPath, "expected %v, got %v", targetFile, foundPath)
}

// stop file exists in parent directory relative to working directory
// target file exists in a sub directory
func TestFindFile_ExistsInSameDirAsStop(t *testing.T) {
	st := "stopFile3.txt"
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	b := "baseDir3"
	err = os.Mkdir(b, 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, b))
	})

	baseDir := filepath.Join(wd, b)
	// create stop file in base dir
	stopFile := filepath.Join(wd, st)
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	t.Cleanup(func() {
		os.Remove(stopFile)
		_ = os.Chdir(wd)
	})

	// create nested file in base dir
	targetFileName := "target3.txt"
	targetFile := filepath.Join(baseDir, targetFileName)
	_, err = os.Create(targetFile)
	require.NoError(t, err, "error creating temp test file")

	err = os.Chdir(baseDir)
	require.NoError(t, err, "error changing working directory")

	foundPath, err := FindFile(targetFileName, st, 2)
	require.NoError(t, err, "error calling FindFile")
	require.Equal(t, targetFile, foundPath, "expected %v, got %v", targetFile, foundPath)
}

// file doesn't exist anywhere
func TestFindFile_FileDoesNotExist(t *testing.T) {
	st := "stopFile4.txt"
	_, err := os.Create(st)
	require.NoError(t, err, "error creating temp stop file")

	t.Cleanup(func() {
		_ = os.Remove(st)
	})

	_, err = FindFile("nonExistentFile.txt", st, 2)
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
	st := "stopFile5.txt"
	b := "baseDir5"
	wd, err := os.Getwd()
	require.NoError(t, err, "error getting working directory")

	err = os.Mkdir(b, 0755)
	require.NoError(t, err, "error creating temp base dir")
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(wd, b))
		_ = os.Chdir(wd)
	})

	baseDir := filepath.Join(wd, b)

	currentDir := baseDir
	for i := 0; i <= 3; i++ {
		currentDir = filepath.Join(currentDir, fmt.Sprintf("nested%d", i))
		err = os.Mkdir(currentDir, 0755)
		require.NoError(t, err, "error creating temp nested dir")
	}

	stopFile := filepath.Join(baseDir, st)
	_, err = os.Create(stopFile)
	require.NoError(t, err, "error creating temp stop file")

	err = os.Chdir(currentDir)
	require.NoError(t, err, "error changing working directory")

	tg := "target5.txt"
	_, err = FindFile(tg, st, 2)
	require.Error(t, err, "expected error calling FindFile")
	require.Contains(t, err.Error(), ErrStopFileNotFoundWithinLimit, "got wrong error")
}

func TestFindDirectoriesContainingFile(t *testing.T) {
	t.Run("file in root directory", func(t *testing.T) {
		// Create a temporary directory using the os package
		tempDir, err := os.MkdirTemp("", "rootcheck")
		if err != nil {
			require.NoError(t, err, "error creating temp directory")
		}
		defer os.RemoveAll(tempDir) // Clean up after yourself

		exampleFile := "examplefile.txt"
		_, err = os.CreateTemp(tempDir, exampleFile)
		if err != nil {
			require.NoError(t, err, "error creating temp file")
		}

		testRegex := regexp.MustCompile(".*examplefile.txt")
		dirs, err := FindDirectoriesContainingFile(tempDir, testRegex)
		require.NoError(t, err, "error calling FindDirectoriesContainingFile")
		require.Equal(t, []string{tempDir}, dirs)
	})

	t.Run("file in root and sub directory", func(t *testing.T) {
		// Create a temporary directory using the os package
		tempDir, err := os.MkdirTemp("", "rootcheck")
		if err != nil {
			require.NoError(t, err, "error creating temp directory")
		}
		defer os.RemoveAll(tempDir) // Clean up after yourself

		exampleFile := "examplefile.txt"
		_, err = os.CreateTemp(tempDir, exampleFile)
		if err != nil {
			require.NoError(t, err, "error creating temp file")
		}

		subDir, err := os.MkdirTemp(tempDir, "subdir")
		if err != nil {
			require.NoError(t, err, "error creating temp directory")
		}

		exampleFile2 := "examplefile.txt"
		_, err = os.CreateTemp(subDir, exampleFile2)
		if err != nil {
			require.NoError(t, err, "error creating temp file")
		}

		testRegex := regexp.MustCompile(".*examplefile.txt")
		dirs, err := FindDirectoriesContainingFile(tempDir, testRegex)
		require.NoError(t, err, "error calling FindDirectoriesContainingFile")
		require.Equal(t, []string{tempDir, subDir}, dirs)
	})

	t.Run("file in sub directory", func(t *testing.T) {
		// Create a temporary directory using the os package
		tempDir, err := os.MkdirTemp("", "rootcheck")
		if err != nil {
			require.NoError(t, err, "error creating temp directory")
		}
		defer os.RemoveAll(tempDir) // Clean up after yourself

		subDir, err := os.MkdirTemp(tempDir, "subdir")
		if err != nil {
			require.NoError(t, err, "error creating temp directory")
		}

		exampleFile2 := "examplefile.txt"
		_, err = os.CreateTemp(subDir, exampleFile2)
		if err != nil {
			require.NoError(t, err, "error creating temp file")
		}

		testRegex := regexp.MustCompile(".*examplefile.txt")
		dirs, err := FindDirectoriesContainingFile(tempDir, testRegex)
		require.NoError(t, err, "error calling FindDirectoriesContainingFile")
		require.Equal(t, []string{subDir}, dirs)
	})

}
