package volume

import (
	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	volcommon "github.com/smartcontractkit/integrations-framework/suite/volume/common"
	"math/big"
	"time"
)

var _ = Describe("Flux monitor volume tests", func() {
	Describe("round completion times", func() {
		s, err := actions.DefaultLocalSetup(
			environment.NewChainlinkCluster(5),
			client.NewNetworkFromConfig,
		)
		Expect(err).ShouldNot(HaveOccurred())
		nodes, err := environment.GetChainlinkClients(s.Env)
		Expect(err).ShouldNot(HaveOccurred())
		nodeAddresses, err := actions.ChainlinkNodeAddresses(nodes)
		Expect(err).ShouldNot(HaveOccurred())
		adapter, err := environment.GetExternalAdapter(s.Env)
		Expect(err).ShouldNot(HaveOccurred())
		err = adapter.SetVariable(5)
		Expect(err).ShouldNot(HaveOccurred())
		s.Client.ParallelTransactions(true)
		err = actions.FundChainlinkNodes(
			nodes,
			s.Client,
			s.Wallets.Default(),
			big.NewFloat(2),
			nil,
		)
		Expect(err).ShouldNot(HaveOccurred())
		err = s.Client.WaitForEvents()
		Expect(err).ShouldNot(HaveOccurred())

		jobPrefix := "flux_monitor"
		spec := &FluxTestSpec{
			TestSpec: volcommon.TestSpec{
				EnvSetup:                s,
				Nodes:                   nodes,
				Adapter:                 adapter,
				NodesAddresses:          nodeAddresses,
				OnChainCheckAttemptsOpt: retry.Attempts(120),
			},
			AggregatorsNum:      100,
			RequiredSubmissions: 5,
			RestartDelayRounds:  0,
			JobPrefix:           jobPrefix,
			// minimum poll time
			NodePollTimePeriod: 15 * time.Second,
			FluxOptions:        contracts.DefaultFluxAggregatorOptions(),
		}
		s.Client.ParallelTransactions(false)
		ft, err := NewFluxTest(spec)
		Expect(err).ShouldNot(HaveOccurred())

		currentRound := 1
		rounds := 10

		err = ft.checkRoundDataOnChain(currentRound, 5)
		Expect(err).ShouldNot(HaveOccurred())

		Measure("Round completion time percentiles", func(b Benchmarker) {
			newVal, err := ft.Adapter.TriggerValueChange(currentRound)
			Expect(err).ShouldNot(HaveOccurred())
			err = ft.checkRoundDataOnChain(currentRound+1, newVal)
			Expect(err).ShouldNot(HaveOccurred())
			err = ft.roundsMetrics(
				currentRound+1,
				spec.RequiredSubmissions,
				big.NewInt(int64(newVal)),
			)
			Expect(err).ShouldNot(HaveOccurred())
			cpu, mem, err := ft.Prom.ResourcesSummary()
			b.RecordValue("CPU", cpu)
			b.RecordValue("MEM", mem)
			Expect(err).ShouldNot(HaveOccurred())
			currentRound += 1
		}, rounds)

		AfterSuite(func() {
			By("Calculating round percentiles", func() {
				percs, err := ft.CalculatePercentiles(ft.roundsDurationData)
				Expect(err).ShouldNot(HaveOccurred())
				ft.PrintPercentileMetrics(percs)
			})
			By("Tearing down the environment", s.TearDown())
		})
	})
})
