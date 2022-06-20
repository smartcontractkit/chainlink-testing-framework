package benchmark

//revive:disable:dot-imports
import (
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/helmenv/environment"
	// "github.com/smartcontractkit/helmenv/tools"
)

var _ = Describe("Keeper benchmark suite @benchmark-keeper", func() {
	var (
		err                 error
		env                 *environment.Environment
		keeperBenchmarkTest *testsetups.KeeperBenchmarkTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
				"/root/test-env.json", // Default location for the soak-test-runner container
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Keeper test", func() {
			keeperBenchmarkTest = testsetups.NewKeeperBenchmarkTest(
				testsetups.KeeperBenchmarkTestInputs{
					NumberOfContracts: 500,
					KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
						PaymentPremiumPPB:    uint32(200000000),
						BlockCountPerTurn:    big.NewInt(100),
						CheckGasLimit:        uint32(2000000),
						StalenessSeconds:     big.NewInt(90000),
						GasCeilingMultiplier: uint16(2),
						MaxPerformGas:        uint32(5000000),
						MinUpkeepSpend:       big.NewInt(0),
						FallbackGasPrice:     big.NewInt(2e11),
						FallbackLinkPrice:    big.NewInt(2e18),
					},
					CheckGasToBurn:       1000000,
					PerformGasToBurn:     150000,
					BlockRange:           3600,
					BlockInterval:        20,
					ChainlinkNodeFunding: big.NewFloat(1000),
				},
			)
			keeperBenchmarkTest.Setup(env)
		})
	})

	Describe("Watching the keeper contracts to ensure they reply in time", func() {
		It("Watches for Upkeep counts", func() {
			keeperBenchmarkTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if err := actions.TeardownRemoteSuite(keeperBenchmarkTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
		})
	})
})
