.PHONY: clean client server

clean:
	rm backups

client: 
	go build -o ~/bin/backups client/main.go
	./client/bootstrap.sh

server:
	#
