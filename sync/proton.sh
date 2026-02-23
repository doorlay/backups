DAY="$(date +%F)"  # e.g. 2026-02-22
RETROOT="/srv/backups/proton/.rsync-history"
RETDEST="$RETROOT/$DAY"

#ssh admin@192.168.1.216 "mkdir -p '$RETDEST'"

rsync -a --delete --backup --backup-dir="$RETDEST" --exclude-from="./.exclusions" --itemize-changes \
  -e "ssh" \
  "/Users/doorlay/Library/CloudStorage/ProtonDrive-nicholas@doorlay.com-folder/Offline/" \
  "admin@192.168.1.216:/srv/backups/proton/"
