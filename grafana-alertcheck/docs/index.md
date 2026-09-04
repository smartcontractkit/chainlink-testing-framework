---
id: grafana-alertcheck-index
title: Grafana Alertcheck
sidebar_label: Overview
sidebar_position: 0
description: A CD quality gate that bookends a release with alert-state observation and answers whether any watched Grafana alert was bad during the release window.
---

# Grafana Alertcheck

`grafana-alertcheck` is a CD quality gate for Grafana alerts. It bookends a release with two commands — `watch` (record) and `check` (classify) — and answers one question:

> A release finished at time T. Was any of these Grafana alerts in a bad state during the next N minutes?

The contract is `watch → your work → check`. Between the two you run whatever you want (deploy, tests, migration); the gate only observes, then classifies.

It **fails closed**: if it cannot get an answer, it stops the release. It never passes an unproven window.

## How it works, in one paragraph

`watch` starts a background recorder that polls each named alert and appends snapshots to a JSONL log. Your work then emits two RFC3339 timestamps — `from` (when the change landed) and `to` (when the work ended). `check` proves continuous coverage of `[from, to]`, builds a state timeline per alert, classifies it, and exits `0`, `1`, or `2`.

## Install

```bash
go install github.com/smartcontractkit/chainlink-testing-framework/grafana-alertcheck/cmd/grafana-alertcheck@latest
```

Connection details come from the environment — the token is env-only, never a flag:

```bash
export GRAFANA_URL=https://grafana.example.com
export GRAFANA_TOKEN=…
```

Requires Grafana >= 13.0.0 and < 14.0.0. Outside that range the gate exits `2`.

## Quickstart — recorder mode

```bash
grafana-alertcheck watch --out /tmp/run.jsonl --alerts alerts.txt
./deploy.sh   # emits deployed_at=<RFC3339> when the rollout is stable
./verify.sh   # emits finished_at=<RFC3339> when the work is done
grafana-alertcheck check --in /tmp/run.jsonl --from "$deployed_at" --to "$finished_at"
```

`alerts.txt` holds one alert name per line. See [Naming alerts](./reference/cli#naming-alerts).

`watch` returns only after the recorder has observed every named, non-paused alert once and reported ready — so auth, name-resolution, and parse failures surface **before** your deploy runs.

## Quickstart — single-step mode

Skip the recorder and observe the window inline, from inside `check` itself:

```bash
grafana-alertcheck check --alerts alerts.txt --to "$finished_at"
```

In single-step mode the window starts at `check`'s first observation; if you give no `--from`, the interval before that first observation is declared as a blind spot with a warning (not an error).

## Exit codes

| Code | Meaning |
| ---- | ------- |
| `0`  | Pass — no violations |
| `1`  | Violations (including a paused-rule-only `--min-observed` shortfall) |
| `2`  | The gate could not check — config, auth, resolution, a coverage gap, health/staleness, the drain limit, transport, … |

An error is never a pass: `2` wins over any violation found alongside it.

## Common surprises

- A **paused rule fails by default** — even one someone else paused. Use `--allow-paused`.
- A fix that **stops emitting a metric is not a recovery** — the instance vanishes, which is a discontinuity, not health.
- The gate checks alert **state and health**, not notification delivery — a silenced alert that still fires fails.
- `recovered` has **no deadline** — a bad-at-`from` alert that clears by `to` passes; set `--preexisting fail` to forbid it.
- A **retry is a new deploy**, not a replay — re-running the job re-records against a new `from`.

## More

- [How alerts are evaluated](./how-alerts-are-evaluated) — the verdict model and coverage proof
- [Check budget and scheduling](./advanced) — why the schedule and budget look the way they do, and why history isn't queried
- [CLI reference](./reference/cli) — every subcommand and flag
