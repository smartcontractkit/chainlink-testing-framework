package testcommon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smartcontractkit/integrations-framework/environment/charts/mockserver"
	"math/big"
	"os"
	"path/filepath"
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
	Adapter        environment.ExternalAdapter
	DefaultWallet  client.BlockchainWallet
	OCRInstance    contracts.OffchainAggregator
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
		i.Adapter, err = environment.GetExternalAdapter(i.SuiteSetup.Env)
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
			i.ChainlinkNodes,
			contracts.DefaultOffChainAggregatorConfig(len(i.ChainlinkNodes)),
		)
		Expect(err).ShouldNot(HaveOccurred())
		err = i.OCRInstance.Fund(i.DefaultWallet, nil, big.NewFloat(2))
		Expect(err).ShouldNot(HaveOccurred())
		err = i.SuiteSetup.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
	})

	By("Sending OCR jobs to chainlink nodes", func() {
		// Initialize bootstrap node
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

		bta := client.BridgeTypeAttributes{
			Name: "variable",
			URL:  fmt.Sprintf("%s/variable", i.Adapter.ClusterURL()),
		}

		// Send OCR job to other nodes
		for index := 1; index < len(i.ChainlinkNodes); index++ {
			nodeP2PIds, err := i.ChainlinkNodes[index].ReadP2PKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := i.ChainlinkNodes[index].PrimaryEthAddress()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeys, err := i.ChainlinkNodes[index].ReadOCRKeys()
			Expect(err).ShouldNot(HaveOccurred())
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

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
		roundTimeout := time.Minute * 2
		// Set adapter answer to 5
		err := i.Adapter.SetVariable(5)
		Expect(err).ShouldNot(HaveOccurred())
		err = i.OCRInstance.RequestNewRound(i.DefaultWallet)
		Expect(err).ShouldNot(HaveOccurred())
		err = i.SuiteSetup.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		// Wait for the first round
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(i.OCRInstance, big.NewInt(1), roundTimeout)
		i.SuiteSetup.Client.AddHeaderEventSubscription(i.OCRInstance.Address(), ocrRound)
		err = i.SuiteSetup.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		// Check answer is as expected
		answer, err := i.OCRInstance.GetLatestAnswer(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(answer.Int64()).Should(Equal(int64(5)), "Latest answer from OCR is not as expected")

		// Change adapter answer to 10
		err = i.Adapter.SetVariable(10)
		Expect(err).ShouldNot(HaveOccurred())

		// Wait for the second round
		ocrRound = contracts.NewOffchainAggregatorRoundConfirmer(i.OCRInstance, big.NewInt(2), roundTimeout)
		i.SuiteSetup.Client.AddHeaderEventSubscription(i.OCRInstance.Address(), ocrRound)
		err = i.SuiteSetup.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		// Check answer is as expected
		answer, err = i.OCRInstance.GetLatestAnswer(context.Background())
		Expect(err).ShouldNot(HaveOccurred())
		Expect(answer.Int64()).Should(Equal(int64(10)), "Latest answer from OCR is not as expected")
	})
}

// WriteDataForOTPEToInitializerFileForMockserver Write to initializerJson mocked weiwatchers data needed for otpe
func WriteDataForOTPEToInitializerFileForMockserver(i *OCRSetupInputs) {
	contractInfo := mockserver.ContractInfoJSON{
		ContractVersion: 4,
		Path:            "test",
		Status:          "live",
		ContractAddress: i.OCRInstance.Address(),
	}

	contractsInfo := []mockserver.ContractInfoJSON{contractInfo}

	contractsInitializer := mockserver.HttpInitializer{
		Request:  mockserver.HttpRequest{Path: "/contracts.json"},
		Response: mockserver.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []mockserver.NodeInfoJSON

	for _, chainlink := range i.ChainlinkNodes {
		ocrKeys, err := chainlink.ReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred())
		nodeInfo := mockserver.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := mockserver.HttpInitializer{
		Request:  mockserver.HttpRequest{Path: "/nodes.json"},
		Response: mockserver.HttpResponse{Body: nodesInfo},
	}
	initializers := []mockserver.HttpInitializer{contractsInitializer, nodesInitializer}

	initializersBytes, err := json.Marshal(initializers)
	Expect(err).ShouldNot(HaveOccurred())

	fileName := filepath.Join(tools.ProjectRoot, "environment/charts/mockserver-config/static/initializerJson.json")
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	Expect(err).ShouldNot(HaveOccurred())

	body := string(initializersBytes)
	_, err = f.WriteString(body)
	Expect(err).ShouldNot(HaveOccurred())

	err = f.Close()
	Expect(err).ShouldNot(HaveOccurred())
}
