# Usva

[![push-validation](https://github.com/romeq/usva/actions/workflows/push-validation.yml/badge.svg)](https://github.com/romeq/usva/actions/workflows/push-validation.yml)

Usva is a feature-rich file cloud with a modern REST API.
Privacy, accessibility and efficiency will always come first, thus usva is built to be highly reliable and easy to set up for anybody.

If you have any questions or comments about usva's security practices, you can open an issue.

### Features

##### For service provider

- TLS support
- Ability to set certain limits to file upload (for example: maximum file size, maximum file size for encrypted files)
- Easy, straightforward and documented setup
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



<span style="font-size: 15px;color:green;background:lightgreen;border-radius:3px;">Recommended!</span>

### With Docker

```sh
% git clone https://github.com/romeq/usva && cd usva
% cp config-example.toml config.toml
% $EDITOR config.toml
% make run-docker
```

##### In configuration

- `server.address` has to be `0.0.0.0` so that the server can be accessed from outside
- if `SVPORT` environment variable is set, `server.port` has to be equal to it
- `database.host` must be equal to`"db"`, which is the database's container name in `docker-compose.yml`
- `database.port` has to be 5432 or same as `DB_PORT`
- `database.user` has to be same as `DB_USERNAME`
- `database.password` has to be same as `DB_PASSWORD`  
- `database.database` has to be same as `DB_NAME`



### Without Docker

```sh
git clone https://github.com/romeq/usva && cd usva

# setup database user
# this step expects that you use "postgres" as the administrator user, if your system diverges from that just specify your system's one
export DB_HOST="127.0.0.1"
export DB_PORT=5432
export DB_OWNER=usva
createuser $DB_OWNER -PU postgres

# create database
DB_USERNAME=postgres \
	DB_PASSWORD=postgres \
	DB_OWNER=$DB_OWNER \
	make db-create

# configure server
cp config-example.toml config.toml
$EDITOR config.toml

# replace "dbownerpassword" with password you provided earlier on the createuser part
DB_USERNAME=$DB_OWNER \
	DB_PASSWORD=dbownerpassword \
	make migrateup setup build run
```

##### Following these steps, if database is located in the same machine, configure Usva as follows:

- `database.user` has to be same as `DB_OWNER` 
- `database.password` has to be the same password you provided when you created the user
- `database.host` has to be `"127.0.0.1"`
- `database.port` has to be the same port where database is listening, default for PostgreSQL is  `5432`
- `database.database` has to be same as `DB_OWNER`

### Server configuration

#### Note for Docker users

Docker image uses config.toml for the server's configuration. By default this file is
shared between the host and the container for easy reconfiguration.

#### Configuration options

Full configuration can be found in `config-example.toml`, and we suggest for you to use
a modified copy of that file as server's configuration.

## API Specification

Most up-to-date API specification can be found in [APISPEC.md](../APISPEC.md)
