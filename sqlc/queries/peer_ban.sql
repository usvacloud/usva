-- name: BanPeer :exec
INSERT INTO
    peer_ban(peer_id)
VALUES
    ($1);

-- name: IsBanned :one
SELECT
    peer_id
FROM
    peer_ban
WHERE
    peer_id = $1;

-- name: RemoveBan :exec
DELETE FROM
    peer_ban
WHERE
    peer_id = $1;