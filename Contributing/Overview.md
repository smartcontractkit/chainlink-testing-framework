# Contributing Guide

Thank you for your interest in contributing to the Chainlink Testing Framework (CTF)! We welcome contributions from the community to improve the framework, add features, fix bugs, and enhance documentation.

---

## How to Contribute

1. **Fork the repository** and create your branch from `main`.
2. **Open an issue** to discuss your proposed change if it's non-trivial.
3. **Write clear, maintainable code** and add tests for new features or bug fixes.
4. **Document your changes** in the code and/or wiki as appropriate.
5. **Open a Pull Request (PR)** with a clear description of your changes and link to any relevant issues.
6. **Participate in code review** and address feedback promptly.

---

## Code of Conduct

We are committed to fostering a welcoming and inclusive environment. Please read and follow our [Code of Conduct](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/CODE_OF_CONDUCT.md).

---

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/smartcontractkit/chainlink-testing-framework.git
   cd chainlink-testing-framework
   ```
2. **Install dependencies**
   - [Go](https://go.dev/doc/install) 1.21+
   - [Docker](https://www.docker.com/)
   - [direnv](https://direnv.net/) (optional, for environment management)
   - [nix](https://nixos.org/manual/nix/stable/installation/installation.html) (optional, for dev environments)
3. **Set up your environment**
   ```bash
   cp Getting-Started/Installation.md .
   # Follow the installation instructions for your OS
   ```
4. **Run tests locally**
   ```bash
   CTF_CONFIGS=smoke.toml go test -v
   ```
5. **Lint and format your code**
   ```bash
   go fmt ./...
   go vet ./...
   golangci-lint run
   ```

---

## Pull Request Process

1. **Open a draft PR** early to get feedback.
2. **Ensure all tests pass** and code is linted.
3. **Describe your changes** clearly in the PR description.
4. **Link to any related issues** (e.g., `Closes #123`).
5. **Request review** from maintainers or relevant code owners.
6. **Address review comments** and update your PR as needed.
7. **Wait for approval and merge** (maintainers will handle merging).

---

## Getting Help

- **Search issues**: [GitHub Issues](https://github.com/smartcontractkit/chainlink-testing-framework/issues)
- **Ask in discussions**: [GitHub Discussions](https://github.com/smartcontractkit/chainlink-testing-framework/discussions)
- **Open a new issue** for bugs, feature requests, or questions
- **Contact maintainers** via GitHub if you need further assistance

---

## Resources
- [Home](../Home)
- [Installation Guide](../Getting-Started/Installation)
- [API Reference](../Reference/API-Reference)
- [Troubleshooting](../Reference/Troubleshooting)
- [Code of Conduct](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/CODE_OF_CONDUCT.md)
- [LICENSE](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/LICENSE) 