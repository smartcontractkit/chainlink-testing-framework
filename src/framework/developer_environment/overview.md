# Generating Chainlink Environment

Our code generation tools automatically build complete Chainlink environments and test templates, which minimizes documentation and provides a framework that is both structured and easily extensible.

We provide a single environment for both quick and local developer environment and production-ready environments in Kubernetes.

## Local Environment

Read `help` first and then build an environment for a single EVM network:

```bash
ctf gen -h
# generate a new Chainlink environment in "devenv" directory with 4 Chainlink nodes and one EVM network. Generate CLI called "pcli" and enter the shell
ctf gen env --cli pcli --product-name MyProduct --output-dir devenv --nodes 4
```

Follow further instructions in `devenv/README.md`

## Remote Environment

Infrastructure is deployed by internal operators so here we provide only configuration and tools to interact with deployed environments.

# Generating Infrastructure Testing Template

Generate performance and chaos testing template for a Kubernetes namespace.

```bash
ctf gen load -h
# generate test suite named ChaosGen, with workload + default chaos experiments (fail + latency) for all the pods that have app.kubernetes.io/instance annotation
ctf gen load -w -n ChaosGen default
```
