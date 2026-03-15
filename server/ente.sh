#!/usr/bin/env bash
set -Eeuo pipefail

command -v ente >/dev/null 2>&1 || {
  echo "ERROR: ente CLI not found in PATH. Install it first, then rerun." >&2
  exit 127
}

# Pull in env vars from .env
set -o allexport
source .env
set +o allexport

# Prevent concurrent runs (cron overlap, manual run, etc.)
LOCK="/run/lock/ente-cli-export.lock"
mkdir -p "$(dirname "$LOCK")"
exec 9>"$LOCK"
flock -n 9 || exit 0

# Headless-friendly secrets storage (no keyring)
install -d -m 700 "$(dirname "$SECRETS_PATH")"
touch "$SECRETS_PATH"
chmod 600 "$SECRETS_PATH"
export ENTE_CLI_SECRETS_PATH="$SECRETS_PATH"   # :contentReference[oaicite:4]{index=4}

# Self-hosted endpoint config (optional)
if [[ -n "$ENTE_API_ENDPOINT" ]]; then
  mkdir -p "$HOME/.ente"
  cat >"$HOME/.ente/config.yaml" <<EOF
endpoint:
  api: ${ENTE_API_ENDPOINT}
EOF
  # config.yaml can live in $HOME/.ente/config.yaml (among other locations). :contentReference[oaicite:5]{index=5}
fi

# Ensure export directory exists (ideally on a mounted disk/NAS)
mkdir -p "$EXPORT_DIR"

# (Optional but nice) reaffirm the export dir in case it changed
# This is the documented way to set where exports land. :contentReference[oaicite:6]{index=6}
ente account update --app photos --email "$ENTE_EMAIL" --dir "$EXPORT_DIR"

# Build export args
args=()
if [[ -n "$ALBUMS" ]]; then
  args+=(--albums "$ALBUMS")  # comma-separated album names :contentReference[oaicite:7]{index=7}
fi
if [[ "$INCLUDE_HIDDEN" == "true" ]]; then
  args+=(--hidden=false)
fi
if [[ "$INCLUDE_SHARED" == "true" ]]; then
  args+=(--shared=false)
fi

# Run the export (incremental + resumable) :contentReference[oaicite:8]{index=8}
{
  echo "[$(date -Is)] Starting Ente export to: $EXPORT_DIR"
  ente export "${args[@]}"
  echo "[$(date -Is)] Ente export finished OK"
} >>"$LOG_FILE" 2>&1
