package benchspy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (l *LocalStorage) Store(testName, commitOrTag string, report interface{}) (string, error) {
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

	return abs, nil
}

func (l *LocalStorage) Load(testName, commitOrTag string, report interface{}) error {
	l.defaultDirectoryIfEmpty()
	if testName == "" {
		return errors.New("test name is empty. Please set it and try again")
	}

	cleanTestName := l.cleanTestName(testName)

	var ref string
	if commitOrTag == "" {
		entries, err := os.ReadDir(l.Directory)
		if err != nil {
			return errors.Wrap(err, "failed to read storage directory")
		}

		// Store both refs and file paths
		var refs []string
		filesByRef := make(map[string]string)
		for _, entry := range entries {
			if !entry.IsDir() && strings.Contains(entry.Name(), cleanTestName) {
				parts := strings.Split(entry.Name(), "-")
				if len(parts) == 2 {
					ref := strings.TrimSuffix(parts[len(parts)-1], ".json")
					refs = append(refs, ref)
					filesByRef[ref] = filepath.Join(l.Directory, entry.Name())
				} else {
					return errors.Errorf("invalid file name: %s. Expected: %s-<ref>.json", entry.Name(), cleanTestName)
				}
			}
		}

		if len(refs) == 0 {
			return fmt.Errorf("no reports found in directory %s", l.Directory)
		}

		if len(refs) > 1 {
			// Find git root
			cmd := exec.Command("git", "rev-parse", "--show-toplevel")
			cmd.Dir = l.Directory
			out, err := cmd.Output()
			if err != nil {
				return errors.Wrap(err, "failed to find git root")
			}
			gitRoot := strings.TrimSpace(string(out))

			// Resolve all refs to commit hashes
			resolvedRefs := make(map[string]string)
			for _, ref := range refs {
				cmd = exec.Command("git", "rev-parse", ref)
				cmd.Dir = gitRoot
				if out, err := cmd.Output(); err == nil {
					resolvedRefs[ref] = strings.TrimSpace(string(out))
				}
			}

			// Find latest among resolved commits
			var commitRefs []string
			for _, hash := range resolvedRefs {
				commitRefs = append(commitRefs, hash)
			}

			args := append([]string{"rev-list", "--topo-order", "--date-order", "--max-count=1"}, commitRefs...)
			cmd = exec.Command("git", args...)
			cmd.Dir = gitRoot
			out, err = cmd.Output()
			if err != nil {
				return errors.Wrap(err, "failed to find latest reference")
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

			if !foundOriginal {
				return fmt.Errorf("no file found for latest commit %s. This should never happen", latestCommit)
			}
		} else {
			ref = refs[0]
		}
	} else {
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

	return nil
}
