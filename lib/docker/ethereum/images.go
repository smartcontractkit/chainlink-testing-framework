package ethereum

const (
	DefaultBesuEth1Image = "hyperledger/besu:22.1.0"
	DefaultBesuEth2Image = "hyperledger/besu:25.2.0"
	BesuBaseImageName    = "hyperledger/besu"
	besuGitRepo          = "hyperledger/besu"

	DefaultErigonEth1Image = "thorax/erigon:v2.40.0"
	DefaultErigonEth2Image = "thorax/erigon:2.60.8"
	ErigonBaseImageName    = "thorax/erigon"
	erigonGitRepo          = "ledgerwatch/erigon"

	DefaultGethEth1Image = "ethereum/client-go:v1.13.8"
	DefaultGethEth2Image = "ethereum/client-go:v1.15.0"
	GethBaseImageName    = "ethereum/client-go"
	gethGitRepo          = "ethereum/go-ethereum"

	DefaultNethermindEth1Image = "nethermind/nethermind:1.16.0"
	DefaultNethermindEth2Image = "nethermind/nethermind:1.31.0"
	NethermindBaseImageName    = "nethermind/nethermind"
	nethermindGitRepo          = "NethermindEth/nethermind"

	DefaultRethEth2Image = "ghcr.io/paradigmxyz/reth:v1.1.5"
	RethBaseImageName    = "ghcr.io/paradigmxyz/reth"
	rethGitRepo          = "paradigmxyz/reth"

	GenesisGeneratorDenebImage    = "tofelb/ethereum-genesis-generator:3.3.5-main-8a8fb99" // latest one, copy of public.ecr.aws/w0i8p0z9/ethereum-genesis-generator:main-8a8fb99
	GenesisGeneratorShanghaiImage = "tofelb/ethereum-genesis-generator:2.0.5"
)
