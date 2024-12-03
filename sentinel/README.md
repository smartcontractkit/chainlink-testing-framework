# Subscription Manager

## Overview

The **Subscription Manager** is a core component of the [Sentinel](https://github.com/smartcontractkit/sentinel) designed to facilitate efficient and organized management of blockchain event subscriptions. It enables seamless subscription and unsubscription to specific blockchain events, broadcasting of event logs to subscribed listeners, and ensures cache consistency through intelligent invalidation mechanisms.

> **Note:** The Subscription Manager is **not** intended to be used as a standalone component. It is tightly integrated within Sentinel to provide subscription tracking for testing applications.

## Features

- **Subscribe to Events:** Allows subscribing to specific addresses and topics, enabling listeners to receive relevant event logs.
- **Unsubscribe from Events:** Facilitates the removal of subscriptions, ensuring that listeners no longer receive event logs they are no longer interested in.
- **Broadcast Event Logs:** Efficiently broadcasts received event logs to all relevant subscribers, ensuring real-time updates.
- **Cache Management:** Maintains an internal cache of active subscriptions and intelligently invalidates the cache upon any subscription changes to ensure data consistency.
- **Thread-Safe Operations:** Utilizes mutexes to ensure that all subscription operations are thread-safe, preventing race conditions in concurrent environments.
- **Comprehensive Logging:** Integrates with a logging system to provide detailed insights into subscription changes, cache invalidations, and broadcasting activities.

## Integration within the Chainlink Testing Framework

The Subscription Manager is a pivotal component within the Sentinel module of the Chainlink Testing Framework. It works in tandem with other modules to simulate and monitor blockchain interactions, making it invaluable for testing decentralized applications.

### How It Works

1. **Subscription Registration:**
   - **Subscribe:** Listeners can subscribe to specific events by providing the contract address and event topic. The Subscription Manager registers these subscriptions and updates its internal cache accordingly.
   - **Unsubscribe:** Subscribers can be removed from specific event subscriptions, prompting the Subscription Manager to update its registry and invalidate caches as necessary.

2. **Event Broadcasting:**
   - When an event log is captured, the Subscription Manager identifies all subscribers interested in that event based on the address and topic.
   - It then broadcasts the log to all relevant subscribers, ensuring they receive timely updates.

3. **Cache Invalidation:**
   - Any change in subscriptions (addition or removal) triggers a cache invalidation to maintain consistency.
   - The Subscription Manager ensures that the cache accurately reflects the current state of active subscriptions.

### Usage within Tests

While the Subscription Manager isn't meant to be used directly in isolation, it plays a crucial role in facilitating tests that require monitoring and reacting to specific blockchain events.

```go
package subscription_manager_test

import (
    "testing"

    "github.com/ethereum/go-ethereum/common"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel/subscription_manager"
)

func TestExample(t *testing.T) {
    manager := subscription_manager.NewSubscriptionManager(/* parameters */)

    address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
    topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcdefabcd")

    // Subscribe to an event
    ch, err := manager.Subscribe(address, topic)
    require.NoError(t, err)
    assert.NotNil(t, ch)

    // Perform actions that trigger the event...

    // Listen for the event log
    log := <-ch
    assert.Equal(t, expectedLog, log)

    // Unsubscribe when done
    err = manager.Unsubscribe(address, topic, ch)
    require.NoError(t, err)
}
```

## API Reference

### `Subscribe(address common.Address, topic common.Hash) (chan Log, error)`

**Description:**  
Subscribes to a specific event based on the provided contract address and event topic.

**Parameters:**
- `address` (`common.Address`): The address of the contract emitting the event.
- `topic` (`common.Hash`): The specific event topic to subscribe to.

**Returns:**
- `chan Log`: A channel through which event logs will be received.
- `error`: An error object if the subscription fails.

### `Unsubscribe(address common.Address, topic common.Hash, ch chan Log) error`

**Description:**  
Unsubscribes from a specific event, removing the provided channel from the subscription list.

**Parameters:**
- `address` (`common.Address`): The address of the contract emitting the event.
- `topic` (`common.Hash`): The specific event topic to unsubscribe from.
- `ch` (`chan Log`): The channel to remove from the subscription list.

**Returns:**
- `error`: An error object if the unsubscription fails.

### `BroadcastLog(eventKey EventKey, log Log)`

**Description:**  
Broadcasts an event log to all subscribers associated with the given event key.

**Parameters:**
- `eventKey` (`EventKey`): A struct containing the address and topic.
- `log` (`Log`): The event log to broadcast.

### `GetAddressesAndTopics() map[common.Address][]common.Hash`

**Description:**  
Retrieves the current mapping of subscribed addresses and their corresponding topics.

**Returns:**
- `map[common.Address][]common.Hash`: A map where each key is an address and the value is a slice of topics subscribed to that address.

### `Close()`

**Description:**  
Closes the Subscription Manager, unsubscribing all listeners and cleaning up resources.

**Usage:**
```go
manager.Close()
```

## Testing

The Subscription Manager is rigorously tested to ensure reliability and correctness. Tests are designed to verify:

- **Subscription and Unsubscription:** Ensures that subscribing and unsubscribing operations correctly modify the internal registry.
- **Event Broadcasting:** Validates that event logs are accurately broadcasted to all relevant subscribers.
- **Cache Invalidation:** Confirms that the internal cache is properly invalidated upon any changes in subscriptions.
- **Thread Safety:** Utilizes Go's race detector to ensure that concurrent operations do not lead to race conditions.

### Running Tests

To execute the tests with the race detector enabled:

```bash
go test -race ./subscription_manager
```

### Mocking and Logging

A `MockLogger` is employed to capture and verify log messages generated by the Subscription Manager. This allows tests to assert that specific actions (like subscribing or unsubscribing) produce the expected log outputs, enhancing the robustness of the test suite.