package flake

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

func TestValidFlakyTestFile(t *testing.T) {
	err := ValidateFileAgainstSchema(filepath.Join(utils.ProjectRoot, "flake", "test_data", "flaky_test_good.json"))
	require.NoError(t, err)
}

func TestInvalidFlakyTestFile(t *testing.T) {
	err := ValidateFileAgainstSchema(filepath.Join(utils.ProjectRoot, "flake", "test_data", "flaky_test_invalid.json"))
	require.Error(t, err)
}
