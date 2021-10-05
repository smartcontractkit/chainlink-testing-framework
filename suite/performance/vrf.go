package performance

import (
	"context"
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
	RoundTimeout time.Duration
	TestDuration time.Duration
}

// VRFTest is the implementation of Test that will configure and execute soak test
// of VRF contracts & jobs
type VRFTest struct {
	TestOptions     VRFTestOptions
	ContractOptions contracts.OffchainOptions
	Environment     environment.Environment
	Blockchain      client.BlockchainClient
	Wallets         client.BlockchainWallets
	Deployer        contracts.ContractDeployer

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	link              contracts.LinkToken
	vrf               contracts.VRF
	blockHashStore    contracts.BlockHashStore
	contractInstances []ConsumerCoordinatorPair
	adapter           environment.ExternalAdapter

	jobMap VRFJobMap
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
		jobMap:      VRFJobMap{},
	}
}

// RecordValues records VRF metrics
func (f *VRFTest) RecordValues(b ginkgo.Benchmarker) error {
	// TODO: collect metrics
	return nil
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
func (f *VRFTest) waitRoundFulfilled(roundID *big.Int) error {
	for _, p := range f.contractInstances {
		confirmer := contracts.NewVRFConsumerRoundConfirmer(p.consumer, roundID, f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(p.consumer.Address(), confirmer)
	}
	return f.Blockchain.WaitForEvents()
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
	ctx, cancel := context.WithTimeout(context.Background(), f.TestOptions.TestDuration)
	defer cancel()
	i := 1
	for {
		select {
		case <-ctx.Done():
			log.Warn().Msg("Test finished")
			return nil
		default:
			log.Warn().Int("RoundID", i).Msg("New round")
			if err := f.requestRandomness(); err != nil {
				return err
			}
			if err := f.waitRoundFulfilled(big.NewInt(int64(i))); err != nil {
				return err
			}
			i++
		}
	}
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
