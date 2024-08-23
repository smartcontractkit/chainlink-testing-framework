package examples_wasp

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"
	network_debug_contract "github.com/smartcontractkit/seth/contracts/bind/debug"
	link_token "github.com/smartcontractkit/seth/contracts/bind/link"
	network_sub_contract "github.com/smartcontractkit/seth/contracts/bind/sub"
)

func init() {
	_ = os.Setenv("SETH_CONFIG_PATH", "seth.toml")
}

var (
	TestEnv TestEnvironment
)

type TestEnvironment struct {
	Client                  *seth.Client
	DebugContract           *network_debug_contract.NetworkDebugContract
	DebugSubContract        *network_sub_contract.NetworkDebugSubContract
	LinkTokenContract       *link_token.LinkToken
	DebugContractAddress    common.Address
	DebugSubContractAddress common.Address
	DebugContractRaw        *bind.BoundContract
	ContractMap             seth.ContractMap
}

func NewDebugContractSetup() (
	*seth.Client,
	*network_debug_contract.NetworkDebugContract,
	common.Address,
	common.Address,
	*bind.BoundContract,
	error,
) {
	cfg, err := seth.ReadConfig()
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	cs, err := seth.NewContractStore(cfg.ABIDir, cfg.BINDir)
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	addrs, pkeys, err := cfg.ParseKeys()
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	contractMap := seth.NewEmptyContractMap()

	abiFinder := seth.NewABIFinder(contractMap, cs)
	tracer, err := seth.NewTracer(cs, &abiFinder, cfg, contractMap, addrs)
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}

	nm, err := seth.NewNonceManager(cfg, addrs, pkeys)
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, errors.Wrap(err, seth.ErrCreateNonceManager)
	}

	c, err := seth.NewClientRaw(cfg, addrs, pkeys, seth.WithContractStore(cs), seth.WithTracer(tracer), seth.WithNonceManager(nm))
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	subData, err := c.DeployContractFromContractStore(c.NewTXOpts(), "NetworkDebugSubContract.abi")
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	data, err := c.DeployContractFromContractStore(c.NewTXOpts(), "NetworkDebugContract.abi", subData.Address)
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	contract, err := network_debug_contract.NewNetworkDebugContract(data.Address, c.Client)
	if err != nil {
		return nil, nil, common.Address{}, common.Address{}, nil, err
	}
	return c, contract, data.Address, subData.Address, data.BoundContract, nil
}

func TestMain(m *testing.M) {
	var err error
	client, debugContract, debugContractAddress, debugSubContractAddress, debugContractRaw, err := NewDebugContractSetup()
	if err != nil {
		panic(err)
	}

	linkTokenAbi, err := link_token.LinkTokenMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	linkDeploymentData, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	if err != nil {
		panic(err)
	}
	linkToken, err := link_token.NewLinkToken(linkDeploymentData.Address, client.Client)
	if err != nil {
		panic(err)
	}
	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	client.ContractStore.AddABI("LinkToken", *linkAbi)

	contractMap := seth.NewEmptyContractMap()
	for k, v := range client.ContractAddressToNameMap.GetContractMap() {
		contractMap.AddContract(k, v)
	}

	TestEnv = TestEnvironment{
		Client:                  client,
		DebugContract:           debugContract,
		LinkTokenContract:       linkToken,
		DebugContractAddress:    debugContractAddress,
		DebugSubContractAddress: debugSubContractAddress,
		DebugContractRaw:        debugContractRaw,
		ContractMap:             contractMap,
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}
