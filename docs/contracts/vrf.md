---
layout: default
title: Verifiable Random Function
nav_order: 5
parent: Contracts
---

# VRF

[Example Usage](./../suite/smoke/contracts_vrf_test.go)

```go
type VRF interface {
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
  ProofLength(context.Context) (*big.Int, error)
}

type VRFCoordinator interface {
  RegisterProvingKey(
    fromWallet client.BlockchainWallet,
    fee *big.Int,
    oracleAddr string,
    publicProvingKey [2]*big.Int,
    jobID [32]byte,
  ) error
  HashOfKey(ctx context.Context, pubKey [2]*big.Int) ([32]byte, error)
  Address() string
}

type VRFConsumer interface {
  Address() string
  RequestRandomness(fromWallet client.BlockchainWallet, hash [32]byte, fee *big.Int) error
  CurrentRoundID(ctx context.Context) (*big.Int, error)
  RandomnessOutput(ctx context.Context) (*big.Int, error)
  WatchPerfEvents(ctx context.Context, eventChan chan<- *PerfEvent) error
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
}
```

VRF enables you to generate provably random numbers and post them to the blockchain.

[Learn More](https://docs.chain.link/docs/chainlink-vrf/)