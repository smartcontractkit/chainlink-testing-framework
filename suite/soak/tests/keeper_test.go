package soak

//revive:disable:dot-imports
import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/geth"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
)

var _ = Describe("Keeper block time soak test @soak-keeper-block-time", func() {
	var (
		err                 error
		env                 *environment.Environment
		keeperBlockTimeTest *testsetups.KeeperBlockTimeTest
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			env = environment.New(
				&environment.Config{InsideK8s: true, TTL: 12 * time.Hour},
			)
			err = env.
				AddHelm(mockservercfg.New(nil)).
				AddHelm(mockserver.New(nil)).
				AddHelm(geth.New(nil)).
				AddHelm(chainlink.New(nil)).
				Run()
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Str("Namespace", env.Cfg.Namespace).Msg("Connected to Soak Environment")
		})

		By("Setup the Keeper test", func() {
			keeperBlockTimeTest = testsetups.NewKeeperBlockTimeTest(
				testsetups.KeeperBlockTimeTestInputs{
					NumberOfContracts: 50,
					KeeperRegistrySettings: &contracts.KeeperRegistrySettings{
						PaymentPremiumPPB:    uint32(200000000),
						FlatFeeMicroLINK:     uint32(0),
						BlockCountPerTurn:    big.NewInt(3),
						CheckGasLimit:        uint32(2500000),
						StalenessSeconds:     big.NewInt(90000),
						GasCeilingMultiplier: uint16(1),
						MinUpkeepSpend:       big.NewInt(0),
						MaxPerformGas:        uint32(5000000),
						FallbackGasPrice:     big.NewInt(2e11),
						FallbackLinkPrice:    big.NewInt(2e18),
					},
					CheckGasToBurn:       2400000,
					PerformGasToBurn:     2400000,
					BlockRange:           2000,
					BlockInterval:        200,
					ChainlinkNodeFunding: big.NewFloat(10),
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
