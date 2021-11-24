# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlink Integration Framework is a blockchain development framework written in Go. Its primary purpose is to help
chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the
chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help
test their contracts, or even for those that aren't using chainlink.

See the [docs](https://smartcontractkit.github.io/integrations-framework/) or our
[go reference](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework) page for more detailed info and
examples. If you just want a quick overview, keep reading.

## WIP

As of now, this framework is still very much a work in progress, and will have frequent changes, many of which will probably be breaking.

**As of Monday, November 22, 2021, there has been a massive overhaul of how the framework works. Namely use of the [helmenv](https://github.com/smartcontractkit/helmenv) library**

## Setup

In order to use this framework, you must have a connection to an actively running Kubernetes cluster. If you don't have
one handy, check out [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests,
but will likely need to be allocated more power, or you'll need to use a more powerful cluster in general to run tests
that require lots of services, like OCR.

## Usage

Here's a simple example on deploying and interacting with a basic storage contract using this framework and
[Ginkgo](https://github.com/onsi/ginkgo), a BDD testing framework we've come to really enjoy. You can use another testing
framework, including Go's default testing if you prefer otherwise.

See our [suite/smoke](suite/smoke) directory for quite a few examples of the framework's usage.

## Chainlink Values

If you would like to change the Chainlink values that are used for environments, you can use JSON to squash them. Have a look over at our [helmenv](https://github.com/smartcontractkit/helmenv/tree/v1.0.5/charts/chainlink) chainlink charts to get a grasp of how things are structured. We'll be writing more on this later, but for now, you can squash values by providing a `CHARTS` environment variable.

```sh
CHARTS='{"chainlink": {"values": {"chainlink": {"image": {"version": "<version>"}}}}}' make test_smoke args="-nodes=5"
```
