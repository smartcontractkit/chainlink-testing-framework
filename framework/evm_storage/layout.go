package evm_storage

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	StorageSlotSizeBytes = 32
)

type StorageEntry struct {
	Label  string `json:"label"`
	Slot   string `json:"slot"`
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

type StorageLayout struct {
	Storage []StorageEntry `json:"storage"`
}

// New creates a new storage layout wrapper
func New(filename string) (*StorageLayout, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage layout file %s: %w", filename, err)
	}
	var layout StorageLayout
	if err := json.Unmarshal(raw, &layout); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage layout: %w", err)
	}
	return &layout, nil
}

func (s *StorageLayout) GetSlots() map[string]string {
	slots := make(map[string]string)
	for _, entry := range s.Storage {
		slots[entry.Label] = entry.Slot
	}
	return slots
}

// Slot calculates a slot in Solidity mapping for storage field and a key
func (s *StorageLayout) Slot(label string) string {
	return s.GetSlots()[label]
}

// MapSlot calculates a slot in Solidity mapping for storage field and a key
func (s *StorageLayout) MapSlot(label, key string) string {
	baseSlot := s.GetSlots()[label]
	return mapSlot(baseSlot, key)
}

// ArraySlot calculates a slot in Solidity array for storage field and a key
func (s *StorageLayout) ArraySlot(label string, index int64) string {
	baseSlot := s.GetSlots()[label]
	return arraySlot(baseSlot, index)
}

// arraySlot calculates a slot in an array for a base slot and index
func arraySlot(baseSlot string, index int64) string {
	slotInt := new(big.Int)
	slotInt.SetString(baseSlot, 10)
	base := crypto.Keccak256(slotInt.FillBytes(make([]byte, StorageSlotSizeBytes)))
	baseInt := new(big.Int).SetBytes(base)
	offset := new(big.Int).Add(baseInt, big.NewInt(index))
	return "0x" + fmt.Sprintf("%064x", offset)
}

// mapSlot calculates a slot in Solidity mapping for a base slot and a key
func mapSlot(baseSlot string, key string) string {
	slotInt := new(big.Int)
	slotInt.SetString(baseSlot, 10)
	keyBytes, _ := hex.DecodeString(strings.TrimPrefix(key, "0x"))
	buf := append(make([]byte, StorageSlotSizeBytes-len(keyBytes)), keyBytes...)
	buf = append(buf, slotInt.FillBytes(make([]byte, StorageSlotSizeBytes))...)
	hash := crypto.Keccak256(buf)
	return "0x" + hex.EncodeToString(hash)
}
