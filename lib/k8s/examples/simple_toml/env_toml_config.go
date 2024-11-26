package main

import (
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

func main() {
	// in actual implementation here you should read the config from TOML file instead of creating structs manually
	chainlinkConfig := ctf_config.ChainlinkImageConfig{
		Image:           ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
		Version:         ptr.Ptr("2.18.0"),
		PostgresVersion: ptr.Ptr("12.0"),
	}

	pyroscope := ctf_config.PyroscopeConfig{
		Enabled: ptr.Ptr(false),
	}

	config := struct {
		Chainlink ctf_config.ChainlinkImageConfig
		Pyroscope ctf_config.PyroscopeConfig
	}{
		Chainlink: chainlinkConfig,
		Pyroscope: pyroscope,
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(&chainlinkConfig, target)
		ctf_config.MightConfigOverridePyroscopeKey(&pyroscope, target)
	}

	err := environment.New(&environment.Config{
		NamespacePrefix:   "ztest",
		KeepConnection:    true,
		RemoveOnInterrupt: true,
	}).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.NewWithOverride(0, map[string]interface{}{
			"replicas": 1,
		}, &config, overrideFn)).
		Run()
	if err != nil {
		panic(err)
	}
}
