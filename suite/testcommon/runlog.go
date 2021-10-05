package testcommon

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

// RunlogSetupInputs inputs needed for a runlog test
type RunlogSetupInputs struct {
	S             *actions.DefaultSuiteSetup
	Adapter       environment.ExternalAdapter
	Nodes         []client.Chainlink
	NodeAddresses []common.Address
	Oracle        contracts.Oracle
	Consumer      contracts.APIConsumer
	JobUUID       uuid.UUID
	Err           error
}

// SetupRunlogEnv does all the environment setup for a run log type test
func SetupRunlogEnv(i *RunlogSetupInputs) {
	By("Deploying the environment", func() {
		i.S, i.Err = actions.DefaultLocalSetup(
			"basic-chainlink",
			environment.NewChainlinkCluster(1),
			client.NewNetworkFromConfig,
			tools.ProjectRoot,
		)
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.Adapter, i.Err = environment.GetExternalAdapter(i.S.Env)
		Expect(i.Err).ShouldNot(HaveOccurred())
	})
}

// SetupRunlogTest does all other test preparations for runlog
func SetupRunlogTest(i *RunlogSetupInputs) {
	By("Funding Chainlink nodes", func() {
		i.Nodes, i.Err = environment.GetChainlinkClients(i.S.Env)
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.NodeAddresses, i.Err = actions.ChainlinkNodeAddresses(i.Nodes)
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.Err = actions.FundChainlinkNodes(i.Nodes, i.S.Client, i.S.Wallets.Default(), big.NewFloat(2), nil)
		Expect(i.Err).ShouldNot(HaveOccurred())
	})
	By("Deploying and funding the contracts", func() {
		i.Oracle, i.Err = i.S.Deployer.DeployOracle(i.S.Wallets.Default(), i.S.Link.Address())
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.Consumer, i.Err = i.S.Deployer.DeployAPIConsumer(i.S.Wallets.Default(), i.S.Link.Address())
		Expect(i.Err).ShouldNot(HaveOccurred())
		i.Err = i.Consumer.Fund(i.S.Wallets.Default(), nil, big.NewFloat(2))
		Expect(i.Err).ShouldNot(HaveOccurred())
	})
	By("Permitting node to fulfill request", func() {
		i.Err = i.Oracle.SetFulfillmentPermission(i.S.Wallets.Default(), i.NodeAddresses[0].Hex(), true)
		Expect(i.Err).ShouldNot(HaveOccurred())
	})
	By("Creating directrequest job", func() {
		i.JobUUID = uuid.NewV4()

		bta := client.BridgeTypeAttributes{
			Name: "five",
			URL:  fmt.Sprintf("%s/five", i.Adapter.ClusterURL()),
		}
		i.Err = i.Nodes[0].CreateBridge(&bta)
		Expect(i.Err).ShouldNot(HaveOccurred())

		os := &client.DirectRequestTxPipelineSpec{
			BridgeTypeAttributes: bta,
			DataPath:             "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())

		_, err = i.Nodes[0].CreateJob(&client.DirectRequestJobSpec{
			Name:              "direct_request",
			ContractAddress:   i.Oracle.Address(),
			ExternalJobID:     i.JobUUID.String(),
			ObservationSource: ost,
		})
		Expect(err).ShouldNot(HaveOccurred())
	})
}

// CallRunlogOracle calls runlog oracle
func CallRunlogOracle(i *RunlogSetupInputs) {
	By("Calling oracle contract", func() {
		jobUUIDReplaces := strings.Replace(i.JobUUID.String(), "-", "", 4)
		Expect(i.Err).ShouldNot(HaveOccurred())
		var jobID [32]byte
		copy(jobID[:], jobUUIDReplaces)
		i.Err = i.Consumer.CreateRequestTo(
			i.S.Wallets.Default(),
			i.Oracle.Address(),
			jobID,
			big.NewInt(1e18),
			i.Adapter.ClusterURL()+"/five",
			"data,result",
			big.NewInt(100),
		)
		Expect(i.Err).ShouldNot(HaveOccurred())
	})
}

// CheckRunlogCompleted checks if oracle send the data on chain
func CheckRunlogCompleted(i *RunlogSetupInputs) {
	By("receives API call data on-chain", func() {
		Eventually(func(g Gomega){
			d, err := i.Consumer.Data(context.Background())
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(d).ShouldNot(BeNil())
			log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
			g.Expect(d.Int64()).Should(BeNumerically("==", 5))
		}, "2m", "1s").Should(Succeed())
	})
}
