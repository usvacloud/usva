-- name: NewReport :exec
INSERT INTO report(file_uuid, reason)
VALUES($1, $2);