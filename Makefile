clean:
	rm backup

client: build-client run-client

build-client:
	go build -o backup client/main.go 

run-client:
	./client/bootstrap.sh

build-server:
	#
	
run-server:
	#
