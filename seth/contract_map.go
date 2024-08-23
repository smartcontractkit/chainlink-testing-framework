package seth

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pelletier/go-toml/v2"
)

type ContractMap struct {
	mu         *sync.RWMutex
	addressMap map[string]string
}

func NewEmptyContractMap() ContractMap {
	return ContractMap{
		mu:         &sync.RWMutex{},
		addressMap: map[string]string{},
	}
}

func NewContractMap(contracts map[string]string) ContractMap {
	return ContractMap{
		mu:         &sync.RWMutex{},
		addressMap: contracts,
	}
}

func (c ContractMap) GetContractMap() map[string]string {
	return c.addressMap
}

func (c ContractMap) IsKnownAddress(addr string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.addressMap[strings.ToLower(addr)] != ""
}

func (c ContractMap) GetContractName(addr string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.addressMap[strings.ToLower(addr)]
}

func (c ContractMap) GetContractAddress(addr string) string {
	if addr == UNKNOWN {
		return UNKNOWN
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.addressMap {
		if v == addr {
			return k
		}
	}
	return UNKNOWN
}

func (c ContractMap) AddContract(addr, name string) {
	if addr == UNKNOWN {
		return
	}

	name = strings.TrimSuffix(name, ".abi")
	c.mu.Lock()
	defer c.mu.Unlock()
	c.addressMap[strings.ToLower(addr)] = name
}

func (c ContractMap) Size() int {
	return len(c.addressMap)
}

func SaveDeployedContract(filename, contractName, address string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}
	defer file.Close()

	v := map[string]string{
		address: contractName,
	}

	marhalled, err := toml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(marhalled))
	return err
}

func LoadDeployedContracts(filename string) (map[string]string, error) {
	tomlFile, err := os.Open(filename)
	if err != nil {
		return map[string]string{}, nil
	}
	defer tomlFile.Close()

	b, _ := io.ReadAll(tomlFile)
	rawContracts := map[common.Address]string{}
	err = toml.Unmarshal(b, &rawContracts)
	if err != nil {
		return map[string]string{}, err
	}

	contracts := map[string]string{}
	for k, v := range rawContracts {
		contracts[k.Hex()] = v
	}

	return contracts, nil
}
