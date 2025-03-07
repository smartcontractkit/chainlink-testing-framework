package benchspy

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testReport struct {
	Data string `json:"data"`
}

func TestBenchSpy_LocalStorage_Load(t *testing.T) {
	tempDir := t.TempDir()

	storage := &LocalStorage{
		Directory: tempDir,
	}

	// Create test data
	sampleReport := testReport{Data: "test data"}
	reportJSON, err := json.MarshalIndent(sampleReport, "", " ")
	require.NoError(t, err)

	t.Run("successful load with specific commit", func(t *testing.T) {
		// Setup
		commitID := "abc123"
		fileName := filepath.Join(tempDir, "test-abc123.json")
		require.NoError(t, os.WriteFile(fileName, reportJSON, 0600))

		// Test
		var loadedReport testReport
		err := storage.Load("test", commitID, &loadedReport)
		require.NoError(t, err)
		assert.Equal(t, sampleReport.Data, loadedReport.Data)
	})

	t.Run("error when test name is empty", func(t *testing.T) {
		var loadedReport testReport
		err := storage.Load("", "abc123", &loadedReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "test name is empty")
	})

	t.Run("error when file doesn't exist", func(t *testing.T) {
		var loadedReport testReport
		err := storage.Load("nonexistent", "abc123", &loadedReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("error with invalid JSON", func(t *testing.T) {
		// Setup
		commitID := "def456"
		fileName := filepath.Join(tempDir, "test-def456.json")
		require.NoError(t, os.WriteFile(fileName, []byte("invalid json"), 0600))

		// Test
		var loadedReport testReport
		err := storage.Load("test", commitID, &loadedReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode file")
	})

	t.Run("error when no reports found in directory", func(t *testing.T) {
		// Setup empty directory
		emptyDir := t.TempDir()
		storage := &LocalStorage{
			Directory: emptyDir,
		}

		// Attempt to load from empty directory
		var loadedReport testReport
		err := storage.Load("test", "", &loadedReport)

		// Verify error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no reports found in directory")
	})

	t.Run("error when no commits found for latest commit", func(t *testing.T) {
		// Setup git repo in temp dir
		gitDir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Create two reports
		fileName := filepath.Join(gitDir, "test-abc123.json")
		require.NoError(t, os.WriteFile(fileName, reportJSON, 0600))

		fileName = filepath.Join(gitDir, "test-abc1234.json")
		require.NoError(t, os.WriteFile(fileName, reportJSON, 0600))

		// Configure git for test
		//nolint
		configCmd := exec.Command("git", "config", "user.email", "test@example.com")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())
		//nolint
		configCmd = exec.Command("git", "config", "user.name", "Test User")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		storage := &LocalStorage{
			Directory: gitDir,
		}
		var loadedReport testReport
		err := storage.Load("test", "", &loadedReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find latest reference")
	})

	t.Run("loads latest git tag", func(t *testing.T) {
		// Setup git repo in temp dir
		gitDir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Configure git for test
		configCmd := exec.Command("git", "config", "user.email", "test@example.com")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())
		configCmd = exec.Command("git", "config", "user.name", "Test User")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		// Disable commit signing
		configCmd = exec.Command("git", "config", "--local", "commit.gpgsign", "false")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		// Create test files and tag them
		storage := &LocalStorage{Directory: gitDir}

		// First tag with older data
		oldReport := testReport{Data: "old data v1.0.0"}
		_, err := storage.Store("test", "v1.0.0", oldReport)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "first")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "tag", "v1.0.0")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Second tag with newer data
		newReport := testReport{Data: "new data v2.0.0"}
		_, err = storage.Store("test", "v2.0.0", newReport)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "second")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "tag", "v2.0.0")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Test loading latest tag - should get v2.0.0 data
		var loadedReport testReport
		err = storage.Load("test", "", &loadedReport)
		require.NoError(t, err)
		assert.Equal(t, newReport, loadedReport)
	})

	t.Run("loads latest commit", func(t *testing.T) {
		// Setup git repo in temp dir
		gitDir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Configure git for test
		configCmd := exec.Command("git", "config", "user.email", "test@example.com")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())
		configCmd = exec.Command("git", "config", "user.name", "Test User")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		// Disable commit signing
		configCmd = exec.Command("git", "config", "--local", "commit.gpgsign", "false")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		// Create test files and commit them
		storage := &LocalStorage{Directory: gitDir}

		// First commit with older data
		oldCommitData := testReport{Data: "old commit data"}
		_, err := storage.Store("test", "commit1", oldCommitData)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "first commit")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Get actual commit hash
		cmd = exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = gitDir
		commitHash1, err := cmd.Output()
		require.NoError(t, err)

		// Rename file to use actual commit hash
		oldPath := filepath.Join(gitDir, "test-commit1.json")
		newPath := filepath.Join(gitDir, fmt.Sprintf("test-%s.json", strings.TrimSpace(string(commitHash1))))
		require.NoError(t, os.Rename(oldPath, newPath))

		// Second commit with newer data
		newCommitData := testReport{Data: "new commit data"}
		_, err = storage.Store("test", "commit2", newCommitData)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "second commit")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Get actual commit hash and rename file
		cmd = exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = gitDir
		commitHash2, err := cmd.Output()
		require.NoError(t, err)
		oldPath = filepath.Join(gitDir, "test-commit2.json")
		newPath = filepath.Join(gitDir, fmt.Sprintf("test-%s.json", strings.TrimSpace(string(commitHash2))))
		require.NoError(t, os.Rename(oldPath, newPath))

		// Test loading latest commit - should get newest data
		var loadedReport testReport
		err = storage.Load("test", "", &loadedReport)
		require.NoError(t, err)
		assert.Equal(t, newCommitData, loadedReport)
	})

	t.Run("prefers newer commits over latest tag", func(t *testing.T) {
		// Setup git repo in temp dir
		gitDir := t.TempDir()
		cmd := exec.Command("git", "init")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Configure git for test
		configCmd := exec.Command("git", "config", "user.email", "test@example.com")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())
		configCmd = exec.Command("git", "config", "user.name", "Test User")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())
		configCmd = exec.Command("git", "config", "--local", "commit.gpgsign", "false")
		configCmd.Dir = gitDir
		require.NoError(t, configCmd.Run())

		storage := &LocalStorage{Directory: gitDir}

		// First tagged commit
		v1Data := testReport{Data: "v1.0.0 data"}
		_, err := storage.Store("test", "v1.0.0", v1Data)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "first tagged commit")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "tag", "v1.0.0")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Second tagged commit
		v2Data := testReport{Data: "v2.0.0 data"}
		_, err = storage.Store("test", "v2.0.0", v2Data)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "second tagged commit")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "tag", "v2.0.0")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())

		// Third untagged commit (newer than v2.0.0)
		newerData := testReport{Data: "newer untagged commit data"}
		_, err = storage.Store("test", "HEAD", newerData)
		require.NoError(t, err)
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "commit", "-m", "newer untagged commit")
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
		cmd = exec.Command("git", "rev-parse", "HEAD")
		cmd.Dir = gitDir
		commitHash, err := cmd.Output()
		require.NoError(t, err)
		oldPath := filepath.Join(gitDir, "test-HEAD.json")
		newPath := filepath.Join(gitDir, fmt.Sprintf("test-%s.json", strings.TrimSpace(string(commitHash))))
		require.NoError(t, os.Rename(oldPath, newPath))

		// Test loading - should newer commit
		var loadedReport testReport
		err = storage.Load("test", "", &loadedReport)
		require.NoError(t, err)
		assert.Equal(t, newerData, loadedReport)
	})
}

func TestBenchSpy_LocalStorage_Store(t *testing.T) {
	tempDir := t.TempDir()

	storage := &LocalStorage{
		Directory: tempDir,
	}

	sampleReport := testReport{Data: "test data"}

	t.Run("successful store with valid data", func(t *testing.T) {
		filePath, err := storage.Store("test", "abc123", sampleReport)
		require.NoError(t, err)
		assert.Contains(t, filePath, "test-abc123.json")

		// Verify file contents
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		var savedReport testReport
		require.NoError(t, json.Unmarshal(data, &savedReport))
		assert.Equal(t, sampleReport.Data, savedReport.Data)
	})

	t.Run("error with invalid JSON marshaling", func(t *testing.T) {
		invalidReport := make(chan int) // channels can't be marshaled
		_, err := storage.Store("test", "abc123", invalidReport)
		require.Error(t, err)
	})

	t.Run("error with invalid directory permissions", func(t *testing.T) {
		// Create read-only directory
		readOnlyDir := filepath.Join(tempDir, "readonly")
		require.NoError(t, os.MkdirAll(readOnlyDir, 0444))

		storage := &LocalStorage{
			Directory: readOnlyDir,
		}
		_, err := storage.Store("test", "abc123", sampleReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create file")
	})

	t.Run("error with directory creation", func(t *testing.T) {
		// Try to create directory in a non-existent parent
		storage := &LocalStorage{
			Directory: "/nonexistent/path/reports",
		}
		_, err := storage.Store("test", "abc123", sampleReport)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create directory")
	})

	t.Run("uses default directory when empty", func(t *testing.T) {
		storage := &LocalStorage{} // empty directory
		_, err := storage.Store("test", "abc123", sampleReport)
		require.NoError(t, err)

		// Verify file was created in default directory
		_, err = os.Stat(filepath.Join(DEFAULT_DIRECTORY, "test-abc123.json"))
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.RemoveAll(DEFAULT_DIRECTORY)
		})
	})

	t.Run("handles special characters in test name and commit", func(t *testing.T) {
		// Use URL-safe special characters instead
		filePath, err := storage.Store("test-with_special.chars", "commit-2.0_RC1", sampleReport)
		require.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(filePath)
		require.NoError(t, err)
	})

	t.Run("error with write permissions", func(t *testing.T) {
		// Create directory with read-only permission after creation
		restrictedDir := filepath.Join(tempDir, "restricted")
		require.NoError(t, os.MkdirAll(restrictedDir, 0755))
		require.NoError(t, os.Chmod(restrictedDir, 0444))

		storage := &LocalStorage{
			Directory: restrictedDir,
		}
		_, err := storage.Store("test", "abc123", sampleReport)
		require.Error(t, err)

		// Restore permissions for cleanup
		t.Cleanup(func() { _ = os.Chmod(restrictedDir, 0755) })
	})
}
func TestBenchSpy_LocalStorage_Load_GitEdgeCases(t *testing.T) {
	gitDir := t.TempDir()

	// Init git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = gitDir
	require.NoError(t, cmd.Run())

	// Configure git
	for _, config := range [][]string{
		{"user.email", "test@example.com"},
		{"user.name", "Test User"},
		{"commit.gpgsign", "false"},
	} {
		//nolint
		cmd := exec.Command("git", "config", "--local", config[0], config[1])
		cmd.Dir = gitDir
		require.NoError(t, cmd.Run())
	}

	storage := &LocalStorage{Directory: gitDir}

	t.Run("works with invalid git ref (1 file)", func(t *testing.T) {
		testData := testReport{Data: "test"}
		_, err := storage.Store("test", "invalid##ref", testData)
		require.NoError(t, err)

		var report testReport
		err = storage.Load("test", "", &report)
		require.NoError(t, err)
	})

	t.Run("error with non-existent commit", func(t *testing.T) {
		var report testReport
		err := storage.Load("test", "nonexistentcommit", &report)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("error with invalid git ref (2 files)", func(t *testing.T) {
		// Create 2 files with invalid ref
		testData := testReport{Data: "test"}
		_, err := storage.Store("test", "invalid##ref", testData)
		require.NoError(t, err)

		testData = testReport{Data: "test"}
		_, err = storage.Store("test", "invalid##ref2", testData)
		require.NoError(t, err)

		var report testReport
		err = storage.Load("test", "", &report)
		require.Error(t, err)
	})

	t.Run("error when git command fails", func(t *testing.T) {
		// Create 2 invalid files and break git repo
		require.NoError(t, os.RemoveAll(filepath.Join(gitDir, ".git")))

		testData := testReport{Data: "test"}
		_, err := storage.Store("test", "invalid##ref", testData)
		require.NoError(t, err)

		testData = testReport{Data: "test"}
		_, err = storage.Store("test", "invalid##ref2", testData)
		require.NoError(t, err)

		var report testReport
		err = storage.Load("test", "", &report)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find git root")
	})
}

func TestBenchSpy_LocalStorage_Load_PathEdgeCases(t *testing.T) {
	t.Run("loads with relative directory path", func(t *testing.T) {
		relDir := "./test_reports"
		t.Cleanup(func() { _ = os.RemoveAll(relDir) })

		storage := &LocalStorage{Directory: relDir}
		testData := testReport{Data: "test"}

		// Store and load with relative path
		_, err := storage.Store("test", "ref", testData)
		require.NoError(t, err)

		var loaded testReport
		err = storage.Load("test", "ref", &loaded)
		require.NoError(t, err)
		assert.Equal(t, testData, loaded)
	})

	t.Run("loads with absolute directory path", func(t *testing.T) {
		absDir, err := filepath.Abs(t.TempDir())
		require.NoError(t, err)

		storage := &LocalStorage{Directory: absDir}
		testData := testReport{Data: "test"}

		_, err = storage.Store("test", "ref", testData)
		require.NoError(t, err)

		var loaded testReport
		err = storage.Load("test", "ref", &loaded)
		require.NoError(t, err)
		assert.Equal(t, testData, loaded)
	})

	t.Run("handles directory with special characters", func(t *testing.T) {
		specialDir := filepath.Join(t.TempDir(), "test dir with spaces!")
		require.NoError(t, os.MkdirAll(specialDir, 0755))

		storage := &LocalStorage{Directory: specialDir}
		testData := testReport{Data: "test"}

		_, err := storage.Store("test", "ref", testData)
		require.NoError(t, err)

		var loaded testReport
		err = storage.Load("test", "ref", &loaded)
		require.NoError(t, err)
		assert.Equal(t, testData, loaded)
	})
}
