---
id: grafana-alertcheck-cli
title: CLI reference
sidebar_label: CLI reference
sidebar_position: 0
description: "Full reference for the grafana-alertcheck CLI: watch, check, list, environment, naming, and output."
---

# CLI reference

```
grafana-alertcheck <list|watch|check|stop>
```

Connection details are always from the environment: `GRAFANA_URL` and `GRAFANA_TOKEN`. The token is never a flag and never logged.

## `list`

Lists every rule — kind, datasource, folder, group, title, key, uid. Grafana-managed rules come from the ruler endpoint; datasource-managed rules are auto-discovered per datasource through the Prometheus API. Useful to check auth and to find `uid:`/`key:` names. The bulk datasource fetch can take seconds per source.

```bash
grafana-alertcheck list
```

## `watch` — record

```bash
grafana-alertcheck watch --out <file> [--pidfile F] [--daemon-log F] \
  (--alerts <file|-> [--folder F] | --include-labels k=v,... [--exclude-labels k=v,...]) \
  [--exclude-alerts <file|->] [--concurrency N] [--until RFC3339]
```

| Flag | Default | Meaning |
| ---- | ------- | ------- |
| `--out` | — | JSONL log path (required) |
| `--pidfile` | `<out>.pid` | Where the recorder's pid is written |
| `--daemon-log` | `<out>.daemon.log` | stdout/stderr sink for the detached recorder |
| `--alerts` | — | File of alert names, one per line, or `-` for stdin (required unless `--include-labels`) |
| `--folder` | — | Default folder to scope unqualified names (with `--alerts` only) |
| `--include-labels` | — | Comma-separated exact-match `key=value` pairs selecting rules by label (cannot be combined with `--alerts`) |
| `--exclude-labels` | — | Comma-separated exact-match `key=value` pairs; a rule carrying any of them is dropped (requires `--include-labels`) |
| `--exclude-alerts` | — | File of alert names, one per line, or `-` for stdin; subtracted from the selected set (works with `--alerts` and with labels) |
| `--concurrency` | `1` | Max concurrent requests to Grafana |
| `--until` | run until signalled | Optional hard stop |

`watch` observes every non-paused rule once, checks the budget, writes the header (with `ready_at` stamped once the observation pass completes) and those observations, then detaches a background recorder and returns. Recording is **unfiltered** — there is no `--states` here, so the same log can be re-classified later under different `--states` without re-recording.

## `stop` — reap the recorder

```bash
grafana-alertcheck stop --out <file> [--pidfile F]
```

| Flag | Default | Meaning |
| ---- | ------- | ------- |
| `--out` | — | JSONL log path whose recorder to stop (required) |
| `--pidfile` | `<out>.pid` | Pidfile of the recorder to stop |

Stops a detached recorder that is still running, or confirms it has already finished. It reads the pidfile, then asks the log's **flock** whether a writer exists right now (the lock is authoritative; a pid can be reused): if a writer is alive it is sent `SIGTERM`, and if it ignores that it is killed, then the pidfile is removed.

It is **idempotent** — after `check` has already stopped the recorder, or after a previous `stop`, it reports that there is nothing to stop and exits `0`. That is what lets an `if: always()` step call it on both the success and failure paths.

Use it when the work failed and the alert verdict no longer matters, but the recorder must still be reaped: the recorder is detached in its own session, so neither `check` nor the runner's cleanup will stop it, and it would keep polling Grafana until its window elapsed.

## `check` — classify

```bash
grafana-alertcheck check [--in <file>] [--pidfile F] --from RFC3339 --to RFC3339 \
  [--alerts ... [--folder F] | --include-labels k=v,... [--exclude-labels k=v,...]] [--exclude-alerts ...] \
  [--states ...] [--preexisting ...] [--min-observed N] \
  [--allow-paused] [--nodata-is-unobservable] [--fail-fast=false] [--concurrency N] [--output json]
```

| Flag | Default | Meaning |
| ---- | ------- | ------- |
| `--in` | — | Log recorded by `watch`; empty selects single-step mode |
| `--pidfile` | `<in>.pid` | Recorder to stop before reading `--in` |
| `--from` | see below | Moment the deploy finished |
| `--to` | — | End of the window (required) |
| `--alerts` | — | Required **without** `--in` (unless `--include-labels`); refused **with** `--in` |
| `--folder` | — | Default folder to scope unqualified names (with `--alerts` only) |
| `--include-labels` | — | Comma-separated exact-match `key=value` pairs selecting rules by label (cannot be combined with `--alerts`) |
| `--exclude-labels` | — | Comma-separated exact-match `key=value` pairs; a rule carrying any of them is dropped (requires `--include-labels`) |
| `--exclude-alerts` | — | File of alert names, one per line, or `-` for stdin; subtracted from the selected set (works with `--alerts` and with labels; refused **with** `--in`) |
| `--states` | `firing,recovering` | Comma-separated bad states: `firing,pending,recovering,nodata,error` |
| `--preexisting` | `fail-unless-recovered` | `fail-unless-recovered` \| `fail` \| `ignore` |
| `--min-observed` | every resolved rule | Minimum rules that must be observed |
| `--allow-paused` | `false` | Don't count pre-window-paused rules against `--min-observed` |
| `--nodata-is-unobservable` | `false` | Treat sustained `health=nodata` as unobservable |
| `--fail-fast` | `true` | Stop as soon as a failure that cannot become a pass is observed; `--fail-fast=false` waits for the full window and its coverage proof |
| `--concurrency` | `1` | Max concurrent requests |
| `--output` | `table` | `json` also writes the machine-readable result to stdout |

By default `check` **exits early** on a failure that cannot become a pass: a post-`from` bad onset, or an inability (a heartbeat gap, a sustained `health=error`, a stale evaluation, an in-window pause, an absent rule). This is a latency optimization, not a weaker gate — it never exits `0` early. The one observable difference is that an early exit can report `1` where a full run would have discovered an inability later and reported `2`. `--fail-fast=false` always waits for `to + transitionGrace` and the full coverage proof; the `Result` then carries no `terminated_early` marker. With early exit the JSON result includes `terminated_early` naming the rule, kind, reason and time.

`--from` and `--to` are RFC3339 with an explicit offset and must come from your work — `from` from the deploy step, `to` from the step that finishes. In recorder mode an absent `--from` is a hard error; in single-step mode it falls back (with a warning) to the start of the step. A window with a subsecond part is rounded up to the next whole second by extending `to`, so the plan never reads `window 9m59.99445781s`.

## Naming alerts

Alert names take one of these forms. Grafana-managed rules use folder/group; datasource-managed rules are auto-discovered (no selection flag) and use datasource/group.

| Form | Meaning |
| ---- | ------- |
| `HighErrorRate` | Title only; scoped by `--folder` for Grafana rules |
| `Platform/HighErrorRate` | Grafana folder + title |
| `Platform/api/HighErrorRate` | Grafana folder + group + title |
| `ExampleMetrics/HighErrorRate` | Datasource group + title |
| `VM Prod/ExampleMetrics/HighErrorRate` | Datasource name + group + title |
| `uid:abc123` | Exact Grafana uid |
| `key:ds:[…]` | Exact rule key across both kinds (copyable from `list`) |

A datasource rule's **name can itself contain `/`** (e.g. `devex-cicd/prod/griddle-github: ContainersNotReady`). The exact name is tried first, so the `TITLE` from `list` always resolves, and `key:` is the unambiguous fallback.

`--folder` scopes a bare Grafana title only. A recording rule, or a datasource rule with no identifiable datasource, is refused with a specific error; a no-match points at `list`; an ambiguous name lists every candidate with its full name, source and `uid:`/`key:`. Duplicate names collapse to one (a note, not an error).

Auto-discovery keeps `/api/datasources` entries with `type == "prometheus"` and `jsonData.manageAlerts == true`, then probes each. The token needs `datasources:read` plus datasource query permission; a failure names the permission.

## Selecting alerts by labels

Instead of naming alerts, `watch` and single-step `check` accept a label selection:

```bash
grafana-alertcheck watch --out /tmp/run.jsonl --include-labels team=bcm,env=stage
grafana-alertcheck check --to "$finished_at" --include-labels team=bcm --exclude-labels severity=info
```

`--include-labels` takes comma-separated exact-match `key=value` pairs; a rule must carry **all** of them. `--exclude-labels` is optional and drops any rule carrying **one** of its pairs. A rule that does not carry the label is never dropped, only never included. Values cannot contain commas; `key=` matches only rules that carry the label with an empty value.

`--exclude-alerts` is an enumerated list (a file, or `-` for stdin) subtracted from whichever set was selected — names or labels. It is the escape hatch for a rule that carries the include labels but must not be watched.

The label flags cannot be combined with `--alerts` or `--folder`, and they are refused with `--in` — the recorded log names its own alert set. A selection that matches no rules, whose matches are all excluded, or that matches a recording rule exits `2`: an empty watch set must never pass. When rules are selected by labels, the resolved set is printed one rule per line before the run starts, so an operator can see exactly what matched.

## Output and exit codes

The human table goes to **stderr**: `RESULTS` (one row per rule: verdict, time broken, check cadence, whether the window was observed, and `SOURCE` — `grafana` or `datasource`), `VIOLATIONS` (one per distinct rule/verdict/state/health/note signature, with an `INSTANCES` count — instance identity is only in the JSON), and `LIMITS USED` (each rule's observation limits in plain words, explained by a legend, plus the extra observation time, the evaluation wait, the largest measured clock difference and the Grafana version; the closing violations count is marked ✅/❌).

`DETAILS` carries only rule-specific notes; kind-level caveats (datasource rules have no pause signal and treat a departure as a recovery) are printed once, before the table. The JSON outcome values are `healthy`, `new_failure`, `still_failing`, `recovered`, `unstable`, `paused`, `not_verified` and the synthetic `not_counted`; `--output json` adds each rule's `source_kind` and the run-level `caveats`, and writes the result to stdout.

| Code | Meaning |
| ---- | ------- |
| `0`  | Pass |
| `1`  | Violations |
| `2`  | Could not check — every library error, never a pass |
