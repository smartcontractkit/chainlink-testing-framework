// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

//revive:disable:dot-imports
import (
	"context"
	"math/big"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/testreporters"
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs *OCRSoakTestInputs

	TestReporter testreporters.OCRSoakTestReporter
	ocrInstances []contracts.OffchainAggregator
	mockServer   *client.MockserverClient

	env            *environment.Environment
	chainlinkNodes []client.Chainlink
	networks       *client.Networks
	defaultNetwork client.BlockchainClient
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	TestDuration         time.Duration // How long to run the test for (assuming things pass)
	NumberOfContracts    int           // Number of OCR contracts to launch
	ChainlinkNodeFunding *big.Float    // Amount of ETH to fund each chainlink node with
	RoundTimeout         time.Duration // How long to wait for a round to update before timing out
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
			Reports: make(map[string]*testreporters.OCRSoakTestReport),
		},
	}
}

// Setup sets up the test environment, deploying contracts and funding chainlink nodes
func (t *OCRSoakTest) Setup(env *environment.Environment) {
	t.ensureInputValues()
	t.env = env
	var err error

	// Make connections to soak test resources
	networkRegistry := client.NewSoakNetworkRegistry()
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

	t.ocrInstances = actions.DeployOCRContracts(
		t.Inputs.NumberOfContracts,
		linkTokenContract,
		contractDeployer,
		t.chainlinkNodes,
		t.networks,
	)
	err = t.defaultNetwork.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())
	for _, ocrInstance := range t.ocrInstances {
		t.TestReporter.Reports[ocrInstance.Address()] = &testreporters.OCRSoakTestReport{
			ContractAddress: ocrInstance.Address(),
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

	// Test Loop
	roundNumber := 1
	for {
		select {
		case <-testContext.Done():
			log.Info().Msg("Soak Test Complete")
			return
		default:
			log.Info().Int("Round Number", roundNumber).Msg("Starting new Round")
			adapterValue := t.changeAdapterValue(roundNumber)
			t.waitForRoundToComplete(roundNumber)
			t.checkLatestRound(adapterValue, roundNumber)
			roundNumber++
		}
	}
}

// Networks returns the networks that the test is running on
func (t *OCRSoakTest) TearDownVals() (*environment.Environment, *client.Networks, []client.Chainlink, testreporters.TestReporter) {
	return t.env, t.networks, t.chainlinkNodes, &t.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (t *OCRSoakTest) ensureInputValues() {
	inputs := t.Inputs
	Expect(inputs.NumberOfContracts).Should(BeNumerically(">=", 1), "Expecting at least 1 OCR contract")
	Expect(inputs.ChainlinkNodeFunding.Float64()).Should(BeNumerically(">", 0), "Expecting non-zero chainlink node funding amount")
	Expect(inputs.TestDuration).Should(BeNumerically(">=", time.Minute*1), "Expected test duration to be more than a minute")
	Expect(inputs.RoundTimeout).Should(BeNumerically(">=", time.Second*15), "Expected test duration to be more than 15 seconds")
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
		t.defaultNetwork.AddHeaderEventSubscription(ocrInstance.Address(), ocrRound)
	}
	err := t.defaultNetwork.WaitForEvents()
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
