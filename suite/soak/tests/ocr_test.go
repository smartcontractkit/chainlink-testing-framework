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

var _ = Describe("OCR Soak Test @soak-ocr", func() {
	var (
		err         error
		env         *environment.Environment
		ocrSoakTest *testsetups.OCRSoakTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
				tools.ChartsRoot,
				"/root/test-env.json", // Default location for the soak-test-runner container
			)
			Expect(err).ShouldNot(HaveOccurred(), "Failed to connect to running soak environment")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setting up Soak Test", func() {
			ocrSoakTest = testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
				TestDuration:         time.Hour * 4,
				NumberOfContracts:    4,
				ChainlinkNodeFunding: big.NewFloat(1),
				RoundTimeout:         time.Minute * 1,
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
