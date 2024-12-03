package comparator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type LocalReportStorage struct{}

func (l *LocalReportStorage) Store(testName, commitOrTag string, report interface{}) (string, error) {
	asJson, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return "", errors.Wrapf(err, "failed to create directory %s", directory)
		}
	}

	reportFilePath := filepath.Join(directory, fmt.Sprintf("%s-%s.json", testName, commitOrTag))
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

func (l *LocalReportStorage) Load(testName, commitOrTag string, report interface{}) error {
	if testName == "" {
		return errors.New("test name is empty. Please set it and try again")
	}

	if commitOrTag == "" {
		tagsOrCommits, tagErr := extractTagsOrCommits(directory)
		if tagErr != nil {
			return tagErr
		}

		latestCommit, commitErr := findLatestCommit(tagsOrCommits)
		if commitErr != nil {
			return commitErr
		}
		commitOrTag = latestCommit
	}
	reportFilePath := filepath.Join(directory, fmt.Sprintf("%s-%s.json", testName, commitOrTag))

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

func extractTagsOrCommits(directory string) ([]string, error) {
	pattern := regexp.MustCompile(`.+-(.+)\.json$`)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read directory %s", directory)
	}

	var tagsOrCommits []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := pattern.FindStringSubmatch(file.Name())
		if len(matches) == 2 {
			tagsOrCommits = append(tagsOrCommits, matches[1])
		}
	}

	return tagsOrCommits, nil
}

func findLatestCommit(references []string) (string, error) {
	if len(references) == 0 {
		return "", fmt.Errorf("no references provided")
	}

	args := append([]string{"rev-list", "--topo-order", "--date-order", "-n", "1"}, references...)
	cmd := exec.Command("git", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run git rev-list: %s, error: %v", stderr.String(), err)
	}

	latestCommit := strings.TrimSpace(stdout.String())
	if latestCommit == "" {
		return "", fmt.Errorf("no latest commit found")
	}

	return latestCommit, nil
}
