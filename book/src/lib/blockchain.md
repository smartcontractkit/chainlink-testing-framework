# Blockchain Clients

This folder contains the bulk of code that handles integrating with different EVM chains. If you're looking to run tests on a new EVM chain, and are having issues with the default implementation, you've come to the right place.

### Some Terminology

- [L2 Chain](https://ethereum.org/en/layer-2/): A Layer 2 chain "branching" off Ethereum.
- [EVM](https://ethereum.org/en/developers/docs/evm/): Ethereum Virtual Machine that underpins the Ethereum blockchain.
- [EVM Compatible](https://blog.thirdweb.com/evm-compatible-blockchains-and-ethereum-virtual-machine/#:~:text=What%20does%20'EVM%20compatibility'%20mean,significant%20changes%20to%20their%20code.): A chain that has some large, underlying differences from how base Ethereum works, but can still be interacted with largely the same way as Ethereum.
- [EIP-1559](https://eips.ethereum.org/EIPS/eip-1559): The Ethereum Improvement Proposal that changed how gas fees are calculated and paid on Ethereum
- Legacy Transactions: Transactions that are sent using the old gas fee calculation method, the one used before EIP-1559.
- Dynamic Fee Transaction: Transactions that are sent using the new gas fee calculation method, the one used after EIP-1559.

## How Client Integrations Work

In order to test Chainlink nodes, the `chainlink-testing-framework` needs to be able to interact with the chain that the node is running on. This is done through the `blockchain.EVMClient` interface. The `EVMClient` interface is a wrapper around [geth](https://geth.ethereum.org/) to interact with the blockchain. We conduct all our testing blockchain operations through this wrapper, like sending transactions and monitoring on-chain events. The primary implementation of this wrapper is built for [Ethereum](./ethereum.go). Most others, like the [Metis](./metis.go) and [Optimism](./optimism.go) integrations, extend and modify the base Ethereum implementation.

## Do I Need a New Integration?

If you're reading this, probably. The default EVM integration is designed to work with mainnet Ethereum, which covers most other EVM chain interactions, but it's not guaranteed to work with all of them. If you're on a new chain and the test framework is throwing errors while doing basic things like send transactions, receive new headers, or deploy contracts, you'll likely need to create a new integration. The most common issue with new chains (especially L2s) is gas estimations and lack of support for dynamic transactions.

## Creating a New Integration

Take a look at the [Metis](./metis.go) integration as an example. Metis is an L2, EVM compatible chain. It's largely the same as the base Ethereum integration, so we'll extend from that.

```go
type MetisMultinodeClient struct {
  *EthereumMultinodeClient
}

type MetisClient struct {
  *EthereumClient
}
```

Now we need to let other libraries (like our tests in the main Chainlink repo) that this integration exists. So we add the new implementation to the [known_networks.go](./known_networks.go) file. We can then add that network to our tests' own [known_networks.go](https://github.com/smartcontractkit/chainlink/blob/develop/integration-tests/known_networks.go) file (it's annoying, there are plans to simplify).

Now our Metis integration is the exact same as our base Ethereum one, which doesn't do us too much good. I'm assuming you came here to make some changes, so first let's find out what we need to change. This is a mix of reading developer documentation on the chain you're testing and trial and error. Mostly the latter in later stages. In the case of Metis, like many L2s, they [have their own spin on gas fees](https://docs.metis.io/dev/protocol-in-detail/transaction-fees-on-the-metis-platform). They also only support Legacy transactions. So we'll need to override any methods that deal with gas estimations, `Fund`, `DeployContract`, and `ReturnFunds`.
