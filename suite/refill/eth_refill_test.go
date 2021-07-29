package refill

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

var _ = Describe("ETH refill suite", func() {
	var s *actions.DefaultSuiteSetup
	var err error

	DescribeTable("Can work after refill", func(
		envInitFunc environment.K8sEnvSpecInit,
		networkInitFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		s, err = actions.DefaultLocalSetup(envInitFunc, networkInitFunc)
		Expect(err).ShouldNot(HaveOccurred())
		fluxInstance, err := s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(s.Wallets.Default(), big.NewInt(0), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.UpdateAvailableFunds(context.Background(), s.Wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())
		clNodes, err := environment.GetChainlinkClients(s.Env)
		Expect(err).ShouldNot(HaveOccurred())
		nodeAddrs, err := actions.ChainlinkNodeAddresses(clNodes)
		Expect(err).ShouldNot(HaveOccurred())

		err = fluxInstance.SetOracles(s.Wallets.Default(),
			contracts.SetOraclesOptions{
				AddList:            nodeAddrs,
				RemoveList:         []common.Address{},
				AdminList:          nodeAddrs,
				MinSubmissions:     3,
				MaxSubmissions:     3,
				RestartDelayRounds: 0,
			})
		Expect(err).ShouldNot(HaveOccurred())
		oracles, err := fluxInstance.GetOracles(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("oracles set")

		adapter, err := environment.GetExternalAdapter(s.Env)
		Expect(err).ShouldNot(HaveOccurred())

		os := &client.PipelineSpec{
			URL:         adapter.ClusterURL() + "/five",
			Method:      "POST",
			RequestData: "{}",
			DataPath:    "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		for _, n := range clNodes {
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
		err = actions.FundChainlinkNodes(clNodes, s.Client, s.Wallets.Default(), big.NewInt(1e16), nil)
		Expect(err).ShouldNot(HaveOccurred())
		err = adapter.SetVariable(6)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		// no ETH
		err = adapter.SetVariable(5)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).Should(HaveOccurred())
		// refill and check if it works
		err = actions.FundChainlinkNodes(clNodes, s.Client, s.Wallets.Default(), big.NewInt(2e18), nil)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.AwaitNextRoundFinalized(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		data, err := fluxInstance.GetContractData(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(data.LatestRoundData.Answer.Int64()).Should(Equal(int64(5)))
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)

	AfterEach(func() {
		s.Env.TearDown()
	})
})
