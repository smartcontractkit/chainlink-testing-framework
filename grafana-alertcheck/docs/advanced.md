---
id: grafana-alertcheck-advanced
title: Check budget and scheduling
sidebar_label: Budget and scheduling
sidebar_position: 2
description: Why grafana-alertcheck schedules per rule, how the request budget works, and why it never queries state history.
---

# Check budget and scheduling

## Per-rule schedules, never a global cycle

Each rule polls at its **own** cadence, `--poll-interval` (default: half the rule's own evaluation interval). There is deliberately no single global minimum-interval cycle.

One rule at `intervalSeconds=10` beside twenty at `300` keeps a 5 s cadence for itself and 150 s for the other twenty — not a 5 s cycle for all of them, which would be a 60× request bloat at ~1.8 s per request and would fail to start on a reasonable fleet.

The scheduler staggers each rule's initial next-due time across its cadence, and serves due rules **earliest-due-first**, so a tight rule never queues behind slack ones.

## The check budget

The gate records one observation of every rule up front and checks the schedule against those **measured** latencies (payload sizes vary ~230× across rules, so a fixed estimate is meaningless). It errors at start — before waiting — if any of three conditions hold:

- **Utilization** — total request rate exceeds `--concurrency`.
- **Per-rule** — one rule's request can't fit its own cadence.
- **Burst bound** — the slowest request exceeds the fleet's tightest cadence, which can open a mid-run gap.

The error names the three levers only: raise `--concurrency`, raise `--poll-interval`, or watch fewer alerts. It never prescribes a single interval.

## Why the gate never queries state history

Querying Grafana's alert state history after the fact fails closed *in the wrong direction* — it returns "pass" when the truth is unknown:

- History stores **transitions**, not states. An alert firing through the whole window has its only record *before* the window.
- The annotations API **does not serve Loki-backed history** at all.
- An empty result is indistinguishable from a healthy one: no alert fired, the backend differs, retention removed data, the token lacked permission — all look identical.
- There is **no coverage signal** — nothing proves the history is complete to time T.
- Artifact transitions (`Paused`, `RuleDeleted`, `Updated`, `MissingSeries`) look like recoveries.

Instead, `watch` records its own evidence live and the log becomes the source of truth. The trade-off: the gate can miss an episode shorter than a rule's poll interval, though `activeAt` still surfaces sub-interval onsets for instances still active at a poll.

A corollary of recording fresh: there is no replay. Re-running a failed job is a new deploy with a new `from` and a new recording — never a re-classification of old evidence.
