# WASP - Alert Checker

> [!WARNING]  
> The API used by `AlertChecker` is unstable, and this section may be out of date.

`AlertChecker` is a simple yet powerful component that allows you to verify which alerts fired during test execution.

It supports the following functionalities:
* Checking if **any alert** fired during a specified time range for a given dashboard.
* Checking if **any group of alerts** fired for a given dashboard.

The first mode is more suitable for existing dashboards, while the second mode is ideal for dashboards created specifically for the test (as it does not support time range selection).

> [!NOTE]  
> To use this component, you need to set certain Grafana-specific variables as described in the [Configuration](../configuration.md) section.

For a practical example of how to use `AlertChecker`, refer to the [Testing Alerts](../testing_alerts.md) section.

> [!WARNING]  
> If you define alerts yourself, ensure they are grouped by adding a label named `requirement_name` with a value representing the alert group.
