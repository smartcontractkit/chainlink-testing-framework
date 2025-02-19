package blockchain

import (
	"strings"

	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

func defaultAnvilZksync(in *Input) {
	if in.ChainID == "" {
		in.ChainID = "260"
	}
	if in.Port == "" {
		in.Port = "8011"
	}
}

func newAnvilZksync(in *Input) (*Output, error) {
	defaultAnvilZksync(in)
	req := baseRequest(in)

	// anvil-zskync has no public image builds.
	// see: https://foundry-book.zksync.io/getting-started/installation#using-foundry-with-docker
	req.FromDockerfile = testcontainers.FromDockerfile{
		Context:    ".",
		Dockerfile: "anvilZksync.Dockerfile",
		KeepImage:  true,
	}

	req.Entrypoint = []string{"/bin/sh", "-c", "/root/.foundry/bin/anvil-zksync --offline run"}

	framework.L.Info().Any("Cmd", strings.Join(req.Entrypoint, " ")).Msg("Creating anvil with command")

	return createGenericEvmContainer(in, req)
}
