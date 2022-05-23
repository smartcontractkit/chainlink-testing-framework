package soak

//revive:disable:dot-imports
import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
)

var _ = Describe("OCR Soak Test @soak-ocr", func() {
	var (
		err         error
		env         *environment.Environment
		ocrSoakTest *testsetups.OCRSoakTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env = environment.New(&environment.Config{InsideK8s: true})
			err = env.
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(ethereum.New(nil)).
				AddHelm(chainlink.New(0, nil)).
				Run()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Namespace", env.Cfg.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setting up Soak Test", func() {
			ocrSoakTest = testsetups.NewOCRSoakTest(&testsetups.OCRSoakTestInputs{
				TestDuration:         time.Hour * 1,
				NumberOfContracts:    4,
				ChainlinkNodeFunding: big.NewFloat(1),
				ExpectedRoundTime:    time.Minute,
				RoundTimeout:         time.Minute * 10,
				TimeBetweenRounds:    time.Minute,
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
