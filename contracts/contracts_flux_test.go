package contracts

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"strings"
	"time"
)

type DefaultSuiteSetup struct {
	Config   *config.Config
	Client   client.BlockchainClient
	Wallets  client.BlockchainWallets
	Deployer ContractDeployer
	Link     LinkToken
}

func DefaultSetup(initFunc client.BlockchainNetworkInit) (*DefaultSuiteSetup, error) {
	conf, err := config.NewWithPath(config.LocalConfig, "../config")
	networkConfig, err := initFunc(conf)
	if err != nil {
		return nil, err
	}
	blockchainClient, err := client.NewBlockchainClient(networkConfig)
	if err != nil {
		return nil, err
	}
	wallets, err := networkConfig.Wallets()
	if err != nil {
		return nil, err
	}
	contractDeployer, err := NewContractDeployer(blockchainClient)
	if err != nil {
		return nil, err
	}
	link, err := contractDeployer.DeployLinkTokenContract(wallets.Default())
	if err != nil {
		return nil, err
	}
	return &DefaultSuiteSetup{
		Config:   conf,
		Client:   blockchainClient,
		Wallets:  wallets,
		Deployer: contractDeployer,
		Link:     link,
	}, nil
}

var _ = Describe("Flux aggregator suite", func() {
	DescribeTable("deploy and interact with the FluxAggregator contract", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions FluxAggregatorOptions,
	) {
		// Setup network and blockchainClient
		s, err := DefaultSetup(initFunc)
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy FluxMonitor contract
		fluxInstance, err := s.Deployer.DeployFluxAggregatorContract(s.Wallets.Default(), fluxOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = fluxInstance.Fund(s.Wallets.Default(), big.NewInt(0), big.NewInt(1e18))
		Expect(err).ShouldNot(HaveOccurred())

		// Update funds
		err = fluxInstance.UpdateAvailableFunds(nil, s.Wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// check funds updated
		payment, err := fluxInstance.PaymentAmount(nil)
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Int("payment", int(payment.Int64())).Msg("payment amount")

		avFunds, err := fluxInstance.AvailableFunds(nil)
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Int("available funds", int(avFunds.Int64())).Msg("funds")

		alFunds, err := fluxInstance.AllocatedFunds(nil)
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Int("allocated funds", int(alFunds.Int64())).Msg("funds")

		clNodes, addrs, err := client.ConnectToTemplateNodes()
		err = client.FundTemplateNodes(s.Client, s.Wallets, clNodes, 2e18, 2e18)
		Expect(err).ShouldNot(HaveOccurred())

		// set oracles and submissions
		err = fluxInstance.SetOracles(s.Wallets.Default(),
			SetOraclesOptions{
				AddList:            addrs[:3],
				RemoveList:         []common.Address{},
				AdminList:          addrs[:3],
				MinSubmissions:     3,
				MaxSubmissions:     3,
				RestartDelayRounds: 0,
			})
		oracles, err := fluxInstance.GetOracles(nil)
		oraclesString := strings.Join(oracles, ",")
		log.Info().Str("Oracles", oraclesString).Msg("oracles set")

		//err = fluxInstance.SetRequesterPermissions(nil, s.Wallets.Default(), common.HexToAddress(s.Wallets.Default().Address()), true, 0)
		//Expect(err).ShouldNot(HaveOccurred())
		//err = fluxInstance.RequestNewRound(nil, s.Wallets.Default())
		//Expect(err).ShouldNot(HaveOccurred())

		go tools.NewExternalAdapter("6644")
		time.Sleep(2 * time.Second)
		_, err = tools.SetVariableMockData(0.05)

		// Send Flux job to other nodes
		for index := 0; index < 3; index++ {
			// TODO: also try with a bridge source
			//bridgeName := fmt.Sprintf("flux-bridge-%d", index)
			//err = clNodes[index].CreateBridge(&client.BridgeTypeAttributes{
			//	Name: bridgeName,
			//	URL:  "http://host.docker.internal:6644/five",
			//})
			//Expect(err).ShouldNot(HaveOccurred())
			observationSource := `fetch    [type=http method=POST url="http://host.docker.internal:6644/variable" requestData="{}"];
			parse    [type=jsonparse path="data,result"];
			fetch -> parse;`
			fluxSpec := &client.FluxMonitorJobSpec{
				Name:              "flux_monitor",
				ContractAddress:   fluxInstance.Address(),
				IdleTimerPeriod:   5 * time.Second,
				IdleTimerDisabled: false,
				PollTimerPeriod:   15 * time.Second, // min 15s
				PollTimerDisabled: false,
				AbsoluteThreshold: float32(0.1),
				ObservationSource: observationSource,
			}
			_, err = clNodes[index].CreateJob(fluxSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}
		time.Sleep(30 * time.Second)

		r2, err := fluxInstance.LatestRound(nil)
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Str("round_id", r2.String()).Msg("latest round")

		data, err := fluxInstance.GetContractData(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().Str("data", data.LatestRoundData.Answer.String()).Msg("Data is here")
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, DefaultFluxAggregatorOptions()),
	)
})
