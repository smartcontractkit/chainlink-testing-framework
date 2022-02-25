package smoke

//revive:disable:dot-imports
import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("Keeper suite @keeper", func() {
	var (
		err              error
		networks         *client.Networks
		contractDeployer contracts.ContractDeployer
		registry         contracts.KeeperRegistry
		consumer         contracts.KeeperConsumer
		checkGasLimit    = uint32(2500000)
		linkToken        contracts.LinkToken
		chainlinkNodes   []client.Chainlink
		nodeAddresses    []common.Address
		env              *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(environment.ChainlinkReplicas(6, nil), ""),
				tools.ChartsRoot,
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = env.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Connecting to launched resources", func() {
			networkRegistry := client.NewDefaultNetworkRegistry()
			networks, err = networkRegistry.GetNetworks(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
			contractDeployer, err = contracts.NewContractDeployer(networks.Default)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
			chainlinkNodes, err = client.ConnectChainlinkNodes(env)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			nodeAddresses, err = actions.ChainlinkNodeAddresses(chainlinkNodes)
			Expect(err).ShouldNot(HaveOccurred(), "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
			networks.Default.ParallelTransactions(true)
		})

		By("Funding Chainlink nodes", func() {
			txCost, err := networks.Default.EstimateCostForChainlinkOperations(10)
			Expect(err).ShouldNot(HaveOccurred(), "Estimating cost for Chainlink Operations shouldn't fail")
			err = actions.FundChainlinkNodes(chainlinkNodes, networks.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred(), "Funding Chainlink nodes shouldn't fail")
			// Edge case where simulated networks need some funds at the 0x0 address in order for keeper reads to work
			if networks.Default.GetNetworkType() == "eth_simulated" {
				err = actions.FundAddresses(networks.Default, big.NewFloat(1), "0x0")
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		By("Deploying Keeper contracts", func() {
			linkToken, err = contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")
			ef, err := contractDeployer.DeployMockETHLINKFeed(big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred(), "Deploying mock ETH-Link feed shouldn't fail")
			gf, err := contractDeployer.DeployMockGasFeed(big.NewInt(2e11))
			Expect(err).ShouldNot(HaveOccurred(), "Deploying mock gas feed shouldn't fail")
			registry, err = contractDeployer.DeployKeeperRegistry(
				&contracts.KeeperRegistryOpts{
					LinkAddr:             linkToken.Address(),
					ETHFeedAddr:          ef.Address(),
					GasFeedAddr:          gf.Address(),
					PaymentPremiumPPB:    uint32(200000000),
					BlockCountPerTurn:    big.NewInt(3),
					CheckGasLimit:        checkGasLimit,
					StalenessSeconds:     big.NewInt(90000),
					GasCeilingMultiplier: uint16(1),
					FallbackGasPrice:     big.NewInt(2e11),
					FallbackLinkPrice:    big.NewInt(2e18),
				},
			)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper registry shouldn't fail")
			err = linkToken.Transfer(registry.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred(), "Funding keeper registry contract shouldn't fail")
			consumer, err = contractDeployer.DeployKeeperConsumer(big.NewInt(5))
			Expect(err).ShouldNot(HaveOccurred(), "Deploying keeper consumer shouldn't fail")
			err = linkToken.Transfer(consumer.Address(), big.NewInt(1e18))
			Expect(err).ShouldNot(HaveOccurred(), "Funding keeper consumer contract shouldn't fail")
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
		})

		By("Registering upkeep target", func() {
			registrar, err := contractDeployer.DeployUpkeepRegistrationRequests(
				linkToken.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying UpkeepRegistrationRequests contract shouldn't fail")
			err = registry.SetRegistrar(registrar.Address())
			Expect(err).ShouldNot(HaveOccurred(), "Registering the registrar address on the registry shouldn't fail")
			err = registrar.SetRegistrarConfig(
				true,
				uint32(999),
				uint16(999),
				registry.Address(),
				big.NewInt(0),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Setting the registrar configuration shouldn't fail")
			req, err := registrar.EncodeRegisterRequest(
				"upkeep_1",
				[]byte("0x1234"),
				consumer.Address(),
				checkGasLimit,
				consumer.Address(),
				[]byte("0x"),
				big.NewInt(9e18),
				0,
			)
			Expect(err).ShouldNot(HaveOccurred(), "Encoding the register request shouldn't fail")
			err = linkToken.TransferAndCall(registrar.Address(), big.NewInt(9e18), req)
			Expect(err).ShouldNot(HaveOccurred(), "Funding registrar with LINK shouldn't fail")
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
		})

		By("Adding Keepers and a job", func() {
			primaryNode := chainlinkNodes[0]
			primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred(), "Reading ETH Keys from Chainlink Client shouldn't fail")
			nodeAddressesStr := make([]string, 0)
			for _, cla := range nodeAddresses {
				nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
			}
			payees := []string{
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
				consumer.Address(),
			}
			err = registry.SetKeepers(nodeAddressesStr, payees)
			Expect(err).ShouldNot(HaveOccurred(), "Setting keepers in the registry shouldn't fail")
			_, err = primaryNode.CreateJob(&client.KeeperJobSpec{
				Name:                     "keeper-test-job",
				ContractAddress:          registry.Address(),
				FromAddress:              primaryNodeAddress,
				MinIncomingConfirmations: 1,
				ObservationSource:        client.ObservationSourceKeeperDefault(),
			})
			Expect(err).ShouldNot(HaveOccurred(), "Creating KeeperV2 Job shouldn't fail")
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for event subscriptions in nodes shouldn't fail")
		})
	})

	Describe("with Keeper job", func() {
		It("performs upkeep of a target contract", func() {
			Eventually(func(g Gomega) {
				cnt, err := consumer.Counter(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred(), "Calling consumer's Counter shouldn't fail")
				g.Expect(cnt.Int64()).Should(BeNumerically(">", int64(0)), "Expected consumer counter to be greater than 0, but got %d", cnt.Int64())
				log.Info().Int64("Upkeep counter", cnt.Int64()).Msg("Upkeeps performed")
			}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			networks.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(env, networks, utils.ProjectRoot, nil)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
