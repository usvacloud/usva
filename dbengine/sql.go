package dbengine

const getPasswordQuery = `
SELECT passwdhash
	FROM files
	WHERE file_uuid = ?;`

const getFileQuery = `
SELECT file_uuid, upload_date, isencrypted, viewcount
	FROM files 
	WHERE file_uuid = ?;`

const incrementFileViewCountQuery = `
UPDATE files
	SET viewcount = viewcount + 1
	WHERE file_uuid = ?;
`

const insertFileQuery = `
INSERT INTO files(file_uuid, title, passwdhash, upload_date, isencrypted, viewcount)
VALUES($1, $2, $3, $4, $5, 0)
`
const deleteFileQuery = `
DELETE FROM files
	WHERE file_uuid = ?;
`
