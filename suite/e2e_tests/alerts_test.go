package e2e_tests

import (
	. "github.com/onsi/ginkgo"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("OCR @alerts suite", func() {
	var testSetup *testcommon.OCRSetupInputs
	//rulesFilePath := "../ocr-telemetry-prometheus-exporter/ocr.rules.yml"
	//file, err := os.Open(rulesFilePath)
	//Expect(err).ShouldNot(HaveOccurred())
	//rules := map[string]*os.File{"ocrRulesYml": file}

	BeforeEach(func() {
		testSetup = &testcommon.OCRSetupInputs{}
	})

	Describe("Telemetry Down Alerts", func() {
		FIt("Doesn't start the OCR protocol", func() {
			testcommon.NewOCRSetupInputForAtlas(testSetup, 3, 1)

		})

	})

	AfterEach(func() {
		//By("Stop chaos", func() {
		//	err := testSetup.SuiteSetup.Environment().StopAllChaos()
		//	Expect(err).ShouldNot(HaveOccurred())
		//})
		//By("Tearing down the environment", testSetup.SuiteSetup.TearDown())
	})
})
