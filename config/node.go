package config

type NodeConfig struct {
	BaseConfigTOML           string            `toml:",omitempty"`
	CommonChainConfigTOML    string            `toml:",omitempty"`
	ChainConfigTOMLByChainID map[string]string `toml:",omitempty"` // key is chainID
	NumberOfNodes            *int              `toml:",omitempty"`
	UseExisting              *bool             `toml:"use_existing,omitempty"`
	Nodes                    []NodeDetails     `toml:"Node,omitempty"`
	Namespace                *string           `toml:"namespace,omitempty"`
}

type NodeDetails struct {
	URL        *string `toml:"url,omitempty"`
	InternalIP *string `toml:"internal_ip,omitempty"`
	Email      *string
	Password   *string
}
