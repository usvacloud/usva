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
- Dockerized environment

##### For client

- Files can be locked with a password
- Files can be encrypted
- Endpoint for viewing file's information, such as
  - View count
  - Size
  - Encryption status
  - Date uploaded
- Downloading a file

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
git clone https://github.com/romeq/usva && cd usva

# setup database user
export DB_OWNER=usva
createuser $DB_OWNER -PU postgres

# create database
DB_USERNAME=postgres \
	DB_PASSWORD=postgres \
	DB_OWNER=$DB_OWNER \
	make db-create

# replace "dbownerpassword" with password you provided earlier on the createuser part
DB_USERNAME=$DB_OWNER \
	DB_PASSWORD=dbownerpassword \
	make migrateup setup build run
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
