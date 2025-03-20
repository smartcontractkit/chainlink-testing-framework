# Chaos Testing

We offer Docker and Kubernetes boilerplates designed to test the resilience of `NodeSet` and `Blockchain`, which you can customize and integrate into your pipeline.


## Goals

We recommend structuring your tests as a linear suite that applies various chaos experiments and verifies the outcomes using a load testing suite. Focus on critical user metrics, such as:

- The ratio of successful responses to failed responses
- The nth percentile of response latency

Next, evaluate observability:

- Ensure proper alerts are triggered during failures (manual or automated)
- Verify the service recovers within the expected timeframe (manual or automated)

In summary, the **primary** focus is on meeting user expectations and maintaining SLAs, while the **secondary** focus is on observability and making operational part smoother.


## Docker

For Docker, we utilize [Pumba](https://github.com/alexei-led/pumba) to conduct chaos experiments, including:

- Container reboots
- Network simulations (such as delays, packet loss, corruption, etc., using the tc tool)
- Stress testing for CPU and memory usage

Additionally, we offer a [resources](../../framework/components/resources.md) API that allows you to test whether your software can operate effectively in low-resource environments.

You can also use [fake](../../framework/components/mocking.md) package to create HTTP chaos experiments.

Given the complexity of `Kubernetes`, we recommend starting with `Docker` first. Identifying faulty behavior in your services early—such as cascading latency—can prevent more severe issues when scaling up. Addressing these problems at a smaller scale can save significant time and effort later.

Check `NodeSet` + `Blockchain` template [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/chaos/chaos_docker_test.go).

## Kubernetes

We utilize a subset of [ChaosMesh](https://chaos-mesh.org/) experiments that can be safely executed on an isolated node group. These include:

- [Pod faults](https://chaos-mesh.org/docs/simulate-pod-chaos-on-kubernetes/)

- [Network faults](https://chaos-mesh.org/docs/simulate-network-chaos-on-kubernetes/) – We focus on delay and partition experiments, as others may impact pods outside the dedicated node group.

- [HTTP faults](https://chaos-mesh.org/docs/simulate-http-chaos-on-kubernetes/)

Check `NodeSet` + `Blockchain` template [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/chaos/chaos_k8s_test.go).

## Blockchain

We also offer a set of blockchain-specific experiments, which typically involve API calls to blockchain simulators to execute certain actions. These include:

- Adjusting gas prices

- Introducing chain reorganizations (setting a new head)

- Utilizing developer APIs (e.g., Anvil)

Check [gas](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/chaos/chaos_blockchain_evm_gas_test.go) and [reorg](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/chaos/chaos_blockchain_evm_reorg_test.go) examples, the same example work for [K8s](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/chaos/chaos_k8s_test.go).

## Debugging

To debug `Docker` applications you can just use `CTFv2` deployments.

To debug `K8s` please use our [simulator](../chaos/debug-k8s.md).
