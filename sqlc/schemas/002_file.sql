CREATE TABLE IF NOT EXISTS file(
    file_uuid TEXT PRIMARY KEY, -- TEXT because file_uuid includes the extension for the file
    title VARCHAR(256),
    passwdhash VARCHAR(512),
    access_token VARCHAR(40) NOT NULL UNIQUE,
    encrypted BOOLEAN NOT NULL DEFAULT FALSE,
    file_size INTEGER NOT NULL DEFAULT 0,
    encryption_iv BYTEA DEFAULT NULL,
    upload_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    viewcount INTEGER NOT NULL,
    file_hash VARCHAR(64) NOT NULL
);