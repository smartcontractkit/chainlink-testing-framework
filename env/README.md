## Chainlink environment
This is an environment library we are using in tests, it provides:
- [cdk8s](https://cdk8s.io/) based wrappers
- High-level k8s API
- Automatic port forwarding

You can also use this library to spin up standalone environments

### Local k8s cluster
Read [here](KUBERNETES.md) about how to spin up a local cluster

#### Install
Set up deps, you need to have `node 14.x.x`, [helm](https://helm.sh/docs/intro/install/) and [yarn](https://classic.yarnpkg.com/lang/en/docs/install/#mac-stable)

Then use
```shell
make install_deps
```

### Usage
Create an env in a separate file and run it
```
export CHAINLINK_IMAGE="public.ecr.aws/chainlink/chainlink"
export CHAINLINK_TAG="1.4.0-root"
export CHAINLINK_ENV_USER="Satoshi"
go run examples/simple/env.go
```
For more features follow [tutorial](./TUTORIAL.md)

# Develop
#### Running standalone example environment
```shell
go run examples/simple/env.go
```
If you have another env of that type, you can connect by overriding environment name
```
ENV_NAMESPACE="..."  go run examples/chainlink/env.go
```

Add more presets [here](./presets)

Add more programmatic examples [here](./examples/)

If you have [chaosmesh]() installed in your cluster you can pull and generated CRD in go like that
```
make chaosmesh
```

If you need to check your system tests coverage, use [that](TUTORIAL.md#coverage)