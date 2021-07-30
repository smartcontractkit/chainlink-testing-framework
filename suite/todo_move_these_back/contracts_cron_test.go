package contracts

import (
	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Cronjob suite", func() {
	var (
		s       *actions.DefaultSuiteSetup
		adapter environment.ExternalAdapter
		nodes   []client.Chainlink
		job     *client.Job
		err     error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			s, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(1),
				client.NewNetworkFromConfig,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
			adapter, err = environment.GetExternalAdapter(s.Env)
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Adding cron job to a node", func() {
			os := &client.PipelineSpec{
				URL:         adapter.ClusterURL() + "/five",
				Method:      "GET",
				RequestData: "{}",
				DataPath:    "data,result",
			}
			ost, err := os.String()
			Expect(err).ShouldNot(HaveOccurred())
			job, err = nodes[0].CreateJob(&client.CronJobSpec{
				Schedule:          "CRON_TZ=UTC * * * * * *",
				ObservationSource: ost,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("with Cron job", func() {
		It("runs 5 times with no errors", func() {
			err = retry.Do(func() error {
				jobRuns, err := nodes[0].ReadRunsByJob(job.Data.ID)
				if err != nil {
					return err
				}
				if len(jobRuns.Data) != 5 {
					return errors.New("not all jobs are completed")
				}
				for _, jr := range jobRuns.Data {
					Expect(jr.Attributes.Errors).Should(Equal([]interface{}{nil}))
				}
				return nil
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", s.TearDown())
	})
})
