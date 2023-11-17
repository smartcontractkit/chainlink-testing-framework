package test_env

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
)

const (
	NETHERMIND_IMAGE_TAG = "1.22.0"
)

type Nethermind struct {
	EnvComponent
	ExternalHttpUrl      string
	InternalHttpUrl      string
	ExternalWsUrl        string
	InternalWsUrl        string
	InternalExecutionURL string
	ExternalExecutionURL string
	ExecutionDir         string
	genesisTime          int
	consensusLayer       ConsensusLayer
	l                    zerolog.Logger
}

func NewNethermind(networks []string, executionDir string, consensusLayer ConsensusLayer, genesisTime int, opts ...EnvComponentOption) *Nethermind {
	g := &Nethermind{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "nethermind", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		ExecutionDir:   executionDir,
		consensusLayer: consensusLayer,
		genesisTime:    genesisTime,
		l:              log.Logger,
	}
	for _, opt := range opts {
		opt(&g.EnvComponent)
	}
	return g
}

func (g *Nethermind) WithLogger(l zerolog.Logger) *Nethermind {
	g.l = l
	return g
}

func (g *Nethermind) StartContainer() (blockchain.EVMNetwork, error) {
	r, err := g.getContainerRequest(g.Networks)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	ct, err := docker.StartContainerWithRetry(g.l, tc.GenericContainerRequest{
		ContainerRequest: *r,
		Reuse:            true,
		Started:          true,
		Logger:           &g.l,
	})
	if err != nil {
		return blockchain.EVMNetwork{}, errors.Wrapf(err, "cannot start nethermind container")
	}

	host, err := GetHost(context.Background(), ct)
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	httpPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_HTTP_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	wsPort, err := ct.MappedPort(context.Background(), NatPort(TX_GETH_WS_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}
	executionPort, err := ct.MappedPort(context.Background(), NatPort(ETH2_EXECUTION_PORT))
	if err != nil {
		return blockchain.EVMNetwork{}, err
	}

	g.Container = ct
	g.ExternalHttpUrl = FormatHttpUrl(host, httpPort.Port())
	g.InternalHttpUrl = FormatHttpUrl(g.ContainerName, TX_GETH_HTTP_PORT)
	g.ExternalWsUrl = FormatWsUrl(host, wsPort.Port())
	g.InternalWsUrl = FormatWsUrl(g.ContainerName, TX_GETH_WS_PORT)
	g.InternalExecutionURL = FormatHttpUrl(g.ContainerName, ETH2_EXECUTION_PORT)
	g.ExternalExecutionURL = FormatHttpUrl(host, executionPort.Port())

	networkConfig := blockchain.SimulatedEVMNetwork
	networkConfig.Name = fmt.Sprintf("nethermind-eth2-%s", g.consensusLayer)
	networkConfig.URLs = []string{g.ExternalWsUrl}
	networkConfig.HTTPURLs = []string{g.ExternalHttpUrl}

	g.l.Info().Str("containerName", g.ContainerName).
		Msg("Started Nethermind container")

	return networkConfig, nil
}

func (g *Nethermind) GetInternalExecutionURL() string {
	return g.InternalExecutionURL
}

func (g *Nethermind) GetExternalExecutionURL() string {
	return g.ExternalExecutionURL
}

func (g *Nethermind) GetInternalHttpUrl() string {
	return g.InternalHttpUrl
}

func (g *Nethermind) GetInternalWsUrl() string {
	return g.InternalWsUrl
}

func (g *Nethermind) GetExternalHttpUrl() string {
	return g.ExternalHttpUrl
}

func (g *Nethermind) GetExternalWsUrl() string {
	return g.ExternalWsUrl
}

func (g *Nethermind) GetContainerName() string {
	return g.ContainerName
}

func (g *Nethermind) GetContainer() *tc.Container {
	return &g.Container
}

func (g *Nethermind) getContainerRequest(networks []string) (*tc.ContainerRequest, error) {
	//this might fail in CI
	// keystoreDir, err := os.MkdirTemp(g.ExecutionDir, "keystore")
	// if err != nil {
	// 	return nil, err
	// }

	passwordFile, err := os.CreateTemp("", "password.txt")
	if err != nil {
		return nil, err
	}

	key1File, err := os.CreateTemp("", "key1")
	if err != nil {
		return nil, err
	}
	_, err = key1File.WriteString(`{"address":"123463a4b065722e99115d6c222f267d9cabb524","crypto":{"cipher":"aes-128-ctr","ciphertext":"93b90389b855889b9f91c89fd15b9bd2ae95b06fe8e2314009fc88859fc6fde9","cipherparams":{"iv":"9dc2eff7967505f0e6a40264d1511742"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"c07503bb1b66083c37527cd8f06f8c7c1443d4c724767f625743bd47ae6179a4"},"mac":"6d359be5d6c432d5bbb859484009a4bf1bd71b76e89420c380bd0593ce25a817"},"id":"622df904-0bb1-4236-b254-f1b8dfdff1ec","version":3}`)
	if err != nil {
		return nil, err
	}

	jwtSecret, err := os.CreateTemp("/tmp", "jwtsecret")
	if err != nil {
		return nil, err
	}
	_, err = jwtSecret.WriteString("0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345")
	if err != nil {
		return nil, err
	}

	netherChainspecFile, err := os.CreateTemp("", "chainspec.json")
	if err != nil {
		return nil, err
	}

	netherChainspec, err := buildNethermindChainspec([]string{}, g.genesisTime)
	if err != nil {
		return nil, err
	}

	_, err = netherChainspecFile.WriteString(netherChainspec)
	if err != nil {
		return nil, err
	}

	// fmt.Print(netherChainspec)

	//$NETHERMIND/Nethermind.Runner
	//--Init.ChainSpecPath ./nethermind.json
	//--JsonRpc.JwtSecretFile $DEVNET/jwt.hex
	//--Init.IsMining true
	//--KeyStore.PasswordFiles $DEVNET/keystore_password.txt
	//--KeyStore.UnlockAccounts 0x123463a4b065722e99115d6c222f267d9cabb524
	//--KeyStore.BlockAuthorAccount 0x123463a4b065722e99115d6c222f267d9cabb524 //what is this?
	//--Network.MaxActivePeers 0
	//--Init.GenesisHash 0x19286b9ba93edd49b64266c1df5e8303ecb3c504528b010b56f5b237aa0896c2 //we can use default
	//-c withdrawals_test

	//--Blocks.RandomizedBlocks -- could be interesting!
	//--Network.EnableUPnP -- enable and see what's the result?

	return &tc.ContainerRequest{
		Name:            g.ContainerName,
		AlwaysPullImage: true,
		Image:           fmt.Sprintf("nethermind/nethermind:%s", NETHERMIND_IMAGE_TAG),
		Networks:        networks,
		ExposedPorts:    []string{NatPortFormat(TX_GETH_HTTP_PORT), NatPortFormat(TX_GETH_WS_PORT), NatPortFormat(ETH2_EXECUTION_PORT)},
		WaitingFor: tcwait.ForAll(
			// NewHTTPStrategy("/", NatPort(TX_GETH_HTTP_PORT)),
			tcwait.ForLog("Nethermind initialization completed").
				WithStartupTimeout(120 * time.Second).
				WithPollInterval(1 * time.Second),
			// NewWebSocketStrategy(NatPort(TX_GETH_WS_PORT), g.l),
		),
		Cmd: []string{"--Init.ChainSpecPath=/nethermind/chainspec.json",
			"--Init.DiscoveryEnabled=false",
			"--Init.IsMining=true",
			// "--Init.GenesisHash=0xf27cb05e82e8851835d22c14c6d0330a27584eb62cb0e2bb2a3c655feabe6b6b",
			// "--Init.GenesisHash=0x602be38796a4746165fa88bd2ec516dd35eaca9eac476bd62832fbd890e9a72b",
			"--HealthChecks.Enabled=true", // default slug /health
			"--JsonRpc.Enabled=true",
			"--JsonRpc.EnabledModules=[eth,subscribe,trace,txPool,web3,personal,proof,net,parity,health,rpc, debug]",
			"--JsonRpc.Host=0.0.0.0",
			//			"--http.corsdomain=*",
			//			"--http.vhosts=*",
			fmt.Sprintf("--JsonRpc.Port=%s", TX_GETH_HTTP_PORT),
			"--JsonRpc.EngineHost=0.0.0.0",
			"--JsonRpc.EnginePort=" + ETH2_EXECUTION_PORT,
			// "--ws",
			// "--ws.api=admin,debug,web3,eth,txpool,net",
			// "--ws.addr=0.0.0.0",
			// "--ws.origins=*",
			fmt.Sprintf("--JsonRpc.WebSocketsPort=%s", TX_GETH_WS_PORT),
			// "--authrpc.vhosts=*",
			// "--authrpc.addr=0.0.0.0",
			"--JsonRpc.JwtSecretFile=/nethermind/shared/jwtsecret",
			// "--datadir=/execution",
			// "--rpc.allow-unprotected-txs",
			// "--rpc.txfeecap=0",
			// "--allow-insecure-unlock",
			"--Network.MaxActivePeers=0",
			"--KeyStore.BlockAuthorAccount=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--KeyStore.UnlockAccounts=0x123463a4b065722e99115d6c222f267d9cabb524",
			"--KeyStore.PasswordFiles=/nethermind/password.txt",
			// "--nodiscover",
			// "--syncmode=full",
			// "--networkid=1337",
			// "--graphql",
			// "--graphql.corsdomain=*",
		},
		Files: []tc.ContainerFile{
			{
				HostFilePath:      passwordFile.Name(),
				ContainerFilePath: "/nethermind/password.txt",
				FileMode:          0644,
			},
			{
				HostFilePath:      jwtSecret.Name(),
				ContainerFilePath: "/nethermind/shared/jwtsecret",
				FileMode:          0644,
			},
			{
				HostFilePath:      netherChainspecFile.Name(),
				ContainerFilePath: "/nethermind/chainspec.json",
				FileMode:          0644,
			},
			{
				HostFilePath:      key1File.Name(),
				ContainerFilePath: "/nethermind/keystore/key-123463a4b065722e99115d6c222f267d9cabb524",
				FileMode:          0644,
			},
		},
		Mounts: tc.ContainerMounts{
			tc.ContainerMount{
				Source: tc.GenericBindMountSource{
					HostPath: g.ExecutionDir,
				},
				Target: "/nethermind/shared",
			},
		},
	}, nil
}

func buildNethermindChainspec(addressesToFund []string, genesisTime int) (string, error) {
	for i := range addressesToFund {
		if has0xPrefix(addressesToFund[i]) {
			addressesToFund[i] = addressesToFund[i][2:]
		}
	}

	data := struct {
		AddressesToFund   []string
		GenesisTimestamp  string
		ShanghaiTimestamp string
	}{
		AddressesToFund:   addressesToFund,
		GenesisTimestamp:  fmt.Sprintf("%x", genesisTime),
		ShanghaiTimestamp: fmt.Sprintf("%x", genesisTime+30),
	}

	t, err := template.New("genesis-json").Funcs(funcMap).Parse(NethermindGeneisJSON)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}

//"timestamp": "0x{{ .GenesisTimestamp }}",

var NethermindGeneisJSON = `
{
	"name": "Private network",
	"engine": {
		"Ethash": {
			"params": {
			}
		}
	},	
	"params": {
		"gasLimitBoundDivisor": "0x400",
		"accountStartNonce": "0x0",
		"maximumExtraDataSize": "0x20",
		"minGasLimit": "0x1388",
		"chainID": 1337,
		"networkID": 1337,
		"eip150Transition": "0x0",
		"eip155Transition": "0x0",
		"eip158Transition": "0x0",
		"eip160Transition": "0x0",
		"eip161abcTransition": "0x0",
		"eip161dTransition": "0x0",
		"eip140Transition": "0x0",
		"eip211Transition": "0x0",
		"eip214Transition": "0x0",
		"eip658Transition": "0x0",
		"eip145Transition": "0x0",
		"eip1014Transition": "0x0",
		"eip1052Transition": "0x0",
		"eip1283Transition": "0x0",
		"eip1283DisableTransition": "0x0",
		"eip152Transition": "0x0",
		"eip1108Transition": "0x0",
		"eip1344Transition": "0x0",
		"eip1884Transition": "0x0",
		"eip2028Transition": "0x0",
		"eip2200Transition": "0x0",
		"eip2565Transition": "0x0",
		"eip2929Transition": "0x0",
		"eip2930Transition": "0x0",
		"eip1559Transition": "0x0",
		"eip3198Transition": "0x0",
		"eip3529Transition": "0x0",
		"eip3541Transition": "0x0",
		"eip4895TransitionTimestamp": "0x6553701d",
		"eip3855TransitionTimestamp": "0x6553701d",
		"eip3651TransitionTimestamp": "0x6553701d",
		"eip3860TransitionTimestamp": "0x6553701d",
		"eip4844TransitionTimestamp": "0x65565e1d",
		"eip4788TransitionTimestamp": "0x65565e1d",
		"eip1153TransitionTimestamp": "0x65565e1d",
		"eip5656TransitionTimestamp": "0x65565e1d",
		"eip6780TransitionTimestamp": "0x65565e1d"
		"terminalTotalDifficulty": "0x0"
		},
		"genesis": {
		"seal": {
			"ethereum": {
			"nonce": "0x42",
			"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
			}
		},
	  "difficulty": "0x1",
	  "author": "0x0000000000000000000000000000000000000000",
	  "timestamp": "0x{{ .GenesisTimestamp }}",
	  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	  "extradata": "0x0000000000000000000000000000000000000000000000000000000000000000123463a4B065722E99115D6c222f267d9cABb5240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	  "gasLimit":"0x1C9C380",
	  "baseFeePerGas":"0x7"
	},
	"accounts": {
	  "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b":
	  {
		"balance":"0x6d6172697573766477000000"
	  },
	  "123463a4b065722e99115d6c222f267d9cabb524": {
		"balance": "0x43c33c1937564800000"
		},
		"0x4242424242424242424242424242424242424242": {
			"balance": "0",
			"code": "0x60806040526004361061003f5760003560e01c806301ffc9a71461004457806322895118146100a4578063621fd130146101ba578063c5f2892f14610244575b600080fd5b34801561005057600080fd5b506100906004803603602081101561006757600080fd5b50357fffffffff000000000000000000000000000000000000000000000000000000001661026b565b604080519115158252519081900360200190f35b6101b8600480360360808110156100ba57600080fd5b8101906020810181356401000000008111156100d557600080fd5b8201836020820111156100e757600080fd5b8035906020019184600183028401116401000000008311171561010957600080fd5b91939092909160208101903564010000000081111561012757600080fd5b82018360208201111561013957600080fd5b8035906020019184600183028401116401000000008311171561015b57600080fd5b91939092909160208101903564010000000081111561017957600080fd5b82018360208201111561018b57600080fd5b803590602001918460018302840111640100000000831117156101ad57600080fd5b919350915035610304565b005b3480156101c657600080fd5b506101cf6110b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102095781810151838201526020016101f1565b50505050905090810190601f1680156102365780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561025057600080fd5b506102596110c7565b60408051918252519081900360200190f35b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a70000000000000000000000000000000000000000000000000000000014806102fe57507fffffffff0000000000000000000000000000000000000000000000000000000082167f8564090700000000000000000000000000000000000000000000000000000000145b92915050565b6030861461035d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118056026913960400191505060405180910390fd5b602084146103b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603681526020018061179c6036913960400191505060405180910390fd5b6060821461040f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806118786029913960400191505060405180910390fd5b670de0b6b3a7640000341015610470576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118526026913960400191505060405180910390fd5b633b9aca003406156104cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260338152602001806117d26033913960400191505060405180910390fd5b633b9aca00340467ffffffffffffffff811115610535576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602781526020018061182b6027913960400191505060405180910390fd5b6060610540826114ba565b90507f649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c589898989858a8a6105756020546114ba565b6040805160a0808252810189905290819060208201908201606083016080840160c085018e8e80828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910187810386528c815260200190508c8c808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690920188810386528c5181528c51602091820193918e019250908190849084905b83811015610648578181015183820152602001610630565b50505050905090810190601f1680156106755780820380516001836020036101000a031916815260200191505b5086810383528881526020018989808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018881038452895181528951602091820193918b019250908190849084905b838110156106ef5781810151838201526020016106d7565b50505050905090810190601f16801561071c5780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1600060028a8a600060801b604051602001808484808284377fffffffffffffffffffffffffffffffff0000000000000000000000000000000090941691909301908152604080517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0818403018152601090920190819052815191955093508392506020850191508083835b602083106107fc57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016107bf565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610859573d6000803e3d6000fd5b5050506040513d602081101561086e57600080fd5b5051905060006002806108846040848a8c6116fe565b6040516020018083838082843780830192505050925050506040516020818303038152906040526040518082805190602001908083835b602083106108f857805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016108bb565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610955573d6000803e3d6000fd5b5050506040513d602081101561096a57600080fd5b5051600261097b896040818d6116fe565b60405160009060200180848480828437919091019283525050604080518083038152602092830191829052805190945090925082918401908083835b602083106109f457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016109b7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610a51573d6000803e3d6000fd5b5050506040513d6020811015610a6657600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610ada57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610a9d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610b37573d6000803e3d6000fd5b5050506040513d6020811015610b4c57600080fd5b50516040805160208101858152929350600092600292839287928f928f92018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b60208310610bd957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b9c565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610c36573d6000803e3d6000fd5b5050506040513d6020811015610c4b57600080fd5b50516040518651600291889160009188916020918201918291908601908083835b60208310610ca957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610c6c565b6001836020036101000a0380198251168184511680821785525050505050509050018367ffffffffffffffff191667ffffffffffffffff1916815260180182815260200193505050506040516020818303038152906040526040518082805190602001908083835b60208310610d4e57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610d11565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610dab573d6000803e3d6000fd5b5050506040513d6020811015610dc057600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610e3457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610df7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610e91573d6000803e3d6000fd5b5050506040513d6020811015610ea657600080fd5b50519050858114610f02576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260548152602001806117486054913960600191505060405180910390fd5b60205463ffffffff11610f60576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806117276021913960400191505060405180910390fd5b602080546001019081905560005b60208110156110a9578160011660011415610fa0578260008260208110610f9157fe5b0155506110ac95505050505050565b600260008260208110610faf57fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061102557805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610fe8565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015611082573d6000803e3d6000fd5b5050506040513d602081101561109757600080fd5b50519250600282049150600101610f6e565b50fe5b50505050505050565b60606110c26020546114ba565b905090565b6020546000908190815b60208110156112f05781600116600114156111e6576002600082602081106110f557fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061116b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161112e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156111c8573d6000803e3d6000fd5b5050506040513d60208110156111dd57600080fd5b505192506112e2565b600283602183602081106111f657fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061126b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161122e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156112c8573d6000803e3d6000fd5b5050506040513d60208110156112dd57600080fd5b505192505b6002820491506001016110d1565b506002826112ff6020546114ba565b600060401b6040516020018084815260200183805190602001908083835b6020831061135a57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161131d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790527fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000095909516920191825250604080518083037ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8018152601890920190819052815191955093508392850191508083835b6020831061143f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101611402565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa15801561149c573d6000803e3d6000fd5b5050506040513d60208110156114b157600080fd5b50519250505090565b60408051600880825281830190925260609160208201818036833701905050905060c082901b8060071a60f81b826000815181106114f457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060061a60f81b8260018151811061153757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060051a60f81b8260028151811061157a57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060041a60f81b826003815181106115bd57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060031a60f81b8260048151811061160057fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060021a60f81b8260058151811061164357fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060011a60f81b8260068151811061168657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060001a60f81b826007815181106116c957fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050919050565b6000808585111561170d578182fd5b83861115611719578182fd5b505082019391909203915056fe4465706f736974436f6e74726163743a206d65726b6c6520747265652066756c6c4465706f736974436f6e74726163743a207265636f6e7374727563746564204465706f7369744461746120646f6573206e6f74206d6174636820737570706c696564206465706f7369745f646174615f726f6f744465706f736974436f6e74726163743a20696e76616c6964207769746864726177616c5f63726564656e7469616c73206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c7565206e6f74206d756c7469706c65206f6620677765694465706f736974436f6e74726163743a20696e76616c6964207075626b6579206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f20686967684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f206c6f774465706f736974436f6e74726163743a20696e76616c6964207369676e6174757265206c656e677468a26469706673582212201dd26f37a621703009abf16e77e69c93dc50c79db7f6cc37543e3e0e3decdc9764736f6c634300060b0033",
		  },		
	}
  }
`
