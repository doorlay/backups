.PHONY: clean build-client client run-client server

clean:
	rm backups

build-client:
	go build -o ~/bin/backups client/main.go

client: build-client
	./client/bootstrap.sh

run-client: build-client
	~/bin/backups

server:
	#
