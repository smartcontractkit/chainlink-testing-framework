# Components Resources

You can use `resources` to limit containers CPU/Memory for `NodeSet`, `Blockchain` and `PostgreSQL` components.

```toml
[blockchain_a.resources]
cpus = 0.5
memory_mb = 1048

[nodeset.db.resources]
cpus = 2
memory_mb = 2048

[nodeset.node_specs.node.resources]
cpus = 1
memory_mb = 1048
```

Read more about resource constraints [here](https://docs.docker.com/engine/containers/resource_constraints/).

We are using `cpu-period` and `cpu-quota` for simplicity, and because it's working with an arbitrary amount of containers, it is absolute.

How quota and period works:

- To allocate `1 CPU`, we set `CPUQuota = 100000` and `CPUPeriod = 100000` (1 full period).
- To allocate `0.5 CPU`, we set `CPUQuota = 50000` and `CPUPeriod = 100000`.
- To allocate `2 CPUs`, we set `CPUQuota = 200000` and `CPUPeriod = 100000`.

Read more about [CFS](https://engineering.squarespace.com/blog/2017/understanding-linux-container-scheduling).

When the `resources.memory_mb` key is not empty, we disable swap, ensuring the container goes OOM when memory is exhausted, allowing for more precise detection of sudden memory spikes.

Full configuration [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/smoke_limited_resources.toml)
