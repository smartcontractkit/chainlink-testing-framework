---
layout: default
title: Mock Adapter
nav_order: 6
has_children: false
---

# Mock Adapter

Most Chainlink jobs involve reaching out to outside APIs (referred to as "adapters") to gather data before making on-chain transactions. So along with a simulated blockchain, we also launch a mock adapter (sometimes referred to as a "mock server") that you can set and retrieve values for. We use an implementation of the [mock-server](https://www.mock-server.com/) project, but have simplified our interaction with it through just a few methods.

Connect to your mock adapter much like you do with your Chainlink nodes.

```go
// Get a connection to the mock server
mockserver, err = client.ConnectMockServer(env)

path := "/test_resp"
// Any GET calls to the path will return a Chainlink-node-readable response of 5
err := mockserver.SetValuePath(path, 5)

// Test some things

// Change the response of path to 15
err := mockserver.SetValuePath(path, 15)
```
