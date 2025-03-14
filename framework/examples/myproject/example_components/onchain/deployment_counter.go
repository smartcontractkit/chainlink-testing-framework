package onchain

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/examples/example_components/onchain/gethwrappers"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func NewCounterDeployment(c *seth.Client, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	counterABI, err := gethwrappers.GethwrappersMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	dd, err := c.DeployContract(c.NewTXOpts(),
		"TestCounter",
		*counterABI,
		common.FromHex(gethwrappers.GethwrappersMetaData.Bin),
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
