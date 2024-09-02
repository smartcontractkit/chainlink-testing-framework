package mirror

import (
	"fmt"
	"os"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
)

// AddMirrorToImageIfSet adds the internal docker repo to the image name if it is not already present.
func AddMirrorToImageIfSet(name string) string {
	ecr := os.Getenv(config.EnvVarInternalDockerRepo)
	if ecr != "" && !strings.HasPrefix(name, ecr) {
		name = fmt.Sprintf("%s/%s", ecr, name)
	}
	return name
}
