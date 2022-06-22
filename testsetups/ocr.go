// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

//revive:disable:dot-imports
import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs *OCRSoakTestInputs

	TestReporter testreporters.OCRSoakTestReporter
	ocrInstances []contracts.OffchainAggregator
	mockServer   *client.MockserverClient

	env            *environment.Environment
	chainlinkNodes []client.Chainlink
	c              blockchain.EVMClient
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	TestDuration         time.Duration // How long to run the test for (assuming things pass)
	NumberOfContracts    int           // Number of OCR contracts to launch
	ChainlinkNodeFunding *big.Float    // Amount of ETH to fund each chainlink node with
	RoundTimeout         time.Duration // How long to wait for a round to update before failing the test
	ExpectedRoundTime    time.Duration // How long each round is expected to take
	TimeBetweenRounds    time.Duration // How long to wait after a completed round to start a new one, set 0 for instant
	StartingAdapterValue int
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(inputs *OCRSoakTestInputs) *OCRSoakTest {
	if inputs.StartingAdapterValue == 0 {
		inputs.StartingAdapterValue = 5
	}
	return &OCRSoakTest{
		Inputs: inputs,
		TestReporter: testreporters.OCRSoakTestReporter{
			Reports:           make(map[string]*testreporters.OCRSoakTestReport),
			ExpectedRoundTime: inputs.ExpectedRoundTime,
		},
	}
}

// Setup sets up the test environment, deploying contracts and funding chainlink nodes
func (t *OCRSoakTest) Setup(env *environment.Environment) {
	t.ensureInputValues()
	t.env = env
	var err error

	// Make connections to soak test resources
	t.c, err = blockchain.NewEthereumMultiNodeClientSetup(DefaultGethSettings)(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to blockchain nodes shouldn't fail")
	contractDeployer, err := contracts.NewContractDeployer(t.c)
	Expect(err).ShouldNot(HaveOccurred(), "Deploying contracts shouldn't fail")
	t.chainlinkNodes, err = client.ConnectChainlinkNodes(env)
	Expect(err).ShouldNot(HaveOccurred(), "Connecting to chainlink nodes shouldn't fail")
	t.mockServer, err = client.ConnectMockServer(env)
	Expect(err).ShouldNot(HaveOccurred(), "Creating mockserver clients shouldn't fail")
	t.c.ParallelTransactions(true)

	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	Expect(err).ShouldNot(HaveOccurred(), "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(t.chainlinkNodes[1:], t.c, t.Inputs.ChainlinkNodeFunding)
	Expect(err).ShouldNot(HaveOccurred(), "Error funding Chainlink nodes")

	t.ocrInstances = actions.DeployOCRContracts(
		t.Inputs.NumberOfContracts,
		linkTokenContract,
		contractDeployer,
		t.chainlinkNodes,
		t.c,
	)
	err = t.c.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range t.ocrInstances {
		t.TestReporter.Reports[ocrInstance.Address()] = &testreporters.OCRSoakTestReport{
			ContractAddress:   ocrInstance.Address(),
			ExpectedRoundtime: t.Inputs.ExpectedRoundTime,
		}
	}
}

// Run starts the OCR soak test
func (t *OCRSoakTest) Run() {
	// Set initial value and create jobs
	By("Setting adapter responses",
		actions.SetAllAdapterResponsesToTheSameValue(t.Inputs.StartingAdapterValue, t.ocrInstances, t.chainlinkNodes, t.mockServer))
	By("Creating OCR jobs", actions.CreateOCRJobs(t.ocrInstances, t.chainlinkNodes, t.mockServer))

	log.Info().
		Str("Test Duration", t.Inputs.TestDuration.Truncate(time.Second).String()).
		Str("Round Timeout", t.Inputs.RoundTimeout.String()).
		Int("Number of OCR Contracts", len(t.ocrInstances)).
		Msg("Starting OCR Soak Test")

	testContext, testCancel := context.WithTimeout(context.Background(), t.Inputs.TestDuration)
	defer testCancel()

	stopTestChannel := make(chan struct{}, 1)
	StartRemoteControlServer("OCR Soak Test", stopTestChannel)

	// Test Loop
	roundNumber := 1
	newRoundTrigger, cancelFunc := context.WithTimeout(context.Background(), 0)
	for {
		select {
		case <-stopTestChannel:
			cancelFunc()
			t.TestReporter.UnexpectedShutdown = true
			log.Warn().Msg("Received shut down signal. Soak test stopping early")
			return
		case <-testContext.Done():
			cancelFunc()
			log.Info().Msg("Soak test complete")
			return
		case <-newRoundTrigger.Done():
			log.Info().Int("Round Number", roundNumber).Msg("Starting new Round")
			adapterValue := t.changeAdapterValue(roundNumber)
			t.waitForRoundToComplete(roundNumber)
			t.checkLatestRound(adapterValue, roundNumber)
			roundNumber++
			log.Info().Str("Time", fmt.Sprint(t.Inputs.TimeBetweenRounds)).Msg("Waiting between OCR Rounds")
			newRoundTrigger, cancelFunc = context.WithTimeout(context.Background(), t.Inputs.TimeBetweenRounds)
		}
	}
}

// Networks returns the networks that the test is running on
func (t *OCRSoakTest) TearDownVals() (*environment.Environment, []client.Chainlink, testreporters.TestReporter, blockchain.EVMClient) {
	return t.env, t.chainlinkNodes, &t.TestReporter, t.c
}

// ensureValues ensures that all values needed to run the test are present
func (t *OCRSoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 OCR contract")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.ExpectedRoundTime).Should(BeNumerically(">=", time.Second*1), "Expected ExpectedRoundTime to be greater than 1 second")
	Expect(inputs.RoundTimeout).Should(BeNumerically(">=", inputs.ExpectedRoundTime), "Expected RoundTimeout to be greater than ExpectedRoundTime")
	Expect(inputs.TimeBetweenRounds).ShouldNot(BeNil(), "You forgot to set TimeBetweenRounds")
	Expect(inputs.TimeBetweenRounds).Should(BeNumerically("<", time.Hour), "TimeBetweenRounds must be less than 1 hour")
}

// changes the mock adapter value for OCR instances to retrieve answers from
func (t *OCRSoakTest) changeAdapterValue(roundNumber int) int {
	adapterValue := 0
	if roundNumber%2 == 1 {
		adapterValue = t.Inputs.StartingAdapterValue
	} else {
		adapterValue = t.Inputs.StartingAdapterValue * 25
	}
	By("Setting adapter responses",
		actions.SetAllAdapterResponsesToTheSameValue(adapterValue, t.ocrInstances, t.chainlinkNodes, t.mockServer))
	log.Debug().
		Int("New Value", adapterValue).
		Int("Round Number", roundNumber).
		Msg("Changed Mock Server Adapter Value for New Round")
	return adapterValue
}

// waits for the specified round number to complete on all deployed OCR instances
func (t *OCRSoakTest) waitForRoundToComplete(roundNumber int) {
	for _, ocrInstance := range t.ocrInstances {
		report := t.TestReporter.Reports[ocrInstance.Address()]
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(
			ocrInstance,
			big.NewInt(int64(roundNumber)),
			t.Inputs.RoundTimeout,
			report,
		)
		t.c.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
	}
	err := t.c.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred(), "Error while waiting for OCR round number %d to complete", roundNumber)
}

// checks on all OCR instances that they all received the correct answer from the latest round
func (t *OCRSoakTest) checkLatestRound(expectedValue, roundNumber int) {
	var roundAnswerGroup sync.WaitGroup
	roundAnswerChannel := make(chan latestRoundAnswer, len(t.ocrInstances))
	for _, ocrInstance := range t.ocrInstances {
		roundAnswerGroup.Add(1)
		ocrInstance := ocrInstance
		go func() {
			defer GinkgoRecover() // This doesn't seem to work properly (ginkgo still panics without recovery). Possible Ginkgo bug?
			defer roundAnswerGroup.Done()

			answer, err := ocrInstance.GetLatestAnswer(context.Background())
			Expect(err).ShouldNot(HaveOccurred(), "Error retrieving latest answer from the OCR contract at %s", ocrInstance.Address())
			log.Info().
				Str("Contract", ocrInstance.Address()).
				Int64("Answer", answer.Int64()).
				Int("Expected Answer", expectedValue).
				Int("Round Number", roundNumber).
				Msg("Latest Round Answer")
			roundAnswerChannel <- latestRoundAnswer{answer: answer.Int64(), contractAddress: ocrInstance.Address()}
		}()
	}
	roundAnswerGroup.Wait()
	close(roundAnswerChannel)
	for latestRound := range roundAnswerChannel {
		Expect(latestRound.answer).Should(BeNumerically(
			"==",
			int64(expectedValue)),
			"Received incorrect answer for OCR round number %d from the OCR contract at %s", latestRound.answer, latestRound.contractAddress,
		)
	}
}

// wrapper around latest answer stats so we can check the answer outside of a go routine
// TODO: I tried doing the assertion inside the go routine, but was met with a possible Ginkgo bug
type latestRoundAnswer struct {
	answer          int64
	contractAddress string
}
