// File: chain_poller/chain_poller_test.go
package chain_poller_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/chain_poller"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/internal"
)

func TestNewChainPoller_Success(t *testing.T) {
	mockClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: mockClient,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)
}

func TestNewChainPoller_NilBlockchainClient(t *testing.T) {
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: nil,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.Error(t, err)
	assert.Nil(t, chainPoller)
	assert.Equal(t, "blockchain client cannot be nil", err.Error())
}

func TestChainPoller_Poll_SingleFilterQueryWithLogs(t *testing.T) {
	mockClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: mockClient,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)

	// Define a filter query
	filterQuery := internal.FilterQuery{
		FromBlock: 101,
		ToBlock:   110,
		Addresses: []common.Address{
			common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			},
		},
	}

	// Define mock logs to return
	testLogs := []internal.Log{
		{
			BlockNumber: 105,
			TxHash:      common.HexToHash("0x1234"),
			Address:     common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Topics:      []common.Hash{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
			Data:        []byte("test data 1"),
			Index:       0,
		},
		{
			BlockNumber: 107,
			TxHash:      common.HexToHash("0x5678"),
			Address:     common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Topics:      []common.Hash{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
			Data:        []byte("test data 2"),
			Index:       1,
		},
	}

	mockClient.On("FilterLogs", mock.Anything, filterQuery).Return(testLogs, nil)

	// Perform polling
	logs, err := chainPoller.Poll(context.Background(), []internal.FilterQuery{filterQuery})
	require.NoError(t, err)
	require.Len(t, logs, 2)

	assert.Equal(t, testLogs, logs)

	// Verify that FilterLogs was called with expected query
	mockClient.AssertCalled(t, "FilterLogs", mock.Anything, filterQuery)
}

func TestChainPoller_Poll_MultipleFilterQueries(t *testing.T) {
	mockClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: mockClient,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)

	// Define multiple filter queries
	filterQuery1 := internal.FilterQuery{
		FromBlock: 101,
		ToBlock:   110,
		Addresses: []common.Address{
			common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			},
		},
	}

	filterQuery2 := internal.FilterQuery{
		FromBlock: 101,
		ToBlock:   110,
		Addresses: []common.Address{
			common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash("0xfeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedface"),
			},
		},
	}

	// Define mock logs for filterQuery1
	testLogs1 := []internal.Log{
		{
			BlockNumber: 103,
			TxHash:      common.HexToHash("0x1111"),
			Address:     common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
			Topics:      []common.Hash{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")},
			Data:        []byte("test data 1"),
			Index:       0,
		},
	}

	// Define mock logs for filterQuery2
	testLogs2 := []internal.Log{
		{
			BlockNumber: 104,
			TxHash:      common.HexToHash("0x2222"),
			Address:     common.HexToAddress("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			Topics:      []common.Hash{common.HexToHash("0xfeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedfacefeedface")},
			Data:        []byte("test data 2"),
			Index:       1,
		},
	}

	mockClient.On("FilterLogs", mock.Anything, filterQuery1).Return(testLogs1, nil)
	mockClient.On("FilterLogs", mock.Anything, filterQuery2).Return(testLogs2, nil)

	// Perform polling
	logs, err := chainPoller.Poll(context.Background(), []internal.FilterQuery{filterQuery1, filterQuery2})
	require.NoError(t, err)
	require.Len(t, logs, 2)

	expectedLogs := append(testLogs1, testLogs2...)
	assert.Equal(t, expectedLogs, logs)

	// Verify that FilterLogs was called with both queries
	mockClient.AssertCalled(t, "FilterLogs", mock.Anything, filterQuery1)
	mockClient.AssertCalled(t, "FilterLogs", mock.Anything, filterQuery2)
}

func TestChainPoller_Poll_NoLogs(t *testing.T) {
	mockClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: mockClient,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)

	// Define a filter query with no matching logs
	filterQuery := internal.FilterQuery{
		FromBlock: 101,
		ToBlock:   110,
		Addresses: []common.Address{
			common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
			},
		},
	}

	// Mock FilterLogs to return no logs
	mockClient.On("FilterLogs", mock.Anything, filterQuery).Return([]internal.Log{}, nil)

	// Perform polling
	logs, err := chainPoller.Poll(context.Background(), []internal.FilterQuery{filterQuery})
	require.NoError(t, err)
	require.Len(t, logs, 0)

	// Verify that FilterLogs was called with expected query
	mockClient.AssertCalled(t, "FilterLogs", mock.Anything, filterQuery)
}

func TestChainPoller_Poll_FilterLogsError(t *testing.T) {
	mockClient := new(internal.MockBlockchainClient)
	testLogger := logging.GetTestLogger(t)

	config := chain_poller.ChainPollerConfig{
		BlockchainClient: mockClient,
		Logger:           testLogger,
		ChainID:          1,
	}

	chainPoller, err := chain_poller.NewChainPoller(config)
	require.NoError(t, err)
	require.NotNil(t, chainPoller)

	// Define a filter query
	filterQuery := internal.FilterQuery{
		FromBlock: 101,
		ToBlock:   110,
		Addresses: []common.Address{
			common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		},
		Topics: [][]common.Hash{
			{
				common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			},
		},
	}

	// Mock FilterLogs to return an error
	mockClient.On("FilterLogs", mock.Anything, filterQuery).Return([]internal.Log{}, errors.New("FilterLogs error"))

	// Perform polling
	logs, err := chainPoller.Poll(context.Background(), []internal.FilterQuery{filterQuery})
	require.NoError(t, err)
	require.Len(t, logs, 0)

	// Verify that FilterLogs was called with expected query
	mockClient.AssertCalled(t, "FilterLogs", mock.Anything, filterQuery)

}
