package smoke

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("OCR Feed @ocr", func() {
	var i *actions.OCRSetupInputs

	DescribeTable("Deploys and watches an OCR feed @ocr", func(
		envInit environment.K8sEnvSpecInit,
	) {
		i = &actions.OCRSetupInputs{}
		By("Deploying environment", actions.DeployOCRForEnv(i, envInit))
		By("Funding nodes", actions.FundNodes(i))
		By("Deploying OCR contracts", actions.DeployOCRContracts(i, 1))
		By("Creating OCR jobs", actions.CreateOCRJobs(i))
		By("Checking OCR rounds", actions.CheckRound(i))
	},
		Entry("all the same version", environment.NewChainlinkCluster(6)),
		Entry("different versions", environment.NewMixedVersionChainlinkCluster(6, 2)),
	)

	AfterEach(func() {
		By("Calculating gas costs", func() {
			i.NetworkInfo.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	})
})
