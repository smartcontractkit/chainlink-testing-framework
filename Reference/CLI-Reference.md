# CLI Reference

This page documents the main commands and options for the Chainlink Testing Framework (CTF) CLI.

---

## Overview

The CTF CLI is used to manage test environments, observability stacks, and component lifecycles. It is distributed as a standalone binary (`framework`).

---

## Main Commands

### 1. Environment Management

- `ctf obs up` – Start the observability stack (Grafana, Loki, Prometheus, Pyroscope)
- `ctf obs down` – Stop the observability stack
- `ctf bs up` – Start the Blockscout stack
- `ctf bs down` – Stop the Blockscout stack

### 2. Component Management

- `ctf d rm` – Remove all running containers and clean up resources
- `ctf d ls` – List running containers

### 3. Test Execution

- `go test -v` – Run tests (uses CTF config and environment)
- `CTF_CONFIGS=smoke.toml go test -v -run TestName` – Run a specific test with a config

### 4. Configuration Validation

- `ctf validate config.toml` – Validate a TOML configuration file

---

## Common Flags

- `--help` – Show help for any command
- `--version` – Show CLI version
- `--log-level` – Set log verbosity (e.g., `info`, `debug`)
- `--config` – Specify a config file (alternative to `CTF_CONFIGS`)

---

## Usage Examples

```bash
# Start observability stack
ctf obs up

# Run a test with a specific config
CTF_CONFIGS=smoke.toml go test -v -run TestSmoke

# Remove all containers
ctf d rm

# Validate a config file
ctf validate myconfig.toml
```

---

## More
- For full CLI documentation, run `ctf --help` or see the [README](https://github.com/smartcontractkit/chainlink-testing-framework#readme) 