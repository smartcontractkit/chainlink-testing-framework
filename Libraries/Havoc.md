# Havoc - Chaos Testing Library

[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/chainlink-testing-framework/havoc.svg)](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc)

## Overview

The `havoc` package is designed to facilitate chaos testing within [Kubernetes](https://kubernetes.io/) environments using [Chaos Mesh](https://chaos-mesh.org/). It offers a structured way to define, execute, and manage chaos experiments as code, directly integrated into Go applications or testing suites, simplifying the creation and control of Chaos Mesh experiments.

## Goals

- **Chaos Object Management** - Create, update, pause, resume, and delete chaos experiments using Go structures and methods
- **Lifecycle Hooks** - Utilize chaos listeners to hook into the lifecycle of chaos experiments
- **Different Experiments** - Create and manage different types of chaos experiments to affect network, IO, K8s pods, and more
- **Active Monitoring** - Monitor and react to the status of chaos experiments programmatically
- **Observability Integration** - Structured logging and Grafana annotations for chaos experiment monitoring

## Key Features

### 1. **Chaos Experiment Management**
- Create and manage Chaos Mesh experiments programmatically
- Full lifecycle control (create, pause, resume, delete)
- Type-safe experiment configuration

### 2. **Lifecycle Hooks**
- Implement custom listeners for chaos experiment events
- React to experiment state changes
- Integrate with monitoring and alerting systems

### 3. **Multiple Experiment Types**
- **Network Chaos** - Network latency, packet loss, bandwidth limits
- **Pod Chaos** - Pod failure, pod kill, container kill
- **IO Chaos** - IO latency, IO error injection
- **Kernel Chaos** - Kernel panic, memory corruption
- **Time Chaos** - Clock skew, time manipulation
- **DNS Chaos** - DNS error injection, DNS spoofing
- **HTTP Chaos** - HTTP error injection, HTTP latency

### 4. **Observability Integration**
- Structured logging with [zerolog](https://github.com/rs/zerolog)
- Grafana dashboard annotations
- Prometheus metrics integration
- Custom event listeners

## Requirements

- [Go](https://go.dev/) 1.21+
- A Kubernetes cluster with [Chaos Mesh installed](https://chaos-mesh.org/docs/quick-start/)
- [k8s.io/client-go](https://github.com/kubernetes/client-go) for Kubernetes API access

## Installation

### Using Go Modules
```bash
go get github.com/smartcontractkit/chainlink-testing-framework/havoc
```

### Dependencies
```bash
go get k8s.io/client-go
go get k8s.io/apimachinery
```

## Quick Start

### Basic Chaos Experiment

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/smartcontractkit/chainlink-testing-framework/havoc"
    "github.com/smartcontractkit/chainlink-testing-framework/havoc/chaosmesh"
)

func main() {
    // Create chaos client
    client, err := havoc.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    // Define network chaos experiment
    networkChaos := &chaosmesh.NetworkChaos{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "network-latency-test",
            Namespace: "default",
        },
        Spec: chaosmesh.NetworkChaosSpec{
            Action: chaosmesh.NetworkDelayAction,
            Mode:   chaosmesh.OneMode,
            Selector: chaosmesh.PodSelectorSpec{
                Namespaces: []string{"default"},
                LabelSelectors: map[string]string{
                    "app": "my-app",
                },
            },
            Delay: &chaosmesh.DelaySpec{
                Latency:     "100ms",
                Correlation: "100",
                Jitter:      "0ms",
            },
            Duration: &chaosmesh.Duration{
                Duration: "30s",
            },
        },
    }
    
    // Create chaos experiment
    chaos, err := havoc.NewChaos(client, networkChaos)
    if err != nil {
        log.Fatal(err)
    }
    
    // Start the experiment
    err = chaos.Create(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Chaos experiment started")
    
    // Wait for experiment to complete
    time.Sleep(35 * time.Second)
    
    // Clean up
    err = chaos.Delete(context.Background())
    if err != nil {
        log.Printf("Failed to delete chaos experiment: %v", err)
    }
}
```

### With Lifecycle Listeners

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/smartcontractkit/chainlink-testing-framework/havoc"
    "github.com/smartcontractkit/chainlink-testing-framework/havoc/chaosmesh"
)

// Custom chaos listener
type MyChaosListener struct{}

func (l *MyChaosListener) OnChaosCreated(chaos *havoc.Chaos) {
    log.Printf("Chaos experiment created: %s", chaos.Name())
}

func (l *MyChaosListener) OnChaosStarted(chaos *havoc.Chaos) {
    log.Printf("Chaos experiment started: %s", chaos.Name())
}

func (l *MyChaosListener) OnChaosPaused(chaos *havoc.Chaos) {
    log.Printf("Chaos experiment paused: %s", chaos.Name())
}

func (l *MyChaosListener) OnChaosResumed(chaos *havoc.Chaos) {
    log.Printf("Chaos experiment resumed: %s", chaos.Name())
}

func (l *MyChaosListener) OnChaosDeleted(chaos *havoc.Chaos) {
    log.Printf("Chaos experiment deleted: %s", chaos.Name())
}

func main() {
    client, err := havoc.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    
    // Create pod chaos experiment
    podChaos := &chaosmesh.PodChaos{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "pod-failure-test",
            Namespace: "default",
        },
        Spec: chaosmesh.PodChaosSpec{
            Action: chaosmesh.PodFailureAction,
            Mode:   chaosmesh.OneMode,
            Selector: chaosmesh.PodSelectorSpec{
                Namespaces: []string{"default"},
                LabelSelectors: map[string]string{
                    "app": "my-app",
                },
            },
            Duration: &chaosmesh.Duration{
                Duration: "10s",
            },
        },
    }
    
    // Create chaos with listener
    chaos, err := havoc.NewChaos(client, podChaos)
    if err != nil {
        log.Fatal(err)
    }
    
    // Add custom listener
    chaos.AddListener(&MyChaosListener{})
    
    // Start experiment
    err = chaos.Create(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    // Wait and clean up
    time.Sleep(15 * time.Second)
    chaos.Delete(context.Background())
}
```

## Experiment Types

### Network Chaos

```go
// Network latency
networkChaos := &chaosmesh.NetworkChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "network-latency",
        Namespace: "default",
    },
    Spec: chaosmesh.NetworkChaosSpec{
        Action: chaosmesh.NetworkDelayAction,
        Mode:   chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
        Delay: &chaosmesh.DelaySpec{
            Latency:     "200ms",
            Correlation: "100",
            Jitter:      "50ms",
        },
        Duration: &chaosmesh.Duration{
            Duration: "60s",
        },
    },
}

// Network packet loss
networkChaos := &chaosmesh.NetworkChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "network-loss",
        Namespace: "default",
    },
    Spec: chaosmesh.NetworkChaosSpec{
        Action: chaosmesh.NetworkLossAction,
        Mode:   chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
        Loss: &chaosmesh.LossSpec{
            Loss:        "50",
            Correlation: "100",
        },
        Duration: &chaosmesh.Duration{
            Duration: "30s",
        },
    },
}
```

### Pod Chaos

```go
// Pod failure
podChaos := &chaosmesh.PodChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "pod-failure",
        Namespace: "default",
    },
    Spec: chaosmesh.PodChaosSpec{
        Action: chaosmesh.PodFailureAction,
        Mode:   chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
        Duration: &chaosmesh.Duration{
            Duration: "30s",
        },
    },
}

// Pod kill
podChaos := &chaosmesh.PodChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "pod-kill",
        Namespace: "default",
    },
    Spec: chaosmesh.PodChaosSpec{
        Action: chaosmesh.PodKillAction,
        Mode:   chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
    },
}
```

### IO Chaos

```go
// IO latency
ioChaos := &chaosmesh.IOChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "io-latency",
        Namespace: "default",
    },
    Spec: chaosmesh.IOChaosSpec{
        Action: chaosmesh.IODelayAction,
        Mode:   chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
        Delay: "100ms",
        Percent: 50,
        Path: "/var/log",
        Methods: []string{"read", "write"},
        Duration: &chaosmesh.Duration{
            Duration: "60s",
        },
    },
}
```

### Time Chaos

```go
// Clock skew
timeChaos := &chaosmesh.TimeChaos{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "clock-skew",
        Namespace: "default",
    },
    Spec: chaosmesh.TimeChaosSpec{
        Mode: chaosmesh.OneMode,
        Selector: chaosmesh.PodSelectorSpec{
            Namespaces: []string{"default"},
            LabelSelectors: map[string]string{"app": "my-app"},
        },
        TimeOffset: "1h",
        Duration: &chaosmesh.Duration{
            Duration: "30s",
        },
    },
}
```

## Observability Integration

### Structured Logging

Havoc provides structured logging using zerolog:

```go
import (
    "github.com/rs/zerolog/log"
)

// Enable debug logging
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// Chaos events are automatically logged
chaos, err := havoc.NewChaos(client, networkChaos)
if err != nil {
    log.Fatal().Err(err).Msg("Failed to create chaos")
}
```

### Grafana Annotations

Integrate with Grafana dashboards:

```go
import (
    "github.com/smartcontractkit/chainlink-testing-framework/havoc"
)

// Create Grafana annotator
grafanaAnnotator := havoc.NewSingleLineGrafanaAnnotator(
    "http://localhost:3000",
    "your-token",
    []string{"dashboard-uid"},
)

// Add to chaos experiment
chaos.AddListener(grafanaAnnotator)

// Start experiment (annotations will be created automatically)
err = chaos.Create(context.Background())
```

### Range Annotations

For experiments with duration, use range annotations:

```go
// Create range annotator
rangeAnnotator := havoc.NewRangeGrafanaAnnotator(
    "http://localhost:3000",
    "your-token",
    []string{"dashboard-uid"},
)

chaos.AddListener(rangeAnnotator)
```

## Advanced Usage

### Multiple Experiments

```go
func runChaosSuite() error {
    client, err := havoc.NewClient()
    if err != nil {
        return err
    }
    
    experiments := []chaosmesh.Chaos{
        createNetworkLatency(),
        createPodFailure(),
        createIOChaos(),
    }
    
    for _, exp := range experiments {
        chaos, err := havoc.NewChaos(client, exp)
        if err != nil {
            return err
        }
        
        // Start experiment
        err = chaos.Create(context.Background())
        if err != nil {
            return err
        }
        
        // Wait for completion
        time.Sleep(30 * time.Second)
        
        // Clean up
        chaos.Delete(context.Background())
    }
    
    return nil
}
```

### Conditional Chaos

```go
func conditionalChaos(condition bool) error {
    if !condition {
        return nil
    }
    
    client, err := havoc.NewClient()
    if err != nil {
        return err
    }
    
    chaos, err := havoc.NewChaos(client, createPodChaos())
    if err != nil {
        return err
    }
    
    return chaos.Create(context.Background())
}
```

### Chaos with Recovery

```go
func chaosWithRecovery() error {
    client, err := havoc.NewClient()
    if err != nil {
        return err
    }
    
    chaos, err := havoc.NewChaos(client, createNetworkChaos())
    if err != nil {
        return err
    }
    
    // Start chaos
    err = chaos.Create(context.Background())
    if err != nil {
        return err
    }
    
    // Monitor system health
    go func() {
        for {
            if isSystemUnhealthy() {
                // Pause chaos if system is unhealthy
                chaos.Pause(context.Background())
                break
            }
            time.Sleep(5 * time.Second)
        }
    }()
    
    // Wait for experiment duration
    time.Sleep(60 * time.Second)
    
    // Clean up
    return chaos.Delete(context.Background())
}
```

## Testing Integration

### With CTF Framework

```go
func TestChaosResilience(t *testing.T) {
    // Setup your application
    app := setupApplication(t)
    
    // Create chaos client
    client, err := havoc.NewClient()
    require.NoError(t, err)
    
    // Create chaos experiment
    chaos, err := havoc.NewChaos(client, createPodChaos())
    require.NoError(t, err)
    
    // Start chaos
    err = chaos.Create(context.Background())
    require.NoError(t, err)
    
    // Test application resilience
    t.Run("application remains functional", func(t *testing.T) {
        // Your resilience test logic here
        require.True(t, app.IsHealthy())
    })
    
    // Clean up
    chaos.Delete(context.Background())
}
```

### With WASP Load Testing

```go
func TestChaosUnderLoad(t *testing.T) {
    // Setup load test
    profile := wasp.NewProfile()
    generator := wasp.NewGenerator(&wasp.Config{
        T: 60 * time.Second,
        RPS: 100,
        LoadType: wasp.RPS,
        Schedule: wasp.Plain(100, 60*time.Second),
    })
    
    // Add load test logic
    generator.AddRequestFn(func(ctx context.Context) error {
        return performRequest(ctx)
    })
    
    profile.Add(generator)
    
    // Start chaos in background
    go func() {
        client, _ := havoc.NewClient()
        chaos, _ := havoc.NewChaos(client, createNetworkChaos())
        chaos.Create(context.Background())
        time.Sleep(30 * time.Second)
        chaos.Delete(context.Background())
    }()
    
    // Run load test
    _, err := profile.Run(true)
    require.NoError(t, err)
}
```

## Configuration

### Kubernetes Client Configuration

```go
// Use default kubeconfig
client, err := havoc.NewClient()

// Use custom kubeconfig
client, err := havoc.NewClientWithConfig(&rest.Config{
    Host: "https://kubernetes.example.com",
    BearerToken: "your-token",
    TLSClientConfig: rest.TLSClientConfig{
        Insecure: true,
    },
})
```

### Chaos Mesh Configuration

```go
// Configure Chaos Mesh namespace
client, err := havoc.NewClient()
if err != nil {
    log.Fatal(err)
}

// Set Chaos Mesh namespace
client.SetChaosMeshNamespace("chaos-mesh")

// Configure experiment defaults
client.SetDefaultDuration("30s")
client.SetDefaultMode(chaosmesh.OneMode)
```

## Best Practices

### 1. **Start Small**
- Begin with simple experiments (pod kill, network latency)
- Gradually increase complexity
- Monitor system impact carefully

### 2. **Use Appropriate Selectors**
```go
Selector: chaosmesh.PodSelectorSpec{
    Namespaces: []string{"production"},
    LabelSelectors: map[string]string{
        "app": "my-app",
        "tier": "backend",
    },
    // Use specific pod names for targeted testing
    Pods: map[string][]string{
        "production": {"my-app-1", "my-app-2"},
    },
},
```

### 3. **Monitor System Health**
```go
// Add health monitoring to chaos experiments
chaos.AddListener(&HealthMonitor{
    threshold: 0.8, // 80% health threshold
    onUnhealthy: func() {
        chaos.Pause(context.Background())
    },
})
```

### 4. **Use Appropriate Durations**
```go
Duration: &chaosmesh.Duration{
    Duration: "30s", // Short for testing
    // Duration: "5m", // Longer for stress testing
},
```

### 5. **Clean Up Properly**
```go
defer func() {
    if chaos != nil {
        chaos.Delete(context.Background())
    }
}()
```

## Troubleshooting

### Common Issues

#### Chaos Mesh Not Installed
```bash
# Check if Chaos Mesh is installed
kubectl get pods -n chaos-mesh

# Install Chaos Mesh if needed
curl -sSL https://mirrors.chaos-mesh.org/v2.6.0/crd.yaml | kubectl apply -f -
kubectl apply -f https://mirrors.chaos-mesh.org/v2.6.0/rbac.yaml
kubectl apply -f https://mirrors.chaos-mesh.org/v2.6.0/chaos-mesh.yaml
```

#### Permission Issues
```bash
# Check RBAC permissions
kubectl auth can-i create networkchaos --namespace default
kubectl auth can-i create podchaos --namespace default

# Create necessary RBAC rules
kubectl apply -f rbac.yaml
```

#### Experiment Not Working
```bash
# Check experiment status
kubectl get networkchaos -n default
kubectl describe networkchaos network-latency-test

# Check Chaos Mesh logs
kubectl logs -n chaos-mesh -l app.kubernetes.io/name=chaos-mesh
```

### Debug Mode

Enable debug logging:

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Enable debug logging
zerolog.SetGlobalLevel(zerolog.DebugLevel)
log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
```

## Examples

Check the [examples directory](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/havoc/examples) for comprehensive examples:

- [Basic chaos experiments](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/havoc/examples/basic)
- [Network chaos testing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/havoc/examples/network)
- [Pod chaos testing](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/havoc/examples/pod)
- [Integration with CTF](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/havoc/examples/ctf-integration)

## API Reference

For detailed API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/havoc).

## Contributing

We welcome contributions! Please see the [Contributing Guide](../../Contributing/Overview) for details on:

- Code of Conduct
- Development Setup
- Testing Guidelines
- Pull Request Process

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/LICENSE) file for details. 