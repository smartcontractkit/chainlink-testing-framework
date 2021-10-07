package testcommon

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/suite/steps"
	"math/big"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

// OCRSetupInputs inputs needed for OCR tests
type OCRSetupInputs struct {
	SuiteSetup     *actions.DefaultSuiteSetup
	ChainlinkNodes []client.Chainlink
	DefaultWallet  client.BlockchainWallet
	OCRInstance    contracts.OffchainAggregator
	Mockserver     *client.MockserverClient
}

// DeployOCRForEnv deploys the environment
func DeployOCRForEnv(i *OCRSetupInputs, envName string, envInit environment.K8sEnvSpecInit) {
	By("Deploying the environment", func() {
		var err error
		i.SuiteSetup, err = actions.DefaultLocalSetup(
			envName,
			envInit,
			client.NewNetworkFromConfig,
			tools.ProjectRoot,
		)
		Expect(err).ShouldNot(HaveOccurred())
		i.Mockserver, err = environment.GetMockserverClientFromEnv(i.SuiteSetup.Env)
		Expect(err).ShouldNot(HaveOccurred())

		i.ChainlinkNodes, err = environment.GetChainlinkClients(i.SuiteSetup.Env)
		Expect(err).ShouldNot(HaveOccurred())
		i.DefaultWallet = i.SuiteSetup.Wallets.Default()
		i.SuiteSetup.Client.ParallelTransactions(true)
	})
}

// SetupOCRTest setup for an ocr test
func SetupOCRTest(i *OCRSetupInputs) {
	By("Funding nodes and deploying OCR contract", func() {
		ethAmount, err := i.SuiteSetup.Deployer.CalculateETHForTXs(i.SuiteSetup.Wallets.Default(), i.SuiteSetup.Network.Config(), 2)
		Expect(err).ShouldNot(HaveOccurred())
		err = actions.FundChainlinkNodes(
			i.ChainlinkNodes,
			i.SuiteSetup.Client,
			i.DefaultWallet,
			ethAmount,
			big.NewFloat(2),
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Deploy and config OCR contract
		deployer, err := contracts.NewContractDeployer(i.SuiteSetup.Client)
		Expect(err).ShouldNot(HaveOccurred())

		i.OCRInstance, err = deployer.DeployOffChainAggregator(i.DefaultWallet, contracts.DefaultOffChainAggregatorOptions())
		Expect(err).ShouldNot(HaveOccurred())
		err = i.OCRInstance.SetConfig(
			i.DefaultWallet,
			i.ChainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(i.ChainlinkNodes[1:])),
		)
		Expect(err).ShouldNot(HaveOccurred())
		err = i.OCRInstance.Fund(i.DefaultWallet, nil, big.NewFloat(2))
		Expect(err).ShouldNot(HaveOccurred())
		err = i.SuiteSetup.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
	})
}

// SendOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters
func SendOCRJobs(i *OCRSetupInputs) {
	By("Sending OCR jobs to chainlink nodes", func() {
		bootstrapNode := i.ChainlinkNodes[0]
		bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
		Expect(err).ShouldNot(HaveOccurred())
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			ContractAddress: i.OCRInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.CreateJob(bootstrapSpec)
		Expect(err).ShouldNot(HaveOccurred())

		for index := 1; index < len(i.ChainlinkNodes); index++ {
			nodeP2PIds, err := i.ChainlinkNodes[index].ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := i.ChainlinkNodes[index].PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeys, err := i.ChainlinkNodes[index].ReadOCRKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("node_%d", index),
				URL:  fmt.Sprintf("%s/node_%d", i.Mockserver.Config.ClusterURL, index),
			}

			err = i.ChainlinkNodes[index].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred())

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    i.OCRInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
			}
			_, err = i.ChainlinkNodes[index].CreateJob(ocrSpec)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})
}

// CheckRound checks the ocr rounds for correctness
func CheckRound(i *OCRSetupInputs) {
	By("Checking OCR rounds", func() {
		// Set adapters answer to 5
		var adapterResults []int
		for index := 1; index < len(i.ChainlinkNodes); index++ {
			result := 5
			adapterResults = append(adapterResults, result)
		}
		SetAdapterResults(i, adapterResults)

		StartNewRound(i, 1)

		// Check answer is as expected
		answer, err := i.OCRInstance.GetLatestAnswer(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(answer.Int64()).Should(Equal(int64(5)), "Latest answer from OCR is not as expected")

		// Change adapters answer to 10
		adapterResults = []int{}
		for index := 1; index < len(i.ChainlinkNodes); index++ {
			result := 10
			adapterResults = append(adapterResults, result)
		}
		SetAdapterResults(i, adapterResults)

		StartNewRound(i, 2)

		// Check answer is as expected
		answer, err = i.OCRInstance.GetLatestAnswer(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(answer.Int64()).Should(Equal(int64(10)), "Latest answer from OCR is not as expected")
	})
}

// StartNewRound requests a new round from the ocr contract and waits for confirmation
func StartNewRound(i *OCRSetupInputs, roundNr int64) {
	roundTimeout := time.Minute * 2

	err := i.OCRInstance.RequestNewRound(i.DefaultWallet)
	Expect(err).ShouldNot(HaveOccurred())
	err = i.SuiteSetup.Client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	// Wait for the second round
	ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(i.OCRInstance, big.NewInt(roundNr), roundTimeout)
	i.SuiteSetup.Client.AddHeaderEventSubscription(i.OCRInstance.Address(), ocrRound)
	err = i.SuiteSetup.Client.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())
}

// SetAdapterResults sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters
func SetAdapterResults(i *OCRSetupInputs, results []int) {
	Expect(len(results)).Should(BeNumerically("==", len(i.ChainlinkNodes[1:])))

	log.Info().Interface("New Adapter results", results).Msg("Setting new values")

	for index := 1; index < len(i.ChainlinkNodes); index++ {
		pathSelector := client.PathSelector{Path: fmt.Sprintf("/node_%d", index)}
		err := i.Mockserver.ClearExpectation(pathSelector)
		Expect(err).ShouldNot(HaveOccurred())
	}

	var initializers []client.HttpInitializer
	for index := 1; index < len(i.ChainlinkNodes); index++ {
		adResp := client.AdapterResponse{
			Id:    "",
			Data:  client.AdapterResult{Result: results[index-1]},
			Error: nil,
		}
		nodesInitializer := client.HttpInitializer{
			Request:  client.HttpRequest{Path: fmt.Sprintf("/node_%d", index)},
			Response: client.HttpResponse{Body: adResp},
		}
		initializers = append(initializers, nodesInitializer)
	}

	err := i.Mockserver.PutExpectations(initializers)
	Expect(err).ShouldNot(HaveOccurred())
}

// NewOCRSetupInputForObservability deploys and setups env and clients for testing observability
func NewOCRSetupInputForObservability(i *OCRSetupInputs, nodeCount int, rules map[string]*os.File) {
	DeployOCRForEnv(
		i,
		"basic-chainlink",
		environment.NewChainlinkClusterForObservabilityTesting(nodeCount),
	)
	SetupOCRTest(i)

	err := i.Mockserver.PutExpectations(steps.GetMockserverInitializerDataForOTPE(
		i.OCRInstance.Address(),
		i.ChainlinkNodes,
	))
	Expect(err).ShouldNot(HaveOccurred())

	err = i.SuiteSetup.Env.DeploySpecs(environment.OtpeGroup())
	Expect(err).ShouldNot(HaveOccurred())

	err = i.SuiteSetup.Env.DeploySpecs(environment.PrometheusGroup(rules))
	Expect(err).ShouldNot(HaveOccurred())
}