#!/usr/bin/env bash
set -euo pipefail

BACKUP_PATH="/data/backups"
RESULTS_FILE="/srv/backups/ente-sync-results.log"

notify() {
  if [[ -n "${NTFY_TOPIC:-}" ]]; then
    curl -s -d "$1" "https://ntfy.sh/${NTFY_TOPIC}" >/dev/null
  fi
}

record_result() {
  echo "$1" >> "${RESULTS_FILE}"
}

echo "Starting restic backup of ${BACKUP_PATH}"

if ! restic backup "${BACKUP_PATH}"; then
  record_result "FAIL"
  notify "Restic backup failed"
  exit 1
fi

echo "Running restic forget/prune"

if ! restic forget --prune --keep-daily 30 --keep-weekly 26 --keep-monthly 12; then
  record_result "FAIL"
  notify "Restic prune failed"
  exit 1
fi

record_result "OK"
echo "Restic backup completed successfully"
