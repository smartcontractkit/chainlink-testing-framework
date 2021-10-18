package e2e_tests

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
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/steps"
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
				_, err := testSetup.SuiteSetup.Environment().ApplyChaos(experimentSpec)
				Expect(err).ShouldNot(HaveOccurred())
			}

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

	Describe("OCR feed close to reporting failure", func() {
		It("Fires when the number of OCR nodes falls below 2*ocr_contract_config_f+1+2 because the oracle is dead", func() {
			testcommon.SendOCRJobs(testSetup)
			testcommon.CheckRound(testSetup)

			experimentSpec := &experiments.PodKill{
				TargetAppLabel: "chainlink-1"}
			_, err := testSetup.SuiteSetup.Environment().ApplyChaos(experimentSpec)
			Expect(err).ShouldNot(HaveOccurred())

			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "7m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("OCR feed close to reporting failure")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})

		It("Fires when the number of OCR nodes falls below 2*ocr_contract_config_f+1+2 because the oracle is "+
			"removed from the contract config", func() {
			testcommon.SendOCRJobs(testSetup)
			testcommon.CheckRound(testSetup)

			err := testSetup.OCRInstance.SetConfig(
				testSetup.DefaultWallet,
				testSetup.ChainlinkNodes[2:],
				contracts.DefaultOffChainAggregatorConfig(len(testSetup.ChainlinkNodes[2:])),
			)
			Expect(err).ShouldNot(HaveOccurred())

			err = testSetup.Mockserver.PutExpectations(steps.GetMockserverInitializerDataForOTPE(
				testSetup.OCRInstance.Address(),
				testSetup.ChainlinkNodes[2:],
			))
			Expect(err).ShouldNot(HaveOccurred())

			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "7m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("OCR feed close to reporting failure")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})
	})

	Describe("No observations from an OCR oracle", func() {
		It("The oracle is down", func() {
			testcommon.SendOCRJobs(testSetup)
			testcommon.CheckRound(testSetup)

			experimentSpec := &experiments.PodKill{
				TargetAppLabel: "chainlink-1"}
			_, err := testSetup.SuiteSetup.Environment().ApplyChaos(experimentSpec)
			Expect(err).ShouldNot(HaveOccurred())

			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "7m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("No observations from an OCR oracle")),
				}),
				"State": Equal(v1.AlertState("firing")),
			})))
		})

		It("The data source is down", func() {
			testcommon.SendOCRJobs(testSetup)
			testcommon.CheckRound(testSetup)

			pathSelector := client.PathSelector{Path: "/node_1"}
			err := testSetup.Mockserver.ClearExpectation(pathSelector)
			Expect(err).ShouldNot(HaveOccurred())

			prometheus, err := environment.GetPrometheusClientFromEnv(testSetup.SuiteSetup.Environment())
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) []v1.Alert {
				alerts, err := prometheus.Alerts(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				return alerts.Alerts
			}, "7m", "15s").Should(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Labels": MatchKeys(IgnoreExtras, Keys{
					model.LabelName("alertname"): Equal(model.LabelValue("No observations from an OCR oracle")),
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
