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

Memory swapping is off if you specify `resources` key.

Full configuration [example]()