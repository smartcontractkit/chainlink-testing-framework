package framework

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	EnvVarIgnoreCriticalLogs = "CTF_IGNORE_CRITICAL_LOGS"
)

func getLogDirectories() ([]string, error) {
	logDirs := make([]string, 0)
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "logs-") {
			logDirs = append(logDirs, filepath.Join(currentDir, entry.Name()))
		}
	}
	return logDirs, nil
}

// CheckCLNodeContainerErrors check if any CL node container logs has errors
func CheckCLNodeContainerErrors() error {
	dirs, err := getLogDirectories()
	if err != nil {
		return err
	}
	for _, dd := range dirs {
		if err := checkNodeLogErrors(dd); err != nil {
			return err
		}
	}
	return nil
}

// checkNodeLogsErrors check Chainlink nodes logs for error levels
func checkNodeLogErrors(dir string) error {
	if os.Getenv(EnvVarIgnoreCriticalLogs) == "true" {
		L.Warn().Msg(`CTF_IGNORE_CRITICAL_LOGS is set to true, we ignore all CRIT|FATAL|PANIC errors in node logs!`)
		return nil
	}
	fileRegex := regexp.MustCompile(`^node.*\.log$`)
	logLevelRegex := regexp.MustCompile(`(CRIT|PANIC|FATAL)`)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && fileRegex.MatchString(info.Name()) {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNumber := 1
			for scanner.Scan() {
				line := scanner.Text()
				if logLevelRegex.MatchString(line) {
					return fmt.Errorf("file %s contains a matching log level at line %d: %s", path, lineNumber, line)
				}
				lineNumber++
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}
		}
		return nil
	})
	return err
}
