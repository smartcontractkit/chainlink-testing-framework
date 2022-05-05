// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/blockchain"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testreporters"
)

// VRFV2SoakTest defines a typical VRFV2 soak test
type VRFV2SoakTest struct {
	Inputs *VRFV2SoakTestInputs

	TestReporter testreporters.VRFV2SoakTestReporter
	mockServer   *client.MockserverClient

	env            *environment.Environment
	Consumer       contracts.VRFConsumerV2
	Coordinator    contracts.VRFCoordinatorV2
	ChainlinkNodes []client.Chainlink
	JobInfo        []VRFV2SoakTestJobInfo
	networks       *blockchain.Networks
	defaultNetwork blockchain.EVMClient

	NumberRequestsToValidate int
	NumberRequestsValidated  int
}

// VRFV2SoakTestJobInfo defines a jobs into and proving key info
type VRFV2SoakTestJobInfo struct {
	Job            *client.Job
	ProvingKey     [2]*big.Int
	ProvingKeyHash [32]byte
}

// VRFV2SoakTestTestFunc function type for the request and validation you want done on each iteration
type VRFV2SoakTestTestFunc func(t *VRFV2SoakTest, requestNumber int)

// VRFV2SoakTestInputs define required inputs to run a vrfv2 soak test
type VRFV2SoakTestInputs struct {
	TestDuration         time.Duration // How long to run the test for (assuming things pass)
	ChainlinkNodeFunding *big.Float    // Amount of ETH to fund each chainlink node with

	RequestsPerSecond  int                   // Number of requests for randomness per minute
	ReadEveryNRequests int                   // Check the randomness output every n number of requests
	TestFunc           VRFV2SoakTestTestFunc // The function that makes the request and validations wanted
}

// NewVRFV2SoakTest creates a new vrfv2 soak test to setup and run
func NewVRFV2SoakTest(inputs *VRFV2SoakTestInputs) *VRFV2SoakTest {
	return &VRFV2SoakTest{
		Inputs: inputs,
		TestReporter: testreporters.VRFV2SoakTestReporter{
			Reports: make(map[string]*testreporters.VRFV2SoakTestReport),
		},
	}
}

// Setup sets up the test environment, deploying contracts and funding chainlink nodes
func (t *VRFV2SoakTest) Setup(env *environment.Environment) {
	t.ensureInputValues()
	t.env = env
	var err error

	// Make connections to soak test resources
	networkRegistry := blockchain.NewSoakNetworkRegistry()
	t.networks, err = networkRegistry.GetNetworks(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
	t.defaultNetwork = t.networks.Default
	contractDeployer, err := contracts.NewContractDeployer(t.defaultNetwork)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
	t.ChainlinkNodes, err = client.ConnectChainlinkNodesSoak(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	t.mockServer, err = client.ConnectMockServerSoak(env)
	Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	t.defaultNetwork.ParallelTransactions(true)
	Expect(err).ShouldNot(HaveOccurred())

	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes
	err = actions.FundChainlinkNodes(t.ChainlinkNodes, t.defaultNetwork, t.Inputs.ChainlinkNodeFunding)
	Expect(err).ShouldNot(HaveOccurred())

	t.Coordinator, t.Consumer = actions.DeployVrfv2Contracts(linkTokenContract, contractDeployer, t.networks)
	jobs, provingKeys := actions.CreateVrfV2Jobs(t.ChainlinkNodes, t.Coordinator)
	Expect(len(jobs)).Should(Equal(len(provingKeys)), "Should have a set of keys for each job")

	// Create proving key hash here so we aren't calculating it in the test run itself.
	for i, pk := range provingKeys {
		keyHash, err := t.Coordinator.HashOfKey(context.Background(), pk)
		Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a keyHash from the proving keys")
		ji := VRFV2SoakTestJobInfo{
			Job:            jobs[i],
			ProvingKey:     provingKeys[i],
			ProvingKeyHash: keyHash,
		}
		t.JobInfo = append(t.JobInfo, ji)
	}

	err = t.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())
}

// Run starts the VRFV2 soak test
func (t *VRFV2SoakTest) Run() {
	log.Info().
		Str("Test Duration", t.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Max number of requests per second wanted", t.Inputs.RequestsPerSecond).
		Msg("Starting VRFV2 Soak Test")

	testContext, testCancel := context.WithTimeout(context.Background(), t.Inputs.TestDuration)
	defer testCancel()

	t.NumberRequestsToValidate = 0
	t.NumberRequestsValidated = 0

	// Test Loop
	requestNumber := 1
	stop := false
	startTime := time.Now()
	ticker := time.NewTicker(time.Second / time.Duration(t.Inputs.RequestsPerSecond))
	for {
		select {
		case <-testContext.Done():
			// stop making requests
			stop = true
			ticker.Stop()
			break
		case <-ticker.C:
			go requestAndValidate(t, requestNumber)
			requestNumber++
		}

		if stop {
			if t.NumberRequestsToValidate == t.NumberRequestsValidated {
				// stop the test loop entirely
				break
			} else {
				sleepTime := time.Duration(5)
				time.Sleep(time.Second * time.Duration(sleepTime))
				log.Info().Int64("Sleeping for ", int64(sleepTime)).Msg("Waiting for test \"Eventually\" statements to complete")
			}
		}
	}
	log.Info().Int("Requests", requestNumber).Msg("Total Completed Requests")
	log.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished VRFV2 Soak Test")
}

func requestAndValidate(
	t *VRFV2SoakTest,
	requestNumber int,
) {
	// defer GinkgoRecover()
	log.Info().Int("Request Number", requestNumber).Msg("Making a Request")
	t.Inputs.TestFunc(t, requestNumber)
	t.NumberRequestsToValidate++
}

// Networks returns the networks that the test is running on
func (t *VRFV2SoakTest) TearDownVals() (*environment.Environment, *blockchain.Networks, []client.Chainlink, testreporters.TestReporter) {
	return t.env, t.networks, t.ChainlinkNodes, &t.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (t *VRFV2SoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.RequestsPerSecond).Should(BeNumerically(">=", 1), "Expecting at least 1 request per second")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.ReadEveryNRequests).Should(BeNumerically(">", 0), "Expected the test to read requests for verification at some point")
	Expect(inputs.TestFunc).ShouldNot(BeNil(), "Expected to have a test to run")
}
