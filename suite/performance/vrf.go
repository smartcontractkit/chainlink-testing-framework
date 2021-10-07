package performance

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"golang.org/x/sync/errgroup"
	"math/big"
	"time"
)

// VRFJobMap is a custom map type that holds the record of jobs by the contract instance and the chainlink node
type VRFJobMap map[ConsumerCoordinatorPair]map[client.Chainlink]VRFProvingData

// VRFProvingData proving key and job ID pair
type VRFProvingData struct {
	ProvingKeyHash [32]byte
	JobID          string
}

// ConsumerCoordinatorPair consumer and coordinator pair
type ConsumerCoordinatorPair struct {
	consumer    contracts.VRFConsumer
	coordinator contracts.VRFCoordinator
}

// VRFTestOptions contains the parameters for the VRF soak test to be executed
type VRFTestOptions struct {
	TestOptions
}

// VRFTest is the implementation of Test that will configure and execute soak test
// of VRF contracts & jobs
type VRFTest struct {
	TestOptions VRFTestOptions
	Environment environment.Environment
	Blockchain  client.BlockchainClient
	Wallets     client.BlockchainWallets
	Deployer    contracts.ContractDeployer

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	link              contracts.LinkToken
	vrf               contracts.VRF
	blockHashStore    contracts.BlockHashStore
	contractInstances []ConsumerCoordinatorPair
	adapter           environment.ExternalAdapter

	testResults PerfRequestIDTestResults
	jobMap      VRFJobMap
}

// NewVRFTest creates new VRF performance/soak test
func NewVRFTest(
	testOptions VRFTestOptions,
	env environment.Environment,
	link contracts.LinkToken,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	adapter environment.ExternalAdapter,
) Test {
	return &VRFTest{
		TestOptions: testOptions,
		Environment: env,
		link:        link,
		Blockchain:  blockchain,
		Wallets:     wallets,
		Deployer:    deployer,
		adapter:     adapter,
		testResults: NewPerfRequestIDTestResults(),
		jobMap:      VRFJobMap{},
	}
}

// Setup setups VRF performance/soak test
func (f *VRFTest) Setup() error {
	chainlinkClients, err := environment.GetChainlinkClients(f.Environment)
	if err != nil {
		return err
	}
	nodeAddresses, err := actions.ChainlinkNodeAddresses(chainlinkClients)
	if err != nil {
		return err
	}
	adapter, err := environment.GetExternalAdapter(f.Environment)
	if err != nil {
		return err
	}
	f.chainlinkClients = chainlinkClients
	f.nodeAddresses = nodeAddresses
	f.adapter = adapter
	return f.deployContracts()
}

// deployConsumerCoordinatorPair deploys consumer + coordinator
// VRF coordinator can't register several contracts with the same proving key, so we splitting them to ease metrics aggregating
func (f *VRFTest) deployConsumerCoordinatorPair(c chan<- ConsumerCoordinatorPair) error {
	coord, err := f.Deployer.DeployVRFCoordinator(f.Wallets.Default(), f.link.Address(), f.blockHashStore.Address())
	if err != nil {
		return err
	}
	consumer, err := f.Deployer.DeployVRFConsumer(f.Wallets.Default(), f.link.Address(), coord.Address())
	if err != nil {
		return err
	}
	if err = consumer.Fund(f.Wallets.Default(), big.NewFloat(0), big.NewFloat(2)); err != nil {
		return err
	}
	c <- ConsumerCoordinatorPair{consumer: consumer, coordinator: coord}
	return nil
}

// deployCommonContracts deploys BlockHashStore/VRFCoordinator/VRF contracts
func (f *VRFTest) deployCommonContracts() error {
	var err error
	f.blockHashStore, err = f.Deployer.DeployBlockhashStore(f.Wallets.Default())
	if err != nil {
		return err
	}
	f.vrf, err = f.Deployer.DeployVRFContract(f.Wallets.Default())
	if err != nil {
		return err
	}
	return f.Blockchain.WaitForEvents()
}

// deployContracts deploys common contracts and required amount of VRF consumers
func (f *VRFTest) deployContracts() error {
	if err := f.deployCommonContracts(); err != nil {
		return err
	}

	contractChan := make(chan ConsumerCoordinatorPair, f.TestOptions.NumberOfContracts)
	g := errgroup.Group{}

	for i := 0; i < f.TestOptions.NumberOfContracts; i++ {
		g.Go(func() error {
			return f.deployConsumerCoordinatorPair(contractChan)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(contractChan)
	for contract := range contractChan {
		f.contractInstances = append(f.contractInstances, contract)
	}
	return f.Blockchain.WaitForEvents()
}

// waitRoundFulfilled awaits randomness round fulfillment,
// there is no "round" in VRF by design, it's artificially introduced to have some checkpoint in soak/perf test
func (f *VRFTest) waitRoundFulfilled(roundID int) error {
	for _, p := range f.contractInstances {
		confirmer := contracts.NewVRFConsumerRoundConfirmer(p.consumer, big.NewInt(int64(roundID)), f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(p.consumer.Address(), confirmer)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *VRFTest) watchPerfEvents() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan *contracts.PerfEvent)
		g := errgroup.Group{}
		for _, p := range f.contractInstances {
			p := p
			g.Go(func() error {
				if err := p.consumer.WatchPerfEvents(context.Background(), ch); err != nil {
					return err
				}
				return nil
			})
		}
		for {
			select {
			case event := <-ch:
				rqID := common.Bytes2Hex(event.RequestID[:])
				r := f.testResults.Get(rqID)
				loc, _ := time.LoadLocation("UTC")
				r.EndTime = time.Unix(event.BlockTimestamp.Int64(), 0).In(loc)
				log.Debug().
					Int64("Round", event.Round.Int64()).
					Str("RequestID", rqID).
					Time("EndTime", r.EndTime).
					Msg("Perf event received")
			case <-ctx.Done():
				return
			}
		}
	}()
	return cancel
}

// requestRandomness requests randomness for every consumer for every node (keyHash)
func (f *VRFTest) requestRandomness() error {
	g := errgroup.Group{}
	for p, provingDataByNode := range f.jobMap {
		p := p
		for _, provingData := range provingDataByNode {
			provingData := provingData
			g.Go(func() error {
				err := p.consumer.RequestRandomness(f.Wallets.Default(), provingData.ProvingKeyHash, big.NewInt(1))
				if err != nil {
					return err
				}
				return nil
			})
		}
	}
	return g.Wait()
}

// Run runs VRF performance/soak test
func (f *VRFTest) Run() error {
	if err := f.createChainlinkJobs(); err != nil {
		return err
	}
	var ctx context.Context
	var testCtxCancel context.CancelFunc
	if f.TestOptions.TestDuration.Seconds() > 0 {
		ctx, testCtxCancel = context.WithTimeout(context.Background(), f.TestOptions.TestDuration)
	} else {
		ctx, testCtxCancel = context.WithCancel(context.Background())
	}
	defer testCtxCancel()
	cancelPerfEvents := f.watchPerfEvents()
	currentRound := 0
	for {
		select {
		case <-ctx.Done():
			log.Warn().Msg("Test finished")
			time.Sleep(f.TestOptions.GracefulStopDuration)
			cancelPerfEvents()
			return nil
		default:
			log.Warn().Int("RoundID", currentRound).Msg("New round")
			if err := f.requestRandomness(); err != nil {
				return err
			}
			if err := f.waitRoundFulfilled(currentRound + 1); err != nil {
				return err
			}
			if f.TestOptions.NumberOfRounds != 0 && currentRound >= f.TestOptions.NumberOfRounds {
				log.Warn().Msg("Final round is reached")
				testCtxCancel()
			}
			currentRound++
		}
	}
}

// RecordValues will query all of the latencies of the VRFConsumer and match them by RequestID
func (f *VRFTest) RecordValues(b ginkgo.Benchmarker) error {
	// can't estimate perf metrics in soak mode
	if f.TestOptions.NumberOfRounds == 0 {
		return nil
	}
	actions.SetChainlinkAPIPageSize(f.chainlinkClients, f.TestOptions.NumberOfRounds*f.TestOptions.NumberOfContracts)
	if err := f.setResultStartTimes(); err != nil {
		return err
	}
	return f.calculateLatencies(b)
}

func (f *VRFTest) calculateLatencies(b ginkgo.Benchmarker) error {
	var latencies []time.Duration
	for rqID, testResult := range f.testResults.GetAll() {
		latency := testResult.EndTime.Sub(testResult.StartTime)
		log.Debug().
			Str("RequestID", rqID).
			Time("StartTime", testResult.StartTime).
			Time("EndTime", testResult.EndTime).
			Dur("Duration", latency).
			Msg("Calculating latencies for request id")
		if testResult.StartTime.IsZero() {
			log.Warn().
				Str("RequestID", rqID).
				Msg("Start time zero")
		}
		if testResult.EndTime.IsZero() {
			log.Warn().
				Str("RequestID", rqID).
				Msg("End time zero")
		}
		if latency.Seconds() < 0 {
			log.Warn().
				Str("RequestID", rqID).
				Msg("Latency below zero")
		} else {
			latencies = append(latencies, latency)
		}
	}
	if err := recordResults(b, "Request latency", latencies); err != nil {
		return err
	}
	return nil
}

func (f *VRFTest) setResultStartTimes() error {
	g := errgroup.Group{}
	for contract := range f.jobMap {
		contract := contract
		g.Go(func() error {
			return f.setResultStartTimeByContract(contract)
		})
	}
	return g.Wait()
}

func (f *VRFTest) setResultStartTimeByContract(contract ConsumerCoordinatorPair) error {
	for _, chainlink := range f.chainlinkClients {
		chainlink := chainlink

		jobRuns, err := chainlink.ReadRunsByJob(f.jobMap[contract][chainlink].JobID)
		if err != nil {
			return err
		}
		log.Debug().
			Str("Node", chainlink.URL()).
			Int("Runs", len(jobRuns.Data)).
			Msg("Total runs")
		for _, jobDecodeData := range jobRuns.Data {
			var taskRun client.TaskRun
			for _, tr := range jobDecodeData.Attributes.TaskRuns {
				if tr.Type == "ethabidecodelog" {
					taskRun = tr
				}
			}
			var decodeLogTaskRun *client.DecodeLogTaskRun
			if err := json.Unmarshal([]byte(taskRun.Output), &decodeLogTaskRun); err != nil {
				return err
			}
			rqInts := decodeLogTaskRun.RequestID
			rqID := common.Bytes2Hex(rqInts)
			loc, _ := time.LoadLocation("UTC")
			startTime := jobDecodeData.Attributes.CreatedAt.In(loc)
			log.Debug().
				Time("StartTime", startTime).
				Str("RequestID", rqID).
				Msg("Request found")
			d := f.testResults.Get(rqID)
			d.StartTime = startTime
		}
	}
	return nil
}

// createChainlinkJobs create and collect VRF jobs for every Chainlink node
func (f *VRFTest) createChainlinkJobs() error {
	jobsChan := make(chan VRFJobMap, len(f.chainlinkClients)*len(f.contractInstances))
	g := errgroup.Group{}
	for _, p := range f.contractInstances {
		p := p
		for _, n := range f.chainlinkClients {
			n := n
			g.Go(func() error {
				nodeKeys, err := n.ReadVRFKeys()
				if err != nil {
					return err
				}
				pubKeyCompressed := nodeKeys.Data[0].ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: p.coordinator.Address(),
				}
				ost, err := os.String()
				if err != nil {
					return err
				}
				jobID, err := n.CreateJob(&client.VRFJobSpec{
					Name:               "vrf",
					CoordinatorAddress: p.coordinator.Address(),
					PublicKey:          pubKeyCompressed,
					Confirmations:      1,
					ExternalJobID:      jobUUID.String(),
					ObservationSource:  ost,
				})
				if err != nil {
					return err
				}
				oracleAddr, err := n.PrimaryEthAddress()
				if err != nil {
					return err
				}
				provingKey, err := actions.EncodeOnChainVRFProvingKey(nodeKeys.Data[0])
				if err != nil {
					return err
				}
				if err = p.coordinator.RegisterProvingKey(
					f.Wallets.Default(),
					big.NewInt(1),
					oracleAddr,
					provingKey,
					actions.EncodeOnChainExternalJobID(jobUUID),
				); err != nil {
					return err
				}
				requestHash, err := p.coordinator.HashOfKey(context.Background(), provingKey)
				if err != nil {
					return err
				}
				jobsChan <- VRFJobMap{p: map[client.Chainlink]VRFProvingData{n: {JobID: jobID.Data.ID, ProvingKeyHash: requestHash}}}
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(jobsChan)

	for jobMap := range jobsChan {
		for contractAddr, m := range jobMap {
			if _, ok := f.jobMap[contractAddr]; !ok {
				f.jobMap[contractAddr] = map[client.Chainlink]VRFProvingData{}
			}
			for k, v := range m {
				f.jobMap[contractAddr][k] = v
			}
		}
	}
	return nil
}
