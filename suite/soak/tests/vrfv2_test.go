package soak

//revive:disable:dot-imports
import (
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/testsetups"
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
				// "/Users/tateexon/git/integrations-framework/suite/soak/chainlink-soak-gnl54.yaml",
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Vrfv2 test", func() {
			vrfv2SoakTest = testsetups.NewVRFV2SoakTest(
				&testsetups.VRFV2SoakTestInputs{
					TestDuration:            time.Minute * 2,
					RequestsPerSecondWanted: 1,
					ChainlinkNodeFunding:    big.NewFloat(1000),
					RoundTimeout:            time.Minute * 2,
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
