cd "$(dirname "$0")"
DAY="$(date +%F)"
BACKUPROOT="/data/backups/.rsync-history/proton"
BACKUPDEST="$BACKUPROOT/$DAY"

ssh admin@192.168.1.216 "mkdir -p '$BACKUPDEST'"

rsync -a --delete --backup --backup-dir="$BACKUPDEST" --exclude-from="./.exclusions" --itemize-changes \
  -e "ssh" \
  "~/Library/CloudStorage/ProtonDrive-nicholas@doorlay.com-folder/Offline/" \
  "admin@192.168.1.216:/data/backups/proton/"
