---
layout: default
title: Chainlink Testing Framework
nav_order: 1
description: "A general blockchain integration testing framework geared towards Chainlink projects"
permalink: /
---

# Chainlink Testing Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework)
![Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlink Testing Framework is a blockchain development and testing framework written in [Go](https://go.dev/). While the framework is designed primarily with testing Chainlink nodes in mind, it's not at all limited to that function. With this framework, blockchain developers can create extensive integration, e2e, performance, and chaos tests for almost anything!

Are you new to [blockchain development](https://ethereum.org/en/developers/docs/), [smart contracts](https://docs.chain.link/docs/beginners-tutorial/), or [Chainlink](https://chain.link/)? Learn more by clicking the links!

Here you'll find some guidelines on writing blockchain tests using this framework, and some tips on contributing to it. In most code examples presented, error checking is omitted for brevity's sake. **Please check your errors**.

Some notable packages we use include:

* [chainlink-env](https://github.com/smartcontractkit/chainlink-env)
* [zerolog](https://github.com/rs/zerolog)
* [Ginkgo and Gomega](https://github.com/onsi/ginkgo)
* [Kubernetes](https://github.com/kubernetes/kubernetes)
