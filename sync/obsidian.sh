DAY="$(date +%F)"
BACKUPROOT="/srv/backups/.rsync-history/obsidian"
BACKUPDEST="$BACKUPROOT/$DAY"

ssh admin@192.168.1.216 "mkdir -p '$BACKUPDEST'"

rsync -a --delete --backup --backup-dir="$BACKUPDEST" --exclude-from="./.exclusions" --itemize-changes \
  -e "ssh" \
  "/Users/doorlay/Documents/Notes/" \
  "admin@192.168.1.216:/srv/backups/obsidian/"
