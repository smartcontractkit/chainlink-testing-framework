package evm_storage

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	StorageSlotSizeBytes = 32
)

var (
	ErrNoSlot = errors.New("no such slot found in layout JSON")
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

// MustSlot calculates a slot in Solidity mapping for storage field and a key
func (s *StorageLayout) MustSlot(label string) string {
	if _, ok := s.GetSlots()[label]; !ok {
		panic(fmt.Errorf("layout label: %s, %w", label, ErrNoSlot))
	}
	return s.GetSlots()[label]
}

// MustMapSlot calculates a slot in Solidity mapping for storage field and a key
func (s *StorageLayout) MustMapSlot(label, key string) string {
	if _, ok := s.GetSlots()[label]; !ok {
		panic(fmt.Errorf("layout label: %s, %w", label, ErrNoSlot))
	}
	baseSlot := s.GetSlots()[label]
	return mapSlot(baseSlot, key)
}

// MustArraySlot calculates a slot in Solidity array for storage field and a key
func (s *StorageLayout) MustArraySlot(label string, index int64) string {
	if _, ok := s.GetSlots()[label]; !ok {
		panic(fmt.Errorf("layout label: %s, %w", label, ErrNoSlot))
	}
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

func ShiftHexByOffset(hexStr string, offset int) string {
	// Strip "0x"
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Parse into big.Int
	n := new(big.Int)
	n.SetString(hexStr, 16)

	// Shift left by offset * 8 bits
	n.Lsh(n, uint(offset*8))

	// Return as 0x-prefixed, 32-byte hex
	return fmt.Sprintf("0x%064x", n)
}

// MustEncodeStorageSlot encodes a value for Solidity storage slots based on type
// Panics if encoding fails
func MustEncodeStorageSlot(solidityType string, value interface{}) string {
	// Handle address type specially
	if solidityType == "address" {
		switch v := value.(type) {
		case common.Address:
			return fmt.Sprintf("0x%064x", v.Big())
		case string:
			if !common.IsHexAddress(v) {
				panic(fmt.Sprintf("invalid address format: %v", v))
			}
			return fmt.Sprintf("0x%064x", common.HexToAddress(v).Big())
		default:
			panic(fmt.Sprintf("unsupported address type: %T", value))
		}
	}

	// Create the ABI type
	typ, err := abi.NewType(solidityType, "", nil)
	if err != nil {
		panic(fmt.Sprintf("invalid solidity type %q: %v", solidityType, err))
	}

	// Special handling for bytes32 strings
	if solidityType == "bytes32" {
		if s, ok := value.(string); ok {
			if len(s) > 32 {
				panic("string too long for bytes32")
			}
			var b [32]byte
			copy(b[:], s)
			value = b
		}
	}

	// Encode the value
	encoded, err := abi.Arguments{{Type: typ}}.Pack(value)
	if err != nil {
		panic(fmt.Sprintf("encoding failed for %v (%T) as %s: %v", value, value, solidityType, err))
	}

	// For uint256 and int256, we need to take the last 32 bytes
	if solidityType == "uint256" || solidityType == "int256" {
		if len(encoded) > 32 {
			encoded = encoded[len(encoded)-32:]
		}
	}

	return fmt.Sprintf("0x%064x", new(big.Int).SetBytes(encoded))
}

// MergeHex merges two hex strings with bitwise "OR"
// should be used when you see values with offsets in smart contract storage layout.json file
// example:
//
//	|----------------+-------------------------------------------+------+--------+-------+-------------------------|
//	| number_uint8   | uint8                                     | 3    | 0      | 1     | src/Counter.sol:Counter |
//	|----------------+-------------------------------------------+------+--------+-------+-------------------------|
//	| boolean        | bool                                      | 3    | 1      | 1     | src/Counter.sol:Counter |
//	|----------------+-------------------------------------------+------+--------+-------+-------------------------|
func MergeHex(a, b string) string {
	ai := new(big.Int)
	bi := new(big.Int)
	ai.SetString(a[2:], 16)
	bi.SetString(b[2:], 16)
	ai.Or(ai, bi)
	return fmt.Sprintf("0x%064x", ai)
}
