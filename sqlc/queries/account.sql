-- name: NewAccount :one
INSERT INTO account(
    username, 
    password
)
VALUES ($1, $2)
RETURNING *;

-- name: GetAccountPasswordHash :one
SELECT 
    password 
FROM 
    account
WHERE 
    username = $1;

-- name: ResetPassword :one
UPDATE 
    account 
SET 
    password = $1
WHERE 
    account_id = $2
RETURNING 
    username;

-- name: DeleteAccount :exec
DELETE FROM 
    account
WHERE 
    account_id = $1;

-- name: NewSession :one
INSERT INTO account_session(
    session_id, 
    account_id
)
VALUES 
    ($1, (
        SELECT account_id 
        FROM account
        WHERE username = $2
    ))
RETURNING 
    session_id;

-- name: GetSessionAccount :one
SELECT 
    a.account_id,
    a.username,
    a.register_date,
    a.last_login,
    a.activity_points
FROM 
    account_session AS ac
JOIN 
    account AS a ON a.account_id = ac.account_id
WHERE ac.session_id = $1;

-- name: GetSessions :many
SELECT 
    session_id, start_date
FROM
    account_session
WHERE account_id = (
    SELECT account_id 
    FROM account_session AS ases
    WHERE ases.session_id = $1
);

-- name: DeleteSession :one
DELETE FROM account_session
WHERE
    account_session.session_id = $2 
    AND account_session.account_id = (
        SELECT account_id
        FROM account_session AS acse
        WHERE acse.session_id = $1
    )
RETURNING *;


-- name: DeleteSessions :many
DELETE FROM account_session
WHERE account_id = (
    SELECT account_id
    FROM account_session AS acse
    WHERE acse.session_id = $1
)
RETURNING *;