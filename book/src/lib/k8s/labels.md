# k8s `chain.link` Labels

## Purpose
Resource labeling has been introduced to better associate Kubernetes (k8s) costs with products and teams. This document describes the labels used in the k8s cluster.

## Required Labels
Labels should be applied to all resources in the k8s cluster at three levels:
- **Namespace**
- **Workload**
- **Pod**

All three levels should include the following labels:
- `chain.link/team` - Name of the team that owns the resource.
- `chain.link/product` - Product that the resource belongs to.
- `chain.link/cost-center` - Product and framework name.

Additionally, pods should include the following label:
- `chain.link/component` - Name of the component.

### `chain.link/team`
This label represents the team responsible for the resource, but it might not be the team of the individual who created the resource. It should reflect the team the environment is **created for**.

For example, if you are a member of the Test Tooling team, but someone from the BIX team requests load tests, the namespace should be labeled as: `chain.link/team: bix`.

### `chain.link/product`
This label specifies the product the resource belongs to. Internally, some products may have alternative names (e.g., OCR instead of Data Feeds). To standardize data analysis, use the following names:

```
automation
bcm
ccip
data-feedsv1.0
data-feedsv2.0
data-feedsv3.0
data-streamsv0.3
data-streamsv1.0
deco
functions
proof-of-reserve
scale
staking
vrf
```

For example:
- OCR version 1: `data-feedsv1.0`
- OCR version 2: `data-feedsv2.0`

### `chain.link/cost-center`
This label serves as an umbrella for specific test or environment types and should rarely change. For load or soak tests using solutions provided by the Test Tooling team, use the convention: `test-tooling-<test-type>-test`

For example: `test-tooling-load-test`.

This allows easy distinction from load tests run using other tools.

### `chain.link/component`
This label identifies different components within the same product. Examples include:
- `chainlink` - Chainlink node.
- `geth` - Go-Ethereum blockchain node.
- `test-runner` - Remote test runner.

## Adding Labels to New Components
Adding a new component to an existing framework is discouraged. The recommended approach is to add the component to CRIB and make these labels part of the deployment templates.

If you need to add a new component, refer to the following sections in the [k8s Tutorial](./TUTORIAL.md):
- **Creating a new deployment part in Helm**
- **Creating a new deployment part in cdk8s**