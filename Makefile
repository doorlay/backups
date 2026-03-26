.PHONY: clean build-client client run-client build-server server

clean:
	rm backups

build-client:
	go build -o ~/bin/backups client/main.go

client: build-client
	./client/bootstrap.sh

run-client: build-client
	~/bin/backups

build-server:
	go build -o /srv/backups/bin/ente-sync server/main.go

server: build-server
	sudo cp server/ente-sync.service server/ente-sync.timer server/restic-backup.service server/restic-backup.timer /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable --now ente-sync.timer
	sudo systemctl enable --now restic-backup.timer
