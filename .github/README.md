# Usva 

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

##### For client

- Files can be locked with a password (hashed with bcrypt)
- Endpoint for viewing file's metadata
- Endpoint for viewing file's content
- Downloading and deleting a file



## Installation and usage

Installation is done in 3 steps: downloading source, installing dependencies and compiling it.

#### With Docker

```sh
% git clone https://github.com/romeq/usva && cd usva
# after configuration (see below)
% make run-docker
```

#### Without docker

```sh
% git clone https://github.com/romeq/usva && cd usva
% make setup build migratesetup 
# after setup (see below)
% make run 
```



## Configuration and setup

### Configuration and migrations using Docker

Migrations are automatically ran when server is started. Define following values in `.env`:

```shell
# Required
DB_USERNAME=my-awesome-username # postgres server username
DB_PASSWORD=my-super-secure-password # postgres server password

# Optional
DB_PORT=5434 # Alternative port for postgres server, default = 5432
DB_NAME=my-awesome-database # Alternative database name, default = usva
SV_PORT=8080 # Exposed port, default = 8080
```

### Migrations without Docker

```shell
% DB_HOST="127.0.0.1" \ # PostgreSQL server host
	DB_PORT=5432 \ # PostgreSQL server port
	DB_USERNAME="usva" \ # Username to log in with *
	DB_PASSWORD="password" \ # Password to log in with *
	DB_NAME="usva" # Database name for usva
	make migrateup # Run migrations
```

### Server configuration

Full configuration will look something near following: 

```toml
[Server]
Address = "127.0.0.1" # address to bind to
Port = 8080 # the port to bind to. don't use with docker.
TrustedProxies = ["127.0.0.1"]
DebugMode = false # use of gin's debug mode (includes logging)
HideRequests = false # requests should be hidden from logs
AllowedOrigins = [ "https://example.com" ] # cors allowed origins

[Server.TLS]
Enabled = true
CertFile = "/path/to/cert"
KeyFile = "/path/to/keyfile"

[Files]
MaxSize = 10 # maximum uploaded file size in megabytes
UploadsDir = "uploads" # directory for uploaded files

[Database]
User = "user"
Password = "password"
Host = "127.0.0.1"
Port = 5432
Database = "usva"
```

You should note, though, that only a few fields from above will be required. 
Smallest possible configuration looks something like below.

```toml
[Server]
Address = "127.0.0.1"
Port = 8080
AllowedOrigins = [ "http://example.com" ]

[Files]
UploadsDir = "uploads"

[Database]
User = "user"
```



## API Specification

Most up-to-date API specification can be found in [APISPEC.md](../APISPEC.md)
