// File: internal/tools.go
package internal

import "github.com/ethereum/go-ethereum/common"

// EventKey uniquely identifies an event subscription based on address and topic.
type EventKey struct {
	Address common.Address
	Topic   common.Hash
}
