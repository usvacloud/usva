-- name: GetEncryptionIV :one
UPDATE file
SET 
    last_seen = CURRENT_TIMESTAMP,
    viewcount = viewcount + 1
WHERE file_uuid = $1
RETURNING encryption_iv;
-- name: GetPasswordHash :one
SELECT passwdhash
FROM file
WHERE file_uuid = $1;
-- name: GetFileInformation :one
SELECT file_uuid,
    title,
    upload_date,
    encrypted,
    file_size,
    viewcount
FROM file
WHERE file_uuid = $1;
-- name: GetEncryptedStatus :one
SELECT encrypted FROM file
WHERE file_uuid = $1;
-- name: GetLastSeenAll :many
SELECT file_uuid,
    last_seen
FROM file;
-- name: GetAccessToken :one
SELECT access_token
FROM file
WHERE file_uuid = $1;
-- name: GetFilename :one
SELECT file_uuid
FROM file
WHERE access_token = $1;

-- name: UpdateViewCount :exec
UPDATE file
SET viewcount = viewcount + 1
WHERE file_uuid = $1;
-- name: UpdateLastSeen :exec
UPDATE file
SET last_seen = CURRENT_TIMESTAMP
WHERE file_uuid = $1;
-- name: NewFile :exec
INSERT INTO file(
    file_uuid,
    title,
    passwdhash,
    access_token,
    encryption_iv,
    file_size,
    viewcount
)
VALUES($1, $2, $3, $4, $5, $6, 0);
-- name: DeleteFile :exec
DELETE FROM file
WHERE file_uuid = $1;