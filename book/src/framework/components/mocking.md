# Mocking Services

The framework aims to equip you with all the necessary tools to write end-to-end system-level tests, while still allowing the flexibility to mock third-party services that are not critical to your testing scope.

## Configuration
```toml
[fake]
  # port to start Gin server
  port = 9111
```

## Usage

See [full](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fake_test.go) example.

<div class="warning">

`host.docker.internal` is docker platform dependent!

Use `framework.HostDockerInternal()` to reference `host.docker.internal` in your tests, so they can work in GHA CI
</div>
