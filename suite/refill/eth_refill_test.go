package refill

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

var _ = Describe("ETH refill suite", func() {
	DescribeTable("Can work after refill", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		s, err := suite.DefaultLocalSetup(initFunc)
		Expect(err).ShouldNot(HaveOccurred())
		fluxInstance, err := s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(s.Wallets.Default(), big.NewInt(0), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())
		clNodes, nodeAddrs, err := suite.ConnectToTemplateNodes()
		oraclesAtTest := nodeAddrs[:3]
		clNodesAtTest := clNodes[:3]
		Expect(err).ShouldNot(HaveOccurred())

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

		adapter := tools.NewExternalAdapter()
		os := &client.PipelineSpec{
			URL:         adapter.InsideDockerAddr + "/five",
			Method:      "POST",
			RequestData: "{}",
			DataPath:    "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		for _, n := range clNodesAtTest {
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
		// fund enough for one round
		err = suite.FundTemplateNodes(s.Client, s.Wallets, clNodes, 1e16, 0)
		Expect(err).ShouldNot(HaveOccurred())
		_, _ = tools.SetVariableMockData(adapter.LocalAddr, 6)
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		// no ETH
		_, _ = tools.SetVariableMockData(adapter.LocalAddr, 5)
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).Should(HaveOccurred())
		// refill and check if it works
		err = suite.FundTemplateNodes(s.Client, s.Wallets, clNodes, 2e18, 0)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		data, err := fluxInstance.GetContractData(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(5)))
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)
})
