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

const insertFileQuery = `
INSERT INTO files(file_uuid, title, uploader, passwdhash, upload_date, isencrypted, viewcount)
	VALUES($1, $2, $3, $4, $5, $6, 0)`

const deleteFileQuery = `
DELETE FROM files
	WHERE file_uuid = $1;`

// Feedback related queries
const getFeedbacksQuery = `
SELECT comment, boxes FROM feedbacks;`

const addFeedbackQuery = `
INSERT INTO feedbacks(comment, boxes)
	VALUES($1, $2)`
