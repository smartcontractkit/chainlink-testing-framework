# Havoc

[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework/havoc.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc)

The `havoc` package is designed to facilitate chaos testing within [Kubernetes](https://kubernetes.io/) environments using [Chaos Mesh](https://chaos-mesh.org/). It offers a structured way to define, execute, and manage chaos experiments as code, directly integrated into Go applications or testing suites, simplifying the creation and control of Chaos Mesh experiments.

## Features

- **Chaos Object Management:** Create, update, pause, resume, and delete chaos experiments using Go structures and methods.
- **Lifecycle Hooks:** Utilize chaos listeners to hook into the lifecycle of chaos experiments.
- **Different Experiments:** Create and manage different types of chaos experiments to affect network, IO, K8s pods, and more.
- **Active Monitoring:** Monitor and react to the status of chaos experiments programmatically.

## Requirements

- [Go](https://go.dev/)
- A Kubernetes cluster with [Chaos Mesh installed](https://chaos-mesh.org/docs/quick-start/)

## Active Monitoring

`havoc` enhances chaos experiment observability through structured logging and Grafana annotations by implementing the [ChaosListener](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc#ChaosListener) interface.

The [ChaosLogger](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc#ChaosLogger) is the default implementation. It uses [zerolog](https://github.com/rs/zerolog) to provide structured, queryable logging of chaos events. It automatically logs key lifecycle events such as creation, start, pause, and termination of chaos experiments with detailed contextual information.

### Grafana Annotations

We recommend using [Grafana dashboards](https://grafana.com/) to monitor your chaos experiments, and provide the [SingleLineGrafanaAnnotator](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc#SingleLineGrafanaAnnotator), a `ChaosListener` that annotates dashboards with chaos experiment events so you can see in real time what your chaos experiment is doing.

You can also use the [RangeGrafanaAnnotator](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc#RangeGrafanaAnnotator) to show the full range of a chaos event's duration rather than a single line.

## Creating a Chaos Experiment

To create a chaos experiment, define the chaos object options, initialize a chaos experiment with NewChaos, and then call Create to start the experiment.

See [this runnable example](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc#ExampleNewChaos) of defining a chaos experiment.
