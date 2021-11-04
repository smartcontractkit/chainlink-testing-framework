package chaos

import (
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client/chaos"
	"github.com/smartcontractkit/integrations-framework/client/chaos/experiments"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = XDescribeTable("OCR chaos tests @chaos-ocr", func(
	envInit environment.K8sEnvSpecInit,
	chaosSpec chaos.Experimentable,
) {
	i := &actions.OCRSetupInputs{}
	Context("Runs OCR test with a chaos modifier", func() {
		actions.DeployOCRForEnv(i, envInit)
		actions.FundNodes(i)
		actions.DeployOCRContracts(i, 1)
		actions.CreateOCRJobs(i)
		_, err := i.SuiteSetup.Environment().ApplyChaos(chaosSpec)
		Expect(err).ShouldNot(HaveOccurred())
		actions.CheckRound(i)
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.SuiteSetup.Environment().StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	})
},
	Entry("One node pod failure",
		environment.NewChainlinkCluster(5),
		&experiments.PodFailure{
			LabelKey:   "app",
			LabelValue: "chainlink-0",
			Duration:   10 * time.Second,
		}),
)
