CREATE TABLE IF NOT EXISTS files(
	id 			 SERIAL 		PRIMARY KEY,
    name         VARCHAR(256),
	filename 	 VARCHAR(256) 	NOT NULL UNIQUE,
	password 	 VARCHAR(512),
	is_encrypted BOOLEAN 		NOT NULL DEFAULT FALSE,
	upload_date  VARCHAR(256) 	NOT NULL,
	viewcount	 INTEGER		NOT NULL,
);
