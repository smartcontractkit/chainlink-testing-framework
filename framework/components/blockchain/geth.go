package blockchain

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/testcontainers/testcontainers-go"
)

const (
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`
	GenesisClique     = `{
  "config": {
    "chainId": 1337,
    "homesteadBlock": 0,
    "daoForkBlock": 0,
    "daoForkSupport": true,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "muirGlacierBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "arrowGlacierBlock": 0,
    "grayGlacierBlock": 0,
    "shanghaiTime": 0,
    "clique": {
      "period": 1,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "800000000",
  "extradata": "0x0000000000000000000000000000000000000000000000000000000000000000f39Fd6e51aad88F6F4ce6aB8827279cffFb922660000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "alloc": {
    "f39Fd6e51aad88F6F4ce6aB8827279cffFb92266": {
      "balance": "20000000000000000000000"
    }
  }
}
`
	DefaultGethPrivateKey = `ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80`
)

var initScript = `
#!/bin/bash
if [ ! -d /root/.ethereum/keystore ]; then
	echo "/root/.ethereum/keystore not found, running 'geth init'..."
	geth init --datadir /root/.ethereum/devchain /root/genesis.json
	echo "...done!"
fi

geth init --datadir /root/.ethereum/devchain /root/genesis.json
geth "$@"
`

func defaultGeth(in *Input) {
	if in.Image == "" {
		in.Image = "ethereum/client-go:v1.13.8"
	}
	if in.ChainID == "" {
		in.ChainID = "1337"
	}
	if in.Port == "" {
		in.Port = "8545"
	}
}

func newGeth(in *Input) (*Output, error) {
	defaultGeth(in)
	req := baseRequest(in, WithoutWsEndpoint)
	defaultCmd := []string{
		"--http.corsdomain=*",
		"--http.vhosts=*",
		"--http",
		"--http.addr",
		"0.0.0.0",
		"--http.port",
		in.Port,
		"--http.api",
		"eth,net,web3,personal,debug",
		"--ws",
		"--ws.addr",
		"0.0.0.0",
		"--ws.port",
		in.Port,
		"--ws.api",
		"eth,net,web3,personal,debug",
		fmt.Sprintf("--networkid=%s", in.ChainID),
		"--ipcdisable",
		"--graphql",
		"-graphql.corsdomain", "*",
		"--allow-insecure-unlock",
		"--vmdebug",
		"--mine",
		"--miner.etherbase", RootFundingAddr,
		"--unlock", RootFundingAddr,
	}
	entryPoint := append(defaultCmd, in.DockerCmdParamsOverrides...)

	bindPort := req.ExposedPorts[0]

	initScriptFile, err := os.CreateTemp("", "init_script")
	if err != nil {
		return nil, err
	}
	_, err = initScriptFile.WriteString(initScript)
	if err != nil {
		return nil, err
	}
	genesisFile, err := os.CreateTemp("", "genesis_json")
	if err != nil {
		return nil, err
	}
	_, err = genesisFile.WriteString(GenesisClique)
	if err != nil {
		return nil, err
	}
	keystoreDir, err := os.MkdirTemp("", "keystore")
	if err != nil {
		return nil, err
	}
	key1File, err := os.CreateTemp(keystoreDir, "key1")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(RootFundingWallet)
	if err != nil {
		return nil, err
	}
	configDir, err := os.MkdirTemp("", "config")
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(configDir+"/password.txt", []byte(""), 0600)
	if err != nil {
		return nil, err
	}

	req.AlwaysPullImage = in.PullImage
	req.Image = in.Image
	req.Entrypoint = []string{
		"sh", "./root/init.sh",
		"--datadir", "/root/.ethereum/devchain",
		"--password", "/root/config/password.txt",
	}
	req.HostConfigModifier = func(h *container.HostConfig) {
		h.PortBindings = framework.MapTheSamePort(bindPort)
		framework.ResourceLimitsFunc(h, in.ContainerResources)
		h.Mounts = append(h.Mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   keystoreDir,
			Target:   "/root/.ethereum/devchain/keystore/",
			ReadOnly: false,
		}, mount.Mount{
			Type:     mount.TypeBind,
			Source:   configDir,
			Target:   "/root/config/",
			ReadOnly: false,
		})
	}
	req.Files = []testcontainers.ContainerFile{
		{
			HostFilePath:      initScriptFile.Name(),
			ContainerFilePath: "/root/init.sh",
			FileMode:          0644,
		},
		{
			HostFilePath:      genesisFile.Name(),
			ContainerFilePath: "/root/genesis.json",
			FileMode:          0644,
		},
	}
	req.Cmd = entryPoint

	return createGenericEvmContainer(in, req, false)
}
