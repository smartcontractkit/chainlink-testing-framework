package testsetups

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/helmenv/environment"
)

// VRFV2SoakTest defines a typical VRFV2 soak test
type VRFV2SoakTest struct {
	Inputs *VRFV2SoakTestInputs

	TestReporter testreporters.VRFV2SoakTestReporter
	mockServer   *client.MockserverClient

	env            *environment.Environment
	ChainlinkNodes []client.Chainlink
	Networks       *blockchain.Networks
	DefaultNetwork blockchain.EVMClient

	NumberOfRequests int

	ErrorOccurred error
	ErrorCount    int
}

// VRFV2SoakTestTestFunc function type for the request and validation you want done on each iteration
type VRFV2SoakTestTestFunc func(t *VRFV2SoakTest, requestNumber int) error

// VRFV2SoakTestInputs define required inputs to run a vrfv2 soak test
type VRFV2SoakTestInputs struct {
	TestDuration         time.Duration // How long to run the test for (assuming things pass)
	ChainlinkNodeFunding *big.Float    // Amount of ETH to fund each chainlink node with
	StopTestOnError      bool          // Do we want the test to stop after any error or just continue on

	RequestsPerMinute int                   // Number of requests for randomness per minute
	TestFunc          VRFV2SoakTestTestFunc // The function that makes the request and validations wanted
}

// VRFV2SoakTestJobInfo defines a jobs into and proving key info
type VRFV2SoakTestJobInfo struct {
	Job            *client.Job
	ProvingKey     [2]*big.Int
	ProvingKeyHash [32]byte
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

// Setup sets up the test environment
func (t *VRFV2SoakTest) Setup(env *environment.Environment, charts []string, isLocal bool) {
	t.ensureInputValues()
	t.env = env
	var err error
	var networkRegistry *blockchain.NetworkRegistry

	// Make connections to soak test resources
	if isLocal {
		err = env.ConnectAll()
		Expect(err).ShouldNot(HaveOccurred())
		networkRegistry = blockchain.NewDefaultNetworkRegistry()
		t.ChainlinkNodes, err = client.ConnectChainlinkNodesByCharts(env, charts)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		t.mockServer, err = client.ConnectMockServer(env)
		Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	} else {
		networkRegistry = blockchain.NewSoakNetworkRegistry()
		t.ChainlinkNodes, err = client.ConnectChainlinkNodesSoakByCharts(env, charts)
		Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
		t.mockServer, err = client.ConnectMockServerSoak(env)
		Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	}

	t.Networks, err = networkRegistry.GetNetworks(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
	t.DefaultNetwork = t.Networks.Default
	t.DefaultNetwork.ParallelTransactions(true)
	Expect(err).ShouldNot(HaveOccurred())
}

// Run starts the VRFV2 soak test
func (t *VRFV2SoakTest) Run() {
	log.Info().
		Str("Test Duration", t.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Max number of requests per minute wanted", t.Inputs.RequestsPerMinute).
		Msg("Starting VRFV2 Soak Test")

	// set the requests to only run for a certain amount of time
	testContext, testCancel := context.WithTimeout(context.Background(), t.Inputs.TestDuration)
	defer testCancel()

	t.NumberOfRequests = 0

	// variables dealing with how often to tick and how to stop the ticker
	stop := false
	startTime := time.Now()
	ticker := time.NewTicker(time.Minute / time.Duration(t.Inputs.RequestsPerMinute))

	for {
		// start the loop by checking to see if any of the TestFunc responses have returned an error
		if t.Inputs.StopTestOnError {
			Expect(t.ErrorOccurred).ShouldNot(HaveOccurred())
		}

		select {
		case <-testContext.Done():
			// stop making requests
			stop = true
			ticker.Stop()
			break // breaks the select block
		case <-ticker.C:
			// make the next request
			t.NumberOfRequests++
			go requestAndValidate(t, t.NumberOfRequests)
		}

		if stop {
			break // breaks the for loop and stops the test
		}
	}
	log.Info().Int("Requests", t.NumberOfRequests).Msg("Total Completed Requests")
	log.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished VRFV2 Soak Test Requests")
	Expect(t.ErrorCount).To(BeNumerically("==", 0), "We had a number of errors")
}

func requestAndValidate(t *VRFV2SoakTest, requestNumber int) {
	defer GinkgoRecover()
	// Errors in goroutines cause some weird behavior with how ginkgo returns the error
	// We are having the TestFunc return any errors it sees so we can then propogate them in
	//  the main thread and get proper ginkgo behavior on test failures
	log.Info().Int("Request Number", requestNumber).Msg("Making a Request")
	err := t.Inputs.TestFunc(t, requestNumber)
	// only set the error to be checked if err is not nil so we avoid race conditions with passing requests
	if err != nil {
		t.ErrorOccurred = err
		log.Error().Err(err).Msg("Error Occurred during test")
		t.ErrorCount++
	}
}

// Networks returns the networks that the test is running on
func (t *VRFV2SoakTest) TearDownVals() (*environment.Environment, *blockchain.Networks, []client.Chainlink, testreporters.TestReporter) {
	return t.env, t.Networks, t.ChainlinkNodes, &t.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (t *VRFV2SoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.RequestsPerMinute).Should(BeNumerically(">=", 1), "Expecting at least 1 request per minute")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.TestFunc).ShouldNot(BeNil(), "Expected to have a test to run")
}
