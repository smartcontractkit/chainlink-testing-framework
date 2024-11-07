# Components Persistence

We use static port ranges and volumes for all components to simplify Docker port management for developers.

This approach allows us to apply chaos testing to any container, ensuring it reconnects and retains the data needed for your tests.

When deploying a component, you can explicitly configure port ranges if the default ports don’t meet your needs.

Defaults are:
- [NodeSet](../components/chainlink/nodeset.md) (Node HTTP API): `10000..100XX`
- [NodeSet](../components/chainlink/nodeset.md) (Node P2P API): `12000..120XX`
```
[nodeset]
  # HTTP API port range start, each new node get port incremented (host machine)
  http_port_range_start = 10000
  # P2P API port range start, each new node get port incremented (host machine)
  p2p_port_range_start = 12000
```
- [PostgreSQL](../components/chainlink/nodeset.md): `13000` (we do not allow to have multiple databases for now, for simplicity)
```
    [nodeset.node_specs.db]
      # PostgreSQL volume name
      volume_name = "a"
      # PostgreSQL port (host machine)
      port = 13000
```

When you run `ctf d rm` database volume will be **removed**.

<div class="warning">

One node set is enough for any kind of testing, if you need more nodes consider extending your existing node set:
```
[nodeset]
  nodes = 10
```
</div>
