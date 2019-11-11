test:
	go test -cover ./...

build:
	go build -o ./bin/chat ./chat

start-server:
	./bin/chat

start:
	make build
	make start-server
