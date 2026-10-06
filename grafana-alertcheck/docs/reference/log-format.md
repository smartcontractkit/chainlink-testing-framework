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
  "ready_at": "2026-09-07T10:00:27Z",
  "rules": [
    {
      "key": "rule0000001",
      "uid": "rule0000001",
      "title": "HighErrorRate",
      "folder": "Platform",
      "group": "api",
      "source_kind": "grafana",
      "for_seconds": 300,
      "interval_seconds": 60,
      "is_paused": false,
      "no_data_state": "OK",
      "exec_err_state": "OK",
      "poll_every_seconds": 30
    },
    {
      "key": "ds:[\"vm\",\"ExampleMetrics\",\"ExampleTargetDown\",\"/etc/vm/rules/example.yml\"]",
      "uid": "",
      "title": "ExampleTargetDown",
      "group": "ExampleMetrics",
      "source_kind": "datasource",
      "datasource_uid": "vm",
      "datasource_name": "VM Prod",
      "file": "/etc/vm/rules/example.yml",
      "for_seconds": 300,
      "interval_seconds": 60,
      "is_paused": false,
      "poll_every_seconds": 30
    }
  ]
}
```

- `url` and `rules` are the log's identity — `check` validates them against the current environment and a fresh read.
- `key` is the rule's identity across both source kinds; `uid` is the API-given uid and is empty for a datasource-managed rule. `source_kind`, `datasource_uid`, `datasource_name` and `file` are additive (schema stays `1`) and let `check` re-resolve a datasource rule without discovery. A v1 log written before these fields existed still reads.
- `started_at` is when the recording opened; `ready_at` is when the first-observation pass completed and every watched, non-paused rule had been observed once. The pass is sequential, so `check` refuses a `from` before `ready_at` (a window opening inside the pass would rest on observations that do not exist). `ready_at` is absent on logs written before the field existed; `check` then falls back to `started_at`.
- `is_paused` records the pause state at record start (the moment `paused` means).
- `poll_every_seconds` is the cadence the recording used: always half the rule's `interval_seconds`. `check` derives `maxGap` from it, never by re-deriving from `interval_seconds`.
- `for_seconds`, `interval_seconds`, `no_data_state`, `exec_err_state` are forensic only — `check` re-resolves definitions and never reads them back.

## Poll

```json
{
  "type": "poll",
  "rule_key": "rule0000001",
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

- `rule_key` is the identity across both kinds; `rule_uid` is kept for compatibility and is empty for a datasource-managed rule. A reader uses `rule_key` when present, else `rule_uid`, so an old v1 log stays readable.
- `grafana_now` is the response's `Date` header — never the runner clock.
- `skew_ms`/`skew_bound_ms` are the per-poll clock-skew estimate and its uncertainty (RTT/2), in milliseconds for compactness only.
- `found: false` is an authoritative `2xx` in which this rule was absent — a transport failure is retried and never becomes a poll.
- `keep_firing_for_ms` is the rule's recovery period as reported by this response; `0`/absent means no instance can be `recovering`. A datasource rule reports it from the backend's keep-firing-for (`keep_firing_for` on vmalert, `keepFiringFor` on Prometheus/Mimir), but such a rule never reaches `recovering` (the backend keeps it firing, then drops it).
- `state`, `health`, `last_error` are raw rule-level strings, reporting-only.
- `histogram` is a verbatim copy of the response `totals`; written, never analysed.
- `reasons` counts non-empty instance reasons (`NoData`, `Error`, `KeepLast`, …); composite states stay visible only here.
- `abnormal` holds only instances whose **canonical** state is not `normal`.
- `cleared`/`vanished` are instance keys that left the bad set, resolved against the same response: `cleared` = a real recovery; `vanished` = a discontinuity, never a recovery.

Instance keys are the JSON encoding of the labels map (with stable key order), so they correlate across polls without hashing.

## Stopped

```json
{ "type": "stopped", "at": "2026-09-07T10:10:30Z" }
```

`at` is the recorder's own stop time. `check` compares it against `to + transitionGrace`; absent or earlier is `not_verified` — never a pass.
