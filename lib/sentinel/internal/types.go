// File: internal/tools.go
package internal

import "github.com/ethereum/go-ethereum/common"

// FilterQuery represents the parameters to filter logs/events.
type FilterQuery struct {
	FromBlock uint64
	ToBlock   uint64
	Topics    [][]common.Hash
	Addresses []common.Address
}

// Log represents a single log event fetched from the blockchain.
type Log struct {
	Address     common.Address
	Topics      []common.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      common.Hash
	Index       uint
}

// EventKey uniquely identifies an event subscription based on address and topic.
type EventKey struct {
	Address common.Address
	Topic   common.Hash
}
