package slice

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// ValidateAndDeduplicateAddresses validates a slice of addresses and deduplicates them. It returns the deduplicated slice and a boolean indicating if there were duplicates.
func ValidateAndDeduplicateAddresses(addresses []string) ([]string, bool, error) {
	addressSet := make(map[common.Address]struct{})
	deduplicated := make([]string, 0)

	hadDuplicates := false

	for _, addr := range addresses {
		if !common.IsHexAddress(addr) {
			return []string{}, false, fmt.Errorf("address %s is not a valid hex address", addr)
		}

		asAddr := common.HexToAddress(addr)

		if _, exists := addressSet[asAddr]; exists {
			hadDuplicates = true
			continue
		}

		addressSet[asAddr] = struct{}{}
		deduplicated = append(deduplicated, addr)
	}

	return deduplicated, hadDuplicates, nil
}
