package types

type ExecutionLayer string

const (
	ExecutionLayer_Geth       ExecutionLayer = "geth"
	ExecutionLayer_Nethermind ExecutionLayer = "nethermind"
	ExecutionLayer_Erigon     ExecutionLayer = "erigon"
	ExecutionLayer_Besu       ExecutionLayer = "besu"
	ExecutionLayer_Reth       ExecutionLayer = "reth"
)
