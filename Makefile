SHELL 		= /bin/sh
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
CGO_ENABLED ?= 1
BUILDDIR ?= ./bin

GOPKG=./cmd/webserver

DB_NAME ?= usva
DB_USERNAME ?= usva
DB_PASSWORD ?= usva
DB_OWNER ?= usva
DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432

DB_NAME_TESTS	?= usva_tests
DB_USERNAME_TESTS ?= usva_tests
DB_PASSWORD_TESTS ?= testrunner

START_TEST_DOCKER ?= 1
DB_CONNECTION_STRING = "postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
DB_TESTS_CONNECTION_STRING = "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable"

.PHONY: all lint test

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BUILDDIR)/$(BINARY) $(GOPKG)
setup:
	go get $(GOPKG)
	cp config-example.toml config.toml
deploy:
	docker compose build
	docker compose restart

migrateup:
	cat $$(ls ./sqlc/schemas/* | sort) | psql -d $(DB_CONNECTION_STRING) >/dev/null
migratedown:
	psql -d $(DB_CONNECTION_STRING) -f ./sqlc/dbdown.sql >/dev/null
migrateup-tests:
	cat $$(ls ./sqlc/schemas/* | sort) | psql -d $(DB_TESTS_CONNECTION_STRING) >/dev/null
migratedown-tests:
	psql -d $(DB_TESTS_CONNECTION_STRING) -f ./sqlc/dbdown.sql >/dev/null

db-create:
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USERNAME) --owner=$(DB_OWNER) $(DB_NAME)
db-shell:
	psql -qd $(DB_CONNECTION_STRING)

# running
run: build
	$(BUILDDIR)/$(BINARY) -c ./config.toml
run-docker:
	docker-compose run --service-ports --rm -d server
run-docker-nodaemon:
	docker-compose run --service-ports --rm server

test:
	echo "--------- GO TESTS -----------"
	- go test ./...
	echo "------------------------------"
preparetests:
	[ "${START_TEST_DOCKER}" = "1" ] \
		&& docker-compose -f docker-compose-tests.yml up -d --remove-orphans \
		&& sleep 3
	make migrateup-tests
tests-cleanup:
	rm -rf postgres-tests .postgres-data
	make migratedown-tests

	docker-compose -f docker-compose-tests.yml down

lint:
	golangci-lint run ./...
format:
	go fmt ./...

act-verify-source:
	act -W ./.github/workflows/push-validation.yml
act-docker:
	act -n -W ./.github/workflows/gh-packages.yml
act: act-verify-source act-docker
