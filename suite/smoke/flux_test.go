package smoke

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/utils"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

var _ = Describe("Flux monitor suite @flux", func() {
	var (
		err              error
		nets             *client.Networks
		cd               contracts.ContractDeployer
		lt               contracts.LinkToken
		fluxInstance     contracts.FluxAggregator
		cls              []client.Chainlink
		mockserver       *client.MockserverClient
		nodeAddresses    []common.Address
		fluxRoundTimeout = 3 * time.Minute
		e                *environment.Environment
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.NewEnvironmentFromPreset(filepath.Join(utils.PresetRoot, "chainlink-cluster-3"))
			Expect(err).ShouldNot(HaveOccurred())
			err = e.Connect()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			nets, err = client.NewNetworks(e, nil)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			nodeAddresses, err = actions.ChainlinkNodeAddresses(cls)
			Expect(err).ShouldNot(HaveOccurred())
			mockserver, err = client.NewMockServerClientFromEnv(e)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})
		By("Deploying and funding contract", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			fluxInstance, err = cd.DeployFluxAggregatorContract(lt.Address(), contracts.DefaultFluxAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(fluxInstance.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred())
			err = fluxInstance.UpdateAvailableFunds()
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Funding Chainlink nodes", func() {
			err = actions.FundChainlinkNodes(cls, nets.Default, big.NewFloat(1))
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Setting oracle options", func() {
			err = fluxInstance.SetOracles(
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            nodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          nodeAddresses,
					MinSubmissions:     3,
					MaxSubmissions:     3,
					RestartDelayRounds: 0,
				})
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
			oracles, err := fluxInstance.GetOracles(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Oracles", strings.Join(oracles, ",")).Msg("Oracles set")
		})
		By("Creating flux jobs", func() {
			err = mockserver.SetVariable(0)
			Expect(err).ShouldNot(HaveOccurred())

			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("variable-%s", uuid.NewV4().String()),
				URL:  fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
			}
			for _, n := range cls {
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
			err = mockserver.SetVariable(1e7)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound := contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(2), fluxRoundTimeout)
			nets.Default.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = nets.Default.WaitForEvents()
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

			err = mockserver.SetVariable(1e8)
			Expect(err).ShouldNot(HaveOccurred())

			fluxRound = contracts.NewFluxAggregatorRoundConfirmer(fluxInstance, big.NewInt(3), fluxRoundTimeout)
			nets.Default.AddHeaderEventSubscription(fluxInstance.Address(), fluxRound)
			err = nets.Default.WaitForEvents()
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
			nets.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
