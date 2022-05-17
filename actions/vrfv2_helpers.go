package actions

import (
	"fmt"
	"math/big"

	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
)

func DeployVrfv2Contracts(
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	networks *blockchain.Networks,
) (contracts.VRFCoordinatorV2, contracts.VRFConsumerV2) {
	linkEthFeedResponse := big.NewInt(1e18)
	bhs, err := contractDeployer.DeployBlockhashStore()
	Expect(err).ShouldNot(HaveOccurred())
	mf, err := contractDeployer.DeployMockETHLINKFeed(linkEthFeedResponse)
	Expect(err).ShouldNot(HaveOccurred())
	coordinator, err := contractDeployer.DeployVRFCoordinatorV2(linkTokenContract.Address(), bhs.Address(), mf.Address())
	Expect(err).ShouldNot(HaveOccurred())
	consumer, err := contractDeployer.DeployVRFConsumerV2(linkTokenContract.Address(), coordinator.Address())
	Expect(err).ShouldNot(HaveOccurred())
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	err = linkTokenContract.Transfer(consumer.Address(), big.NewInt(0).Mul(big.NewInt(1e4), big.NewInt(1e18)))
	Expect(err).ShouldNot(HaveOccurred())
	err = coordinator.SetConfig(
		1,
		2.5e6,
		86400,
		33825,
		linkEthFeedResponse,
		ethereum.VRFCoordinatorV2FeeConfig{
			FulfillmentFlatFeeLinkPPMTier1: 1,
			FulfillmentFlatFeeLinkPPMTier2: 1,
			FulfillmentFlatFeeLinkPPMTier3: 1,
			FulfillmentFlatFeeLinkPPMTier4: 1,
			FulfillmentFlatFeeLinkPPMTier5: 1,
			ReqsForTier2:                   big.NewInt(10),
			ReqsForTier3:                   big.NewInt(20),
			ReqsForTier4:                   big.NewInt(30),
			ReqsForTier5:                   big.NewInt(40)},
	)
	Expect(err).ShouldNot(HaveOccurred())
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	err = consumer.CreateFundedSubscription(big.NewInt(0).Mul(big.NewInt(30), big.NewInt(1e18)))
	Expect(err).ShouldNot(HaveOccurred())
	err = networks.Default.WaitForEvents()
	Expect(err).ShouldNot(HaveOccurred())

	return coordinator, consumer
}

func CreateVrfV2Jobs(
	chainlinkNodes []client.Chainlink,
	coordinator contracts.VRFCoordinatorV2,
) ([]*client.Job, [][2]*big.Int) {
	jobs := make([]*client.Job, 0)
	encodedProvingKeys := make([][2]*big.Int, 0)
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
			EVMChainID:               "1337",
			MinIncomingConfirmations: 1,
			PublicKey:                pubKeyCompressed,
			ExternalJobID:            jobUUID.String(),
			ObservationSource:        ost,
			BatchFulfillmentEnabled:  false,
		})
		Expect(err).ShouldNot(HaveOccurred())
		jobs = append(jobs, job)
		provingKey := Vrfv2RegisterProvingKey(vrfKey, oracleAddr, coordinator)
		encodedProvingKeys = append(encodedProvingKeys, provingKey)
	}
	return jobs, encodedProvingKeys
}

func Vrfv2RegisterProvingKey(
	vrfKey *client.VRFKey,
	oracleAddress string,
	coordinator contracts.VRFCoordinatorV2) [2]*big.Int {
	provingKey, err := EncodeOnChainVRFProvingKey(*vrfKey)
	Expect(err).ShouldNot(HaveOccurred())
	err = coordinator.RegisterProvingKey(
		oracleAddress,
		provingKey,
	)
	Expect(err).ShouldNot(HaveOccurred())
	return provingKey
}
