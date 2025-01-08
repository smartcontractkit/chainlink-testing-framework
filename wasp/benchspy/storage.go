package benchspy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const DEFAULT_DIRECTORY = "performance_reports"

type LocalStorage struct {
	Directory string `json:"directory"`
}

func (l *LocalStorage) defaultDirectoryIfEmpty() {
	if l.Directory == "" {
		l.Directory = DEFAULT_DIRECTORY
	}
}

func (l *LocalStorage) cleanTestName(testName string) string {
	// nested tests might contain slashes, replace them with underscores
	return strings.ReplaceAll(testName, "/", "_")
}

// Store saves a test report as a JSON file in local storage.
// It organizes reports by test name and commit/tag, ensuring easy retrieval and management.
// Returns the absolute path of the stored report or an error if the operation fails.
func (l *LocalStorage) Store(testName, commitOrTag string, report interface{}) (string, error) {
	L.Debug().
		Str("Test name", testName).
		Str("Reference", commitOrTag).
		Msg("Storing report on local storage")

	l.defaultDirectoryIfEmpty()
	asJson, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(l.Directory); os.IsNotExist(err) {
		if err := os.MkdirAll(l.Directory, 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory %s", l.Directory)
		}
	}

	cleanTestName := l.cleanTestName(testName)
	reportFilePath := filepath.Join(l.Directory, fmt.Sprintf("%s-%s.json", cleanTestName, commitOrTag))
	reportFile, err := os.Create(reportFilePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", reportFilePath)
	}
	defer func() { _ = reportFile.Close() }()

	reader := bytes.NewReader(asJson)
	_, err = io.Copy(reportFile, reader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to write to file %s", reportFilePath)
	}

	abs, err := filepath.Abs(reportFilePath)
	if err != nil {
		return reportFilePath, nil
	}

	L.Info().
		Str("Test name", testName).
		Str("Reference", commitOrTag).
		Str("Report path", abs).
		Msg("Report stored successfully")

	return abs, nil
}

// Load retrieves a report from local storage based on the specified test name and optional commit or tag.
// It decodes the report into the provided interface, enabling users to access stored test results.
func (l *LocalStorage) Load(testName, commitOrTag string, report interface{}) error {
	l.defaultDirectoryIfEmpty()
	if testName == "" {
		return errors.New("test name is empty. Please set it and try again")
	}

	cleanTestName := l.cleanTestName(testName)

	var ref string
	if commitOrTag == "" {
		L.Info().
			Str("Test name", testName).
			Msg("Loading latest report from local storage")

		var refErr error
		ref, refErr = l.findRef(cleanTestName)
		if refErr != nil {
			return refErr
		}
	} else {
		L.Info().
			Str("Test name", testName).
			Str("Reference", commitOrTag).
			Msg("Loading report from local storage")
		ref = commitOrTag
	}

	reportFilePath := filepath.Join(l.Directory, fmt.Sprintf("%s-%s.json", cleanTestName, ref))

	reportFile, err := os.Open(reportFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", reportFilePath)
	}

	decoder := json.NewDecoder(reportFile)
	if err := decoder.Decode(report); err != nil {
		return errors.Wrapf(err, "failed to decode file %s", reportFilePath)
	}

	L.Info().
		Str("Test name", testName).
		Str("Reference", ref).
		Str("Report path", reportFilePath).
		Msg("Report loaded successfully")

	return nil
}

func (l *LocalStorage) findAllGitlikeReferences(cleanTestName string, entries []fs.DirEntry) ([]string, error) {
	var refs []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.Contains(entry.Name(), cleanTestName) {
			parts := strings.Split(entry.Name(), "-")
			if len(parts) == 2 {
				ref := strings.TrimSuffix(parts[len(parts)-1], ".json")
				refs = append(refs, ref)
			} else {
				return nil, errors.Errorf("invalid file name: %s. Expected: %s-<ref>.json", entry.Name(), cleanTestName)
			}
		}
	}

	return refs, nil
}

func (l *LocalStorage) findLatestGitRef(refs []string) (string, error) {
	L.Debug().
		Int("References found", len(refs)).
		Msg("Finding latest report based on Git history")

	L.Trace().
		Str("References", strings.Join(refs, ", ")).
		Msg("Resolving references to commit hashes")

	var ref string
	// Find git root
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = l.Directory
	out, err := cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to find git root")
	}
	gitRoot := strings.TrimSpace(string(out))

	// Resolve all refs to commit hashes
	resolvedRefs := make(map[string]string)
	for _, ref := range refs {
		cmd = exec.Command("git", "rev-parse", ref)
		cmd.Dir = gitRoot
		if out, err := cmd.Output(); err == nil {
			resolvedRefs[ref] = strings.TrimSpace(string(out))
			L.Trace().
				Str("Reference", ref).
				Str("Resolved", resolvedRefs[ref]).
				Msg("Reference resolved to commit hash")
		} else {
			L.Warn().
				Str("Reference", ref).
				Msg("Failed to resolve reference to commit hash")
		}
	}

	// Find latest among resolved commits
	var commitRefs []string
	for _, hash := range resolvedRefs {
		commitRefs = append(commitRefs, hash)
	}

	L.Debug().
		Str("Commits found", strings.Join(commitRefs, ", ")).
		Msg("Finding latest commit among resolved references")

	args := append([]string{"rev-list", "--topo-order", "--date-order", "--max-count=1"}, commitRefs...)
	cmd = exec.Command("git", args...)
	cmd.Dir = gitRoot
	out, err = cmd.Output()
	if err != nil {
		return "", errors.Wrap(err, "failed to find latest reference")
	}
	latestCommit := strings.TrimSpace(string(out))

	// Find original ref for this commit
	foundOriginal := false
	for origRef, hash := range resolvedRefs {
		if hash == latestCommit {
			ref = origRef
			foundOriginal = true
			break
		}
	}

	L.Debug().
		Str("Latest commit", latestCommit).
		Str("Original ref", ref).
		Msg("Found original reference for latest commit")

	if !foundOriginal {
		return "", fmt.Errorf("no file found for latest commit %s. This should never happen", latestCommit)
	}

	return ref, nil
}

func (l *LocalStorage) findRef(cleanTestName string) (string, error) {
	L.Debug().
		Str("Test name", cleanTestName).
		Msg("Finding available reports in local storage")

	entries, err := os.ReadDir(l.Directory)
	if err != nil {
		return "", errors.Wrap(err, "failed to read storage directory")
	}

	refs, refErr := l.findAllGitlikeReferences(cleanTestName, entries)
	if refErr != nil {
		return "", refErr
	}

	L.Debug().
		Str("Test name", cleanTestName).
		Int("References found", len(refs)).
		Msg("Found available reports in local storage")

	switch len(refs) {
	case 0:
		return "", fmt.Errorf("no reports found in directory %s", l.Directory)
	case 1:
		return refs[0], nil
	default:
		return l.findLatestGitRef(refs)
	}
}
