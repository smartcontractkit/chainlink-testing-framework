package chaos

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
	"time"
)

var _ = Describe("Reorg example test @reorg", func() {
	i := &testcommon.RunlogSetupInputs{}
	It("Performs reorg", func() {
		testcommon.SetupRunlogTest(i)
		rc, err := NewReorgChecker(i.S.Client, i.S.Env, i.S.Config, environment.MinersRPCPort)
		Expect(err).ShouldNot(HaveOccurred())
		err = rc.Fork(90 * time.Second)
		Expect(err).ShouldNot(HaveOccurred())
		err = rc.Verify()
		Expect(err).ShouldNot(HaveOccurred())
		testcommon.CallRunlogOracle(i)
		testcommon.CheckRunlogCompleted(i)
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.S.Env.StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.S.TearDown())
	})
},
)
