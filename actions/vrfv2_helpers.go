package actions

//revive:disable:dot-imports
import (
	"fmt"
	"math/big"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
)

func DeployVRFV2Contracts(
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
	linkEthFeedAddress string,
) (contracts.VRFCoordinatorV2, contracts.VRFConsumerV2, contracts.BlockHashStore) {
	bhs, err := contractDeployer.DeployBlockhashStore()
	Expect(err).ShouldNot(HaveOccurred())
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), linkEthFeedAddress)
	Expect(err).ShouldNot(HaveOccurred())
	consumer, err := contractDeployer.DeployVRFConsumerV2(linkTokenContract.Address(), coordinator.Address())
	Expect(err).ShouldNot(HaveOccurred())
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	return coordinator, consumer, bhs
}

type VRFV2EncodedProvingKey [2]*big.Int

func CreateVRFV2Jobs(
	chainlinkNodes []client.Chainlink,
	coordinator contracts.VRFCoordinatorV2,
	networks *blockchain.Networks,
) ([]*client.Job, []VRFV2EncodedProvingKey) {
	jobs := make([]*client.Job, 0)
	encodedProvingKeys := make([]VRFV2EncodedProvingKey, 0)
	for _, n := range chainlinkNodes {
		vrfKey, err := n.CreateVRFKey()
		Expect(err).ShouldNot(HaveOccurred())
		log.Debug().Interface("Key JSON", vrfKey).Msg("Created proving key")
		pubKeyCompressed := vrfKey.Data.ID
		jobUUID := uuid.NewV4()
		os := &client.VRFV2TxPipelineSpec{
			Address: coordinator.Address(),
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		oracleAddr, err := n.PrimaryEthAddress()
		Expect(err).ShouldNot(HaveOccurred())
		job, err := n.CreateJob(&client.VRFV2JobSpec{
			Name:                     fmt.Sprintf("vrf-%s", jobUUID),
			CoordinatorAddress:       coordinator.Address(),
			FromAddress:              oracleAddr,
			EVMChainID:               networks.Default.GetChainID().String(),
			MinIncomingConfirmations: 1,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		jobs = append(jobs, job)
		provingKey := VRFV2RegisterProvingKey(vrfKey, oracleAddr, coordinator)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)
	}
	return jobs, encodedProvingKeys
}

func VRFV2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2,
) VRFV2EncodedProvingKey {
	provingKey, err := EncodeOnChainVRFProvingKey(*vrfKey)
	Expect(err).ShouldNot(HaveOccurred())
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	Expect(err).ShouldNot(HaveOccurred())
	return provingKey
}
