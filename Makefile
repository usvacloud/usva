SHELL 		= /bin/bash
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
identifier	= usva
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
		-database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST:-127.0.0.1):$(DB_PORT:-5432)/$(DB_NAME:-usva)?sslmode=disable" \
		up

migratedown:
	$(GOPATH)/bin/migrate \
		-source file://migrations \
		-database "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST:-127.0.0.1):$(DB_PORT:-5432)/$(DB_NAME:-usva)?sslmode=disable" \
		down

run:
	./$(BINARY) -c ./config.toml

run-docker:
	sudo docker-compose run --service-ports --rm -d server

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

