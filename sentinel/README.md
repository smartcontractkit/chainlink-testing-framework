# Sentinel

Sentinel is a robust, modular, and extensible Go-based framework designed for monitoring blockchain events and managing event subscriptions. It orchestrates multiple `ChainPollerService` instances, each tied to a specific blockchain network, to fetch logs and notify relevant subscribers seamlessly.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [System Architecture](#system-architecture)
  - [How Components Interact](#how-components-interact)
  - [Core Components](#core-components)
- [Usage](#usage)
  - [Initialize Sentinel](#initialize-sentinel)
  - [Add a Chain](#add-a-chain)
  - [Subscribe to Events](#subscribe-to-events)
  - [Unsubscribe](#unsubscribe)
  - [Remove a Chain](#remove-a-chain)
- [API Reference](#api-reference)
  - [Sentinel](#sentinel)
  - [ChainPollerService](#chainpollerservice)
  - [SubscriptionManager](#subscriptionmanager)
- [Testing](#testing)
  - [Run Tests](#run-tests)
- [Contributing](#contributing)
- [License](#license)

## Overview

Sentinel is a centralized orchestrator that manages multiple blockchain poller services, each responsible for a specific blockchain network (e.g., Ethereum, Polygon, Arbitrum). It provides a unified interface for subscribing to blockchain events, ensuring efficient log polling and event broadcasting to subscribers.


## Key Features

- **Multi-Chain Support**: Manage multiple blockchain networks concurrently.
- **Event Broadcasting**: Relay blockchain events to subscribers via a thread-safe subscription system.
- **Flexible Subscriptions**: Dynamically subscribe and unsubscribe to events based on addresses and topics.
- **Graceful Lifecycle Management**: Start, stop, and clean up resources across services effortlessly.
- **Comprehensive Testing**: Ensures reliability through extensive unit and integration tests.
- **Scalable Architecture**: Designed to handle multiple chains and high-frequency event broadcasting.

## System Architecture

### How Components Interact

## System Architecture

### How Components Interact

```mermaid
graph TD
    Sentinel["Sentinel<br/>(Coordinator)"]

    subgraph Ethereum
        ChainPollerSvc_Ethereum["ChainPollerSvc<br/>(Ethereum)"]
        ChainPoller_Ethereum["ChainPoller<br/>(Log Fetching)"]
        SubscriptionManager_Ethereum["Subscription Manager"]
        ChainPollerSvc_Ethereum --> ChainPoller_Ethereum
        ChainPollerSvc_Ethereum --> SubscriptionManager_Ethereum
        ChainPoller_Ethereum --> Blockchain_Ethereum["Blockchain<br/>(Ethereum)"]
    end

    subgraph Polygon
        ChainPollerSvc_Polygon["ChainPollerSvc<br/>(Polygon)"]
        ChainPoller_Polygon["ChainPoller<br/>(Log Fetching)"]
        SubscriptionManager_Polygon["Subscription Manager"]
        ChainPollerSvc_Polygon --> ChainPoller_Polygon
        ChainPollerSvc_Polygon --> SubscriptionManager_Polygon
        ChainPoller_Polygon --> Blockchain_Polygon["Blockchain<br/>(Polygon)"]
    end

    subgraph Arbitrum
        ChainPollerSvc_Arbitrum["ChainPollerSvc<br/>(Arbitrum)"]
        ChainPoller_Arbitrum["ChainPoller<br/>(Log Fetching)"]
        SubscriptionManager_Arbitrum["Subscription Manager"]
        ChainPollerSvc_Arbitrum --> ChainPoller_Arbitrum
        ChainPollerSvc_Arbitrum --> SubscriptionManager_Arbitrum
        ChainPoller_Arbitrum --> Blockchain_Arbitrum["Blockchain<br/>(Arbitrum)"]
    end

    Sentinel --> Ethereum
    Sentinel --> Polygon
    Sentinel --> Arbitrum
```

### Core Components

1. **Sentinel**:
   - **Role**: Central coordinator managing multiple `ChainPollerService` instances.
   - **Responsibilities**:
     - Handles adding and removing blockchain chains.
     - Manages global subscriptions.
     - Orchestrates communication between components.

2. **ChainPollerService**:
   - **Role**: Manages the polling process for a specific blockchain.
   - **Responsibilities**:
     - Polls blockchain logs based on filter queries.
     - Integrates `ChainPoller` and `SubscriptionManager`.
     - Broadcasts fetched logs to relevant subscribers.

3. **ChainPoller**:
   - **Role**: Fetches logs from blockchain networks.
   - **Responsibilities**:
     - Interacts with the blockchain client to retrieve logs.
     - Processes filter queries to fetch relevant logs.

4. **SubscriptionManager**:
   - **Role**: Manages event subscriptions for a specific chain.
   - **Responsibilities**:
     - Tracks subscriptions to blockchain events.
     - Ensures thread-safe management of subscribers.
     - Broadcasts logs to all relevant subscribers.

## Usage

### Initialize Sentinel

Set up a `Sentinel` instance:

```go
package main

import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/rs/zerolog"
    "os"

    "github.com/smartcontractkit/chainlink-testing-framework/sentinel"
)

func main() {
    // Initialize logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Initialize Sentinel
    sentinelCoordinator := sentinel.NewSentinel(sentinel.SentinelConfig{
        Logger: logger,
    })
    defer sentinelCoordinator.Close()
}
```

### Add a Chain

Add a blockchain to monitor:

```go
package main

import (
    "time"

    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel/blockchain_client_wrapper"
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel/sentinel"
)

func main() {
    // Initialize logger and Sentinel as shown above

    // Setup blockchain client (e.g., Geth)
    client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
    if err != nil {
        panic("Failed to connect to blockchain client: " + err.Error())
    }
    wrappedClient := blockchain_client_wrapper.NewGethClientWrapper(client)

    // Add a new chain to Sentinel
    err = sentinelCoordinator.AddChain(sentinel.AddChainConfig{
        ChainID:          1, // Ethereum Mainnet
        PollInterval:     10 * time.Second,
        BlockchainClient: wrappedClient,
    })
    if err != nil {
        panic("Failed to add chain: " + err.Error())
    }
}
```

### Subscribe to Events

Subscribe to blockchain events:

```go
package main

import (
    "fmt"
    "github.com/ethereum/go-ethereum/common"
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel/api"
)

func main() {
    // Initialize logger, Sentinel, and add a chain as shown above

    // Define the address and topic to subscribe to
    address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
    topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

    // Subscribe to the event
    logCh, err := sentinelCoordinator.Subscribe(1, address, topic)
    if err != nil {
        panic("Failed to subscribe: " + err.Error())
    }
    defer sentinelCoordinator.Unsubscribe(1, address, topic, logCh)

    // Listen for logs in a separate goroutine
    go func() {
        for log := range logCh {
            fmt.Printf("Received log: %+v\n", log)
        }
    }()
}
```

### Unsubscribe

Unsubscribe from events:

```go
package main

func main() {
    // Initialize logger, Sentinel, add a chain, and subscribe as shown above

    // Assume logCh is the channel obtained from Subscribe
    err = sentinelCoordinator.Unsubscribe(1, address, topic, logCh)
    if err != nil {
        panic("Failed to unsubscribe: " + err.Error())
    }
}
```

### Remove a Chain

Remove a blockchain from monitoring:

```go
package main

func main() {
    // Initialize logger, Sentinel, add a chain, and subscribe as shown above

    // Remove the chain
    err = sentinelCoordinator.RemoveChain(1)
    if err != nil {
        panic("Failed to remove chain: " + err.Error())
    }
}
```

## API Reference

### Sentinel

- **`NewSentinel(config SentinelConfig) *Sentinel`**  
  Initializes a new Sentinel instance.

- **`AddChain(config AddChainConfig) error`**  
  Adds a new blockchain chain to Sentinel.

- **`RemoveChain(chainID int64) error`**  
  Removes an existing chain from Sentinel.

- **`Subscribe(chainID int64, address common.Address, topic common.Hash) (chan api.Log, error)`**  
  Subscribes to a specific event on a given chain.

- **`Unsubscribe(chainID int64, address common.Address, topic common.Hash, ch chan api.Log) error`**  
  Unsubscribes from a specific event.

- **`GetService(chainID int64) (*chain_poller_service.ChainPollerService, bool)`**  
  Retrieves the ChainPollerService for a given chain ID.

- **`HasServices() bool`**  
  Checks if there are any active services.

### ChainPollerService

- **`NewChainPollerService(config ChainPollerServiceConfig) (*ChainPollerService, error)`**  
  Initializes a new ChainPollerService.

- **`Start()`**  
  Starts the polling loop.

- **`Stop()`**  
  Stops the polling loop gracefully.

- **`SubscriptionManager() *subscription_manager.SubscriptionManager`**  
  Retrieves the SubscriptionManager.

### SubscriptionManager

- **`Subscribe(address common.Address, topic common.Hash) (chan api.Log, error)`**  
  Registers a new subscription.

- **`Unsubscribe(address common.Address, topic common.Hash, ch chan api.Log) error`**  
  Removes an existing subscription.

- **`BroadcastLog(eventKey internal.EventKey, log api.Log)`**  
  Broadcasts a log to all relevant subscribers.

- **`GetAddressesAndTopics() []internal.EventKey`**  
  Retrieves all unique EventKeys.

## Testing

### Run Tests

Run the comprehensive test suite using:

```bash
go test -race ./... -v
```