package smoke

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/utils"
	"path/filepath"
)

var _ = Describe("Cronjob suite @cron", func() {
	var (
		err        error
		job        *client.Job
		cls        []client.Chainlink
		mockserver *client.MockserverClient
		e          *environment.Environment
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.NewEnvironmentFromPreset(filepath.Join(utils.PresetRoot, "chainlink-cluster-6"))
			Expect(err).ShouldNot(HaveOccurred())
			err = e.Connect()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			mockserver, err = client.NewMockServerClientFromEnv(e)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Adding cron job to a node", func() {
			err = mockserver.SetVariable(5)
			Expect(err).ShouldNot(HaveOccurred())

			bta := client.BridgeTypeAttributes{
				Name:        fmt.Sprintf("variable-%s", uuid.NewV4().String()),
				URL:         fmt.Sprintf("%s/variable", mockserver.Config.ClusterURL),
				RequestData: "{}",
			}
			err = cls[0].CreateBridge(&bta)
			Expect(err).ShouldNot(HaveOccurred())

			job, err = cls[0].CreateJob(&client.CronJobSpec{
				Schedule:          "CRON_TZ=UTC * * * * * *",
				ObservationSource: client.ObservationSourceSpecBridge(bta),
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with Cron job", func() {
		It("runs 5 times with no errors", func() {
			Eventually(func(g Gomega) {
				jobRuns, err := cls[0].ReadRunsByJob(job.Data.ID)
				g.Expect(err).ShouldNot(HaveOccurred())

				g.Expect(len(jobRuns.Data)).Should(BeNumerically("==", 5))

				for _, jr := range jobRuns.Data {
					g.Expect(jr.Attributes.Errors).Should(Equal([]interface{}{nil}))
				}
			}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
