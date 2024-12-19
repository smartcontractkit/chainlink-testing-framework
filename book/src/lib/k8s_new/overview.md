# Kubernetes

> [!WARNING]
> It is highly enouraged that your used [CRIB](https://github.com/smartcontractkit/crib) for `k8s` deployments
> and that you do not use long running tests that are too long to be executed from your local machine or from the CI.
>
> **Proceed at your own risk and peril**.

First of all, `CTFv1` builds `k8s` environments **programmatically** using either `Helm` or `cdk8s` charts. That adds
non-trivial complexity to the deployment process.

Second of all, to manage long-running test it uses a `remote runner`, which is a Docker container with tests that
executes as a `cdk8s`-based chart that creates a `k8s` resource of `job` type to execute the test in a detached manner.
And that necessiates adding some custom logic to the test.

Here we will first look at creation of a very simplified `k8s` environment and then adding an even simpler test,
that will deploy a smart contract and will add support a `remote runner` capability. And finally we will build a
Docker image with our test and set all required environment variables.

Ready?