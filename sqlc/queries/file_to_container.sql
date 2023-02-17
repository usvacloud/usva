-- name: AddFileToContainer :one
INSERT INTO file_to_container(
    file_uuid,
    container_uuid
)
VALUES ($1, $2)
RETURNING file_to_container_uuid;

-- name: RemoveFileFromContainer :exec
DELETE FROM file_to_container 
WHERE file_to_container_uuid = $1;