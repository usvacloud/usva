// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: feedbacks.sql

package db

import (
	"context"
	"database/sql"
)

const getFeedbacks = `-- name: GetFeedbacks :many
SELECT comment,
    boxes
FROM feedback
ORDER BY id DESC
LIMIT $1
`

type GetFeedbacksRow struct {
	Comment sql.NullString `json:"comment"`
	Boxes   string         `json:"boxes"`
}

func (q *Queries) GetFeedbacks(ctx context.Context, limit int32) ([]GetFeedbacksRow, error) {
	rows, err := q.db.Query(ctx, getFeedbacks, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFeedbacksRow{}
	for rows.Next() {
		var i GetFeedbacksRow
		if err := rows.Scan(&i.Comment, &i.Boxes); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const newFeedback = `-- name: NewFeedback :exec
INSERT INTO feedback(comment, boxes)
VALUES($1, $2)
`

type NewFeedbackParams struct {
	Comment sql.NullString `json:"comment"`
	Boxes   string         `json:"boxes"`
}

func (q *Queries) NewFeedback(ctx context.Context, arg NewFeedbackParams) error {
	_, err := q.db.Exec(ctx, newFeedback, arg.Comment, arg.Boxes)
	return err
}
