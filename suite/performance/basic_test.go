package performance

//revive:disable:dot-imports
import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
)

var _ = Describe("Profiling suite @profile", func() {
	var (
		err             error
		nets            *blockchain.Networks
		chainlinkNodes  []client.Chainlink
		mockserver      *client.MockserverClient
		testEnvironment *environment.Environment
		profileTest     *testsetups.ChainlinkProfileTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			testEnvironment, err = environment.DeployOrLoadEnvironment(
				environment.NewChainlinkConfig(
					config.ChainlinkVals(),
					"chainlink-profiling",
					config.GethNetworks()...,
				),
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			err = testEnvironment.ConnectAll()
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to all nodes shouldn't fail")
		})

		By("Setting up the test", func() {
			chainlinkNodes, err = client.ConnectChainlinkNodes(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
			mockserver, err = client.ConnectMockServer(testEnvironment)
			Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")

			profileFunction := func(chainlinkNode client.Chainlink) {
				defer GinkgoRecover()
				bta := client.BridgeTypeAttributes{
					Name:        fmt.Sprintf("variable-%s", uuid.NewV4().String()),
					URL:         fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
					RequestData: "{}",
				}
				err = chainlinkNode.CreateBridge(&bta)
				Expect(err).ShouldNot(HaveOccurred(), "Creating bridge in chainlink node shouldn't fail")

				_, err = chainlinkNode.CreateJob(&client.CronJobSpec{
					Schedule:          "CRON_TZ=UTC * * * * * *",
					ObservationSource: client.ObservationSourceSpecBridge(bta),
				})
				Expect(err).ShouldNot(HaveOccurred(), "Creating Cron Job in chainlink node shouldn't fail")
			}

			profileTest = testsetups.NewChainlinkProfileTest(testsetups.ChainlinkProfileTestInputs{
				ProfileFunction: profileFunction,
				ProfileDuration: time.Second,
				ChainlinkNodes:  chainlinkNodes,
			})
			profileTest.Setup(testEnvironment)
		})
	})

	Describe("checking Chainlink node's PPROF", func() {
		It("queries PPROF on node", func() {
			err = mockserver.SetValuePath("/variable", 5)
			Expect(err).ShouldNot(HaveOccurred(), "Setting value path in mockserver shouldn't fail")

			profileTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(testEnvironment, nets, utils.ProjectRoot, chainlinkNodes, &profileTest.TestReporter)
			Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
		})
	})
})
