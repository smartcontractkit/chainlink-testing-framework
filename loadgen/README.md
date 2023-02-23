# Loadgen

A simple protocol agnostic load tool for `Go` code

### Motivation

We are using a lot of custom protocols and `blockchain` interactions in our tests

A lot of wide-used tools like `ddosify` `k6` `JMeter` `Locust` can't cover our cases or using them going to create additional efforts

They contain a lot of functionality that we don't need:
- Extensive `HTTP` testing
- `JSON/UI` configuration
- Complicated load scheduling code
- Vendor locked clusterization and scaling solutions

They can't allow us to:
- Reuse our `Go` code for custom protocols and `blockchains`
- Reuse our test setup logic and `chaos` experiments
- Unify all testing/setup logic in `Go`

So we've implemented a simple tool with goals in mind:
- Reuse our `Go` code, write consistent setup/test scripts
- Have slim codebase (500-1k loc)
- Be able to perform synthetic load testing for request-based protocols in `Go` with `RPS bound load` (http, tcp, etc.)
- Be able to perform synthetic load testing for streaming protocols in `Go` with `Instances bound load` (ws, wsrpc, etc.)
- Be scalable in `k8s` without complicated configuration or vendored UI interfaces
- Be non-opinionated about reporting, be able to push any arbitrary data to `Loki` without sacrificing performance

## How to use
Docs is `TBD` for the time being, check test implementations

Dashboard is private, you can find it in `K8s tests -> Loadgen`

- [examples](loadgen_example_test.go) for full-fledged generator tests with `Loki`, you can also use them to validate performance of `loadgen`

- [implementation](loadgen_gun_mock.go) of an `RPS` type `gun`

- [implementation](loadgen_instance_mock.go) of an `Instance` type `gun`