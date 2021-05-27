package contracts

import (
	"context"
	"fmt"
	"integrations-framework/client"
	"integrations-framework/config"
	"integrations-framework/tools"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = Describe("Chainlink Node", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewWithPath(config.LocalConfig, "../config")
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("deploy and use basic functionality", func(
		initFunc client.BlockchainNetworkInit,
	) {
		// Setup
		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())
		// Hardcoded one used by Lorenz's protoype to configure the OCR contract
		ocrContractWallet, err := client.NewEthereumWallet("000000000000000000000000000000000000000000000000015c9a825f7c0878")
		Expect(err).ShouldNot(HaveOccurred())
		ethClient, err := client.NewEthereumClient(networkConfig)
		Expect(err).ShouldNot(HaveOccurred())
		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		extraFundingWallet, err := wallets.Wallet(2)
		Expect(err).ShouldNot(HaveOccurred())
		linkInstance, err := DeployLinkTokenContract(ethClient, wallets.Default())
		Expect(err).ShouldNot(HaveOccurred())

		// Launch Nodes
		chainlinkNodes, err := client.CreateTemplateNodes(ethClient, linkInstance.Address())
		Expect(err).ShouldNot(HaveOccurred())
		for index := range chainlinkNodes {
			err = chainlinkNodes[index].Fund(wallets.Default(), big.NewInt(2000000000000000000), big.NewInt(2000000000000000000))
			Expect(err).ShouldNot(HaveOccurred())
		}
		// Fund the default wallet used by Lorenz's prototype script for configuration
		// TODO: Either recreate the prototype code for this section, or get it in an API form
		txhash, err := ethClient.SendTransaction(extraFundingWallet,
			common.HexToAddress(ocrContractWallet.Address()), big.NewInt(7900000000000000000), nil)
		Expect(err).ShouldNot(HaveOccurred())
		err = ethClient.WaitForTransaction(*txhash)
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy and config OCR contract
		ocrInstance, err := DeployOffChainAggregator(ethClient, ocrContractWallet)
		Expect(err).ShouldNot(HaveOccurred())
		err = ocrInstance.SetConfig(context.Background(), ocrContractWallet, chainlinkNodes)
		Expect(err).ShouldNot(HaveOccurred())

		// Create external adapter
		go tools.NewExternalAdapter("6644")

		// Initialize bootstrap node
		bootstrapP2PIds, err := chainlinkNodes[0].ReadP2PKeys()
		Expect(err).ShouldNot(HaveOccurred())
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := buildBootstrapSpec(ocrInstance.Address(), bootstrapP2PId)
		_, err = chainlinkNodes[0].CreateJob(bootstrapSpec)
		Expect(err).ShouldNot(HaveOccurred())

		// Send OCR job to other nodes
		for index := 1; index < len(chainlinkNodes); index++ {
			nodeP2PIds, err := chainlinkNodes[index].ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddresses, err := chainlinkNodes[index].ReadETHKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeTransmitterAddress := nodeTransmitterAddresses.Data[0].Attributes.Address
			nodeOCRKeys, err := chainlinkNodes[index].ReadOCRKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			ocrSpec := buildOCRSpec(ocrInstance.Address(), nodeP2PId, bootstrapP2PId, nodeOCRKeyId, nodeTransmitterAddress)
			_, err = chainlinkNodes[index].CreateJob(ocrSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}

		// Quickly create a bunch of new blocks
		for {
			time.Sleep(time.Second * 2)
			ethClient.SendTransaction(wallets.Default(), common.HexToAddress(extraFundingWallet.Address()),
				big.NewInt(123456789), nil)
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
		}

		// Cleanup
		// err = client.CleanTemplateNodes()
		// Expect(err).ShouldNot(HaveOccurred())

	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork),
		// Tested locally successfully. We need to implement secrets system as well as testing wallets for CI use
		// Entry("on Ethereum Kovan", NewKovanNetwork),
		// Entry("on Ethereum Goerli", NewGoerliNetwork),
	)

})

// var _ = Describe("Client", func() {
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
// 		// Deploy contract
// 		networkConfig, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		client, err := client.NewEthereumClient(networkConfig)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		wallets, err := networkConfig.Wallets()
// 		Expect(err).ShouldNot(HaveOccurred())
// 		storeInstance, err := DeployStorageContract(client, wallets.Default())
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Interact with contract
// 		err = storeInstance.Set(context.Background(), value)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		val, err := storeInstance.Get(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(val).To(Equal(value))
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, big.NewInt(5)),
// 		// Tested locally successfully. We need to implement secrets system as well as testing wallets for CI use
// 		// Entry("on Ethereum Kovan", client.NewKovanNetwork, big.NewInt(5)),
// 		// Entry("on Ethereum Goerli", client.NewGoerliNetwork, big.NewInt(5)),
// 	)

// 	DescribeTable("deploy and interact with the FluxAggregator contract", func(
// 		initFunc client.BlockchainNetworkInit,
// 		fluxOptions FluxAggregatorOptions,
// 	) {
// 		// Setup network and client
// 		networkConfig, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		client, err := client.NewEthereumClient(networkConfig)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		wallets, err := networkConfig.Wallets()
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Deploy LINK contract
// 		linkInstance, err := DeployLinkTokenContract(client, wallets.Default())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		name, err := linkInstance.Name(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(name).To(Equal("ChainLink Token"))

// 		// Deploy FluxMonitor contract
// 		fluxInstance, err := DeployFluxAggregatorContract(client, wallets.Default(), fluxOptions)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		err = fluxInstance.Fund(wallets.Default(), big.NewInt(0), big.NewInt(50000000000))
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Interact with contract
// 		desc, err := fluxInstance.Description(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(desc).To(Equal(fluxOptions.Description))
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, FluxAggregatorOptions{
// 			PaymentAmount: big.NewInt(1),
// 			Timeout:       uint32(5),
// 			MinSubValue:   big.NewInt(1),
// 			MaxSubValue:   big.NewInt(10),
// 			Decimals:      uint8(8),
// 			Description:   "Hardhat Flux Aggregator",
// 		}),
// 		// Tested locally successfully. We need to implement secrets system as well as testing wallets for CI use
// 		// Entry("on Ethereum Kovan", client.NewKovanNetwork, big.NewInt(5)),
// 		// Entry("on Ethereum Goerli", client.NewGoerliNetwork, big.NewInt(5)),
// 	)

// 	DescribeTable("deploy and interact with the OffChain Aggregator contract", func(
// 		initFunc client.BlockchainNetworkInit,
// 		offchainOptions OffchainOptions,
// 	) {
// 		// Setup network and client
// 		networkConfig, err := initFunc(conf)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		client, err := client.NewEthereumClient(networkConfig)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		wallets, err := networkConfig.Wallets()
// 		Expect(err).ShouldNot(HaveOccurred())

// 		// Deploy LINK contract
// 		linkInstance, err := DeployLinkTokenContract(client, wallets.Default())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		name, err := linkInstance.Name(context.Background())
// 		Expect(err).ShouldNot(HaveOccurred())
// 		Expect(name).To(Equal("ChainLink Token"))

// 		// Deploy Offchain contract
// 		offChainInstance, err := DeployOffChainAggregator(client, wallets.Default(), offchainOptions)
// 		Expect(err).ShouldNot(HaveOccurred())
// 		err = offChainInstance.Fund(wallets.Default(), nil, big.NewInt(50000000000))
// 		Expect(err).ShouldNot(HaveOccurred())
// 	},
// 		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, OffchainOptions{
// 			// Some defaults from Lorenz's project. Not sure if we want this in a config file in future?
// 			MaximumGasPrice:         uint32(1000),
// 			ReasonableGasPrice:      uint32(200),
// 			MicroLinkPerEth:         3.6e7,
// 			LinkGweiPerObservation:  1e8,
// 			LinkGweiPerTransmission: 4e8,
// 			MinimumAnswer:           big.NewInt(1),
// 			MaximumAnswer:           big.NewInt(100),
// 			Decimals:                uint8(8),
// 			Description:             "Hardhat OffChain Aggregator",
// 		}),
// 		// Tested locally successfully. We need to implement secrets system as well as testing wallets for CI use
// 		// Entry("on Ethereum Kovan", client.NewKovanNetwork, big.NewInt(5)),
// 		// Entry("on Ethereum Goerli", client.NewGoerliNetwork, big.NewInt(5)),
// 	)
// })

// TODO: Templatize these
func buildOCRSpec(contractAddress, p2pId, bootstrapP2PId, keyBundleId, transmitterAddress string) string {
	return fmt.Sprintf(`type = "offchainreporting"
schemaVersion = 1
contractAddress = "%v"
p2pPeerID = "%v"
p2pBootstrapPeers = [
		"/dns4/chainlink-node-1/tcp/6690/p2p/%v"  
]
isBootstrapPeer = false
keyBundleID = "%v"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "%v"
observationTimeout = "10s"
blockchainTimeout  = "20s"
contractConfigTrackerSubscribeInterval = "2m"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
	fetch    [type=http method=POST url="http://host.docker.internal:6644/random" requestData="{}"];
	parse    [type=jsonparse path="data,result"];    
	fetch -> parse;
	"""`, contractAddress, p2pId, bootstrapP2PId, keyBundleId, transmitterAddress)
}

func buildBootstrapSpec(contractAddress string, p2pID string) string {
	return fmt.Sprintf(`blockchainTimeout = "20s"
contractAddress = "%v"
contractConfigConfirmations = 3
contractConfigTrackerPollInterval = "1m"
contractConfigTrackerSubscribeInterval = "2m"
isBootstrapPeer = true
p2pBootstrapPeers = []
p2pPeerID = "%v"
schemaVersion = 1
type = "offchainreporting"`, contractAddress, p2pID)
}
