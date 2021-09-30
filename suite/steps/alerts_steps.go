package steps

import (
	"encoding/json"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/environment/charts/mockserver"
	"github.com/smartcontractkit/integrations-framework/tools"
	"os"
	"path/filepath"
)

// WriteDataForOTPEToInitializerFileForMockserver Write to initializerJson mocked weiwatchers data needed for otpe
func WriteDataForOTPEToInitializerFileForMockserver(offChainAggregatorInstanceAddress string, chainlinkNodes []client.Chainlink) {
	contractInfo := mockserver.ContractInfoJSON{
		ContractVersion: 4,
		Path:            "test",
		Status:          "live",
		ContractAddress: offChainAggregatorInstanceAddress,
	}

	contractsInfo := []mockserver.ContractInfoJSON{contractInfo}

	contractsInitializer := mockserver.HttpInitializer{
		Request:  mockserver.HttpRequest{Path: "/contracts.json"},
		Response: mockserver.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []mockserver.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.ReadOCRKeys()
		Expect(err).ShouldNot(HaveOccurred())
		nodeInfo := mockserver.NodeInfoJSON{
			NodeAddress: []string{ocrKeys.Data[0].Attributes.OnChainSigningAddress},
			ID:          ocrKeys.Data[0].ID,
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}

	nodesInitializer := mockserver.HttpInitializer{
		Request:  mockserver.HttpRequest{Path: "/nodes.json"},
		Response: mockserver.HttpResponse{Body: nodesInfo},
	}
	initializers := []mockserver.HttpInitializer{contractsInitializer, nodesInitializer}

	initializersBytes, err := json.Marshal(initializers)
	Expect(err).ShouldNot(HaveOccurred())

	fileName := filepath.Join(tools.ProjectRoot, "environment/charts/mockserver-config/static/initializerJson.json")
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	Expect(err).ShouldNot(HaveOccurred())

	body := string(initializersBytes)
	_, err = f.WriteString(body)
	Expect(err).ShouldNot(HaveOccurred())

	err = f.Close()
	Expect(err).ShouldNot(HaveOccurred())
}
