---
layout: default
title: Keeper
nav_order: 6
parent: Contracts
---

# Keeper

[Example Usage](./../suite/smoke/contracts_keeper_test.go)

```go
type UpkeepRegistrar interface {
  Address() string
  SetRegistrarConfig(
    fromWallet client.BlockchainWallet,
    autoRegister bool,
    windowSizeBlocks uint32,
    allowedPerWindow uint16,
    registryAddr string,
    minLinkJuels *big.Int,
  ) error
  EncodeRegisterRequest(
    name string,
    email []byte,
    upkeepAddr string,
    gasLimit uint32,
    adminAddr string,
    checkData []byte,
    amount *big.Int,
    source uint8,
  ) ([]byte, error)
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
}

type KeeperRegistry interface {
  Address() string
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
  SetRegistrar(fromWallet client.BlockchainWallet, registrarAddr string) error
  AddUpkeepFunds(fromWallet client.BlockchainWallet, id *big.Int, amount *big.Int) error
  GetUpkeepInfo(ctx context.Context, id *big.Int) (*UpkeepInfo, error)
  GetKeeperInfo(ctx context.Context, keeperAddr string) (*KeeperInfo, error)
  SetKeepers(fromWallet client.BlockchainWallet, keepers []string, payees []string) error
  GetKeeperList(ctx context.Context) ([]string, error)
  RegisterUpkeep(fromWallet client.BlockchainWallet, target string, gasLimit uint32, admin string, checkData []byte) error
}

type KeeperConsumer interface {
  Address() string
  Fund(fromWallet client.BlockchainWallet, ethAmount, linkAmount *big.Float) error
  Counter(ctx context.Context) (*big.Int, error)
}
```

[Learn More](https://docs.chain.link/docs/chainlink-keepers/introduction/)