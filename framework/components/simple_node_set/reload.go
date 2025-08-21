package simple_node_set

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/chaos"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

// UpgradeNodeSet updates nodes configuration TOML files
// this API is discouraged, however, you can use it if nodes require restart or configuration updates, temporarily!
func UpgradeNodeSet(t *testing.T, in *Input, bc *blockchain.Output, wait time.Duration) (*Output, error) {
	uniq := fmt.Sprintf("%s-%s-%s", framework.DefaultCTFLogsDir, t.Name(), uuid.NewString()[0:4])
	if _, err := framework.SaveContainerLogs(uniq); err != nil {
		return nil, err
	}
	_, err := chaos.ExecPumba(fmt.Sprintf("rm --volumes=false re2:^%s-node.*|%s-ns-postgresql.*", in.Name, in.Name), wait)
	if err != nil {
		return nil, err
	}
	in.Out = nil
	out, err := NewSharedDBNodeSet(in, bc)
	in.Out = out
	return out, err
}
