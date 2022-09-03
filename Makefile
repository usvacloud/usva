SHELL 		= /bin/bash
GO			= go
BINARY		= tapsa
identifier	= tapsa

.PHONY: all lint test

install: build setup
beforecommit: test lint build clean
	go mod tidy
	go fmt ./...

build: 
	$(GO) build -o $(BINARY)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

