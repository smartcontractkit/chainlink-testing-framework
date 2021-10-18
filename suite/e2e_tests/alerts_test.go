package e2e_tests

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
	"os"
)

var _ = Describe("OCR @alerts suite", func() {
	var testSetup *testcommon.OCRSetupInputs
	rulesFilePath := "../ocr-telemetry-prometheus-exporter/ocr.rules.yml"
	file, err := os.Open(rulesFilePath)
	Expect(err).ShouldNot(HaveOccurred())
	rules := map[string]*os.File{"ocrRulesYml": file}

	BeforeEach(func() {
		testSetup = &testcommon.OCRSetupInputs{}
		testcommon.NewOCRSetupInputForObservability(testSetup, 6, rules)
	})

	Describe("Telemetry Down Alerts", func() {
		It("Doesn't start the OCR protocol", func() {
			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "7m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("Telemetry Down (infra)")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})
	})

	AfterEach(func() {
		By("Stop chaos", func() {
			err := testSetup.SuiteSetup.Environment().StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", testSetup.SuiteSetup.TearDown())
	})
})
