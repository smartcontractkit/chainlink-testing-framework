---
layout: default
title: Off-Chain Reporting
nav_order: 4
parent: Contracts
---

# OCR

[Example Usage](./../suite/smoke/contracts_ocr_test.go)

```go
type OffchainAggregator interface {
  Address() string
  Fund(fromWallet client.BlockchainWallet, nativeAmount, linkAmount *big.Float) error
  GetContractData(ctxt context.Context) (*OffchainAggregatorData, error)
  SetConfig(
    fromWallet client.BlockchainWallet,
    chainlinkNodes []client.Chainlink,
    ocrConfig OffChainAggregatorConfig,
  ) error
  SetPayees(client.BlockchainWallet, []common.Address, []common.Address) error
  RequestNewRound(fromWallet client.BlockchainWallet) error
  GetLatestAnswer(ctxt context.Context) (*big.Int, error)
  GetLatestRound(ctxt context.Context) (*RoundData, error)
}
```

OCR (Off-Chain Reporting) enables chainlink oracles to securely come to consensus on an answer and post it on chain,
greatly lowering the expense of posting multiple answers to the blockchain.

[Learn More](https://docs.chain.link/docs/off-chain-reporting/)