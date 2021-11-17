// Package actions enables common chainlink interactions
package actions

import (
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/utils"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"github.com/smartcontractkit/integrations-framework/contracts"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/integrations-framework/client"
)

const (
	// DefaultArtifactsDir default artifacts dir
	DefaultArtifactsDir = "logs"
)

// FundChainlinkNodes will fund all of the Chainlink nodes with a given amount of ETH in wei
func FundChainlinkNodes(
	nodes []client.Chainlink,
	blockchain client.BlockchainClient,
	amount *big.Float,
) error {
	for _, cl := range nodes {
		toAddress, err := cl.PrimaryEthAddress()
		if err != nil {
			return err
		}
		err = blockchain.Fund(toAddress, amount)
		if err != nil {
			return err
		}
	}
	return blockchain.WaitForEvents()
}

// ChainlinkNodeAddresses will return all the on-chain wallet addresses for a set of Chainlink nodes
func ChainlinkNodeAddresses(nodes []client.Chainlink) ([]common.Address, error) {
	addresses := make([]common.Address, 0)
	for _, node := range nodes {
		primaryAddress, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, common.HexToAddress(primaryAddress))
	}
	return addresses, nil
}

// SetChainlinkAPIPageSize specifies the page size from the Chainlink API, useful for high volume testing
func SetChainlinkAPIPageSize(nodes []client.Chainlink, pageSize int) {
	for _, n := range nodes {
		n.SetPageSize(pageSize)
	}
}

// EncodeOnChainExternalJobID encodes external job uuid to on-chain representation
func EncodeOnChainExternalJobID(jobID uuid.UUID) [32]byte {
	var ji [32]byte
	copy(ji[:], strings.Replace(jobID.String(), "-", "", 4))
	return ji
}

// ExtractRequestIDFromJobRun extracts RequestID from job runs response
func ExtractRequestIDFromJobRun(jobDecodeData client.RunsResponseData) ([]byte, error) {
	var taskRun client.TaskRun
	for _, tr := range jobDecodeData.Attributes.TaskRuns {
		if tr.Type == "ethabidecodelog" {
			taskRun = tr
		}
	}
	var decodeLogTaskRun *client.DecodeLogTaskRun
	if err := json.Unmarshal([]byte(taskRun.Output), &decodeLogTaskRun); err != nil {
		return nil, err
	}
	rqInts := decodeLogTaskRun.RequestID
	return rqInts, nil
}

// EncodeOnChainVRFProvingKey encodes uncompressed public VRF key to on-chain representation
func EncodeOnChainVRFProvingKey(vrfKey client.VRFKey) ([2]*big.Int, error) {
	uncompressed := vrfKey.Data.Attributes.Uncompressed
	provingKey := [2]*big.Int{}
	var set1 bool
	var set2 bool
	// strip 0x to convert to int
	provingKey[0], set1 = new(big.Int).SetString(uncompressed[2:66], 16)
	if !set1 {
		return [2]*big.Int{}, errors.New("can not convert VRF key to *big.Int")
	}
	provingKey[1], set2 = new(big.Int).SetString(uncompressed[66:], 16)
	if !set2 {
		return [2]*big.Int{}, errors.New("can not convert VRF key to *big.Int")
	}
	return provingKey, nil
}

// GetMockserverInitializerDataForOTPE creates mocked weiwatchers data needed for otpe
func GetMockserverInitializerDataForOTPE(
	OCRInstances []contracts.OffchainAggregator,
	chainlinkNodes []client.Chainlink,
) (interface{}, error) {
	var contractsInfo []client.ContractInfoJSON

	for index, OCRInstance := range OCRInstances {
		contractInfo := client.ContractInfoJSON{
			ContractVersion: 4,
			Path:            fmt.Sprintf("contract_%d", index),
			Status:          "live",
			ContractAddress: OCRInstance.Address(),
		}

		contractsInfo = append(contractsInfo, contractInfo)
	}

	contractsInitializer := client.HttpInitializer{
		Request:  client.HttpRequest{Path: "/contracts.json"},
		Response: client.HttpResponse{Body: contractsInfo},
	}

	var nodesInfo []client.NodeInfoJSON

	for _, chainlink := range chainlinkNodes {
		ocrKeys, err := chainlink.ReadOCRKeys()
		if err != nil {
			return nil, err
		}
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
	return initializers, nil
}

// TeardownSuite tears down networks/clients and environment
func TeardownSuite(env *environment.Environment, nets *client.Networks) error {
	if ginkgo.CurrentSpecReport().Failed() {
		testFilename := strings.Split(ginkgo.CurrentSpecReport().FileName(), ".")[0]
		_, testName := filepath.Split(testFilename)
		logsPath := filepath.Join(utils.ProjectRoot, DefaultArtifactsDir, fmt.Sprintf("%s-%d", testName, time.Now().Unix()))
		if err := env.Artifacts.DumpTestResult(logsPath, "chainlink"); err != nil {
			return err
		}
	}
	if nets != nil {
		if err := nets.Teardown(); err != nil {
			return err
		}
	}
	if !env.Config.PersistentConnection {
		if err := env.Disconnect(); err != nil {
			return err
		}
	}
	if !env.Config.Persistent {
		if err := env.Teardown(); err != nil {
			return err
		}
	}
	return nil
}
