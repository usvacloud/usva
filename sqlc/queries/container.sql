-- name: AddContainer :one
INSERT INTO container(name, password)
VALUES ($1, $2)
RETURNING container_uuid;

-- name: DeleteContainer :exec
DELETE FROM container
WHERE container_uuid = $1;

-- name: UpdateContainerName :one
UPDATE container
SET name = $1
WHERE container_uuid = $2
RETURNING container_uuid;

-- name: UpdateContainerPassword :one
UPDATE container
SET password = $1
WHERE container_uuid = $2
RETURNING container_uuid;