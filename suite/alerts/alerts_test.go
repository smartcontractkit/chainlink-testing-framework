package alerts

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/environment/charts/mockserver"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"os"
	"path/filepath"
)

var _ = Describe("Alerts suite", func() {
	var (
		suiteSetup *actions.DefaultSuiteSetup
		//adapter        environment.ExternalAdapter
		chainlinkNodes []client.Chainlink
		//explorer       *client.ExplorerClient
		err error
	)

	Describe("Alerts", func() {
		It("Test 1", func() {
			suiteSetup, err = actions.DefaultLocalSetup(
				"basic-chainlink",
				environment.NewChainlinkClusterForAlertsTesting(5),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())

			//explorer, err = environment.GetExplorerClientFromEnv(suiteSetup.Env)
			//Expect(err).ShouldNot(HaveOccurred())
			//fmt.Println(explorer.BaseURL)
			//
			chainlinkNodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(chainlinkNodes[0].URL())
			//
			//adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			//Expect(err).ShouldNot(HaveOccurred())
			//fmt.Println(adapter.ClusterURL())

			/* ####################################################### */
			/* ####################################################### */
			/* ####################################################### */
			/* ####################################################### */

			err := actions.FundChainlinkNodes(
				chainlinkNodes,
				suiteSetup.Client,
				suiteSetup.Wallets.Default(),
				big.NewFloat(0.05),
				big.NewFloat(2),
			)
			Expect(err).ShouldNot(HaveOccurred())

			// Deploy and config OCR contract

			OCRInstance, err := suiteSetup.Deployer.DeployOffChainAggregator(suiteSetup.Wallets.Default(), contracts.DefaultOffChainAggregatorOptions())
			Expect(err).ShouldNot(HaveOccurred())
			err = OCRInstance.SetConfig(
				suiteSetup.Wallets.Default(),
				chainlinkNodes,
				contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes)),
			)
			Expect(err).ShouldNot(HaveOccurred())
			err = OCRInstance.Fund(suiteSetup.Wallets.Default(), nil, big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())

			/* Write to initializerJson the stuff needed for otpe */

			contractInfo := mockserver.ContractInfoJSON{
				ContractVersion: 4,
				Path:            "test",
				Status:          "live",
				ContractAddress: OCRInstance.Address(),
			}

			contractsInfo := []mockserver.ContractInfoJSON{contractInfo}

			contractsInitializer := mockserver.HttpInitializer{
				Request: mockserver.HttpRequest{Path: "/contracts"},
				Response: mockserver.HttpResponse{Body: contractsInfo},
			}

			var nodesInfo []mockserver.NodeInfoJSON

			for _, chainlink := range chainlinkNodes {
				ocrKeys, err := chainlink.ReadOCRKeys()
				Expect(err).ShouldNot(HaveOccurred())
				nodeInfo := mockserver.NodeInfoJSON{
					NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
					ID: ocrKeys.Data[0].ID,
				}
				nodesInfo = append(nodesInfo, nodeInfo)
			}

			nodesInitializer := mockserver.HttpInitializer{
				Request: mockserver.HttpRequest{Path: "/nodes"},
				Response: mockserver.HttpResponse{Body: nodesInfo},
			}
			initializers := []mockserver.HttpInitializer{contractsInitializer, nodesInitializer}

			initializersBytes, err := json.Marshal(initializers)
			Expect(err).ShouldNot(HaveOccurred())

			fileName := filepath.Join(tools.ProjectRoot, "environment/charts/mockserver-config/static/initializerJson.json")
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			Expect(err).ShouldNot(HaveOccurred())

			body := fmt.Sprintf(string(initializersBytes))
			_, err = f.WriteString(body)
			Expect(err).ShouldNot(HaveOccurred())

			err = f.Close()
			Expect(err).ShouldNot(HaveOccurred())

			// Deploy mockserver and otpe

			err = suiteSetup.Env.DeploySpecs(environment.OtpeGroup())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
