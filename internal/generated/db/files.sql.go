// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: files.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM file
WHERE file_uuid = $1
`

func (q *Queries) DeleteFile(ctx context.Context, fileUuid string) error {
	_, err := q.db.Exec(ctx, deleteFile, fileUuid)
	return err
}

const getAccessToken = `-- name: GetAccessToken :one
SELECT access_token
FROM file
WHERE file_uuid = $1
`

func (q *Queries) GetAccessToken(ctx context.Context, fileUuid string) (string, error) {
	row := q.db.QueryRow(ctx, getAccessToken, fileUuid)
	var access_token string
	err := row.Scan(&access_token)
	return access_token, err
}

const getEncryptedStatus = `-- name: GetEncryptedStatus :one
SELECT encrypted FROM file
WHERE file_uuid = $1
`

func (q *Queries) GetEncryptedStatus(ctx context.Context, fileUuid string) (bool, error) {
	row := q.db.QueryRow(ctx, getEncryptedStatus, fileUuid)
	var encrypted bool
	err := row.Scan(&encrypted)
	return encrypted, err
}

const getEncryptionIV = `-- name: GetEncryptionIV :one
UPDATE file
SET 
    last_seen = CURRENT_TIMESTAMP,
    viewcount = viewcount + 1
WHERE file_uuid = $1
RETURNING encryption_iv
`

func (q *Queries) GetEncryptionIV(ctx context.Context, fileUuid string) ([]byte, error) {
	row := q.db.QueryRow(ctx, getEncryptionIV, fileUuid)
	var encryption_iv []byte
	err := row.Scan(&encryption_iv)
	return encryption_iv, err
}

const getFileInformation = `-- name: GetFileInformation :one
SELECT file_uuid,
    title,
    upload_date,
    encrypted,
    file_size,
    viewcount
FROM file
WHERE file_uuid = $1
`

type GetFileInformationRow struct {
	FileUuid   string         `json:"file_uuid"`
	Title      sql.NullString `json:"title"`
	UploadDate time.Time      `json:"upload_date"`
	Encrypted  bool           `json:"encrypted"`
	FileSize   sql.NullInt32  `json:"file_size"`
	Viewcount  int32          `json:"viewcount"`
}

func (q *Queries) GetFileInformation(ctx context.Context, fileUuid string) (GetFileInformationRow, error) {
	row := q.db.QueryRow(ctx, getFileInformation, fileUuid)
	var i GetFileInformationRow
	err := row.Scan(
		&i.FileUuid,
		&i.Title,
		&i.UploadDate,
		&i.Encrypted,
		&i.FileSize,
		&i.Viewcount,
	)
	return i, err
}

const getLastSeenAll = `-- name: GetLastSeenAll :many
SELECT file_uuid,
    last_seen
FROM file
`

type GetLastSeenAllRow struct {
	FileUuid string    `json:"file_uuid"`
	LastSeen time.Time `json:"last_seen"`
}

func (q *Queries) GetLastSeenAll(ctx context.Context) ([]GetLastSeenAllRow, error) {
	rows, err := q.db.Query(ctx, getLastSeenAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLastSeenAllRow
	for rows.Next() {
		var i GetLastSeenAllRow
		if err := rows.Scan(&i.FileUuid, &i.LastSeen); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPasswordHash = `-- name: GetPasswordHash :one
SELECT passwdhash
FROM file
WHERE file_uuid = $1
`

func (q *Queries) GetPasswordHash(ctx context.Context, fileUuid string) (sql.NullString, error) {
	row := q.db.QueryRow(ctx, getPasswordHash, fileUuid)
	var passwdhash sql.NullString
	err := row.Scan(&passwdhash)
	return passwdhash, err
}

const newFile = `-- name: NewFile :exec
INSERT INTO file(
    file_uuid,
    title,
    passwdhash,
    access_token,
    encryption_iv,
    file_size,
    viewcount
)
VALUES($1, $2, $3, $4, $5, $6, 0)
`

type NewFileParams struct {
	FileUuid     string         `json:"file_uuid"`
	Title        sql.NullString `json:"title"`
	Passwdhash   sql.NullString `json:"passwdhash"`
	AccessToken  string         `json:"access_token"`
	EncryptionIv []byte         `json:"encryption_iv"`
	FileSize     sql.NullInt32  `json:"file_size"`
}

func (q *Queries) NewFile(ctx context.Context, arg NewFileParams) error {
	_, err := q.db.Exec(ctx, newFile,
		arg.FileUuid,
		arg.Title,
		arg.Passwdhash,
		arg.AccessToken,
		arg.EncryptionIv,
		arg.FileSize,
	)
	return err
}

const updateLastSeen = `-- name: UpdateLastSeen :exec
UPDATE file
SET last_seen = CURRENT_TIMESTAMP
WHERE file_uuid = $1
`

func (q *Queries) UpdateLastSeen(ctx context.Context, fileUuid string) error {
	_, err := q.db.Exec(ctx, updateLastSeen, fileUuid)
	return err
}

const updateViewCount = `-- name: UpdateViewCount :exec
UPDATE file
SET viewcount = viewcount + 1
WHERE file_uuid = $1
`

func (q *Queries) UpdateViewCount(ctx context.Context, fileUuid string) error {
	_, err := q.db.Exec(ctx, updateViewCount, fileUuid)
	return err
}
