# Fake Services

The framework aims to equip you with all the necessary tools to write end-to-end system-level tests, while still allowing the flexibility to fake third-party services that are not critical to your testing scope.

## Local Usage without Docker (Go runtime)
```toml
[fake]
  # port to start Gin server
  port = 9111
```

See [full](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fake_test.go) example.

Run it
```
CTF_CONFIGS=fake.toml go test -v -run TestFakes
```

<div class="warning">

`host.docker.internal` is Docker platform dependent!

Use `framework.HostDockerInternal()` to reference `host.docker.internal` in your tests, so they can work in GHA CI
</div>

## Dockerized Usage

Copy this example into your project, write the logic of fake using `fake.JSON` and `fake.Func`, build and upload it and run.

## Install

To handle some utility command please install `Taskfile`
```
brew install go-task
```

## Private Repositories (Optional)

If your tests are in a private repository please generate a new SSH key and add it on [GitHub](https://github.com/settings/keys). Don't forget to click `Configure SSO` in UI
```
task new-ssh
```

## Usage

Build it and run locally when developing fakes
```
task build -- ${product-name}-${tag} # ex. myproduct-1.0
task run
```

Test it
```
curl "http://localhost:9111/static-fake"
curl "http://localhost:9111/dynamic-fake"
```
Publish it
```
task publish -- $tag
```

See full [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/fake_docker_test.go)

Run it
```toml
[fake]
  # the image you've built
  image = "test-fakes:myproduct-1.0"
  # port for Gin server
  port = 9111
```

```
CTF_CONFIGS=fake_docker.toml go test -v -run TestDockerFakes
```