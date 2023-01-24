# Usva

[![push-validation](https://github.com/romeq/usva/actions/workflows/push-validation.yml/badge.svg)](https://github.com/romeq/usva/actions/workflows/push-validation.yml)

Usva is a feature-rich file cloud with a modern REST API.
Privacy, accessibility and efficiency will always come first, thus usva is built to be highly reliable and easy to set up for anybody.

If you have any questions or comments about usva's security practices, you can open an issue.

### Features

##### For service provider

- TLS support
- Ability to set certain limits to file upload (for example. maximum file size)
- Really easy configuration and setup with zero hassle
- Source code is always free to download.
- Option for disabling all request logging to enhance privacy on client
- Ratelimits

##### For client

- Files can be locked with a password (hashed with bcrypt)
- Endpoint for viewing file's metadata
- Endpoint for viewing file's content
- Downloading and deleting a file

## Installation and usage

Installation is done in 3 steps: downloading source, installing dependencies and compiling it.

### With Docker

```sh
% git clone https://github.com/romeq/usva && cd usva
% cp config-example.toml config.toml
% $EDITOR config.toml
% make run-docker
```

### Without Docker

```sh
% git clone https://github.com/romeq/usva && cd usva
% make setup build migratesetup
% # configure (see below)
% make run
```

#### Configuration for users without docker

```shell
% DB_HOST="127.0.0.1" # PostgreSQL server host
	DB_PORT=5432 # PostgreSQL server port
	DB_USERNAME="usva" # Username to log in with
	DB_PASSWORD="password" # Password to log in with
	DB_NAME="usva" # Database name for usva
	make migrateup # Run migrations
```

### Server configuration

#### Note for Docker users

Docker image uses config.toml for the server's configuration. By default this file is
shared between the host and the container for easy reconfiguration.

#### Configuration options

Full configuration can be found in `config-example.toml`, and we suggest for you to use
a modified copy of that file as server's configuration.

## API Specification

Most up-to-date API specification can be found in [APISPEC.md](../APISPEC.md)
