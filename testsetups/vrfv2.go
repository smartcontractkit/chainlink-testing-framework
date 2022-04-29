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

// VRFV2SoakTest defines a typical OCR soak test
type VRFV2SoakTest struct {
	Inputs *VRFV2SoakTestInputs

	TestReporter testreporters.VRFV2SoakTestReporter
	mockServer   *client.MockserverClient

	env            *environment.Environment
	consumer       contracts.VRFConsumerV2
	coordinator    contracts.VRFCoordinatorV2
	chainlinkNodes []client.Chainlink
	jobInfo        []VRFV2SoakTestJobInfo
	networks       *blockchain.Networks
	defaultNetwork blockchain.EVMClient
}

type VRFV2SoakTestJobInfo struct {
	job            *client.Job
	provingKey     [2]*big.Int
	provingKeyHash [32]byte
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type VRFV2SoakTestInputs struct {
	TestDuration            time.Duration // How long to run the test for (assuming things pass)
	RequestsPerSecondWanted int           // Number of requests for randomness per minute
	ChainlinkNodeFunding    *big.Float    // Amount of ETH to fund each chainlink node with
	RoundTimeout            time.Duration // How long to wait for a round to update before timing out
}

type requestedRandomnessData struct {
	jobInfo       VRFV2SoakTestJobInfo
	requestNumber int
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
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
	t.chainlinkNodes, err = client.ConnectChainlinkNodesSoak(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	t.mockServer, err = client.ConnectMockServerSoak(env)
	Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	t.defaultNetwork.ParallelTransactions(true)
	Expect(err).ShouldNot(HaveOccurred())

	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes
	err = actions.FundChainlinkNodes(t.chainlinkNodes, t.defaultNetwork, t.Inputs.ChainlinkNodeFunding)
	Expect(err).ShouldNot(HaveOccurred())

	t.coordinator, t.consumer = actions.DeployVrfv2Contracts(linkTokenContract, contractDeployer, t.networks)
	jobs, provingKeys := actions.CreateVrfV2Jobs(t.chainlinkNodes, t.coordinator)
	Expect(len(jobs)).Should(Equal(len(provingKeys)), "Should have a set of keys for each job")

	// Create proving key hash here so we aren't calculating it in the test run itself.
	for i, pk := range provingKeys {
		keyHash, err := t.coordinator.HashOfKey(context.Background(), pk)
		Expect(err).ShouldNot(HaveOccurred(), "Should be able to create a keyHash from the proving keys")
		ji := VRFV2SoakTestJobInfo{
			job:            jobs[i],
			provingKey:     provingKeys[i],
			provingKeyHash: keyHash,
		}
		t.jobInfo = append(t.jobInfo, ji)
	}

	err = t.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())
}

// Run starts the OCR soak test
func (t *VRFV2SoakTest) Run() {
	// durationBetweenRequests := (time.Second / time.Duration(t.Inputs.RequestsPerSecondWanted)).Milliseconds()
	requestedRandomnessChannel := make(chan requestedRandomnessData)
	timeout := time.Minute * 2
	go readRandomness(requestedRandomnessChannel, timeout, t.chainlinkNodes)

	log.Info().
		Str("Test Duration", t.Inputs.TestDuration.Truncate(time.Second).String()).
		Str("Round Timeout", t.Inputs.RoundTimeout.String()).
		Int("Max number of requests per second wanted", t.Inputs.RequestsPerSecondWanted).
		Msg("Starting OCR Soak Test")

	testContext, testCancel := context.WithTimeout(context.Background(), t.Inputs.TestDuration)
	defer testCancel()

	// Test Loop
	requestNumber := 1
	stop := false
	startTime := time.Now()
	// nextRequestTime := startTime.UnixMilli()
	ticker := time.NewTicker(time.Second / time.Duration(t.Inputs.RequestsPerSecondWanted))
	for {
		select {
		case <-testContext.Done():
			stop = true
			ticker.Stop()
			break
		case <-ticker.C:
			log.Info().Int("Request Number", requestNumber).Msg("Making a request")
			go requestRandomness(t.coordinator, t.consumer, t.jobInfo[0], requestNumber, requestedRandomnessChannel)
			requestNumber++
		}

		if stop {
			break
		}
	}
	log.Info().Int("Requests", requestNumber).Msg("Total Completed Requests")
	log.Info().Str("Run Time", time.Since(startTime).String()).Msg("Finished VRFV2 Soak Test")
}

func requestRandomness(
	coordinator contracts.VRFCoordinatorV2,
	consumer contracts.VRFConsumerV2,
	jobInfo VRFV2SoakTestJobInfo,
	requestNumber int,
	ch chan requestedRandomnessData,
) {
	words := uint32(10)
	err := consumer.RequestRandomness(jobInfo.provingKeyHash, 1, 1, 300000, words)
	Expect(err).ShouldNot(HaveOccurred())
	data := requestedRandomnessData{
		jobInfo:       jobInfo,
		requestNumber: requestNumber,
	}
	ch <- data
}

func readRandomness(ch chan requestedRandomnessData, timeout time.Duration, cls []client.Chainlink) {
	for data := range ch {
		Eventually(func(g Gomega) {
			jobRuns, err := cls[0].ReadRunsByJob(data.jobInfo.job.Data.ID)
			g.Expect(err).ShouldNot(HaveOccurred())
			g.Expect(len(jobRuns.Data)).Should(BeNumerically(">=", data.requestNumber))
			// if data.requestNumber%10 == 0 {
			// 	randomness, err := consumer.GetAllRandomWords(context.Background(), int(words))
			// 	g.Expect(err).ShouldNot(HaveOccurred())
			// 	for _, w := range randomness {
			// 		log.Debug().Uint64("Output", w.Uint64()).Msg("Randomness fulfilled")
			// 		g.Expect(w.Uint64()).Should(Not(BeNumerically("==", 0)), "Expected the VRF job give an answer other than 0")
			// 	}
			// }
		}, timeout, "1s").Should(Succeed())
	}
}

// Networks returns the networks that the test is running on
func (t *VRFV2SoakTest) TearDownVals() (*environment.Environment, *blockchain.Networks, []client.Chainlink, testreporters.TestReporter) {
	return t.env, t.networks, t.chainlinkNodes, &t.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (t *VRFV2SoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.RequestsPerSecondWanted).Should(BeNumerically(">=", 1), "Expecting at least 1 request per second")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.RoundTimeout).Should(BeNumerically(">=", time.Second*15), "Expected test duration to be more than 15 seconds")
}
