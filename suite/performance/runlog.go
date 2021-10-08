package performance

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"golang.org/x/sync/errgroup"
)

// RunlogJobMap is a custom map type that holds the record of jobs by the contract instance and the chainlink node
type RunlogJobMap map[ConsumerOraclePair]map[client.Chainlink]string

// ConsumerOraclePair consumer and oracle pair
type ConsumerOraclePair struct {
	consumer contracts.APIConsumer
	oracle   contracts.Oracle
	jobUUID  string
}

// RunlogTestOptions contains the parameters for the Runlog soak test to be executed
type RunlogTestOptions struct {
	TestOptions
	RoundTimeout time.Duration
	AdapterValue int
	TestDuration time.Duration
}

// RunlogTest is the implementation of Test that will configure and execute soak test
// of Runlog contracts & jobs
type RunlogTest struct {
	TestOptions RunlogTestOptions
	Environment environment.Environment
	Blockchain  client.BlockchainClient
	Wallets     client.BlockchainWallets
	Deployer    contracts.ContractDeployer
	Link        contracts.LinkToken

	chainlinkClients  []client.Chainlink
	nodeAddresses     []common.Address
	contractInstances []*ConsumerOraclePair
	adapter           environment.ExternalAdapter

	jobMap RunlogJobMap
}

// NewRunlogTest creates new Runlog performance/soak test
func NewRunlogTest(
	testOptions RunlogTestOptions,
	env environment.Environment,
	link contracts.LinkToken,
	blockchain client.BlockchainClient,
	wallets client.BlockchainWallets,
	deployer contracts.ContractDeployer,
	adapter environment.ExternalAdapter,
) Test {
	return &RunlogTest{
		TestOptions: testOptions,
		Environment: env,
		Link:        link,
		Blockchain:  blockchain,
		Wallets:     wallets,
		Deployer:    deployer,
		adapter:     adapter,
		jobMap:      RunlogJobMap{},
	}
}

// RecordValues records Runlog metrics
func (f *RunlogTest) RecordValues(b ginkgo.Benchmarker) error {
	// TODO: collect metrics
	return nil
}

// Setup setups Runlog performance/soak test
func (f *RunlogTest) Setup() error {
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

func (f *RunlogTest) deployContract(c chan<- *ConsumerOraclePair) error {
	oracle, err := f.Deployer.DeployOracle(f.Wallets.Default(), f.Link.Address())
	if err != nil {
		return err
	}
	if err = oracle.SetFulfillmentPermission(f.Wallets.Default(), f.nodeAddresses[0].Hex(), true); err != nil {
		return err
	}
	consumer, err := f.Deployer.DeployAPIConsumer(f.Wallets.Default(), f.Link.Address())
	if err != nil {
		return err
	}
	err = consumer.Fund(f.Wallets.Default(), nil, big.NewFloat(20000))
	if err != nil {
		return err
	}
	c <- &ConsumerOraclePair{consumer: consumer, oracle: oracle}
	return nil
}

func (f *RunlogTest) deployContracts() error {
	contractChan := make(chan *ConsumerOraclePair, f.TestOptions.NumberOfContracts)
	g := errgroup.Group{}

	for i := 0; i < f.TestOptions.NumberOfContracts; i++ {
		g.Go(func() error {
			return f.deployContract(contractChan)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(contractChan)
	for contract := range contractChan {
		f.contractInstances = append(f.contractInstances, contract)
	}
	log.Warn().Int("Pairs", len(f.contractInstances)).Msg("Pairs")
	return f.Blockchain.WaitForEvents()
}

func (f *RunlogTest) requestData() error {
	g := errgroup.Group{}
	for _, p := range f.contractInstances {
		p := p
		g.Go(func() error {
			jobUUIDReplaces := strings.Replace(p.jobUUID, "-", "", 4)
			var jobID [32]byte
			copy(jobID[:], jobUUIDReplaces)
			if err := p.consumer.CreateRequestTo(
				f.Wallets.Default(),
				p.oracle.Address(),
				jobID,
				big.NewInt(1e18),
				fmt.Sprintf("%s/five", f.adapter.ClusterURL()),
				"data,result",
				big.NewInt(100),
			); err != nil {
				return err
			}
			return nil
		})
	}
	return g.Wait()
}

// Run runs Runlog performance/soak test
func (f *RunlogTest) Run() error {
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
			if err := f.requestData(); err != nil {
				return err
			}
			if err := f.waitRoundEnd(i); err != nil {
				return err
			}
			i++
		}
	}
}

func (f *RunlogTest) waitRoundEnd(roundID int) error {
	for _, p := range f.contractInstances {
		rc := contracts.NewRunlogRoundConfirmer(p.consumer, big.NewInt(int64(roundID)), f.TestOptions.RoundTimeout)
		f.Blockchain.AddHeaderEventSubscription(p.consumer.Address(), rc)
	}
	return f.Blockchain.WaitForEvents()
}

func (f *RunlogTest) createChainlinkJobs() error {
	jobsChan := make(chan RunlogJobMap, len(f.contractInstances))
	g := errgroup.Group{}

	bta := client.BridgeTypeAttributes{
		Name: "five",
		URL:  fmt.Sprintf("%s/five", f.adapter.ClusterURL()),
	}
	if err := f.chainlinkClients[0].CreateBridge(&bta); err != nil {
		return err
	}
	os := &client.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	ost, err := os.String()
	if err != nil {
		return err
	}

	for _, p := range f.contractInstances {
		p := p
		g.Go(func() error {
			jobUUID := uuid.NewV4()
			p.jobUUID = jobUUID.String()
			_, err := f.chainlinkClients[0].CreateJob(&client.DirectRequestJobSpec{
				Name:              "direct_request",
				ContractAddress:   p.oracle.Address(),
				ExternalJobID:     jobUUID.String(),
				ObservationSource: ost,
			})
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	close(jobsChan)

	for jobMap := range jobsChan {
		for contractAddr, m := range jobMap {
			if _, ok := f.jobMap[contractAddr]; !ok {
				f.jobMap[contractAddr] = map[client.Chainlink]string{}
			}
			for k, v := range m {
				f.jobMap[contractAddr][k] = v
			}
		}
	}
	return nil
}
