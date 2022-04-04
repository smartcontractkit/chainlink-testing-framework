---
layout: default
title: Contracts
nav_order: 5
has_children: false
---

# Smart Contracts

Tests require deploying at least a few contracts to whatever chain you're running on. You can find the code used to deploy and interact with these contracts in the [contracts](https://github.com/smartcontractkit/integrations-framework/tree/main/contracts) package. Most of these are self-explanatory, but Keeper and OCR can be a bit more complicated.

## Contract Deployer

Each network has a contract deployer

```go
// See the previous setup code on how to get your default network
// The contract deployer will use the first private key listed to deploy contracts from
contractDeployer, err := contracts.NewContractDeployer(defaultNetwork)
```

From here, you can use the `contractDeployer` on its own, as defined [here](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework/contracts#ContractDeployer). For a lot of one off contracts, like the Link Token contract, that's all you need.

```go
contractDeployer.DeployLinkTokenContract()
```

From there, contract interactions are defined in the [contracts package](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework/contracts#pkg-overview).
