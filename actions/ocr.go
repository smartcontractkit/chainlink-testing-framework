package actions

import (
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"math/big"
	"time"
)

// OCRSetup contains the components needed for testing OCR
type OCRSetup struct {
	Networks          *client.Networks
	ContractDeployer  contracts.ContractDeployer
	LinkTokenContract contracts.LinkToken
	ChainlinkNodes    []client.Chainlink
	Mockserver        *client.MockserverClient
	Env               *environment.Environment
	OCRInstances      []contracts.OffchainAggregator
}

// NewOCRSetup returns a freshly created setup for OCR tests
func NewOCRSetup(e *environment.Environment, chainlinkCharts []string) (*OCRSetup, error) {
	o := &OCRSetup{}
	o.Env = e
	networkRegistry := client.NewNetworkRegistry()
	var err error
	o.Networks, err = networkRegistry.GetNetworks(e)
	if err != nil {
		return nil, err
	}
	o.ContractDeployer, err = contracts.NewContractDeployer(o.Networks.Default)
	if err != nil {
		return nil, err
	}
	o.ChainlinkNodes, err = client.NewChainlinkClients2(e, chainlinkCharts)
	if err != nil {
		return nil, err
	}
	o.Mockserver, err = client.NewMockServerClientFromEnv(e)
	if err != nil {
		return nil, err
	}
	o.Networks.Default.ParallelTransactions(true)
	return o, nil
}

// FundNodes funds all chainlink nodes
func (o *OCRSetup) FundNodes() error {
	txCost, err := o.Networks.Default.CalculateTXSCost(200)
	if err != nil {
		return err
	}
	err = FundChainlinkNodes(o.ChainlinkNodes, o.Networks.Default, txCost)
	if err != nil {
		return err
	}
	return nil
}

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func (o *OCRSetup) DeployOCRContracts(nrOfOCRContracts int) error {
	var err error
	o.LinkTokenContract, err = o.ContractDeployer.DeployLinkTokenContract()
	if err != nil {
		return err
	}

	for nr := 0; nr < nrOfOCRContracts; nr++ {
		OCRInstance, err := o.ContractDeployer.DeployOffChainAggregator(o.LinkTokenContract.Address(), contracts.DefaultOffChainAggregatorOptions())
		if err != nil {
			return err
		}
		err = OCRInstance.SetConfig(
			o.ChainlinkNodes[1:],
			contracts.DefaultOffChainAggregatorConfig(len(o.ChainlinkNodes[1:])),
		)
		if err != nil {
			return err
		}
		err = o.LinkTokenContract.Transfer(OCRInstance.Address(), big.NewInt(2e18))
		if err != nil {
			return err
		}
		err = o.Networks.Default.WaitForEvents()
		if err != nil {
			return err
		}
		o.OCRInstances = append(o.OCRInstances, OCRInstance)
	}
	return nil
}

// CreateOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters
func (o *OCRSetup) CreateOCRJobs() error {
	for OCRInstanceIndex, OCRInstance := range o.OCRInstances {
		bootstrapNode := o.ChainlinkNodes[0]
		bootstrapP2PIds, err := bootstrapNode.ReadP2PKeys()
		if err != nil {
			return err
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			Name:            fmt.Sprintf("bootstrap-%s", uuid.NewV4().String()),
			ContractAddress: OCRInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.CreateJob(bootstrapSpec)
		if err != nil {
			return err
		}

		for nodeIndex := 1; nodeIndex < len(o.ChainlinkNodes); nodeIndex++ {
			nodeP2PIds, err := o.ChainlinkNodes[nodeIndex].ReadP2PKeys()
			if err != nil {
				return err
			}
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := o.ChainlinkNodes[nodeIndex].PrimaryEthAddress()
			if err != nil {
				return err
			}
			nodeOCRKeys, err := o.ChainlinkNodes[nodeIndex].ReadOCRKeys()
			if err != nil {
				return err
			}
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			bta := client.BridgeTypeAttributes{
				Name: fmt.Sprintf("node_%d_contract_%d", nodeIndex, OCRInstanceIndex),
				URL:  fmt.Sprintf("%s/node_%d_contract_%d", o.Mockserver.Config.ClusterURL, nodeIndex, OCRInstanceIndex),
			}

			err = o.ChainlinkNodes[nodeIndex].CreateBridge(&bta)
			if err != nil {
				return err
			}

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    OCRInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
			}
			_, err = o.ChainlinkNodes[nodeIndex].CreateJob(ocrSpec)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// StartNewRound requests a new round from the ocr contract and waits for confirmation
func (o *OCRSetup) StartNewRound(roundNr int64) error {
	roundTimeout := time.Minute * 2
	for _, OCRInstance := range o.OCRInstances {
		err := OCRInstance.RequestNewRound()
		if err != nil {
			return err
		}
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(OCRInstance, big.NewInt(roundNr), roundTimeout)
		o.Networks.Default.AddHeaderEventSubscription(OCRInstance.Address(), ocrRound)
		err = o.Networks.Default.WaitForEvents()
		if err != nil {
			return err
		}
	}
	return nil
}

// SetAdapterResults sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters
func (o *OCRSetup) SetAdapterResults(results []int) error {
	if len(results) != len(o.ChainlinkNodes[1:]) {
		return errors.New("Number of results should equal number of nodes")
	}

	for OCRInstanceIndex := range o.OCRInstances {
		for nodeIndex := 1; nodeIndex < len(o.ChainlinkNodes); nodeIndex++ {
			path := fmt.Sprintf("/node_%d_contract_%d", nodeIndex, OCRInstanceIndex)
			pathSelector := client.PathSelector{Path: path}
			err := o.Mockserver.ClearExpectation(pathSelector)
			if err != nil {
				return err
			}
			err = o.Mockserver.SetValuePath(path, results[nodeIndex-1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CheckRound checks the ocr rounds for correctness
//func (o *OCRSetup) CheckRound() error {
//	// Set adapters answer to 5
//	var adapterResults []int
//	for index := 1; index < len(i.ChainlinkNodes); index++ {
//		result := 5
//		adapterResults = append(adapterResults, result)
//	}
//	SetAdapterResults(i, adapterResults)
//
//	StartNewRound(i, 1)
//
//	// Check answer is as expected
//	for _, OCRInstance := range i.OCRInstances {
//		answer, err := OCRInstance.GetLatestAnswer(context.Background())
//		Expect(err).ShouldNot(HaveOccurred())
//		Expect(answer.Int64()).Should(Equal(int64(5)), "Latest answer from OCR is not as expected")
//	}
//
//	// Change adapters answer to 10
//	adapterResults = []int{}
//	for index := 1; index < len(i.ChainlinkNodes); index++ {
//		result := 10
//		adapterResults = append(adapterResults, result)
//	}
//	SetAdapterResults(i, adapterResults)
//
//	StartNewRound(i, 2)
//
//	// Check answer is as expected
//	for _, OCRInstance := range i.OCRInstances {
//		answer, err := OCRInstance.GetLatestAnswer(context.Background())
//		Expect(err).ShouldNot(HaveOccurred())
//		Expect(answer.Int64()).Should(Equal(int64(10)), "Latest answer from OCR is not as expected")
//	}
//	return nil
//}
