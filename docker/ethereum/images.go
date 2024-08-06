package ethereum

const (
	DefaultBesuEth1Image = "hyperledger/besu:22.1.0"
	DefaultBesuEth2Image = "hyperledger/besu:24.5.1"
	BesuBaseImageName    = "hyperledger/besu"
	besuGitRepo          = "hyperledger/besu"

	DefaultErigonEth1Image = "thorax/erigon:v2.40.0"
	DefaultErigonEth2Image = "thorax/erigon:v2.59.3" // v.2.60.0 is the latest, but gas estimations using zero address are broken
	ErigonBaseImageName    = "thorax/erigon"
	erigonGitRepo          = "ledgerwatch/erigon"

	DefaultGethEth1Image = "ethereum/client-go:v1.13.8"
	DefaultGethEth2Image = "ethereum/client-go:v1.14.3"
	GethBaseImageName    = "ethereum/client-go"
	gethGitRepo          = "ethereum/go-ethereum"

	DefaultNethermindEth1Image = "nethermind/nethermind:1.16.0"
	DefaultNethermindEth2Image = "nethermind/nethermind:1.26.0"
	NethermindBaseImageName    = "nethermind/nethermind"
	nethermindGitRepo          = "NethermindEth/nethermind"

	DefaultRethEth2Image = "ghcr.io/paradigmxyz/reth:v1.0.0"
	RethBaseImageName    = "ghcr.io/paradigmxyz/reth"
	rethGitRepo          = "paradigmxyz/reth"
)
