package steps

import (
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
)

// GetMockserverInitializerDataForOTPE crates mocked weiwatchers data needed for otpe
func GetMockserverInitializerDataForOTPE(offChainAggregatorInstanceAddress string, chainlinkNodes []client.Chainlink) interface{} {
	contractInfo := client.ContractInfoJSON{
		ContractVersion: 4,
		Path:            "test",
		Status:          "live",
		ContractAddress: offChainAggregatorInstanceAddress,
	}

	contractsInfo := []client.ContractInfoJSON{contractInfo}

	contractsInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/contracts.json"},
		Response: client.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []client.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.ReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred())
		nodeInfo := client.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/nodes.json"},
		Response: client.HttpResponse{Body: nodesInfo},
	}
	initializers := []client.HttpInitializer{contractsInitializer, nodesInitializer}
	return initializers
}
