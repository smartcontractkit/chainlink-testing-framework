package seth

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ABIFinder struct {
	ContractMap   ContractMap
	ContractStore *ContractStore
}

type ABIFinderResult struct {
	ABI            abi.ABI
	Method         *abi.Method
	DuplicateCount int
	contractName   string
}

func (a *ABIFinderResult) ContractName() string {
	return strings.TrimSuffix(a.contractName, ".abi")
}

func NewABIFinder(contractMap ContractMap, contractStore *ContractStore) ABIFinder {
	return ABIFinder{
		ContractMap:   contractMap,
		ContractStore: contractStore,
	}
}

// FindABIByMethod finds the ABI method and instance for the given contract address and signature
// If the contract address is known, it will use the ABI instance that is known to be at the address.
// If the contract address is not known, it will iterate over all known ABIs and check if any of them
// has a method with the given signature. If there are duplicates we will use the first ABI that matched.
func (a *ABIFinder) FindABIByMethod(address string, signature []byte) (ABIFinderResult, error) {
	result := ABIFinderResult{}
	stringSignature := common.Bytes2Hex(signature)

	// let's start by checking if we already know what contract is at the address being called,
	// so that we don't have to search all known ABIs. If we have a match, let's double check
	// that it's correct. If it's not we will stop and return an error
	if a.ContractMap.IsKnownAddress(address) {
		contractName := a.ContractMap.GetContractName(address)
		abiInstanceCandidate, ok := a.ContractStore.ABIs[contractName+".abi"]
		if !ok {
			err := errors.New(ErrNoAbiFound)
			L.Err(err).
				Str("Contract", contractName).
				Str("Address", address).
				Msg("ABI not found, even though contract is known. This should not happen. Contract map might be corrupted")
			return ABIFinderResult{}, err
		}

		methodCandidate, err := abiInstanceCandidate.MethodById(signature)
		if err != nil {
			// this might still be a valid, case when contract A and B share the same method signature,
			// and we have already traced that call to contract A (when in fact it was contract B). If
			// the next call we are tracing is to contract B, then the ABI we have selected (belonging to A),
			// won't have it. In this case we should just continue and try to find the method in other ABIs.
			// In that case we should update our mapping, as now we came across a method that's (hopefully)
			// unique to contract B.
			for correctedContractName, correctedAbi := range a.ContractStore.ABIs {
				correctedMethod, abiErr := correctedAbi.MethodById(signature)
				if abiErr == nil {
					L.Debug().
						Str("Address", address).
						Str("Old ABI", contractName).
						Str("New ABI", correctedContractName).
						Str("Signature", stringSignature).
						Msgf("Updating contract mapping as previous one was based on non-unique method signature")

					a.ContractMap.AddContract(address, correctedContractName)

					result.Method = correctedMethod
					result.ABI = correctedAbi
					result.contractName = correctedContractName
					result.DuplicateCount = a.getDuplicateCount(signature)

					return result, nil
				}
			}

			L.Err(err).
				Str("Signature", stringSignature).
				Str("Supposed contract", contractName).
				Str("Supposed address", address).
				Msg("Method not found in known ABI instance. This should not happen. Contract map might be corrupted")

			return ABIFinderResult{}, err
		}

		result.Method = methodCandidate
		result.ABI = abiInstanceCandidate
		result.contractName = contractName
		result.DuplicateCount = 0 // we know the exact contract, so the duplicates here do not matter

		return result, nil
	} else {
		// if we do not know what contract is at given address we need to iterate over all known ABIs
		// and check if any of them has a method with the given signature (this might gave false positives,
		// when more than one contract has the same method signature, but we can't do anything about it)
		// In any case this should happen only when we did not deploy the contract via Seth (as otherwise we
		// know the address of the contract and can map it to the correct ABI instance).
		// If there are duplicates we will use the first ABI that matched.
		for abiName, abiInstanceCandidate := range a.ContractStore.ABIs {
			methodCandidate, err := abiInstanceCandidate.MethodById(signature)
			if err != nil {
				L.Trace().
					Err(err).
					Str("Signature", stringSignature).
					Msg("Method not found")
				continue
			}

			a.ContractMap.AddContract(address, abiName)

			result.ABI = abiInstanceCandidate
			result.Method = methodCandidate
			result.contractName = abiName
			result.DuplicateCount = a.getDuplicateCount(signature)

			break
		}

		if result.Method == nil {
			return ABIFinderResult{}, errors.New(ErrNoABIMethod)
		}
	}

	return result, nil
}

func (a *ABIFinder) getDuplicateCount(signature []byte) int {
	count := 0
	for _, abiInstance := range a.ContractStore.ABIs {
		_, err := abiInstance.MethodById(signature)
		if err == nil {
			count++
		}
	}

	return count - 1
}
