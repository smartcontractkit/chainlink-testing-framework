# Installation Guide

This guide will help you set up the Chainlink Testing Framework (CTF) on your local machine.

## Prerequisites

### Required Software

1. **Docker**
   - **Recommended**: [OrbStack](https://orbstack.dev/) (faster, smaller memory footprint)
   - **Alternative**: [Docker Desktop](https://www.docker.com/products/docker-desktop/)
   
   Tested with:
   ```
   Docker version 27.3.1
   OrbStack Version: 1.8.2 (1080200)
   ```

2. **Golang**
   - Install the latest stable version from [go.dev](https://go.dev/doc/install)
   - Minimum version: Go 1.21+
   - Recommended: Go 1.22+

3. **Git**
   - Required for cloning repositories and managing dependencies

### Optional Software

- **[direnv](https://direnv.net/)** - Automatically load environment variables
- **[nix](https://nixos.org/manual/nix/stable/installation/installation.html)** - For development dependencies (used by some libraries)

## Installation Steps

### 1. Install Docker

#### macOS
```bash
# Install OrbStack (recommended)
brew install --cask orbstack

# Or install Docker Desktop
brew install --cask docker
```

#### Linux
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER
```

#### Windows
- Download and install [Docker Desktop for Windows](https://docs.docker.com/desktop/install/windows-install/)

### 2. Install Golang

#### macOS
```bash
# Using Homebrew
brew install go

# Or download from go.dev
curl -L https://go.dev/dl/go1.22.0.darwin-amd64.pkg -o go.pkg
sudo installer -pkg go.pkg -target /
```

#### Linux
```bash
# Download and install
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Windows
- Download from [go.dev](https://go.dev/dl/) and run the installer

### 3. Download CTF CLI

The CTF CLI provides helpful commands for managing test environments and observability stacks.

#### macOS ARM64 (M1/M2/M3)
```bash
curl -L https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-darwin-arm64.tar.gz | tar -xz
```

#### macOS AMD64 (Intel)
```bash
curl -L https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-darwin-amd64.tar.gz | tar -xz
```

#### Linux ARM64
```bash
curl -L https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-linux-arm64.tar.gz | tar -xz
```

#### Linux AMD64
```bash
curl -L https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-linux-amd64.tar.gz | tar -xz
```

#### Windows
```bash
# Using PowerShell
Invoke-WebRequest -Uri "https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/framework/v0.1.8/framework-v0.1.8-windows-amd64.tar.gz" -OutFile "framework.tar.gz"
tar -xzf framework.tar.gz
```

### 4. Set Up Environment

Create a `.envrc` file in your project directory:

```bash
# Create .envrc file
cat > .envrc << EOF
export TESTCONTAINERS_RYUK_DISABLED=true
EOF

# Load environment variables
source .envrc
```

**Note**: If you have `direnv` installed, the environment variables will be automatically loaded when you enter the directory.

### 5. Verify Installation

Test that everything is working:

```bash
# Check Docker
docker --version
docker run hello-world

# Check Go
go version

# Check CTF CLI
./framework --help
```

## Project Setup

### 1. Create a New Project

```bash
# Create project directory
mkdir my-ctf-project
cd my-ctf-project

# Initialize Go module
go mod init my-ctf-project

# Add framework dependency
go get github.com/smartcontractkit/chainlink-testing-framework/framework
```

### 2. Create Basic Configuration

Create a `smoke.toml` file:

```toml
[blockchain_a]
  type = "anvil"
```

### 3. Create Your First Test

Create `smoke_test.go`:

```go
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

func TestSmoke(t *testing.T) {
    in, err := framework.Load[Config](t)
    require.NoError(t, err)

    bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
    require.NoError(t, err)

    t.Run("blockchain is running", func(t *testing.T) {
        require.NotEmpty(t, bc.Nodes[0].ExternalHTTPUrl)
    })
}
```

### 4. Run Your First Test

```bash
# Run the test
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke

# Clean up containers
ctf d rm
```

## Troubleshooting

### Common Issues

#### Docker Permission Issues
```bash
# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker

# Restart Docker Desktop (macOS/Windows)
```

#### Port Conflicts
If you get port binding errors:
```bash
# Check what's using the port
lsof -i :8545  # Example for Ethereum RPC port

# Kill the process or change ports in your config
```

#### Memory Issues
If you encounter memory issues with Docker:
- Increase Docker memory limit in Docker Desktop settings
- Use OrbStack instead of Docker Desktop (better performance)

#### Go Module Issues
```bash
# Clean module cache
go clean -modcache

# Tidy dependencies
go mod tidy
```

### Getting Help

- Check the [Troubleshooting Guide](../Reference/Troubleshooting)
- Search [existing issues](https://github.com/smartcontractkit/chainlink-testing-framework/issues)
- Create a [new issue](https://github.com/smartcontractkit/chainlink-testing-framework/issues/new) with detailed information

## Next Steps

Now that you have CTF installed, you can:

1. [Write your first test](First-Test)
2. [Learn about configuration](Configuration)
3. [Set up observability](Environment-Setup)
4. [Explore the framework components](../Framework/Components)

## System Requirements

### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 10GB free space
- **OS**: macOS 10.15+, Ubuntu 18.04+, Windows 10+

### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 20GB+ free space
- **OS**: Latest stable versions
- **Docker**: OrbStack (macOS) or Docker Desktop with increased resources 