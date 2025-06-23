# Chainlink Testing Framework (CTF) Wiki

<div align="center">

[![Documentation](https://img.shields.io/badge/Documentation-MDBook-blue?style=for-the-badge)](https://smartcontractkit.github.io/chainlink-testing-framework/overview.html)
[![Framework tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Aframework%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Lib tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Alib%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![WASP tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Awasp%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Seth tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Aseth%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)
[![Havoc tag](https://img.shields.io/github/v/tag/smartcontractkit/chainlink-testing-framework?filter=%2Ahavoc%2A)](https://github.com/smartcontractkit/chainlink-testing-framework/tags)

[![Tests](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/test.yaml)
[![Run all linters](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/linters.yml/badge.svg)](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/linters.yml)

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink-testing-framework)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink-testing-framework)
![Go Version](https://img.shields.io/github/go-mod/go-version/smartcontractkit/chainlink-testing-framework?filename=./lib/go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

</div>

## Overview

The **Chainlink Testing Framework (CTF)** is a comprehensive blockchain development framework written in Go, designed to help Chainlink developers create extensive integration, end-to-end, performance, and chaos tests to ensure the stability of the Chainlink project. It can also be helpful for developers who want to use Chainlink oracles in their projects or even for those not using Chainlink at all.

## 🎯 Primary Goals

- **Reduce complexity** of end-to-end testing
- **Enable tests to run in any environment**
- **Serve as a single source of truth** for system behavior
- **Provide modular, data-driven testing** capabilities
- **Support comprehensive observability** and monitoring

## 🏗️ Architecture

The CTF monorepository contains two major pieces:

### 1. **Framework** - Core Testing Infrastructure
- **Modular component system** for blockchain networks, Chainlink nodes, and other services
- **Configuration-driven testing** with TOML-based configs
- **Component isolation** and replaceability
- **Integrated observability stack** (metrics, logs, traces, profiles)
- **Caching system** for faster test development
- **Quick local environments** (15-second setup with caching)

### 2. **Libraries** - Specialized Testing Tools
- **[WASP](Libraries/WASP)** - Scalable protocol-agnostic load testing library
- **[Seth](Libraries/Seth)** - Reliable and debug-friendly Ethereum client
- **[Havoc](Libraries/Havoc)** - Chaos testing library for Kubernetes environments

## 🚀 Key Features

### Framework Features
- **Straightforward and sequential test composition** - Tests are readable with precise control
- **Modular configuration** - No arcane knowledge required, config reflects components used
- **Component isolation** - Components decoupled via input/output structs
- **Replaceability and extensibility** - Any deployment component can be swapped
- **Quick local environments** - Common setup in just 15 seconds
- **Caching** - Skip setup for faster test development
- **Integrated observability** - Metrics, logs, traces, and profiles

### Library Features
- **WASP**: Protocol-agnostic load testing with Grafana integration
- **Seth**: Transaction decoding, tracing, gas bumping, and multi-key support
- **Havoc**: Chaos engineering with Chaos Mesh integration

## 📚 Quick Start

### Prerequisites
- **Docker** ([OrbStack](https://orbstack.dev/) recommended) or Docker Desktop
- **Golang** (latest stable version)
- **CTF CLI** (download from releases)

### Basic Setup
```bash
# Create project directory
mkdir my-ctf-project && cd my-ctf-project

# Initialize Go module
go mod init my-ctf-project

# Add framework dependency
go get github.com/smartcontractkit/chainlink-testing-framework/framework

# Download CTF CLI (example for macOS ARM64)
curl -L https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-darwin-arm64.tar.gz | tar -xz

# Set up environment
echo 'export TESTCONTAINERS_RYUK_DISABLED=true' > .envrc
source .envrc
```

### Your First Test
```toml
# smoke.toml
[blockchain_a]
  type = "anvil"
```

```go
// smoke_test.go
package mymodule_test

import (
    "github.com/smartcontractkit/chainlink-testing-framework/framework"
    "github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
    "github.com/stretchr/testify/require"
    "testing"
)

type Config struct {
    BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestMe(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    t.Run("test something", func(t *testing.T) {
        require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
    })
}
```

```bash
# Run the test
CTF_CONFIGS=smoke.toml go test -v -run TestMe

# Clean up
ctf d rm
```

## 📖 Documentation Structure

### Getting Started
- [Installation Guide](Getting-Started/Installation)
- [First Test](Getting-Started/First-Test)
- [Configuration Basics](Getting-Started/Configuration)
- [Environment Setup](Getting-Started/Environment-Setup)

### Framework
- [Framework Overview](Framework/Overview)
- [Component System](Framework/Components)
- [Configuration Management](Framework/Configuration)
- [Observability Stack](Framework/Observability)
- [Caching System](Framework/Caching)
- [Test Patterns](Framework/Test-Patterns)

### Libraries
- [WASP - Load Testing](Libraries/WASP)
- [Seth - Ethereum Client](Libraries/Seth)
- [Havoc - Chaos Testing](Libraries/Havoc)

### Advanced Topics
- [Component Development](Advanced/Component-Development)
- [Custom Components](Advanced/Custom-Components)
- [Performance Testing](Advanced/Performance-Testing)
- [Chaos Engineering](Advanced/Chaos-Engineering)
- [CI/CD Integration](Advanced/CI-CD-Integration)

### Examples & Tutorials
- [Basic Examples](Examples/Basic-Examples)
- [Advanced Examples](Examples/Advanced-Examples)
- [Real-world Use Cases](Examples/Use-Cases)

### Reference
- [API Reference](Reference/API-Reference)
- [Configuration Reference](Reference/Configuration-Reference)
- [CLI Reference](Reference/CLI-Reference)
- [Troubleshooting](Reference/Troubleshooting)

## 🎯 Use Cases

### For Chainlink Developers
- **Integration Testing**: Test Chainlink nodes with various blockchain networks
- **End-to-End Testing**: Complete workflow testing from contract deployment to oracle responses
- **Performance Testing**: Load testing with WASP library
- **Chaos Testing**: Failure scenario testing with Havoc
- **Upgrade Testing**: Version compatibility and migration testing

### For Blockchain Developers
- **Smart Contract Testing**: Deploy and test contracts with Seth
- **Multi-Chain Testing**: Support for Ethereum, Solana, TON, and more
- **Network Simulation**: Local blockchain networks for development
- **Gas Optimization**: Transaction tracing and gas analysis

### For DevOps Engineers
- **Infrastructure Testing**: Test deployment configurations
- **Observability**: Integrated monitoring and logging
- **CI/CD Integration**: Automated testing pipelines
- **Environment Management**: Consistent test environments

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](Contributing/Overview) for details on:
- Code of Conduct
- Development Setup
- Testing Guidelines
- Pull Request Process
- Issue Reporting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/LICENSE) file for details.

## 🔗 Links

- **[Repository](https://github.com/smartcontractkit/chainlink-testing-framework)**
- **[Documentation](https://smartcontractkit.github.io/chainlink-testing-framework/)**
- **[Issues](https://github.com/smartcontractkit/chainlink-testing-framework/issues)**
- **[Discussions](https://github.com/smartcontractkit/chainlink-testing-framework/discussions)**
- **[Releases](https://github.com/smartcontractkit/chainlink-testing-framework/releases)**

---

*This wiki provides comprehensive documentation for the Chainlink Testing Framework. For the most up-to-date information, always refer to the [official documentation](https://smartcontractkit.github.io/chainlink-testing-framework/).* 