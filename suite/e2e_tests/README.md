# e2e tests

## Alerting stack used in tests
chainlink nodes -> explorer -> kafka -> otpe -> prometheus

## Dependencies
[integrations-framework](github.com/smartcontractkit/integrations-framework) 

A framework for interacting with chainlink nodes, environments, and other blockchain systems.

[ginkgo](github.com/onsi/ginkgo) 

Expressive Behavior-Driven Development ("BDD") style tests.

[gomega](github.com/onsi/gomega)

A matcher/assertion library. It is best paired with the Ginkgo BDD test framework.


## Running

[instructions](https://onsi.github.io/ginkgo/#running-tests)

The tests expect a connection to a k8s cluster in which the tests will run.