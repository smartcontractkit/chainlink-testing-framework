---
id: grafana-alertcheck-log-format
title: Log format
sidebar_label: Log format
sidebar_position: 1
description: The JSONL log schema written by watch and read by check, for debugging the forensic artifact.
---

# Log format

`watch` records evidence to a JSONL log — one JSON object per line. A poll record *is* the heartbeat; there is no separate heartbeat type.

## Record types

Exactly three:

| `type` | Meaning |
| ------ | ------- |
| `header` | Line 1 — identity and the alert set |
| `poll` | One reduced observation of one rule |
| `stopped` | The sentinel, written on a clean stop only |

The header must be line 1, appear once, and carry `schema_version` `1` (any other value is a read error). Any unparseable line — including the last, or one after the sentinel — makes the log unreadable: a truncated log is evidence the recorder was killed, and must not pass.

## Header

```json
{
  "type": "header",
  "schema_version": 1,
  "url": "https://grafana.example.com",
  "grafana_version": "13.1.0",
  "started_at": "2026-09-07T10:00:00Z",
  "rules": [
    {
      "uid": "rule0000001",
      "title": "HighErrorRate",
      "folder": "Platform",
      "group": "api",
      "for_seconds": 300,
      "interval_seconds": 60,
      "is_paused": false,
      "no_data_state": "OK",
      "exec_err_state": "OK",
      "poll_every_seconds": 30
    }
  ]
}
```

- `url` and `rules` are the log's identity — `check` validates them against the current environment and a fresh ruler read.
- `is_paused` records the pause state at record start (the moment `skipped` means).
- `poll_every_seconds` is the cadence the recording **actually used** (after any `--poll-interval` override). `check` derives `maxGap` from it, never from `interval_seconds`.
- `for_seconds`, `interval_seconds`, `no_data_state`, `exec_err_state` are forensic only — `check` re-resolves definitions and never reads them back.

## Poll

```json
{
  "type": "poll",
  "rule_uid": "rule0000001",
  "grafana_now": "2026-09-07T10:00:30Z",
  "skew_ms": 20,
  "skew_bound_ms": 40,
  "latency_ms": 123,
  "found": true,
  "state": "inactive",
  "health": "ok",
  "last_evaluation": "2026-09-07T10:00:28Z",
  "is_paused": false,
  "histogram": { "alerting": 0, "normal": 2004 },
  "reasons": { "NoData": 1091 },
  "abnormal": [ { "labels": { "env": "prod" }, "state": "firing", "active_at": "2026-09-07T09:50:00Z", "value": "1.5" } ],
  "cleared": [ "env=prod\u0001..." ],
  "vanished": []
}
```

Field notes:

- `grafana_now` is the response's `Date` header — never the runner clock.
- `skew_ms`/`skew_bound_ms` are the per-poll clock-skew estimate and its uncertainty (RTT/2), in milliseconds for compactness only.
- `found: false` is an authoritative `2xx` in which this rule was absent — a transport failure is retried and never becomes a poll.
- `state`, `health`, `last_error` are raw rule-level strings, reporting-only.
- `histogram` is a verbatim copy of the response `totals`; written, never analysed.
- `reasons` counts non-empty instance reasons (`NoData`, `Error`, `KeepLast`, …); composite states stay visible only here.
- `abnormal` holds only instances whose **canonical** state is not `normal`.
- `cleared`/`vanished` are instance keys that left the bad set, resolved against the same response: `cleared` = a real recovery; `vanished` = a discontinuity, never a recovery.

Instance keys are a sorted `k=v\n` join of labels, so they correlate across polls without hashing.

## Stopped

```json
{ "type": "stopped", "at": "2026-09-07T10:10:30Z" }
```

`at` is the recorder's own stop time. `check` compares it against `to + transitionGrace`; absent or earlier is `unobservable` — never a pass.
