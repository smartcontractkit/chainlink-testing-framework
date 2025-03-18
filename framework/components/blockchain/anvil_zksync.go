package blockchain

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var AnvilZKSyncRichAccountPks = []string{
	"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
	"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
	"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
	"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
	"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
	"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
	"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
	"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
	"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
	"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
}

func defaultAnvilZksync(in *Input) {
	if in.ChainID == "" {
		in.ChainID = "260"
	}
	if in.Port == "" {
		in.Port = "8011"
	}
}

const dockerFile = `FROM ubuntu:latest
RUN apt update
RUN apt install -y curl git
RUN curl -L https://raw.githubusercontent.com/matter-labs/foundry-zksync/main/install-foundry-zksync | bash`

// anvil-zskync has no public image builds. this method will build the image from source
// creating a dockerfile in a temporary directory with the necessary commands to install
// foundry-zksync.
// see: https://foundry-book.zksync.io/getting-started/installation#using-foundry-with-docker
func newAnvilZksync(in *Input) (*Output, error) {
	defaultAnvilZksync(in)
	req := baseRequest(in, WithoutWsEndpoint)

	tempDir, err := os.MkdirTemp(".", "anvil-zksync-dockercontext")
	if err != nil {
		return nil, err
	}

	defer func() {
		// delete the folder wether it was successful or not
		_ = os.RemoveAll(tempDir)
	}()

	dockerfilePath := filepath.Join(tempDir, "anvilZksync.Dockerfile")

	if err := os.WriteFile(dockerfilePath, []byte(dockerFile), 0600); err != nil {
		return nil, err
	}

	req.FromDockerfile = testcontainers.FromDockerfile{
		Context:    tempDir,
		Dockerfile: "anvilZksync.Dockerfile",
		KeepImage:  true,
	}

	req.Entrypoint = []string{
		"/bin/sh",
		"-c",
		"/root/.foundry/bin/anvil-zksync" +
			" --chain-id " + in.ChainID +
			" --port " + in.Port}

	framework.L.Info().Any("Cmd", strings.Join(req.Entrypoint, " ")).Msg("Creating anvil zkSync with command")

	output, err := createGenericEvmContainer(in, req, false)
	if err != nil {
		return nil, err
	}

	return output, nil
}
