# Containers Network Isolation

Some components can be isolated from internet. Since some of them doesn't have `iptables` and must be accessible from local host for debugging we isolate network on DNS level

```
[jd]
  # JobDistributor is isolated by default, this flag MUST not be changed!
  no_dns = true

[[nodesets]]
  # NodeSet DNS can be isolated if needed
  no_dns = true
```