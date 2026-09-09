# grafana-alertcheck

A CD quality gate for Grafana alerts. It bookends a release with two commands — `watch` (record) and
`check` (classify) — and answers whether any watched alert was in a bad state during the release window.

```
watch  → your work → check
```

`watch` starts a background recorder that polls each named alert into a JSONL log. After the work emits a
`from`/`to` pair, `check` proves continuous coverage of that window, classifies each alert's state
timeline, and exits `0`, `1`, or `2`.

It **fails closed**: if it cannot get an answer, it stops the release — never a pass on an unproven window.

## Quickstart

```bash
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=…

grafana-alertcheck watch --out /tmp/run.jsonl --alerts alerts.txt
./deploy.sh   # emits deployed_at=<RFC3339>
./verify.sh   # emits finished_at=<RFC3339>
grafana-alertcheck check --in /tmp/run.jsonl --from "$deployed_at" --to "$finished_at"
```

Requires Grafana >= 13.0.0 and < 14.0.0. Connection details come from the environment only — the token is
never a flag.

## Documentation

| Doc | Covers |
| --- | ------ |
| [`docs/index.md`](./docs/index.md) | Overview, quickstarts, exit codes, common surprises |
| [`docs/how-alerts-are-evaluated.md`](./docs/how-alerts-are-evaluated.md) | Verdict model, coverage proof, health/liveness |
| [`docs/advanced.md`](./docs/advanced.md) | Check budget, scheduling, why history isn't queried |
| [`docs/architecture.md`](./docs/architecture.md) | Design invariants, the pure-function seam, recorder lifecycle |
| [`docs/reference/cli.md`](./docs/reference/cli.md) | Full CLI reference — subcommands, flags, naming |
| [`docs/reference/log-format.md`](./docs/reference/log-format.md) | The JSONL log schema, for debugging artifacts |

## Build

```bash
go build ./... && go test ./...
```
