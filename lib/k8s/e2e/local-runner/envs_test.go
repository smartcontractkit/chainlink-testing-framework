package env_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/e2e/common"
)

func TestMultiStageMultiManifestConnection(t *testing.T) {
	common.TestMultiStageMultiManifestConnection(t)
}

func TestConnectWithoutManifest(t *testing.T) {
	common.TestConnectWithoutManifest(t)
}

func Test5NodesSoakEnvironmentWithPVCs(t *testing.T) {
	common.Test5NodesSoakEnvironmentWithPVCs(t)
}

func TestWithSingleNodeEnv(t *testing.T) {
	common.TestWithSingleNodeEnvParallel(t)
}

func TestWithSingleNodeEnvLocalCharts(t *testing.T) {
	t.Setenv(config.EnvVarLocalCharts, "true")
	common.TestWithSingleNodeEnv(t)
}

func TestMultipleNodeWithDiffDBVersionEnv(t *testing.T) {
	common.TestMultipleNodeWithDiffDBVersionEnv(t)
}

func TestMinResources5NodesEnv(t *testing.T) {
	common.TestMinResources5NodesEnv(t)
}

func TestMinResources5NodesEnvWithBlockscout(t *testing.T) {
	common.TestMinResources5NodesEnvWithBlockscout(t)
}

func TestMultipleInstancesOfTheSameType(t *testing.T) {
	common.TestMultipleInstancesOfTheSameType(t)
}

func Test5NodesPlus2MiningGethsReorgEnv(t *testing.T) {
	common.Test5NodesPlus2MiningGethsReorgEnv(t)
}

func TestWithChaos(t *testing.T) {
	common.TestWithChaos(t)
}

func TestEmptyEnvironmentStartup(t *testing.T) {
	common.TestEmptyEnvironmentStartup(t)
}

func TestRolloutRestartUpdate(t *testing.T) {
	t.Skip("Always fails")
	common.TestRolloutRestart(t, true)
}

func TestRolloutRestartBySelector(t *testing.T) {
	common.TestRolloutRestart(t, false)
}

func TestReplaceHelm(t *testing.T) {
	common.TestReplaceHelm(t)
}

func TestRunTimeout(t *testing.T) {
	common.TestRunTimeout(t)
}

func TestReallyLongLogs(t *testing.T) {
	common.TestReallyLongLogs(t)
}

func TestWithAnvil(t *testing.T) {
	common.TestAnvil(t)
}
