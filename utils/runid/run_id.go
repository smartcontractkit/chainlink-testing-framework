package runid

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// GetOrGenerateRunId returns the runId if it is not nil, otherwise it reads the .run.id file and returns the value, or
// creates it with a new UUID, saves to file and returns the value.
func GetOrGenerateRunId(maybeRunId *string) (string, error) {
	if maybeRunId != nil {
		return *maybeRunId, nil
	}

	file, err := os.OpenFile(".run.id", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var runId string

	for scanner.Scan() {
		runId = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if runId != "" {
		return runId, nil
	}

	runId = fmt.Sprintf("local_%s", uuid.NewString())

	if _, err := file.WriteString(runId); err != nil {
		return "", err
	}

	return runId, nil
}

// RemoveLocalRunId removes the .run.id file if it exists and the runId contains 'local' substring indicating local execution.
// In GHA we get run_id from TOML config, so we don't need to remove the file.
func RemoveLocalRunId(runId *string) error {
	if runId != nil && !strings.Contains(*runId, "local") {
		return nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return err
	}

	possiblePath := workingDir + "/.run.id"
	_, err = os.Stat(possiblePath)

	if err != nil {
		return err
	}

	err = os.Remove(possiblePath)
	if err != nil {
		return err
	}

	return nil
}
