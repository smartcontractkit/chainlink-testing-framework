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
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Vrfv2 test", func() {
			// defer GinkgoRecover()
			vrfv2SoakTest = testsetups.NewVRFV2SoakTest(
				&testsetups.VRFV2SoakTestInputs{
					TestDuration:         time.Minute * 1,
					ChainlinkNodeFunding: big.NewFloat(1000),

					RequestsPerSecond:  1,
					ReadEveryNRequests: 10,
					TestFunc: func(t *testsetups.VRFV2SoakTest, requestNumber int) {
						words := uint32(10)
						err := t.Consumer.RequestRandomness(t.JobInfo[0].ProvingKeyHash, 1, 1, 300000, words)
						Expect(err).ShouldNot(HaveOccurred())
						if requestNumber%t.Inputs.ReadEveryNRequests == 0 {
							Eventually(func(g Gomega) {
								log.Info().Int("Request Number", requestNumber).Msg("Validation attempt for request")
								jobRuns, err := t.ChainlinkNodes[0].ReadRunsByJob(t.JobInfo[0].Job.Data.ID)
								g.Expect(err).ShouldNot(HaveOccurred())
								g.Expect(len(jobRuns.Data)).Should(BeNumerically(">=", requestNumber))
								// randomness, err := t.Consumer.GetAllRandomWords(context.Background(), int(10))
								// g.Expect(err).ShouldNot(HaveOccurred())
								// for _, w := range randomness {
								// 	log.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
								// 	g.Expect(w.Uint64()).Should(Not(BeNumerically("==", 0)), "Expected the VRF job give an answer other than 0")
								// }
							}, time.Minute*5, "1s").Should(Succeed())
						}
						t.NumberRequestsValidated++
					},
				})
			vrfv2SoakTest.Setup(env)
		})
	})
	Describe("Run the test", func() {
		It("Makes requests for randomness and veriies number of jobs have been run", func() {
			defer GinkgoRecover()
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
