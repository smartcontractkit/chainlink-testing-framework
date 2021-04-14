package client

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

// KeyType specifies a few possible methods of retrieving configurations - namely wallets - from
type KeyType string

// ChainType specifies the possible underlying blockchain technologies, e.g. Ethereum
type ChainType string

const (
	NetworkEthereumHardhat   = "EthereumHardhat"
	WalletConfigFileLocation = "../wallets.yml" // Can make this an optional command line param?

	EnvKeyType            KeyType = "env"
	FileKeyType           KeyType = "file"
	SecretsManagerKeyType KeyType = "secret"

	EthereumChainType ChainType = "ethereum"
)

// Generalized blockchain client for interaction with multiple different blockchains
type BlockchainClient interface {
	DeployStorageContract(wallet BlockchainWallet) error
}

// NewBlockchainClient returns an implementation of a BlockchainClient based on the given network
func NewBlockchainClient(network BlockchainNetwork) (BlockchainClient, error) {
	switch network.(type) {
	case *EthereumHardhat:
		return NewEthereumClient(network)
	}
	return nil, fmt.Errorf("invalid blockchain network was given")
}

// BlockchainNetwork is the interface that when implemented, defines a new blockchain network that can be tested against
type BlockchainNetwork interface {
	Name() string
	URL() string
	ChainID() *big.Int
	ChainType() ChainType
	Wallets(KeyType) (BlockchainWallets, error)
}

// EthereumHardhat is the implementation of BlockchainNetwork for the local ETH dev server
type EthereumHardhat struct{}

// Name returns the readable name of the hardhat network
func (e *EthereumHardhat) Name() string {
	return NetworkEthereumHardhat
}

// URL returns the RPC URL used for connecting to hardhat
func (e *EthereumHardhat) URL() string {
	return ethereumURL(e)
}

// ChainID returns the on-chain ID of the network being connected to, returning hardhat's default
func (e *EthereumHardhat) ChainID() *big.Int {
	return big.NewInt(31337)
}

// ChainType returns the type of infrastrucructure the blockchain is built on, returning ethereum
func (e *EthereumHardhat) ChainType() ChainType {
	return EthereumChainType
}

// Wallets returns all the viable wallets used for testing on chain, returning hardhat's default
func (e *EthereumHardhat) Wallets(keyType KeyType) (BlockchainWallets, error) {
	walletString, err := retrieveWalletStrings(keyType, e.Name())
	if err != nil {
		return &Wallets{}, err
	}
	return processWalletStrings(walletString, e.ChainType())
}

// BlockchainWallets is an interface that when implemented is a representation of a slice of wallets for
// a specific network
type BlockchainWallets interface {
	Default() BlockchainWallet
	SetDefault(i int) error
	Wallet(i int) (BlockchainWallet, error)
}

// Wallets is the default implementation of BlockchainWallets that holds a slice of wallets with the default
type Wallets struct {
	defaultWallet int
	wallets       []BlockchainWallet
}

// Default returns the default wallet to be used for a transaction on-chain
func (w *Wallets) Default() BlockchainWallet {
	return w.wallets[w.defaultWallet]
}

// SetDefault changes the default wallet to be used for on-chain transactions
func (w *Wallets) SetDefault(i int) error {
	if err := walletSliceIndexInRange(w.wallets, i); err != nil {
		return err
	}
	w.defaultWallet = i
	return nil
}

// Wallet returns a wallet based on a given index in the slice
func (w *Wallets) Wallet(i int) (BlockchainWallet, error) {
	if err := walletSliceIndexInRange(w.wallets, i); err != nil {
		return nil, err
	}
	return w.wallets[i], nil
}

// BlockchainWallet when implemented is the interface to allow multiple wallet implementations for each
// BlockchainNetwork that is supported
type BlockchainWallet interface {
	PrivateKey() string
	Address() string
}

// EthereumWallet is the implementation to allow testing with ETH based wallets
type EthereumWallet struct {
	privateKey string
	address    common.Address
}

// NewEthereumWallet returns the instantiated ETH wallet based on a given private key
func NewEthereumWallet(pk string) (*EthereumWallet, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}
	return &EthereumWallet{
		privateKey: pk,
		address:    crypto.PubkeyToAddress(privateKey.PublicKey),
	}, nil
}

// PrivateKey returns the private key for a given Ethereum wallet
func (e *EthereumWallet) PrivateKey() string {
	return e.privateKey
}

// Address returns the ETH address for a given wallet
func (e *EthereumWallet) Address() string {
	return e.address.String()
}

func ethereumURL(network BlockchainNetwork) string {
	env := getNetworkURLEnv(network.ChainID())
	if len(env) > 0 {
		return env
	}
	return "http://localhost:8545"
}

func getNetworkURLEnv(chainID *big.Int) string {
	return os.Getenv(fmt.Sprintf("NETWORK_%d_URL", chainID))
}

func walletSliceIndexInRange(wallets []BlockchainWallet, i int) error {
	if i > len(wallets)-1 {
		return fmt.Errorf("invalid index in list of wallets")
	}
	return nil
}

// Retrieves a comma separated list of wallet private keys, depending on the config source and blockchain name
func retrieveWalletStrings(keyType KeyType, networkName string) (string, error) {
	var toProcess string
	var err error
	switch keyType {
	case EnvKeyType:
		toProcess = os.Getenv(networkName)
	case FileKeyType:
		keyFile, err := ioutil.ReadFile(WalletConfigFileLocation)
		if err != nil {
			return "", err
		}

		var config map[string]string
		err = yaml.Unmarshal(keyFile, &config)
		if err != nil {
			return "", err
		}
		toProcess = config[networkName]
	case SecretsManagerKeyType:
		// Get from whichever secrets manager we choose
	}
	return toProcess, err
}

// Processes a comma separated list of wallet private keys and gives back actual wallets based on blockchain type
func processWalletStrings(walletKeys string, blockchainType ChainType) (*Wallets, error) {
	var processedWallets []BlockchainWallet
	splitKeys := strings.Split(walletKeys, ",")

	for _, key := range splitKeys {
		switch blockchainType {
		case EthereumChainType:
			wallet, err := NewEthereumWallet(strings.TrimSpace(key))
			if err != nil {
				return &Wallets{}, err
			}
			processedWallets = append(processedWallets, wallet)
		}
	}

	return &Wallets{
		defaultWallet: 0,
		wallets:       processedWallets,
	}, nil
}
