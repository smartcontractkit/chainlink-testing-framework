package suite

import (
	"context"
	"errors"
	"github.com/avast/retry-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/contracts"
	"github.com/smartcontractkit/integrations-framework/suite"
	"github.com/smartcontractkit/integrations-framework/tools"
	"math/big"
	"strings"
)

var _ = Describe("Direct request suite", func() {
	DescribeTable("Runs direct request job, checks data on-chain", func(
		initFunc client.BlockchainNetworkInit,
		fluxOptions contracts.FluxAggregatorOptions,
	) {
		s, err := suite.DefaultLocalSetup(initFunc)
		Expect(err).ShouldNot(HaveOccurred())
		oracle, err := s.Deployer.DeployOracle(s.Wallets.Default(), s.Link.Address())
		Expect(err).ShouldNot(HaveOccurred())
		consumer, err := s.Deployer.DeployAPIConsumer(s.Wallets.Default(), s.Link.Address())
		Expect(err).ShouldNot(HaveOccurred())
		err = consumer.Fund(s.Wallets.Default(), nil, big.NewInt(2e18))
		Expect(err).ShouldNot(HaveOccurred())
		clNodes, _, err := suite.ConnectToTemplateNodes()
		Expect(err).ShouldNot(HaveOccurred())
		err = suite.FundTemplateNodes(s.Client, s.Wallets, clNodes, 2e18, 0)
		Expect(err).ShouldNot(HaveOccurred())
		adapter := tools.NewExternalAdapter()
		keysData, err := clNodes[0].ReadETHKeys()
		Expect(err).ShouldNot(HaveOccurred())
		// permit the node to call fulfill contract method
		err = oracle.SetFulfillmentPermission(s.Wallets.Default(), keysData.Data[0].Attributes.Address, true)
		Expect(err).ShouldNot(HaveOccurred())

		jobUUID := uuid.NewV4()
		os := &client.DirectRequestTxPipelineSpec{
			URL:         adapter.InsideDockerAddr + "/five",
			Method:      "POST",
			RequestData: "{}",
			DataPath:    "data,result",
		}
		ost, err := os.String()
		Expect(err).ShouldNot(HaveOccurred())
		_, err = clNodes[0].CreateJob(&client.DirectRequestJobSpec{
			Name:              "direct_request",
			ContractAddress:   oracle.Address(),
			ExternalJobID:     jobUUID.String(),
			ObservationSource: ost,
		})
		// job uuid must be 32 byte, without hyphens
		jobUUIDReplaces := strings.Replace(jobUUID.String(), "-", "", 4)
		Expect(err).ShouldNot(HaveOccurred())
		var jobID [32]byte
		copy(jobID[:], jobUUIDReplaces)
		err = consumer.CreateRequestTo(s.Wallets.Default(), oracle.Address(), jobID, big.NewInt(1e18), adapter.InsideDockerAddr+"/five", "data,result", big.NewInt(100))
		Expect(err).ShouldNot(HaveOccurred())
		err = retry.Do(func() error {
			d, err := consumer.Data(context.Background())
			if d == nil {
				return errors.New("no data")
			}
			log.Debug().Int64("Data", d.Int64()).Send()
			if d.Int64() != 5 {
				return errors.New("data is not on chain still")
			}
			if err != nil {
				return err
			}
			return nil
		})
		Expect(err).ShouldNot(HaveOccurred())
	},
		Entry("on Ethereum Hardhat", client.NewHardhatNetwork, contracts.DefaultFluxAggregatorOptions()),
	)
})
