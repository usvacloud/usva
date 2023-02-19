-- name: FileToAccount :exec
INSERT INTO file_to_account(file_uuid, account_id) 
VALUES($1, $2);