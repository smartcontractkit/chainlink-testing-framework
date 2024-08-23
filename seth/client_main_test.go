package seth_test

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"
	network_debug_contract "github.com/smartcontractkit/seth/contracts/bind/debug"
	network_sub_contract "github.com/smartcontractkit/seth/contracts/bind/sub"
	"github.com/stretchr/testify/require"

	link_token "github.com/smartcontractkit/seth/contracts/bind/link"
)

/*
	Some tests should be run on testnets/mainnets, so we are deploying the contract only once,
	for these types of tests it's always a choice between funds/speed of tests

	If you need unique setup, just use NewDebugContractSetup in tests
*/

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

func newClient(t *testing.T) *seth.Client {
	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initialize seth")

	return c
}

func newClientWithEphemeralAddresses(t *testing.T) *seth.Client {
	cfg, err := seth.ReadConfig()
	require.NoError(t, err, "failed to read config")

	var sixty int64 = 60
	cfg.EphemeralAddrs = &sixty

	c, err := seth.NewClientWithConfig(cfg)
	require.NoError(t, err, "failed to initialize seth")

	return c
}

func TestDeploymentLinkTokenFromGethWrapperExample(t *testing.T) {
	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initialize seth")
	abi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")
	contractData, err := c.DeployContract(c.NewTXOpts(), "LinkToken", *abi, []byte(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy link token contract from wrapper's ABI/BIN")

	contract, err := link_token.NewLinkToken(contractData.Address, c.Client)
	require.NoError(t, err, "failed to create debug contract instance")

	_, err = c.Decode(contract.Mint(c.NewTXOpts(), common.Address{}, big.NewInt(1)))
	require.NoError(t, err, "failed to decode transaction")
}

func TestDeploymentAbortedWhenContextHasError(t *testing.T) {
	c, err := seth.NewClient()
	require.NoError(t, err, "failed to initialize seth")
	abi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get ABI")

	opts := c.NewTXOpts()
	opts.Context = context.WithValue(context.Background(), seth.ContextErrorKey{}, errors.New("context error"))

	_, err = c.DeployContract(opts, "LinkToken", *abi, []byte(link_token.LinkTokenMetaData.Bin))
	require.Error(t, err, "did not abort deployment of link token contract due to context error")
	require.Contains(t, err.Error(), "aborted contract deployment for", "incorrect context error")
}

func newClientWithContractMapFromEnv(t *testing.T) *seth.Client {
	c := newClient(t)
	if TestEnv.ContractMap.Size() == 0 {
		t.Fatal("contract map is empty")
	}

	// create a copy of the map, so we don't have problem with side effects of modifying client's map
	// impacting the global, underlying one
	contractMap := seth.NewEmptyContractMap()
	for k, v := range TestEnv.ContractMap.GetContractMap() {
		contractMap.AddContract(k, v)
	}

	c.ContractAddressToNameMap = contractMap

	// now let's recreate the Tracer, so that it has the same contract map
	tracer, err := seth.NewTracer(c.ContractStore, c.ABIFinder, c.Cfg, contractMap, c.Addresses)
	require.NoError(t, err, "failed to create tracer")

	c.Tracer = tracer
	c.ABIFinder.ContractMap = contractMap

	return c
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
	cs, err := seth.NewContractStore("./contracts/abi", "./contracts/bin")
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
	if skip := os.Getenv("SKIP_MAIN_CONFIG"); skip == "" {
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
	} else {
		seth.L.Warn().Msg("Skipping main suite setup")
	}

	os.Exit(m.Run())
}
