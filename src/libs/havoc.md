## Havoc

The `havoc` package is a Go library designed to facilitate chaos testing within Kubernetes environments using Chaos Mesh. It offers a structured way to define, execute, and manage chaos experiments as code, directly integrated into Go applications or testing suites. This package simplifies the creation and control of Chaos Mesh experiments, including network chaos, pod failures, and stress testing on Kubernetes clusters.

### Features

- **Chaos Object Management:** Easily create, update, pause, resume, and delete chaos experiments using Go structures and methods.
- **Lifecycle Hooks:** Utilize chaos listeners to hook into lifecycle events of chaos experiments, such as creation, start, pause, resume, and finish.
- **Support for Various Chaos Experiments:** Create and manage different types of chaos experiments like NetworkChaos, IOChaos, StressChaos, PodChaos, and HTTPChaos.
- **Chaos Experiment Status Monitoring:** Monitor and react to the status of chaos experiments programmatically.

### Installation

To use `havoc` in your project, ensure you have a Go environment setup. Then, install the package using go get:

```
go get -u github.com/smartcontractkit/chainlink-testing-framework/havoc
```

Ensure your Kubernetes cluster is accessible and that you have Chaos Mesh installed and configured.

### Monitoring and Observability in Chaos Experiments

`havoc` enhances chaos experiment observability through structured logging and Grafana annotations, facilitated by implementing the ChaosListener interface. This approach allows for detailed monitoring, debugging, and visual representation of chaos experiments' impact.

#### Structured Logging with ChaosLogger

`ChaosLogger` leverages the zerolog library to provide structured, queryable logging of chaos events. It automatically logs key lifecycle events such as creation, start, pause, and termination of chaos experiments, including detailed contextual information.

Instantiate `ChaosLogger` and register it as a listener to your chaos experiments:

```
logger := havoc.NewChaosLogger()
chaos.AddListener(logger)
```

### Default package logger

`havoc/logger.go` contains default `Logger` instance for the package.

#### Visual Monitoring with Grafana Annotations

`SingleLineGrafanaAnnotator` is a `ChaosListener` that annotates Grafana dashboards with chaos experiment events. This visual representation helps correlate chaos events with their effects on system metrics and logs.

Initialize `SingleLineGrafanaAnnotator` with your Grafana instance details and register it alongside `ChaosLogger`:

```
annotator := havoc.NewSingleLineGrafanaAnnotator(
    "http://grafana-instance.com",
    "grafana-access-token",
    "dashboard-uid",
)
chaos.AddListener(annotator)
```

### Creating a Chaos Experiment

To create a chaos experiment, define the chaos object options, initialize a chaos experiment with NewChaos, and then call Create to start the experiment.

Here is an example of creating and starting a PodChaos experiment:

```
package main

import (
    "context"
    "github.com/smartcontractkit/chainlink-testing-framework/havoc"
    "github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "time"
)

func main() {
    // Initialize dependencies
    client, err := havoc.NewChaosMeshClient()
    if err != nil {
        panic(err)
    }
    logger := havoc.NewChaosLogger()
    annotator := havoc.NewSingleLineGrafanaAnnotator(
        "http://grafana-instance.com",
        "grafana-access-token",
        "dashboard-uid",
    )

    // Define chaos experiment
    podChaos := &v1alpha1.PodChaos{ /* PodChaos spec */ }
    chaos, err := havoc.NewChaos(havoc.ChaosOpts{
        Object:      podChaos,
        Description: "Pod failure example",
        DelayCreate: 5 * time.Second,
        Client:      client,
    })
    if err != nil {
        panic(err)
    }

    // Register listeners
    chaos.AddListener(logger)
    chaos.AddListener(annotator)

    // Start chaos experiment
    chaos.Create(context.Background())

    // Manage chaos lifecycle...
}
```

### Test Example

```
func TestChaosDON(t *testing.T) {
	testDuration := time.Minute * 60

    // Load test config
	cfg := &config.MercuryQAEnvChaos{}

	// Define chaos experiments and their schedule

	k8sClient, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	// Test 3.2: Disable 2 nodes simultaneously

	podFailureChaos4, err := k8s_chaos.MercuryPodChaosSchedule(k8s_chaos.MercuryScheduledPodChaosOpts{
		Name:        "schedule-don-ocr-node-failure-4",
		Description: "Disable 2 nodes (clc-ocr-mercury-arb-testnet-qa-nodes-3 and clc-ocr-mercury-arb-testnet-qa-nodes-4)",
		DelayCreate: time.Minute * 0,
		Duration:    time.Minute * 20,
		Namespace:   cfg.ChaosNodeNamespace,
		PodSelector: v1alpha1.PodSelector{
			Mode: v1alpha1.AllMode,
			Selector: v1alpha1.PodSelectorSpec{
				GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
					Namespaces: []string{cfg.ChaosNodeNamespace},
					ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
						{
							Key:      "app.kubernetes.io/instance",
							Operator: "In",
							Values: []string{
								"clc-ocr-mercury-arb-testnet-qa-nodes-3",
								"clc-ocr-mercury-arb-testnet-qa-nodes-4",
							},
						},
					},
				},
			},
		},
		Client: k8sClient,
	})
	require.NoError(t, err)

	// Test 3.3: Disable 3 nodes simultaneously

	podFailureChaos5, err := k8s_chaos.MercuryPodChaosSchedule(k8s_chaos.MercuryScheduledPodChaosOpts{
		Name:        "schedule-don-ocr-node-failure-5",
		Description: "Disable 3 nodes (clc-ocr-mercury-arb-testnet-qa-nodes-3, clc-ocr-mercury-arb-testnet-qa-nodes-4 and clc-ocr-mercury-arb-testnet-qa-nodes-5)",
		DelayCreate: time.Minute * 40,
		Duration:    time.Minute * 20,
		Namespace:   cfg.ChaosNodeNamespace,
		PodSelector: v1alpha1.PodSelector{
			Mode: v1alpha1.AllMode,
			Selector: v1alpha1.PodSelectorSpec{
				GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
					Namespaces: []string{cfg.ChaosNodeNamespace},
					ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
						{
							Key:      "app.kubernetes.io/instance",
							Operator: "In",
							Values: []string{
								"clc-ocr-mercury-arb-testnet-qa-nodes-3",
								"clc-ocr-mercury-arb-testnet-qa-nodes-4",
								"clc-ocr-mercury-arb-testnet-qa-nodes-5",
							},
						},
					},
				},
			},
		},
		Client: k8sClient,
	})
	require.NoError(t, err)

	chaosList := []havoc.ChaosEntity{
		podFailureChaos4,
		podFailureChaos5,
	}

	for _, chaos := range chaosList {
		chaos.AddListener(havoc.NewChaosLogger())
		chaos.AddListener(havoc.NewSingleLineGrafanaAnnotator(cfg.GrafanaURL, cfg.GrafanaToken, cfg.GrafanaDashboardUID))

		// Fail the test if the chaos object already exists
		exists, err := havoc.ChaosObjectExists(chaos.GetObject(), k8sClient)
		require.NoError(t, err)
		require.False(t, exists, "chaos object already exists: %s. Delete it before starting the test", chaos.GetChaosName())

		chaos.Create(context.Background())
	}

	t.Cleanup(func() {
		for _, chaos := range chaosList {
			// Delete chaos object if it still exists
			chaos.Delete(context.Background())
		}
	})

	// Simulate user activity/load for the duration of the chaos experiments
	runUserLoad(t, cfg, testDuration)
}
```
