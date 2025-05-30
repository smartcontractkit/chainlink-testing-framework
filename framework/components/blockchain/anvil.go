package blockchain

import (
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultAnvilPrivateKey = `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
	AnvilPrivateKey1       = `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
	AnvilPrivateKey2       = `0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a`
	AnvilPrivateKey3       = `0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6`
	AnvilPrivateKey4       = `0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a`
	DefaultAnvilPublicKey  = `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`
	AnvilPublicKey1        = `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
	AnvilPublicKey2        = `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
	AnvilPublicKey3        = `0x90F79bf6EB2c4f870365E785982E1f101E93b906`
	AnvilPublicKey4        = `0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65`
)

func defaultAnvil(in *Input) {
	if in.Image == "" {
		in.Image = "ghcr.io/foundry-rs/foundry:stable"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
}

// newAnvil deploy foundry anvil node
func newAnvil(in *Input) (*Output, error) {
	defaultAnvil(in)
	req := baseRequest(in, WithoutWsEndpoint)

	req.Image = in.Image
	req.AlwaysPullImage = in.PullImage

	entryPoint := []string{"anvil"}
	defaultCmd := []string{"--host", "0.0.0.0", "--port", in.Port, "--chain-id", in.ChainID}
	entryPoint = append(entryPoint, defaultCmd...)
	entryPoint = append(entryPoint, in.DockerCmdParamsOverrides...)
	req.Entrypoint = entryPoint

	framework.L.Info().Any("Cmd", strings.Join(entryPoint, " ")).Msg("Creating anvil with command")

	return createGenericEvmContainer(in, req, false)
}
