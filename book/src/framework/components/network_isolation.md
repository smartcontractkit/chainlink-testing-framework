# Containers Network Isolation

Some components can be isolated from internet. Since some of them doesn't have `iptables` and must be accessible from local host for debugging we isolate network on DNS level.

JobDistributor is isolated from internet by default to prevent applying manifest changes when run in tests.

NodeSet DNS isolation can be controlled with a flag
```
[[nodesets]]
  # NodeSet DNS can be isolated if needed
  no_dns = true
```