package soak_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"path/filepath"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/stretchr/testify/require"
)

func TestOCRSoak(t *testing.T) {
	err := actions.RunSoakTest(
		filepath.Join(utils.SoakRoot, "tests"),
		filepath.Join(utils.ProjectRoot, "remote.test"),
		"@soak-ocr",
		"chainlink-soak-ocr",
		6,
	)
	require.NoError(t, err, "Failed to run the test")
}

func TestKeeperSoak(t *testing.T) {
	err := actions.RunSoakTest(
		filepath.Join(utils.SoakRoot, "tests"),
		filepath.Join(utils.ProjectRoot, "remote.test"),
		"@soak-keeper-block-time",
		"chainlink-soak-keeper",
		6,
	)
	require.NoError(t, err, "Failed to run the test")
}
