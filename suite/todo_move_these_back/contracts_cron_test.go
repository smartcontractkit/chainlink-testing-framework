package contracts

import (
	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
)

var _ = Describe("Cronjob suite", func() {
	var s *actions.DefaultSuiteSetup
	var err error

	DescribeTable("use cron job", func(
		envInitFunc environment.K8sEnvSpecInit,
		networkInitFunc client.BlockchainNetworkInit,
		ocrOptions contracts.OffchainOptions,
	) {
		s, err = actions.DefaultLocalSetup(envInitFunc, networkInitFunc)
		Expect(err).ShouldNot(HaveOccurred())
		chainlinkNodes, err := environment.GetChainlinkClients(s.Env)
		Expect(err).ShouldNot(HaveOccurred())
		adapter, err := environment.GetExternalAdapter(s.Env)
		Expect(err).ShouldNot(HaveOccurred())

		os := &client.PipelineSpec{
			URL:         adapter.ClusterURL() + "/five",
			Method:      "POST",
			RequestData: "{}",
			DataPath:    "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		job, err := chainlinkNodes[0].CreateJob(&client.CronJobSpec{
			Schedule:          "CRON_TZ=UTC * * * * * *",
			ObservationSource: ost,
		})
		Expect(err).ShouldNot(HaveOccurred())
		err = retry.Do(func() error {
			jobRuns, err := chainlinkNodes[0].ReadRunsByJob(job.Data.ID)
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
	},
		Entry(
			"on Ethereum Hardhat",
			environment.NewChainlinkCluster(1),
			client.NewNetworkFromConfig,
			contracts.DefaultOffChainAggregatorOptions(),
		),
	)

	AfterEach(func() {
		s.Env.TearDown()
	})
})
