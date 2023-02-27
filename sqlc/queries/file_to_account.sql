-- name: FileToAccount :exec
INSERT INTO file_to_account(file_uuid, account_id) 
VALUES($1, get_userid_by_session($2));

-- name: GetSessionOwnerFiles :many
SELECT 
    f.file_uuid AS filename,
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
WHERE account_id = get_userid_by_session($1)
LIMIT $2;

-- name: GetAllSessionOwnerFiles :many
SELECT 
    f.file_uuid
FROM
    file_to_account
    JOIN file AS f
    USING(file_uuid)
WHERE account_id = get_userid_by_session($1);