package dbengine

import (
	"context"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/romeq/usva/utils"
)

var DbConnection *pgx.Conn

type File struct {
	IncrementalId int    `database:"id"`
	Title         string `database:"title"`
	PasswordHash  string `database:"passwdhash"`
	FileUUID      string `database:"file_uuid"`
	IsEncrypted   bool   `database:"isencrypted"`
	Uploader      string `database:"uploader"`
	UploadDate    string `database:"upload_date"`
	ViewCount     int    `database:"viewcount"`
}

func Init(port uint16, host, database, user, password string) {
	var err error
	DbConnection, err = pgx.Connect(pgx.ConnConfig{
		Host:      host,
		Port:      port,
		User:      user,
		Password:  password,
		Database:  database,
		TLSConfig: nil,
	})
	utils.Check(err)
	utils.Check(DbConnection.Ping(context.Background()))
}
