# WASP - Try It Out Quickly

---

## Prerequisites

Ensure you have the following installed:
* [Golang](https://go.dev/doc/install)
* [Docker](https://docs.docker.com/get-docker/)
* [Nix](https://nixos.org/manual/nix/stable/installation/installation.html)
* [Loki and Grafana](./start_local_observability_stack.md)
* [Pyroscope](./start_local_observability_stack.md)

---

### Install Dependencies with Nix

To install dependencies with Nix, simply execute:

```bash
nix develop
```

---

## Running Grafana and Loki Tests

Assuming you have already started and configured the local observability stack, you can run sample Loki and Grafana tests with:

```bash
make test_loki
```

Once the test is running, you can view the results in the [dashboard](http://localhost:3000/d/wasp/Wasp?orgId=1&refresh=5s&from=now-5m&to=now).

> [!WARNING]  
> If deploying to your own Grafana instance, verify the `DASHBOARD_FOLDER` and `DASHBOARD_NAME`.  
> Defaults are the `LoadTests` directory, and the dashboard is named `Wasp`.

---

## Running Pyroscope Tests

If Grafana, Loki, and Pyroscope are already running, you can execute sample Pyroscope tests with:

```bash
make test_pyro_rps
```

or

```bash
make test_pyro_vu
```

During the test, you can view the results in the [Pyroscope dashboard](http://localhost:4040/).

---

### Additional Debugging Option

You can also add the `-trace trace.out` flag when running any of your tests with `go test` for additional tracing information.