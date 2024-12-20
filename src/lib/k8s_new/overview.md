# Kubernetes

> [!WARNING]
> It is highly recommended to use [CRIB](https://github.com/smartcontractkit/crib) for `k8s` deployments.
> Avoid running long tests that are impractical to execute locally or through CI pipelines.
>
> **Proceed at your own risk.**

---

## Overview

The `CTFv1` tool builds `k8s` environments **programmatically** using either `Helm` or `cdk8s` charts. This approach introduces significant complexity to the deployment process.

To manage long-running tests, `CTFv1` utilizes a `remote runner`, which is essentially a Docker container containing the test logic. This container is deployed as a `cdk8s`-based chart, creating a `k8s` resource of type `job` that runs the test in a detached manner. This setup requires custom logic to integrate with the test framework.

---

## What We’ll Cover

1. Creating a simplified `k8s` environment.
2. Adding a basic test that:
   - Deploys a smart contract.
   - Supports the `remote runner` capability.
3. Building a Docker image for the test and configuring the required environment variables.

---

Are you ready to get started?
