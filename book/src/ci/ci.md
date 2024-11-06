# Continuous Integration

Here we describe our good practices for structuring different types of tests in Continuous Integration (GitHub Actions).

The simplest flow can look like:

Set up secrets in your GitHub repository
```
gh secret set CTF_SIMULATED_KEY_1 --body "..."
```

Add a workflow
```yaml
name: Run E2E tests

on:
  push:

jobs:
  test:
    runs-on: ubuntu-latest
    env:
      CTF_CONFIGS: smoke.toml
      CTF_LOG_LEVEL: info
      CTF_LOKI_STREAM: "false"
      PRIVATE_KEY: ${{ secrets.CTF_SIMULATED_KEY_1 }}
    steps:
      - name: Check out code
        uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22.8
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: |
            ~/.cache/go-build
            ~/go/pkg/mod
          key: go-modules-${{ hashFiles('**/go.sum') }}-${{ runner.os }}
          restore-keys: |
            go-modules-${{ hashFiles('**/go.sum') }}-${{ runner.os }}
      - name: Install dependencies
        run: go mod download
      - name: Run tests
        working-directory: e2e/capabilities
        run: go test -v -run TestDON
```

If you need to structure a lot of different end-to-end tests follow [this](https://github.com/smartcontractkit/.github/tree/main/.github/workflows) guide.