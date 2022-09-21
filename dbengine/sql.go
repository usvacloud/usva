package dbengine

const getPasswordQuery = `
SELECT password
	FROM files
	WHERE filename = ?;`

const getFileQuery = `
SELECT filename, upload_date, is_encrypted, file_size, viewcount
	FROM files 
	WHERE filename = ?`

const incrementFileViewCountQuery = `
UPDATE files
	SET viewcount = viewcount + 1
	WHERE filename = ?
`

const insertFileQuery = `
INSERT INTO files(
	filename,
	password,
	upload_date,
	file_size,
	is_encrypted,
	viewcount,
	owner_id
) VALUES(?, ?, ?, ?, ?, ?, ?);
`
const deleteFileQuery = `
DELETE FROM files
	WHERE filename = ?;
`
