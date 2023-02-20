-- name: FileToAccount :exec
INSERT INTO file_to_account(file_uuid, account_id) 
VALUES($1, (
    SELECT account_id FROM account_session WHERE session_id = $2
));

-- name: GetSessionOwnerFiles :many
SELECT 
    f.file_uuid,
    f.title,
    f.file_size,
    f.viewcount,
    f.encrypted,
    f.upload_date,
    f.last_seen
FROM
    file_to_account 
    JOIN file AS f
    USING(file_uuid)
WHERE account_id = (
    SELECT account_id 
    FROM account_session
    WHERE session_id = $1
)
LIMIT $2;

-- name: GetAllSessionOwnerFiles :many
SELECT 
    f.file_uuid
FROM
    file_to_account
    JOIN file AS f
    USING(file_uuid)
WHERE account_id = (
    SELECT account_id 
    FROM account_session
    WHERE session_id = $1
);