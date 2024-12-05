package simple_node_set

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"time"
)

// UpgradeNodeSet updates nodes configuration TOML files
// this API is discouraged, however, you can use it if nodes require restart or configuration updates, temporarily!
func UpgradeNodeSet(in *Input, bc *blockchain.Output, wait time.Duration) (*Output, error) {
	uniq := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, uuid.NewString()[0:4])
	if err := framework.WriteAllContainersLogs(uniq); err != nil {
		return nil, err
	}
	_, err := chaos.ExecPumba("rm --volumes=false re2:^node.*|ns-postgresql.*", wait)
	if err != nil {
		return nil, err
	}
	in.Out = nil
	out, err := NewSharedDBNodeSet(in, bc)
	in.Out = out
	return out, err
}
