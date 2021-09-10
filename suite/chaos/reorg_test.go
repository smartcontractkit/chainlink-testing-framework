package chaos

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/suite/testcommon"
	"time"
)

var _ = Describe("Reorg example test @reorg", func() {
	i := &testcommon.RunlogSetupInputs{}
	It("Performs reorg and verifies it", func() {
		testcommon.SetupRunlogEnv(i)
		rc, err := NewReorgConfirmer(i.S.Client, i.S.Env)
		Expect(err).ShouldNot(HaveOccurred())
		err = rc.Fork(120 * time.Second)
		Expect(err).ShouldNot(HaveOccurred())
		err = rc.Verify(0, 1)
		Expect(err).ShouldNot(HaveOccurred())
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
