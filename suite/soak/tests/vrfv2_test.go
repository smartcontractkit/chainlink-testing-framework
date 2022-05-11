package soak

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

// VRFV2SoakTestJobInfo defines a jobs into and proving key info
type VRFV2SoakTestJobInfo struct {
	Job            *client.Job
	ProvingKey     [2]*big.Int
	ProvingKeyHash [32]byte
}

var _ = Describe("Vrfv2 soak test suite @soak_vrfv2", func() {
	var (
		err           error
		env           *environment.Environment
		vrfv2SoakTest *testsetups.VRFV2SoakTest
		coordinator   contracts.VRFCoordinatorV2
		consumer      contracts.VRFConsumerV2
		jobInfo       []VRFV2SoakTestJobInfo
	)
	local_runner := len(os.Getenv("LOCAL_RUNNER")) > 0

	BeforeEach(func() {
		By("Deploying the environment", func() {
			if local_runner {
				log.Info().Str("Where", "Locally").Msg("Runner running")
				env, err = environment.DeployOrLoadEnvironment(
					environment.NewChainlinkConfig(
						config.ChainlinkVals(),
						"vrfv2-local-soak",
						// works only on perf Geth
						environment.PerformanceGeth,
					),
					tools.ChartsRoot,
				)
			} else {
				log.Info().Str("Where", "Remotely").Msg("Runner running")
				env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
					tools.ChartsRoot,
					"/root/test-env.json", // Default location for the soak-test-runner container
				)
			}
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Vrfv2 test", func() {
			vrfv2SoakTest = testsetups.NewVRFV2SoakTest(
				&testsetups.VRFV2SoakTestInputs{
					TestDuration:         time.Minute * 1,
					ChainlinkNodeFunding: big.NewFloat(1000),
					StopTestOnError:      false,

					RequestsPerSecond:  25,
					ReadEveryNRequests: 1,

					// Make the test simple and just request randomness and return any errors
					TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) error {
						words := uint32(10)
						err := consumer.RequestRandomness(jobInfo[0].ProvingKeyHash, 1, 1, 300000, words)
						return err
					},
				})
			vrfv2SoakTest.Setup(env, local_runner)

			// With the environment setup now we can deploy contracts and jobs
			contractDeployer, err := contracts.NewContractDeployer(vrfv2SoakTest.DefaultNetwork)
			Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")

			// Deploy LINK
			linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

			// Fund Chainlink nodes
			err = actions.FundChainlinkNodes(vrfv2SoakTest.ChainlinkNodes, vrfv2SoakTest.DefaultNetwork, vrfv2SoakTest.Inputs.ChainlinkNodeFunding)
			Expect(err).ShouldNot(HaveOccurred())

			coordinator, consumer = actions.DeployVrfv2Contracts(linkTokenContract, contractDeployer, vrfv2SoakTest.Networks)
			jobs, provingKeys := actions.CreateVrfV2Jobs(vrfv2SoakTest.ChainlinkNodes, coordinator)
			Expect(len(jobs)).Should(Equal(len(provingKeys)), "Should have a set of keys for each job")

			// Create proving key hash here so we aren't calculating it in the test run itself.
			for i, pk := range provingKeys {
				keyHash, err := coordinator.HashOfKey(context.Background(), pk)
				Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a keyHash from the proving keys")
				ji := VRFV2SoakTestJobInfo{
					Job:            jobs[i],
					ProvingKey:     provingKeys[i],
					ProvingKeyHash: keyHash,
				}
				jobInfo = append(jobInfo, ji)
			}

			err = vrfv2SoakTest.DefaultNetwork.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Describe("Run the test", func() {
		It("Makes requests for randomness and verifies number of jobs have been run", func() {
			vrfv2SoakTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if local_runner {
				_, nets, _, _ := vrfv2SoakTest.TearDownVals()
				err = actions.TeardownSuite(env, nets, utils.ProjectRoot, nil, nil)
				Expect(err).ShouldNot(HaveOccurred())
			} else if err := actions.TeardownRemoteSuite(vrfv2SoakTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
		})
	})
})
