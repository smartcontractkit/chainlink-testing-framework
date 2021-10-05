package observability_test

import (
	"github.com/montanaflynn/stats"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
	"math/rand"
	"time"
)

var _ = Describe("OTPE metrics suite", func() {
	var testSetup *testcommon.OCRSetupInputs

	BeforeEach(func() {
		testSetup = &testcommon.OCRSetupInputs{}
		testcommon.NewOCRSetupInputForObservability(testSetup, 6)
	})

	It("Computes correctly the median feed price from all nodes in current round", func() {
		By("Initializing mockserver adapter endpoints with random data")

		rand.Seed(time.Now().UnixNano())
		min := 10000
		max := 10100
		var adapterResults []int
		for index := 1; index < len(testSetup.ChainlinkNodes); index++ {
			result := rand.Intn(max-min+1) + min
			adapterResults = append(adapterResults, result)
		}
		median, _ := stats.Median(stats.LoadRawData(adapterResults))
		testcommon.SetAdapterResults(testSetup, adapterResults)
		testcommon.SendOCRJobs(testSetup)

		By("Kicking off first round")

		testcommon.StartNewRound(testSetup, 1)

		By("Comparing the metrics value with the expected value")

		prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Env)
		Expect(err).ShouldNot(HaveOccurred())

		Eventually(func(g Gomega) int {
			value, err := prometheus.GetQuery("ocr_telemetry_feed_message_report_req_median")
			g.Expect(err).ShouldNot(HaveOccurred())
			return len(value.(model.Vector))
		}, "2m", "5s").Should(BeNumerically(">", 0))

		Eventually(func(g Gomega) model.SampleValue {
			value, err := prometheus.GetQuery("ocr_telemetry_feed_message_report_req_median")
			g.Expect(err).ShouldNot(HaveOccurred())
			return value.(model.Vector)[0].Value
		}, "2m", "5s").Should(BeNumerically("==", median))
	})

	AfterEach(func() {
		By("Stop chaos", func() {
			err := testSetup.SuiteSetup.Env.StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", testSetup.SuiteSetup.TearDown())
	})
})
