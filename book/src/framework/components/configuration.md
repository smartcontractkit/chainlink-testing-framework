## Test Configuration

In our framework, components control the configuration. Each component defines an Input that you embed into your test configuration. We automatically ensure consistency between TOML and struct definitions during validation.

An example of your component configuration:
```
// Input is a blockchain network configuration params declared by the component
type Input struct {
	Type                     string   `toml:"type" validate:"required,oneof=anvil geth" envconfig:"net_type"`
	Image                    string   `toml:"image" validate:"required"`
	Tag                      string   `toml:"tag" validate:"required"`
	Port                     string   `toml:"port" validate:"required"`
	ChainID                  string   `toml:"chain_id" validate:"required"`
	DockerCmdParamsOverrides []string `toml:"docker_cmd_params"`
	Out                      *Output  `toml:"out"`
}
```

How you use it in tests:
```
type Config struct {
	BlockchainA        *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestDON(t *testing.T) {
	in, err := framework.Load[Config](t)
	require.NoError(t, err)

	// deploy component
	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)
```

In `TOML`:
```
[blockchain_a]
chain_id = "1337"
image = "f4hrenh9it/foundry:latest"
port = "8500"
type = "anvil"
docker_cmd_params = ["-b", "1"]
```

### Best practices for configuration and validation
- Avoid stateful types (e.g., loggers, clients) in your config.
- All `input` fields should include validate: "required", ensuring consistency between TOML and struct definitions.
- Add extra validation rules for URLs or "one-of" variants. Learn more here: go-playground/validator.

### Overriding configuration
To override any configuration, we merge multiple files into a single struct.

You can specify multiple file paths using `CTF_CONFIGS=path1,path2,path3`.

The framework will apply these configurations from right to left.