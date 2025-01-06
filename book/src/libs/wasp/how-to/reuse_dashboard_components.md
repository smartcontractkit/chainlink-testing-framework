# WASP - How to Reuse WASP Components in Your Own Dashboard

You can reuse components from the default WASP dashboard to build your custom dashboard. Here’s an example:

```go
import (
    waspdashboard "github.com/smartcontractkit/wasp/dashboard"
    "github.com/K-Phoen/grabana/dashboard"
)

func BuildCustomLoadTestDashboard(dashboardName string) (dashboard.Builder, error) {
    // Custom key-value pairs used to query panels
    panelQuery := map[string]string{
        "branch":       `=~"${branch:pipe}"`,
        "commit":       `=~"${commit:pipe}"`,
        "network_type": `="testnet"`,
    }

    return dashboard.New(
        dashboardName,
        waspdashboard.WASPLoadStatsRow("Loki", panelQuery),      // WASP component
        waspdashboard.WASPDebugDataRow("Loki", panelQuery, true), // WASP component
        // Your custom panels and rows go here
    )
}
```

---

### Available Components

You can find all reusable components in the following files:
* [grafanasdk/panels.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/dashboard/grafanasdk/panels.go)
* [panels.go](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/dashboard/panels.go)

---

### Key Points

- WASP components like `WASPLoadStatsRow` and `WASPDebugDataRow` can be directly included in your dashboard.
- Customize panel queries using key-value pairs specific to your setup.
- Extend the default dashboard by adding your own rows and panels alongside WASP components.