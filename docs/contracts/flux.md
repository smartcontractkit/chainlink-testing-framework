---
layout: default
title: Flux Aggregator
nav_order: 3
parent: Contracts
---

# Flux Aggregator

[Example Usage](../../suite/smoke/contracts_flux_test.go)

```go
type FluxAggregator interface {
  Address() string
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
  LatestRoundID(ctx context.Context, blockNumber *big.Int) (*big.Int, error)
  LatestRoundData(ctx context.Context) (RoundData, error)
  GetContractData(ctxt context.Context) (*FluxAggregatorData, error)
  UpdateAvailableFunds(ctx context.Context, fromWallet client.BlockchainWallet) error
  PaymentAmount(ctx context.Context) (*big.Int, error)
  RequestNewRound(ctx context.Context, fromWallet client.BlockchainWallet) error
  WithdrawPayment(ctx context.Context, caller client.BlockchainWallet, from common.Address, to common.Address, amount *big.Int) error
  WithdrawablePayment(ctx context.Context, addr common.Address) (*big.Int, error)
  GetOracles(ctx context.Context) ([]string, error)
  SetOracles(client.BlockchainWallet, FluxAggregatorSetOraclesOptions) error
  Description(ctxt context.Context) (string, error)
  SetRequesterPermissions(ctx context.Context, fromWallet client.BlockchainWallet, addr common.Address, authorized bool, roundsDelay uint32) error
  WatchSubmissionReceived(ctx context.Context, eventChan chan<- *SubmissionEvent) error
}
```

The original data aggregation strategy for chainlink oracles, this setup is being phased out a bit in favor of its more
gas-efficient cousin, OCR.

[Learn More](https://docs.chain.link/docs/architecture-decentralized-model/)