package contracts

import (
	"context"
	"math/big"

	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

var _ = Describe("VRF suite @vrf", func() {

	var (
		suiteSetup         *actions.DefaultSuiteSetup
		nodes              []client.Chainlink
		consumer           contracts.VRFConsumer
		coordinator        contracts.VRFCoordinator
		encodedProvingKeys = make([][2]*big.Int, 0)
		err                error
	)

	BeforeEach(func() {
		By("Deploying the environment", func() {
			suiteSetup, err = actions.DefaultLocalSetup(
				environment.NewChainlinkCluster(1),
				client.NewNetworkFromConfig,
				tools.ProjectRoot,
			)
			Expect(err).ShouldNot(HaveOccurred())
			nodes, err = environment.GetChainlinkClients(suiteSetup.Env)
			Expect(err).ShouldNot(HaveOccurred())

			suiteSetup.Client.ParallelTransactions(true)
		})
		By("Funding Chainlink nodes", func() {
			ethAmount, err := suiteSetup.Deployer.CalculateETHForTXs(suiteSetup.Wallets.Default(), suiteSetup.Network.Config(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(nodes, suiteSetup.Client, suiteSetup.Wallets.Default(), ethAmount, nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying VRF contracts", func() {
			bhs, err := suiteSetup.Deployer.DeployBlockhashStore(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = suiteSetup.Deployer.DeployVRFCoordinator(suiteSetup.Wallets.Default(), suiteSetup.Link.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = suiteSetup.Deployer.DeployVRFConsumer(suiteSetup.Wallets.Default(), suiteSetup.Link.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.Fund(suiteSetup.Wallets.Default(), big.NewFloat(0), big.NewFloat(2))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = suiteSetup.Deployer.DeployVRFContract(suiteSetup.Wallets.Default())
			Expect(err).ShouldNot(HaveOccurred())
			err = suiteSetup.Client.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Creating jobs and registering proving keys", func() {
			for _, n := range nodes {
				nodeKeys, err := n.ReadVRFKeys()
				Expect(err).ShouldNot(HaveOccurred())
				log.Debug().Interface("Key JSON", nodeKeys).Msg("Created proving key")
				pubKeyCompressed := nodeKeys.Data[0].ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred())
				_, err = n.CreateJob(&client.VRFJobSpec{
					Name:               "vrf",
					CoordinatorAddress: coordinator.Address(),
					PublicKey:          pubKeyCompressed,
					Confirmations:      1,
					ExternalJobID:      jobUUID.String(),
					ObservationSource:  ost,
				})
				Expect(err).ShouldNot(HaveOccurred())

				oracleAddr, err := n.PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				provingKey, err := actions.EncodeOnChainVRFProvingKey(nodeKeys.Data[0])
				Expect(err).ShouldNot(HaveOccurred())
				err = coordinator.RegisterProvingKey(
					suiteSetup.Wallets.Default(),
					big.NewInt(1),
					oracleAddr,
					provingKey,
					actions.EncodeOnChainExternalJobID(jobUUID),
				)
				Expect(err).ShouldNot(HaveOccurred())
				encodedProvingKeys = append(encodedProvingKeys, provingKey)
			}
		})
	})

	Describe("with VRF job", func() {
		It("fulfills randomness", func() {
			requestHash, err := coordinator.HashOfKey(context.Background(), encodedProvingKeys[0])
			Expect(err).ShouldNot(HaveOccurred())
			err = consumer.RequestRandomness(suiteSetup.Wallets.Default(), requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred())
			err = retry.Do(func() error {
				out, err := consumer.RandomnessOutput(context.Background())
				if err != nil {
					return err
				}
				if out.Uint64() == 0 {
					return errors.New("randomness has not fulfilled yet")
				}
				log.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
				return nil
			})
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	AfterEach(func() {
		By("Printing gas stats", func() {
			suiteSetup.Client.GasStats().PrintStats()
		})
		By("Tearing down the environment", suiteSetup.TearDown())
	})
})
