### Overview
Simple backup orchestrator from my Mac to my homelab, using rsync and macOS's launchd to sync files every hour (and on boot).

### Setup
Run `./bootstrap.sh` to create and load a new launchd configuration. Run this script any time updates are made to the launchd-config file.
 
Grant Full Disk Access: 
- Navigate to System Settings → Privacy & Security → Full Disk Access.
- Click +, then press Cmd+Shift+G and paste the path outputted from `which rsync` and click on the application. Do the same for `/bin/bash` and `/bin/sh`.

### Notes
- stdout is written to `/Users/doorlay/Library/Logs/backups.out.log`
- stderr is wrriten to `/Users/doorlay/Library/Logs/backups.err.log`
- These scripts assume the repo is cloned to `$HOME/Projects/`
- There are probably other assumptions made in this code that are specific to my setup that will break if you try to run these, YMMV 
