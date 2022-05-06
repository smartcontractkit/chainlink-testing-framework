package testsetups

//revive:disable:dot-imports
import (
	"path/filepath"
	"time"

	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/helmenv/environment"
	"golang.org/x/sync/errgroup"
)

// KeeperBlockTimeTest builds a test to check that chainlink nodes are able to upkeep a specified amount of Upkeep
// contracts within a certain block time
type ChainlinkProfileTest struct {
	Inputs       ChainlinkProfileTestInputs
	TestReporter testreporters.ChainlinkProfileTestReporter

	env            *environment.Environment
	networks       *blockchain.Networks
	defaultNetwork blockchain.EVMClient
}

// KeeperBlockTimeTestInputs are all the required inputs for a Keeper Block Time Test
type ChainlinkProfileTestInputs struct {
	ProfileFunction       func(client.Chainlink)
	ProfileDuration       time.Duration
	ProfileFolderLocation string
	ChainlinkNodes        []client.Chainlink
}

// NewKeeperBlockTimeTest prepares a new keeper block time test to be run
func NewChainlinkProfileTest(inputs ChainlinkProfileTestInputs) *ChainlinkProfileTest {
	return &ChainlinkProfileTest{
		Inputs: inputs,
	}
}

// Setup prepares contracts for the test
func (c *ChainlinkProfileTest) Setup(env *environment.Environment) {
	c.ensureInputValues()
	c.env = env
}

// Run runs the keeper block time test
func (c *ChainlinkProfileTest) Run() {
	profileGroup := new(errgroup.Group)
	for ni, cl := range c.Inputs.ChainlinkNodes {
		chainlinkNode := cl
		nodeIndex := ni
		profileGroup.Go(func() error {
			profileResults, err := chainlinkNode.Profile(c.Inputs.ProfileDuration, c.Inputs.ProfileFunction)
			profileResults.NodeIndex = nodeIndex
			if err != nil {
				return err
			}
			c.TestReporter.Results = append(c.TestReporter.Results, profileResults)
			return nil
		})
	}
	Expect(profileGroup.Wait()).ShouldNot(HaveOccurred(), "Error while gathering chainlink Profile tests")
}

// Networks returns the networks that the test is running on
func (c *ChainlinkProfileTest) TearDownVals() (*environment.Environment, *blockchain.Networks, []client.Chainlink, testreporters.TestReporter) {
	return c.env, c.networks, c.Inputs.ChainlinkNodes, &c.TestReporter
}

// ensureValues ensures that all values needed to run the test are present
func (c *ChainlinkProfileTest) ensureInputValues() {
	Expect(c.Inputs.ProfileFunction).ShouldNot(BeNil(), "Forgot to provide a function to profile")
	Expect(c.Inputs.ProfileDuration.Seconds()).Should(BeNumerically(">=", 1), "Time to profile should be at least 1 second")
	var err error
	c.Inputs.ProfileFolderLocation, err = filepath.Abs(c.Inputs.ProfileFolderLocation)
	Expect(err).ShouldNot(HaveOccurred(), "Error marshalling file path '%s' to absolute path", c.Inputs.ProfileFolderLocation)
	Expect(c.Inputs.ProfileFolderLocation).Should(BeADirectory(), "Provided folder location %s is not a valid directory", c.Inputs.ProfileFolderLocation)
	Expect(c.Inputs.ChainlinkNodes).ShouldNot(BeNil(), "Chainlink nodes you want to profile should be provided")
	Expect(len(c.Inputs.ChainlinkNodes)).Should(BeNumerically(">", 0), "No Chainlink nodes provided to profile")
}
