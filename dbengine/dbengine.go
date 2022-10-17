package dbengine

import (
	"context"

	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/romeq/usva/utils"
)

var DbConnection *pgx.Conn

type File struct {
	ID          int
	Filename    string
	Password    string
	IsEncrypted bool
	UploadDate  string `db:"upload_date"`
	OwnerId     int    `db:"owner_id"`
	ViewCount   int    `db:"viewcount"`
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
