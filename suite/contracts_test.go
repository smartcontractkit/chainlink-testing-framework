package suite

import (
	"context"
	"github.com/smartcontractkit/integrations-framework/actions"
	"math/big"
	"time"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("OCR Feed", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("deploy and use basic functionality", func(
		initFunc client.BlockchainNetworkInit,
		ocrOptions contracts.OffchainOptions,
	) {
		// Setup
		network, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		env, err := environment.NewK8sEnvironment(environment.NewChainlinkCluster("../", 7), conf, network)
		Expect(err).ShouldNot(HaveOccurred())
		defer env.TearDown()

		chainlinkNodes, err := environment.GetChainlinkClients(env)
		Expect(err).ShouldNot(HaveOccurred())
		blockchain, err := environment.NewBlockchainClient(env, network)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := network.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		adapter, err := environment.GetExternalAdapter(env)
		Expect(err).ShouldNot(HaveOccurred())

		// Fund each chainlink node
		err = actions.FundChainlinkNodes(
			chainlinkNodes,
			blockchain,
			wallets.Default(),
			big.NewInt(2^18),
			big.NewInt(2^18),
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy and config OCR contract
		deployer, err := contracts.NewContractDeployer(blockchain)
		Expect(err).ShouldNot(HaveOccurred())

		ocrInstance, err := deployer.DeployOffChainAggregator(wallets.Default(), ocrOptions)
		Expect(err).ShouldNot(HaveOccurred())
		err = ocrInstance.SetConfig(
			wallets.Default(),
			chainlinkNodes,
			contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes)),
		)
		Expect(err).ShouldNot(HaveOccurred())
		err = ocrInstance.Fund(wallets.Default(), big.NewInt(2^18), big.NewInt(2^18))
		Expect(err).ShouldNot(HaveOccurred())

		// Initialize bootstrap node
		bootstrapNode := chainlinkNodes[0]
		bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
		Expect(err).ShouldNot(HaveOccurred())
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			ContractAddress: ocrInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.CreateJob(bootstrapSpec)
		Expect(err).ShouldNot(HaveOccurred())

		// Send OCR job to other nodes
		for index := 1; index < len(chainlinkNodes); index++ {
			nodeP2PIds, err := chainlinkNodes[index].ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := chainlinkNodes[index].PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeys, err := chainlinkNodes[index].ReadOCRKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []string{bootstrapP2PId},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpec(adapter.ClusterURL() + "/five"),
			}
			_, err = chainlinkNodes[index].CreateJob(ocrSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}

		// Request a new round from the OCR
		err = ocrInstance.RequestNewRound(wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// Wait for a round
		for i := 0; i < 30; i++ {
			round, err := ocrInstance.GetLatestRound(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().
				Str("Contract Address", ocrInstance.Address()).
				Str("Answer", round.Answer.String()).
				Str("Round ID", round.RoundId.String()).
				Str("Answered in Round", round.AnsweredInRound.String()).
				Str("Started At", round.StartedAt.String()).
				Str("Updated At", round.UpdatedAt.String()).
				Msg("Latest Round Data")
			if round.RoundId.Cmp(big.NewInt(0)) > 0 {
				break // Break when OCR round processes
			}
			time.Sleep(time.Second)
		}

		// Check answer is as expected
		answer, err := ocrInstance.GetLatestAnswer(context.Background())
		log.Info().Str("Answer", answer.String()).Msg("Final Answer")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(answer.Int64()).Should(Equal(int64(5)))
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultOffChainAggregatorOptions()),
	)
})

// var _ = Describe("Contracts", func() {
// 	var conf *config.Config

// 	BeforeEach(func() {
// 		var err error
// 		conf, err = config.NewWithPath(config.LocalConfig, "../config")
// 		Expect(err).ShouldNot(HaveOccurred())
// 	})

// 	DescribeTable("deploy and interact with the storage contract", func(
// 		initFunc client.BlockchainNetworkInit,
// 		value *big.Int,
// 	) {
// 		network, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		testEnv, err := environment.NewK8sEnvironment("storage-contract", 0, network)
// 		Expect(err).ShouldNot(HaveOccurred())

// 		storeInstance, err := testEnv.ContractDeployer().DeployStorageContract(testEnv.Wallets().Default())
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Interact with contract
// 		err = storeInstance.Set(value)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		val, err := storeInstance.Get(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(val).To(Equal(value))

// 		err = testEnv.TearDown()
// 		Expect(err).ShouldNot(HaveOccurred())
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, big.NewInt(5)),
// 	)

// 	DescribeTable("deploy and interact with the FluxAggregator contract", func(
// 		initFunc client.BlockchainNetworkInit,
// 		fluxOptions contracts.FluxAggregatorOptions,
// 	) {
// 		network, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		testEnv, err := environment.NewK8sEnvironment("flux-aggregator-contract", 0, network)
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Deploy LINK contract
// 		linkInstance, err := testEnv.ContractDeployer().DeployLinkTokenContract(testEnv.Wallets().Default())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		name, err := linkInstance.Name(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(name).To(Equal("ChainLink Token"))

// 		// Deploy FluxMonitor contract
// 		fluxInstance, err := testEnv.ContractDeployer().DeployFluxAggregatorContract(testEnv.Wallets().Default(), fluxOptions)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		err = fluxInstance.Fund(testEnv.Wallets().Default(), big.NewInt(0), big.NewInt(50000000000))
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Interact with contract
// 		desc, err := fluxInstance.Description(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(desc).To(Equal(fluxOptions.Description))

// 		err = testEnv.TearDown()
// 		Expect(err).ShouldNot(HaveOccurred())
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
// 	)

// 	DescribeTable("deploy and interact with the OffChain Aggregator contract", func(
// 		initFunc client.BlockchainNetworkInit,
// 		ocrOptions contracts.OffchainOptions,
// 	) {
// 		network, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		testEnv, err := environment.NewK8sEnvironment("ocr-contract", 0, network)
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Deploy LINK contract
// 		linkInstance, err := testEnv.ContractDeployer().DeployLinkTokenContract(testEnv.Wallets().Default())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		name, err := linkInstance.Name(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(name).To(Equal("ChainLink Token"))

// 		// Deploy Offchain contract
// 		offChainInstance, err := testEnv.ContractDeployer().DeployOffChainAggregator(testEnv.Wallets().Default(), ocrOptions)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		err = offChainInstance.Fund(testEnv.Wallets().Default(), nil, big.NewInt(50000000000))
// 		Expect(err).ShouldNot(HaveOccurred())

// 		err = testEnv.TearDown()
// 		Expect(err).ShouldNot(HaveOccurred())
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultOffChainAggregatorOptions()),
// 	)
// })
