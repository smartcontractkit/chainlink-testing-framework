package runid

import (
	"bufio"
	"os"

	"github.com/google/uuid"
)

func GetOrGenerateRunId() (string, error) {
	inOs := os.Getenv("RUN_ID")

	if inOs != "" {
		return inOs, nil
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

	runId = uuid.NewString()

	if _, err := file.WriteString(runId); err != nil {
		return "", err
	}

	return runId, nil
}

func RemoveLocalRunId() error {
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
