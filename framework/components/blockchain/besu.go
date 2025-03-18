package blockchain

const (
	DefaultBesuPrivateKey1 = "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"
	DefaultBesuPrivateKey2 = "c87509a1c067bbde78beb793e6fa76530b6382a4c0241e5e4a9ec0a0f44dc0d3"
	DefaultBesuPrivateKey3 = "ae6ae8e5ccbfb04590405997ee2d52d2b330726137b875053c36d94e974d162f"
)

func defaultBesu(in *Input) {
	if in.Image == "" {
		in.Image = "hyperledger/besu:24.9.1"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
	if in.WSPort == "" {
		in.WSPort = "8546"
	}
}

func newBesu(in *Input) (*Output, error) {
	defaultBesu(in)
	req := baseRequest(in, WithWsEndpoint)

	req.Image = in.Image
	req.AlwaysPullImage = in.PullImage

	defaultCmd := []string{
		"--network=dev",
		"--miner-enabled",
		"--miner-coinbase=0xfe3b557e8fb62b89f4916b721be55ceb828dbd73",
		"--rpc-http-cors-origins=all",
		"--host-allowlist=*",
		"--rpc-ws-enabled",
		"--rpc-http-enabled",
		"--rpc-http-host", "0.0.0.0",
		"--rpc-ws-host", "0.0.0.0",
		"--rpc-http-port", in.Port,
		"--rpc-ws-port", in.WSPort,
		"--data-path=/tmp/tmpDatdir",
	}
	entryPoint := append(defaultCmd, in.DockerCmdParamsOverrides...)
	req.Cmd = entryPoint

	return createGenericEvmContainer(in, req, true)
}
