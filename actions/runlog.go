package actions

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/hooks"
	"github.com/smartcontractkit/integrations-framework/utils"
)

// RunlogSetupInputs inputs needed for a runlog test
type RunlogSetupInputs struct {
	SuiteSetup  SuiteSetup
	NetworkInfo NetworkInfo
	Adapter     environment.ExternalAdapter
	Nodes         []client.Chainlink
	NodeAddresses []common.Address
	Oracle        contracts.Oracle
	Consumer      contracts.APIConsumer
	JobUUID       uuid.UUID
	Err           error
}

// SetupRunlogEnv does all the environment setup for a run log type test
func SetupRunlogEnv(i *RunlogSetupInputs) func() {
	return func() {
		i.SuiteSetup, i.Err = SingleNetworkSetup(
			environment.NewChainlinkCluster(1),
			hooks.EVMNetworkFromConfigHook,
			hooks.EthereumDeployerHook,
			hooks.EthereumClientHook,
			utils.ProjectRoot,
		)
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.Adapter, i.Err = environment.GetExternalAdapter(i.SuiteSetup.Environment())
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.NetworkInfo = i.SuiteSetup.DefaultNetwork()
	}
}
