package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Props struct {
	Name string `envconfig:"MY_NAME" yaml:"name"`
}

func TestOverrideCodeEnv(t *testing.T) {
	t.Run("CL env and version", func(t *testing.T) {
		defaultCodeProps := map[string]interface{}{
			"replicas": "1",
			"env": map[string]interface{}{
				"database_url": "postgresql://postgres:node@0.0.0.0/chainlink?sslmode=disable",
			},
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   "public.ecr.aws/chainlink/chainlink",
					"version": "1.4.1-root",
				},
				"web_port": "6688",
				"p2p_port": "6690",
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "350m",
						"memory": "1024Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "350m",
						"memory": "1024Mi",
					},
				},
			},
			"db": map[string]interface{}{
				"stateful": false,
				"capacity": "1Gi",
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"cpu":    "250m",
						"memory": "256Mi",
					},
					"limits": map[string]interface{}{
						"cpu":    "250m",
						"memory": "256Mi",
					},
				},
			},
		}
		t.Setenv(EnvVarCLImage, "abc")
		t.Setenv(EnvVarCLTag, "def")
		MustEnvOverrideVersion(&defaultCodeProps)
		require.Equal(t, "abc", defaultCodeProps["chainlink"].(map[string]interface{})["image"].(map[string]interface{})["image"])
		require.Equal(t, "def", defaultCodeProps["chainlink"].(map[string]interface{})["image"].(map[string]interface{})["version"])
	})
}
