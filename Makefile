    # Borrowed from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
    .PHONY: help
    help: ## Display this help section
        @awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-38s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
    .DEFAULT_GOAL := help

test:
	go test -race -cover ./...

coverage:
	go test -race -covermode=count -coverprofile=coverage.out ./...

coverage-txt: coverage
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out

build:
	go build -race -o ./bin/chat ./chat

lint:
	golangci-lint run
	golint chat trace

start:
	go run ./chat


# TODO: go test all - before building a release, fetches and runs tests for all dependencies as well.
