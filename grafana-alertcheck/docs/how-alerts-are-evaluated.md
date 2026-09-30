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
| `healthy` | Good throughout, observed throughout | → 0 |
| `new_failure` | Entered a bad state **inside** the window | → 1 |
| `still_failing` | Bad at `from`, still bad at `to` | → 1 |
| `recovered` | Bad at `from`, cleared before `to`, stayed clear | → 0 |
| `unstable` | Cleared, then became bad again | → 1 |
| `paused` | Paused **before** the window opened | counts against `--min-observed` unless `--allow-paused` |
| `not_verified` | The window could not be observed: a gap, sustained `health=error`, a stale evaluation, or an absent rule | → 2 |

A `--min-observed` deficit that no rule explains is reported as `not_counted`: it is not a verdict on any alert.

`recovered` has **no deadline** — an alert that clears at minute 58 of a 60-minute window still passes. The total bad time is reported as `BadFor`.

### Preexisting policy

For an instance already bad when `from` opened, `--preexisting` decides:

- `fail-unless-recovered` (default) — clears and stays clear → pass; never clears → fail.
- `fail` — any preexisting instance fails, recovered or not.
- `ignore` — preexisting instances are disregarded; only new episodes fail.

## Cleared vs vanished

When an instance leaves the bad set, the gate looks it up **in the same response**:

- Present as `normal` → `cleared` (a real recovery).
- Absent, or present as `normal (MissingSeries)` → `vanished` (a discontinuity, **not** a recovery).

A vanished instance that was bad stays `still_failing`. A metric that stops being emitted is not evidence of health — this is deliberate and can surprise users whose fix is to remove a metric rather than drive it to a good value.

## Coverage proof

Before classifying, `check` must **prove** continuous coverage of `[from, to]` for each alert. Nine checks run; any failure makes the rule `not_verified`:

1. **Sentinel** — a clean recorder stop, timestamped at or after `to + transitionGrace`. A recorder that died mid-window looks exactly like a coverage gap and is one.
2. **`from` bounds** — `from` earlier than the recording start is unprovable.
3. **Heartbeat gap** — any gap larger than `maxGap` (= 2 × poll cadence) inside the window. Data at both ends with a hole between is not enough.
4. **`health=error`** — a contiguous run longer than `healthGrace` consumes coverage; a short blip is a note.
5. **`health=nodata`** — a note, never fatal (unless `--nodata-is-unobservable`).
6. **Liveness** — `grafana_now − lastEvaluation` must not exceed `evalStaleAfter`. This is an **absolute** check, never a "did it increase since the last poll" delta.
7. **In-window pause** — a poll reporting `isPaused` mid-window is `not_verified` (the primary pause detector).
8. **Rule absent** — an authoritative `2xx` with no matching rule.
9. **`KeepLast`** — a note naming a stale-state blind spot.

## Health: `error` vs `nodata`

- `health=error` means the query **failed** — a malfunction. Sustained past `healthGrace`, it makes the rule `not_verified`.
- `health=nodata` means the query **ran and returned no series** — indistinguishable from a quiet system. It is not fatal by default; most of a fleet runs `no_data_state: OK`.

## Early exit (fail-fast)

`check` does not have to wait for the whole window to know the run has failed. As soon as it observes a condition that cannot become a pass, it stops and classifies the sub-window it did see:

- a **post-`from` bad onset** — the full classifier would call it `new_failure` (or `unstable`), which fails whether or not it later clears; or
- an **inability** — a heartbeat gap, a sustained `health=error` run, a stale evaluation, an in-window pause, or an absent rule.

A preexisting bad instance is deliberately **not** terminal: if it clears before `to` the full run would call it `recovered`, which passes.

Fail-fast is on by default and always preserves the failure: an early run can exit `1` or `2`, never `0`. The one difference from a full run is that an early exit may report `1` before an inability surfaces that would have made it `2`. `--fail-fast=false` disables the guard and always waits for the full window and its coverage proof.

## The drain wait and `transitionGrace`

A condition that arises just before `to` becomes `firing` only at the first evaluation after its `for` elapses. `transitionGrace` (derived from the watched rules' `for` values) extends the classification bound past `to` so such a surfacing condition is caught. After collection, a **drain wait** polls until each rule has evaluated through `to + transitionGrace` (bounded by `drainTimeout`); a rule that never does is `not_verified`.

Run time = `(to − from) + transitionGrace + drainTimeout`. This is printed at start. A requested window with a subsecond part is rounded up to the next whole second by extending `to`, so the plan never reads a window like `9m59.99445781s`.
