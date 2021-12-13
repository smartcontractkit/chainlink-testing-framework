package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"
	"time"

	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

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
		for i := 0; i < len(ocrInstances); i++ {
			bootstrapNode := chainlinkNodes[0]
			bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
			bootstrapSpec := &client.OCRBootstrapJobSpec{
				Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
				ContractAddress: ocrInstances[i].Address(),
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

				bta := client.BridgeTypeAttributes{
					Name: fmt.Sprintf("node_%d_contract_%d", nodeIndex, i),
					URL:  fmt.Sprintf("%s/node_%d_contract_%d", mockserver.Config.ClusterURL, nodeIndex, i),
				}

				err = chainlinkNodes[nodeIndex].CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred())

				ocrSpec := &client.OCRTaskJobSpec{
					ContractAddress:    ocrInstances[i].Address(),
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

// SetAdapterResponses sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters, to be used in combination with CreateOCRJobs
func SetAdapterResponses(
	results []int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
	mockserver *client.MockserverClient,
) func() {
	return func() {
		Expect(len(results)).Should(BeNumerically("==", len(chainlinkNodes[1:])))
		for OCRInstanceIndex := range ocrInstances {
			for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
				path := fmt.Sprintf("/node_%d_contract_%d", nodeIndex, OCRInstanceIndex)
				err := mockserver.SetValuePath(path, results[nodeIndex-1])
				Expect(err).ShouldNot(HaveOccurred())
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
