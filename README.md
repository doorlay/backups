### Overview
Simple, bespoke backup orchestrator from all my digital things. 
<!--- Follows 3-2-1 backup rule (three copies, two media types, one off-site)-->
- Syncs files hourly from my Mac to my homelab
- Syncs photos hourly from my photo provider's cloud to my homelab
<!--- Uploads point-in-time encrypted backups to S3 using restic  -->

### Setup
Run `cp stub.env .env` and fill out all of the environment variables within `.env`.

Run `./bootstrap.sh` to create and load a new launchd configuration. Run this script any time updates are made to the launchd-config file.
 
Grant Full Disk Access: 
- Navigate to System Settings → Privacy & Security → Full Disk Access.
- Click +, then press Cmd+Shift+G and paste the path outputted from `which rsync` and click on the application. Do the same for `/bin/bash` and `/bin/sh`.

Run `ente account add` on your pi.

### Notes
- stdout is written to `~/Library/Logs/backups.out.log`
- stderr is wrriten to `~/Library/Logs/backups.err.log`
- These scripts assume the repo is cloned to `$HOME/Projects/`
- There are probably other assumptions made in this code that are specific to my setup that will break if you try to run these, YMMV 
