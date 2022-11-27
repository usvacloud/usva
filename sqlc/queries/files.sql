-- name: GetPasswordHash :one
SELECT passwdhash
FROM file
WHERE file_uuid = $1;
-- name: FileInformation :one
SELECT file_uuid,
    title,
    upload_date,
    isencrypted,
    viewcount
FROM file
WHERE file_uuid = $1;
-- name: UpdateViewCount :exec
UPDATE file
SET viewcount = viewcount + 1
WHERE file_uuid = $1;
-- name: UpdateLastSeen :exec
UPDATE file
SET last_seen = $2
WHERE file_uuid = $1;
-- name: NewFile :exec
INSERT INTO file(
        file_uuid,
        title,
        uploader,
        passwdhash,
        access_token,
        upload_date,
        isencrypted,
        last_seen,
        viewcount
    )
VALUES($1, $2, $3, $4, $5, $6, $7, $8, 0);
-- name: DeleteFile :exec
DELETE FROM file
WHERE file_uuid = $1;
-- name: GetLastSeenAll :many
SELECT file_uuid,
    last_seen
FROM file;
-- name: GetAccessToken :one
SELECT access_token
FROM file
WHERE file_uuid = $1;