# Mono Repository Tooling

In our multi-module Go repository, we use tools like:
- DevBox
- Just
- Matrix CI pattern

Open DevBox shell:
```
devbox shell
```

Install pre-commit hooks first:
```
just install
```

Each package has tests, run using commands in the justfile, examples:
```
# run all the tests (cache)
just test-all
# run package tests with regex
just test wasp TestSmoke
# run all package tests
just test tools/ghlatestreleasechecker ./...
```

Use linters:
```
# all packages
just lint-all
# one package
just lint wasp
```

For extra dependencies, we use [NixHub](https://www.nixhub.io/) to add them to [DevBox](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/devbox.json), which also works in [CI](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/.github/workflows/seth-test.yml#L62).
