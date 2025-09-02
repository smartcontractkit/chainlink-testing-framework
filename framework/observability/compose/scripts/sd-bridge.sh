#!/usr/bin/env bash
# sd-bridge.sh (simplified)
# Discover Docker containers by label, pull each container's discovery JSON,
# add labels, merge + de-duplicate, and write a single file_sd JSON.

set -Eeuo pipefail

# --- Config (env) ---
LABEL_MATCH="${LABEL_MATCH:-framework=ctf}"
DEFAULT_PATH="${DISCOVERY_PATH:-/discovery}"
DEFAULT_PORT="${DISCOVERY_PORT:-6688}"
DEFAULT_SCHEME="${DISCOVERY_SCHEME:-http}"
PREFER_NETWORK="${NETWORK_NAME:-}"
OUT="${OUT:-/out/merged.json}"
SLEEP="${SLEEP:-15}"
REQUEST_TIMEOUT="${REQUEST_TIMEOUT:-5}"
REWRITE_TO_IP="${REWRITE_TO_IP:-0}"   # set to 1 to replace host with container IP

# --- Helpers ---
log(){ printf '[sd-bridge] %s\n' "$*" >&2; }
get_ip(){
  local cid="$1" net="$2"
  if [[ -n "$net" ]]; then
    docker inspect "$cid" | jq -r --arg n "$net" '.[0].NetworkSettings.Networks[$n].IPAddress // empty'
  else
    docker inspect "$cid" | jq -r '.[0].NetworkSettings.Networks | to_entries[0].value.IPAddress // empty'
  fi
}
get_label(){ docker inspect "$1" | jq -r --arg k "$2" '.[0].Config.Labels[$k] // empty'; }
get_name(){  docker inspect "$1" | jq -r '.[0].Name | ltrimstr("/")'; }
merge_and_dedupe(){
  jq -s '
    add // []
    | map({targets: (.targets // []), labels: (.labels // {})})
    | group_by(.labels)
    | map({labels: (.[0].labels), targets: ([.[].targets[]] | unique | sort)})
  '
}
atomic_write(){ local p="$1" t="$1.tmp"; cat >"$t" && mv "$t" "$p"; }

# --- Init ---
mkdir -p "$(dirname "$OUT")"
echo '[]' | atomic_write "$OUT"

# --- Main loop ---
while true; do
  mapfile -t cids < <(docker ps -q --filter "label=$LABEL_MATCH" || true)
  if (( ${#cids[@]} == 0 )); then
    echo '[]' | atomic_write "$OUT"
    log "no matching containers; wrote empty array"
    sleep "$SLEEP"; continue
  fi

  files=()
  for cid in "${cids[@]}"; do
    ip="$(get_ip "$cid" "$PREFER_NETWORK")"
    [[ -z "$ip" ]] && { log "skip ${cid:0:12}: no IP"; continue; }
    name="$(get_name "$cid")"

    # Per-container overrides (optional)
    path="$(get_label "$cid" prom_sd_path)";  path="${path:-$DEFAULT_PATH}"
    port="$(get_label "$cid" prom_sd_port)";  port="${port:-$DEFAULT_PORT}"
    scheme="$(get_label "$cid" prom_sd_scheme)"; scheme="${scheme:-$DEFAULT_SCHEME}"

    url="${scheme}://${ip}:${port}${path}"
    f="$(mktemp)"; files+=("$f")
    if curl -fsSL --max-time "$REQUEST_TIMEOUT" "$url" | jq '.' > "$f" 2>/dev/null; then
      # Add labels (and optionally rewrite host -> container IP)
      if [[ "$REWRITE_TO_IP" == "1" ]]; then
        jq --arg ip "$ip" --arg name "$name" '
          map(
            .targets |= map($ip + ":" + (split(":")[1])) |
            .labels = ((.labels // {}) + {
              container_name: $name,
              scrape_path: (.labels.__metrics_path__ // "")
            })
          )
        ' "$f" > "$f.tmp" && mv "$f.tmp" "$f"
      else
        jq --arg name "$name" '
          map(
            .labels = ((.labels // {}) + {
              container_name: $name,
              scrape_path: (.labels.__metrics_path__ // "")
            })
          )
        ' "$f" > "$f.tmp" && mv "$f.tmp" "$f"
      fi
      log "ok ${url}"
    else
      log "fail ${url}; using []"
      echo '[]' > "$f"
    fi
  done

  if (( ${#files[@]} > 0 )); then
    cat "${files[@]}" | merge_and_dedupe | atomic_write "$OUT"
    rm -f "${files[@]}"
    log "merged ${#files[@]} lists into $(wc -c < "$OUT") bytes"
  else
    echo '[]' | atomic_write "$OUT"
  fi

  sleep "$SLEEP"
done
