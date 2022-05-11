package soak_runner

import (
	"path/filepath"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

func TestSoak(t *testing.T) {
	actions.TestSoak(
		utils.ProjectRoot,
		filepath.Join(utils.ProjectRoot, "networks.yaml"),
		t,
		environment.NewChainlinkConfig(
			environment.ChainlinkReplicas(6, config.ChainlinkVals()),
			"chainlink-soak",
			config.GethNetworks()...,
		))
}
