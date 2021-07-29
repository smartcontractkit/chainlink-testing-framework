package refill

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"math/big"
	"strings"
	"time"
)

var _ = Describe("FluxAggregator ETH Refill", func() {
	var s *actions.DefaultSuiteSetup
	var adapter environment.ExternalAdapter
	var nodes []client.Chainlink
	var nodeAddresses []common.Address
	var err error
	var fluxInstance contracts.FluxAggregator

	BeforeSuite(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster("../../", 3),
				client.NewHardhatNetwork,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	BeforeEach(func() {
		By("Deploying and funding the contract", func() {
			fluxInstance, err = s.Deployer.DeployFluxAggregatorContract(
				s.Wallets.Default(),
				contracts.DefaultFluxAggregatorOptions(),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.Fund(s.Wallets.Default(), big.NewInt(0), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting FluxAggregator options", func() {
			err = fluxInstance.SetOracles(s.Wallets.Default(),
				contracts.SetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				})
			Expect(err).ShouldNot(HaveOccurred())
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("oracles set")
		})

		By("Adding FluxAggregator jobs to nodes", func() {
			os := &client.PipelineSpec{
				URL:         adapter.ClusterURL() + "/five",
				Method:      "POST",
				RequestData: "{}",
				DataPath:    "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())
			for _, n := range nodes {
				fluxSpec := &client.FluxMonitorJobSpec{
					Name:              "flux_monitor",
					ContractAddress:   fluxInstance.Address(),
					PollTimerPeriod:   15 * time.Second, // min 15s
					PollTimerDisabled: false,
					ObservationSource: ost,
				}
				_, err = n.CreateJob(fluxSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Funding ETH for a single round", func() {
			err = actions.FundChainlinkNodes(nodes, s.Client, s.Wallets.Default(), big.NewInt(1e16), nil)
			Expect(err).ShouldNot(HaveOccurred())
			err = adapter.SetVariable(6)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.AwaitNextRoundFinalized(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Draining ETH on the nodes", func() {
			err = adapter.SetVariable(5)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.AwaitNextRoundFinalized(context.Background())
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("with FluxAggregator", func() {
		It("should refill and await the next round", func() {
			err = actions.FundChainlinkNodes(nodes, s.Client, s.Wallets.Default(), big.NewInt(2e18), nil)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.AwaitNextRoundFinalized(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			data, err := fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(5)))
		})
	})

	AfterSuite(func() {
		By("Tearing down the environment", func() {
			s.Env.TearDown()
		})
	})
})
