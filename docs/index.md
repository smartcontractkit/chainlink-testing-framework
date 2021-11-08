---
layout: default
title: Home
nav_order: 1
description: "A general blockchain integration testing framework geared towards Chainlink projects"
permalink: /
---

# Chainlink Integrations Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlnk Integrations Framework is a blockchain development and testing framework written in Go. While the framework
is designed primarily with testing Chainlink nodes in mind, it's not at all limited to that function. With this
framework, blockchain developers can create extensive integration, e2e, performance, and chaos tests for almost anything!

Are you new to [blockchain development](https://ethereum.org/en/developers/docs/),
[smart contracts](https://docs.chain.link/docs/beginners-tutorial/),
or [Chainlink](https://chain.link/)? Learn more by clicking the links!

## Kubernetes Setup

In order to use this framework, you must have a connection to an actively running
[Kubernetes cluster](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/). If you don't have
one handy, check out [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests,
but if you write tests that make use of multiple Chainlink nodes, or try to run many tests in parallel, you'll likely
find minikube inadequeate. A Kubernetes cluster with 4 vCPU and 10 GB RAM is a good starting point for when you start
to notice issues.

**The framework will use whatever your primary KUBECONFIG cluster is.** Learn more
about setting up Kubernetes [here](https://kubernetes.io/docs/setup/).

### Why?

There's a lot of different components to bring up for each test, most of which involve:

* A simulated blockchain
* Some number of Chainlink nodes
* An equal number of postgres DBs to support the chainlink nodes
* At least one external adapter

Following the good testing practice of having clean, non-dependent test environments means we're creating a lot of these
components for each test, and tearing them down soon after. In order to organize these test environments, and after
finding `docker compose` to be woefully inadequate after a certain point, kubernetes was the obvious choice.
