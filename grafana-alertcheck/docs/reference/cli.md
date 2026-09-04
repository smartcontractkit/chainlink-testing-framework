---
id: grafana-alertcheck-cli
title: CLI reference
sidebar_label: CLI reference
sidebar_position: 0
description: Full reference for the grafana-alertcheck CLI: watch, check, list, environment, naming, and output.
---

# CLI reference

```
grafana-alertcheck <list|watch|check>
```

Connection details are always from the environment: `GRAFANA_URL` and `GRAFANA_TOKEN`. The token is never a flag and never logged.

## `list`

Lists every rule from the ruler endpoint — kind, folder, group, title, uid. Useful to check auth and to find `uid:` names.

```bash
grafana-alertcheck list
```

## `watch` — record

```bash
grafana-alertcheck watch --out <file> [--pidfile F] [--daemon-log F] \
  --alerts <file|-> [--folder F] [--poll-interval D] [--concurrency N] [--until RFC3339]
```

| Flag | Default | Meaning |
| ---- | ------- | ------- |
| `--out` | — | JSONL log path (required) |
| `--pidfile` | `<out>.pid` | Where the recorder's pid is written |
| `--daemon-log` | `<out>.daemon.log` | stdout/stderr sink for the detached recorder |
| `--alerts` | — | File of alert names, one per line, or `-` for stdin (required) |
| `--folder` | — | Default folder to scope unqualified names |
| `--poll-interval` | half the rule's interval | Override every rule's cadence (never clamped) |
| `--concurrency` | `1` | Max concurrent requests to Grafana |
| `--until` | run until signalled | Optional hard stop |

`watch` writes the header, observes every non-paused rule once, checks the budget, then detaches a background recorder and returns. Recording is **unfiltered** — there is no `--states` here, so the same log can be re-classified later under different `--states` without re-recording.

## `check` — classify

```bash
grafana-alertcheck check [--in <file>] [--pidfile F] --from RFC3339 --to RFC3339 \
  [--alerts ...] [--folder F] [--states ...] [--preexisting ...] [--min-observed N] \
  [--allow-paused] [--nodata-is-unobservable] [--concurrency N] [--output json]
```

| Flag | Default | Meaning |
| ---- | ------- | ------- |
| `--in` | — | Log recorded by `watch`; empty selects single-step mode |
| `--pidfile` | `<in>.pid` | Recorder to stop before reading `--in` |
| `--from` | see below | Moment the deploy finished |
| `--to` | — | End of the window (required) |
| `--alerts` | — | Required **without** `--in`; refused **with** `--in` |
| `--states` | `firing` | Comma-separated bad states: `firing,pending,nodata,error` |
| `--preexisting` | `fail-unless-recovered` | `fail-unless-recovered` \| `fail` \| `ignore` |
| `--min-observed` | every resolved rule | Minimum rules that must be observed |
| `--allow-paused` | `false` | Don't count pre-window-paused rules against `--min-observed` |
| `--nodata-is-unobservable` | `false` | Treat sustained `health=nodata` as unobservable |
| `--concurrency` | `1` | Max concurrent requests |
| `--output` | `table` | `json` also writes the machine-readable result to stdout |

`--from` and `--to` are RFC3339 with an explicit offset and must come from your work — `from` from the deploy step, `to` from the step that finishes. In recorder mode an absent `--from` is a hard error; in single-step mode it falls back (with a warning) to the start of the step.

## Naming alerts

Alert names take one of four forms:

| Form | Meaning |
| ---- | ------- |
| `HighErrorRate` | Title only, scoped by `--folder` |
| `Platform/HighErrorRate` | Folder + title |
| `Platform/api/HighErrorRate` | Folder + group + title (always unique) |
| `uid:abc123` | Exact uid (present on both endpoints) |

Datasource-managed and recording rules are refused with a specific error. A name matching multiple rules errors listing every candidate with the copyable `Folder/Group/Title` and its `uid:` form. A no-match errors with case-insensitive substring suggestions and points at `list`. Duplicate names that resolve to the same uid collapse to one (a note, not an error).

## Output and exit codes

The human table goes to **stderr**: `RESULTS` (one row per rule), `VIOLATIONS` (one per violation), and `THRESHOLDS` (each rule's `maxGap`/`healthGrace`/`evalStaleAfter` plus global `transitionGrace`/`drainTimeout` and the largest measured clock skew). `--output json` writes the result to stdout.

| Code | Meaning |
| ---- | ------- |
| `0`  | Pass |
| `1`  | Violations |
| `2`  | Could not check — every library error, never a pass |
