package testcommon

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/suite/steps"

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
	SuiteSetup     actions.SuiteSetup
	NetworkInfo    actions.NetworkInfo
	ChainlinkNodes []client.Chainlink
	DefaultWallet  client.BlockchainWallet
	OCRInstances   []contracts.OffchainAggregator
	Mockserver     *client.MockserverClient
}

// DeployOCRForEnv deploys the environment
func DeployOCRForEnv(i *OCRSetupInputs, envInit environment.K8sEnvSpecInit) {
	By("Deploying the environment", func() {
		var err error
		i.SuiteSetup, err = actions.SingleNetworkSetup(
			envInit,
			actions.EVMNetworkFromConfigHook,
			actions.EthereumDeployerHook,
			actions.EthereumClientHook,
			tools.ProjectRoot,
		)
		Expect(err).ShouldNot(HaveOccurred())
		i.Mockserver, err = environment.GetMockserverClientFromEnv(i.SuiteSetup.Environment())
		Expect(err).ShouldNot(HaveOccurred())

		i.ChainlinkNodes, err = environment.GetChainlinkClients(i.SuiteSetup.Environment())
		Expect(err).ShouldNot(HaveOccurred())
		i.NetworkInfo = i.SuiteSetup.DefaultNetwork()
		i.DefaultWallet = i.NetworkInfo.Wallets.Default()
		i.NetworkInfo.Client.ParallelTransactions(true)
	})
}

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRContracts(i *OCRSetupInputs, nrOfOCRContracts int) {
	deployer, err := contracts.NewContractDeployer(i.NetworkInfo.Client)
	Expect(err).ShouldNot(HaveOccurred())

	for nr := 0; nr < nrOfOCRContracts; nr++ {
		OCRInstance, err := deployer.DeployOffChainAggregator(i.DefaultWallet, contracts.DefaultOffChainAggregatorOptions())
		Expect(err).ShouldNot(HaveOccurred())
		err = OCRInstance.SetConfig(
			i.DefaultWallet,
			i.ChainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(i.ChainlinkNodes[1:])),
		)
		Expect(err).ShouldNot(HaveOccurred())
		err = OCRInstance.Fund(i.DefaultWallet, nil, big.NewFloat(2))
		Expect(err).ShouldNot(HaveOccurred())
		err = i.NetworkInfo.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
		i.OCRInstances = append(i.OCRInstances, OCRInstance)
	}
}

// FundNodes funds all chainlink nodes
func FundNodes(i *OCRSetupInputs) {
	ethAmount, err := i.NetworkInfo.Deployer.CalculateETHForTXs(i.NetworkInfo.Wallets.Default(), i.NetworkInfo.Network.Config(), 2)
	Expect(err).ShouldNot(HaveOccurred())
	err = actions.FundChainlinkNodes(
		i.ChainlinkNodes,
		i.NetworkInfo.Client,
		i.DefaultWallet,
		ethAmount,
		big.NewFloat(2),
	)
	Expect(err).ShouldNot(HaveOccurred())
}

// SendOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters
func SendOCRJobs(i *OCRSetupInputs) {
	By("Sending OCR jobs to chainlink nodes", func() {
		for OCRInstanceIndex, OCRInstance := range i.OCRInstances {
			bootstrapNode := i.ChainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				ContractAddress: OCRInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred())

			for nodeIndex := 1; nodeIndex < len(i.ChainlinkNodes); nodeIndex++ {
				nodeP2PIds, err := i.ChainlinkNodes[nodeIndex].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := i.ChainlinkNodes[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeys, err := i.ChainlinkNodes[nodeIndex].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				bta := client.BridgeTypeAttributes{
					Name: fmt.Sprintf("node_%d_contract_%d", nodeIndex, OCRInstanceIndex),
					URL:  fmt.Sprintf("%s/node_%d_contract_%d", i.Mockserver.Config.ClusterURL, nodeIndex, OCRInstanceIndex),
				}

				err = i.ChainlinkNodes[nodeIndex].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    OCRInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = i.ChainlinkNodes[nodeIndex].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}
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
		for _, OCRInstance := range i.OCRInstances {
			answer, err := OCRInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(5)), "Latest answer from OCR is not as expected")
		}

		// Change adapters answer to 10
		adapterResults = []int{}
		for index := 1; index < len(i.ChainlinkNodes); index++ {
			result := 10
			adapterResults = append(adapterResults, result)
		}
		SetAdapterResults(i, adapterResults)

		StartNewRound(i, 2)

		// Check answer is as expected
		for _, OCRInstance := range i.OCRInstances {
			answer, err := OCRInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(answer.Int64()).Should(Equal(int64(10)), "Latest answer from OCR is not as expected")
		}
	})
}

// StartNewRound requests a new round from the ocr contract and waits for confirmation
func StartNewRound(i *OCRSetupInputs, roundNr int64) {
	roundTimeout := time.Minute * 2
	for _, OCRInstance := range i.OCRInstances {
		err := OCRInstance.RequestNewRound(i.DefaultWallet)
		Expect(err).ShouldNot(HaveOccurred())
		err = i.SuiteSetup.DefaultNetwork().Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		// Wait for the second round
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(OCRInstance, big.NewInt(roundNr), roundTimeout)
		i.SuiteSetup.DefaultNetwork().Client.AddHeaderEventSubscription(OCRInstance.Address(), ocrRound)
		err = i.SuiteSetup.DefaultNetwork().Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
	}
}

// SetAdapterResults sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters
func SetAdapterResults(i *OCRSetupInputs, results []int) {
	Expect(len(results)).Should(BeNumerically("==", len(i.ChainlinkNodes[1:])))

	log.Info().Interface("New Adapter results", results).Msg("Setting new values")

	for OCRInstanceIndex := range i.OCRInstances {
		for nodeIndex := 1; nodeIndex < len(i.ChainlinkNodes); nodeIndex++ {
			pathSelector := client.PathSelector{Path: fmt.Sprintf("/node_%d_contract_%d", nodeIndex, OCRInstanceIndex)}
			err := i.Mockserver.ClearExpectation(pathSelector)
			Expect(err).ShouldNot(HaveOccurred())
		}

	}
	var initializers []client.HttpInitializer

	for OCRInstanceIndex := range i.OCRInstances {
		for nodeIndex := 1; nodeIndex < len(i.ChainlinkNodes); nodeIndex++ {
			adResp := client.AdapterResponse{
				Id:    "",
				Data:  client.AdapterResult{Result: results[nodeIndex-1]},
				Error: nil,
			}
			nodesInitializer := client.HttpInitializer{
				Request:  client.HttpRequest{Path: fmt.Sprintf("/node_%d_contract_%d", nodeIndex, OCRInstanceIndex)},
				Response: client.HttpResponse{Body: adResp},
			}
			initializers = append(initializers, nodesInitializer)
		}
	}

	err := i.Mockserver.PutExpectations(initializers)
	Expect(err).ShouldNot(HaveOccurred())
}

// NewOCRSetupInputForObservability deploys and setups env and clients for testing observability
func NewOCRSetupInputForObservability(i *OCRSetupInputs, nodeCount int, contractCount int, rules map[string]*os.File) {
	DeployOCRForEnv(
		i,
		environment.NewChainlinkClusterForObservabilityTesting(nodeCount),
	)
	FundNodes(i)
	DeployOCRContracts(i, contractCount)

	err := i.Mockserver.PutExpectations(steps.GetMockserverInitializerDataForOTPE(
		i.OCRInstances,
		i.ChainlinkNodes,
	))
	Expect(err).ShouldNot(HaveOccurred())

	err = i.SuiteSetup.Environment().DeploySpecs(environment.OtpeGroup())
	Expect(err).ShouldNot(HaveOccurred())

	err = i.SuiteSetup.Environment().DeploySpecs(environment.PrometheusGroup(rules))
	Expect(err).ShouldNot(HaveOccurred())
}
