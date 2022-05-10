package soak

//revive:disable:dot-imports
import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("Vrfv2 soak test suite @soak_vrfv2", func() {
	var (
		err           error
		env           *environment.Environment
		vrfv2SoakTest *testsetups.VRFV2SoakTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
				tools.ChartsRoot,
				"/root/test-env.json", // Default location for the soak-test-runner container
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Vrfv2 test", func() {
			vrfv2SoakTest = testsetups.NewVRFV2SoakTest(
				&testsetups.VRFV2SoakTestInputs{
					TestDuration:         time.Minute * 10,
					ChainlinkNodeFunding: big.NewFloat(1000),
					StopTestOnError:      false,

					RequestsPerSecond:  25,
					ReadEveryNRequests: 1,
					TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) error {
						words := uint32(10)
						err := t.Consumer.RequestRandomness(t.JobInfo[0].ProvingKeyHash, 1, 1, 300000, words)
						return err
					},
				})
			vrfv2SoakTest.Setup(env)
		})
	})
	Describe("Run the test", func() {
		It("Makes requests for randomness and veriies number of jobs have been run", func() {
			vrfv2SoakTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if err := actions.TeardownRemoteSuite(vrfv2SoakTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
		})
	})
})
