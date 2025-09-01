#!/usr/bin/env bash
# sd-bridge.sh
# Discover Docker containers by label, pull each container's discovery JSON, merge, de-duplicate,
# and write a single Prometheus file_sd JSON. Optionally serve the same JSON over HTTP.

set -Eeuo pipefail

# ---------- Configuration via env vars ----------
LABEL_MATCH="${LABEL_MATCH:-prom_sd=true}"      # filter for worker containers
DEFAULT_PATH="${DISCOVERY_PATH:-/discovery}"
DEFAULT_PORT="${DISCOVERY_PORT:-6688}"
DEFAULT_SCHEME="${DISCOVERY_SCHEME:-http}"
PREFER_NETWORK="${NETWORK_NAME:-}"              # optional Docker network to prefer for IP
OUT="${OUT:-/out/merged.json}"                  # file_sd output
SLEEP="${SLEEP:-15}"                            # seconds between scans
REQUEST_TIMEOUT="${REQUEST_TIMEOUT:-5}"         # curl timeout seconds

# Optional lightweight HTTP serving of OUT for http_sd
SERVE_ADDR="${SERVE_ADDR:-}"                    # example: ":8080" to serve /targets
SERVE_PATH="${SERVE_PATH:-/targets}"            # URL path

# ---------- Helpers ----------
log(){ printf '[sd-bridge] %s\n' "$*" >&2; }
get_ip(){
  local cid="$1"; local net="$2"
  if [[ -n "$net" ]]; then
    docker inspect "$cid" | jq -r --arg n "$net" '.[0].NetworkSettings.Networks[$n].IPAddress // empty'
  else
    docker inspect "$cid" | jq -r '.[0].NetworkSettings.Networks | to_entries[0].value.IPAddress // empty'
  fi
}
get_label(){
  local cid="$1" key="$2"
  docker inspect "$cid" | jq -r --arg k "$key" '.[0].Config.Labels[$k] // empty'
}
merge_and_dedupe(){
  # stdin: many JSON arrays of target groups
  # out: one array, grouped by identical labels with unique sorted targets
  jq -s '
    add // []
    | map({targets: (.targets // []), labels: (.labels // {})})
    | sort_by(.labels)
    | group_by(.labels)
    | map({labels: (.[0].labels), targets: ([.[].targets[]] | unique | sort)})
  '
}
atomic_write(){
  local path="$1"; local tmp="$1.tmp"
  cat > "$tmp" && mv "$tmp" "$path"
}

# ---------- Optional HTTP server ----------
serve_http(){
  # Serves OUT at SERVE_PATH. Requires busybox httpd or python3 in the container.
  if command -v busybox >/dev/null 2>&1; then
    log "serving http_sd at ${SERVE_ADDR}${SERVE_PATH} using busybox httpd"
    # Busybox serves a directory. Symlink requested path to OUT.
    local root; root="$(dirname "$OUT")"; mkdir -p "$root"
    # Keep a symlink named targets.json and rewrite on change
    ln -sf "$(basename "$OUT")" "$root/targets.json"
    exec busybox httpd -f -p "${SERVE_ADDR#:}" -h "$root"
  elif command -v python3 >/dev/null 2>&1; then
    log "serving http_sd at ${SERVE_ADDR}${SERVE_PATH} using python http.server"
    cd "$(dirname "$OUT")" && exec python3 -m http.server "${SERVE_ADDR#:}"
  else
    log "no http server available in image; install busybox or python3"
    sleep infinity
  fi
}

# Ensure output dir and initial file
mkdir -p "$(dirname "$OUT")"
echo '[]' | atomic_write "$OUT"

# If SERVE_ADDR is set, background a tiny polling loop that keeps a shadow file named targets.json
if [[ -n "$SERVE_ADDR" ]]; then
  # Run server in background subshell so main loop continues to update OUT
  serve_http &
fi

# ---------- Main loop ----------
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
    if [[ -z "$ip" ]]; then log "skip ${cid:0:12}: no IP"; continue; fi

    # Per container overrides via labels
    path="$(get_label "$cid" prom_sd_path)"; path="${path:-$DEFAULT_PATH}"
    port="$(get_label "$cid" prom_sd_port)"; port="${port:-$DEFAULT_PORT}"
    scheme="$(get_label "$cid" prom_sd_scheme)"; scheme="${scheme:-$DEFAULT_SCHEME}"

    url="${scheme}://${ip}:${port}${path}"
    f="$(mktemp)"; files+=("$f")
    if curl -fsSL --max-time "$REQUEST_TIMEOUT" "$url" | jq '.' > "$f" 2>/dev/null; then
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
