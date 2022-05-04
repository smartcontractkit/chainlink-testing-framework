---
layout: default
title: Chainlink Nodes
nav_order: 4
has_children: false
---

# Chainlink Nodes

Make sure to have your environment setup code in place, which will deploy a few Chainlink nodes, and give you an `env` environment.

```go
// Get a list of all Chainlink nodes deployed in your test environment
chainlinkNodes, err := client.ConnectChainlinkNodes(env)
```

From here, you can [interact with](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/client#Chainlink) each Chainlink node to manage keys, jobs, bridges, and transactions.

## Chainlink Jobs

The most common interaction you'll have with Chainlink nodes will likely be creating jobs, using the `chainlinkNode.CreateJob(JobSpec)` method. Chainlink jobs are how the Chainlink nodes know what actions they're expected to perform on chain, and how thy should perform them. A typical test consists of launching your resources, deploying contracts to the blockchain, and telling the Chainlink node to interact with those contracts by creating a job. Read more about Chainlink jobs and the specifics on using them [here](https://docs.chain.link/docs/jobs/).

There are plenty of built in [JobSpecs](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/client#JobSpec) like the [Keeper Job Spec](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/client#KeeperJobSpec) and the [OCR Job Spec](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/client#OCRTaskJobSpec) that you can use. But if for whatever reason, those don't do the job for you, you can create a raw TOML job with `CreateJobRaw(string)` like below.

```go
jobData, err := chainlinkNode.CreateJobRaw(`
schemaVersion = 1
otherField    = true
`)
```
