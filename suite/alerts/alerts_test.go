package alerts

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"os"
	"path/filepath"
)

var _ = FDescribe("Alerts suite", func() {
	var (
		suiteSetup *actions.DefaultSuiteSetup
		//adapter        environment.ExternalAdapter
		chainlinkNodes []client.Chainlink
		//explorer       *client.ExplorerClient
		err error
	)

	Describe("Alerts", func() {
		It("Test 1", func() {
			//suiteSetup, err = actions.DefaultLocalSetup(
			//	environment.NewChainlinkClusterForAlertsTesting(0),
			//	client.NewNetworkFromConfig,
			//	tools.ProjectRoot,
			//)
			//Expect(err).ShouldNot(HaveOccurred())

			suiteSetup, err = actions.DefaultLocalSetup3(
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

			fileName := filepath.Join(tools.ProjectRoot, "environment/charts/mockserver-config/static/initializerJson.json")
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			Expect(err).ShouldNot(HaveOccurred())

			body := fmt.Sprintf(`[
  {
    "httpRequest": {
      "path": "/opel"
    },
    "httpResponse": {
      "body": "%s"
    }
  }
]`, OCRInstance.Address())
			_, err = f.WriteString(body)
			Expect(err).ShouldNot(HaveOccurred())

			err = f.Close()
			Expect(err).ShouldNot(HaveOccurred())

		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
