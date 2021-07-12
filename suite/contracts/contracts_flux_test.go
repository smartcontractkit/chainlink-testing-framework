package suite

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/suite"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"strings"
	"time"
)

var _ = Describe("Flux monitor suite", func() {
	DescribeTable("Answering to deviation in rounds", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		s, err := suite.DefaultLocalSetup(initFunc)
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy FluxMonitor contract
		fluxInstance, err := s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(s.Wallets.Default(), big.NewInt(0), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// get nodes and their addresses
		clNodes, nodeAddrs, err := suite.ConnectToTemplateNodes()
		oraclesAtTest := nodeAddrs[:3]
		clNodesAtTest := clNodes[:3]
		Expect(err).ShouldNot(HaveOccurred())
		err = suite.FundTemplateNodes(s.Client, s.Wallets, clNodes, 2e18, 0)
		Expect(err).ShouldNot(HaveOccurred())

		// set oracles and submissions
		err = fluxInstance.SetOracles(s.Wallets.Default(),
			contracts.SetOraclesOptions{
				AddList:            oraclesAtTest,
				RemoveList:         []common.Address{},
				AdminList:          oraclesAtTest,
				MinSubmissions:     3,
				MaxSubmissions:     3,
				RestartDelayRounds: 0,
			})
		Expect(err).ShouldNot(HaveOccurred())
		oracles, err := fluxInstance.GetOracles(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("oracles set")

		// set variable adapter
		adapter := tools.NewExternalAdapter()

		// Send Flux job to chainlink nodes
		for _, n := range clNodesAtTest {
			fluxSpec := &client.FluxMonitorJobSpec{
				Name:              "flux_monitor",
				ContractAddress:   fluxInstance.Address(),
				PollTimerPeriod:   15 * time.Second, // min 15s
				PollTimerDisabled: false,
				ObservationSource: client.ObservationSourceSpec(adapter.InsideDockerAddr + "/variable"),
			}
			_, err = n.CreateJob(fluxSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}
		// first change
		_, _ = tools.SetVariableMockData(adapter.LocalAddr, 5)
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		{
			data, err := fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Interface("data", data).Msg("round data")
			Expect(len(data.Oracles)).Should(Equal(3))
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(5)))
			Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(1)))
			Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(1)))
			Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999997)))
			Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(3)))
		}
		// second change + 20%
		_, _ = tools.SetVariableMockData(adapter.LocalAddr, 6)
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		{
			data, err := fluxInstance.GetContractData(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(data.Oracles)).Should(Equal(3))
			Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(6)))
			Expect(data.LatestRoundData.RoundId.Int64()).Should(Equal(int64(2)))
			Expect(data.LatestRoundData.AnsweredInRound.Int64()).Should(Equal(int64(2)))
			Expect(data.AvailableFunds.Int64()).Should(Equal(int64(999999999999999994)))
			Expect(data.AllocatedFunds.Int64()).Should(Equal(int64(6)))
			log.Info().Interface("data", data).Msg("round data")
		}
		// check available payments for oracles
		for _, oracleAddr := range oraclesAtTest {
			payment, _ := fluxInstance.WithdrawablePayment(context.Background(), oracleAddr)
			Expect(payment.Int64()).Should(Equal(int64(2)))
		}
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Check removing/adding oracles, check new rounds is correct", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Check oracle cooldown when add", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Adapter went offline, come online, round data received in suggested round", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Different sources, only one have flux", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Bridge source", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Check withdrawal with respect to RESERVE_ROUNDS", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	DescribeTable("Person other than oracles starting a round", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		// TODO
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)
})
