-- name: GetFeedbacks :many
SELECT comment,
    boxes
FROM feedback
ORDER BY id DESC
LIMIT $1;
-- name: NewFeedback :exec
INSERT INTO feedback(comment, boxes)
VALUES($1, $2);