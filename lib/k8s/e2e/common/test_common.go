package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/client"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/presets"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	TestEnvType = "chainlink-testing-framework-k8s-test"
)

var (
	testSelector = fmt.Sprintf("envType=%s", TestEnvType)
)

func GetTestEnvConfig(t *testing.T) *environment.Config {
	testType := "functional"
	product := "internal-k8s-test-runner"

	nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(product, testType)
	require.NoError(t, err, "Error creating required chain.link labels for namespace")

	workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(product, testType)
	require.NoError(t, err, "Error creating required chain.link labels for workloads and pods")

	return &environment.Config{
		NamespacePrefix:                       TestEnvType,
		Test:                                  t,
		Labels:                                nsLabels,
		WorkloadLabels:                        workloadPodLabels,
		PodLabels:                             workloadPodLabels,
		SkipRequiredChainLinkLabelsValidation: true,
	}
}

func TestMultiStageMultiManifestConnection(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	testEnvConfig := GetTestEnvConfig(t)

	ethChart := ethereum.New(nil)
	ethNetworkName := ethChart.GetProps().(*ethereum.Props).NetworkName

	// we adding the same chart with different index and executing multi-stage deployment
	// connections should be renewed
	e := environment.New(testEnvConfig)
	err := e.AddHelm(ethChart).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
			"toml":     presets.BaseToml,
		})).
		Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
	require.Len(t, e.URLs[chainlink.NodesLocalURLsKey], 1)
	require.Len(t, e.URLs[chainlink.NodesInternalURLsKey], 1)
	require.Len(t, e.URLs[chainlink.DBsLocalURLsKey], 1)
	require.Len(t, e.URLs, 7)

	err = e.AddHelm(chainlink.New(1, map[string]any{
		"replicas": 1,
		"toml":     presets.BaseToml,
	})).
		Run()
	require.NoError(t, err)
	require.Len(t, e.URLs[chainlink.NodesLocalURLsKey], 2)
	require.Len(t, e.URLs[chainlink.NodesInternalURLsKey], 2)
	require.Len(t, e.URLs[chainlink.DBsLocalURLsKey], 2)
	require.Len(t, e.URLs, 7)

	urls := make([]string, 0)
	if e.Cfg.InsideK8s {
		urls = append(urls, e.URLs[chainlink.NodesInternalURLsKey]...)
		urls = append(urls, e.URLs[ethNetworkName+"_internal_http"]...)
	} else {
		urls = append(urls, e.URLs[chainlink.NodesLocalURLsKey]...)
		urls = append(urls, e.URLs[ethNetworkName+"_http"]...)
	}

	r := resty.New()
	for _, u := range urls {
		l.Info().Str("URL", u).Send()
		res, err := r.R().Get(u)
		require.NoError(t, err)
		require.Equal(t, "200 OK", res.Status())
	}
}

func TestConnectWithoutManifest(t *testing.T) {
	l := logging.GetTestLogger(t)
	existingEnvConfig := GetTestEnvConfig(t)
	testEnvConfig := GetTestEnvConfig(t)
	existingEnvAlreadySetupVar := "ENV_ALREADY_EXISTS"
	var existingEnv *environment.Environment

	// only run this section if we don't already have an existing environment
	// needed for remote runner based tests to prevent duplicate envs from being created
	if os.Getenv(existingEnvAlreadySetupVar) == "" {
		existingEnv = environment.New(existingEnvConfig)
		l.Info().Str("Namespace", existingEnvConfig.Namespace).Msg("Existing Env Namespace")
		// deploy environment to use as an existing one for the test
		require.NotNil(t, existingEnv, "existingEnv is nil")
		require.NotNil(t, existingEnv.Cfg, "existingEnv.Cfg is nil %v", existingEnv)
		existingEnv.Cfg.JobImage = ""
		existingEnv.AddHelm(ethereum.New(nil)).
			AddHelm(chainlink.New(0, map[string]any{
				"replicas": 1,
				"toml":     presets.BaseToml,
			}))
		err := existingEnv.Run()
		require.NoError(t, err)
		// propagate the existing environment to the remote runner
		t.Setenv(fmt.Sprintf("TEST_%s", existingEnvAlreadySetupVar), "abc")
		// set the namespace to the existing one for local runs
		testEnvConfig.Namespace = existingEnv.Cfg.Namespace
	} else {
		l.Info().Msg("Environment already exists, verifying it is correct")
		require.NotEmpty(t, os.Getenv(config.EnvVarNamespace))
		skipManifestUpdate, err := strconv.ParseBool(os.Getenv(config.EnvVarSkipManifestUpdate))
		require.NoError(t, err, "Failed to parse the no manifest update env var")
		require.True(t, skipManifestUpdate)
	}

	// Now run an environment without a manifest like a normal test
	testEnvConfig.SkipManifestUpdate = true
	testEnv := environment.New(testEnvConfig)
	l.Info().Msgf("Testing Env Namespace %s", testEnv.Cfg.Namespace)
	err := testEnv.AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
			"toml":     presets.BaseToml,
		})).
		Run()
	require.NoError(t, err)
	if testEnv.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, testEnv.Shutdown())
	})

	connection := client.LocalConnection
	if testEnv.Cfg.InsideK8s {
		connection = client.RemoteConnection
	}
	url, err := testEnv.Fwd.FindPort("chainlink-0:node-1", "node", "access").As(connection, client.HTTP)
	require.NoError(t, err)
	urlGeth, err := testEnv.Fwd.FindPort("geth:0", "geth-network", "http-rpc").As(connection, client.HTTP)
	require.NoError(t, err)
	r := resty.New()
	l.Info().Msgf("getting url: %s", url)
	res, err := r.R().Get(url)
	require.NoError(t, err)
	require.Equal(t, "200 OK", res.Status())
	l.Info().Msgf("getting url: %s", url)
	res, err = r.R().Get(urlGeth)
	require.NoError(t, err)
	require.Equal(t, "200 OK", res.Status())
	l.Info().Msgf("done getting url: %s", url)
}

func Test5NodesSoakEnvironmentWithPVCs(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := presets.EVMSoak(testEnvConfig)
	err := e.Run()
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestWithSingleNodeEnvParallel(t *testing.T) {
	t.Parallel()
	TestWithSingleNodeEnv(t)
}

func TestWithSingleNodeEnv(t *testing.T) {
	testEnvConfig := GetTestEnvConfig(t)
	e, err := presets.EVMOneNode(testEnvConfig)
	require.NoError(t, err)
	err = e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestMultipleNodeWithDiffDBVersionEnv(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := presets.EVMMultipleNodesWithDiffDBVersion(testEnvConfig)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestMinResources5NodesEnv(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := presets.EVMMinimalLocal(testEnvConfig)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestMinResources5NodesEnvWithBlockscout(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e, err := presets.EVMMinimalLocalBS(testEnvConfig)
	require.NoError(t, err)
	err = e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func Test5NodesPlus2MiningGethsReorgEnv(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e, err := presets.EVMReorg(testEnvConfig)
	require.NoError(t, err)
	err = e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestMultipleInstancesOfTheSameType(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := environment.New(testEnvConfig).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]any{
			"replicas": 1,
			"toml":     presets.BaseToml,
		})).
		AddHelm(chainlink.New(1, map[string]any{
			"replicas": 1,
			"toml":     presets.BaseToml,
		}))
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

// TestWithChaos runs a test with chaos injected into the environment.
func TestWithChaos(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	appLabel := "chainlink-0"
	testCase := struct {
		chaosFunc  chaos.ManifestFunc
		chaosProps *chaos.Props
	}{
		chaos.NewFailPods,
		&chaos.Props{
			LabelsSelector: &map[string]*string{client.AppLabel: ptr.Ptr(appLabel)},
			DurationStr:    "30s",
		},
	}
	testEnvConfig := GetTestEnvConfig(t)
	cd := chainlink.New(0, map[string]any{
		"replicas": 1,
		"toml":     presets.BaseToml,
	})

	e := environment.New(testEnvConfig).
		AddHelm(ethereum.New(nil)).
		AddHelm(cd)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})

	url := e.URLs["chainlink_local"][0]
	r := resty.New()
	res, err := r.R().Get(url)
	require.NoError(t, err)
	require.Equal(t, "200 OK", res.Status())

	// start chaos
	_, err = e.Chaos.Run(testCase.chaosFunc(e.Cfg.Namespace, testCase.chaosProps))
	require.NoError(t, err)
	gom := gomega.NewGomegaWithT(t)
	gom.Eventually(func(g gomega.Gomega) {
		res, err = r.R().Get(url)
		g.Expect(err).Should(gomega.HaveOccurred())
		l.Info().Msg("Expected error was found")
	}, "1m", "3s").Should(gomega.Succeed())

	l.Info().Msg("Waiting for Pod to start back up")
	err = e.Run()
	require.NoError(t, err)

	// verify that the node can receive requests again
	url = e.URLs["chainlink_local"][0]
	res, err = r.R().Get(url)
	require.NoError(t, err)
	require.Equal(t, "200 OK", res.Status())
}

func TestEmptyEnvironmentStartup(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := environment.New(testEnvConfig)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}

func TestRolloutRestart(t *testing.T, statefulSet bool) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	cd := chainlink.New(0, map[string]any{
		"replicas": 5,
		"toml":     presets.BaseToml,
		"db": map[string]any{
			"stateful": true,
			"capacity": "1Gi",
		},
	})

	e := environment.New(testEnvConfig).
		AddHelm(ethereum.New(nil)).
		AddHelm(cd)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})

	if statefulSet {
		err = e.RolloutStatefulSets()
		require.NoError(t, err, "failed to rollout statefulsets")
	} else {
		err = e.RolloutRestartBySelector("deployment", testSelector)
		require.NoError(t, err, "failed to rollout restart deployment")
	}

	err = e.Run()
	require.NoError(t, err, "failed to run environment")
}

func TestReplaceHelm(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	cd := chainlink.New(0, map[string]any{
		"toml": presets.BaseToml,
		"chainlink": map[string]any{
			"resources": map[string]any{
				"requests": map[string]any{
					"cpu": "350m",
				},
				"limits": map[string]any{
					"cpu": "350m",
				},
			},
		},
	})

	e := environment.New(testEnvConfig).AddHelm(cd)
	err := e.Run()
	require.NoError(t, err)
	if e.WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
	require.NoError(t, err)
	cd = chainlink.New(1, map[string]any{
		"toml": presets.BaseToml,
		"chainlink": map[string]any{
			"resources": map[string]any{
				"requests": map[string]any{
					"cpu": "345m",
				},
				"limits": map[string]any{
					"cpu": "345m",
				},
			},
		},
	})
	require.NoError(t, err)
	e, err = e.
		ReplaceHelm("chainlink-0", cd)
	require.NoError(t, err)
	err = e.Run()
	require.NoError(t, err)
}

func TestRunTimeout(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e, err := presets.EVMOneNode(testEnvConfig)
	require.NoError(t, err)
	e.Cfg.ReadyCheckData.Timeout = 5 * time.Second
	err = e.Run()
	require.Error(t, err)
}

func TestReallyLongLogs(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)
	testEnvConfig := GetTestEnvConfig(t)
	val, _ := os.LookupEnv(config.EnvVarJobImage)
	if val != "" {
		env := environment.New(testEnvConfig)
		err := env.Run()
		require.NoError(t, err)
	}
	s := strings.Repeat("a", 500000)
	// this shouldn't hang
	l.Info().Int("len", len(s)).Str("string", s).Msg("string")
}

func TestAnvil(t *testing.T) {
	t.Parallel()
	testEnvConfig := GetTestEnvConfig(t)
	e := presets.FoundryNetwork(testEnvConfig)
	err := e.Run()
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, e.Shutdown())
	})
}
