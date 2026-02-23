### Overview
This folder contains scripts used to sync files from my mac to my Raspberry Pi in my homelab. Every hour (and on boot), launchd will kick off scripts that will rsync over any updates.

### Setup
Grant Full Disk Access to rsync: 
- Navigate to System Settings → Privacy & Security → Full Disk Access.
- Click +, then press Cmd+Shift+G and paste the path outputted from `which rsync` and click on the application.

### Notes
- stdout is written to `/Users/doorlay/Library/Logs/run-backups.out.log`
- stderr is wrriten to `/Users/doorlay/Library/Logs/run-backups.err.log` 
