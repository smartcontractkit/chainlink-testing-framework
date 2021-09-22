package alerts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("Alerts suite", func() {
	Describe("Alerts", func() {
		It("Deploys the alerts stack up to OTPE", func() {
			i := &testcommon.OCRSetupInputs{}
			testcommon.DeployOCRForEnv(i, "basic-chainlink", environment.NewChainlinkClusterForAlertsTesting(5))
			testcommon.SetupOCRTest(i)
			testcommon.CheckRound(i)
			testcommon.WriteDataForOTPEToInitializerFileForMockserver(i)

			err := i.SuiteSetup.Env.DeploySpecs(environment.OtpeGroup())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
