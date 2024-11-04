# Interactive


For non-technical users or those building with Chainlink products outside of Golang, we offer an interactive method to deploy a NodeSet.

The only requirement is to have Docker running. 

If you're on OS X, we recommend to use [OrbStack](https://orbstack.dev/).

For other platforms use [Docker Desktop](https://www.docker.com/products/docker-desktop/).

Download the latest CLI [here](https://github.com/smartcontractkit/chainlink-testing-framework/releases/tag/framework%2Fv0.1.6)

Allow it to run in `System Settings -> Security Settings` (OS X)

![img.png](images/img.png)

```
./framework build node_set
```
Press `Ctrl+C` to remove the stack.

![img.png](images/interactive-node-set.png)