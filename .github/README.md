# Usva 

Usva is a feature-rich file cloud with a modern REST API. 
Privacy, accessibility and efficiency will always come first, thus usva is built to be highly reliable and easy to set up for anybody.

If you have any questions or comments about usva's security practices, you can open an issue.

### Features

##### For service provider

- TLS support
- Ability to set certain limits to file upload (for example. maximum file size)
- Really easy configuration and hosting with **zero bs** - source code is always free to download.
- Option for disabling all request logging to enhance privacy on client

##### For client

- Files can be locked with a password (hashed with bcrypt)
- Endpoint for viewing file's data
- Downloading and deleting a file

## Installation

Installation is done in 3 steps: downloading source, installing dependencies and compiling it.

```sh
% git clone https://github.com/romeq/usva && cd usva
% go get -u
% make build
# If everything went well,
# binary is now compiled in current working directory.
# You can launch it now with ./usva -c config_example.toml  
```

### Configuration

Full configuration will look something near following: 

```toml
[Server]
Address = "127.0.0.1" # address to bind to
Port = 8080 # the port to bind to
TrustedProxies = ["127.0.0.1"]
DebugMode = false # use of gin's debug mode (includes logging)
HideRequests = false # requests should be hidden from logs
# cors allowed origins
AllowedOrigins = [ "example.com" ]

[Server.TLS]
Enabled = true
CertFile = "/path/to/cert"
KeyFile = "/path/to/keyfile"

[Files]
MaxSize = 10 # maximum uploaded file size in megabytes
UploadsDir = "uploads" # directory for uploaded files
```

You should note, though, that only a few fields from above will be required.
Smallest possible configuration looks something like below.

```toml
[Server]
Address = "127.0.0.1"
Port = 8080
# cors allowed origins
AllowedOrigins = [ "example.com" ]

[Files]
UploadsDir = "uploads"
```

## API Specification

Most up-to-date API specification can be found in `APISPEC.md`
