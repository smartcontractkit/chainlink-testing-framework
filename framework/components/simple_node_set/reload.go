package simple_node_set

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"time"
)

// UpgradeNodeSet updates nodes configuration TOML files
// this API is discouraged, however, you can use it if nodes require restart or configuration updates, temporarily!
func UpgradeNodeSet(in *Input, bc *blockchain.Output, url string, wait time.Duration) (*Output, error) {
	_, err := chaos.ExecPumba("rm --volumes=false re2:node.*|postgresql.*", wait)
	if err != nil {
		return nil, err
	}
	in.Out = nil
	out, err := NewSharedDBNodeSet(in, bc, url)
	in.Out = out
	return out, err
}
