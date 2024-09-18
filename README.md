<div align="center">

# Chainlink Testing Framework

[![Main branch breaking changes check](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/rc-breaking-changes.yaml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/rc-breaking-changes.yaml)
[![Lib tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Alib%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![WASP tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Awasp%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Seth tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Aseth%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Havoc tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Ahavoc%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/lib)
![Go Version](https://img.shields.io/github/go-mod/go-version/smartcontractkit/chainlink-testing-framework?filename=./lib/go.mod)
![Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

The Chainlink Testing Framework (CTF) is a blockchain development framework written in Go. Its primary purpose is to help chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help test their contracts, or even for those that aren't using chainlink.

If you're looking to implement a new chain integration for the testing framework, head over to the [blockchain](./blockchain/) directory for more info.

# Content

1. [Libraries](#libraries)

## Libraries

CTF monorepository contains a set of libraries:

- [Harness](lib/README.md) - Library to interact with different blockchains, create CL node jobs and use k8s and docker.
- [WASP](wasp/README.md) - Scalable protocol-agnostic load testing library for `Go`
- [Havoc](havoc/README.md) - Chaos testing library
- [Seth](seth/README.md) - Ethereum client library with transaction tracing and gas bumping