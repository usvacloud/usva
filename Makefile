SHELL 		= /bin/sh
GO			= go
GOPATH 		= $(shell go env GOPATH)
BINARY		= usva
CGO_ENABLED ?= 0

DB_NAME_TESTS	?= usva_tests
DB_USERNAME_TESTS ?= usva_tests
DB_PASSWORD_TESTS ?= testrunner

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

migrateup-tests:
	$(GOPATH)/bin/migrate \
		-source file://migrations \
		-database "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable" \
		up

migratedown-tests:
	$(GOPATH)/bin/migrate \
		-source file://migrations \
		-database "postgres://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME_TESTS)?sslmode=disable" \
		down -all

db-create:
	@ psql "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/" \
		-q -c "CREATE DATABASE $(DB_NAME) OWNER postgres ENCODING UTF8;"

db-create-tests:
	@- psql "postgresql://$(DB_USERNAME_TESTS):$(DB_PASSWORD_TESTS)@$(DB_HOST):$(DB_PORT)/" \
		-q -c "CREATE DATABASE $(DB_NAME_TESTS) OWNER $(DB_USERNAME_TESTS) ENCODING UTF8;"

run:
	@./$(BINARY) -c ./config.toml

run-docker:
	@docker-compose run --service-ports --rm -d server

run-docker-nodaemon:
	@docker-compose run --service-ports --rm server

build:
	@CGO_ENABLED=$(CGO_ENABLED) $(GO) build -o $(BINARY)

test: preparetests
	@-echo "--------- GO TESTS -----------" 
	@ go test ./...
	@-echo "------------------------------" 

	@make tests-cleanup clean

preparetests:
	@- mkdir test-uploads postgres-tests 2>/dev/null
	@ docker-compose -f docker-compose-dev.yml up -d
	@ sleep 1
	@ make db-create-tests migrateup-tests


tests-cleanup:
	@- rm -r test-uploads
	@- make migratedown-tests

	@- docker-compose -f docker-compose-dev.yml down

lint:
	@golangci-lint run ./...

format:
	@go fmt ./...

clean:
	@rm -rf $(BINARY) test-uploads postgres-tests
