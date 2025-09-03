#!/usr/bin/env bash
# sd-merge.sh
# Discover Docker containers by label, fetch each container's /discovery JSON,
# add labels (container_name, scrape_path), merge + dedupe, and write a single file_sd JSON.
# TODO: This script should be removed once we convert the prom metrics to OTEL

set -Eeuo pipefail

# -------------------- Configuration (via env) --------------------
LABEL_MATCH="${LABEL_MATCH:-framework=ctf}"   # docker ps --filter "label=$LABEL_MATCH"
DEFAULT_PATH="${DISCOVERY_PATH:-/discovery}"  # default discovery path inside each container
DEFAULT_PORT="${DISCOVERY_PORT:-6688}"        # default discovery port
DEFAULT_SCHEME="${DISCOVERY_SCHEME:-http}"    # http or https
PREFER_NETWORK="${NETWORK_NAME:-}"            # prefer IP from this Docker network (optional)
OUT="${OUT:-/out/merged.json}"                # file_sd output path
SLEEP="${SLEEP:-15}"                          # seconds between scans
REQUEST_TIMEOUT="${REQUEST_TIMEOUT:-5}"       # curl timeout (s)
REWRITE_TO_IP="${REWRITE_TO_IP:-0}"           # 1 = replace host with container IP in targets

# -------------------- Helpers --------------------
log(){ printf '[sd-merge] %s\n' "$*" >&2; }

# Atomic writer: reads stdin, writes to $1.tmp, then mv -> $1
atomic_write(){
  local path="$1" tmp="$1.tmp"
  cat > "$tmp" && mv "$tmp" "$path"
}

# -------------------- Init --------------------
mkdir -p "$(dirname "$OUT")"
echo '[]' | atomic_write "$OUT"

# -------------------- Main loop --------------------
while true; do
  # List container IDs matching the label
  mapfile -t cids < <(docker ps -q --filter "label=$LABEL_MATCH" || true)

  if ((${#cids[@]} == 0)); then
    echo '[]' | atomic_write "$OUT"
    log "no matching containers; wrote empty array"
    sleep "$SLEEP"
    continue
  fi

  # Emit each container's (possibly empty) discovery array, then merge once with jq -s
  {
    for cid in "${cids[@]}"; do
      # Inspect once, reuse for IP, name, and labels
      inspect="$(docker inspect "$cid" 2>/dev/null || true)"
      [[ -z "$inspect" ]] && { log "skip ${cid:0:12}: inspect failed"; echo '[]'; continue; }

      # Resolve container IP (optionally prefer a specific network)
      if [[ -n "$PREFER_NETWORK" ]]; then
        ip="$(jq -r --arg n "$PREFER_NETWORK" '.[0].NetworkSettings.Networks[$n].IPAddress // ""' <<<"$inspect")"
      else
        ip="$(jq -r '.[0].NetworkSettings.Networks | to_entries[0].value.IPAddress // ""' <<<"$inspect")"
      fi
      [[ -z "$ip" ]] && { log "skip ${cid:0:12}: no IP"; echo '[]'; continue; }

      # Container name and optional per-container overrides
      name="$(jq -r '.[0].Name | ltrimstr("/")' <<<"$inspect")"
      path="$(jq -r '.[0].Config.Labels.prom_sd_path // empty' <<<"$inspect")";   path="${path:-$DEFAULT_PATH}"
      port="$(jq -r '.[0].Config.Labels.prom_sd_port // empty' <<<"$inspect")";   port="${port:-$DEFAULT_PORT}"
      scheme="$(jq -r '.[0].Config.Labels.prom_sd_scheme // empty' <<<"$inspect")"; scheme="${scheme:-$DEFAULT_SCHEME}"

      url="${scheme}://${ip}:${port}${path}"

      # Fetch discovery JSON; treat errors as empty array
      payload="$(curl -fsSL --max-time "$REQUEST_TIMEOUT" "$url" 2>/dev/null || echo '[]')"

      # Normalize to array, add labels, optional host->IP rewrite while keeping port from targets
      if [[ "$REWRITE_TO_IP" == "1" ]]; then
        jq --arg ip "$ip" --arg name "$name" '
          (if type=="array" then . else [] end)
          | map(
              .targets |= map( $ip + ":" + (split(":")[1]) ) |
              .labels = ((.labels // {}) + {
                container_name: $name,
                scrape_path: (.labels.__metrics_path__ // "")
              })
            )
        ' <<<"$payload"
      else
        jq --arg name "$name" '
          (if type=="array" then . else [] end)
          | map(
              .labels = ((.labels // {}) + {
                container_name: $name,
                scrape_path: (.labels.__metrics_path__ // "")
              })
            )
        ' <<<"$payload"
      fi

      log "ok $url"
    done
  } \
  | jq -s '
      # Merge all arrays, coerce to {targets,labels}, then group by labels and dedupe targets
      add // []
      | map({targets: (.targets // []), labels: (.labels // {})})
      | group_by(.labels)
      | map({ labels: (.[0].labels)
            , targets: ([.[].targets[]] | unique | sort)
            })
    ' \
  | atomic_write "$OUT"

  log "wrote $(wc -c < "$OUT") bytes to $OUT"
  sleep "$SLEEP"
done
