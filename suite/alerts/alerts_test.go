package alerts

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/integrations-framework/chaos"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/steps"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("OCR Alerts suite", func() {
	var (
		testSetup *testcommon.OCRSetupInputs
	)

	BeforeEach(func() {
		By("Deploying the basic environment", func() {
			testSetup = &testcommon.OCRSetupInputs{}
			testcommon.DeployOCRForEnv(
				testSetup,
				"basic-chainlink",
				environment.NewChainlinkClusterForAlertsTesting(6),
			)
			testcommon.SetupOCRTest(testSetup)
		})

		By("Creating initializer file for mockserver", func() {
			steps.WriteDataForOTPEToInitializerFileForMockserver(
				testSetup.OCRInstance.Address(),
				testSetup.ChainlinkNodes,
			)
		})

		By("Deploying otpe and prometheus", func() {
			err := testSetup.SuiteSetup.Env.DeploySpecs(environment.OtpeGroup())
			Expect(err).ShouldNot(HaveOccurred())

			err = testSetup.SuiteSetup.Env.DeploySpecs(environment.PrometheusGroup())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Telemetry Down Alerts", func() {
		It("Doesn't start the OCR protocol", func() {
			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "5m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("Telemetry Down (infra)")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})

		It("Shuts down all chainlink nodes after some successful rounds", func() {
			testcommon.SendOCRJobs(testSetup)
			testcommon.CheckRound(testSetup)

			experimentSpecs := make([]chaos.Experimentable, len(testSetup.ChainlinkNodes[1:]))

			for i := 0; i < len(testSetup.ChainlinkNodes[1:]); i++ {
				experimentSpecs[i] = &experiments.PodKill{
					TargetAppLabel: fmt.Sprintf("chainlink-%d", i+1),
				}
			}

			for _, experimentSpec := range experimentSpecs {
				_, err := testSetup.SuiteSetup.Env.ApplyChaos(experimentSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}

			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "5m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("Telemetry Down (infra)")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})
	})

	AfterEach(func() {
		By("Stop chaos", func() {
			err := testSetup.SuiteSetup.Env.StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", testSetup.SuiteSetup.TearDown())
	})
})
