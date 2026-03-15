#!/usr/bin/env bash
set -euo pipefail

LABEL="com.doorlay.backups"


# Must be run somewhere inside your git repo
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || true)"
if [[ -z "${REPO_ROOT}" ]]; then
  echo "Error: not inside a git repo (git rev-parse failed). Run this script from within your repo." >&2
  exit 1
fi

SRC="${REPO_ROOT}/client/launchd-config"
DST_DIR="${HOME}/Library/LaunchAgents"
DST="${DST_DIR}/${LABEL}.plist"

if [[ ! -f "${SRC}" ]]; then
  echo "Error: source file not found: ${SRC}" >&2
  exit 1
fi

mkdir -p "${DST_DIR}"

# Create/replace symlink
ln -sfn "${SRC}" "${DST}"
echo "Symlinked:"
echo "  ${DST} -> ${SRC}"

# Validate plist syntax if available
if command -v plutil >/dev/null 2>&1; then
  plutil -lint "${DST}" || true
fi

# Reload + start new symlinked launchd file
launchctl unload ~/Library/LaunchAgents/com.doorlay.backups.plist 2>/dev/null || true
launchctl load   ~/Library/LaunchAgents/com.doorlay.backups.plist
launchctl start  com.doorlay.backups
