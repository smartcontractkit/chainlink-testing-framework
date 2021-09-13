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

		reorgConfirmer, err := NewReorgConfirmer(
			i.S.Client,
			i.S.Env,
			5,
			10,
			time.Second*600,
		)
		Expect(err).ShouldNot(HaveOccurred())
		i.S.Client.AddHeaderEventSubscription("reorg", reorgConfirmer)
		err = i.S.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
		err = reorgConfirmer.Verify()
		Expect(err).ShouldNot(HaveOccurred())
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.S.Env.StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.S.TearDown())
	})
})
