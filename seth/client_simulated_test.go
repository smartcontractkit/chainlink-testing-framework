package seth_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"
	"github.com/stretchr/testify/require"
)

func newClientWithEthClient(ethClient simulated.Client) (*seth.Client, error) {
	cfg, err := seth.ReadConfig()
	if err != nil {
		return nil, err
	}
	cs, err := seth.NewContractStore("./contracts/abi", "./contracts/bin", nil)
	if err != nil {
		return nil, err
	}
	addrs, pkeys, err := cfg.ParseKeys()
	if err != nil {
		return nil, err
	}
	contractMap := seth.NewEmptyContractMap()

	abiFinder := seth.NewABIFinder(contractMap, cs)
	tracer, err := seth.NewTracer(cs, &abiFinder, cfg, contractMap, addrs)
	if err != nil {
		return nil, err
	}

	nm, err := seth.NewNonceManager(cfg, addrs, pkeys)
	if err != nil {
		return nil, errors.Wrap(err, seth.ErrCreateNonceManager)
	}

	c, err := seth.NewClientRaw(cfg, addrs, pkeys, seth.WithContractStore(cs), seth.WithTracer(tracer), seth.WithNonceManager(nm), seth.WithEthClient(ethClient))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func Test_SimulatedBackend(t *testing.T) {
	alloc := map[common.Address]types.Account{
		common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"): {
			Balance: big.NewInt(1000000000000000000), // 1 Ether
		},
	}
	backend := simulated.NewBackend(alloc)
	client, err := newClientWithEthClient(backend.Client())
	require.NoError(t, err)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				backend.Commit()
			case <-ctx.Done():
				backend.Close()
				return
			}
		}
	}()

	linkAbi, err := link_token.LinkTokenMetaData.GetAbi()
	require.NoError(t, err, "failed to get LINK ABI")

	_, err = client.DeployContract(client.NewTXOpts(), "LinkToken", *linkAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	require.NoError(t, err, "failed to deploy LINK contract")
}
