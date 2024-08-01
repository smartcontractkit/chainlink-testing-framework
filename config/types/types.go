package types

type ExecutionLayer string

const (
	ExecutionLayer_Geth       ExecutionLayer = "geth"
	ExecutionLayer_Nethermind ExecutionLayer = "nethermind"
	ExecutionLayer_Erigon     ExecutionLayer = "erigon"
	ExecutionLayer_Besu       ExecutionLayer = "besu"
	ExecutionLayer_Reth       ExecutionLayer = "reth"
)

type EthereumVersion string

const (
	EthereumVersion_Eth2 EthereumVersion = "eth2"
	// Deprecated: use EthereumVersion_Eth2 instead
	EthereumVersion_Eth2_Legacy EthereumVersion = "pos"
	EthereumVersion_Eth1        EthereumVersion = "eth1"
	// Deprecated: use EthereumVersion_Eth1 instead
	EthereumVersion_Eth1_Legacy EthereumVersion = "pow"
)
