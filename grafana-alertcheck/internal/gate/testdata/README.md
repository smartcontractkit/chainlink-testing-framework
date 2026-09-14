# Fixture provenance

All fixtures are sanitized slices of real Grafana 13.1.0 payloads captured into
`tmp/` (`tmp/state_all.json`, `tmp/ruler_all.json`, `tmp/health.json` — gitignored, never committed).
Renames are consistent across files: the same real folder/rule keeps the same fake identity everywhere
it appears (e.g. `folder0000002`/`rule0000002` is the same real paused rule in both
`state_paused.json` and `ruler_rules.json`).

An earlier revision embedded these notes as a top-level `_fixture_note` JSON key. That works for the
state-endpoint fixtures (the parser decodes the top level generically, so a stray key is just ignored),
but breaks the ruler-endpoint fixtures: `ParseDefinitions` decodes the whole body directly as
`map[string][]group`, and a `_fixture_note` string value can't unmarshal as `[]group`. Notes now live
here instead, for every fixture, for consistency.

## State endpoint (`/api/prometheus/grafana/api/v1/rules`)

- **state_one_instance.json** — real rule (`bfhp23rgt18u8f`, folder `BCM`, "[PROD][ACE] High CPU
  Utilization") with exactly one live instance. Renamed to folder `ExampleTeam`/`folder0000001`, rule
  `rule0000001`/"Example High CPU Usage", datasource `datasource000001`. Also the template the
  high-cardinality test helper (`synthesizeHighCardinalityState`) clones from.
- **state_paused.json** — one of the 3 real `isPaused:true` rules (`ac153fee-...`, folder `Mercury`,
  "Rubberbanding"). Renamed to folder `ExampleMetrics`/`folder0000002`, rule
  `rule0000002`/"Example Paused Rule". Unmodified: `isPaused:true`, zero `lastEvaluation`, `health:ok`,
  `state:inactive`, absent `alerts`/`labels`.
- **state_health_error.json** — real `health:error` rule ("[JD] No Job Proposals", folder
  `job-distributor`). Renamed to folder `ExampleService`/`folder0000003`,
  rule `rule0000003`/"Example No Data Source". Unmodified: `health:error`, `lastError` text, the single
  `Error` instance.
- **state_health_nodata.json** — real `health:nodata` rule ("ARE test", folder `diegos_playground`).
  Renamed to folder `ExamplePlayground`/`folder0000004`, rule `rule0000004`/"Example NoData Rule".
  Unmodified: `health:nodata`, the single `NoData` instance.
- **state_reason_composite.json** — two real instances combined under one rule to cover composite
  state parsing: a real `"Normal (Error)"` instance (from a Flux-reconciliation rule; 14 of that state exist
  in the capture) and a real `"Normal (NoData)"` instance (from a pod-liveness rule; 1091 of that state
  exist), plus one plain `"Normal"` instance for contrast. Renamed to folder
  `ExampleInfra`/`folder0000005`, rule `rule0000005`/"Example Composite Reasons".
- **state_missing_optional.json** — derived from `state_one_instance.json`: `alerts`, `totals`,
  `totalsFiltered` and `labels` all removed. Must parse with `Instances=nil`, `Totals=nil`.
- **state_missing_health.json**, **state_missing_lasteval.json**, **state_missing_state.json**,
  **state_missing_interval.json** — derived from `state_one_instance.json`, each with one of the four
  required keys removed (`health`, `lastEvaluation`, rule-level `state`, group-level `interval`). Each
  must be a parse error.
- **state_missing_file.json** / **state_missing_name.json** — derived from `state_one_instance.json`:
  the group-level `file`/`name` keys removed respectively. Not among the four required fields above,
  but the parser treats group identity as strict too.
- **state_zerotime_unpaused.json** — derived from `state_one_instance.json`: `lastEvaluation` set to
  the zero time while `isPaused` stays `false`. Must be a parse error — only a paused rule may report
  the zero time.
- **state_unknown_state.json** — derived from `state_one_instance.json`: the instance state hand-edited
  to `"Weird (NoData)"`, a syntactically valid composite whose base isn't in the 5-value allowlist. Must
  be a parse error.
- **state_only_active_instances.json** — derived from a real rule that genuinely had 1 `Alerting` + 22
  `Normal` instances (`totals: {alerting:1, normal:22}`, rule `dfhp1t5pkosu8f`, folder `BCM`). `alerts[]`
  trimmed to the single `Alerting` instance only, while `totals` is left **unchanged** — the shape a
  state endpoint that stopped returning normal instances would produce (the instance list says "only
  active" while totals disagrees). Renamed to folder `ExampleTeam`/`folder0000001`, rule `rule0000006`.
  `ParseState` itself parses this fine; `VerifyNormalInstancesVisible` is what rejects it.

## Ruler endpoint (`/api/ruler/grafana/api/v1/rules`)

- **ruler_rules.json** — contains:
  - The real true 2-way title collision: namespace `CRE-BCM-Prod-Zone-A`, group `Gateway`, identical
    folder+group+title, distinct UIDs (`ffvabtvvbozcwf`/`efvabtwbxlvk0b`) — renamed to namespace
    `Example-Zone-A`, rules `rule0000006a`/`rule0000006b`, both titled "Example No Gateways Available".
    Folder/Group/Title alone does **not** disambiguate this pair; only `uid:` does.
  - The 3 real `is_paused:true` rules, renamed to `rule0000002`/`rule0000007`/`rule0000008`.
    `rule0000002` intentionally shares its identity (`folder0000002`) with `state_paused.json`.
  - A real `for:1d` rule (`afs438kjd4v7kd` → `rule0000009`).
  - **`rule0000010` is DERIVED**: no `for:1w` rule exists anywhere in the capture (verified). Built by
    copying the `for:1d` rule and changing `for` to `1w` and its identity, to exercise the `w` unit.
- **ruler_datasource_managed.json** — **DERIVED**, no datasource-managed rule exists in the capture
  (verified: 0 rules lack a `grafana_alert` block). Hand-built minimal shape: a rule object with no
  `grafana_alert` key at all and an `alert` field carrying its Prometheus-format name — exactly how
  Grafana represents a datasource-managed (native Prometheus-format) alerting rule. `ParseDefinitions`
  must classify it as `KindDatasourceManaged`, parse `Title` from `alert`, and leave `UID` empty (this
  shape has no uid at all — inventing one would be inventing shape) without rejecting the rule
  (rejection is `Resolve`'s job, and only for rules a user actually named).
- **ruler_recording.json** — **DERIVED**, no recording rule exists in the capture (verified: 0 rules
  carry `grafana_alert.record`). Hand-built: a `grafana_alert` block with a `record` sub-object but
  deliberately *without* `no_data_state`/`exec_err_state`/`is_paused`/`intervalSeconds`/`namespace_uid`
  — those are alerting-only concepts a recording rule may not carry, and since none exist in the
  capture that shape is unverified either way. `ParseDefinitions` must classify it as `KindRecording`
  and must not require those fields for this Kind (requiring them bricks `ParseDefinitions` for every
  named rule in the same response over one recording rule elsewhere in the fleet).
