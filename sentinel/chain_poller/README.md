# **Chain Poller**

## **Overview**

The **Chain Poller** is a lightweight utility in the [Sentinel](https://github.com/smartcontractkit/sentinel) framework designed to fetch blockchain logs based on filter queries. It serves as a bridge between the blockchain client and higher-level services like the `ChainPollerService`.

---

## **Features**

- **Flexible Queries**: Fetch logs based on block ranges, addresses, and topics.
- **Error Logging**: Captures errors during polling for troubleshooting.

---

## **Usage**

### **Initialization**
Create a new `ChainPoller` with the required configuration:
```go
config := chain_poller.ChainPollerConfig{
    BlockchainClient: client,
    Logger:           logger,
    ChainID:          1,
}

chainPoller, err := chain_poller.NewChainPoller(config)
if err != nil {
    panic("Failed to initialize Chain Poller: " + err.Error())
}
```

### **Polling**
Fetch logs using filter queries:
```go
queries := []internal.FilterQuery{
    {
        FromBlock: 100,
        ToBlock:   200,
        Addresses: []common.Address{
            common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
        },
        Topics: [][]common.Hash{
            {common.HexToHash("0xabcdefabcdefabcdefabcdefabcdefabcdefabcdef")},
        },
    },
}

logs, err := chainPoller.Poll(context.Background(), queries)
if err != nil {
    logger.Error("Failed to fetch logs", "error", err)
}
fmt.Println("Fetched logs:", logs)
```

---

## **API Reference**

### **`NewChainPoller(config ChainPollerConfig) (*ChainPoller, error)`**
- Initializes a new Chain Poller with the specified blockchain client and logger.

### **`Poll(ctx context.Context, filterQueries []FilterQuery) ([]Log, error)`**
- Fetches logs from the blockchain based on the given filter queries.

---

## **Testing**

Run the tests:
```bash
go test ./chain_poller
```