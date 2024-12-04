# **Sentinel**

## **Overview**

The **Sentinel** is a modular and extensible system for monitoring blockchain events and managing event subscriptions. Sentinel orchestrates multiple `ChainPollerService` instances, each tied to a specific blockchain network, to fetch logs, process events, and notify subscribers.

---

## **Key Features**

- **Multi-Chain Support**: Manage multiple blockchain networks (e.g., Ethereum, Polygon, Arbitrum) concurrently.
- **Event Broadcasting**: Relay blockchain events to subscribers via a thread-safe subscription system.
- **Flexible Subscriptions**: Dynamically subscribe and unsubscribe to events based on addresses and topics.
- **Graceful Lifecycle Management**: Start, stop, and clean up resources across services effortlessly.

---

## **System Architecture**

### **How Components Interact**

```
                              +----------------+
                              |    Sentinel    |
                              | (Coordinator)  |
                              +----------------+
                                       |
     +---------------------------------+--------------------------------+
     |                                 |                                |
+------------+               +-------------------+           +-------------------+
| Chain ID 1 |               |    Chain ID 2     |           |    Chain ID 3     |
| (Ethereum) |               |   (Polygon)       |           |   (Arbitrum)      |
+------------+               +-------------------+           +-------------------+
       |                               |                                |
+----------------+           +-------------------+           +-------------------+
| ChainPollerSvc |           | ChainPollerSvc    |           | ChainPollerSvc    |
|  (Service 1)   |           |   (Service 2)     |           |   (Service 3)     |
+----------------+           +-------------------+           +-------------------+
       |                               |                                |
+----------------+           +-------------------+           +-------------------+
|   ChainPoller  |           |   ChainPoller     |           |   ChainPoller     |
| (Log Fetching) |           | (Log Fetching)    |           | (Log Fetching)    |
+----------------+           +-------------------+           +-------------------+
       |                               |                                |
+----------------+           +-------------------+           +-------------------+
| Subscription   |           | Subscription      |           | Subscription      |
|   Manager      |           |   Manager         |           |   Manager         |
+----------------+           +-------------------+           +-------------------+
       |                               |                                |
+----------------+           +-------------------+           +-------------------+
| Blockchain     |           | Blockchain        |           | Blockchain        |
| (Ethereum)     |           | (Polygon)         |           | (Arbitrum)        |
+----------------+           +-------------------+           +-------------------+

```

### **Core Components**
1. **Sentinel**:
   - Central coordinator managing multiple `ChainPollerService` instances.
   - Handles multi-chain subscriptions, lifecycle management, and configuration.

2. **ChainPollerService**:
   - Polls blockchain logs and broadcasts events to subscribers.
   - Integrates `ChainPoller` and `SubscriptionManager`.

3. **ChainPoller**:
   - Fetches logs from blockchain networks based on filter queries.

4. **SubscriptionManager**:
   - Tracks subscriptions to blockchain events.
   - Broadcasts logs to subscribers.

---

## **Usage**

### **Initialize Sentinel**
Set up a `Sentinel` instance:
```go
import (
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel"
    "github.com/smartcontractkit/chainlink-testing-framework/sentinel/chain_poller_service"
)

logger := internal.NewDefaultLogger()
sentinelInstance := sentinel.NewSentinel(sentinel.SentinelConfig{
    Logger: logger,
})
```

### **Add a Chain**
Add a blockchain to monitor:
```go
config := chain_poller_service.ChainPollerServiceConfig{
    PollInterval:     100 * time.Millisecond,
    ChainPoller:      chainPoller,
    Logger:           logger,
    BlockchainClient: client,
    ChainID:          1,
}

err := sentinelInstance.AddChain(config)
if err != nil {
    panic("Failed to add chain: " + err.Error())
}
```

### **Subscribe to Events**
Subscribe to blockchain events:
```go
address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

logCh, err := sentinelInstance.Subscribe(1, address, topic)
if err != nil {
    panic("Failed to subscribe: " + err.Error())
}

// Listen for logs
go func() {
    for log := range logCh {
        fmt.Println("Received log:", log)
    }
}()
```

### **Unsubscribe**
Unsubscribe from events:
```go
err = sentinelInstance.Unsubscribe(1, address, topic, logCh)
if err != nil {
    panic("Failed to unsubscribe: " + err.Error())
}
```

### **Remove a Chain**
Remove a blockchain from monitoring:
```go
err = sentinelInstance.RemoveChain(1)
if err != nil {
    panic("Failed to remove chain: " + err.Error())
}
```

---

## **Testing**

### **Run Tests**
Run tests for Sentinel and its components:
```bash
go test ./sentinel
```