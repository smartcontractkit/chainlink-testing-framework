package runid

import (
	"os"

	"github.com/google/uuid"
)

func GetOrGenerateRunId() string {
	inOs := os.Getenv("RUN_ID")

	if inOs != "" {
		return inOs
	}

	return uuid.NewString()[0:16]
}
