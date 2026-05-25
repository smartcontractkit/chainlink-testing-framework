# Generating Chainlink Environment

Our code generation tools automatically build complete Chainlink environments and test templates, which minimizes documentation and provides a framework that is both structured and easily extensible.

Let's read `help` first and then build an environment for a single EVM network:
```bash
ctf gen -h
# generate a new Chainlink environment in "devenv" directory with 4 Chainlink nodes and one EVM network. Generate CLI called "pcli" and enter the shell
ctf gen env --cli pcli --product-name MyProduct --output-dir devenv --nodes 4
```

Follow further instructions in `devenv/README.md`

# Generating Infrastructure Testing Template

If you have Chainlink infrastructure already deployed it's useful to generate a workload + chaos suite template.
```bash
ctf gen load -h
# generate test suite named ChaosGen, with workload + default chaos experiments (fail + latency) for all the pods that have app.kubernetes.io/instance annotation
ctf gen load -w -n ChaosGen default
```