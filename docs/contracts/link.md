---
layout: default
title: LINK Token
nav_order: 2
parent: Contracts
---

# LINK Token

```go
type LinkToken interface {
  Address() string
  Approve(fromWallet client.BlockchainWallet, to string, amount *big.Int) error
  Transfer(fromWallet client.BlockchainWallet, to string, amount *big.Int) error
  BalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
  TransferAndCall(fromWallet client.BlockchainWallet, to string, amount *big.Int, data []byte) error
  Fund(fromWallet client.BlockchainWallet, ethAmount *big.Float) error
  Name(context.Context) (string, error)
}
```

The LINK token used to power chainlink oracles. We deploy a new token on each test to simplify wallet management.
You can use the link token contract to transfer tokens to your chainlink nodes.

[Learn More](https://docs.chain.link/docs/architecture-request-model/#link-token)
