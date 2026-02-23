#!/usr/bin/env bash
set -euo pipefail

# Prevent overlapping runs (portable lock)
LOCKDIR="/tmp/run-backups.lock"
if ! mkdir "$LOCKDIR" 2>/dev/null; then
  exit 0
fi
trap 'rmdir "$LOCKDIR"' EXIT

# Optional: skip if Pi isn't reachable (avoids noisy errors off-network)
PI_HOST="admin@192.168.1.216"
if ! /usr/bin/ssh -o BatchMode=yes -o ConnectTimeout=5 "$PI_HOST" "true" >/dev/null 2>&1; then
  exit 0
fi

"$HOME/Projects/scripts/backup/obsidian.sh"
"$HOME/Projects/scripts/backup/proton.sh"
