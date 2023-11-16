package mirror

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/projectpath"
)

// findImageByName takes the name to search for as a string.
// It returns the name and version of the name if found, or an empty string with an error otherwise.
func findImageByName(name string) (string, error) {
	// Read the mirror json file from chainlink-testing-framework/mirror/mirror.json
	data, err := os.ReadFile(filepath.Join(projectpath.MirrorDir, "mirror.json"))
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON data into a slice of strings
	var images []string
	err = json.Unmarshal(data, &images)
	if err != nil {
		return "", err
	}

	// Iterate through each image name to find a match
	for _, image := range images {
		// Check if the name is a prefix of image (excluding version part)
		if strings.HasPrefix(image, name) {
			return image, nil
		}
	}

	// If the name is not found, return an error
	return "", fmt.Errorf("image '%s' not found in the mirrored images", name)
}

// GetImage gets the internal image name and version.
// If the internal docker repo is not set, it returns the original name and version.
func GetImage(name string) (string, error) {
	internalDockerRepo := os.Getenv(config.EnvVarInternalDockerRepo)
	// append a '/' if the internal docker repo is set
	if internalDockerRepo != "" {
		internalDockerRepo = fmt.Sprintf("%s/", internalDockerRepo)
	}
	image, err := findImageByName(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", internalDockerRepo, image), nil
}
