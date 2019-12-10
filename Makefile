test:
	go test -cover ./...

coverage:
	go test -covermode=count -coverprofile=coverage.out ./...

coverage-txt: coverage
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out

build:
	go build -o ./bin/chat ./chat

lint:
	golangci-lint run
	golint chat trace

start-server:
	./bin/chat

start:
	make build
	make start-server
