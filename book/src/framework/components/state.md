# Exposing Components (Data and ports)

We use static port ranges and volumes for all components to simplify Docker port management for developers.

This approach allows us to apply chaos testing to any container, ensuring it reconnects and retains the data needed for your tests.

When deploying a component, you can explicitly configure port ranges if the default ports don’t meet your needs.

Defaults are:
- [NodeSet](../components/chainlink/nodeset.md) (Node HTTP API): `10000..100XX`
- [NodeSet](../components/chainlink/nodeset.md) (Node P2P API): `12000..120XX`
- [NodeSet](../components/chainlink/nodeset.md) (Delve debugger): `40000..400XX` (if you are using debug image)
- Shared `PostgreSQL` volume is called `postgresql_data`
```
[[nodesets]]
  # HTTP API port range start, each new node get port incremented (host machine)
  http_port_range_start = 10000
  # P2P API port range start, each new node get port incremented (host machine)
  p2p_port_range_start = 12000
```
- [PostgreSQL](../components/chainlink/nodeset.md): `13000` (we do not allow to have multiple databases for now, for simplicity)
```
    [nodesets.node_specs.db]
      # PostgreSQL volume name
      volume_name = "a"
      # PostgreSQL port (host machine)
      port = 13000
```

When you run `ctf d rm` database volume will be **removed**.


<div class="warning">

One node set is enough for any kind of testing, if you need more nodes consider extending your existing node set:
```
[[nodesets]]
  nodes = 10
```
</div>

## Custom ports

You can also define a custom set of ports for any node.
```toml
[[nodesets]]
  name = "don"
  nodes = 5
  override_mode = "each"
  
  [nodesets.db]
    image = "postgres:12.0"

  [[nodesets.node_specs]]

    [nodesets.node_specs.node]
      # here we defined 2 new ports to listen and mapped them to our host machine
      # syntax is "host:docker"
      custom_ports = ["14000:15000"]
      image = "public.ecr.aws/chainlink/chainlink:v2.16.0"
```
