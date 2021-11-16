package smoke

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/helmenv/environment"
	"math/big"
	"path/filepath"

	"github.com/smartcontractkit/integrations-framework/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/actions"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
)

var _ = Describe("VRF suite @vrf", func() {
	var (
		err                error
		nets               *client.Networks
		cd                 contracts.ContractDeployer
		consumer           contracts.VRFConsumer
		coordinator        contracts.VRFCoordinator
		encodedProvingKeys = make([][2]*big.Int, 0)
		lt                 contracts.LinkToken
		cls                []client.Chainlink
		e                  *environment.Environment
	)
	BeforeEach(func() {
		By("Deploying the environment", func() {
			e, err = environment.NewEnvironmentFromPreset(filepath.Join(utils.PresetRoot, "chainlink-cluster-3"))
			Expect(err).ShouldNot(HaveOccurred())
			err = e.Connect()
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Getting the clients", func() {
			nets, err = client.NewNetworks(e, nil)
			Expect(err).ShouldNot(HaveOccurred())
			cd, err = contracts.NewContractDeployer(nets.Default)
			Expect(err).ShouldNot(HaveOccurred())
			cls, err = client.NewChainlinkClients(e)
			Expect(err).ShouldNot(HaveOccurred())
			nets.Default.ParallelTransactions(true)
		})
		By("Funding Chainlink nodes", func() {
			txCost, err := nets.Default.CalculateTXSCost(1)
			Expect(err).ShouldNot(HaveOccurred())
			err = actions.FundChainlinkNodes(cls, nets.Default, txCost)
			Expect(err).ShouldNot(HaveOccurred())
		})
		By("Deploying VRF contracts", func() {
			lt, err = cd.DeployLinkTokenContract()
			Expect(err).ShouldNot(HaveOccurred())
			bhs, err := cd.DeployBlockhashStore()
			Expect(err).ShouldNot(HaveOccurred())
			coordinator, err = cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
			Expect(err).ShouldNot(HaveOccurred())
			consumer, err = cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
			Expect(err).ShouldNot(HaveOccurred())
			err = lt.Transfer(consumer.Address(), big.NewInt(2e18))
			Expect(err).ShouldNot(HaveOccurred())
			_, err = cd.DeployVRFContract()
			Expect(err).ShouldNot(HaveOccurred())
			err = nets.Default.WaitForEvents()
			Expect(err).ShouldNot(HaveOccurred())
		})

		By("Creating jobs and registering proving keys", func() {
			for _, n := range cls {
				nodeKey, err := n.CreateVRFKey()
				Expect(err).ShouldNot(HaveOccurred())
				log.Debug().Interface("Key JSON", nodeKey).Msg("Created proving key")
				pubKeyCompressed := nodeKey.Data.ID
				jobUUID := uuid.NewV4()
				os := &client.VRFTxPipelineSpec{
					Address: coordinator.Address(),
				}
				ost, err := os.String()
				Expect(err).ShouldNot(HaveOccurred())
				_, err = n.CreateJob(&client.VRFJobSpec{
					Name:               fmt.Sprintf("vrf-%s", jobUUID),
					CoordinatorAddress: coordinator.Address(),
					PublicKey:          pubKeyCompressed,
					Confirmations:      1,
					ExternalJobID:      jobUUID.String(),
					ObservationSource:  ost,
				})
				Expect(err).ShouldNot(HaveOccurred())

				oracleAddr, err := n.PrimaryEthAddress()
				Expect(err).ShouldNot(HaveOccurred())
				provingKey, err := actions.EncodeOnChainVRFProvingKey(*nodeKey)
				Expect(err).ShouldNot(HaveOccurred())
				err = coordinator.RegisterProvingKey(
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
			err = consumer.RequestRandomness(requestHash, big.NewInt(1))
			Expect(err).ShouldNot(HaveOccurred())

			Eventually(func(g Gomega) {
				out, err := consumer.RandomnessOutput(context.Background())
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(out.Uint64()).Should(Not(BeNumerically("==", 0)))
				log.Debug().Uint64("Output", out.Uint64()).Msg("Randomness fulfilled")
			}, "2m", "1s").Should(Succeed())
		})
	})

	AfterEach(func() {
		By("Printing gas stats", func() {
			nets.Default.GasStats().PrintStats()
		})
		By("Tearing down the environment", func() {
			err = actions.TeardownSuite(e, nets)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
