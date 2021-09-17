package contracts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
)

var _ = Describe("OCR Feed @ocr", func() {

	DescribeTable("Deploys and watches an OCR feed @ocr", func(
		envName string,
		envInit environment.K8sEnvSpecInit,
	) {
		i := &testcommon.OCRSetupInputs{}
		testcommon.DeployOCRForEnv(i, envName, envInit)
		testcommon.SetupOCRTest(i)
		testcommon.CheckRound(i)
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	},
		Entry("all the same version", "basic-chainlink", environment.NewChainlinkCluster(5)),
		Entry("different versions", "mixed-version-chainlink", environment.NewMixedVersionChainlinkCluster(5, 2)),
	)
})
