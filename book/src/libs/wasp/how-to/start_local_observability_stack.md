# WASP - How to Start Local Observability Stack

To execute all examples or tests locally, you need to configure a local observability stack. We have 2 versions of popular stacks: LGTM (Grafana, Loki, Prometheus) and VictoriaMetrics + OTEL (Grafana, VictoriaMetrics, VictoriaLogs, OTEL).


---

## LGTM

```bash
just lgtm-up
just lgtm-down
```

* [Dashboard](http://localhost:3000/d/wasp-victorialogs/wasp-victorialogs?orgId=1&from=now-5m&to=now&timezone=browser&var-go_test_name=$__all&var-gen_name=$__all&var-call_group=$__all&var-branch=$__all&var-commit=$__all&refresh=5s)

Don't forget to remove the stack when you are done!

## VictoriaMetrics + OTEL

```bash
just victoria-up
just victoria-down
```

* [Dashboard](http://localhost:3000/d/wasp-victorialogs/wasp-victorialogs?orgId=1&from=now-5m&to=now&timezone=browser&var-go_test_name=$__all&var-gen_name=$__all&var-call_group=$__all&var-branch=$__all&var-commit=$__all&refresh=5s)
