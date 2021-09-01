package contracts

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
	//"math/big"
)

var _ = Describe("Flux monitor suite @flux", func() {
	var (
		suiteSetup *actions.DefaultSuiteSetup
		adapter    environment.ExternalAdapter
		chainlinkNodes []client.Chainlink
		explorer       *client.ExplorerClient
		err           error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			//suiteSetup, err = actions.DefaultLocalSetup(
			//	environment.NewChainlinkCluster(1),
			//	client.NewNetworkFromConfig,
			//	tools.ProjectRoot,
			//)


			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkNodesGroups,
				4,
				client.NewNetworkFromConfig,
				tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())

			explorer, err = environment.GetExplorerClient(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(explorer.BaseURL)

			credentials, err := explorer.PostAdminNodes("node_")
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("AccessKey", credentials.AccessKey).Msg("AccessKey")

			chainlinkNodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(chainlinkNodes[0].URL())

			adapter, err = environment.GetExternalAdapter(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println(adapter.ClusterURL())
		})
	})

	Describe("Alerts", func() {
		It("Test 1", func() {

		})
	})

	AfterEach(func() {
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
