test:
	go test -cover ./...

build:
	go build -o ./bin/chat ./chat

start-server:
	./bin/chat

start: build start-server
