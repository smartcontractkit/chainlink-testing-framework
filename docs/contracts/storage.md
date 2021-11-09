---
layout: default
title: Storage
nav_order: 1
parent: Contracts
---

# Storage

[Example Usage](./../suite/smoke/contracts_test.go)

```go
type Storage interface {
  Get(ctxt context.Context) (*big.Int, error)
  Set(*big.Int) error
}
```

Super basic, we use it for minimum viability tests. Not much more to say here.