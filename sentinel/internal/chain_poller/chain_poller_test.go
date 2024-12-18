// File: internal/chain_poller/chain_poller_test.go
package chain_poller

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

// Helper function to initialize a new ChainPoller for testing
func setupChainPoller(t *testing.T, blockchainClient *internal.MockBlockchainClient, logger *zerolog.Logger, chainID int64) *ChainPoller {
	t.Helper()
	config := ChainPollerConfig{
		BlockchainClient: blockchainClient,
		Logger:           logger,
		ChainID:          chainID,
	}

	chainPoller, err := NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)
	return chainPoller
}

// Helper function to create a FilterQuery
func createFilterQuery(fromBlock, toBlock uint64, addresses []common.Address, topics [][]common.Hash) api.FilterQuery {
	return api.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
		Topics:    topics,
	}
}

// Helper function to create mock logs
func createMockLogs(blockNumbers []uint64, txHashes []common.Hash, addresses []common.Address, topics [][]common.Hash, data [][]byte, indexes []uint) []api.Log {
	logs := make([]api.Log, len(blockNumbers))
	for i := range logs {
		logs[i] = api.Log{
			BlockNumber: blockNumbers[i],
			TxHash:      txHashes[i],
			Address:     addresses[i],
			Topics:      topics[i],
			Data:        data[i],
			Index:       indexes[i],
		}
	}
	return logs
}

func TestNewChainPoller(t *testing.T) {
	testLogger := logging.GetTestLogger(t)

	tests := []struct {
		name              string
		config            ChainPollerConfig
		expectedError     string
		expectChainPoller bool
	}{
		{
			name: "Success",
			config: ChainPollerConfig{
				BlockchainClient: new(internal.MockBlockchainClient),
				Logger:           &testLogger,
				ChainID:          1,
			},
			expectedError:     "",
			expectChainPoller: true,
		},
		{
			name: "NilBlockchainClient",
			config: ChainPollerConfig{
				BlockchainClient: nil,
				Logger:           &testLogger,
				ChainID:          1,
			},
			expectedError:     "blockchain client cannot be nil",
			expectChainPoller: false,
		},
		{
			name: "NoLoggerPassed",
			config: ChainPollerConfig{
				BlockchainClient: new(internal.MockBlockchainClient),
				Logger:           nil,
				ChainID:          1,
			},
			expectedError:     "no logger passed",
			expectChainPoller: false,
		},
		{
			name: "ChainIDNotSet",
			config: ChainPollerConfig{
				BlockchainClient: new(internal.MockBlockchainClient),
				Logger:           &testLogger,
				ChainID:          0, // Assuming 0 is invalid
			},
			expectedError:     "chain ID not set",
			expectChainPoller: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chainPoller, err := NewChainPoller(tt.config)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Nil(t, chainPoller)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, chainPoller)
			}
		})
	}
}

func TestChainPoller_FilterLogs(t *testing.T) {
	mockBlockchainClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)
	chainID := int64(1)

	chainPoller := setupChainPoller(t, mockBlockchainClient, &testLogger, chainID)

	tests := []struct {
		name            string
		filterQueries   []api.FilterQuery
		mockReturnLogs  [][]api.Log
		mockReturnError []error
		expectedLogs    []api.Log
		expectedError   bool
	}{
		{
			name: "SingleFilterQueryWithLogs",
			filterQueries: []api.FilterQuery{
				createFilterQuery(
					101,
					110,
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{
							common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
						},
					},
				),
			},
			mockReturnLogs: [][]api.Log{
				createMockLogs(
					[]uint64{105, 107},
					[]common.Hash{common.HexToHash("0x1234"), common.HexToHash("0x5678")},
					[]common.Address{
						common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
						common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
					},
					[][]common.Hash{
						{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
						{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
					},
					[][]byte{
						[]byte("test data 1"),
						[]byte("test data 2"),
					},
					[]uint{0, 1},
				),
			},
			mockReturnError: []error{nil},
			expectedLogs: createMockLogs(
				[]uint64{105, 107},
				[]common.Hash{common.HexToHash("0x1234"), common.HexToHash("0x5678")},
				[]common.Address{
					common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
					common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
				},
				[][]common.Hash{
					{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
					{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
				},
				[][]byte{
					[]byte("test data 1"),
					[]byte("test data 2"),
				},
				[]uint{0, 1},
			),
			expectedError: false,
		},
		{
			name: "MultipleFilterQueries",
			filterQueries: []api.FilterQuery{
				createFilterQuery(
					101,
					110,
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{
							common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
						},
					},
				),
				createFilterQuery(
					101,
					110,
					[]common.Address{common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")},
					[][]common.Hash{
						{
							common.HexToHash("0xfeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedface"),
						},
					},
				),
			},
			mockReturnLogs: [][]api.Log{
				createMockLogs(
					[]uint64{103},
					[]common.Hash{common.HexToHash("0x1111")},
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
					},
					[][]byte{[]byte("test data 1")},
					[]uint{0},
				),
				createMockLogs(
					[]uint64{104},
					[]common.Hash{common.HexToHash("0x2222")},
					[]common.Address{common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")},
					[][]common.Hash{
						{common.HexToHash("0xfeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedface")},
					},
					[][]byte{[]byte("test data 2")},
					[]uint{1},
				),
			},
			mockReturnError: []error{nil, nil},
			expectedLogs: append(
				createMockLogs(
					[]uint64{103},
					[]common.Hash{common.HexToHash("0x1111")},
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
					},
					[][]byte{[]byte("test data 1")},
					[]uint{0},
				),
				createMockLogs(
					[]uint64{104},
					[]common.Hash{common.HexToHash("0x2222")},
					[]common.Address{common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")},
					[][]common.Hash{
						{common.HexToHash("0xfeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedface")},
					},
					[][]byte{[]byte("test data 2")},
					[]uint{1},
				)...,
			),
			expectedError: false,
		},
		{
			name: "NoLogs",
			filterQueries: []api.FilterQuery{
				createFilterQuery(
					101,
					110,
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
						},
					},
				),
			},
			mockReturnLogs: [][]api.Log{
				{}, // No logs returned
			},
			mockReturnError: []error{nil},
			expectedLogs:    []api.Log{}, // Expecting an empty slice
			expectedError:   false,
		},
		{
			name: "FilterLogsError",
			filterQueries: []api.FilterQuery{
				createFilterQuery(
					101,
					110,
					[]common.Address{common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")},
					[][]common.Hash{
						{
							common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
						},
					},
				),
			},
			mockReturnLogs: [][]api.Log{
				{}, // No logs returned due to error
			},
			mockReturnError: []error{errors.New("FilterLogs error")},
			expectedLogs:    []api.Log{}, // Expecting an empty slice
			expectedError:   false,       // According to original test, no error is propagated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			for i, fq := range tt.filterQueries {
				var logs []api.Log
				if i < len(tt.mockReturnLogs) {
					logs = tt.mockReturnLogs[i]
				} else {
					logs = []api.Log{}
				}
				var err error
				if i < len(tt.mockReturnError) {
					err = tt.mockReturnError[i]
				} else {
					err = nil
				}
				mockBlockchainClient.On("FilterLogs", mock.Anything, fq).Return(logs, err).Once()
			}

			// Execute the method under test
			logs, err := chainPoller.FilterLogs(context.Background(), tt.filterQueries)

			// Assertions
			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if len(tt.expectedLogs) == 0 {
				assert.Empty(t, logs, "Expected logs to be empty")
			} else {
				assert.Equal(t, tt.expectedLogs, logs, "Logs should match expected logs")
			}

			// Assert that all expectations were met
			mockBlockchainClient.AssertExpectations(t)
		})
	}
}
