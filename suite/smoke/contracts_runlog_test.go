package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/utils"
	"path/filepath"
)

var _ = Describe("Direct request suite @runlog", func() {
	var (
		err error
		e   *environment.Environment
	)
	//var (
	//	suiteSetup    actions.SuiteSetup
	//	networkInfo   actions.NetworkInfo
	//	adapter       environment.ExternalAdapter
	//	nodes         []client.Chainlink
	//	nodeAddresses []common.Address
	//	oracle        contracts.Oracle
	//	consumer      contracts.APIConsumer
	//	jobUUID       uuid.UUID
	//	err           error
	//)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.NewEnvironmentFromPreset(filepath.Join(utils.ProjectRoot, "preset", "chainlink-cluster"))
			Expect(err).ShouldNot(HaveOccurred())
			//err = e.Connect()
			//Expect(err).ShouldNot(HaveOccurred())
			//err = e.Teardown()
			//Expect(err).ShouldNot(HaveOccurred())

			//suiteSetup, err = actions.SingleNetworkSetup(
			//	environment.NewChainlinkCluster(1),
			//	hooks.EVMNetworkFromConfigHook,
			//	hooks.EthereumDeployerHook,
			//	hooks.EthereumClientHook,
			//	utils.ProjectRoot,
			//)
			//Expect(err).ShouldNot(HaveOccurred())
			//networkInfo = suiteSetup.DefaultNetwork()
			//adapter, err = environment.GetExternalAdapter(suiteSetup.Environment())
			//Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating networks", func() {
			//nets, err := client.NewNetworks(e,
			//	map[string]client.ExternalClientImplFunc{
			//		"terra": func(networkName string, networkConfig map[string]interface{}, env *environment.Environment) (client.BlockchainClient, error) {
			//			d, err := yaml.Marshal(networkConfig)
			//			if err != nil {
			//				return nil, err
			//			}
			//			var cfg *config.ETHNetwork
			//			if err := yaml.Unmarshal(d, &cfg); err != nil {
			//				return nil, err
			//			}
			//			if !cfg.External {
			//				cfg.URLs = env.Config.NetworksURLs[networkName]
			//			}
			//			cfg.ID = cfg.Type
			//			return client.NewEthereumMultiNodeClient(cfg)
			//		},
			//	},
			//)
			nets, err := client.NewNetworks(e, nil)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err := contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			lt, err := cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Address", lt.Address()).Send()
		})

		By("Funding Chainlink nodes", func() {
			//nodes, err = environment.GetChainlinkClients(suiteSetup.Environment())
			//Expect(err).ShouldNot(HaveOccurred())
			//nodeAddresses, err = actions.ChainlinkNodeAddresses(nodes)
			//Expect(err).ShouldNot(HaveOccurred())
			//ethAmount, err := networkInfo.Deployer.CalculateETHForTXs(networkInfo.Wallets.Default(), networkInfo.NetworkConfig.Config(), 1)
			//Expect(err).ShouldNot(HaveOccurred())
			//err = actions.FundChainlinkNodes(nodes, networkInfo.Client, networkInfo.Wallets.Default(), ethAmount, nil)
			//Expect(err).ShouldNot(HaveOccurred())
		})

		By("Deploying and funding the contracts", func() {
			//oracle, err = networkInfo.Deployer.DeployOracle(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			//Expect(err).ShouldNot(HaveOccurred())
			//consumer, err = networkInfo.Deployer.DeployAPIConsumer(networkInfo.Wallets.Default(), networkInfo.Link.Address())
			//Expect(err).ShouldNot(HaveOccurred())
			//err = consumer.Fund(networkInfo.Wallets.Default(), nil, big.NewFloat(2))
			//Expect(err).ShouldNot(HaveOccurred())
		})

		By("Permitting node to fulfill request", func() {
			//err = oracle.SetFulfillmentPermission(networkInfo.Wallets.Default(), nodeAddresses[0].Hex(), true)
			//Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating directrequest job", func() {
			//jobUUID = uuid.NewV4()
			//
			//bta := client.BridgeTypeAttributes{
			//	Name: "five",
			//	URL:  fmt.Sprintf("%s/five", adapter.ClusterURL()),
			//}
			//err = nodes[0].CreateBridge(&bta)
			//Expect(err).ShouldNot(HaveOccurred())
			//
			//os := &client.DirectRequestTxPipelineSpec{
			//	BridgeTypeAttributes: bta,
			//	DataPath:             "data,result",
			//}
			//ost, err := os.String()
			//Expect(err).ShouldNot(HaveOccurred())
			//
			//_, err = nodes[0].CreateJob(&client.DirectRequestJobSpec{
			//	Name:              "direct_request",
			//	ContractAddress:   oracle.Address(),
			//	ExternalJobID:     jobUUID.String(),
			//	ObservationSource: ost,
			//})
			//Expect(err).ShouldNot(HaveOccurred())
		})

		By("Calling oracle contract", func() {
			//jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
			//var jobID [32]byte
			//copy(jobID[:], jobUUIDReplaces)
			//err = consumer.CreateRequestTo(
			//	networkInfo.Wallets.Default(),
			//	oracle.Address(),
			//	jobID,
			//	big.NewInt(1e18),
			//	fmt.Sprintf("%s/five", adapter.ClusterURL()),
			//	"data,result",
			//	big.NewInt(100),
			//)
			//Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with DirectRequest job", func() {
		It("receives API call data on-chain", func() {
			//Eventually(func(g Gomega) {
			//	d, err := consumer.Data(context.Background())
			//	g.Expect(err).ShouldNot(HaveOccurred())
			//	g.Expect(d).ShouldNot(BeNil())
			//	log.Debug().Int64("Data", d.Int64()).Msg("Found on chain")
			//	g.Expect(d.Int64()).Should(BeNumerically("==", 5))
			//}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Calculating gas costs", func() {
			//networkInfo.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			//err = e.Artifacts.DumpTestResult(filepath.Join(utils.ProjectRoot, "logs", "test_N"), "chainlink")
			//Expect(err).ShouldNot(HaveOccurred())
			//err = e.Teardown()
			//Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
