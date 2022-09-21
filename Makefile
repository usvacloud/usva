SHELL 		= /bin/bash
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
identifier	= usva
BUILDDIR	= ${BUILDDIR:?target}
CGO_ENABLED ?= 0


.PHONY: all lint test

setup: clean
	go get -u
	cp config-example.toml config.toml

migratesetup:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrateup: 
	$(GOPATH)/bin/migrate \
		-source file://migrations \
		-database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		up

migratedown: 
	$(GOPATH)/bin/migrate \
		-source file://migrations \
		-database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		down

run: migrateup
	./$(BINARY) -c config.toml

build: 
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BINARY)

test:
	go test ./...

lint:
	golangci-lint run ./...

format:
	go fmt ./...

clean:
	rm -f $(BINARY)

