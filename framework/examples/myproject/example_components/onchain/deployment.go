package onchain

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	testToken "github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/gethwrappers"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

type Input struct {
	URL string  `toml:"url"`
	Out *Output `toml:"out"`
}

type Output struct {
	UseCache  bool             `toml:"use_cache"`
	Addresses []common.Address `toml:"addresses"`
}

func NewProductOnChainDeployment(in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}

	// deploy your contracts here, example

	t := seth.Duration{D: 2 * time.Minute}

	c, err := seth.NewClientBuilder().
		WithRpcUrl(in.URL).
		WithProtections(true, false, &t).
		WithGasPriceEstimations(true, 0, seth.Priority_Fast).
		WithTracing(seth.TracingLevel_All, []string{seth.TraceOutput_Console}).
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	if err != nil {
		return nil, err
	}

	contractABI, err := testToken.BurnMintERC677MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	dd, err := c.DeployContract(c.NewTXOpts(),
		"TestToken",
		*contractABI,
		common.FromHex(testToken.BurnMintERC677MetaData.Bin),
		"TestToken",
		"TestToken",
		uint8(18),
		big.NewInt(1000),
	)
	if err != nil {
		return nil, err
	}

	out := &Output{
		UseCache: true,
		// save all the addresses to output, so it can be cached
		Addresses: []common.Address{dd.Address},
	}
	in.Out = out
	return out, nil
}
