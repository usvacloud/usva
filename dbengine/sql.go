package dbengine

const getPasswordQuery = `
SELECT passwdhash
	FROM files
	WHERE file_uuid = $1;`

const getFileInformationQuery = `
SELECT file_uuid, title, upload_date, isencrypted, viewcount
	FROM files 
	WHERE file_uuid = $1;`

const incrementFileViewCountQuery = `
UPDATE files
	SET viewcount = viewcount + 1
	WHERE file_uuid = $1;`

const insertFileQuery = `
INSERT INTO files(file_uuid, title, passwdhash, upload_date, isencrypted, viewcount)
	VALUES($1, $2, $3, $4, $5, 0)
`
const deleteFileQuery = `
DELETE FROM files
	WHERE file_uuid = $1;`
