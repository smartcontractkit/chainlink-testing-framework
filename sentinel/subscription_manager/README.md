# **Subscription Manager**

## **Overview**

The **Subscription Manager** is a utility in the [Sentinel](https://github.com/smartcontractkit/sentinel) module of the Chainlink Testing Framework. It efficiently manages blockchain event subscriptions, broadcasts event logs, and maintains an optimized cache of active subscriptions.

---

## **Features**

- **Dynamic Subscriptions**: Subscribe and unsubscribe from blockchain events.
- **Log Broadcasting**: Relay logs to relevant subscribers.
- **Cache Optimization**: Maintains a mapping of active subscriptions for efficient filtering.
- **Thread Safety**: Ensures concurrency safety with mutex locks.
- **Integration**: Seamlessly works with `ChainPollerService` and `Sentinel`.

---

## **How It Works**

1. **Subscribe**: 
   - Add a subscription by providing an address and topic.
   - Returns a channel for receiving logs.

2. **Unsubscribe**:
   - Remove a subscription and safely close the associated channel.

3. **Broadcast Logs**:
   - Sends logs to all subscribers of a specific address and topic.

4. **Cache Management**:
   - Maintains a cache for quick access to active subscriptions.
   - Automatically invalidates the cache when subscriptions change.

---

## **API Reference**

### **`Subscribe(address common.Address, topic common.Hash) (chan Log, error)`**
Adds a subscription for a blockchain address and topic. Returns a channel for logs.

### **`Unsubscribe(address common.Address, topic common.Hash, ch chan Log) error`**
Removes a subscription and closes the associated channel.

### **`BroadcastLog(eventKey EventKey, log Log)`**
Broadcasts a log to all subscribers for a specific event.

### **`GetAddressesAndTopics() map[common.Address][]common.Hash`**
Retrieves the current mapping of addresses and topics in the subscription cache.

### **`Close()`**
Gracefully shuts down the manager, unsubscribes all listeners, and clears the registry.

---

## **Usage Example**

```go
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/sentinel/subscription_manager"
)

func main() {
	manager := subscription_manager.NewSubscriptionManager(/* logger, chainID */)

	address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

	// Subscribe to an event
	ch, err := manager.Subscribe(address, topic)
	if err != nil {
		fmt.Println("Subscription failed:", err)
		return
	}

	// Simulate log broadcast
	go func() {
		log := internal.Log{Address: address, Topics: []common.Hash{topic}, Data: []byte("event data")}
		manager.BroadcastLog(internal.EventKey{Address: address, Topic: topic}, log)
	}()

	// Receive logs
	go func() {
		for log := range ch {
			fmt.Println("Received log:", log)
		}
	}()

	// Unsubscribe
	if err := manager.Unsubscribe(address, topic, ch); err != nil {
		fmt.Println("Unsubscribe failed:", err)
	}
}
```

---

## **Testing**

### **Run Tests**
Execute unit tests with race detection:
```bash
go test -race ./subscription_manager
```

### **Test Coverage**
- Thread-safe operations.
- Subscription and unsubscription behavior.
- Log broadcasting accuracy.
- Cache invalidation and consistency.
