cd "$(dirname "$0")"
DAY="$(date +%F)"
DEST="/data/backups/obsidian"
HIST=".rsync-history/$DAY"

ssh admin@192.168.1.216 "mkdir -p '$DEST/$HISTREL'"

rsync -a --delete \
  --backup \
  --backup-dir="$HISTREL" \
  --exclude="/.rsync-history/" \
  -e "ssh" \
  "$HOME/Notes/" \
  "admin@192.168.1.216:$DEST/"
