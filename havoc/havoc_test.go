package havoc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	// We are not testing with real k8s, namespace is just a placeholder that should match in snapshots/results
	Namespace      = "cl-cluster"
	TestDataDir    = "testdata"
	SnapshotDir    = filepath.Join(TestDataDir, "snapshot")
	ResultsDir     = filepath.Join(TestDataDir, "results")
	DeploymentsDir = filepath.Join(TestDataDir, "deployments")
	ConfigsDir     = filepath.Join(TestDataDir, "configs")
	OAPISpecs      = filepath.Join(TestDataDir, "openapi_specs")
)

var (
	AllExperimentTypes = []string{
		ChaosTypeFailure,
		ChaosTypeLatency,
		ChaosTypeGroupFailure,
		ChaosTypeGroupLatency,
		ChaosTypeStressMemory,
		ChaosTypeStressGroupMemory,
		ChaosTypeStressCPU,
		ChaosTypeStressGroupCPU,
		ChaosTypePartitionGroup,
		ChaosTypeHTTP,
		ChaosTypePartitionExternal,
		ChaosTypeBlockchainSetHead,
	}
)

func init() {
	InitDefaultLogging()
}

func setup(t *testing.T, podsInfoPath string, configPath string, resultsDir string) (*Controller, *PodsListResponse) {
	d, err := os.ReadFile(filepath.Join(DeploymentsDir, podsInfoPath))
	require.NoError(t, err)
	var plr *PodsListResponse
	err = json.Unmarshal(d, &plr)
	require.NoError(t, err)
	var cfg *Config
	if configPath != "" {
		cfg, err = ReadConfig(filepath.Join(ConfigsDir, configPath))
		require.NoError(t, err)
	} else {
		cfg = DefaultConfig()
		cfg.Havoc.Dir = filepath.Join(ResultsDir, resultsDir)
	}
	m, err := NewController(cfg)
	require.NoError(t, err)
	return m, plr
}

func TestSmokeParsingGenerating(t *testing.T) {
	type test struct {
		name         string
		podsDumpName string
		configName   string
		snapshotDir  string
		resultsDir   string
	}
	tests := []test{
		{
			name:         "can generate for 1 pod without groups",
			podsDumpName: "deployment_single_pod.json",
			configName:   "",
			snapshotDir:  "single_pod",
			resultsDir:   "single_pod",
		},
		{
			name:         "can generate for a component group",
			podsDumpName: "deployment_single_group.json",
			configName:   "",
			snapshotDir:  "single_group",
			resultsDir:   "single_group",
		},
		{
			name:         "standalone pods + component group + network group + blockchain experiments",
			podsDumpName: "deployment_crib_block_rewind.json",
			configName:   "crib-all.toml",
			snapshotDir:  "all",
			resultsDir:   "all",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, plr := setup(t, tc.podsDumpName, tc.configName, tc.resultsDir)
			_, _, err := m.generateSpecs(Namespace, plr)
			require.NoError(t, err)
			snapshotData, err := m.ReadExperimentsFromDir(AllExperimentTypes, filepath.Join(SnapshotDir, tc.snapshotDir))
			require.NoError(t, err)
			generatedData, err := m.ReadExperimentsFromDir(AllExperimentTypes, filepath.Join(ResultsDir, tc.resultsDir))
			require.NoError(t, err)
			require.Equal(t, len(snapshotData), len(generatedData))
			for i := range snapshotData {
				// Replace snapshot dir name to match it with expected results path
				snapshotData[i].Path = strings.ReplaceAll(snapshotData[i].Path, SnapshotDir, ResultsDir)
				require.Equal(t, snapshotData[i], generatedData[i])
			}
		})
	}
}

/*
These are just an easy way to enter debug with arbitrary config, or some tweaks, run it manually
*/
func TestManualGenerate(t *testing.T) {
	cfg, err := ReadConfig("havoc.toml")
	require.NoError(t, err)
	m, err := NewController(cfg)
	require.NoError(t, err)
	err = m.GenerateSpecs("cl-cluster")
	require.NoError(t, err)
}

func TestManualRun(t *testing.T) {
	cfg, err := ReadConfig("havoc.toml")
	require.NoError(t, err)
	m, err := NewController(cfg)
	require.NoError(t, err)
	err = m.Run()
	require.NoError(t, err)
}
