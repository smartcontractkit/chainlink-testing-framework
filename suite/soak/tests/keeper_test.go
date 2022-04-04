package soak

//revive:disable:dot-imports
import (
	"math/big"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/helmenv/tools"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/testsetups"
)

var _ = Describe("Keeper performance suite @block-time-keeper", func() {
	var (
		err                 error
		env                 *environment.Environment
		keeperBlockTimeTest *testsetups.KeeperBlockTimeTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env, err = environment.DeployOrLoadEnvironmentFromConfigFile(
				tools.ChartsRoot,
				"/root/test-env.json", // Default location for the soak-test-runner container
			)
			Expect(err).ShouldNot(HaveOccurred(), "Environment deployment shouldn't fail")
			log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Keeper test", func() {
			keeperBlockTimeTest = testsetups.NewKeeperBlockTimeTest(
				testsetups.KeeperBlockTimeTestInputs{
					NumberOfContracts: 501,
					KeeperContractSettings: &testsetups.KeeperContractSettings{
						PaymentPremiumPPB:    uint32(200000000),
						BlockCountPerTurn:    big.NewInt(3),
						CheckGasLimit:        uint32(2500000),
						StalenessSeconds:     big.NewInt(90000),
						GasCeilingMultiplier: uint16(1),
						FallbackGasPrice:     big.NewInt(2e11),
						FallbackLinkPrice:    big.NewInt(2e18),
					},
					BlockRange:           100,
					BlockInterval:        20,
					ChainlinkNodeFunding: big.NewFloat(.001),
				},
			)
			keeperBlockTimeTest.Setup(env)
		})
	})

	Describe("Watching the keeper contracts to ensure they reply in time", func() {
		It("Watches for Upkeep counts", func() {
			keeperBlockTimeTest.Run()
		})
	})

	AfterEach(func() {
		By("Tearing down the environment", func() {
			if err := actions.TeardownRemoteSuite(keeperBlockTimeTest.TearDownVals()); err != nil {
				log.Error().Err(err).Msg("Error tearing down environment")
			}
		})
	})
})
