# WASP - Configuration

WASP can be configured using environment variables for most commonly used settings. However, to fully leverage the flexibility of WASP, you may need to use programmatic configuration for more advanced features.

### Required Environment Variables

At a minimum, you need to provide the following environment variables, as WASP requires Loki:

* `LOKI_URL` - The Loki **endpoint** to which logs are pushed (e.g., [http://localhost:3100/loki/api/v1/push](http://localhost:3100/loki/api/v1/push)).
* `LOKI_TOKEN` - The authorization token.

Optionally, you can also provide the following:
* `LOKI_TENANT_ID` - A tenant ID that acts as a bucket identifier for logs, logically separating them from other sets of logs. If the tenant ID doesn't exist, Loki will create it. Can be empty if log separation is not a concern.
* `LOKI_BASIC_AUTH` -  Basic authentication credentials.

---

### Alert Configuration

To enable alert checking, you need to provide the following additional environment variables, as alerts are an integral part of Grafana:

* `GRAFANA_URL` - The base URL of the Grafana instance.
* `GRAFANA_TOKEN` - An API token with permissions to access the following namespaces:  
  `/api/alertmanager/`, `/api/annotations/`, `/api/dashboard/`, `/api/datasources/`, `/api/v1/provisioning/`, and `/api/ruler/`.

---

### Grafana Dashboard Creation

If you want WASP to create a Grafana dashboard for you, provide the following environment variables:

* `DATA_SOURCE_NAME` - The name of the data source (currently, only `Loki` is supported).
* `DASHBOARD_FOLDER` - The folder in which to create the dashboard.
* `DASHBOARD_NAME` - The name of the dashboard.

---

### Log Level Control

You can control the log level using this environment variable:

* `WASP_LOG_LEVEL` - Sets the log level (`trace`, `debug`, `info`, `warn`, `error`; defaults to `info`).

---

And that's it! You are now ready to start using WASP to its full potential.

> [!NOTE]  
> Remember, WASP offers much more configurability through its programmatic API.
