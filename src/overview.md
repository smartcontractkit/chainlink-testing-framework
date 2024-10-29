# Intro


The Chainlink Testing Framework (CTF) is a blockchain development framework written in Go. Its primary purpose is to help chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help test their contracts, or even for those that aren't using chainlink.


# Content

1. [Libraries](#libraries)
2. [Releasing](#releasing)

## Libraries

CTF monorepository contains a set of libraries:

- [Framework](framework.md) - Library to interact with different blockchains, create CL node jobs and use k8s and docker.
- [WASP](wasp.md) - Scalable protocol-agnostic load testing library for `Go`
- [Havoc](havoc.md) - Chaos testing library
- [Seth](seth.md) - Ethereum client library with transaction tracing and gas bumping

## Releasing

We follow [SemVer](https://semver.org/) and follow best Go practices for releasing our modules, please follow the [instruction](RELEASE.md)