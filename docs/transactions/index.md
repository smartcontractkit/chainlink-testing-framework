---
layout: default
title: Transactions
nav_order: 5
has_children: false
---

# Transactions

In order to fund contracts and Chainlink nodes, you need to make a few transactions on the blockchain you're working with. By default, sending a transaction will block the test until the transaction is confirmed on the blockchain, then move on. If you want to change this behavior, check the section on [parallel transactions](#parallel-transactions).

## Funding Contracts

Most on-chain objects like contracts have some sort of `Fund(...)` method that can be used to send an amount of funds to that contract's address.

```go
// Funds the ocrContract from the default private key.
// Fund this contract with 1 ETH
err := ocrContract.Fund(big.NewFloat(1))
// Funds the same contract with .001 ETH
err = ocrContract.Fund(big.NewFloat(.001))
```

## Funding Chainlink Nodes

Chainlink nodes are easy to fund with the handy [actions package](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework/actions).

```go
// Funds each node in the `chainlinkNodes` array with .01 of the network's native currency
err := actions.FundChainlinkNodes(chainlinkNodes, defaultNetwork, big.NewFloat(.01))
```

## Parallel Transactions

Sometimes you'll want to fund a bunch of addresses or launch a few contracts that have no relation to each other. Rather than waiting for each transaction or contract deployment to be confirmed before moving on to the next one, you can make use of parallel transactions.

```go
// Any transactions made after this line will be submitted instantly to the blockchain without waiting for previous ones.
defaultNetwork.ParallelTransactions(true)
```

This can seriously speed up some test setups and executions, for example, if you want to fund your Chainlink nodes and a couple deployed contracts at the beginning of your test. But be wary, as some events depend on previous ones being confirmed, like funding a contract only after it has been deployed and confirmed on chain. For that, utilize the `defaultNetwork.WaitForEvents()` method call to halt running.

```go
// Set parallel transactions
defaultNetwork.ParallelTransactions(true)

// A pretend method to deploy a bunch of contracts and return them as a list
someListOfContracts := deploySomeContracts()
// Fund Chainlink nodes, which don't depend on contracts being confirmed on chain
err := actions.FundChainlinkNodes(chainlinkNodes, defaultNetwork, big.NewFloat(.01))

// Wait for all on-chain events started above to complete before moving on
err = defaultNetwork.WaitForEvents()

// Fund the deployed contracts
for _, contract := range someListOfContracts {
  contract.Fund(big.NewFloat(1))
}

// Wait for the contract funding to go through before moving on
err = defaultNetwork.WaitForEvents()
```

If you see errors in funding or interacting with contracts that imply the contracts don't exist, it's likely you set parallel transactions to `true` and failed to `WaitForEvents()` at an appropriate time. This feature is handy to speed up your test setups and executions, but can trip you up if not properly monitored, so take care.