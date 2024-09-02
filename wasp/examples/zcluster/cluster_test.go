package main

import (
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClusterEntrypoint(t *testing.T) {
	p, err := wasp.NewClusterProfile(&wasp.ClusterConfig{
		Namespace: "wasp",
		// Builds and publishes test image if set to true
		UpdateImage:       true,
		DockerCmdExecPath: "..",
		BuildCtxPath:      ".",
		HelmValues: map[string]string{
			"env.loki.url":        os.Getenv("LOKI_URL"),
			"env.loki.token":      os.Getenv("LOKI_TOKEN"),
			"env.loki.basic_auth": os.Getenv("LOKI_BASIC_AUTH"),
			"env.loki.tenant_id":  os.Getenv("LOKI_TENANT_ID"),
			"image":               os.Getenv("WASP_TEST_IMAGE"),
			"test.binaryName":     os.Getenv("WASP_TEST_BIN"),
			"test.name":           os.Getenv("WASP_TEST_NAME"),
			"env.wasp.log_level":  "debug",
			"jobs":                "1",
			// other test vars pass through
			"test.MY_CUSTOM_VAR": "abc",
		},
	})
	require.NoError(t, err)
	err = p.Run()
	require.NoError(t, err)
}
