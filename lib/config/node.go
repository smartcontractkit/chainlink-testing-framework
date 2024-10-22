package config

type NodeConfig struct {
	BaseConfigTOML           string            `toml:",omitempty"`
	CommonChainConfigTOML    string            `toml:",omitempty"`
	ChainConfigTOMLByChainID map[string]string `toml:",omitempty"` // key is chainID
}
