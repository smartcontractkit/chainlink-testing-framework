package refill

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("FluxAggregator ETH Refill @refill", func() {
	var (
		suiteSetup    *actions.DefaultSuiteSetup
		adapter       environment.ExternalAdapter
		nodes         []client.Chainlink
		nodeAddresses []common.Address
		err           error
		fluxInstance  contracts.FluxAggregator
	)
	fluxRoundTimeout := 30 * time.Second

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(3),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			Expect(err).ShouldNot(HaveOccurred())

			suiteSetup.Client.ParallelTransactions(true)
		})
	})

	JustBeforeEach(func() {
		By("Deploying and funding the contract", func() {
			fluxInstance, err = suiteSetup.Deployer.DeployFluxAggregatorContract(
				suiteSetup.Wallets.Default(),
				contracts.DefaultFluxAggregatorOptions(),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.Fund(suiteSetup.Wallets.Default(), nil, big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.UpdateAvailableFunds(context.Background(), suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Setting FluxAggregator options", func() {
			err = fluxInstance.SetOracles(suiteSetup.Wallets.Default(),
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				},
			)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})

		By("Adding FluxAggregator jobs to nodes", func() {
			bta := client.BridgeTypeAttributes{
				Name:        "variable",
				URL:         fmt.Sprintf("%s/variable", adapter.ClusterURL()),
				RequestData: "{}",
			}

			os := &client.PipelineSpec{
				BridgeTypeAttributes: bta,
				DataPath:             "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())

			for _, n := range nodes {
				err = n.CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				fluxSpec := &client.FluxMonitorJobSpec{
					Name:              "flux_monitor",
					ContractAddress:   fluxInstance.Address(),
					PollTimerPeriod:   15 * time.Second, // min 15s
					PollTimerDisabled: false,
					ObservationSource: ost,
				}
				_, err := n.CreateJob(fluxSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Funding ETH for a single round", func() {
			submissionGasUsed, err := suiteSetup.Network.FluxMonitorSubmissionGasUsed()
			Expect(err).ShouldNot(HaveOccurred())
			txCost, err := suiteSetup.Client.CalculateTxGas(submissionGasUsed)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(nodes, suiteSetup.Client, suiteSetup.Wallets.Default(), txCost, nil)
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
			err = adapter.SetVariable(6)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(1), fluxRoundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Draining ETH on the nodes", func() {
			err = adapter.SetVariable(5)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err.Error()).Should(ContainSubstring("timeout waiting for flux round to confirm"))
		})
	})

	Describe("with FluxAggregator", func() {
		It("should refill and await the next round", func() {
			err = actions.FundChainlinkNodes(nodes, suiteSetup.Client, suiteSetup.Wallets.Default(), big.NewFloat(2), nil)
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
			suiteSetup.Client.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			data, err := fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(5)))
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
