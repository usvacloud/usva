package dbengine

// File related queries
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

const updateLastSeenQuery = `
UPDATE files
	SET last_seen = $2
	WHERE file_uuid = $1;`

const insertFileQuery = `
INSERT INTO files(
	file_uuid, 
	title, 
	uploader, 
	passwdhash, 
	upload_date, 
	isencrypted, 
	last_seen,
	viewcount
)
VALUES($1, $2, $3, $4, $5, $6, $7, 0);`

const deleteFileQuery = `
DELETE FROM files
	WHERE file_uuid = $1;`

const lastSeenAllQuery = `
SELECT file_uuid, last_seen FROM files;`

const reportQuery = `
INSERT INTO reports(file_uuid, reason) 
	VALUES($1, $2);`

// Feedback related queries
const getFeedbacksQuery = `
SELECT comment, boxes FROM feedbacks 
    ORDER BY id DESC
    LIMIT $1;`

const addFeedbackQuery = `
INSERT INTO feedbacks(comment, boxes)
	VALUES($1, $2);`
