package chaos

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/chaos"
	"github.com/smartcontractkit/integrations-framework/chaos/experiments"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
	"time"
)

var _ = XDescribeTable("OCR chaos tests @chaos-ocr", func(
	envName string,
	envInit environment.K8sEnvSpecInit,
	chaosSpec chaos.Experimentable,
) {
	i := &testcommon.OCRSetupInputs{}
	Context("Runs OCR test with a chaos modifier", func() {
		testcommon.DeployOCRForEnv(i, envName, envInit)
		testcommon.SetupOCRTest(i)
		_, err := i.SuiteSetup.Env.ApplyChaos(chaosSpec)
		Expect(err).ShouldNot(HaveOccurred())
		testcommon.CheckRound(i)
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.SuiteSetup.Env.StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	})
},
	Entry("One node pod failure",
		"basic-chainlink",
		environment.NewChainlinkCluster(5),
		&experiments.PodFailure{
			LabelKey:   "app",
			LabelValue: "chainlink-0",
			Duration:   10 * time.Second,
		}),
)
