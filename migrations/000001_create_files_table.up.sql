BEGIN;

CREATE TABLE IF NOT EXISTS files(
	id 			 SERIAL 		PRIMARY KEY,
    title        VARCHAR(256),
	passwdhash 	 VARCHAR(512),
	uploader     VARCHAR(256),
	file_uuid 	 VARCHAR(256) 	NOT NULL UNIQUE,
	isencrypted  BOOLEAN 		NOT NULL DEFAULT FALSE,
	upload_date  VARCHAR(256) 	NOT NULL,
	viewcount	 INTEGER		NOT NULL
);

COMMIT;
