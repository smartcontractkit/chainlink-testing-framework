package volume

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	volumecommon "github.com/smartcontractkit/integrations-framework/suite/volume/common"
	"github.com/smartcontractkit/integrations-framework/tools"
	"golang.org/x/sync/errgroup"
	"math/big"
	"sync"
	"time"
)

// FluxTestSpec flux aggregator volume test spec
type FluxTestSpec struct {
	volumecommon.TestSpec
	AggregatorsNum              int
	RequiredSubmissions         int
	RestartDelayRounds          int
	JobPrefix                   string
	ObservedValueChangeInterval time.Duration
	NodePollTimePeriod          time.Duration
	FluxOptions                 contracts.FluxAggregatorOptions
}

// FluxTest flux test runtime data
type FluxTest struct {
	volumecommon.Test
	// round durations, calculated as a difference from earliest chainlink run for contract across all nodes,
	// until all confirmations are found on-chain (block_timestamp)
	roundsDurationData []time.Duration
	FluxInstances      []contracts.FluxAggregator
	ContractsToJobsMap map[string][]contracts.JobByInstance
	NodesByHostPort    map[string]client.Chainlink
}

// FluxInstanceDeployment data required by flux instance to calculate per round metrics
type FluxInstanceDeployment struct {
	volumecommon.InstanceDeployment
	Spec              *FluxTestSpec
	ContractToJobsMap map[string][]contracts.JobByInstance
	NodesByHostPort   map[string]client.Chainlink
}

// NewFluxTest deploys AggregatorsNum flux aggregators concurrently
func NewFluxTest(spec *FluxTestSpec) (*FluxTest, error) {
	fluxInstances := make([]contracts.FluxAggregator, 0)
	nodesByHostPort := make(map[string]client.Chainlink)
	contractToJobsMap := make(map[string][]contracts.JobByInstance)
	mu := &sync.Mutex{}
	g := &errgroup.Group{}
	for i := 0; i < spec.AggregatorsNum; i++ {
		d := &FluxInstanceDeployment{
			InstanceDeployment: volumecommon.InstanceDeployment{
				Index:         i,
				Suite:         spec.EnvSetup,
				NodeAddresses: spec.NodesAddresses,
				Nodes:         spec.Nodes,
				Adapter:       spec.Adapter,
			},
			Spec:              spec,
			NodesByHostPort:   nodesByHostPort,
			ContractToJobsMap: contractToJobsMap,
		}
		g.Go(func() error {
			log.Info().Int("Instance ID", d.Index).Msg("Deploying contracts instance")
			fluxInstance, err := d.Suite.Deployer.DeployFluxAggregatorContract(d.Suite.Wallets.Default(), d.Spec.FluxOptions)
			if err != nil {
				return err
			}
			err = fluxInstance.Fund(d.Suite.Wallets.Default(), big.NewFloat(0), big.NewFloat(1))
			if err != nil {
				return err
			}
			err = fluxInstance.UpdateAvailableFunds(context.Background(), d.Suite.Wallets.Default())
			if err != nil {
				return err
			}
			// set oracles and submissions
			err = fluxInstance.SetOracles(d.Suite.Wallets.Default(),
				contracts.FluxAggregatorSetOraclesOptions{
					AddList:            d.NodeAddresses,
					RemoveList:         []common.Address{},
					AdminList:          d.NodeAddresses,
					MinSubmissions:     uint32(d.Spec.RequiredSubmissions),
					MaxSubmissions:     uint32(d.Spec.RequiredSubmissions),
					RestartDelayRounds: uint32(d.Spec.RestartDelayRounds),
				})
			if err != nil {
				return err
			}
			for _, n := range d.Nodes {
				fluxSpec := &client.FluxMonitorJobSpec{
					Name:            fmt.Sprintf("%s_%d", d.Spec.JobPrefix, d.Index),
					ContractAddress: fluxInstance.Address(),
					PollTimerPeriod: d.Spec.NodePollTimePeriod,
					// it's crucial not to skew rounds schedule for that particular volume test
					IdleTimerDisabled: true,
					PollTimerDisabled: false,
					ObservationSource: client.ObservationSourceSpec(d.Adapter.ClusterURL() + "/variable"),
				}
				job, err := n.CreateJob(fluxSpec)
				if err != nil {
					return err
				}
				mu.Lock()
				d.NodesByHostPort[n.URL()] = n
				d.ContractToJobsMap[fluxInstance.Address()] = append(d.ContractToJobsMap[fluxInstance.Address()],
					contracts.JobByInstance{
						ID:       job.Data.ID,
						Instance: n.URL(),
					})
				mu.Unlock()
			}
			mu.Lock()
			fluxInstances = append(fluxInstances, fluxInstance)
			mu.Unlock()
			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("Contracts to jobs", contractToJobsMap).Msg("Debug data for per round metrics")
	prom, err := tools.NewPrometheusClient(spec.EnvSetup.Config.Prometheus.URL)
	if err != nil {
		return nil, err
	}
	return &FluxTest{
		Test: volumecommon.Test{
			DefaultSetup:            spec.EnvSetup,
			OnChainCheckAttemptsOpt: spec.OnChainCheckAttemptsOpt,
			Nodes:                   spec.Nodes,
			Adapter:                 spec.Adapter,
			Prom:                    prom,
		},
		FluxInstances:      fluxInstances,
		ContractsToJobsMap: contractToJobsMap,
		NodesByHostPort:    nodesByHostPort,
		roundsDurationData: make([]time.Duration, 0),
	}, nil
}

// roundsStartTimes gets run start time for every contract, ns
func (vt *FluxTest) roundsStartTimes() (map[string]int64, error) {
	mu := &sync.Mutex{}
	g := &errgroup.Group{}
	contractsStartTimes := make(map[string]int64)
	for contractAddr, jobs := range vt.ContractsToJobsMap {
		contractAddr := contractAddr
		jobs := jobs
		g.Go(func() error {
			if err := actions.GetRoundStartTimesAcrossNodes(contractAddr, jobs, vt.NodesByHostPort, mu, contractsStartTimes); err != nil {
				return err
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	for contract, rst := range contractsStartTimes {
		log.Debug().Str("Contract", contract).Int64("Start time", rst/1e9).Send()
	}
	log.Debug().Interface("Round start times", contractsStartTimes).Send()
	return contractsStartTimes, nil
}

// checkRoundSubmissionsEvents checks whether all submissions is found, gets last event block as round end time
func (vt *FluxTest) checkRoundSubmissionsEvents(ctx context.Context, roundID int, submissions int, submissionVal *big.Int) (map[string]int64, error) {
	g := errgroup.Group{}
	mu := &sync.Mutex{}
	endTimes := make(map[string]int64)
	for _, fi := range vt.FluxInstances {
		fi := fi
		g.Go(func() error {
			err := actions.GetRoundCompleteTimestamps(fi, roundID, submissions, submissionVal, mu, endTimes, vt.DefaultSetup.Client)
			if err != nil {
				return err
			}
			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return nil, err
	}
	return endTimes, nil
}

// roundsMetrics gets start times from runs via API, awaits submissions for each contract and get round end times
func (vt *FluxTest) roundsMetrics(roundID int, submissions int, submissionVal *big.Int) error {
	startTimes, err := vt.roundsStartTimes()
	if err != nil {
		return err
	}
	endTimes, err := vt.checkRoundSubmissionsEvents(context.Background(), roundID, submissions, submissionVal)
	if err != nil {
		return err
	}
	for contract := range startTimes {
		duration := time.Duration(endTimes[contract] - startTimes[contract])
		vt.roundsDurationData = append(vt.roundsDurationData, duration)
		log.Info().Str("Contract", contract).Str("Round duration", duration.String()).Send()
	}
	return nil
}

// checkRoundDataOnChain check that ContractData about last round is correct
func (vt *FluxTest) checkRoundDataOnChain(roundID int, newVal int) error {
	if err := retry.Do(func() error {
		finished, err := vt.isContractDataOnChain(roundID, newVal)
		if err != nil {
			return err
		}
		if !finished {
			return errors.New("round is not finished")
		}
		return nil
	}, vt.OnChainCheckAttemptsOpt); err != nil {
		return errors.Wrap(err, "round is not fully finished on chain")
	}
	return nil
}

// isContractDataOnChain check all answers in particular round is correct
func (vt *FluxTest) isContractDataOnChain(roundID int, newVal int) (bool, error) {
	log.Debug().Int("Round ID", roundID).Msg("Checking round completion on chain")
	var rounds []contracts.RoundData
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(len(vt.FluxInstances))
	for _, flux := range vt.FluxInstances {
		flux := flux
		go func() {
			defer wg.Done()
			cd, err := flux.LatestRoundData(context.Background())
			if err != nil {
				log.Err(err).Msg("Failed to get contract data")
				return
			}
			mu.Lock()
			rounds = append(rounds, cd)
			mu.Unlock()
		}()
	}
	wg.Wait()
	for _, r := range rounds {
		if r.RoundId.Int64() != int64(roundID) || r.Answer.Int64() != int64(newVal) {
			return false, nil
		}
	}
	return true, nil
}
