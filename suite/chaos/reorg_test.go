package chaos

import (
	"github.com/smartcontractkit/integrations-framework/actions"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reorg example test @reorg", func() {
	i := &actions.RunlogSetupInputs{}
	It("Performs reorg and verifies it", func() {
		By("Deploying the environment", actions.SetupRunlogEnv(i))

		reorgConfirmer, err := NewReorgConfirmer(
			i.SuiteSetup.DefaultNetwork().Client,
			i.SuiteSetup.Environment(),
			5,
			10,
			time.Second*600,
		)
		Expect(err).ShouldNot(HaveOccurred())
		i.SuiteSetup.DefaultNetwork().Client.AddHeaderEventSubscription("reorg", reorgConfirmer)
		err = i.SuiteSetup.DefaultNetwork().Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())
		err = reorgConfirmer.Verify()
		Expect(err).ShouldNot(HaveOccurred())
	})
	AfterEach(func() {
		By("Restoring chaos", func() {
			err := i.SuiteSetup.Environment().StopAllChaos()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Tearing down the environment", i.SuiteSetup.TearDown())
	})
})
