package e2e_remote_runner_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/e2e/common"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/presets"
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

func TestFundReturnShutdownLogic(t *testing.T) {
	t.Parallel()
	testEnvConfig := common.GetTestEnvConfig(t)
	e := presets.EVMMinimalLocal(testEnvConfig)
	err := e.Run()
	if e.WillUseRemoteRunner() {
		require.Error(t, err, "Should return an error")
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
	require.NoError(t, err)
	fmt.Println(environment.FAILED_FUND_RETURN)
}

func TestFailedTestLogic(t *testing.T) {
	t.Skip("This test is meant to fail, and can only be evaluated by looking at the logs. Only turn on if checking this specific logic.")
	t.Parallel()
	testEnvConfig := common.GetTestEnvConfig(t)
	e := presets.OnlyRemoteRunner(testEnvConfig)
	err := e.Run()
	if e.WillUseRemoteRunner() {
		fmt.Println("Inside K8s?", e.Cfg.InsideK8s)
		fmt.Println("Test Failed?", e.Cfg.Test.Failed())
		require.True(t, e.Cfg.Test.Failed(), "Test should have failed")
		fmt.Println("This is a test-of-a-test and is confusing. The test that this tests should fail. But that also means this tests fails. If you're reading this, the test has actually passed.")
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
	require.NoError(t, err)
	fmt.Println("Inside K8s?", e.Cfg.InsideK8s)
	fmt.Println(environment.TEST_FAILED)
}

func TestRemoteRunnerOneSetupWithMultipleTests(t *testing.T) {
	t.Parallel()
	testEnvConfig := common.GetTestEnvConfig(t)
	ethChart := ethereum.New(nil)
	e := environment.New(testEnvConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethChart).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 5,
			"toml":     presets.BaseToml,
		}))
	err := e.Run()
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}

	// setup of variables to use for verification in a t.Run
	ethNetworkName := ethChart.GetProps().(*ethereum.Props).NetworkName
	urls := make([]string, 0)
	if e.Cfg.InsideK8s {
		urls = append(urls, e.URLs[chainlink.NodesInternalURLsKey]...)
		urls = append(urls, e.URLs[ethNetworkName+"_internal_http"]...)
	} else {
		urls = append(urls, e.URLs[chainlink.NodesLocalURLsKey]...)
		urls = append(urls, e.URLs[ethNetworkName+"_http"]...)
	}

	log.Info().Str("Test", "Before").Msg("Before Tests")

	// This test will verify we can connect a t.Run to the env of the parent test
	t.Run("do one", func(t1 *testing.T) {
		t1.Parallel()
		test1EnvConfig := common.GetTestEnvConfig(t1)
		test1EnvConfig.Namespace = e.Cfg.Namespace
		test1EnvConfig.SkipManifestUpdate = true
		e1 := presets.EVMMinimalLocal(test1EnvConfig)
		err := e1.Run()
		require.NoError(t1, err)
		log.Info().Str("Test", "One").Msg("Inside test")
		time.Sleep(1 * time.Second)
	})

	// This test will verify the sub t.Run properly uses the variables from the parent test
	t.Run("do two", func(t2 *testing.T) {
		t2.Parallel()
		log.Info().Str("Test", "Two").Msg("Inside test")
		r := resty.New()
		for _, u := range urls {
			log.Info().Str("URL", u).Send()
			res, err := r.R().Get(u)
			require.NoError(t2, err)
			require.Equal(t2, "200 OK", res.Status())
		}
	})

	log.Info().Str("Test", "After").Msg("After Tests")
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
