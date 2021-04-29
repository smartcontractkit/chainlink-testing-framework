package client

type ClientStore interface {
	Client(ContractInstance) BlockchainClient
}

type ContractInstance interface {
	NetworkType() string
}
