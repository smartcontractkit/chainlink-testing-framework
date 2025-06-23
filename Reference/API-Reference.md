# API Reference

This page provides a high-level API reference for the main packages and components in the Chainlink Testing Framework (CTF).

---

## Framework Core

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/framework`
- **Key Types:**
  - `func Load[T any](t *testing.T) (T, error)` – Loads TOML config into Go struct
  - `type Logger` – Structured logging
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/framework)

---

## Blockchain Components

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain`
- **Key Types:**
  - `type Input` – Blockchain config
  - `func NewBlockchainNetwork(input *Input) (*Network, error)`
  - `type Network` – Blockchain network instance
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain)

---

## Chainlink Node Components

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode`
- **Key Types:**
  - `type Input` – Chainlink node config
  - `func NewChainlinkNode(input *Input) (*Node, error)`
  - `type Node` – Chainlink node instance
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode)

---

## WASP (Load Testing)

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/wasp`
- **Key Types:**
  - `type Profile` – Load test profile
  - `type Generator` – Load generator
  - `func NewProfile() *Profile`
  - `func NewGenerator(cfg *Config) *Generator`
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/wasp)

---

## Seth (Ethereum Client)

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/seth`
- **Key Types:**
  - `type Client` – Ethereum client
  - `func NewClient(url string) (*Client, error)`
  - `func NewClientWithConfig(cfg *Config) (*Client, error)`
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/seth)

---

## Havoc (Chaos Testing)

- **Package:** `github.com/smartcontractkit/chainlink-testing-framework/havoc`
- **Key Types:**
  - `type Client` – Chaos Mesh client
  - `type Chaos` – Chaos experiment
  - `func NewClient() (*Client, error)`
  - `func NewChaos(client *Client, exp interface{}) (*Chaos, error)`
- **Docs:** [GoDoc](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc)

---

## More
- For additional packages (e.g., postgres, s3provider, testreporters), see the [GoDoc index](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework) 