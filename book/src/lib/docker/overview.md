# Docker

Docker package represents low-level wrappers for Docker containers built using [testcontainers-go](https://golang.testcontainers.org/) library.
It is strongly advised that you use CTFv2's [Framework](../framework/overview.md) to build your enviroment instead of using it directly. **Consider yourself warned!**

Supported Docker containers can be divided in a couple of categories:
* [Chainlink ecosystem-related](./chainlink_ecosystem.md)
* [Blockchain nodes](./blockchain_nodes.md)
* Third party apps
* [Test helpers](./test_helpers.md)
* helper containers

Following Chainlink-related containers are available:
* Job Distributor

Blockchain nodes:
* Besu
* Erigon
* Geth
* Nethermind
* Reth
* Prysm Beacon client (PoS Ethereum)

Third party apps:
* Kafka
* Postgres
* Zookeeper
* Schema registry

Test helpers (mocking solutions):
* Killgrave
* Mockserver

Helper containers:
* Ethereum genesis generator
* Validator keys generator

# Basic structure

All of our Docker container wrappers are composed of some amount of specific elements and an universal part called `EnvComponent` defined as:
```go
type EnvComponent struct {
	ContainerName      string             `json:"containerName"`
	ContainerImage     string             `json:"containerImage"`
	ContainerVersion   string             `json:"containerVersion"`
	ContainerEnvs      map[string]string  `json:"containerEnvs"`
	WasRecreated       bool               `json:"wasRecreated"`
	Networks           []string           `json:"networks"`
	Container          tc.Container       `json:"-"`
	PostStartsHooks    []tc.ContainerHook `json:"-"`
	PostStopsHooks     []tc.ContainerHook `json:"-"`
	PreTerminatesHooks []tc.ContainerHook `json:"-"`
	LogLevel           string             `json:"-"`
	StartupTimeout     time.Duration      `json:"-"`
}
```

It comes with a bunch of functional options that allow you to:
* set container name
* set image name and tag
* set startup timeout
* set log level
* set post start/stop hooks

And following chaos-testing related functions:
* `ChaosNetworkDelay()` - introducing networking delay
* `ChaosNetworkLoss()` - simulating network package loss
* `ChaosPause()` - pausing the container