SHELL 		= /bin/sh
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
CGO_ENABLED ?= 0
BUILDDIR ?= ./bin

DB_NAME ?= usva
DB_USERNAME ?= dev
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_NAME_TESTS	?= usva_tests
DB_USERNAME_TESTS ?= usva_tests
DB_PASSWORD_TESTS ?= testrunner
START_TEST_DOCKER ?= 1

.PHONY: all lint test

setup-and-build: setup build

build:
	@-mkdir $(BUILDDIR)
	@-CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BUILDDIR)/$(BINARY) cmd/server

setup: clean
	go get -u
	cp config-example.toml config.toml

migrateup:
	cat ./sqlc/schemas/* | psql \
		-d "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-f -

migratedown:
	psql \
		-d "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" \
		-f ./sqlc/dbdown.sql

migratedown-tests:
	psql \
		-d "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable" \
		-f ./sqlc/dbdown.sql

migrateup-tests:
	cat ./sqlc/schemas/* | psql \
		-d "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable" \
		-f -

db-create:
	 psql "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/" \
		-q -c "CREATE DATABASE $(DB_NAME) OWNER postgres ENCODING UTF8;"

run:
	$(BUILDDIR)/$(BINARY) -c ./config.toml

run-docker:
	docker-compose run --service-ports --rm -d server

run-docker-nodaemon:
	docker-compose run --service-ports --rm server

test:
	- echo "--------- GO TESTS -----------" 
	- go test ./...
	- echo "------------------------------" 

preparetests:
	- mkdir test-uploads postgres-tests
	- [ "${START_TEST_DOCKER}" = "1" ] \
		&& docker-compose -f docker-compose-dev.yml up -d --remove-orphans \
		&& sleep 3
	make migrateup-tests


tests-cleanup:
	- rm -r test-uploads postgres-tests
	- make migratedown-tests

	- docker-compose -f docker-compose-dev.yml down

lint:
	golangci-lint run ./...

format:
	go fmt ./...

clean:
	rm -rf $(BINARY) test-uploads postgres-tests
