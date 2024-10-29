package reports

import (
	"log"
	"time"

	"github.com/prometheus/common/model"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

func SendResultsToLoki(lc *client.LokiPromtailClient, testResults []TestResult) {
	for _, result := range testResults {
		labels := model.LabelSet{
			"job":          model.LabelValue("flaky_tests"),
			"test_name":    model.LabelValue(result.TestName),
			"test_package": model.LabelValue(result.TestPackage),
			// "pass_ratio":   model.LabelValue(fmt.Sprintf("%.2f", result.PassRatio)),
		}
		err := lc.HandleStruct(labels, time.Now(), result)
		if err != nil {
			log.Printf("Failed to send test result to Loki: %v", err)
		}
	}
}
