// Package networks holds all known network information for the tests
package networks

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// Pre-configured test networks and their connections
// Some networks with public RPC endpoints are already filled out, but make use of environment variables to use info like
// private RPC endpoints and private keys.
var (
	// To create replica of simulated EVM network, with different chain ids
	AdditionalSimulatedChainIds = []int64{3337, 4337, 5337, 6337, 7337, 8337, 9337, 9338}
	AdditionalSimulatedPvtKeys  = []string{
		"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
		"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
		"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		"8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba",
		"92db14e403b83dfe3df233f83dfa3a0d7096f21ca9b0d6d6b8d88b2b4ec1564e",
		"4bbbf85ce3377467afe5d46f804f221813b2bb87f24d81f60f1fcdbf7cbf4356",
		"dbda1821b80551c9d65939329250298aa3472ba22feea921c0cf5d620ea67b97",
		"2a871d0798f97d79848a013d4936a73bf4cc922c825d33c1cf7073dff6d409c6",
		"f214f2b2cd398c806f84e317254e0f0b801d0643303237d97a22a48e01628897",
		"701b615bbdfb9de65240bc28bd21bbc0d996645a3dd57e7b12bc2bdf6f192c82",
		"a267530f49f8280200edf313ee7af6b827f2a8bce2897751d06a843f644967b1",
		"47c99abed3324a2707c28affff1267e45918ec8c3f20b8aa892e8b065d2942dd",
		"c526ee95bf44d8fc405a158bb884d9d1238d99f0612e9f33d006bb0789009aaa",
		"8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61",
		"ea6c44ac03bff858b476bba40716402b03e41b8e97e276d1baec7c37d42484a0",
		"689af8efa8c651a91ad287602527f3af2fe9f6501a7ac4b061667b5a93e037fd",
		"de9be858da4a475276426320d5e9262ecfc3ba460bfac56360bfa6c4c28b4ee0",
		"df57089febbacf7ba0bc227dafbffa9fc08a93fdc68e1e42411a14efcf23656e",
	}

	// SimulatedEVM represents a simulated network
	SimulatedEVM blockchain.EVMNetwork = blockchain.SimulatedEVMNetwork
	// generalEVM is a customizable network through environment variables
	// This is getting little use, and causes some confusion. Can re-enable if people want it.
	// generalEVM blockchain.EVMNetwork = blockchain.LoadNetworkFromEnvironment()

	// SimulatedevmNonDev1 represents a simulated network which can be used to deploy a non-dev geth node
	SimulatedEVMNonDev1 = blockchain.EVMNetwork{
		Name:                 "source-chain",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		SupportsEIP1559:      true,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://source-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	// SimulatedEVM_NON_DEV_2 represents a simulated network with chain id 2337 which can be used to deploy a non-dev geth node
	SimulatedEVMNonDev2 = blockchain.EVMNetwork{
		Name:                 "dest-chain",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              2337,
		PrivateKeys: []string{
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	// SimulatedBesuNonDev1 represents a simulated network which can be used to deploy a non-dev besu node
	// in a CCIP source-chain -> dest-chain communication
	SimulatedBesuNonDev1 = blockchain.EVMNetwork{
		Name:                 "source-chain",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		SupportsEIP1559:      false,
		SimulationType:       "besu",
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://source-chain-ethereum-besu:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-besu:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	// SimulatedBesuNonDev2 represents a simulated network which can be used to deploy a non-dev besu node
	// in a CCIP source-chain -> dest-chain communication
	SimulatedBesuNonDev2 = blockchain.EVMNetwork{
		Name:                 "dest-chain",
		Simulated:            true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		SupportsEIP1559:      false,
		SimulationType:       "besu",
		ChainID:              2337,
		PrivateKeys: []string{
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-besu:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-besu:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	SimulatedEVMNonDev = blockchain.EVMNetwork{
		Name:                 "geth",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://geth-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://geth-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	EthereumMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Ethereum Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   1,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// sepoliaTestnet https://sepolia.dev/
	SepoliaTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Sepolia Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   11155111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// goerliTestnet https://goerli.net/
	GoerliTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Goerli Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   5,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 5 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	KlaytnMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Klaytn Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   8217,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// klaytnBaobab https://klaytn.foundation/
	KlaytnBaobab blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Klaytn Baobab",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.KlaytnClientImplementation,
		ChainID:                   1001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	MetisAndromeda blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Andromeda",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   1088,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	// metisStardust https://www.metis.io/
	MetisStardust blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   588,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	ArbitrumMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   42161,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           100000000,
	}

	// arbitrumGoerli https://developer.offchainlabs.com/docs/public_chains
	ArbitrumGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   421613,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           100000000,
	}

	ArbitrumSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   421614,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           100000000,
	}

	OptimismMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   10,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// optimismGoerli https://dev.optimism.io/kovan-to-goerli/
	OptimismGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	RSKMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "RSK Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.RSKClientImplementation,
		ChainID:                   30,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	// rskTestnet https://www.rsk.co/
	RSKTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "RSK Testnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.RSKClientImplementation,
		ChainID:                   31,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	PolygonMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   137,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             550,
		DefaultGasLimit:           6000000,
	}

	// PolygonMumbai https://mumbai.polygonscan.com/
	PolygonMumbai blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Mumbai",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   80001,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             550,
		DefaultGasLimit:           6000000,
	}

	AvalancheMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Avalanche Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43114,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             35,
		DefaultGasLimit:           6000000,
	}

	AvalancheFuji = blockchain.EVMNetwork{
		Name:                      "Avalanche Fuji",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43113,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             35,
		DefaultGasLimit:           6000000,
	}

	Quorum = blockchain.EVMNetwork{
		Name:                      "Quorum",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.QuorumClientImplementation,
		ChainID:                   1337,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	BaseGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   84531,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
	}

	CeloAlfajores = blockchain.EVMNetwork{
		Name:                      "Celo Alfajores",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   44787,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	ScrollSepolia = blockchain.EVMNetwork{
		Name:                      "Scroll Sepolia",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534351,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	ScrollMainnet = blockchain.EVMNetwork{
		Name:                      "Scroll Mainnet",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534352,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
	}

	CeloMainnet = blockchain.EVMNetwork{
		Name:                      "Celo",
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   42220,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	BaseMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   8453,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
	}

	BSCTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   97,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
	}

	BSCMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   56,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
	}

	LineaGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Linea Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.LineaClientImplementation,
		ChainID:                   59140,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	LineaMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Linea Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.LineaClientImplementation,
		ChainID:                   59144,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	PolygonZkEvmGoerli = blockchain.EVMNetwork{
		Name:                      "Polygon zkEVM Goerli",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.PolygonZkEvmClientImplementation,
		ChainID:                   1442,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
	}

	PolygonZkEvmMainnet = blockchain.EVMNetwork{
		Name:                      "Polygon zkEVM Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.PolygonZkEvmClientImplementation,
		ChainID:                   1101,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
	}

	WeMixTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "WeMix Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.WeMixClientImplementation,
		ChainID:                   1112,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             1,
	}

	WeMixMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "WeMix Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.WeMixClientImplementation,
		ChainID:                   1111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             1,
	}

	FantomTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Fantom Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.FantomClientImplementation,
		ChainID:                   4002,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	FantomMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Fantom Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.FantomClientImplementation,
		ChainID:                   250,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	KromaMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Kroma Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.KromaClientImplementation,
		ChainID:                   255,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
	}

	KromaSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Kroma Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.KromaClientImplementation,
		ChainID:                   2358,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.JSONStrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
	}

	MappedNetworks = map[string]blockchain.EVMNetwork{
		"SIMULATED":               SimulatedEVM,
		"SIMULATED_1":             SimulatedEVMNonDev1,
		"SIMULATED_2":             SimulatedEVMNonDev2,
		"SIMULATED_BESU_NONDEV_1": SimulatedBesuNonDev1,
		"SIMULATED_BESU_NONDEV_2": SimulatedBesuNonDev2,
		"SIMULATED_NONDEV":        SimulatedEVMNonDev,
		// "GENERAL":         generalEVM, // See above
		"ETHEREUM_MAINNET":      EthereumMainnet,
		"GOERLI":                GoerliTestnet,
		"SEPOLIA":               SepoliaTestnet,
		"KLAYTN_MAINNET":        KlaytnMainnet,
		"KLAYTN_BAOBAB":         KlaytnBaobab,
		"METIS_ANDROMEDA":       MetisAndromeda,
		"METIS_STARDUST":        MetisStardust,
		"ARBITRUM_MAINNET":      ArbitrumMainnet,
		"ARBITRUM_GOERLI":       ArbitrumGoerli,
		"ARBITRUM_SEPOLIA":      ArbitrumSepolia,
		"OPTIMISM_MAINNET":      OptimismMainnet,
		"OPTIMISM_GOERLI":       OptimismGoerli,
		"BASE_GOERLI":           BaseGoerli,
		"CELO_ALFAJORES":        CeloAlfajores,
		"CELO_MAINNET":          CeloMainnet,
		"RSK":                   RSKTestnet,
		"POLYGON_MUMBAI":        PolygonMumbai,
		"POLYGON_MAINNET":       PolygonMainnet,
		"AVALANCHE_FUJI":        AvalancheFuji,
		"AVALANCHE_MAINNET":     AvalancheMainnet,
		"QUORUM":                Quorum,
		"SCROLL_SEPOLIA":        ScrollSepolia,
		"SCROLL_MAINNET":        ScrollMainnet,
		"BASE_MAINNET":          BaseMainnet,
		"BSC_TESTNET":           BSCTestnet,
		"BSC_MAINNET":           BSCMainnet,
		"LINEA_GOERLI":          LineaGoerli,
		"LINEA_MAINNET":         LineaMainnet,
		"POLYGON_ZKEVM_GOERLI":  PolygonZkEvmGoerli,
		"POLYGON_ZKEVM_MAINNET": PolygonZkEvmMainnet,
		"FANTOM_TESTNET":        FantomTestnet,
		"FANTOM_MAINNET":        FantomMainnet,
		"WEMIX_TESTNET":         WeMixTestnet,
		"WEMIX_MAINNET":         WeMixMainnet,
		"KROMA_SEPOLIA":         KromaSepolia,
		"KROMA_MAINNET":         KromaMainnet,
	}
)

// Get []blockchain.EVMNetwork from env vars. Panic if env vars not set or no networks found
func MustGetSelectedNetworksFromEnv() []blockchain.EVMNetwork {
	selectedNetworksEnv := os.Getenv("SELECTED_NETWORKS")
	emptyEnvErr := errors.Errorf("env var 'SELECTED_NETWORKS' is not set or is empty. Use valid network(s) separated by comma from %v", getValidNetworkKeys())
	if selectedNetworksEnv == "" {
		panic(emptyEnvErr)
	}
	networkKeys := strings.Split(selectedNetworksEnv, ",")
	if len(networkKeys) == 0 {
		panic(emptyEnvErr)
	}
	networks := make([]blockchain.EVMNetwork, 0)
	for i := range networkKeys {
		var walletKeys, httpUrls, wsUrls []string
		if !strings.Contains(networkKeys[i], "SIMULATED") {
			// Get network RPC WS URL from env var
			wsEnvVar := fmt.Sprintf("%s_URLS", networkKeys[i])
			wsEnvVal := os.Getenv(wsEnvVar)
			if wsEnvVal == "" {
				// Get default value
				defaultUrls, err := utils.GetEnv("EVM_URLS")
				if err != nil {
					panic(errors.Errorf("error getting %s EVM_URLS var", err))
				}
				if defaultUrls == "" {
					panic(errors.Errorf("set %s or EVM_URLS env var", wsEnvVar))
				}
				log.Warn().Msgf("%s not set, defaulting to EVM_URLS", wsEnvVar)
				wsUrls = strings.Split(defaultUrls, ",")
			} else {
				wsUrls = strings.Split(wsEnvVal, ",")
			}

			// Get network RPC HTTP URL from env var
			httpEnvVar := fmt.Sprintf("%s_HTTP_URLS", networkKeys[i])
			httpEnvVal := os.Getenv(httpEnvVar)
			if httpEnvVal == "" {
				// Get default value
				defaultUrls, err := utils.GetEnv("EVM_HTTP_URLS")
				if err != nil {
					panic(errors.Errorf("error getting %s EVM_HTTP_URLS var", err))
				}
				if defaultUrls == "" {
					panic(errors.Errorf("set %s or EVM_HTTP_URLS env var", httpEnvVar))
				}
				log.Warn().Msgf("%s not set, defaulting to EVM_HTTP_URLS", httpEnvVar)
				httpUrls = strings.Split(defaultUrls, ",")
			} else {
				httpUrls = strings.Split(httpEnvVal, ",")
			}

			// Get network wallet key from env var
			walletKeysEnvVar := fmt.Sprintf("%s_KEYS", networkKeys[i])
			walletKeysEnvVal := os.Getenv(walletKeysEnvVar)
			if walletKeysEnvVal == "" {
				// Get default value
				defaultKeys, err := utils.GetEnv("EVM_KEYS")
				if err != nil {
					panic(errors.Errorf("error getting EVM_KEYS var: %s", err))
				}
				if defaultKeys == "" {
					panic(errors.Errorf("set %s or EVM_KEYS env var", walletKeysEnvVar))
				}
				log.Warn().Msgf("%s not set, defaulting to EVM_KEYS", walletKeysEnvVar)
				walletKeys = strings.Split(defaultKeys, ",")
			} else {
				walletKeys = strings.Split(walletKeysEnvVal, ",")
			}
		}
		network, err := NewEVMNetwork(networkKeys[i], walletKeys, httpUrls, wsUrls)
		if err != nil {
			panic(err)
		}
		networks = append(networks, network)
	}
	return networks
}

func NewEVMNetwork(networkKey string, walletKeys, httpUrls, wsUrls []string) (blockchain.EVMNetwork, error) {
	if network, valid := MappedNetworks[networkKey]; valid {
		// Overwrite network default values
		if len(httpUrls) > 0 {
			network.HTTPURLs = httpUrls
		}
		if len(wsUrls) > 0 {
			network.URLs = wsUrls
		}
		if len(walletKeys) > 0 {
			setKeys(&network, walletKeys)
		}
		return network, nil
	}
	return blockchain.EVMNetwork{}, errors.Errorf("network key: '%v' is invalid. Use a valid network(s) separated by comma from %v", networkKey, getValidNetworkKeys())
}

func getValidNetworkKeys() []string {
	validKeys := make([]string, 0)
	for validNetwork := range MappedNetworks {
		validKeys = append(validKeys, validNetwork)
	}
	return validKeys
}

// setKeys sets a network's private key(s) based on env vars
func setKeys(network *blockchain.EVMNetwork, walletKeys []string) {
	for keyIndex := range walletKeys { // Sanitize keys of possible `0x` prefix
		// Trim some common addons
		walletKeys[keyIndex] = strings.Trim(walletKeys[keyIndex], "\"'")
		walletKeys[keyIndex] = strings.TrimSpace(walletKeys[keyIndex])
		walletKeys[keyIndex] = strings.TrimPrefix(walletKeys[keyIndex], "0x")
	}
	network.PrivateKeys = walletKeys

	// log public keys for debugging
	publicKeys := []string{}
	for _, key := range network.PrivateKeys {
		publicKey, err := privateKeyToAddress(key)
		if err != nil {
			log.Fatal().Err(err).Msg("Error reading private key")
		}
		publicKeys = append(publicKeys, publicKey)
	}
	log.Info().Interface("Funding Addresses", publicKeys).Msg("Read Network Keys")
}

func privateKeyToAddress(privateKeyString string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting private key to public ECDSA key")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex(), nil
}
