// Package networks holds all known network information for the tests
package networks

import (
	"crypto/ecdsa"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/config"
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
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://source-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
		FinalityDepth:             1,
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
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
		FinalityDepth:             1,
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
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://source-chain-ethereum-besu:8546"},
		HTTPURLs:                  []string{"http://source-chain-ethereum-besu:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
		FinalityDepth:             1,
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
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
		},
		URLs:                      []string{"ws://dest-chain-ethereum-besu:8546"},
		HTTPURLs:                  []string{"http://dest-chain-ethereum-besu:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
		FinalityDepth:             1,
	}

	SimulatedEVMNonDev = blockchain.EVMNetwork{
		Name:                 "geth",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              1337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
		},
		URLs:                      []string{"ws://geth-ethereum-geth:8546"},
		HTTPURLs:                  []string{"http://geth-ethereum-geth:8544"},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
	}

	Anvil = blockchain.EVMNetwork{
		Name:                 "Anvil",
		Simulated:            true,
		SupportsEIP1559:      true,
		ClientImplementation: blockchain.EthereumClientImplementation,
		ChainID:              31337,
		PrivateKeys: []string{
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			"7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6",
			"47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a",
		},
		ChainlinkTransactionLimit: 500000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 5 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 5 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 5 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
	}

	// metisStardust https://www.metis.io/
	MetisStardust blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Stardust",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   588,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
	}

	MetisSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Metis Sepolia",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.MetisClientImplementation,
		ChainID:                   59902,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
	}

	ArbitrumMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Arbitrum Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.ArbitrumClientImplementation,
		ChainID:                   42161,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// OptimismGoerli https://dev.optimism.io/kovan-to-goerli/
	OptimismGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// OptimismSepolia https://community.optimism.io/docs/useful-tools/networks/#parameters-for-node-operators-2
	OptimismSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Optimism Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   11155420,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
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
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	PolygonAmoy blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Polygon Amoy",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.PolygonClientImplementation,
		ChainID:                   80002,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	AvalancheMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Avalanche Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43114,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	AvalancheFuji = blockchain.EVMNetwork{
		Name:                      "Avalanche Fuji",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   43113,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	Quorum = blockchain.EVMNetwork{
		Name:                      "Quorum",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.QuorumClientImplementation,
		ChainID:                   1337,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	// BaseSepolia https://base.mirror.xyz/kkz1-KFdUwl0n23PdyBRtnFewvO48_m-fZNzPMJehM4
	BaseSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   84532,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	ScrollSepolia = blockchain.EVMNetwork{
		Name:                      "Scroll Sepolia",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534351,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             400,
		DefaultGasLimit:           6000000,
	}

	ScrollMainnet = blockchain.EVMNetwork{
		Name:                      "Scroll Mainnet",
		ClientImplementation:      blockchain.ScrollClientImplementation,
		ChainID:                   534352,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityDepth:             400,
		DefaultGasLimit:           6000000,
	}

	CeloAlfajores = blockchain.EVMNetwork{
		Name:                      "Celo Alfajores",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   44787,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             10,
	}

	CeloMainnet = blockchain.EVMNetwork{
		Name:                      "Celo",
		ClientImplementation:      blockchain.CeloClientImplementation,
		ChainID:                   42220,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             10,
	}

	BaseMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Base Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.OptimismClientImplementation,
		ChainID:                   8453,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	BSCTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   97,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	BSCMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "BSC Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.BSCClientImplementation,
		ChainID:                   56,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      3,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	LineaGoerli blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Linea Goerli",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.LineaClientImplementation,
		ChainID:                   59140,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
	}

	LineaSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Linea Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.LineaClientImplementation,
		ChainID:                   59141,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityDepth:             200,
		DefaultGasLimit:           6000000,
	}

	LineaMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Linea Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.LineaClientImplementation,
		ChainID:                   59144,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityDepth:             200,
		DefaultGasLimit:           6000000,
	}

	PolygonZkEvmGoerli = blockchain.EVMNetwork{
		Name:                      "Polygon zkEVM Goerli",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.PolygonZkEvmClientImplementation,
		ChainID:                   1442,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           8000000,
	}

	PolygonZkEvmCardona = blockchain.EVMNetwork{
		Name:                      "Polygon zkEVM Cardona",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.PolygonZkEvmClientImplementation,
		ChainID:                   2442,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           8000000,
	}

	WeMixTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "WeMix Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.WeMixClientImplementation,
		ChainID:                   1112,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	WeMixMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "WeMix Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.WeMixClientImplementation,
		ChainID:                   1111,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       0,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	FantomTestnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Fantom Testnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.FantomClientImplementation,
		ChainID:                   4002,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
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
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	KromaSepolia blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Kroma Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.KromaClientImplementation,
		ChainID:                   2358,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	NexonDev blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Nexon Dev",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   955081,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	NexonTest blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Nexon Test",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   595581,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	NexonQa blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Nexon Qa",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   807424,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	NexonStage blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Nexon Stage",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   847799,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      0,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	GnosisChiado blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Gnosis Chiado",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.GnosisClientImplementation,
		ChainID:                   10200,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	GnosisMainnet blockchain.EVMNetwork = blockchain.EVMNetwork{
		Name:                      "Gnosis Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.GnosisClientImplementation,
		ChainID:                   100,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	ModeSepolia = blockchain.EVMNetwork{
		Name:                      "Mode Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   919,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	ModeMainnet = blockchain.EVMNetwork{
		Name:                      "Mode Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   34443,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       1000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	BlastSepolia = blockchain.EVMNetwork{
		Name:                      "Blast Sepolia",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   168587773,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	BlastMainnet = blockchain.EVMNetwork{
		Name:                      "Blast Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   81457,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	ZKSyncSepolia = blockchain.EVMNetwork{
		Name:                      "ZKSync Sepolia",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   300,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityDepth:             200,
		DefaultGasLimit:           6000000,
	}

	ZKSyncMainnet = blockchain.EVMNetwork{
		Name:                      "ZKSync Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   324,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityDepth:             1200,
		DefaultGasLimit:           6000000,
	}

	AstarShibuya = blockchain.EVMNetwork{
		Name:                      "Astar Shibuya",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   81,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           8000000,
	}

	AstarMainnet = blockchain.EVMNetwork{
		Name:                      "Astar Mainnet",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   592,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           8000000,
	}

	HederaTestnet = blockchain.EVMNetwork{
		Name:                      "Hedera Testnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   296,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	XLayerSepolia = blockchain.EVMNetwork{
		Name:                      "XLayer Sepolia",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   195,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	XLayerMainnet = blockchain.EVMNetwork{
		Name:                      "XLayer Mainnet",
		SupportsEIP1559:           false,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   196,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		FinalityTag:               true,
		DefaultGasLimit:           6000000,
	}

	TreasureRuby = blockchain.EVMNetwork{
		Name:                      "Treasure Ruby",
		SupportsEIP1559:           true,
		ClientImplementation:      blockchain.EthereumClientImplementation,
		ChainID:                   978657,
		Simulated:                 false,
		ChainlinkTransactionLimit: 5000,
		Timeout:                   blockchain.StrDuration{Duration: 3 * time.Minute},
		MinimumConfirmations:      1,
		GasEstimationBuffer:       10000,
		DefaultGasLimit:           6000000,
	}

	MappedNetworks = map[string]blockchain.EVMNetwork{
		"SIMULATED":               SimulatedEVM,
		"ANVIL":                   Anvil,
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
		"METIS_SEPOLIA":         MetisSepolia,
		"ARBITRUM_MAINNET":      ArbitrumMainnet,
		"ARBITRUM_GOERLI":       ArbitrumGoerli,
		"ARBITRUM_SEPOLIA":      ArbitrumSepolia,
		"OPTIMISM_MAINNET":      OptimismMainnet,
		"OPTIMISM_GOERLI":       OptimismGoerli,
		"OPTIMISM_SEPOLIA":      OptimismSepolia,
		"BASE_GOERLI":           BaseGoerli,
		"BASE_SEPOLIA":          BaseSepolia,
		"CELO_ALFAJORES":        CeloAlfajores,
		"CELO_MAINNET":          CeloMainnet,
		"RSK":                   RSKTestnet,
		"POLYGON_MUMBAI":        PolygonMumbai,
		"POLYGON_AMOY":          PolygonAmoy,
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
		"LINEA_SEPOLIA":         LineaSepolia,
		"LINEA_MAINNET":         LineaMainnet,
		"POLYGON_ZKEVM_GOERLI":  PolygonZkEvmGoerli,
		"POLYGON_ZKEVM_CARDONA": PolygonZkEvmCardona,
		"POLYGON_ZKEVM_MAINNET": PolygonZkEvmMainnet,
		"FANTOM_TESTNET":        FantomTestnet,
		"FANTOM_MAINNET":        FantomMainnet,
		"WEMIX_TESTNET":         WeMixTestnet,
		"WEMIX_MAINNET":         WeMixMainnet,
		"KROMA_SEPOLIA":         KromaSepolia,
		"KROMA_MAINNET":         KromaMainnet,
		"NEXON_DEV":             NexonDev,
		"NEXON_TEST":            NexonTest,
		"NEXON_QA":              NexonQa,
		"NEXON_STAGE":           NexonStage,
		"GNOSIS_CHIADO":         GnosisChiado,
		"GNOSIS_MAINNET":        GnosisMainnet,
		"BLAST_SEPOLIA":         BlastSepolia,
		"BLAST_MAINNET":         BlastMainnet,
		"MODE_SEPOLIA":          ModeSepolia,
		"MODE_MAINNET":          ModeMainnet,
		"ZKSYNC_SEPOLIA":        ZKSyncSepolia,
		"ZKSYNC_MAINNET":        ZKSyncMainnet,
		"ASTAR_SHIBUYA":         AstarShibuya,
		"ASTAR_MAINNET":         AstarMainnet,
		"HEDERA_TESTNET":        HederaTestnet,
		"XLAYER_SEPOLIA":        XLayerSepolia,
		"XLAYER_MAINNET":        XLayerMainnet,
		"TREASURE_RUBY":         TreasureRuby,
	}
)

// Get []blockchain.EVMNetwork from TOML config. Panic if no networks are found
func MustGetSelectedNetworkConfig(networkCfg *config.NetworkConfig) []blockchain.EVMNetwork {
	if networkCfg == nil || len(networkCfg.SelectedNetworks) == 0 {
		panic(fmt.Errorf("network config has no or empty selected networks. Use valid network(s) separated by comma from %v", getValidNetworkKeys()))
	}
	nets, err := SetNetworks(*networkCfg)
	if err != nil {
		panic(err)
	}
	return nets
}

func SetNetworks(networkCfg config.NetworkConfig) ([]blockchain.EVMNetwork, error) {
	networks := make([]blockchain.EVMNetwork, 0)
	selectedNetworks := networkCfg.SelectedNetworks
	for i := range selectedNetworks {
		var walletKeys, httpUrls, wsUrls []string
		networkName := strings.ToUpper(selectedNetworks[i])
		forked := false
		if networkCfg.AnvilConfigs != nil {
			_, forked = networkCfg.AnvilConfigs[networkName]
		}

		// if network is not simulated or forked, use the rpc urls and wallet keys from config
		if !strings.Contains(networkName, "SIMULATED") && !forked {
			var ok, wsOk, httpOk bool
			// Check for WS URLs
			wsUrls, wsOk = networkCfg.RpcWsUrls[selectedNetworks[i]]
			// Check for HTTP URLs
			httpUrls, httpOk = networkCfg.RpcHttpUrls[selectedNetworks[i]]

			// WS can be present but only if HTTP is also available, the CL node cannot function only on WS
			if wsOk && !httpOk {
				return nil, fmt.Errorf("WS RPC endpoint for %s network is set without an HTTP endpoint; only HTTP or both HTTP and WS are allowed", selectedNetworks[i])
			}

			// Validate that there is at least one HTTP endpoint
			if !httpOk {
				return nil, fmt.Errorf("at least one HTTP RPC endpoint for %s network must be set", selectedNetworks[i])
			}

			walletKeys, ok = networkCfg.WalletKeys[selectedNetworks[i]]
			if !ok {
				return nil, fmt.Errorf("no wallet keys found in config for '%s' network", selectedNetworks[i])
			}
		}
		// if evm_network config is found, use it
		if networkCfg.EVMNetworks != nil {
			if network, ok := networkCfg.EVMNetworks[networkName]; ok && network != nil {
				if err := NewEVMNetwork(network, walletKeys, httpUrls, wsUrls); err != nil {
					return nil, err
				}
				networks = append(networks, *network)
				continue
			}
		}
		// if there is no evm_network config, use the known networks to find the network config from the map
		if knownNetwork, valid := MappedNetworks[networkName]; valid {
			err := NewEVMNetwork(&knownNetwork, walletKeys, httpUrls, wsUrls)
			if err != nil {
				return nil, err
			}
			networks = append(networks, knownNetwork)
			continue
		}
		// if network is not found in known networks or in toml's evm_network config, throw an error
		return nil, fmt.Errorf("no evm_network config found in network config. "+
			"network '%s' is not a valid network. "+
			"Use valid network(s) separated by comma from %v "+
			"or add the evm_network details to the network config file",
			selectedNetworks[i], getValidNetworkKeys())
	}
	return networks, nil
}

// NewEVMNetwork sets the network's private key(s) and rpc urls
func NewEVMNetwork(network *blockchain.EVMNetwork, walletKeys, httpUrls, wsUrls []string) error {
	if len(httpUrls) > 0 {
		network.HTTPURLs = httpUrls
	}
	if len(wsUrls) > 0 {
		network.URLs = wsUrls
	}
	if len(walletKeys) > 0 {
		if err := setKeys(network, walletKeys); err != nil {
			return err
		}
	}
	return nil
}

func getValidNetworkKeys() []string {
	validKeys := make([]string, 0)
	for validNetwork := range MappedNetworks {
		validKeys = append(validKeys, validNetwork)
	}
	return validKeys
}

// setKeys sets a network's private key(s) based on env vars
func setKeys(network *blockchain.EVMNetwork, walletKeys []string) error {
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
			return fmt.Errorf("error converting private key to public key: %w", err)
		}
		publicKeys = append(publicKeys, publicKey)
	}
	log.Info().Interface("Funding Addresses", publicKeys).Msg("Read Network Keys")
	return nil
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
