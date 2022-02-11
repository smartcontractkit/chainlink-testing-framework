package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

// This actions file often returns functions, rather than just values. These are used as common test helpers, and are
// handy to have returning as functions so that Ginkgo can use them in an aesthetically pleasing way.

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRContracts(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []client.Chainlink,
	networks *client.Networks,
) []contracts.OffchainAggregator {
	var ocrInstances []contracts.OffchainAggregator
	for i := 0; i < numberOfContracts; i++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		Expect(err).ShouldNot(HaveOccurred(), "Deploying OCR instance %d shouldn't fail", i+1)
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			chainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
		)
		ocrInstances = append(ocrInstances, ocrInstance)
		Expect(err).ShouldNot(HaveOccurred())
		err = linkTokenContract.Transfer(ocrInstance.Address(), big.NewInt(2e18))
		Expect(err).ShouldNot(HaveOccurred(), "Transfering LINK token to OCR instance %d shouldn't fail", i+1)
		err = networks.Default.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred(), "Waiting for Event subscriptions of OCR instance %d shouldn't fail", i+1)
	}
	return ocrInstances
}

// CreateOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRJobs(
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
	mockserver *client.MockserverClient,
) func() {
	return func() {
		for _, ocrInstance := range ocrInstances {
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from bootstrap node")
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bootstrap job on bootstrap node")

			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				nodeP2PIds, err := chainlinkNodes[nodeIndex].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail reading P2P keys from OCR node %d", nodeIndex+1)
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
				nodeOCRKeys, err := chainlinkNodes[nodeIndex].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				nodeContractPairID := BuildNodeContractPairID(chainlinkNodes[nodeIndex], ocrInstance)
				Expect(err).ShouldNot(HaveOccurred())
				bta := client.BridgeTypeAttributes{
					Name: nodeContractPairID,
					URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, nodeContractPairID),
				}

				// This sets a default value for all node and ocr instances in order to avoid 404 issues
				SetAllAdapterResponsesToTheSameValue(0, ocrInstances, chainlinkNodes, mockserver)

				err = chainlinkNodes[nodeIndex].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating bridge in OCR node %d", nodeIndex+1)

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = chainlinkNodes[nodeIndex].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred(), "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
			}
		}
	}
}

// SetAdapterResponse sets a single adapter response that correlates with an ocr contract and a chainlink node
func SetAdapterResponse(
	response int,
	ocrInstance contracts.OffchainAggregator,
	chainlinkNode client.Chainlink,
	mockserver *client.MockserverClient,
) func() {
	return func() {
		nodeContractPairID := BuildNodeContractPairID(chainlinkNode, ocrInstance)
		path := fmt.Sprintf("/%s", nodeContractPairID)
		err := mockserver.SetValuePath(path, response)
		Expect(err).ShouldNot(HaveOccurred(), "Setting mockserver value path shouldn't fail")
	}
}

// SetAllAdapterResponsesToTheSameValue sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to the same response
func SetAllAdapterResponsesToTheSameValue(
	response int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
	mockserver *client.MockserverClient,
) func() {
	return func() {
		for _, ocrInstance := range ocrInstances {
			for _, node := range chainlinkNodes {
				SetAdapterResponse(response, ocrInstance, node, mockserver)()
			}
		}
	}
}

// SetAllAdapterResponsesToDifferentValues sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to different responses
func SetAllAdapterResponsesToDifferentValues(
	responses []int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
	mockserver *client.MockserverClient,
) func() {
	return func() {
		Expect(len(responses)).Should(BeNumerically("==", len(chainlinkNodes[1:])))
		for _, ocrInstance := range ocrInstances {
			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				SetAdapterResponse(responses[nodeIndex-1], ocrInstance, chainlinkNodes[nodeIndex], mockserver)()
			}
		}
	}
}

// StartNewRound requests a new round from the ocr contract and waits for confirmation
func StartNewRound(
	roundNr int64,
	ocrInstances []contracts.OffchainAggregator,
	networks *client.Networks,
) func() {
	return func() {
		roundTimeout := time.Minute * 2
		for i := 0; i < len(ocrInstances); i++ {
			err := ocrInstances[i].RequestNewRound()
			Expect(err).ShouldNot(HaveOccurred(), "Requesting new round in OCR instance %d shouldn't fail", i+1)
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstances[i], big.NewInt(roundNr), roundTimeout)
			networks.Default.AddHeaderEventSubscription(ocrInstances[i].Address(), ocrRound)
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred(), "Waiting for Event subscriptions of OCR instance %d shouldn't fail", i+1)
		}
	}
}

func BuildNodeContractPairID(node client.Chainlink, ocrInstance contracts.OffchainAggregator) string {
	Expect(node).ShouldNot(BeNil())
	Expect(ocrInstance).ShouldNot(BeNil())
	nodeAddress, err := node.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred(), "Getting chainlink node's primary ETH address shouldn't fail")
	shortNodeAddr := nodeAddress[2:12]
	shortOCRAddr := ocrInstance.Address()[2:12]
	return strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr))
}
