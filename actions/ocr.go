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
	ocrInstances := []contracts.OffchainAggregator{}
	for i := 0; i < numberOfContracts; i++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		Expect(err).ShouldNot(HaveOccurred())
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			chainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
		)
		ocrInstances = append(ocrInstances, ocrInstance)
		Expect(err).ShouldNot(HaveOccurred())
		err = linkTokenContract.Transfer(ocrInstance.Address(), big.NewInt(2e18))
		Expect(err).ShouldNot(HaveOccurred())
		err = networks.Default.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
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
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocrInstance.Address(),
				P2PPeerID:       bootstrapP2PId,
				IsBootstrapPeer: true,
			}
			_, err = bootstrapNode.CreateJob(bootstrapSpec)
			Expect(err).ShouldNot(HaveOccurred())

			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				nodeP2PIds, err := chainlinkNodes[nodeIndex].ReadP2PKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
				nodeTransmitterAddress, err := chainlinkNodes[nodeIndex].PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeys, err := chainlinkNodes[nodeIndex].ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeOCRKeyId := nodeOCRKeys.Data[0].ID

				nodeContractPairID := buildNodeContractPairID(chainlinkNodes[nodeIndex], ocrInstance)
				Expect(err).ShouldNot(HaveOccurred())
				bta := client.BridgeTypeAttributes{
					Name: nodeContractPairID,
					URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, nodeContractPairID),
				}

				// This sets a default value for all node and ocr instances in order to avoid 404 issues
				SetAllAdapterResponses(0, ocrInstances, chainlinkNodes, mockserver)

				err = chainlinkNodes[nodeIndex].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstance.Address(),
					P2PPeerID:          nodeP2PId,
					P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
					KeyBundleID:        nodeOCRKeyId,
					TransmitterAddress: nodeTransmitterAddress,
					ObservationSource:  client.ObservationSourceSpecBridge(bta),
				}
				_, err = chainlinkNodes[nodeIndex].CreateJob(ocrSpec)
				Expect(err).ShouldNot(HaveOccurred())
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
		nodeContractPairID := buildNodeContractPairID(chainlinkNode, ocrInstance)
		path := fmt.Sprintf("/%s", nodeContractPairID)
		err := mockserver.SetValuePath(path, response)
		Expect(err).ShouldNot(HaveOccurred())
	}
}

// SetAllAdapterResponses sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to the same response
func SetAllAdapterResponses(
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
			Expect(err).ShouldNot(HaveOccurred())
			ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstances[i], big.NewInt(roundNr), roundTimeout)
			networks.Default.AddHeaderEventSubscription(ocrInstances[i].Address(), ocrRound)
			err = networks.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		}
	}
}

func buildNodeContractPairID(node client.Chainlink, ocrInstance contracts.OffchainAggregator) string {
	Expect(node).ShouldNot(BeNil())
	Expect(ocrInstance).ShouldNot(BeNil())
	nodeAddress, err := node.PrimaryEthAddress()
	Expect(err).ShouldNot(HaveOccurred())
	shortNodeAddr := nodeAddress[2:12]
	shortOCRAddr := ocrInstance.Address()[2:12]
	return strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr))
}
