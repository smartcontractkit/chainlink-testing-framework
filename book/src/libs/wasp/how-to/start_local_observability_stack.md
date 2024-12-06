# WASP - How to Start Local Observability Stack

To execute all examples or tests locally, you need to configure a local observability stack. This stack includes **Loki**, **Grafana**, and **Pyroscope**.

---

## Prerequisites

Ensure you have the following installed:
* [Docker](https://docs.docker.com/get-docker/)

---

## Grafana and Loki

To start the local observability stack, run the following command:

```bash
make start
```

This command will download and run Docker images for Loki and Grafana. Once completed, it will output something like:

```bash
Service account id: 2
Grafana token: "<grafana token>"
```

### Setting Up Environment Variables

Next, set up the required environment variables.
> [!WARNING]  
> Replace `<Grafana token>` with the token provided in the previous step.

```bash
export LOKI_TOKEN=
export LOKI_URL=http://localhost:3030/loki/api/v1/push
export GRAFANA_URL=http://localhost:3000
export GRAFANA_TOKEN=<Grafana token>
export DATA_SOURCE_NAME=Loki
export DASHBOARD_FOLDER=LoadTests
export DASHBOARD_NAME=Wasp
```

### Accessing Services

* Grafana: [http://localhost:3000/](http://localhost:3000/)
* Loki: [http://localhost:3030/](http://localhost:3030/)

### Stopping the Containers

To stop both containers, run:

```bash
make stop
```

---

## Pyroscope

To start Pyroscope, execute:

```bash
make pyro_start
```

> [!NOTE]  
> Pyroscope is available at: [http://localhost:4040/](http://localhost:4040/)

### Stopping Pyroscope

To stop Pyroscope, run:

```bash
make pyro_stop
```