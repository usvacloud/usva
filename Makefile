SHELL 		= /bin/sh
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
CGO_ENABLED ?= 0
BUILDDIR ?= ./bin

DB_NAME ?= usva
DB_USERNAME ?= dev
DB_PASSWORD ?= dev
DB_OWNER ?= dev
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
DB_NAME_TESTS	?= usva_tests
DB_USERNAME_TESTS ?= usva_tests
DB_PASSWORD_TESTS ?= testrunner
START_TEST_DOCKER ?= 1
DB_CONNECTION_STRING = "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
DB_TESTS_CONNECTION_STRING = "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable"

.PHONY: all lint test

setup-and-build: setup build

build:
	@-mkdir $(BUILDDIR)
	@-CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BUILDDIR)/$(BINARY) .

setup:
	@-go get -u
	@-cp config-example.toml config.toml

migrateup:
	cat ./sqlc/schemas/* | psql -d $(DB_CONNECTION_STRING)

migratedown:
	psql -d $(DB_CONNECTION_STRING) -f ./sqlc/dbdown.sql

migratedown-tests:
	psql -d $(DB_TESTS_CONNECTION_STRING) -f ./sqlc/dbdown.sql

migrateup-tests:
	cat ./sqlc/schemas/* | psql -d $(DB_TESTS_CONNECTION_STRING)

db-create:
	createdb -U $(DB_USERNAME) --owner=$(DB_OWNER) $(DB_NAME)

db-shell:
	@psql -qd $(DB_CONNECTION_STRING)

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
		&& docker-compose -f docker-compose-tests.yml up -d --remove-orphans \
		&& sleep 3
	make migrateup-tests


tests-cleanup:
	- rm -r test-uploads postgres-tests .postgres-data
	- make migratedown-tests

	- docker-compose -f docker-compose-tests.yml down

lint:
	golangci-lint run ./...

format:
	go fmt ./...

act-verify-source:
	act -n -W ./.github/workflows/push-validation.yml

act-docker:
	act -n -W ./.github/workflows/gh-packages.yml

act: act-verify-source act-docker
