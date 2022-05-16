package soak_runner

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
)

var _ = Describe("OCR Soak Test @soak-ocr", func() {
	var (
		err         error
		env         *environment.Environment
		ocrSoakTest *testsetups.OCRSoakTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
				"/root/test-env.json", // Default location for the soak-test-runner container
			)
			Expect(err).ShouldNot(HaveOccurred(), "Failed to connect to running soak environment")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setting up Soak Test", func() {
			ocrSoakTest = testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
				TestDuration:         time.Hour * 168,
				NumberOfContracts:    2,
				ChainlinkNodeFunding: big.NewFloat(10),
				ExpectedRoundTime:    time.Minute,
				RoundTimeout:         time.Minute * 10,
				TimeBetweenRounds:    50 * time.Millisecond,
				StartingAdapterValue: 5,
			})
			ocrSoakTest.Setup(env)
		})
	})

	Describe("With soak test contracts deployed", func() {
		It("runs the soak test until error or timeout", func() {
			ocrSoakTest.Run()
		})
	})

	AfterEach(func() {
		if err = actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals()); err != nil {
			log.Error().Err(err).Msg("Error when tearing down remote suite")
		}
		log.Info().Msg("Soak Test Concluded")
	})
})
