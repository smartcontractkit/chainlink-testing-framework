package seth

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	ErrOpenABIFile = "failed to open ABI file"
	ErrParseABI    = "failed to parse ABI file"
	ErrOpenBINFile = "failed to open BIN file"
)

// ContractStore contains all ABIs that are used in decoding. It might also contain contract bytecode for deployment
type ContractStore struct {
	ABIs ABIStore
	BINs map[string][]byte
	mu   *sync.RWMutex
}

type ABIStore map[string]abi.ABI

func (c *ContractStore) GetABI(name string) (*abi.ABI, bool) {
	if !strings.HasSuffix(name, ".abi") {
		name = name + ".abi"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	abi, ok := c.ABIs[name]
	return &abi, ok
}

func (c *ContractStore) AddABI(name string, abi abi.ABI) {
	if !strings.HasSuffix(name, ".abi") {
		name = name + ".abi"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ABIs[name] = abi
}

func (c *ContractStore) GetBIN(name string) ([]byte, bool) {
	if !strings.HasSuffix(name, ".bin") {
		name = name + ".bin"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	bin, ok := c.BINs[name]
	return bin, ok
}

func (c *ContractStore) AddBIN(name string, bin []byte) {
	if !strings.HasSuffix(name, ".bin") {
		name = name + ".bin"
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.BINs[name] = bin
}

// NewContractStore creates a new Contract store
func NewContractStore(abiPath, binPath string) (*ContractStore, error) {
	cs := &ContractStore{ABIs: make(ABIStore), BINs: make(map[string][]byte), mu: &sync.RWMutex{}}

	if abiPath != "" {
		files, err := os.ReadDir(abiPath)
		if err != nil {
			return nil, err
		}
		var foundABI bool
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".abi") {
				L.Debug().Str("File", f.Name()).Msg("ABI file loaded")
				ff, err := os.Open(filepath.Join(abiPath, f.Name()))
				if err != nil {
					return nil, errors.Wrap(err, ErrOpenABIFile)
				}
				a, err := abi.JSON(ff)
				if err != nil {
					return nil, errors.Wrap(err, ErrParseABI)
				}
				cs.ABIs[f.Name()] = a
				foundABI = true
			}
		}
		if !foundABI {
			L.Warn().Msg("No ABI files found")
			L.Warn().Msg("You will need to provide the bytecode manually, when deploying contracts")
		}
	}

	if binPath != "" {
		files, err := os.ReadDir(binPath)
		if err != nil {
			return nil, err
		}
		var foundBIN bool
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".bin") {
				L.Debug().Str("File", f.Name()).Msg("BIN file loaded")
				bin, err := os.ReadFile(filepath.Join(binPath, f.Name()))
				if err != nil {
					return nil, errors.Wrap(err, ErrOpenBINFile)
				}
				cs.BINs[f.Name()] = common.FromHex(string(bin))
				foundBIN = true
			}
		}
		if !foundBIN {
			L.Warn().Msg("No BIN files found")
			L.Warn().Msg("You will need to provide the bytecode manually, when deploying contracts")
		}
	}

	return cs, nil
}
