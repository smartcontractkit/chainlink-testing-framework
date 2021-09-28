package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/tools"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Flux monitor suite @flux", func() {
	var (
		s             *actions.DefaultSuiteSetup
		adapter       environment.ExternalAdapter
		nodes         []client.Chainlink
		nodeAddresses []common.Address
		fluxInstance  contracts.FluxAggregator
		err           error
	)
	fluxRoundTimeout := time.Minute * 2

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				"basic-chainlink",
				environment.NewChainlinkCluster(3),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())

			s.Client.ParallelTransactions(true)
		})

		By("Deploying and funding contract", func() {
			fluxInstance, err = s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), contracts.DefaultFluxAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.Fund(s.Wallets.Default(), nil, big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Funding Chainlink nodes", func() {
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())
			ethAmount, err := s.Deployer.CalculateETHForTXs(s.Wallets.Default(), s.Network.Config(), 3)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(
				nodes,
				s.Client,
				s.Wallets.Default(),
				ethAmount,
				nil,
			)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting oracle options", func() {
			err = fluxInstance.SetOracles(s.Wallets.Default(),
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				})
			Expect(err).ShouldNot(HaveOccurred())
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})

		By("Creating flux jobs", func() {
			bta := client.BridgeTypeAttributes{
				Name: "variable",
				URL:  fmt.Sprintf("%s/variable", adapter.ClusterURL()),
			}
			for _, n := range nodes {
				err = n.CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				fluxSpec := &client.FluxMonitorJobSpec{
					Name:              "flux_monitor",
					ContractAddress:   fluxInstance.Address(),
					PollTimerPeriod:   15 * time.Second, // min 15s
					PollTimerDisabled: false,
					ObservationSource: client.ObservationSourceSpecBridge(bta),
				}
				_, err = n.CreateJob(fluxSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	})

	Describe("with Flux job", func() {
		It("performs two rounds and has withdrawable payments for oracles", func() {
			err = adapter.SetVariable(1e7)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
			s.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Interface("data", data).Msg("Round data")
			Expect(len(data.Oracles)).Should(Equal(3))
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(1e7)))
			Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(2)))
			Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(2)))
			Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999994)))
			Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(6)))

			err = adapter.SetVariable(1e8)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(3), fluxRoundTimeout)
			s.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = s.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			data, err = fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(data.Oracles)).Should(Equal(3))
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(1e8)))
			Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(3)))
			Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(3)))
			Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999991)))
			Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(9)))
			log.Info().Interface("data", data).Msg("Round data")

			for _, oracleAddr := range nodeAddresses {
				payment, _ := fluxInstance.WithdrawablePayment(context.Background(), oracleAddr)
				Expect(payment.Int64()).Should(Equal(int64(3)))
			}
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			s.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", s.TearDown())
	})
})
