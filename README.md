### Overview
Simple backup orchestrator for all my digital things. 
- Hourly sync from my mac to my server
- Hourly sync from my photo provider to my server
- Point-in-time encrypted backups in S3
- Push notifications for success/errors via ntfy

### Setup
#### Client
The client portion of the codebase runs on my mac, handling syncing from my mac to my server, kicked off hourly via launchd. Synced files end up in subdirectories inside `/data/backups/` on the server. To setup: 

1. Grant rsync Full Disk Access by navigating to System Settings → Privacy & Security → Full Disk Access. Click +, then press Cmd+Shift+G and paste the path outputted from `which rsync` and click on the application. Do the same for `/bin/bash` and `/bin/sh`.
2. Add source and destination paths to `client/backups.conf`.
3. Run `cd client && cp stub.env .env` and set `NTFY_TOPIC` to a random UUID. Download the ntfy app and subscribe to that UUID.
4. Run `make client`.

#### Server
The server portion of the codebase runs on my server, handling syncing from my photo provider (Ente) to my server, kicked off hourly via systemd. These files end up in `/data/backups/photos/` on the server. To setup:

1. Install the Ente CLI:
```
wget https://github.com/ente-io/ente/releases/latest/download/ente-linux-arm64 -O ente
chmod +x ente
sudo mv ente /usr/local/bin/
```
2. Install Golang if you haven't already. On a headless Pi, the easiest way is `wget https://dl.google.com/go/go{version}.linux-arm64.tar.gz`.
3. Perform the initial Ente login to generate the secrets file:
```
export ENTE_CLI_SECRETS_PATH="/srv/backups/ente-secrets/secrets.db"
ente account add
```
4. Run `cd server && cp stub.env .env` and fill out the environment variables.
5. Run `make server`.

To start the client backups, run `make client` on your mac to build and install the launchd agent.
To start the server backups, run `make server` on your server.

### Development
- `make build-client` — compile the client binary without deploying to launchd
- `make run-client` — build and run the client manually (useful for testing changes)
- `make client` — build and deploy to launchd (runs hourly via `RunAtLoad`)
- `make build-server` — compile the server binary without deploying to systemd
- `make server` — build and deploy the systemd service and timer

### Notes
- Client logs are written to `~/Library/Logs/backups.out.log` and `~/Library/Logs/backups.err.log`
- Server logs are in journald: `journalctl -u ente-sync.service`
- The client assumes the repo is cloned to `$HOME/Projects/` on your mac
- The server assumes the repo is cloned to `/srv/backups/` on your server
- There are probably other assumptions made in this code that are specific to my setup that will break if you try to run these, YMMV
