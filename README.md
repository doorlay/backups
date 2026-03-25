### Overview
Simple backup orchestrator for all my digital things. 
- Hourly sync from my mac to my server
- Hourly sync from my photo provider to my server
- Point-in-time encrypted backups in S3

### Setup
#### Client
The client portion of the codebase runs on my mac, handling syncing from my mac to my server, kicked off hourly via launchd. Synced files end up in subdirectories inside `/data/backups/` on the server. To setup: 

1. Grant rsync Full Disk Access by navigating to System Settings → Privacy & Security → Full Disk Access. Click +, then press Cmd+Shift+G and paste the path outputted from `which rsync` and click on the application. Do the same for `/bin/bash` and `/bin/sh`.
2. Add source and destination paths to `client/backups.conf`.
3. Run `cd client && cp stub.env .env` and set `NTFY_TOPIC` to a random UUID. Download the ntfy app and subscribe to that UUID.
4. Run `make client`.

#### Server 
The server portion of the codebase runs on my server, handling syncing from my photo provider to my server, kicked off hourly via systemd. These files end up in `/data/backups/photos/` on the server. To setup:

1. Download the latest Linux ARM64 release from Ente's GitHub:
```
wget https://github.com/ente-io/ente/releases/latest/download/ente-linux-arm64 -O ente
chmod +x ente
sudo mv ente /usr/local/bin/
```

2. Create a directory for your backup tool and a .env file. We will use a systemd environment file for the most "Linux-native" approach:
```
ENTE_EMAIL=your-email@example.com
EXPORT_DIR=/mnt/external_drive/ente_photos
SECRETS_PATH=/home/pi/.ente/secrets.db
# Optional: comma separated list of album names
ALBUMS=Family,Vacation
INCLUDE_HIDDEN=true
```

To setup the `server` portion:
3. Since you are on a headless Pi, you must perform the initial login once to generate the secrets file:
```
export ENTE_CLI_SECRETS_PATH="/home/pi/.ente/secrets.db"
ente account login --email your-email@example.com
```

4. Automate with Systemd:
```
```
- If you haven't already downloaded Golang on your server, do so. If you're on a headless Pi (like myself), the easiest way is to run `wget https://dl.google.com/go/go{}.linux-arm64.tar.gz`, adding the semantic version into the above URL. 
- Run `cp stub.env .env` and fill out all of the environment variables within `.env`.
- Run `ente account add` on your pi.

To start the client backups, run `make client` on your mac to build and install the launchd agent.
To start the server backups, run `make run-server` on your server.

### Development
- `make build-client` — compile the client binary without deploying to launchd
- `make run-client` — build and run the client manually (useful for testing changes)
- `make client` — build and deploy to launchd (runs hourly via `RunAtLoad`)

### Notifications
The client supports push notifications via [ntfy.sh](https://ntfy.sh). Set `NTFY_TOPIC` in `client/.env` to enable (see setup step 3).

### Notes
- stdout is written to `~/Library/Logs/backups.out.log`
- stderr is written to `~/Library/Logs/backups.err.log`
- These scripts assume the repo is cloned to `$HOME/Projects/`
- There are probably other assumptions made in this code that are specific to my setup that will break if you try to run these, YMMV
