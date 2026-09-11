---
id: grafana-alertcheck-architecture
title: Architecture
sidebar_label: Architecture
sidebar_position: 3
description: The design invariants, pure-function seam, and recorder lifecycle of grafana-alertcheck, for maintainers.
---

# Architecture

This page documents the invariants and seams a maintainer must not break. It exists because most of them are the difference between a gate that fails closed and one that silently passes broken windows.

## Fail-closed invariants

The gate must fail if it cannot get an answer. Every rule below is a specific instance of that:

- **An error is never a pass.** A pass is exactly `len(Violations) == 0 && err == nil`. Every error path leaves `err` non-nil, and the CLI maps that to exit `2` unconditionally.
- **Inability beats violation.** Any `unobservable` rule is exit `2`, even alongside a real violation found first.
- **Absent never means normal.** An instance that leaves the bad set is looked up in the *same* response: present as `normal` → cleared; absent (or `MissingSeries`) → vanished (a discontinuity, not a recovery).
- **Staleness is absolute.** `grafana_now − lastEvaluation` is compared against a threshold, never "did it increase since the last poll" — a delta check reports stale on ~half the polls of a healthy rule (we poll at half of `intervalSeconds` of each rule).
- **`grafana_now` is the response `Date` header.** Never the runner clock, in any comparison against a Grafana timestamp.
- **No early exit.** `check` collects to `to + transitionGrace` before classifying once.
- **No replay.** No run-id key, no artifact download, no state between attempts. A retry is a new piece of work and observation.

## The pure-function seam

All correctness lives in two phases written as **pure functions** over a flat list of polls — no HTTP, no files, no clock, no goroutines:

```
HTTP ──> Source ──> []StateRule ──> reduce ──> []Poll ──> proveCoverage ──> decide ──> Result
                            │
                       JSONL log ──> ReadLog ──┘
```

- `proveCoverage` (the nine coverage checks) and `decide` (the instance timelines and outcomes) are pure; tests drive them with `[]Poll` literals and a fake `Clock`, with no sleeping or fixture server.
- `Check`/`Watch` are I/O shells: HTTP, signals, the pidfile, file reads, the countdown print. The only test doubles needed are the `Source` and `Clock` interfaces.
- `Policy` is the narrowed view of `Config` that reaches the pure layer — classification knobs and the window, no URL and no token. The token must never cross that line, which is the cheapest guarantee it never lands in an error string or a result.

## Strict parsing as the version guard

Both API responses are parsed strictly: a missing or unparseable **required** field (`health`, `state`, `lastEvaluation`, `interval`) is an error, never a zero value. Optional keys (`alerts`, `totals`, `labels`, `keepFiringFor`) are absent-tolerant, and unknown keys are ignored — so Grafana can add fields without breaking the parser, but removing one fails loudly.

This, plus the declared supported range (Grafana >= 13.0.0, < 14.0.0), is how a deprecation or schema change is caught instead of silently misread.

## The recorder lifecycle

`watch` detaches a background recorder so observation survives the step boundary:

1. Parent resolves names, writes the header, observes every non-paused rule once, checks the budget.
2. Parent re-execs itself as the child (`--daemon-child`) under a new session/process group, stdout/stderr to the daemon log.
3. Child re-reads the header, reopens the log `O_APPEND`, takes the exclusive `flock`, and writes one readiness byte on `--ready-fd`.
4. Parent writes the pidfile **after** the readiness report, then returns.

Two authorities, only one of which is evidence:

- The **pidfile** says a recording ever started (written only after ready, removed on failure). It can go stale — a pid gets reused.
- The **flock** says a writer exists *now*. The kernel drops it on exit, so the lock is always authoritative.

On a clean stop (SIGTERM/SIGINT/`--until`) the child finishes the in-flight write, appends the `stopped` sentinel, fsyncs, and exits. A hard error writes no sentinel — so a recorder that died reads exactly like a coverage gap, because it is one.

`check` signals via the pidfile, waits for the **lock** to release (never the pid), and only then reads the log once. Reading while a writer can still append can only produce a shorter window than was recorded.

## The log is the source of truth

`watch` records raw evidence, so nothing trusts a state that could become unreachable. Two consequences a maintainer must preserve:

- The **header is authoritative for recording facts** (the cadence actually used, the URL, the alert set); the ruler API is authoritative for **rule facts** (`for`, `intervalSeconds`, kind). `check` always re-resolves definitions fresh and never reconstructs them from the header — the header duplicates `for`/`interval` only so the uploaded artifact is self-describing.
- The **cadence authority** is the header's `poll_every_seconds`, not the definitions. Re-deriving it would compare gaps recorded at an override cadence against default-cadence thresholds — fail-open in the faster-override direction.
