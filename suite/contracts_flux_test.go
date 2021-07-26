package suite

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Flux monitor suite", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("Answering to deviation in rounds", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		network, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		testEnv, err := environment.NewK8sEnvironment("basic-flux-monitor", 5, network)
		Expect(err).ShouldNot(HaveOccurred())
		defaultWallet := testEnv.Wallets().Default()

		// Deploy FluxMonitor contract
		fluxInstance, err := testEnv.ContractDeployer().DeployFluxAggregatorContract(defaultWallet, fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(defaultWallet, big.NewInt(0), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.UpdateAvailableFunds(context.Background(), defaultWallet)
		Expect(err).ShouldNot(HaveOccurred())

		// get nodes and their addresses
		clNodes := testEnv.ChainlinkNodes()
		nodeAddrs, err := testEnv.ChainlinkNodeETHAddresses()
		Expect(err).ShouldNot(HaveOccurred())
		oraclesAtTest := nodeAddrs[:3]
		clNodesAtTest := clNodes[:3]
		Expect(err).ShouldNot(HaveOccurred())
		err = testEnv.FundAllNodes(defaultWallet, big.NewInt(2e18), nil)
		Expect(err).ShouldNot(HaveOccurred())

		// set oracles and submissions
		err = fluxInstance.SetOracles(defaultWallet,
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
		log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")

		// Send Flux job to chainlink nodes
		for _, n := range clNodesAtTest {
			fluxSpec := &client.FluxMonitorJobSpec{
				Name:              "flux_monitor",
				ContractAddress:   fluxInstance.Address(),
				PollTimerPeriod:   15 * time.Second, // min 15s
				PollTimerDisabled: false,
				ObservationSource: client.ObservationSourceSpec(testEnv.Adapter().ClusterURL() + "/variable"),
			}
			_, err = n.CreateJob(fluxSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}
		// first change
		err = testEnv.Adapter().SetVariable(5)
		Expect(err).ShouldNot(HaveOccurred())
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
		err = testEnv.Adapter().SetVariable(6)
		Expect(err).ShouldNot(HaveOccurred())
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

		err = testEnv.TearDown()
		Expect(err).ShouldNot(HaveOccurred())
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
