# **Chain Poller Service**

## **Overview**

The **Chain Poller Service** is a higher-level abstraction in the [Sentinel](https://github.com/smartcontractkit/lib/sentinel) framework that manages periodic blockchain log polling and event broadcasting. It integrates with the `ChainPoller` and `SubscriptionManager` to provide a complete solution for monitoring blockchain activity and notifying subscribers in real time.

---

## **Features**

- **Automated Polling**: Periodically fetches logs based on active subscriptions.
- **Real-Time Broadcasting**: Sends logs to subscribers through the `SubscriptionManager`.
- **Multi-Chain Support**: Handles subscriptions and polling for multiple blockchain networks.
- **Graceful Start/Stop**: Ensures clean resource management when starting or stopping the service.

---

## **Usage**

### **Initialization**
Create a `ChainPollerService` with the necessary configuration:
```go
config := chain_poller_service.ChainPollerServiceConfig{
    PollInterval:     100 * time.Millisecond,
    ChainPoller:      chainPoller,       // Instance of ChainPoller
    Logger:           logger,           // Logger instance
    BlockchainClient: client,           // Blockchain client
    ChainID:          1,                // Chain ID
}

pollerService, err := chain_poller_service.NewChainPollerService(config)
if err != nil {
    panic("Failed to initialize Chain Poller Service: " + err.Error())
}
```

### **Start and Stop**
Start and stop the polling service:
```go
pollerService.Start()
defer pollerService.Stop()
```

### **Subscriptions**
Manage subscriptions through the integrated `SubscriptionManager`:
```go
address := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
topic := common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")

logCh, err := pollerService.SubscriptionMgr.Subscribe(address, topic)
if err != nil {
    panic("Failed to subscribe: " + err.Error())
}

// Handle incoming logs
go func() {
    for log := range logCh {
        fmt.Println("Received log:", log)
    }
}()

// Unsubscribe when done
err = pollerService.SubscriptionMgr.Unsubscribe(address, topic, logCh)
if err != nil {
    panic("Failed to unsubscribe: " + err.Error())
}
```

---

## **API Reference**

### **`NewChainPollerService(config ChainPollerServiceConfig) (*ChainPollerService, error)`**
- Initializes a new Chain Poller Service with the specified configuration.

### **`Start()`**
- Begins the periodic polling process.

### **`Stop()`**
- Gracefully stops the polling process and releases resources.

### **`SubscriptionMgr`**
- Access the `SubscriptionManager` to manage subscriptions.

---

## **Testing**

### **Test Coverage**
The `ChainPollerService` includes tests for:
1. **Initialization**: Ensures valid configurations are required.
2. **Polling**: Verifies that logs are fetched and broadcasted correctly.
3. **Lifecycle Management**: Confirms proper start/stop behavior.
4. **Subscriptions**: Tests integration with the `SubscriptionManager`.

### **How to Run Tests**
To execute tests:
```bash
go test ./chain_poller_service
```