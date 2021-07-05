package suite

import (
	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/suite"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("Cronjob suite", func() {
	DescribeTable("use cron job", func(
		initFunc client.BlockchainNetworkInit,
		ocrOptions contracts.OffchainOptions,
	) {
		_, err := suite.DefaultLocalSetup(initFunc)
		Expect(err).ShouldNot(HaveOccurred())
		chainlinkNodes, _, err := suite.ConnectToTemplateNodes()
		Expect(err).ShouldNot(HaveOccurred())

		adapter := tools.NewExternalAdapter()

		job, err := chainlinkNodes[0].CreateJob(&client.CronJobSpec{
			Schedule:          "CRON_TZ=UTC * * * * * *",
			ObservationSource: client.ObservationSourceSpec(adapter.InsideDockerAddr + "/five"),
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
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultOffChainAggregatorOptions()),
	)
})
