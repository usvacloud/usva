CREATE TABLE IF NOT EXISTS files(
	id 			 SERIAL 		PRIMARY KEY,
	filename 	 VARCHAR(256) 	NOT NULL UNIQUE,
	password 	 VARCHAR(512),
	is_encrypted BOOLEAN 		NOT NULL DEFAULT FALSE,
	upload_date  VARCHAR(256) 	NOT NULL,
	file_size 	 INTEGER 		NOT NULL,
	viewcount	 INTEGER		NOT NULL,
	owner_id 	 INTEGER
);
