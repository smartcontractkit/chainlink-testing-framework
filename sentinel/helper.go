// File: helper.go
package sentinel

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
)

// ConvertAPILogToTypesLog maps an api.Log to a types.Log.
// Returns a pointer to types.Log and an error if mapping fails.
func ConvertAPILogToTypesLog(log api.Log) (*types.Log, error) {
	// Validate required fields
	if log.Address == (common.Address{}) {
		return nil, fmt.Errorf("api.Log Address is empty")
	}
	if len(log.Topics) == 0 {
		return nil, fmt.Errorf("api.Log Topics are empty")
	}
	if log.BlockNumber == 0 {
		return nil, fmt.Errorf("api.Log BlockNumber is zero")
	}

	// Map fields
	mappedLog := &types.Log{
		Address:     log.Address,
		Topics:      log.Topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash,
		Index:       log.Index,
	}

	return mappedLog, nil
}
