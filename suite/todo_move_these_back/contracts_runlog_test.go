package contracts

import (
	"context"
	"errors"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"math/big"
	"strings"
)

var _ = Describe("Direct request suite", func() {
	var (
		s             *actions.DefaultSuiteSetup
		adapter       environment.ExternalAdapter
		nodes         []client.Chainlink
		nodeAddresses []common.Address
		oracle        contracts.Oracle
		consumer      contracts.APIConsumer
		jobUUID       uuid.UUID
		err           error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(1),
				client.NewNetworkFromConfig,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Funding Chainlink nodes", func() {
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(nodes, s.Client, s.Wallets.Default(), big.NewInt(2e18), nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying and funding the contracts", func() {
			oracle, err = s.Deployer.DeployOracle(s.Wallets.Default(), s.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = s.Deployer.DeployAPIConsumer(s.Wallets.Default(), s.Link.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(s.Wallets.Default(), nil, big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Permitting node to fulfill request", func() {
			err = oracle.SetFulfillmentPermission(s.Wallets.Default(), nodeAddresses[0].Hex(), true)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Creating directrequest job", func() {
			jobUUID = uuid.NewV4()
			os := &client.DirectRequestTxPipelineSpec{
				URL:         adapter.ClusterURL() + "/five",
				Method:      "GET",
				RequestData: "{}",
				DataPath:    "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())
			_, err = nodes[0].CreateJob(&client.DirectRequestJobSpec{
				Name:              "direct_request",
				ContractAddress:   oracle.Address(),
				ExternalJobID:     jobUUID.String(),
				ObservationSource: ost,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Calling oracle contract", func() {
			jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
			Expect(err).ShouldNot(HaveOccurred())
			var jobID [32]byte
			copy(jobID[:], jobUUIDReplaces)
			err = consumer.CreateRequestTo(
				s.Wallets.Default(),
				oracle.Address(),
				jobID,
				big.NewInt(1e18),
				adapter.ClusterURL()+"/five",
				"data,result",
				big.NewInt(100),
			)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with DirectRequest job", func() {
		It("receives API call data on-chain", func() {
			err = retry.Do(func() error {
				d, err := consumer.Data(context.Background())
				if d == nil {
					return errors.New("no data")
				}
				log.Debug().Int64("Data", d.Int64()).Send()
				if d.Int64() != 5 {
					return errors.New("data is not on chain")
				}
				if err != nil {
					return err
				}
				return nil
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", s.TearDown())
	})
})
