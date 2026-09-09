---
id: grafana-alertcheck-evaluation
title: How alerts are evaluated
sidebar_label: How alerts are evaluated
sidebar_position: 1
description: The verdict model, instance timelines, and coverage proof behind grafana-alertcheck.
---

# How alerts are evaluated

Both `watch`+`check` (recorder mode) and `check` alone (single-step mode) converge on the same input: a flat list of polls. Everything below runs over that list; the mode only changes where the polls came from.

## Instance states

Grafana reports instance states in two vocabularies (`Alerting`/`Normal` at instance level, `firing`/`inactive` at rule level). The gate normalizes every instance to one canonical set:

| Canonical | Meaning |
| --------- | ------- |
| `normal`  | Healthy |
| `firing`  | The condition is true and `for` has elapsed |
| `pending` | The condition is true, `for` has not elapsed |
| `nodata`  | The query returned no series (synthetic instance) |
| `error`   | The query failed (synthetic instance) |

A rule's **rule-level** `state` and `health` are kept verbatim and only reported — they are never classified. The **instance** state is what the classifier reasons about.

A "bad" instance is one whose canonical state is in `--states` (default `firing`). `pending` and `nodata` are excluded by default.

## Verdict model

For each instance the gate builds a timeline of bad spans over `[from, to]`, then takes the worst outcome across a rule's instances as the rule's outcome.

| Outcome | Shape | Exit |
| ------- | ----- | ---- |
| `clean` | Good throughout, observed throughout | → 0 |
| `newly_bad` | Entered a bad state **inside** the window | → 1 |
| `persistently_bad` | Bad at `from`, still bad at `to` | → 1 |
| `recovered` | Bad at `from`, cleared before `to`, stayed clear | → 0 |
| `flapping` | Cleared, then became bad again | → 1 |
| `skipped` | Paused **before** the window opened | reported, not observable |
| `unobservable` | Coverage gap / sustained `health=error` / stale / absent | → 2 |

`recovered` has **no deadline** — an alert that clears at minute 58 of a 60-minute window still passes. The total bad time is reported as `BadFor`; the removed deadline is replaced by that measured value rather than a derived limit.

### Preexisting policy

For an instance already bad when `from` opened, `--preexisting` decides:

- `fail-unless-recovered` (default) — clears and stays clear → pass; never clears → fail.
- `fail` — any preexisting instance fails, recovered or not.
- `ignore` — preexisting instances are disregarded; only new episodes fail.

## Cleared vs vanished

When an instance leaves the bad set, the gate looks it up **in the same response**:

- Present as `normal` → `cleared` (a real recovery).
- Absent, or present as `normal (MissingSeries)` → `vanished` (a discontinuity, **not** a recovery).

A vanished instance that was bad stays `persistently_bad`. A metric that stops being emitted is not evidence of health — this is deliberate and can surprise users whose fix is to remove a metric rather than drive it to a good value.

## Coverage proof

Before classifying, `check` must **prove** continuous coverage of `[from, to]` for each alert. Nine checks run; any failure makes the rule `unobservable`:

1. **Sentinel** — a clean recorder stop, timestamped at or after `to + transitionGrace`. A recorder that died mid-window looks exactly like a coverage gap and is one.
2. **`from` bounds** — `from` earlier than the recording start is unprovable.
3. **Heartbeat gap** — any gap larger than `maxGap` (= 2 × poll cadence) inside the window. Data at both ends with a hole between is not enough.
4. **`health=error`** — a contiguous run longer than `healthGrace` consumes coverage; a short blip is a note.
5. **`health=nodata`** — a note, never fatal (unless `--nodata-is-unobservable`).
6. **Liveness** — `grafana_now − lastEvaluation` must not exceed `evalStaleAfter`. This is an **absolute** check, never a "did it increase since the last poll" delta.
7. **In-window pause** — a poll reporting `isPaused` mid-window is `unobservable` (the primary pause detector).
8. **Rule absent** — an authoritative `2xx` with no matching rule.
9. **`KeepLast`** — a note naming a stale-state blind spot.

## Health: `error` vs `nodata`

- `health=error` means the query **failed** — a malfunction. Sustained past `healthGrace`, it makes the rule `unobservable`.
- `health=nodata` means the query **ran and returned no series** — indistinguishable from a quiet system. It is not fatal by default; most of a fleet runs `no_data_state: OK`.

## The drain wait and `transitionGrace`

A condition that arises just before `to` becomes `firing` only at the first evaluation after its `for` elapses. `transitionGrace` (derived from the watched rules' `for` values) extends the classification bound past `to` so such a surfacing condition is caught. After collection, a **drain wait** polls until each rule has evaluated through `to + transitionGrace` (bounded by `drainTimeout`); a rule that never does is `unobservable`.

Run time = `(to − from) + transitionGrace + drainTimeout`. This is printed at start, and the grace is warned about when it exceeds a quarter of the window — the window may be too short for the alert's `for`.
