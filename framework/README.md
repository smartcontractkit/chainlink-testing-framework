## Chainlink Testing Framework Harness

<!-- TOC -->
* [CLI](./cmd/README.md)
* [Components](./COMPONENTS.md)
* [Configuration](./CONFIGURATION.md)
* [Caching](./CACHING.md)
* [Local Observability Stack](./cmd/observability/README.md)
* [Examples](https://github.com/smartcontractkit/chainlink/tree/8e8597aa14c39c48ed4b3261f6080fa43b5d7cd0/e2e/capabilities)
<!-- TOC -->

This module includes the CTFv2 harness, a lightweight, modular, and data-driven framework designed for combining off-chain and on-chain components while implementing best practices for end-to-end system-level testing:

- **Non-nil configuration**: All test variables must have defaults, automatic validation.
- **Component isolation**: Components are decoupled via input/output structs, without exposing internal details.
- **Modular configuration**: No arcane knowledge of framework settings is required; the config is simply a reflection of the components being used in the test. Components declare their own configuration—'what you see is what you get.'
- **Replaceability and extensibility**: Since components are decoupled via outputs, any deployment component can be swapped with a real service without altering the test code.
- **Caching**: any component can use cached configs to skip environment setup for faster test development
- **Integrated observability stack**: use `ctf obs up` to spin up a local observability stack.


### Environment variables (Tests, when using in Go code)
|             Name             |                                                                      Description                                                                       |          Possible values | Default |        Required?         |
|:----------------------------:|:------------------------------------------------------------------------------------------------------------------------------------------------------:|-------------------------:|:-------:|:------------------------:|
|         CTF_CONFIGS          | Path(s) to test config files. <br/>Can be more than one, ex.: smoke.toml,smoke_1.toml,smoke_2.toml.<br/>First filepath will hold all the merged values | Any valid TOML file path |         |            ✅             |
|        CTF_LOG_LEVEL         |                                                                   Harness log level                                                                    | `info`, `debug`, `trace` | `info`  |            🚫            |
|       CTF_LOKI_STREAM        |                                                Streams all components logs to `Loki`, see params below                                                 |          `true`, `false` | `false` |            🚫            |
|           LOKI_URL           |                                            URL to `Loki` push api, should be like`${host}/loki/api/v1/push`                                            |                      URL |    -    | If you use `Loki` then ✅ |
|        LOKI_TENANT_ID        |                                                Streams all components logs to `Loki`, see params below                                                 |          `true`, `false` |    -    | If you use `Loki` then ✅ |
| TESTCONTAINERS_RYUK_DISABLED |                                   Testcontainers-Go reaper container, removes all the containers after the test exit                                   |          `true`, `false` | `false` |            🚫            |
|         RESTY_DEBUG          |                                                            Log all Resty client HTTP calls                                                             |          `true`, `false` | `false` |            🚫            |