# 📦 Pods

High-level K8s API for developer envnrionments.

## Why

This framework introduces a lightweight abstraction layer, allowing developers to focus on product
configuration while abstracting away Kubernetes complexities.
It meant to be used with `CTFv2` framework to run components on `K8s`.

### Real world example (Chainlink Node Set)

In this example we'll spin up a local `Kind` cluster and deploy a Chainlink cluster,
just [40 lines](https://github.com/smartcontractkit/chainlink-testing-framework/pods/blob/master/examples/nodeset_test.go) of code (without product
configuration).

Follow the [README](./environment/README.md)

### Developing

Install pre-commit hooks and check available actions (lint, test ,etc)

```
just install
just
```

Add new features to `pods.go`, add new tests to `pods_test.go` and make pre-commit hooks and then CI pass.

Run `just test-deploy-cover` to check coverage.

Create additional directories for product-specific deployments if needed.

### Importing CRDs and K8s manifests

```
devbox shell
# filter current cluster CRDs, grep 'monitoring'
just crds monitoring
# save it as YAML
just crd servicemonitors.monitoring.coreos.com crds/monitoring.coreos.com
# import as Go code
just import crds/monitoring.coreos.com.yml
```

### K8s bindings versions

Check available `cdk8s` bindings [versions](https://github.com/cdk8s-team/cdk8s/tree/master/kubernetes-schemas).