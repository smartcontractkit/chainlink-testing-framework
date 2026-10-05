---
id: grafana-alertcheck-advanced
title: Check budget and scheduling
sidebar_label: Budget and scheduling
sidebar_position: 2
description: Why grafana-alertcheck schedules per rule, how the request budget works, and why it never queries state history.
---

# Check budget and scheduling

## Per-rule schedules, never a global cycle

Each rule polls at its **own** cadence (default: half the rule's own evaluation interval). There is deliberately no single global minimum-interval cycle. Overwrite with `--poll-interval`.

One rule at `intervalSeconds=10` beside twenty at `300` keeps a 5 s cadence for itself and 150 s for the other twenty — not a 5 s cycle for all of them, which would be a 60× request bloat at ~1.8 s per request and would fail to start on a reasonable fleet.

The scheduler staggers each rule's initial next-due time across its cadence, and serves due rules **earliest-due-first**, so a tight rule never queues behind slack ones.

## The check budget

The gate records one observation of every rule up front and checks the schedule against those **measured** latencies (payload sizes varied ~230× across existing rules, so a fixed estimate would be meaningless). It errors at start — before waiting — if any of four conditions hold:

- **Utilization** — total request rate exceeds `--concurrency`.
- **Per-rule** — one rule's request can't fit its own cadence.
- **Burst bound** — the slowest request exceeds the fleet's tightest cadence, which can open a mid-run gap.
- **Startup handoff** — draining the first-observation pass's backlog at `--concurrency` would leave some rule unpolled past its own `maxGap`. A rule the pass observed early is seeded overdue, and a tight rule observed late can queue behind every rule due before it. The gate simulates the poller's first cycles from the recorded observation times and measured latencies — each wake takes every rule due at that instant, polls the batch at `--concurrency`, and wakes again when it ends — and refuses if any rule's first poll would land past its `maxGap`. Steady-state utilization cannot see this — a long pass at low concurrency is exactly the case it passes.

The error names only the levers that can fix it: the minimum `--concurrency` when the schedule is concurrency-bound, and `--poll-interval` or a smaller alert set for single-request shapes concurrency cannot shorten. It never prescribes a single interval.

## The startup pass and `ready_at`

Before detaching, `watch` observes every non-paused rule once, sequentially bounded by `--concurrency`. With many alerts and a low concurrency that pass takes real time (120 rules at ~230 ms each and the default concurrency of 1 is ~27 s). The header's `ready_at` stamps the moment the pass completed, and `check` refuses a `from` before it: a window opening inside the pass names observations that do not exist yet, and the earliest rules have no next poll until the detached recorder starts. This is a startup validation, checked from the immutable header before the wait, so a `from` emitted before `watch` returns fails immediately with a named reason instead of surfacing as a heartbeat gap mid-window. Emit `from` only after `watch` returns.

The detached recorder then continues the schedule the first observations were on (each rule's next poll is one cadence after its last recorded observation) rather than drawing fresh phases, so the handoff adds no extra up-to-one-cadence delay to the rules the pass observed first. Continuing the schedule is necessary but not sufficient: the startup-handoff budget above proves the backlog can actually be drained before any rule's `maxGap`, and refuses the run at startup when it cannot.

Single-step `check` runs the same pass itself. It cannot watch before it started, so a `from` inside the pass is a declared blind interval: the run warns, classifies from the pass completion, and the live poller continues the pass's schedule. It never classifies a window that opens before every rule has been observed.

## Datasource-managed rules: discovery and cost

Datasource-managed rules are auto-discovered — there is no selection flag. The candidate filter is strict: `type == "prometheus"` **and** `jsonData.manageAlerts == true`. The strict `true` matters: the `AlertStateHistoryBackend` datasource shares VictoriaMetrics' backend and also reports `manageAlerts` truthy, so a `!= false` filter would make every rule name ambiguous. Loki is deferred: its ruler API is broken/disabled in our Grafana, so only Prometheus-flavored rules are in scope.

Each candidate is probed with a `rule_name[]=__probe__` request; a `manageAlerts=true` source whose probe fails is a hard error naming the source, never a silently dropped source.

Cost differs by mode. `list` and label selection take the **bulk** response — one request per datasource, several MB and several seconds for a large ruler. Name selection takes a **filtered** request, ~1 KB and ~1 s. The filter uses vmalert's `[]`-suffixed parameters (`rule_name[]`, `rule_group[]`, `file[]`): vmalert reads only those and ignores plain `rule_name=`, an upstream quirk pinned by tests. `limit_alerts` is a Grafana parameter that vmalert ignores and is therefore omitted.

## Why the gate never queries state history

Querying Grafana's alert state history after the fact fails closed *in the wrong direction* — it returns "pass" when the truth is unknown:

- History stores **transitions**, not states. An alert firing through the whole window has its only record *before* the window.
- The annotations API **does not serve Loki-backed history** at all.
- An empty result is indistinguishable from a healthy one: no alert fired, the backend differs, retention removed data, the token lacked permission — all look identical.
- There is **no coverage signal** — nothing proves the history is complete to time T.
- Artifact transitions (`Paused`, `RuleDeleted`, `Updated`, `MissingSeries`) look like recoveries.

Instead, `watch` records its own evidence live and the log becomes the source of truth. The trade-off: the gate can miss an episode shorter than a rule's poll interval, though `activeAt` still surfaces sub-interval onsets for instances still active at a poll.

A corollary of recording fresh: there is no replay. Re-running a failed job is a new piece of work with a new `from` and a new recording — never a re-classification of old evidence.
